# WordRush Server

Colyseus-based WebSocket server for multiplayer WordRush.

## 🚀 Quick Start

```bash
cd server
npm install
npm run dev
```

Server runs on `ws://localhost:2567`

## 📡 API

### Client → Server Messages

- `join` `{ name: string }` - Join/create room
- `select_tile` `{ row: number, col: number }` - Select/deselect tile
- `submit_word` `{}` - Submit current word
- `use_powerup` `{ powerUp: string }` - Use a power-up
- `clear_selection` `{}` - Clear all selected tiles
- `ready` `{}` - Mark player as ready

### Server → Client Messages

- `game_start` `{ players, currentPlayerIndex, timeLeft }`
- `timer_update` `{ timeLeft }`
- `word_update` `{ word }`
- `word_submitted` `{ playerId, word, score, newBoard, newScore, remainingPowerUps }`
- `powerup_used` `{ powerUp, playerId?, message?, newBoard? }`
- `selection_cleared` `{}`
- `player_ready` `{ playerId }`
- `game_over` `{ winner, finalScores }`
- `error` `{ message }`

## 🏗️ Architecture

- **Room**: `GameRoom` (max 4 clients)
- **State**: `GameState` (reactive, automatically synced)
- **Schemas**: `PlayerState`, `GameState` (using @colyseus/schema)

## 🧪 Testing

```bash
# Run server
npm run dev

# In another terminal, run client
cd ..
npm run dev
```

## 📦 Dependencies

- `colyseus` - Game server framework
- `express` - HTTP server (for monitor)
- `cors` - CORS headers

## 🔧 Configuration

- Port: `2567` (hardcoded for MVP)
- Max clients: `4`
- Dictionary: 100 common words (expandable)

## 📊 Scaling

For production:
- Use Redis adapter for horizontal scaling
- Deploy on Railway/Render/Heroku
- Add authentication
- Persist player profiles to database
- Add matchmaking queue