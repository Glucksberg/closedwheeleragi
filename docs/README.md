# 📚 ClosedWheeler AGI - Documentation

Complete technical documentation for ClosedWheeler AGI v2.0

---

## 🚀 Quick Start Guides

Perfect for first-time users and common tasks.

| Guide | Description |
|-------|-------------|
| **[Getting Started](GETTING_STARTED.md)** | First-time setup walkthrough |
| **[Common Tasks](COMMON_TASKS.md)** | Examples of typical usage |
| **[Commands Reference](COMMANDS.md)** | All available commands |
| **[Test Guide](../TESTE-ISTO-AGORA.md)** | How to test new features |
| **[Quick Start](../LEIA-ME-PRIMEIRO.md)** | Portuguese quick start |

---

## 🏗️ Core Features

Deep dives into major system components.

### Context Optimization
- **[Context Optimization Guide](CONTEXT_OPTIMIZATION.md)** - Full technical guide (500+ lines)
  - Session-based caching
  - 60-80% token savings
  - Smart change detection
  - Auto-compression

### Model Configuration
- **[Model Interview System](MODEL_INTERVIEW.md)** - Self-configuration (500+ lines)
  - How models configure themselves
  - Interview process explained
  - JSON format and validation
  - Example responses

- **[Model Parameters Guide](MODEL_PARAMETERS.md)** - Parameter detection (500+ lines)
  - Auto-detection system
  - Known model profiles
  - Best practices

### Browser Automation
- **[Browser Complete Guide](BROWSER_COMPLETE_GUIDE.md)** - All 9 browser tools
  - Element mapping
  - Coordinate clicking
  - AI-optimized screenshots
  - Complete workflow examples

- **[Browser Navigation](BROWSER_NAVIGATION_GUIDE.md)** - Navigation basics
  - Quick reference
  - Common patterns

### Fallback System
- **[Fallback Models Guide](FALLBACK_MODELS_GUIDE.md)** - Automatic failover
  - Zero context loss
  - Timeout configuration
  - How it works

### Communication
- **[Telegram Integration](TELEGRAM_CHAT_GUIDE.md)** - Remote control
  - Setup instructions
  - Available commands
  - Approval workflow

- **[Telegram Quick Setup](TELEGRAM_SETUP_QUICK.md)** - Fast setup

### Setup & Configuration
- **[Enhanced Setup Guide](ENHANCED_SETUP_GUIDE.md)** - Complete setup wizard
  - All 8 steps explained
  - Presets available
  - Configuration options

- **[Fallback Quick Start](QUICK_START_FALLBACK.md)** - Fallback setup

---

## 📖 Implementation Details

For developers and advanced users who want to understand internals.

### Implementation Summaries
| Document | Topic |
|----------|-------|
| **[Final Implementation](FINAL_IMPLEMENTATION.md)** | Complete system overview |
| **[Implementation Complete](IMPLEMENTATION_COMPLETE.md)** | Feature completion status |
| **[Fallback Implementation](FALLBACK_IMPLEMENTATION.md)** | Fallback system internals |

---

## 📝 Release Notes & Changelogs

Track what changed and when.

| Date | Document | Highlights |
|------|----------|------------|
| **2026-02-08** | **[Release Notes v2](RELEASE_NOTES_v2.md)** | Version 2.0 release |

---

## 🏗️ Architecture

Understanding how the system works.

### System Components

```
ClosedWheelerAGI/
│
├── pkg/agent/              # Core agent with session manager
├── pkg/llm/                # LLM client + model profiles + interview
├── pkg/memory/             # Tiered memory (STM/WM/LTM)
├── pkg/browser/            # Playwright integration
├── pkg/telegram/           # Telegram bot
├── pkg/permissions/        # RBAC and auditing
├── pkg/tools/              # Tool registry
├── pkg/prompts/            # Rules manager
└── pkg/tui/                # Terminal UI
```

### Data Flow

```
User Input
    │
    ▼
Agent (pkg/agent/)
    │
    ├─→ Session Manager ────→ Context Caching
    │                          (60-80% token savings)
    │
    ├─→ Memory Manager ─────→ STM/WM/LTM
    │                          (Auto-compression)
    │
    ├─→ LLM Client ─────────→ API Calls
    │                          (With fallback)
    │
    ├─→ Tools Registry ─────→ Execute Actions
    │   │
    │   ├─→ Browser Tools
    │   ├─→ File Operations
    │   └─→ Git Commands
    │
    └─→ Response ───────────→ User
```

---

## 🎯 Key Concepts

### Context Optimization

**Problem**: Sending full context (system prompt + rules + project info) every message wastes tokens.

**Solution**: Session-based caching
1. First message: Send full context
2. Store hashes of components
3. Next messages: Only conversation
4. Refresh only when components change

**Result**: 60-80% token savings, 2-3x faster

### Model Self-Configuration

**Problem**: Different models support different parameters.

**Solution**: Ask the model to configure itself!
1. Send interview question
2. Model responds with JSON config
3. System validates and saves

**Result**: Zero manual configuration, always optimal

### Memory Tiers

**Problem**: Can't keep everything in context forever.

**Solution**: Three-tier system
- **STM** (Short-term): Recent messages (20)
- **WM** (Working): Active files (50)
- **LTM** (Long-term): Compressed summaries (100)

**Result**: Efficient context management

---

## 🎓 Best Practices

### Configuration

1. **Run model interview** during setup
2. **Use recommended parameters** (temp=0.7, top_p=0.9)
3. **Set compression trigger** to 15 for balance
4. **Configure fallback models** for reliability

### Usage

1. **Check context indicator** in TUI (● = cached, ○ = refresh)
2. **Monitor message count** (CTX: N msgs)
3. **Let compression happen** automatically
4. **Reload config** after changes (/config reload)

### Performance

1. **Keep rules concise** - Large rules increase context
2. **Use fallbacks** - Prevents timeout failures
3. **Tune max_tokens** - Balance detail vs speed
4. **Monitor token usage** - Watch for excessive use

---

## 🐛 Troubleshooting

### Common Issues

| Problem | Solution | Doc Reference |
|---------|----------|---------------|
| High token usage | Check compression trigger | [Context Optimization](CONTEXT_OPTIMIZATION.md) |
| Browser timeout | Increase timeout in config | [Browser Guide](BROWSER_COMPLETE_GUIDE.md) |
| Model errors | Check if parameters supported | [Model Interview](MODEL_INTERVIEW.md) |
| Slow responses | Enable fallback models | [Fallback Guide](FALLBACK_MODELS_GUIDE.md) |
| Context too large | Lower compression trigger | [Context Optimization](CONTEXT_OPTIMIZATION.md) |

### Debug Steps

1. **Check logs**: `.agi/logs/latest.log`
2. **Verify config**: `.agi/config.json`
3. **Test model**: Run interview during setup
4. **Monitor TUI**: Watch context indicators

---

## 📊 Performance Metrics

### Token Savings

```
Before Optimization:
Message 1: 2100 tokens (full context)
Message 2: 2150 tokens (full context)
Message 10: 4000 tokens (full context)
Total: ~28,000 tokens

After Optimization:
Message 1: 2100 tokens (full context)
Message 2: 150 tokens (messages only!)
Message 10: 1000 tokens (messages only!)
Total: ~8,000 tokens

Savings: 20,000 tokens (71%)
```

### Response Time

```
Before: 3-5 seconds (large prompts)
After:  1-2 seconds (cached context)
Improvement: 2-3x faster
```

### Session Length

```
Before: 40 messages (context limit)
After:  120 messages (compression)
Improvement: 3x more messages
```

---

## 📦 Document Index

### By Category

**Setup & Getting Started**:
- Getting Started Guide
- Enhanced Setup Guide
- Quick Start (PT-BR)
- Test Guide

**Core Features**:
- Context Optimization
- Model Interview
- Model Parameters
- Browser Automation
- Fallback Models
- Telegram Integration

**Technical**:
- Architecture
- Implementation Details
- System Components
- Data Flow

**Reference**:
- Commands
- API Reference
- Configuration Options
- Troubleshooting

---

## 🔗 Quick Links

- **[Main README](../README.md)** - Project overview
- **[Changelog](../CHANGELOG.md)** - Version history
- **[License](../LICENSE)** - MIT License
- **[Source Code](../pkg/)** - Implementation

---

## 📞 Need Help?

1. **Check relevant guide** above
2. **See [Troubleshooting](#troubleshooting)** section
3. **Check logs** in `.agi/logs/`
4. **Review config** in `.agi/config.json`

---

**Documentation Version**: 2.0
**Last Updated**: 2026-02-08
**Completeness**: ✅ Full Coverage

*Complete technical documentation for ClosedWheeler AGI!* 📚🚀
