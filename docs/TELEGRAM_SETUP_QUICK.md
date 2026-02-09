# 📱 Telegram Setup - Quick Guide

**Setup Time**: < 3 minutes

---

## 🚀 During Initial Setup

### Option 1: Configure During Setup (Recommended)

```
📱 Telegram Integration (Optional)

Telegram allows you to control the agent remotely:
  • Chat with the agent from anywhere
  • Execute commands (/status, /logs, /model)
  • Approve sensitive operations

To get a bot token:
  1. Open Telegram and find @BotFather
  2. Send: /newbot
  3. Follow instructions to create your bot
  4. Copy the token (looks like: 1234567890:ABC...)

Configure Telegram now? (y/N): y

Enter Telegram Bot Token []: 1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
✅ Telegram token saved!
```

**Setup will save**:
- Token to `.env` (secure)
- `enabled: true` in `.agi/config.json`

### Option 2: Skip and Configure Later

```
Configure Telegram now? (y/N): n
⏭️  Skipping Telegram setup (you can configure it later in .agi/config.json)
```

---

## 🔗 Pairing Process

After setup completes with Telegram enabled:

### Step 1: Start the Agent

```bash
./ClosedWheeler
```

### Step 2: Open Your Bot in Telegram

Find your bot by the username you created with @BotFather.

### Step 3: Send `/start`

```
You: /start

Bot: 👋 Olá! Bem-vindo ao ClosedWheelerAGI

     Seu Chat ID: 123456789

     Configure este ID no config.json (campo `telegram.chat_id`)
     para ativar o controle remoto.

     Use /help para ver os comandos disponíveis.
```

### Step 4: Copy Your Chat ID

From the bot's response: `123456789`

### Step 5: Edit `.agi/config.json`

```json
{
  "telegram": {
    "enabled": true,
    "bot_token": "",
    "chat_id": 123456789,    // ← Add your Chat ID here
    "notify_on_tool_start": true
  }
}
```

### Step 6: Restart the Agent

```bash
# Stop with Ctrl+C
# Start again
./ClosedWheeler
```

### Step 7: Test It!

```
You: /status

Bot: 📊 AGI Status
     Memory: STM: 5 │ WM: 12 │ LTM: 45
     Project: ClosedWheelerAGI (27 files, Go)
```

**✅ Done! You're now connected!**

---

## 🎯 Complete Flow Example

```
┌─────────────────────────────────────────────┐
│ 1. Setup Wizard                             │
│    Configure Telegram? (y/N): y             │
│    Enter Token: 1234567890:ABC...           │
│    ✅ Token saved to .env                   │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│ 2. Start Agent                              │
│    ./ClosedWheeler                          │
│    [Agent starts with Telegram enabled]     │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│ 3. Open Telegram                            │
│    Find your bot                            │
│    Send: /start                             │
│    Get Chat ID: 123456789                   │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│ 4. Edit Config                              │
│    .agi/config.json                         │
│    Set: "chat_id": 123456789                │
│    Save file                                │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│ 5. Restart Agent                            │
│    Ctrl+C to stop                           │
│    ./ClosedWheeler to start                 │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│ 6. Test Connection                          │
│    Send: /status                            │
│    Bot responds with agent status           │
│    ✅ CONNECTED!                            │
└─────────────────────────────────────────────┘
```

---

## 🔧 Manual Configuration (If Skipped)

If you skipped Telegram during setup, configure it manually:

### 1. Get Bot Token

```
1. Open Telegram
2. Find @BotFather
3. Send: /newbot
4. Follow instructions
5. Copy token: 1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
```

### 2. Add to `.env`

```bash
# Add at the end of .env
TELEGRAM_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
```

### 3. Enable in `.agi/config.json`

```json
{
  "telegram": {
    "enabled": true,    // ← Change to true
    "bot_token": "",
    "chat_id": 0,
    "notify_on_tool_start": true
  }
}
```

### 4. Follow Pairing Steps Above

Continue from "Step 1: Start the Agent"

---

## 💡 Tips

### Security

- ✅ Token is in `.env` (gitignored)
- ✅ Only your Chat ID can control the agent
- ✅ All commands are logged in audit.log

### Testing

```bash
# Test bot responds
/start

# Test commands work
/status
/logs
/model

# Test chat works
"What files are in this project?"
```

### Troubleshooting

**Bot doesn't respond to /start**:
- Check token is correct in .env
- Restart the agent
- Check .agi/agent.log for errors

**Bot says "Access denied"**:
- You need to set `chat_id` in config.json
- Send /start to get your Chat ID
- Edit config and restart

**Bot responds but commands don't work**:
- Your Chat ID may be wrong
- Check config.json has correct chat_id
- Send /start again to verify

---

## 📱 Available Commands

Once connected:

```
/start   - Show your Chat ID
/help    - List all commands
/status  - Agent status and memory
/logs    - Recent logs
/diff    - Git diff
/model   - Show/change model
```

**Chat mode**:
Send any message without "/" to chat with the agent!

---

## ⏱️ Time Estimate

| Step | Time |
|------|------|
| Get bot token from BotFather | 30s |
| Configure in setup | 10s |
| Start agent | 5s |
| Get Chat ID | 10s |
| Edit config | 20s |
| Restart agent | 5s |
| Test connection | 10s |
| **Total** | **~90s** |

---

## 🎉 Summary

**Setup Flow**:
1. ✅ Setup wizard asks about Telegram
2. ✅ Enter token (or skip)
3. ✅ Agent shows pairing instructions
4. ✅ Send /start to bot
5. ✅ Copy Chat ID to config
6. ✅ Restart agent
7. ✅ **Connected!**

**You can now**:
- 💬 Chat with agent from anywhere
- 🔧 Execute commands remotely
- ✅ Approve operations via Telegram
- 📊 Monitor status on the go

---

**Status**: ✅ **IMPLEMENTED**
**Complexity**: 🟢 **Simple**
**Time**: ⏱️ **< 3 minutes**

*Control your agent from anywhere! 🌎📱*
