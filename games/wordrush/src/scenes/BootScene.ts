import Phaser from 'phaser';

export class BootScene extends Phaser.Scene {
  constructor() {
    super({ key: 'BootScene' });
  }

  preload() {
    // Create loading bar
    const width = this.cameras.main.width;
    const height = this.cameras.main.height;
    
    const progressBar = this.add.graphics();
    const progressBox = this.add.graphics();
    progressBox.fillStyle(0x222222, 0.8);
    progressBox.fillRect(width / 2 - 160, height / 2 - 25, 320, 50);
    
    const loadingText = this.add.text(width / 2, height / 2 - 50, 'Loading...', {
      font: '20px monospace',
      fill: '#ffffff'
    }).setOrigin(0.5, 0.5);
    
    this.load.on('progress', (value: number) => {
      progressBar.clear();
      progressBar.fillStyle(0x00ff00, 1);
      progressBar.fillRect(width / 2 - 150, height / 2 - 15, 300 * value, 30);
    });
    
    this.load.on('complete', () => {
      progressBar.destroy();
      progressBox.destroy();
      loadingText.destroy();
    });

    // Load assets (using generated graphics for MVP)
    this.load.image('tile', 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAYAAABytg0kAAAAFklEQVQYV2NkYGD4z4AE/lkY/8bGJiQAABp8BZa5cRcAAAAASUVORK5CYII=');
    this.load.image('letter-bg', 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAA8AAAAPCAYAAAA71pVKAAAAFklEQVQoU2NkYGBg+OaA4O8BAOMCAF8AADtQB8Cv7dYcAAAAAElFTkSuQmCC');
  }

  create() {
    // Initialize game settings
    this.registry.set('settings', {
      soundEnabled: true,
      musicEnabled: true,
      difficulty: 'normal'
    });

    // Load dictionary (minimal for MVP)
    this.loadDictionary();

    this.scene.start('MenuScene');
  }

  async loadDictionary() {
    // Minimal dictionary for MVP - 100 most common English words
    const commonWords = [
      'the','be','to','of','and','a','in','that','have','I','it','for','not','on','with','he','as','you','do','at','this','but','his','by','from','they','we','say','her','she','or','an','will','my','one','all','would','there','their','what','so','up','out','if','about','who','get','which','go','me','when','make','can','like','time','no','just','him','know','take','people','into','year','your','good','some','could','them','see','other','than','then','now','look','only','come','its','over','think','also','back','after','use','two','how','our','work','first','well','way','even','new','want','because','any','these','give','day','most','us'
    ];

    this.registry.set('dictionary', new Set(commonWords.map(w => w.toLowerCase())));
  }
}