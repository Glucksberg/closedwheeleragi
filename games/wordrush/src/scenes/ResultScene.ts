import Phaser from 'phaser';

export class ResultScene extends Phaser.Scene {
  private scores: Array<{name: string, score: number}> = [];
  private winner: any = null;

  constructor() {
    super({ key: 'ResultScene' });
  }

  create() {
    this.scores = this.registry.get('finalScores') as Array<{name: string, score: number}>;
    this.winner = this.registry.get('winner');
    
    this.cameras.main.fadeIn(500);
    
    this.drawResults();
  }

  private drawResults() {
    const centerX = this.cameras.main.width / 2;
    
    // Title
    this.add.text(centerX, 80, 'GAME OVER', {
      fontSize: '48px',
      fontFamily: 'Arial',
      color: '#ffffff',
      fontStyle: 'bold'
    }).setOrigin(0.5);

    // Winner announcement
    if (this.winner) {
      const winnerText = this.add.text(centerX, 150, `🏆 ${this.winner.name} Wins! 🏆`, {
        fontSize: '32px',
        fontFamily: 'Arial',
        color: '#ffd700',
        fontStyle: 'bold'
      }).setOrigin(0.5);
      
      this.tweens.add({
        targets: winnerText,
        scale: 1.1,
        duration: 1000,
        yoyo: true,
        repeat: -1
      });
    }

    // Scores table
    const startY = 220;
    const sortedScores = [...this.scores].sort((a, b) => b.score - a.score);
    
    sortedScores.forEach((player, index) => {
      const y = startY + index * 60;
      const isFirst = index === 0;
      
      // Background for top 3
      if (index < 3) {
        const bg = this.add.rectangle(centerX, y + 20, 400, 50, 0x000000, 0.5);
        bg.setStrokeStyle(2, isFirst ? 0xffd700 : 0x667eea);
      }
      
      const text = this.add.text(centerX, y, 
        `${index + 1}. ${player.name}: ${player.score} points`, {
        fontSize: isFirst ? '28px' : '24px',
        fontFamily: 'Arial',
        color: isFirst ? '#ffd700' : '#ffffff',
        fontStyle: isFirst ? 'bold' : 'normal'
      }).setOrigin(0.5);
    });

    // Buttons
    this.createButton(centerX - 120, 500, 'PLAY AGAIN', () => this.restartGame());
    this.createButton(centerX + 120, 500, 'MAIN MENU', () => this.goToMenu());
    
    // Stats
    const totalWords = this.calculateTotalWords();
    const avgWordLength = this.calculateAvgWordLength();
    
    this.add.text(centerX, 580, `Words played: ${totalWords} | Avg length: ${avgWordLength.toFixed(1)}`, {
      fontSize: '14px',
      fontFamily: 'Arial',
      color: '#aaaaaa'
    }).setOrigin(0.5);
  }

  private createButton(x: number, y: number, text: string, callback: () => void) {
    const container = this.add.container(x, y);
    
    const bg = this.add.rectangle(0, 0, 150, 50, 0x4a5568, 1);
    bg.setInteractive({ useHandCursor: true });
    
    const label = this.add.text(0, 0, text, {
      fontSize: '18px',
      fontFamily: 'Arial',
      color: '#ffffff',
      fontStyle: 'bold'
    }).setOrigin(0.5);

    container.add([bg, label]);

    bg.on('pointerover', () => bg.setFillStyle(0x667eea));
    bg.on('pointerout', () => bg.setFillStyle(0x4a5568));
    bg.on('pointerdown', callback);

    return container;
  }

  private restartGame() {
    this.scene.start('GameScene');
  }

  private goToMenu() {
    this.scene.start('MenuScene');
  }

  private calculateTotalWords(): number {
    // In full version, would track words submitted
    return Math.floor(Math.random() * 20) + 10; // Mock for MVP
  }

  private calculateAvgWordLength(): number {
    // Mock for MVP
    return 4.2;
  }
}