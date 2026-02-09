// Package context provides project context management.
package context

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ProjectContext holds all information about the project
type ProjectContext struct {
	RootPath     string
	Files        map[string]*FileInfo
	Dependencies map[string][]string
	Metrics      *Metrics
	mu           sync.RWMutex
}

// FileInfo holds information about a single file
type FileInfo struct {
	Path      string
	RelPath   string
	Content   string
	Language  string
	Size      int64
	LineCount int
	Functions []FunctionInfo
	Imports   []string
}

// FunctionInfo holds information about a function
type FunctionInfo struct {
	Name       string
	StartLine  int
	EndLine    int
	Complexity int
	Signature  string
}

// Metrics holds project-wide metrics
type Metrics struct {
	TotalFiles     int
	TotalLines     int
	TotalFunctions int
	Languages      map[string]int
}

// NewProjectContext creates a new project context
func NewProjectContext(rootPath string) *ProjectContext {
	absPath, _ := filepath.Abs(rootPath)
	return &ProjectContext{
		RootPath:     absPath,
		Files:        make(map[string]*FileInfo),
		Dependencies: make(map[string][]string),
		Metrics:      &Metrics{Languages: make(map[string]int)},
	}
}

// Load loads all files in the project
func (pc *ProjectContext) Load(ignorePatterns []string) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	// Reset
	pc.Files = make(map[string]*FileInfo)
	pc.Metrics = &Metrics{Languages: make(map[string]int)}

	err := filepath.Walk(pc.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Get relative path
		relPath, _ := filepath.Rel(pc.RootPath, path)

		// Check ignore patterns
		for _, pattern := range ignorePatterns {
			if matchIgnorePattern(relPath, pattern) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		if info.IsDir() {
			return nil
		}

		// Detect language and skip binary files
		lang := detectLanguage(path)
		if lang == "binary" {
			return nil
		}

		// Create file info
		fi := &FileInfo{
			Path:     path,
			RelPath:  relPath,
			Language: lang,
			Size:     info.Size(),
		}

		// Load content for smaller files
		if info.Size() < 1024*1024 { // 1MB limit
			content, err := os.ReadFile(path)
			if err == nil {
				fi.Content = string(content)
				fi.LineCount = strings.Count(fi.Content, "\n") + 1
			}
		}

		// Analyze Go files
		if lang == "go" {
			fi.analyzeGo()
		}

		pc.Files[path] = fi
		pc.Metrics.TotalFiles++
		pc.Metrics.TotalLines += fi.LineCount
		pc.Metrics.TotalFunctions += len(fi.Functions)
		pc.Metrics.Languages[lang]++

		// Track dependencies
		if len(fi.Imports) > 0 {
			pc.Dependencies[path] = fi.Imports
		}

		return nil
	})

	return err
}

// GetFile returns file info for a path
func (pc *ProjectContext) GetFile(path string) (*FileInfo, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	// Try absolute path first
	if fi, ok := pc.Files[path]; ok {
		return fi, true
	}

	// Try relative path
	absPath := filepath.Join(pc.RootPath, path)
	fi, ok := pc.Files[absPath]
	return fi, ok
}

// GetFileContent returns the content of a file
func (pc *ProjectContext) GetFileContent(path string) (string, error) {
	fi, ok := pc.GetFile(path)
	if !ok {
		return "", fmt.Errorf("file not found: %s", path)
	}
	return fi.Content, nil
}

// GetFilesByLanguage returns files of a specific language
func (pc *ProjectContext) GetFilesByLanguage(lang string) []*FileInfo {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	files := make([]*FileInfo, 0)
	for _, fi := range pc.Files {
		if fi.Language == lang {
			files = append(files, fi)
		}
	}
	return files
}

// GetSummary returns a text summary of the project
func (pc *ProjectContext) GetSummary() string {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Project: %s\n", filepath.Base(pc.RootPath)))
	sb.WriteString(fmt.Sprintf("Files: %d\n", pc.Metrics.TotalFiles))
	sb.WriteString(fmt.Sprintf("Lines: %d\n", pc.Metrics.TotalLines))
	sb.WriteString(fmt.Sprintf("Functions: %d\n", pc.Metrics.TotalFunctions))
	sb.WriteString("\nLanguages:\n")
	for lang, count := range pc.Metrics.Languages {
		sb.WriteString(fmt.Sprintf("  %s: %d files\n", lang, count))
	}
	return sb.String()
}

// GetFileList returns a list of all files
func (pc *ProjectContext) GetFileList() []string {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	files := make([]string, 0, len(pc.Files))
	for _, fi := range pc.Files {
		files = append(files, fi.RelPath)
	}
	return files
}

// analyzeGo performs basic Go file analysis
func (fi *FileInfo) analyzeGo() {
	lines := strings.Split(fi.Content, "\n")

	// Find imports
	inImport := false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "import (") {
			inImport = true
			continue
		}
		if inImport {
			if line == ")" {
				inImport = false
				continue
			}
			if line != "" && !strings.HasPrefix(line, "//") {
				pkg := strings.Trim(line, "\"")
				fi.Imports = append(fi.Imports, pkg)
			}
		}
		if strings.HasPrefix(line, "import \"") {
			pkg := strings.TrimPrefix(line, "import \"")
			pkg = strings.TrimSuffix(pkg, "\"")
			fi.Imports = append(fi.Imports, pkg)
		}
	}

	// Find functions
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "func ") {
			fn := FunctionInfo{
				StartLine: i + 1,
				Signature: line,
			}

			// Extract name
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				fnName := parts[1]
				if idx := strings.Index(fnName, "("); idx != -1 {
					fn.Name = fnName[:idx]
				}
			}

			// Find end of function
			braceCount := 0
			for j := i; j < len(lines); j++ {
				l := lines[j]
				braceCount += strings.Count(l, "{") - strings.Count(l, "}")

				// Count complexity
				trimmed := strings.TrimSpace(l)
				if strings.HasPrefix(trimmed, "if ") ||
					strings.HasPrefix(trimmed, "for ") ||
					strings.HasPrefix(trimmed, "switch ") ||
					strings.HasPrefix(trimmed, "case ") {
					fn.Complexity++
				}

				if braceCount == 0 && j > i {
					fn.EndLine = j + 1
					break
				}
			}

			fn.Complexity++ // Base complexity
			fi.Functions = append(fi.Functions, fn)
		}
	}
}

// Helper functions

func matchIgnorePattern(path, pattern string) bool {
	// Normalize path to use forward slashes
	path = filepath.ToSlash(path)
	pattern = filepath.ToSlash(strings.TrimSuffix(pattern, "/"))

	parts := strings.Split(path, "/")
	for _, part := range parts {
		if part == "" {
			continue
		}
		// Exact match for a component
		if part == pattern {
			return true
		}
		// Glob match for a component (e.g., *.log)
		if matched, _ := filepath.Match(pattern, part); matched {
			return true
		}
	}

	return false
}

func detectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	langMap := map[string]string{
		".go":    "go",
		".py":    "python",
		".js":    "javascript",
		".ts":    "typescript",
		".jsx":   "javascript",
		".tsx":   "typescript",
		".java":  "java",
		".c":     "c",
		".cpp":   "cpp",
		".h":     "c",
		".hpp":   "cpp",
		".rs":    "rust",
		".rb":    "ruby",
		".php":   "php",
		".cs":    "csharp",
		".sh":    "shell",
		".bash":  "shell",
		".yaml":  "yaml",
		".yml":   "yaml",
		".json":  "json",
		".xml":   "xml",
		".html":  "html",
		".css":   "css",
		".md":    "markdown",
		".txt":   "text",
		".sql":   "sql",
		".proto": "protobuf",
	}

	if lang, ok := langMap[ext]; ok {
		return lang
	}

	return "binary"
}
