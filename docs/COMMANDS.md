# 📟 Commands Reference

Complete reference for all available commands in ClosedWheeler AGI.

---

## 🎮 Basic Commands

### `/help`

**Description**: Display help information and available commands

**Usage**:
```
/help
```

**Output**:
```
🦅 ClosedWheeler AGI - Available Commands

Basic:
  /help          - Show this help
  /clear         - Clear conversation
  /quit          - Exit agent

Configuration:
  /config reload - Reload configuration
  /config show   - Show current config
  /model         - Show/switch model
  ...
```

---

### `/clear`

**Description**: Clear the current conversation history

**Usage**:
```
/clear
```

**Effect**:
- Clears all messages from short-term memory
- Resets working memory
- Preserves long-term memory
- Resets session context (triggers context refresh on next message)

**When to use**:
- Starting a new topic
- Conversation getting too long
- Want to reset context

---

### `/quit` or `/exit`

**Description**: Exit the agent gracefully

**Usage**:
```
/quit
```
or
```
/exit
```

**Effect**:
- Saves current memory state
- Closes browser if open
- Shuts down Telegram bot if running
- Exits application

---

## ⚙️ Configuration Commands

### `/config reload`

**Description**: Reload configuration from `.agi/config.json`

**Usage**:
```
/config reload
```

**Use cases**:
- After manually editing config.json
- After changing .env file
- To apply new settings without restart

**Output**:
```
✅ Configuration reloaded successfully
Model: gpt-4o-mini
Temperature: 0.70
Max Tokens: 4096
```

---

### `/config show`

**Description**: Display current configuration

**Usage**:
```
/config show
```

**Output**:
```
Current Configuration:
━━━━━━━━━━━━━━━━━━━━━━

Model Settings:
  Primary Model:   gpt-4o-mini
  Fallback Models: gpt-3.5-turbo
  Temperature:     0.70
  Top-P:           0.90
  Max Tokens:      4096
  Context Window:  128000 tokens

Memory Settings:
  Compression Trigger: 15 messages
  STM Size:            20 items
  Working Memory:      50 items
  Long-term Memory:    100 items

API Settings:
  Base URL:        https://api.openai.com/v1
  API Key:         sk-...***

Browser Settings:
  Headless:        true
  Default Timeout: 60s

Telegram:
  Enabled:         false
```

---

### `/config set <key> <value>`

**Description**: Change a configuration value

**Usage**:
```
/config set temperature 0.8
/config set max_tokens 2048
/config set compression_trigger 20
```

**Available keys**:
- `temperature` (0.0-1.0)
- `top_p` (0.0-1.0)
- `max_tokens` (256-8192)
- `compression_trigger` (5-50)

---

## 🤖 Model Commands

### `/model`

**Description**: Show current model information

**Usage**:
```
/model
```

**Output**:
```
Current Model: gpt-4o-mini
━━━━━━━━━━━━━━━━━━━━━━━━━

Configuration:
  Temperature:     0.70
  Top-P:           0.90
  Max Tokens:      4096
  Context Window:  128000 tokens
  Agent-Ready:     ✅

Fallback Models:
  1. gpt-3.5-turbo

Session Stats:
  Total Tokens:    8,432
  API Calls:       12
  Avg per call:    702 tokens
```

---

### `/model switch <name>`

**Description**: Switch to a different model

**Usage**:
```
/model switch gpt-4o
/model switch claude-sonnet-4
```

**Effect**:
- Switches to specified model
- Loads model-specific parameters
- Resets session context
- Preserves conversation history

**Example**:
```
You: /model switch claude-sonnet-4

Agent:
✅ Switched to model: claude-sonnet-4

Model Configuration:
  Temperature:     0.70
  Top-P:           0.95
  Max Tokens:      4096
  Context Window:  200000 tokens

Session reset. Context will be refreshed on next message.
```

---

### `/model list`

**Description**: List available models from API

**Usage**:
```
/model list
```

**Output**:
```
Available Models:
━━━━━━━━━━━━━━━━

OpenAI:
  1. gpt-4o
  2. gpt-4o-mini ✅ (current)
  3. gpt-4-turbo
  4. gpt-3.5-turbo

Use: /model switch <name>
```

---

## 🧠 Memory Commands

### `/memory stats`

**Description**: Show memory statistics

**Usage**:
```
/memory stats
```

**Output**:
```
Memory Statistics:
━━━━━━━━━━━━━━━━━━

Short-term Memory (STM):
  Items:     12 / 20
  Usage:     60%
  Oldest:    45 minutes ago

Working Memory (WM):
  Items:     8 / 50
  Usage:     16%
  Files:     config.json, agent.go, memory.go

Long-term Memory (LTM):
  Items:     3 / 100
  Usage:     3%
  Summaries: 3 compressed sessions

Context:
  Messages:  14
  Trigger:   15 (compression pending)
  Status:    ● Cached
```

---

### `/memory clear`

**Description**: Clear memory (with confirmation)

**Usage**:
```
/memory clear
```

**Prompts for confirmation**:
```
⚠️  This will clear ALL memory (STM, WM, LTM)
Are you sure? (yes/no): yes

✅ Memory cleared
```

---

### `/memory compress`

**Description**: Manually trigger memory compression

**Usage**:
```
/memory compress
```

**Effect**:
- Compresses current conversation to summary
- Moves summary to long-term memory
- Clears short-term memory
- Resets session context

---

## 🌐 Browser Commands

### `/browser status`

**Description**: Check browser status

**Usage**:
```
/browser status
```

**Output when inactive**:
```
Browser Status: Inactive
No browser session running
```

**Output when active**:
```
Browser Status: Active
━━━━━━━━━━━━━━━━━━━━

Current Page:  https://example.com
Tabs Open:     3
Session Time:  5 minutes
```

---

### `/browser close`

**Description**: Close browser session

**Usage**:
```
/browser close
```

**Effect**:
- Closes all browser tabs
- Shuts down browser instance
- Frees resources

---

## 💬 Telegram Commands

### `/telegram status`

**Description**: Check Telegram bot status

**Usage**:
```
/telegram status
```

**Output**:
```
Telegram Bot: Active
━━━━━━━━━━━━━━━━━━━━

Bot Username:  @YourBot
Authorized:    1 admin
Active Chats:  1
Uptime:        2h 34m
```

---

### `/telegram send <message>`

**Description**: Send message via Telegram (from TUI)

**Usage**:
```
/telegram send Task completed successfully!
```

**Note**: Only works if Telegram integration is enabled

---

## 📊 Session Commands

### `/session stats`

**Description**: Show session statistics

**Usage**:
```
/session stats
```

**Output**:
```
Session Statistics:
━━━━━━━━━━━━━━━━━━━

Session ID:       abc123...
Started:          2h 15m ago
Messages:         24
Context Status:   ● Cached

Token Usage:
  Total Prompt:   18,432 tokens
  Avg per call:   768 tokens
  Savings:        ~65% (cached context)

API Calls:
  Total:          24
  Successful:     24
  Failed:         0
  Fallback Used:  0
```

---

### `/session reset`

**Description**: Reset session (forces context refresh)

**Usage**:
```
/session reset
```

**Effect**:
- Resets session tracking
- Clears context cache hashes
- Forces full context send on next message
- Preserves conversation history

**When to use**:
- After major config changes
- When context seems stale
- Debugging context issues

---

## 🔍 Debug Commands

### `/debug on`

**Description**: Enable debug logging

**Usage**:
```
/debug on
```

**Effect**:
- Shows detailed API requests/responses
- Displays token usage per message
- Shows tool execution details
- Logs to `.agi/logs/debug.log`

---

### `/debug off`

**Description**: Disable debug logging

**Usage**:
```
/debug off
```

---

### `/debug context`

**Description**: Show current context structure

**Usage**:
```
/debug context
```

**Output**:
```
Context Structure:
━━━━━━━━━━━━━━━━━━

System Prompt:
  Hash:    a3f5c2...
  Size:    1,234 tokens
  Cached:  ✅

Rules:
  Hash:    b7d9e1...
  Size:    567 tokens
  Cached:  ✅

Project Info:
  Hash:    c4a8f3...
  Size:    123 tokens
  Cached:  ✅

Messages:
  Count:   14
  Size:    2,100 tokens
  Cached:  ❌ (always fresh)

Total Context: ~4,024 tokens
Next Message:  ~150 tokens (cached context)
```

---

## 🛠️ Tool Commands

### `/tools list`

**Description**: List all available tools

**Usage**:
```
/tools list
```

**Output**:
```
Available Tools:
━━━━━━━━━━━━━━━━

File Operations:
  - ReadFile          Read file contents
  - WriteFile         Write to file
  - DeleteFile        Delete file
  - ListFiles         List directory

Browser:
  - OpenBrowser       Open browser
  - Navigate          Navigate to URL
  - Click             Click element
  - Type              Type text
  - Screenshot        Take screenshot
  - GetText           Extract text
  - MapElements       Map clickable elements
  - CloseBrowser      Close browser

Code Execution:
  - ExecuteCommand    Run shell command
  - RunTests          Run test suite

Git:
  - GitStatus         Check git status
  - GitCommit         Create commit
  - GitPush           Push changes

Total: 18 tools
```

---

## 📝 Keyboard Shortcuts

### In TUI Mode

**Ctrl+D**: Send message
**Ctrl+C**: Cancel input or exit
**Ctrl+L**: Clear screen
**Up/Down**: Navigate message history

---

## 🎯 Special Commands

### `/setup`

**Description**: Re-run interactive setup wizard

**Usage**:
```
/setup
```

**Effect**:
- Starts setup wizard
- Allows reconfiguring all settings
- Can re-interview models
- Overwrites existing config

---

### `/version`

**Description**: Show version information

**Usage**:
```
/version
```

**Output**:
```
🦅 ClosedWheeler AGI
Version:        2.0
Build:          2026-02-08
Go Version:     1.21.0
Platform:       windows/amd64
Binary Size:    13 MB
```

---

## 🆘 Help Commands

### `/help <topic>`

**Description**: Get help on specific topic

**Usage**:
```
/help config
/help memory
/help browser
```

**Available topics**:
- config
- memory
- browser
- telegram
- models
- tools
- shortcuts

---

## 💡 Usage Tips

### Combining Commands

```bash
# Reset and configure
/session reset
/config set temperature 0.8
/config reload

# Debug workflow
/debug on
[ask question]
/debug context
/debug off
```

### Command Aliases

Some commands have shorter aliases:
- `/q` → `/quit`
- `/h` → `/help`
- `/c` → `/clear`
- `/m` → `/model`
- `/s` → `/session stats`

---

## 🔐 Admin Commands (Telegram)

These commands work only in Telegram if you're authorized:

### `/admin_status`

**Description**: Get full system status

**Usage**: Send in Telegram chat
```
/admin_status
```

---

### `/admin_approve <id>`

**Description**: Approve pending permission request

**Usage**: Reply to permission request
```
/admin_approve 12345
```

---

### `/admin_deny <id>`

**Description**: Deny permission request

**Usage**: Reply to permission request
```
/admin_deny 12345
```

---

## 📚 Command Categories

### Essential (Use Daily)
- `/help`
- `/clear`
- `/model`
- `/config reload`

### Configuration (Setup & Tuning)
- `/config show`
- `/config set`
- `/session reset`

### Memory (Optimization)
- `/memory stats`
- `/memory compress`

### Debugging (When Issues Occur)
- `/debug on/off`
- `/debug context`
- `/session stats`

### Integration (Telegram, Browser)
- `/telegram status`
- `/browser status`

---

**Need more help?** Check [Common Tasks](COMMON_TASKS.md) for usage examples! 🚀

*Complete command reference for ClosedWheeler AGI!* 📟
