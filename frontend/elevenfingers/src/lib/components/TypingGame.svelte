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
  let correctWordPositions: number[] = []; // Track positions where words were completed correctly
  
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
    correctWordPositions = [];
  }
  
  $: if (totalChars > 0) {
    accuracy = Math.round((correctChars / totalChars) * 100);
  }
  
  // Helper function to check if we're trying to backspace into a correctly typed word
  function isBackspaceBlockedPosition(position: number): boolean {
    return correctWordPositions.some(wordEndPos => position <= wordEndPos);
  }
  
  async function handleKeydown(event: KeyboardEvent) {
    if (!isActive) return;
    
    // Only process printable characters and special keys like backspace
    if (event.key.length === 1 || event.key === 'Backspace') {
      event.preventDefault();
      
      if (event.key === 'Backspace') {
        // Check if we're trying to backspace into a correctly completed word
        if (currentPosition > 0 && !isBackspaceBlockedPosition(currentPosition - 1)) {
          currentPosition--;
          typedText = typedText.substring(0, currentPosition);
          
          // If we're backspacing out of an error state, clear it
          if (errorState && currentPosition >= 0) {
            // Check if we've backspaced to a correct state
            if (currentPosition === 0 || typedText.substring(0, currentPosition) === gameText.substring(0, currentPosition)) {
              errorState = false;
            }
          }
        } else if (isBackspaceBlockedPosition(currentPosition - 1)) {
          // Display feedback that backspace is blocked for correct words
          showBackspaceBlockedMessage();
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
              
              // Mark this position as the end of a correctly typed word
              correctWordPositions.push(currentPosition - 1);
              
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
  
  // Message display for blocked backspace
  let showingBackspaceMessage = false;
  let backspaceMessageTimeout: ReturnType<typeof setTimeout> | null = null;
  
  function showBackspaceBlockedMessage() {
    showingBackspaceMessage = true;
    
    // Clear any existing timeout
    if (backspaceMessageTimeout) {
      clearTimeout(backspaceMessageTimeout);
    }
    
    // Hide the message after 2 seconds
    backspaceMessageTimeout = setTimeout(() => {
      showingBackspaceMessage = false;
    }, 2000);
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
      if (backspaceMessageTimeout) {
        clearTimeout(backspaceMessageTimeout);
      }
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
    
    <!-- Backspace blocked message -->
    {#if showingBackspaceMessage}
      <div class="info-message">
        Backspace is disabled for correctly typed words!
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
          class={`char ${getCharClass(char, index)} ${index === currentPosition ? 'current' : ''} ${correctWordPositions.includes(index) ? 'locked' : ''}`}
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
  
  .info-message {
    background-color: #e8f5e9;
    color: #2e7d32;
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    margin-bottom: 1rem;
    font-weight: bold;
    border-left: 4px solid #2e7d32;
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
  
  .locked {
    font-weight: bold;
    background-color: rgba(46, 204, 113, 0.1);
    border-radius: 2px;
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
