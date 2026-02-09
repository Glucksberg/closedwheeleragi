import { Client, Room, type Schema } from 'colyseus.js';

export interface WordRushPlayer {
  name: string;
  score: number;
  board: string; // 16 chars
  selectedMask: number;
  powerUps: string[];
  isReady: boolean;
}

export interface WordRushState {
  players: Map<string, WordRushPlayer>;
  timeLeft: number;
  currentPlayerIndex: number;
  currentWord: string;
  gameOver: boolean;
  winnerName: string;
}

export class WordRushClient {
  private client: Client;
  private room: Room<WordRushState> | null = null;
  private onStateChange?: (state: WordRushState) => void;
  private onMessage?: (type: string, data: any) => void;
  private onError?: (error: any) => void;

  constructor(serverUrl: string = 'ws://localhost:2567') {
    this.client = new Client(serverUrl);
  }

  async joinOrCreateRoom(roomName: string, playerName: string): Promise<Room<WordRushState>> {
    try {
      this.room = await this.client.joinOrCreate<WordRushState>(roomName, {
        name: playerName
      });
      
      this.setupListeners();
      return this.room;
    } catch (error) {
      console.error('Failed to join room:', error);
      throw error;
    }
  }

  async joinExistingRoom(roomId: string, playerName: string): Promise<Room<WordRushState>> {
    try {
      this.room = await this.client.joinById<WordRushState>(roomId, {
        name: playerName
      });
      
      this.setupListeners();
      return this.room;
    } catch (error) {
      console.error('Failed to join room:', error);
      throw error;
    }
  }

  private setupListeners() {
    if (!this.room) return;

    this.room.onStateChange((state) => {
      this.onStateChange?.(state);
    });

    this.room.onMessage('game_start', (data) => {
      this.onMessage?.('game_start', data);
    });

    this.room.onMessage('timer_update', (data) => {
      this.onMessage?.('timer_update', data);
    });

    this.room.onMessage('word_update', (data) => {
      this.onMessage?.('word_update', data);
    });

    this.room.onMessage('word_submitted', (data) => {
      this.onMessage?.('word_submitted', data);
    });

    this.room.onMessage('powerup_used', (data) => {
      this.onMessage?.('powerup_used', data);
    });

    this.room.onMessage('selection_cleared', (data) => {
      this.onMessage?.('selection_cleared', data);
    });

    this.room.onMessage('player_ready', (data) => {
      this.onMessage?.('player_ready', data);
    });

    this.room.onMessage('game_over', (data) => {
      this.onMessage?.('game_over', data);
    });

    this.room.onMessage('error', (data) => {
      this.onMessage?.('error', data);
    });

    this.room.onError((error) => {
      this.onError?.(error);
    });

    this.room.onLeave(() => {
      this.room = null;
    });
  }

  // Message sending
  selectTile(row: number, col: number) {
    this.room?.send('select_tile', { row, col });
  }

  submitWord() {
    this.room?.send('submit_word', {});
  }

  usePowerUp(powerUp: string) {
    this.room?.send('use_powerup', { powerUp });
  }

  clearSelection() {
    this.room?.send('clear_selection', {});
  }

  setReady() {
    this.room?.send('ready', {});
  }

  // Callbacks
  onStateChangeCallback(callback: (state: WordRushState) => void) {
    this.onStateChange = callback;
  }

  onMessageCallback(callback: (type: string, data: any) => void) {
    this.onMessage = callback;
  }

  onErrorCallback(callback: (error: any) => void) {
    this.onError = callback;
  }

  // Getters
  getRoom(): Room<WordRushState> | null {
    return this.room;
  }

  getSessionId(): string | null {
    return this.room?.sessionId ?? null;
  }

  async leave() {
    await this.room?.leave();
  }

  async disconnect() {
    await this.client.leaveAll();
  }
}