# 🚀 Quick Start: Fallback Models

**5-Minute Setup Guide**

---

## 🎯 What is Fallback?

If your primary model is slow or fails, the system automatically tries backup models - **without losing context or memory**.

---

## ⚡ Quick Setup

### Step 1: Edit `.agi/config.json`

Find these lines:
```json
"model": "gpt-4o-mini",
"fallback_models": [],
"fallback_timeout": 30,
```

### Step 2: Add Fallback Models

Change to:
```json
"model": "gpt-4o-mini",
"fallback_models": ["gpt-3.5-turbo"],
"fallback_timeout": 30,
```

### Step 3: Done!

That's it! Now if `gpt-4o-mini` takes more than 30 seconds or fails, the system will automatically try `gpt-3.5-turbo`.

---

## 📊 Common Configurations

### For Reliability (Recommended)

```json
"model": "gpt-4o-mini",
"fallback_models": ["gpt-3.5-turbo"],
"fallback_timeout": 30
```

### For Quality (Expensive → Cheap)

```json
"model": "gpt-4o",
"fallback_models": ["gpt-4o-mini", "gpt-3.5-turbo"],
"fallback_timeout": 30
```

### For Speed (Fast Failover)

```json
"model": "gpt-4o-mini",
"fallback_models": ["gpt-3.5-turbo"],
"fallback_timeout": 15
```

### No Fallback (Default)

```json
"model": "gpt-4o-mini",
"fallback_models": [],
"fallback_timeout": 30
```

---

## 🔍 How to Know It's Working?

### Logs

When fallback triggers, you'll see in `.agi/agent.log`:

```
[WARN] Primary model gpt-4o-mini failed: timeout
[INFO] Attempting fallback model 1/1: gpt-3.5-turbo
[INFO] Fallback model gpt-3.5-turbo succeeded!
```

### Check Logs

```bash
tail -f .agi/agent.log | grep -i fallback
```

If you see nothing, it means your primary model is working perfectly! (Which is good)

---

## 💡 Tips

1. **Always configure at least 1 fallback** for production
2. **Use 30s timeout** as default (works for most models)
3. **Order matters**: List models from best to worst quality
4. **Same API key**: All models must be available with your API key

---

## ❓ FAQ

**Q: Does fallback affect memory or tasks?**
A: No! Context is preserved 100%.

**Q: Does fallback cost more?**
A: Only if primary fails. If primary works, there's zero overhead.

**Q: Can I use different providers?**
A: Currently, all models must be from the same provider (same API endpoint).

**Q: How many fallback models can I add?**
A: As many as you want, but 1-2 is usually enough.

**Q: What if all models fail?**
A: You'll get an error, just like before. Fallback increases reliability but doesn't guarantee 100% uptime.

---

## 📚 Full Documentation

For detailed information, see:
- [FALLBACK_MODELS_GUIDE.md](FALLBACK_MODELS_GUIDE.md) - Complete guide (640+ lines)
- [FALLBACK_IMPLEMENTATION.md](FALLBACK_IMPLEMENTATION.md) - Technical details

---

**That's it!** Your AGI now has automatic failover protection. 🛡️

Configure it once and forget about it - you're protected! 🚀
