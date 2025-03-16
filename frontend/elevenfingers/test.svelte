// src/routes/game/+page.svelte
<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import WaitingRoom from '$lib/components/WaitingRoom.svelte';
  import TypingGame from '$lib/components/TypingGame.svelte';
  
  let ws: WebSocket;
  let currentView = 'waiting'; // 'waiting', 'game'
  let roomData = { players: {} };
  let gameData: any = null;
  
  onMount(() => {
    // Connect to WebSocket server
    ws = new WebSocket('ws://localhost:8080/ws'); // Replace with your actual WebSocket URL
    
    ws.onopen = () => {
      console.log('Connected to the server');
      // Automatically join the test room when the page loads
      joinGame();
    };
    
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      
      switch(data.type) {
        case 'roomStatus':
          roomData = data;
          currentView = 'waiting';
          break;
        case 'gameStart':
          gameData = data;
          currentView = 'game';
          break;
        default:
          console.log('Received unknown message type:', data.type);
      }
    };
    
    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
    
    ws.onclose = () => {
      console.log('Disconnected from the server');
    };
  });
  
  onDestroy(() => {
    if (ws) {
      ws.close();
    }
  });
  
  function joinGame() {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'join',
        content: { room: 'test' }
      }));
    }
  }
  
  function sendReady() {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'ready'
      }));
    }
  }
  
  function sendWordComplete(word: string) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'wordComplete',
        content: { word }
      }));
    }
  }
</script>

<svelte:head>
  <title>Tenfinger Typing Game</title>
  <meta name="description" content="Multiplayer typing game" />
</svelte:head>

<div class="container">
  {#if currentView === 'waiting'}
    <WaitingRoom 
      players={roomData.players} 
      onReady={sendReady}
    />
  {:else if currentView === 'game'}
    <TypingGame 
      gameData={gameData} 
      onWordComplete={sendWordComplete}
    />
  {/if}
</div>

<style>
  .container {
    max-width: 1000px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>

// src/routes/+page.svelte - Main landing page

<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import WaitingRoom from '$lib/components/WaitingRoom.svelte';
  import TypingGame from '$lib/components/TypingGame.svelte';
  
  let ws: WebSocket;
  let currentView = 'waiting'; // 'waiting', 'game'
  let roomData = { players: {} };
  let gameData: any = null;
  let typingGameComponent: TypingGame;
  
  onMount(() => {
    // Connect to WebSocket server
    ws = new WebSocket('ws://localhost:8080/ws'); // Replace with your actual WebSocket URL
    
    ws.onopen = () => {
      console.log('Connected to the server');
      // Automatically join the test room when the page loads
      joinGame();
    };
    
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      
      switch(data.type) {
        case 'roomStatus':
          roomData = data;
          currentView = 'waiting';
          break;
        case 'startGame':
          gameData = data;
          currentView = 'game';
          break;
        case 'userProgress':
          // Handle progress updates from server
          if (typingGameComponent && data.userid && typeof data.percentage === 'number') {
            typingGameComponent.updateProgress(data.userid, data.percentage);
          }
          break;
        default:
          console.log('Received unknown message type:', data.type);
      }
    };
    
    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
    
    ws.onclose = () => {
      console.log('Disconnected from the server');
    };
  });
  
  onDestroy(() => {
    if (ws) {
      ws.close();
    }
  });
  
  function joinGame() {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'join',
        content: { room: 'test' }
      }));
    }
  }
  
  function sendReady() {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'ready'
      }));
    }
  }
  
  function sendWordComplete(word: string) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'wordComplete',
        content: { word }
      }));
    }
  }
</script>

<svelte:head>
  <title>Tenfinger Typing Game</title>
  <meta name="description" content="Multiplayer typing game" />
</svelte:head>

<div class="container">
  {#if currentView === 'waiting'}
    <WaitingRoom 
      players={roomData.players} 
      onReady={sendReady}
    />
  {:else if currentView === 'game'}
    <TypingGame 
      gameData={gameData} 
      onWordComplete={sendWordComplete}
      bind:this={typingGameComponent}
    />
  {/if}
</div>

<style>
  .container {
    max-width: 1000px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>

// src/lib/components/WaitingRoom.svelte - Waiting room component
<script lang="ts">
  export let players: Record<string, boolean> = {};
  export let onReady: () => void;
  
  $: readyCount = Object.values(players).filter(ready => ready).length;
  $: totalPlayers = Object.keys(players).length;
</script>

<div class="waiting-room">
  <h2>Waiting Room</h2>
  
  <div class="player-list">
    <h3>Players ({readyCount}/{totalPlayers} ready)</h3>
    {#if totalPlayers === 0}
      <p>No players in room yet. Waiting for others to join...</p>
    {:else}
      <ul>
        {#each Object.entries(players) as [playerName, isReady]}
          <li class={isReady ? 'ready' : 'not-ready'}>
            {playerName} {isReady ? '✓' : ''}
          </li>
        {/each}
      </ul>
    {/if}
  </div>
  
  <div class="ready-section">
    <p>Once you're ready to play, click the button below:</p>
    <button on:click={onReady}>I'm Ready</button>
    <p class="hint">Game will start when all players are ready</p>
  </div>
</div>

<style>
  .waiting-room {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 2rem;
  }
  
  h2 {
    font-size: 2.5rem;
    margin-bottom: 2rem;
    text-align: center;
  }
  
  .player-list {
    width: 100%;
    max-width: 500px;
    margin-bottom: 2rem;
    padding: 1.5rem;
    border: 1px solid #ddd;
    border-radius: 8px;
    background-color: #f9f9f9;
  }
  
  h3 {
    margin-top: 0;
    margin-bottom: 1rem;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid #eee;
  }
  
  ul {
    list-style-type: none;
    padding: 0;
    margin: 0;
  }
  
  li {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid #eee;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  li:last-child {
    border-bottom: none;
  }
  
  .ready {
    color: #2ecc71;
    font-weight: 600;
  }
  
  .not-ready {
    color: #7f8c8d;
  }
  
  .ready-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  button {
    padding: 12px 28px;
    font-size: 1rem;
    font-weight: 600;
    color: white;
    background-color: #4a56e2;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s;
    margin: 1rem 0;
  }
  
  button:hover {
    background-color: #3643cf;
  }
  
  .hint {
    font-style: italic;
    color: #666;
    font-size: 0.9rem;
    margin-top: 0.5rem;
  }
</style>


// TypingGame.svelte

<script lang="ts">
  import { onMount, tick } from 'svelte';
  
  export let gameData: {
    text: string,
    StartTime: string,
    IsActive: string
  };
  export let onWordComplete: (word: string) => void;
  
  let gameText = '';
  let textArray: string[] = [];
  let currentPosition = 0;
  let typedText = '';
  
  let isActive = true;
  let startTime: Date;
  let currentWordStart = 0;
  let wordsTyped = 0;
  let accuracy = 100;
  let correctChars = 0;
  let totalChars = 0;
  let errorState = false; // Flag to track if user is in error state
  
  // Progress tracking - from server only
  let playerProgress: Record<string, number> = {};
  let currentUserId = ''; // This would be set from localStorage or session
  
  $: if (gameData) {
    gameText = gameData.text;
    textArray = gameText.split('');
    isActive = "TRUE";
    startTime = new Date(gameData.StartTime);
    currentPosition = 0;
    typedText = '';
    currentWordStart = 0;
    wordsTyped = 0;
    correctChars = 0;
    totalChars = 0;
    errorState = false;
  }
  
  $: if (totalChars > 0) {
    accuracy = Math.round((correctChars / totalChars) * 100);
  }
  
  async function handleKeydown(event: KeyboardEvent) {
    if (!isActive) return;
    
    // Only process printable characters and special keys like backspace
    if (event.key.length === 1 || event.key === 'Backspace') {
      event.preventDefault();
      
      if (event.key === 'Backspace') {
        if (currentPosition > 0) {
          currentPosition--;
          typedText = typedText.substring(0, currentPosition);
          
          // If we're backspacing out of an error state, clear it
          if (errorState && currentPosition >= 0) {
            // Check if we've backspaced to a correct state
            if (currentPosition === 0 || typedText.substring(0, currentPosition) === gameText.substring(0, currentPosition)) {
              errorState = false;
            }
          }
        }
      } else {
        // If in error state, don't allow typing except backspace
        if (errorState) {
          return;
        }
        
        // Regular character typed
        const isCorrect = event.key === textArray[currentPosition];
        
        if (isCorrect) {
          typedText = typedText + event.key;
          currentPosition++;
          
          // Update statistics
          totalChars++;
          correctChars++;
          
          // Check if a word is completed (space or end of text)
          if (event.key === ' ' || currentPosition === textArray.length) {
            const completedWord = gameText.substring(currentWordStart, currentPosition).trim();
            if (completedWord) {
              onWordComplete(completedWord);
              currentWordStart = currentPosition;
              wordsTyped++;
            }
          }
        } else {
          // Incorrect character - enter error state
          typedText = typedText + event.key;
          currentPosition++;
          errorState = true;
          
          // Update statistics
          totalChars++;
          // Don't increment correctChars
        }
      }
      
      await tick();
    }
  }
  
  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    
    // Get username from localStorage if available
    const username = localStorage.getItem('tenfinger_username');
    if (username) {
      currentUserId = username;
    } else {
      // Fallback to a random guest ID
      currentUserId = `Guest_${Math.floor(Math.random() * 10000)}`;
    }
    
    return () => {
      window.removeEventListener('keydown', handleKeydown);
    };
  });
  
  // Handle new websocket message for progress updates
  export function updateProgress(userId: string, percentage: number) {
    playerProgress[userId] = percentage;
    playerProgress = {...playerProgress}; // Trigger reactivity
  }
  
  function getCharClass(char: string, index: number) {
    if (index >= typedText.length) {
      return 'not-typed';
    }
    
    return typedText[index] === char ? 'correct' : 'incorrect';
  }
  
  // Sort players by progress for the leaderboard
  $: sortedPlayers = Object.entries(playerProgress)
    .sort(([, progressA], [, progressB]) => progressB - progressA);
</script>

<div class="game-container">
  <h2>Typing Game</h2>
  
  {#if isActive}
    <div class="stats">
      <div class="stat">
        <span class="label">Words:</span>
        <span class="value">{wordsTyped}</span>
      </div>
      <div class="stat">
        <span class="label">Accuracy:</span>
        <span class="value">{accuracy}%</span>
      </div>
      <div class="stat">
        <span class="label">Progress:</span>
        <span class="value">{playerProgress[currentUserId] || 0}%</span>
      </div>
    </div>
    
    <!-- Error state notification -->
    {#if errorState}
      <div class="error-message">
        Typo detected! Press backspace to correct the error.
      </div>
    {/if}
    
    <!-- Progress race track - based entirely on server data -->
    <div class="progress-race">
      <h3>Race Progress</h3>
      <div class="progress-container">
        {#each sortedPlayers as [playerId, progress]}
          <div class="player-progress">
            <div class="player-name">{playerId === currentUserId ? `${playerId} (You)` : playerId}</div>
            <div class="progress-bar-container">
              <div class="progress-bar" style="width: {progress}%"></div>
              <div class="progress-value">{progress}%</div>
            </div>
          </div>
        {/each}
      </div>
    </div>
    
    <div class="text-display" class:error-state={errorState} tabindex="0">
      {#each textArray as char, index}
        <span 
          class={`char ${getCharClass(char, index)} ${index === currentPosition ? 'current' : ''}`}
        >{char}</span>
      {/each}
    </div>
  {:else}
    <div class="waiting">
      <p>Waiting for the game to start...</p>
    </div>
  {/if}
</div>

<style>
  .game-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 2rem;
  }
  
  h2 {
    font-size: 2.5rem;
    margin-bottom: 2rem;
    text-align: center;
  }
  
  h3 {
    font-size: 1.5rem;
    margin-bottom: 1rem;
    text-align: center;
  }
  
  .stats {
    display: flex;
    justify-content: space-around;
    width: 100%;
    max-width: 800px;
    margin-bottom: 2rem;
  }
  
  .stat {
    padding: 0.5rem 1rem;
    background-color: #f5f5f5;
    border-radius: 4px;
    display: flex;
    flex-direction: column;
    align-items: center;
    min-width: 100px;
  }
  
  .label {
    font-size: 0.8rem;
    color: #666;
  }
  
  .value {
    font-size: 1.2rem;
    font-weight: bold;
    color: #333;
  }
  
  .error-message {
    background-color: #ffebee;
    color: #d32f2f;
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    margin-bottom: 1rem;
    font-weight: bold;
    border-left: 4px solid #d32f2f;
    animation: pulse 1.5s infinite;
  }
  
  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.7; }
  }
  
  .progress-race {
    width: 100%;
    max-width: 800px;
    margin-bottom: 2rem;
    padding: 1.5rem;
    background-color: #f8f9fa;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  .progress-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .player-progress {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .player-name {
    flex: 0 0 120px;
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .progress-bar-container {
    flex: 1;
    height: 24px;
    background-color: #e9ecef;
    border-radius: 12px;
    overflow: hidden;
    position: relative;
  }
  
  .progress-bar {
    height: 100%;
    background-color: #4a56e2;
    border-radius: 12px;
    transition: width 0.3s ease;
  }
  
  .progress-value {
    position: absolute;
    right: 10px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 0.8rem;
    font-weight: 600;
    color: #fff;
    text-shadow: 0 0 2px rgba(0, 0, 0, 0.5);
  }
  
  .text-display {
    width: 100%;
    max-width: 800px;
    min-height: 200px;
    background-color: #f8f9fa;
    padding: 2rem;
    border-radius: 8px;
    font-size: 1.5rem;
    line-height: 2;
    text-align: left;
    font-family: monospace;
    white-space: pre-wrap;
    position: relative;
    outline: none;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    transition: background-color 0.3s;
  }
  
  .text-display.error-state {
    background-color: #fff8f8;
    box-shadow: 0 0 0 2px #ffcdd2, 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  .char {
    position: relative;
    display: inline-block;
  }
  
  .current {
    position: relative;
  }
  
  .current::after {
    content: '';
    position: absolute;
    left: 0;
    bottom: 0;
    height: 2px;
    width: 100%;
    background-color: #4a56e2;
    animation: blink 1s infinite;
  }
  
  .correct {
    color: #2ecc71;
  }
  
  .incorrect {
    color: #e74c3c;
    text-decoration: underline;
  }
  
  .not-typed {
    color: #333;
  }
  
  .waiting {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    max-width: 800px;
    min-height: 200px;
    background-color: #f8f9fa;
    padding: 2rem;
    border-radius: 8px;
    font-size: 1.5rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  @keyframes blink {
    0%, 100% { opacity: 1; }
    50% { opacity: 0; }
  }
</style>

// src/routes/login/+page.svelte
<script lang="ts">
  import { goto } from '$app/navigation';
  
  let username = '';
  let password = '';
  let error = '';
  
  function handleSubmit() {
    // In a real app, you would handle authentication here
    // For now, we'll just do a basic validation and redirect
    if (!username) {
      error = 'Username is required';
      return;
    }
    
    // Store username in localStorage or a store
    localStorage.setItem('tenfinger_username', username);
    goto('/');
  }
</script>

<svelte:head>
  <title>Login - Tenfinger</title>
</svelte:head>

<div class="login-container">
  <div class="login-card">
    <h1>Login</h1>
    
    <form on:submit|preventDefault={handleSubmit}>
      {#if error}
        <div class="error">{error}</div>
      {/if}
      
      <div class="form-group">
        <label for="username">Username</label>
        <input 
          type="text" 
          id="username" 
          bind:value={username} 
          placeholder="Enter your username"
        />
      </div>
      
      <div class="form-group">
        <label for="password">Password (optional)</label>
        <input 
          type="password" 
          id="password" 
          bind:value={password} 
          placeholder="Enter your password"
        />
      </div>
      
      <button type="submit">Login</button>
    </form>
    
    <p class="note">
      Don't have an account? You can play as a guest by just entering a username.
    </p>
  </div>
</div>

<style>
  .login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 80vh;
    padding: 2rem;
  }
  
  .login-card {
    background-color: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 400px;
  }
  
  h1 {
    margin-top: 0;
    margin-bottom: 1.5rem;
    text-align: center;
    color: #333;
  }
  
  .form-group {
    margin-bottom: 1.5rem;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: #555;
  }
  
  input {
    width: 100%;
    padding: 0.75rem;
    font-size: 1rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    background-color: #f9f9f9;
  }
  
  input:focus {
    outline: none;
    border-color: #4a56e2;
    box-shadow: 0 0 0 2px rgba(74, 86, 226, 0.2);
  }
  
  button {
    width: 100%;
    padding: 0.75rem;
    background-color: #4a56e2;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  button:hover {
    background-color: #3643cf;
  }
  
  .error {
    background-color: #fdedee;
    color: #e74c3c;
    padding: 0.75rem;
    border-radius: 4px;
    margin-bottom: 1rem;
    font-size: 0.9rem;
  }
  
  .note {
    margin-top: 1.5rem;
    text-align: center;
    font-size: 0.9rem;
    color: #666;
  }
</style>

// src/app.css - Global styles
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen-Sans, Ubuntu, Cantarell, 'Helvetica Neue', sans-serif;
  line-height: 1.6;
  color: #333;
  background-color: #fff;
}

a {
  color: #4a56e2;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

h1, h2, h3, h4, h5, h6 {
  line-height: 1.2;
}

button {
  cursor: pointer;
}

// src/lib/websocket.ts - WebSocket service
import { writable } from 'svelte/store';

// Store for WebSocket connection status
export const wsStatus = writable('disconnected');

// Store for player data
export const players = writable({});

// Store for game data
export const gameData = writable(null);

let ws: WebSocket | null = null;
let reconnectTimeout: ReturnType<typeof setTimeout> | null = null;

export function connectWebSocket(url: string) {
  // Clear any existing reconnect timeout
  if (reconnectTimeout) {
    clearTimeout(reconnectTimeout);
    reconnectTimeout = null;
  }
  
  // Close existing connection
  if (ws) {
    ws.close();
  }
  
  wsStatus.set('connecting');
  
  ws = new WebSocket(url);
  
  ws.onopen = () => {
    wsStatus.set('connected');
    console.log('WebSocket connected');
  };
  
  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      
      switch (data.type) {
        case 'roomStatus':
          players.set(data.players || {});
          break;
        case 'gameStart':
          gameData.set(data);
          break;
        default:
          console.log('Received unknown message type:', data.type);
      }
    } catch (error) {
      console.error('Error parsing WebSocket message:', error);
    }
  };
  
  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
    wsStatus.set('error');
  };
  
  ws.onclose = () => {
    wsStatus.set('disconnected');
    console.log('WebSocket disconnected');
    
    // Set a reconnect timeout
    reconnectTimeout = setTimeout(() => {
      connectWebSocket(url);
    }, 5000);
  };
  
  return {
    send: (data: any) => {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify(data));
      } else {
        console.error('WebSocket not connected');
      }
    },
    close: () => {
      if (ws) {
        ws.close();
      }
    }
  };
}

// src/types.ts - TypeScript types
export interface Player {
  name: string;
  isReady: boolean;
}

export interface GameData {
  text: string;
  StartTime: string;
  IsActive: string;
}

export interface WordCompletePayload {
  type: 'wordComplete';
  content: {
    word: string;
  };
}

export interface JoinPayload {
  type: 'join';
  content: {
    room: string;
  };
}

export interface ReadyPayload {
  type: 'ready';
}

// src/hooks.server.ts - SvelteKit server hooks
// This file might be used for server-side authentication, but we'll keep it minimal
export async function handle({ event, resolve }) {
  // You can add server-side authentication here if needed
  return await resolve(event);
}
