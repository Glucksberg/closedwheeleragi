import Phaser from 'phaser';

interface Player {
  id: number;
  name: string;
  score: number;
  board: string[][];
  selected: boolean[][];
  powerUps: string[];
}

export class GameScene extends Phaser.Scene {
  private players: Player[] = [];
  private currentPlayerIndex: number = 0;
  private gridSize: number = 4;
  private tileSize: number = 80;
  private tileSpacing: number = 4;
  private boardStartX!: number;
  private boardStartY!: number;
  private timer!: Phaser.Time.TimerEvent;
  private timeLeft: number = 120; // 2 minutes
  private isGameOver: boolean = false;
  private dictionary: Set<string> = new Set();
  private currentWord: string = '';
  private selectedTiles: {row: number, col: number}[] = [];
  private powerUpButtons: Phaser.GameObjects.Container[] = [];

  constructor() {
    super({ key: 'GameScene' });
  }

  create() {
    this.dictionary = this.registry.get('dictionary') as Set<string>;
    
    // Initialize players (local multiplayer for MVP)
    this.initializePlayers();

    // Calculate board position
    this.boardStartX = (this.cameras.main.width - (this.gridSize * (this.tileSize + this.tileSpacing))) / 2;
    this.boardStartY = 100;

    // Draw UI
    this.drawUI();

    // Draw boards for each player
    this.drawBoards();

    // Setup input
    this.setupInput();

    // Start timer
    this.startTimer();

    // Setup power-ups UI
    this.setupPowerUps();

    // Instructions
    this.showTemporaryText('Select adjacent letters to form words!', 3000);
  }

  private initializePlayers() {
    // For MVP: 2 local players
    this.players = [
      {
        id: 0,
        name: 'Player 1',
        score: 0,
        board: this.generateBoard(),
        selected: Array(this.gridSize).fill(null).map(() => Array(this.gridSize).fill(false)),
        powerUps: ['double', 'clear']
      },
      {
        id: 1,
        name: 'Player 2',
        score: 0,
        board: this.generateBoard(),
        selected: Array(this.gridSize).fill(null).map(() => Array(this.gridSize).fill(false)),
        powerUps: ['freeze', 'double']
      }
    ];
  }

  private generateBoard(): string[][] {
    const letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
    const board: string[][] = [];
    
    for (let row = 0; row < this.gridSize; row++) {
      board[row] = [];
      for (let col = 0; col < this.gridSize; col++) {
        // Weighted random - common letters more frequent
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
        
        board[row][col] = letters[letterIndex];
      }
    }
    
    return board;
  }

  private drawUI() {
    const centerX = this.cameras.main.width / 2;

    // Timer at top
    this.timerText = this.add.text(centerX, 30, this.formatTime(this.timeLeft), {
      fontSize: '36px',
      fontFamily: 'Arial',
      color: '#ffffff',
      fontStyle: 'bold'
    }).setOrigin(0.5);

    // Player scores
    this.updateScoreDisplay();

    // Current player indicator
    this.currentPlayerText = this.add.text(centerX, 70, `${this.players[this.currentPlayerIndex].name}'s Turn`, {
      fontSize: '20px',
      fontFamily: 'Arial',
      color: '#ffd700'
    }).setOrigin(0.5);
  }

  private drawBoards() {
    // Clear existing boards
    this.children.list.forEach(child => {
      if (child.texture?.key === 'tile-bg' || child.texture?.key === 'tile-letter') {
        child.destroy();
      }
    });

    // Draw both players' boards
    this.players.forEach((player, playerIndex) => {
      const offsetX = playerIndex === 0 ? -200 : 200;
      const boardX = this.boardStartX + offsetX;
      
      // Player label
      this.add.text(boardX + (this.gridSize * this.tileSize) / 2, this.boardStartY - 30, 
        `${player.name} (${player.score})`, {
        fontSize: '18px',
        fontFamily: 'Arial',
        color: playerIndex === this.currentPlayerIndex ? '#ffd700' : '#aaaaaa'
      }).setOrigin(0.5);

      // Draw tiles
      for (let row = 0; row < this.gridSize; row++) {
        for (let col = 0; col < this.gridSize; col++) {
          const x = boardX + col * (this.tileSize + this.tileSpacing);
          const y = this.boardStartY + row * (this.tileSize + this.tileSpacing);
          
          // Tile background
          const bg = this.add.rectangle(x + this.tileSize/2, y + this.tileSize/2, this.tileSize, this.tileSize, 0x4a5568, 1);
          bg.setStrokeStyle(2, 0x667eea);
          
          // Letter
          const letter = player.board[row][col];
          const letterText = this.add.text(x + this.tileSize/2, y + this.tileSize/2, letter, {
            fontSize: '32px',
            fontFamily: 'Arial',
            color: '#ffffff',
            fontStyle: 'bold'
          }).setOrigin(0.5);
          
          // Store reference
          (bg as any).row = row;
          (bg as any).col = col;
          (bg as any).playerId = playerIndex;
          (bg as any).isTile = true;
          
          bg.setInteractive({ useHandCursor: true });
          
          // Click handler
          bg.on('pointerdown', () => this.onTileClick(playerIndex, row, col));
          
          // Selection visual
          if (player.selected[row][col]) {
            bg.setFillStyle(0x667eea);
          }
        }
      }
    });
  }

  private onTileClick(playerId: number, row: number, col: number) {
    if (this.isGameOver) return;
    if (playerId !== this.currentPlayerIndex) return;
    
    const player = this.players[playerId];
    
    // Check if adjacent to last selected tile
    if (this.selectedTiles.length > 0) {
      const last = this.selectedTiles[this.selectedTiles.length - 1];
      const isAdjacent = Math.abs(last.row - row) <= 1 && Math.abs(last.col - col) <= 1;
      if (!isAdjacent) {
        this.showTemporaryText('Tiles must be adjacent!', 1500);
        return;
      }
    }
    
    // Toggle selection
    if (player.selected[row][col]) {
      // Deselect
      player.selected[row][col] = false;
      this.selectedTiles = this.selectedTiles.filter(t => t.row !== row || t.col !== col);
      this.currentWord = this.currentWord.replace(player.board[row][col], '');
    } else {
      // Select
      player.selected[row][col] = true;
      this.selectedTiles.push({row, col});
      this.currentWord += player.board[row][col];
    }
    
    this.drawBoards();
    this.updateCurrentWordDisplay();
  }

  private updateCurrentWordDisplay() {
    // Remove existing word display
    this.children.list.forEach(child => {
      if ((child as any).isWordDisplay) {
        child.destroy();
      }
    });

    if (this.currentWord.length > 0) {
      const wordText = this.add.text(this.cameras.main.width / 2, 550, this.currentWord, {
        fontSize: '28px',
        fontFamily: 'Arial',
        color: '#ffffff',
        backgroundColor: '#000000',
        padding: { x: 10, y: 5 }
      }).setOrigin(0.5);
      (wordText as any).isWordDisplay = true;
    }
  }

  private setupInput() {
    // Submit word button (Enter key)
    this.input.keyboard?.on('keydown-ENTER', () => this.submitWord());
    
    // Clear selection button (Escape)
    this.input.keyboard?.on('keydown-ESC', () => this.clearSelection());
  }

  private setupPowerUps() {
    const player = this.players[this.currentPlayerIndex];
    const startX = 50;
    const y = 500;
    
    player.powerUps.forEach((powerUp, index) => {
      const x = startX + index * 120;
      const btn = this.createPowerUpButton(x, y, powerUp);
      this.powerUpButtons.push(btn);
    });
  }

  private createPowerUpButton(x: number, y: number, powerUp: string) {
    const container = this.add.container(x, y);
    
    const bg = this.add.rectangle(0, 0, 100, 60, 0x805ad5, 1);
    bg.setInteractive({ useHandCursor: true });
    
    const label = this.add.text(0, 0, this.getPowerUpName(powerUp), {
      fontSize: '14px',
      fontFamily: 'Arial',
      color: '#ffffff',
      align: 'center'
    }).setOrigin(0.5);

    container.add([bg, label]);
    
    bg.on('pointerover', () => bg.setFillStyle(0x6b46c1));
    bg.on('pointerout', () => bg.setFillStyle(0x805ad5));
    bg.on('pointerdown', () => this.usePowerUp(powerUp));
    
    return container;
  }

  private getPowerUpName(key: string): string {
    const names: Record<string, string> = {
      'double': '2X Points',
      'clear': 'Clear Board',
      'freeze': 'Freeze Opponent'
    };
    return names[key] || key;
  }

  private usePowerUp(powerUp: string) {
    if (this.isGameOver) return;
    
    const player = this.players[this.currentPlayerIndex];
    const powerUpIndex = player.powerUps.indexOf(powerUp);
    if (powerUpIndex === -1) {
      this.showTemporaryText('No more of this power-up!', 1500);
      return;
    }
    
    // Remove power-up
    player.powerUps.splice(powerUpIndex, 1);
    
    switch (powerUp) {
      case 'clear':
        this.clearBoard();
        break;
      case 'freeze':
        this.freezeOpponent();
        break;
      case 'double':
        this.doublePointsNext = true;
        this.showTemporaryText('Double points next word!', 2000);
        break;
    }
    
    this.setupPowerUps(); // Refresh power-up buttons
  }

  private clearBoard() {
    const player = this.players[this.currentPlayerIndex];
    for (let row = 0; row < this.gridSize; row++) {
      for (let col = 0; col < this.gridSize; col++) {
        player.board[row][col] = this.getRandomLetter();
        player.selected[row][col] = false;
      }
    }
    this.selectedTiles = [];
    this.currentWord = '';
    this.drawBoards();
    this.updateCurrentWordDisplay();
    this.showTemporaryText('Board cleared!', 1500);
  }

  private freezeOpponent() {
    // For MVP: just show message
    this.showTemporaryText('Opponent frozen for 5 seconds!', 2000);
    // In full version: would disable opponent input temporarily
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

  private submitWord() {
    if (this.currentWord.length < 3) {
      this.showTemporaryText('Word too short! (min 3 letters)', 1500);
      return;
    }
    
    const lowerWord = this.currentWord.toLowerCase();
    if (!this.dictionary.has(lowerWord)) {
      this.showTemporaryText(`"${this.currentWord}" not in dictionary!`, 1500);
      return;
    }
    
    // Calculate score
    let baseScore = this.calculateWordScore(this.currentWord);
    if (this.doublePointsNext) {
      baseScore *= 2;
      this.doublePointsNext = false;
    }
    
    const player = this.players[this.currentPlayerIndex];
    player.score += baseScore;
    
    // Clear used tiles and generate new ones
    for (const tile of this.selectedTiles) {
      player.board[tile.row][tile.col] = this.getRandomLetter();
      player.selected[tile.row][tile.col] = false;
    }
    
    this.selectedTiles = [];
    this.currentWord = '';
    
    this.updateScoreDisplay();
    this.drawBoards();
    this.updateCurrentWordDisplay();
    
    this.showTemporaryText(`+${baseScore} points!`, 1000);
    
    // Check for win condition (score threshold)
    if (player.score >= 500) {
      this.endGame(player);
    }
  }

  private calculateWordScore(word: string): number {
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
    
    // Bonus for longer words
    if (word.length >= 5) score += 5;
    if (word.length >= 7) score += 10;
    
    return score;
  }

  private clearSelection() {
    const player = this.players[this.currentPlayerIndex];
    for (const tile of this.selectedTiles) {
      player.selected[tile.row][tile.col] = false;
    }
    this.selectedTiles = [];
    this.currentWord = '';
    this.drawBoards();
    this.updateCurrentWordDisplay();
  }

  private startTimer() {
    this.timer = this.time.addEvent({
      delay: 1000,
      callback: this.onTimerTick,
      callbackScope: this,
      loop: true
    });
  }

  private onTimerTick() {
    this.timeLeft--;
    this.timerText.setText(this.formatTime(this.timeLeft));
    
    if (this.timeLeft <= 10) {
      this.tweens.add({
        targets: this.timerText,
        scale: 1.2,
        duration: 200,
        yoyo: true
      });
      this.timerText.setColor('#ff4444');
    }
    
    if (this.timeLeft <= 0) {
      this.endGame(null); // Time's up, winner is highest score
    }
  }

  private formatTime(seconds: number): string {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  }

  private updateScoreDisplay() {
    // Remove old score texts
    this.children.list.forEach(child => {
      if ((child as any).isScoreDisplay) {
        child.destroy();
      }
    });

    // Draw new scores
    this.players.forEach((player, index) => {
      const offsetX = index === 0 ? -200 : 200;
      const x = this.boardStartX + offsetX + (this.gridSize * this.tileSize) / 2;
      const y = this.boardStartY + this.gridSize * (this.tileSize + this.tileSpacing) + 40;
      
      const scoreText = this.add.text(x, y, `Score: ${player.score}`, {
        fontSize: '20px',
        fontFamily: 'Arial',
        color: '#ffffff',
        fontStyle: 'bold'
      }).setOrigin(0.5);
      (scoreText as any).isScoreDisplay = true;
    });
  }

  private showTemporaryText(message: string, duration: number) {
    const text = this.add.text(this.cameras.main.width / 2, this.cameras.main.height - 50, message, {
      fontSize: '18px',
      fontFamily: 'Arial',
      color: '#ffff00',
      backgroundColor: '#000000',
      padding: { x: 10, y: 5 }
    }).setOrigin(0.5);
    
    this.tweens.add({
      targets: text,
      alpha: 0,
      y: text.y - 20,
      duration: duration,
      onComplete: () => text.destroy()
    });
  }

  private endGame(winner: Player | null) {
    this.isGameOver = true;
    this.timer.remove();
    
    if (!winner) {
      // Time's up - find winner
      winner = this.players.reduce((prev, current) => 
        prev.score > current.score ? prev : current
      );
    }
    
    this.registry.set('finalScores', this.players.map(p => ({ name: p.name, score: p.score })));
    this.registry.set('winner', winner);
    
    this.cameras.main.fade(500, 0, 0, 0, false, (camera, progress) => {
      if (progress === 1) {
        this.scene.start('ResultScene');
      }
    });
  }

  private doublePointsNext: boolean = false;
  private timerText!: Phaser.GameObjects.Text;
  private currentPlayerText!: Phaser.GameObjects.Text;
}