import { Room, Client } from 'colyseus';
import { Schema, type, MapSchema } from '@colyseus/schema';

export class PlayerState extends Schema {
  @type('string') name: string = '';
  @type('number') score: number = 0;
  @type('string') board: string = ''; // 16 chars for 4x4
  @type('number') selectedMask: number = 0; // 16-bit mask for selected tiles
  @type('string[]') powerUps: string[] = [];
  @type('boolean') isReady: boolean = false;
}

export class GameState extends Schema {
  @type({ map: PlayerState }) players = new MapSchema<PlayerState>();
  @type('number') timeLeft: number = 120;
  @type('number') currentPlayerIndex: number = 0;
  @type('string') currentWord: string = '';
  @type('boolean') gameOver: boolean = false;
  @type('string') winnerName: string = '';
}

export class GameRoom extends Room<GameState> {
  maxClients: number = 4;
  autoDispose: boolean = true;
  
  private timerInterval: NodeJS.Timeout | null = null;
  private dictionary: Set<string> = new Set();
  private powerUpPool: string[] = ['double', 'clear', 'freeze'];

  onCreate(options: any) {
    this.setState(new GameState());
    this.loadDictionary();
    
    this.onMessage('join', (client, message) => {
      this.handlePlayerJoin(client, message.name);
    });
    
    this.onMessage('select_tile', (client, message) => {
      this.handleTileSelect(client, message.row, message.col);
    });
    
    this.onMessage('submit_word', (client) => {
      this.handleSubmitWord(client);
    });
    
    this.onMessage('use_powerup', (client, message) => {
      this.handleUsePowerUp(client, message.powerUp);
    });
    
    this.onMessage('clear_selection', (client) => {
      this.handleClearSelection(client);
    });
    
    this.onMessage('ready', (client) => {
      this.handleReady(client);
    });

    this.onDispose(() => {
      if (this.timerInterval) {
        clearInterval(this.timerInterval);
      }
    });
  }

  private async loadDictionary() {
    // In production, load from file or database
    const commonWords = [
      'the','be','to','of','and','a','in','that','have','I','it','for','not','on','with','he','as','you','do','at','this','but','his','by','from','they','we','say','her','she','or','an','will','my','one','all','would','there','their','what','so','up','out','if','about','who','get','which','go','me','when','make','can','like','time','no','just','him','know','take','people','into','year','your','good','some','could','them','see','other','than','then','now','look','only','come','its','over','think','also','back','after','use','two','how','our','work','first','well','way','even','new','want','because','any','these','give','day','most','us'
    ];
    this.dictionary = new Set(commonWords.map(w => w.toLowerCase()));
  }

  private handlePlayerJoin(client: Client, name: string) {
    const player = new PlayerState();
    player.name = name || `Player_${this.state.players.size + 1}`;
    player.score = 0;
    player.board = this.generateBoardString();
    player.selectedMask = 0;
    player.powerUps = this.getRandomPowerUps(2);
    player.isReady = false;
    
    this.state.players.set(client.sessionId, player);
    
    console.log(`${player.name} joined. Total players: ${this.state.players.size}`);
    
    // Start game when we have at least 2 players or after 10 seconds
    if (this.state.players.size >= 2) {
      this.startGame();
    }
  }

  private generateBoardString(): string {
    const letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
    const weights = [13, 9, 12, 15, 10, 8, 9, 6, 5, 4, 4, 8, 7, 9, 15, 12, 6, 8, 7, 9, 15, 12, 6, 8, 7, 2];
    const totalWeight = weights.reduce((a, b) => a + b, 0);
    
    let board = '';
    for (let i = 0; i < 16; i++) {
      let random = Math.random() * totalWeight;
      let letterIndex = 0;
      for (let j = 0; j < weights.length; j++) {
        random -= weights[j];
        if (random <= 0) {
          letterIndex = j;
          break;
        }
      }
      board += letters[letterIndex];
    }
    return board;
  }

  private getRandomPowerUps(count: number): string[] {
    const powerUps: string[] = [];
    for (let i = 0; i < count; i++) {
      const index = Math.floor(Math.random() * this.powerUpPool.length);
      powerUps.push(this.powerUpPool[index]);
    }
    return powerUps;
  }

  private startGame() {
    if (this.timerInterval) return; // Already started
    
    console.log('Game starting!');
    
    // Set first player
    const firstClientId = this.state.players.keys().next().value;
    this.state.currentPlayerIndex = 0;
    
    // Start timer
    this.timerInterval = setInterval(() => {
      this.state.timeLeft--;
      this.broadcast('timer_update', { timeLeft: this.state.timeLeft });
      
      if (this.state.timeLeft <= 0) {
        this.endGame();
      }
    }, 1000);
    
    this.broadcast('game_start', {
      players: Array.from(this.state.players.values()),
      currentPlayerIndex: this.state.currentPlayerIndex,
      timeLeft: this.state.timeLeft
    });
  }

  private handleTileSelect(client: Client, row: number, col: number) {
    if (!this.isCurrentPlayer(client)) return;
    
    const player = this.state.players.get(client.sessionId);
    if (!player) return;
    
    const index = row * 4 + col;
    const mask = player.selectedMask;
    
    // Toggle selection
    if (mask & (1 << index)) {
      player.selectedMask = mask & ~(1 << index); // Deselect
    } else {
      // Check adjacency if not first selection
      if (mask !== 0) {
        if (!this.isAdjacent(mask, row, col)) {
          this.send(client, 'error', { message: 'Tiles must be adjacent' });
          return;
        }
      }
      player.selectedMask = mask | (1 << index); // Select
    }
    
    this.updateCurrentWord(client);
  }

  private isAdjacent(mask: number, row: number, col: number): boolean {
    // Check if any selected tile is adjacent to (row, col)
    for (let r = 0; r < 4; r++) {
      for (let c = 0; c < 4; c++) {
        const index = r * 4 + c;
        if (mask & (1 << index)) {
          if (Math.abs(r - row) <= 1 && Math.abs(c - col) <= 1) {
            return true;
          }
        }
      }
    }
    return false;
  }

  private updateCurrentWord(client: Client) {
    const player = this.state.players.get(client.sessionId);
    if (!player) return;
    
    let word = '';
    for (let i = 0; i < 16; i++) {
      if (player.selectedMask & (1 << i)) {
        word += player.board[i];
      }
    }
    
    this.state.currentWord = word;
    this.send(client, 'word_update', { word });
  }

  private handleSubmitWord(client: Client) {
    if (!this.isCurrentPlayer(client)) return;
    
    const player = this.state.players.get(client.sessionId);
    if (!player) return;
    
    const word = this.state.currentWord;
    if (word.length < 3) {
      this.send(client, 'error', { message: 'Word too short (min 3 letters)' });
      return;
    }
    
    if (!this.dictionary.has(word.toLowerCase())) {
      this.send(client, 'error', { message: `"${word}" not in dictionary` });
      return;
    }
    
    // Calculate score
    let score = this.calculateScore(word);
    
    // Check for double power-up
    const hasDouble = player.powerUps.includes('double');
    if (hasDouble) {
      score *= 2;
      player.powerUps = player.powerUps.filter(p => p !== 'double');
    }
    
    player.score += score;
    
    // Clear selected tiles and generate new letters
    const newBoard = player.board.split('');
    for (let i = 0; i < 16; i++) {
      if (player.selectedMask & (1 << i)) {
        newBoard[i] = this.getRandomLetter();
      }
    }
    player.board = newBoard.join('');
    player.selectedMask = 0;
    
    this.state.currentWord = '';
    
    // Broadcast updated state
    this.broadcast('word_submitted', {
      playerId: client.sessionId,
      word,
      score,
      newBoard: player.board,
      newScore: player.score,
      remainingPowerUps: player.powerUps
    });
    
    // Check win condition
    if (player.score >= 500) {
      this.endGame(player.name);
    }
  }

  private calculateScore(word: string): number {
    const letterScores: Record<string, number> = {
      'A': 1, 'B': 3, 'C': 3, 'D': 2, 'E': 1,
      'F': 4, 'G': 2, 'H': 4, 'I': 1, 'J': 8,
      'K': 5, 'L': 1, 'M': 3, 'N': 1, 'O': 1,
      'P': 3, 'Q': 10, 'R': 1, 'S': 1, 'T': 1,
      'U': 1, 'V': 4, 'W': 4, 'X': 8, 'Y': 4, 'Z': 10
    };
    
    let score = 0;
    for (const letter of word) {
      score += letterScores[letter] || 1;
    }
    
    if (word.length >= 5) score += 5;
    if (word.length >= 7) score += 10;
    
    return score;
  }

  private getRandomLetter(): string {
    const letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
    const weights = [13, 9, 12, 15, 10, 8, 9, 6, 5, 4, 4, 8, 7, 9, 15, 12, 6, 8, 7, 9, 15, 12, 6, 8, 7, 2];
    const totalWeight = weights.reduce((a, b) => a + b, 0);
    let random = Math.random() * totalWeight;
    
    let letterIndex = 0;
    for (let i = 0; i < weights.length; i++) {
      random -= weights[i];
      if (random <= 0) {
        letterIndex = i;
        break;
      }
    }
    
    return letters[letterIndex];
  }

  private handleUsePowerUp(client: Client, powerUp: string) {
    if (!this.isCurrentPlayer(client)) return;
    
    const player = this.state.players.get(client.sessionId);
    if (!player) return;
    
    const index = player.powerUps.indexOf(powerUp);
    if (index === -1) {
      this.send(client, 'error', { message: 'No more of this power-up' });
      return;
    }
    
    player.powerUps.splice(index, 1);
    
    switch (powerUp) {
      case 'clear':
        // Clear board for this player
        player.board = this.generateBoardString();
        player.selectedMask = 0;
        this.send(client, 'powerup_used', { powerUp, newBoard: player.board });
        break;
      case 'freeze':
        // For MVP: just notify
        this.broadcast('powerup_used', { 
          powerUp, 
          playerId: client.sessionId,
          message: `${player.name} used Freeze!` 
        });
        break;
      case 'double':
        // Will be applied on next word submission
        this.send(client, 'powerup_used', { powerUp, message: 'Double points next word!' });
        break;
    }
  }

  private handleClearSelection(client: Client) {
    if (!this.isCurrentPlayer(client)) return;
    
    const player = this.state.players.get(client.sessionId);
    if (!player) return;
    
    player.selectedMask = 0;
    this.state.currentWord = '';
    this.send(client, 'selection_cleared', {});
  }

  private handleReady(client: Client) {
    const player = this.state.players.get(client.sessionId);
    if (player) {
      player.isReady = true;
      this.broadcast('player_ready', { playerId: client.sessionId });
    }
  }

  private isCurrentPlayer(client: Client): boolean {
    const clientIds = Array.from(this.state.players.keys());
    if (clientIds.length === 0) return false;
    
    const currentIndex = this.state.currentPlayerIndex;
    return client.sessionId === clientIds[currentIndex];
  }

  private endGame(winnerName?: string) {
    if (this.timerInterval) {
      clearInterval(this.timerInterval);
      this.timerInterval = null;
    }
    
    this.state.gameOver = true;
    
    if (!winnerName) {
      // Find highest score
      let maxScore = 0;
      let winner = '';
      this.state.players.forEach(player => {
        if (player.score > maxScore) {
          maxScore = player.score;
          winner = player.name;
        }
      });
      this.state.winnerName = winner;
    } else {
      this.state.winnerName = winnerName;
    }
    
    this.broadcast('game_over', {
      winner: this.state.winnerName,
      finalScores: Array.from(this.state.players.values()).map(p => ({
        name: p.name,
        score: p.score
      }))
    });
  }

  onDispose() {
    console.log('GameRoom disposed');
  }
}