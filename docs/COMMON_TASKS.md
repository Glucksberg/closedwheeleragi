# 📋 Common Tasks

Quick examples of typical usage patterns for ClosedWheeler AGI.

---

## 💬 Conversational Tasks

### Ask Questions

```
You: "Explain how binary search works"

Agent: [Provides clear explanation with examples and time complexity analysis]
```

### Get Recommendations

```
You: "What's the best way to handle API rate limiting in Go?"

Agent: [Suggests token bucket algorithm, provides implementation example]
```

### Understand Code

```
You: "Explain what this function does"
[Paste code]

Agent: [Analyzes code, explains purpose, logic flow, and potential issues]
```

---

## 💻 Code Tasks

### Write New Code

```
You: "Write a function to validate email addresses using regex"

Agent: [Provides implementation with explanation]
```

**Result**:
```go
func isValidEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched
}
```

### Refactor Code

```
You: "Refactor this to use channels instead of callbacks"
[Paste callback-based code]

Agent: [Analyzes code, rewrites using Go channels with explanation]
```

### Debug Code

```
You: "Why is this code causing a deadlock?"
[Paste problematic code]

Agent: [Identifies issue, explains why deadlock occurs, provides fix]
```

### Optimize Performance

```
You: "How can I make this faster?"
[Paste code]

Agent: [Profiles code, suggests optimizations with benchmarks]
```

---

## 📁 File Operations

### Read Files

```
You: "Read config.json and tell me what the current model is"

Agent: [Uses ReadFile tool, reports model name]
```

### Edit Files

```
You: "Change the timeout in config.json from 30 to 60 seconds"

Agent: [Reads file, updates timeout, confirms change]
```

### Create Files

```
You: "Create a .gitignore file for a Go project"

Agent: [Creates .gitignore with common Go patterns]
```

### Search Files

```
You: "Find all files that use the database connection"

Agent: [Searches codebase, lists files with line numbers]
```

---

## 🌐 Browser Research

### Web Search

```
You: "Research the latest features in Go 1.22"

Agent:
  1. Opens browser
  2. Searches for "Go 1.22 features"
  3. Navigates to official release notes
  4. Summarizes key features
  5. Closes browser
```

### Extract Information

```
You: "Get the latest version of React from npmjs.com"

Agent:
  1. Opens browser
  2. Navigates to npmjs.com/package/react
  3. Extracts version number
  4. Reports: "React latest version: 18.2.0"
```

### Fill Forms

```
You: "Fill out this contact form with test data"
[Provides form URL]

Agent:
  1. Opens browser
  2. Maps form elements
  3. Fills fields with test data
  4. Takes screenshot for confirmation
```

---

## 🔧 Project Tasks

### Setup New Project

```
You: "Initialize a new Go module called 'myapi'"

Agent:
  1. Creates directory structure
  2. Runs `go mod init myapi`
  3. Creates main.go with basic structure
  4. Creates .gitignore
  5. Confirms setup complete
```

### Add Dependencies

```
You: "Add the Gin web framework to this project"

Agent: [Runs `go get -u github.com/gin-gonic/gin`, confirms installation]
```

### Run Tests

```
You: "Run all tests in the project"

Agent: [Executes `go test ./...`, reports results]
```

---

## 📊 Data Analysis

### Parse JSON

```
You: "Parse this JSON and tell me how many users are active"
[Paste JSON data]

Agent: [Analyzes JSON, counts active users, provides summary]
```

### Generate Reports

```
You: "Analyze these log files and summarize errors"

Agent:
  1. Reads log files
  2. Categorizes errors
  3. Counts occurrences
  4. Provides formatted report
```

### Compare Data

```
You: "Compare old_config.json and new_config.json, show differences"

Agent: [Reads both files, highlights changes, explains impact]
```

---

## 🛠️ Configuration Tasks

### Check Configuration

```
You: "Show me my current configuration"

Agent: [Reads .agi/config.json, displays formatted summary]
```

**Output**:
```
Current Configuration:
  Model:           gpt-4o-mini
  Temperature:     0.70
  Top-P:           0.90
  Max Tokens:      4096
  Context Window:  128000
  Fallback Models: gpt-3.5-turbo
  Memory:
    - Compression:  15 messages
    - STM:          20 items
    - WM:           50 items
    - LTM:          100 items
```

### Modify Configuration

```
You: "Change compression trigger to 20"

Agent: [Updates config.json, reloads, confirms change]
```

### Reload Configuration

```
You: "/config reload"

Agent: ✅ Configuration reloaded successfully
```

---

## 🎯 Advanced Tasks

### Multi-Step Workflows

```
You: "Research Go logging libraries, compare the top 3, and recommend one"

Agent:
  1. Opens browser
  2. Searches for "best Go logging libraries 2026"
  3. Identifies: zap, logrus, zerolog
  4. Compares features, performance, popularity
  5. Provides detailed recommendation with reasoning
```

### Complex Refactoring

```
You: "Convert this synchronous API to use async/await pattern"
[Paste large codebase section]

Agent:
  1. Analyzes current structure
  2. Identifies blocking operations
  3. Refactors to async pattern
  4. Updates error handling
  5. Provides migration guide
```

### Documentation Generation

```
You: "Generate API documentation for all exported functions in this package"

Agent:
  1. Scans package files
  2. Extracts function signatures
  3. Generates markdown documentation
  4. Includes examples and parameter descriptions
```

---

## 🔄 Git Operations

### Check Status

```
You: "What's my git status?"

Agent: [Runs `git status`, reports untracked files and changes]
```

### Commit Changes

```
You: "Commit these changes with a meaningful message"

Agent:
  1. Reviews changes
  2. Stages relevant files
  3. Creates descriptive commit message
  4. Commits with co-author attribution
```

### Branch Management

```
You: "Create a new feature branch for authentication"

Agent: [Creates `feature/authentication` branch, switches to it]
```

---

## 📝 Documentation Tasks

### Write README

```
You: "Create a README for this project"

Agent:
  1. Analyzes project structure
  2. Identifies key features
  3. Creates comprehensive README with:
     - Description
     - Installation
     - Usage examples
     - Configuration
     - License
```

### Generate Changelog

```
You: "Generate a changelog from git commits since last release"

Agent: [Parses commits, categorizes changes, formats as CHANGELOG.md]
```

---

## 🧪 Testing Tasks

### Write Tests

```
You: "Write unit tests for this function"
[Paste function]

Agent:
  1. Analyzes function logic
  2. Identifies edge cases
  3. Creates comprehensive test suite
  4. Includes table-driven tests
```

### Debug Test Failures

```
You: "Why is TestUserLogin failing?"

Agent:
  1. Runs test
  2. Analyzes error output
  3. Identifies root cause
  4. Suggests fix
```

---

## 💡 Tips for Effective Usage

### Be Specific

❌ **Vague**: "Fix this"
✅ **Specific**: "Fix the null pointer error on line 42 in user.go"

### Provide Context

❌ **No context**: "Update the API"
✅ **With context**: "Update the /users API endpoint to return email field"

### Ask Follow-ups

```
You: "Implement user authentication"
Agent: [Provides JWT-based implementation]

You: "Now add refresh token support"
Agent: [Extends implementation with refresh tokens]
```

### Use Commands

```
/help          - When you need guidance
/model         - To switch models for specific tasks
/clear         - To start fresh conversation
/config reload - After changing configuration
```

---

## 🎯 Task Patterns

### Research → Implement

```
1. You: "Research best practices for Go error handling"
2. Agent: [Provides comprehensive summary]
3. You: "Now refactor my error handling to follow those practices"
4. Agent: [Applies best practices to your code]
```

### Analyze → Fix → Test

```
1. You: "Analyze why this function is slow"
2. Agent: [Profiles code, identifies bottleneck]
3. You: "Fix the performance issue"
4. Agent: [Optimizes code]
5. You: "Write benchmarks to verify improvement"
6. Agent: [Creates benchmarks with before/after comparison]
```

### Plan → Build → Document

```
1. You: "Design a caching layer for this API"
2. Agent: [Provides architectural design]
3. You: "Implement it"
4. Agent: [Builds caching layer]
5. You: "Document how to use it"
6. Agent: [Creates usage documentation]
```

---

## 🚀 Performance Notes

### Context Optimization

When you see **● (green)** in the status bar:
- Context is cached
- Faster responses (~1-2s)
- Lower token usage (~800 tokens/msg)

When you see **○ (orange)**:
- Context refreshing
- Slightly slower (~3s)
- Higher token usage (~2500 tokens)

### Compression

When context grows past trigger (default: 15 messages):
- Agent automatically compresses older messages
- Long-term memory preserves key information
- Session resets for fresh context

---

**Ready to build?** Start with simple tasks and progress to complex workflows! 🚀

*Need more help? Check [Commands Reference](COMMANDS.md) for full command list.*
