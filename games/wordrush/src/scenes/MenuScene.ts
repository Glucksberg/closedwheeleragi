import Phaser from 'phaser';

export class MenuScene extends Phaser.Scene {
  private playButton!: Phaser.GameObjects.Container;
  private howToPlayButton!: Phaser.GameObjects.Container;
  private settingsButton!: Phaser.GameObjects.Container;

  constructor() {
    super({ key: 'MenuScene' });
  }

  create() {
    const centerX = this.cameras.main.width / 2;
    const centerY = this.cameras.main.height / 2;

    // Title
    this.add.text(centerX, 120, 'WORD RUSH', {
      fontSize: '64px',
      fontFamily: 'Arial Black',
      color: '#ffffff',
      stroke: '#000000',
      strokeThickness: 6
    }).setOrigin(0.5);

    this.add.text(centerX, 180, 'Competitive Word Game', {
      fontSize: '20px',
      fontFamily: 'Arial',
      color: '#cccccc'
    }).setOrigin(0.5);

    // Create buttons
    this.createButton(centerX, 280, 'PLAY', () => this.startGame());
    this.createButton(centerX, 360, 'HOW TO PLAY', () => this.showHowToPlay());
    this.createButton(centerX, 440, 'SETTINGS', () => this.showSettings());

    // Footer
    this.add.text(centerX, 560, 'v0.1.0 | Browser Game', {
      fontSize: '14px',
      fontFamily: 'Arial',
      color: '#888888'
    }).setOrigin(0.5);

    // Animated background elements
    this.createBackgroundAnimation();
  }

  private createButton(x: number, y: number, text: string, callback: () => void) {
    const container = this.add.container(x, y);
    
    const bg = this.add.rectangle(0, 0, 200, 60, 0x4a5568, 1);
    bg.setInteractive({ useHandCursor: true });
    
    const label = this.add.text(0, 0, text, {
      fontSize: '24px',
      fontFamily: 'Arial',
      color: '#ffffff',
      fontStyle: 'bold'
    }).setOrigin(0.5);

    container.add([bg, label]);
    container.setSize(200, 60);

    // Hover effects
    bg.on('pointerover', () => {
      bg.setFillStyle(0x667eea);
      this.tweens.add({
        targets: container,
        scaleX: 1.05,
        scaleY: 1.05,
        duration: 100
      });
    });

    bg.on('pointerout', () => {
      bg.setFillStyle(0x4a5568);
      this.tweens.add({
        targets: container,
        scaleX: 1,
        scaleY: 1,
        duration: 100
      });
    });

    bg.on('pointerdown', callback);

    // Store references for later
    if (text === 'PLAY') this.playButton = container;
    else if (text === 'HOW TO PLAY') this.howToPlayButton = container;
    else if (text === 'SETTINGS') this.settingsButton = container;
  }

  private createBackgroundAnimation() {
    // Floating letters in background
    const letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
    for (let i = 0; i < 20; i++) {
      const x = Phaser.Math.Between(50, this.cameras.main.width - 50);
      const y = Phaser.Math.Between(50, this.cameras.main.height - 50);
      const letter = this.add.text(x, y, letters[Phaser.Math.Between(0, 25)], {
        fontSize: Phaser.Math.Between(12, 24) + 'px',
        color: 'rgba(255, 255, 255, 0.1)'
      }).setOrigin(0.5);

      this.tweens.add({
        targets: letter,
        y: y - 100,
        alpha: 0,
        duration: Phaser.Math.Between(3000, 6000),
        repeat: -1,
        yoyo: false,
        onRepeat: () => {
          letter.y = y + 100;
          letter.alpha = 0.1;
        }
      });
    }
  }

  private startGame() {
    this.scene.start('GameScene');
  }

  private showHowToPlay() {
    this.showModal(
      'How to Play',
      [
        '• Form words by selecting adjacent letters',
        '• Words must be 3+ letters long',
        '• Each letter has a point value',
        '• Use power-ups to gain advantage',
        '• Score highest in 2 minutes!',
        '• Multiplayer: Beat other players'
      ],
      'Got it!'
    );
  }

  private showSettings() {
    const settings = this.registry.get('settings') as any;
    
    this.showModal(
      'Settings',
      [
        `Sound: ${settings.soundEnabled ? 'ON' : 'OFF'}`,
        `Music: ${settings.musicEnabled ? 'ON' : 'OFF'}`,
        `Difficulty: ${settings.difficulty}`,
        '',
        '(Settings will be saved)'
      ],
      'Close',
      () => {}
    );
  }

  private showModal(title: string, lines: string[], buttonText: string, onClose?: () => void) {
    const overlay = this.add.rectangle(400, 300, 800, 600, 0x000000, 0.7);
    overlay.setInteractive();

    const modal = this.add.container(400, 300);
    
    const bg = this.add.rectangle(0, 0, 500, 400, 0x2d3748, 1, 4);
    bg.setStrokeStyle(2, 0x667eea);
    
    const titleText = this.add.text(0, -150, title, {
      fontSize: '28px',
      fontFamily: 'Arial',
      color: '#ffffff',
      fontStyle: 'bold'
    }).setOrigin(0.5);

    let yOffset = -80;
    lines.forEach(line => {
      const text = this.add.text(0, yOffset, line, {
        fontSize: '16px',
        fontFamily: 'Arial',
        color: '#e2e8f0',
        align: 'center'
      }).setOrigin(0.5);
      yOffset += 30;
    });

    const closeBtn = this.createModalButton(0, 120, buttonText, () => {
      overlay.destroy();
      modal.destroy();
      onClose?.();
    });

    modal.add([bg, titleText, closeBtn]);
    lines.forEach((_, i) => {
      const y = -80 + i * 30;
      const text = modal.list[i + 1] as Phaser.GameObjects.Text;
      if (text) modal.add(text);
    });

    this.input.on('pointerdown', () => {
      // Prevent clicks outside from closing (could add that feature)
    });
  }

  private createModalButton(x: number, y: number, text: string, callback: () => void) {
    const container = this.add.container(x, y);
    
    const bg = this.add.rectangle(0, 0, 120, 40, 0x667eea, 1);
    bg.setInteractive({ useHandCursor: true });
    
    const label = this.add.text(0, 0, text, {
      fontSize: '18px',
      fontFamily: 'Arial',
      color: '#ffffff',
      fontStyle: 'bold'
    }).setOrigin(0.5);

    container.add([bg, label]);

    bg.on('pointerover', () => bg.setFillStyle(0x5a67d8));
    bg.on('pointerout', () => bg.setFillStyle(0x667eea));
    bg.on('pointerdown', callback);

    return container;
  }
}