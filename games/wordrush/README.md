# WordRush 🎮

Competitive word game - Scrabble meets Battle Royale. Build words from falling letters and beat your opponents!

## 🚀 Quick Start

### Prerequisites
- Node.js 18+
- npm or yarn

### Installation
```bash
cd games/wordrush
npm install
```

### Development
```bash
npm run dev
```
Open http://localhost:3000 in your browser.

### Build
```bash
npm run build
```
Production files will be in `dist/` folder.

## 🎯 Gameplay

- **4x4 grid** per player
- **2-minute** rounds
- Form words by selecting **adjacent** letters
- Words must be **3+ letters**
- Use **power-ups** for advantages:
  - **2X Points**: Double score for next word
  - **Clear Board**: Refresh all letters
  - **Freeze Opponent**: Stop opponent for 5 seconds
- Highest score wins!

## 🛠️ Tech Stack

- **Frontend**: Phaser 3.60 + TypeScript
- **Build**: Vite + esbuild
- **Multiplayer**: Colyseus.js (WebSocket)
- **Deploy**: Cloudflare Pages (frontend) + Railway (backend)

## 📁 Project Structure

```
wordrush/
├── src/
│   ├── scenes/
│   │   ├── BootScene.ts    # Asset loading, dictionary init
│   │   ├── MenuScene.ts    # Main menu, settings
│   │   ├── GameScene.ts    # Core gameplay
│   │   └── ResultScene.ts  # Game over, scores
│   ├── main.ts             # Phaser config
│   └── utils/              # Helpers (to be added)
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── README.md
```

## 🎨 MVP Features (Current)

- ✅ Single-player vs AI (local 2-player hotseat)
- ✅ 4x4 letter grid with weighted random distribution
- ✅ Word validation (100-word dictionary)
- ✅ Scoring system (Scrabble-like)
- ✅ 3 power-ups (double, clear, freeze)
- ✅ 2-minute timer
- ✅ Basic UI (menu, game, results)
- ✅ Responsive design (FIT scaling)

## 🚧 Planned Features

### Phase 1 (Browser Multiplayer)
- [ ] Colyseus server integration
- [ ] Online matchmaking (2-4 players)
- [ ] Player profiles & persistent data
- [ ] Leaderboards
- [ ] Replay system

### Phase 2 (Monetization)
- [ ] Google AdSense (banner)
- [ ] Unity Ads (rewarded for power-ups)
- [ ] PWA manifest (install prompt)
- [ ] Email capture for updates

### Phase 3 (Mobile Port)
- [ ] Flutter/React Native port
- [ ] Push notifications
- [ ] In-app purchases (skins, themes)
- [ ] Battle pass system
- [ ] Store optimization (icons, screenshots)

## 🧪 Testing

```bash
# Run lint
npm run lint

# Build check
npm run build
```

## 📊 Metrics to Track

- Daily Active Users (DAU)
- Retention D1, D7
- Average session length
- Words per minute
- Power-up usage
- Ad impressions & CTR
- Conversion to mobile app

## 🐛 Known Issues

- Dictionary limited to 100 words (expand to 10k+)
- No sound effects/music (add Howler.js)
- No animations for word submission (add particle effects)
- Power-ups don't affect opponent (freeze not implemented)
- No server-side validation (vulnerable to cheating)

## 🤝 Contributing

This is a demo project. Fork and experiment!

## 📄 License

MIT

---

**Built with ❤️ using Phaser, TypeScript, and Vite**