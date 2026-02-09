// Package ignore provides ignore pattern management.
package ignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const DefaultIgnoreFile = ".agiignore"

// Patterns holds ignore patterns
type Patterns struct {
	patterns []string
	filePath string
}

// Load loads ignore patterns from a file
func Load(projectPath string) *Patterns {
	p := &Patterns{
		patterns: make([]string, 0),
		filePath: filepath.Join(projectPath, DefaultIgnoreFile),
	}

	p.load()
	return p
}

// LoadFromFile loads patterns from a specific file
func LoadFromFile(filePath string) *Patterns {
	p := &Patterns{
		patterns: make([]string, 0),
		filePath: filePath,
	}

	p.load()
	return p
}

func (p *Patterns) load() {
	file, err := os.Open(p.filePath)
	if err != nil {
		// No ignore file, use minimal defaults
		p.patterns = []string{".git/", ".agi/"}
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		p.patterns = append(p.patterns, line)
	}
}

// ShouldIgnore checks if a path should be ignored
func (p *Patterns) ShouldIgnore(relPath string) bool {
	// Normalize path separators
	relPath = filepath.ToSlash(relPath)

	for _, pattern := range p.patterns {
		pattern = filepath.ToSlash(pattern)

		// Directory pattern (ends with /)
		if strings.HasSuffix(pattern, "/") {
			dir := strings.TrimSuffix(pattern, "/")
			if strings.HasPrefix(relPath, dir+"/") || relPath == dir {
				return true
			}
			// Check if any path component matches
			parts := strings.Split(relPath, "/")
			for _, part := range parts {
				if part == dir {
					return true
				}
			}
			continue
		}

		// Glob pattern
		if strings.ContainsAny(pattern, "*?") {
			matched, _ := filepath.Match(pattern, filepath.Base(relPath))
			if matched {
				return true
			}
			continue
		}

		// Exact match or contains
		if relPath == pattern || strings.Contains(relPath, pattern) {
			return true
		}
	}

	return false
}

// List returns all patterns
func (p *Patterns) List() []string {
	return p.patterns
}

// Add adds a pattern
func (p *Patterns) Add(pattern string) {
	p.patterns = append(p.patterns, pattern)
}

// Save saves patterns to the file
func (p *Patterns) Save() error {
	var sb strings.Builder

	sb.WriteString("# Coder AGI Ignore File\n")
	sb.WriteString("# Files and directories matching these patterns will be ignored\n\n")

	for _, pattern := range p.patterns {
		sb.WriteString(pattern + "\n")
	}

	return os.WriteFile(p.filePath, []byte(sb.String()), 0644)
}

// CreateDefault creates a default .agiignore file if it doesn't exist
func CreateDefault(projectPath string) error {
	filePath := filepath.Join(projectPath, DefaultIgnoreFile)

	if _, err := os.Stat(filePath); err == nil {
		return nil // Already exists
	}

	content := `# Coder AGI Ignore File
# Files and directories matching these patterns will be ignored
# Supports glob patterns

# Version control
.git/

# Dependencies
vendor/
node_modules/

# Coder AGI data
.agi/

# Build outputs
*.exe
*.dll
*.so
*.dylib
bin/
dist/

# Logs and temp files
*.log
*.tmp

# IDE/Editor
.idea/
.vscode/
`
	return os.WriteFile(filePath, []byte(content), 0644)
}
