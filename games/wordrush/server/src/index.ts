import express from 'express';
import { createServer } from 'http';
import { Server } from 'colyseus';
import { GameRoom } from './rooms/GameRoom';
import { monitor } from '@colyseus/monitor';

const app = express();
app.use(express.json());

const server = createServer(app);
const gameServer = new Server({ server });

// Register game room
gameServer.define('wordrush_game', GameRoom);

// Monitor endpoint (for debugging)
app.use('/colyseus', monitor());

// Start server
gameServer.listen(2567).then(() => {
  console.log('🚀 WordRush Server running on ws://localhost:2567');
  console.log('📊 Monitor available at http://localhost:2567/colyseus');
});

// Graceful shutdown
process.on('SIGTERM', () => {
  gameServer.disconnect();
  process.exit(0);
});