# 🚀 Getting Started with ClosedWheeler AGI

Quick guide to get you up and running in minutes.

---

## Step 1: First Run

```bash
.\ClosedWheeler.exe
```

The interactive setup wizard will start automatically.

---

## Step 2: Setup Wizard

### 2.1 Agent Name
```
Give your agent a name [ClosedWheeler]: MyAssistant
```

Give your agent a friendly name.

### 2.2 API Configuration
```
API Base URL [https://api.openai.com/v1]:
API Key: sk-...
```

**Popular APIs**:
- OpenAI: `https://api.openai.com/v1`
- Anthropic: `https://api.anthropic.com/v1`
- Local (Ollama): `http://localhost:11434/v1`

### 2.3 Model Selection
```
🔍 Fetching available models...
✅ Found 15 models

Select primary model (1-15 or name): gpt-4o-mini
```

Choose your preferred model.

### 2.4 Model Self-Configuration ⭐ NEW!
```
🎤 Model Self-Configuration
Asking 'gpt-4o-mini' to configure itself...

✅ Model configured!
  Context Window:  128000 tokens
  Temperature:     0.70
  Top-P:           0.90
  Max Tokens:      4096
  Agent-Ready:     true
```

The model interviews itself and sets optimal parameters automatically!

### 2.5 Fallback Models
```
Add fallback models? (y/N): y
Select fallback model: gpt-3.5-turbo
Add another? (y/N): n
```

Optional: Add backup models for reliability.

### 2.6 Permissions
```
Choose preset:
1. Full Access
2. Restricted (read/write files only)
3. Read-Only

Select: 1
```

Choose your permission level.

### 2.7 Rules Preset
```
Choose rules preset:
1. None
2. Code Quality
3. Security First
4. Performance
5. Personal Assistant
6. Cybersecurity
7. Data Science
8. DevOps

Select: 2
```

Choose behavioral rules for your agent.

### 2.8 Memory Configuration
```
Choose memory preset:
1. Minimal
2. Balanced (recommended)
3. Extended

Select: 2
```

Configure memory tiers.

### 2.9 Telegram (Optional)
```
Enable Telegram integration? (y/N): n
```

Skip for now, can enable later.

---

## Step 3: First Conversation

After setup completes:

```
🦅 ClosedWheelerAGI v2.0

Type your message (Ctrl+D to send, /help for commands):
> Hello! Can you help me understand quicksort?

Agent: [Provides detailed explanation...]

> Thanks! Now show me Python code for it.

Agent: [Provides implementation...]
```

---

## Step 4: Basic Commands

```bash
/help          # Show all commands
/clear         # Clear conversation
/model         # Show/switch model
/config reload # Reload configuration
/quit          # Exit agent
```

---

## Step 5: Understanding the TUI

### Status Bar (Top)
```
[IDLE] 🦅 ClosedWheelerAGI v2.0  Tokens: 8000 (8 calls)

       ● STM: 6 │ WM: 8 │ LTM: 3 │ CTX: 10 msgs
```

**Indicators**:
- `●` (green) = Context cached (saving tokens!)
- `○` (orange) = Context refreshing
- `CTX: N msgs` = Messages in context
- `⚠️` (orange) = Warning (large context)

### Memory Stats
- **STM**: Short-term memory (recent messages)
- **WM**: Working memory (active files)
- **LTM**: Long-term memory (compressed)

---

## Next Steps

### Learn More
- **[Common Tasks](COMMON_TASKS.md)** - Usage examples
- **[Commands](COMMANDS.md)** - Full command reference
- **[Context Optimization](CONTEXT_OPTIMIZATION.md)** - How caching works

### Configuration
Your config is saved in:
- `.agi/config.json` - Main configuration
- `workplace/.agirules` - Agent rules
- `.env` - API credentials

### Customization
Edit `workplace/.agirules` to customize agent behavior.

---

## Troubleshooting

### Setup Fails
- **Check API key**: Make sure it's valid
- **Check internet**: Connection required
- **Try again**: Setup is repeatable

### Model Interview Fails
```
⚠️  Model self-configuration failed
Use model anyway with fallback config? (Y/n): y
```

Choose Y to use known profiles as fallback.

### Agent Not Responding
1. Check logs: `.agi/logs/latest.log`
2. Verify config: `.agi/config.json`
3. Reload: `/config reload`

---

## Tips

✅ **Let model configure itself** - Most accurate
✅ **Use fallback models** - Prevents timeouts
✅ **Watch context indicator** - ● = efficient!
✅ **Let compression happen** - Automatic
✅ **Read logs** if issues - Very helpful

---

**Ready to use!** Check [Common Tasks](COMMON_TASKS.md) for examples.

*You're all set!* 🚀
