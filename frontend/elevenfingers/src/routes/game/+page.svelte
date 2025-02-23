<script>
  import { onMount, onDestroy } from 'svelte';

  const sampleText = "The quick brown fox jumps over the lazy dog";
  let isGameActivate = false
  let currentIndex = 0;
  let mustCorrectError = false;
  let errorPositions = new Set();
  let ws;
  let currentWord = '';
  let wordStart = 0;

  let characters = sampleText.split('').map(char => ({
    char,
    status: 'pending',
    hasRedHighlight: false
  }));

  onMount(() => {
    containerRef.focus();
    ws = new WebSocket('ws://127.0.0.1:9000/ws');
    
    ws.onopen = () => {
      console.log('Connected to WebSocket server');
    };
  });

  onDestroy(() => {
    if (ws) ws.close();
  });
    
  function startGame(){
      ws.send(JSON.stringify({
          type: 'startGame'
      }))
  }

  function sendWordProgress(word) {
    if (!ws || ws.readyState !== WebSocket.OPEN) return;
    
    ws.send(JSON.stringify({
      type: 'word_complete',
      content: word,
    }));
  }

  function handleKeyPress(event) {
    // Handle backspace - only when there's an error
    if (event.key === 'Backspace') {
      if (mustCorrectError) {
        characters[currentIndex].status = 'pending';
        characters[currentIndex].hasRedHighlight = false;
        mustCorrectError = false;
      }
      event.preventDefault();
      return;
    }

    // Prevent arrow keys
    if (event.key === 'ArrowLeft' || event.key === 'ArrowRight') {
      event.preventDefault();
      return;
    }

    // Don't allow proceeding if there's an uncorrected error
    if (mustCorrectError) return;

    if (currentIndex >= sampleText.length) return;
    
    const typed = event.key;
    
    if (typed.length === 1) {
      if (typed === sampleText[currentIndex]) {
        // Correct character
        characters[currentIndex].status = 'correct';
        if (characters[currentIndex].hasRedHighlight) {
          characters[currentIndex].status = 'corrected';
        }

        // Build current word
        if (sampleText[currentIndex] === ' ' || currentIndex === sampleText.length - 1) {
          // Get the completed word
          let completedWord = sampleText.slice(wordStart, 
            currentIndex === sampleText.length - 1 ? currentIndex + 1 : currentIndex);
          sendWordProgress(completedWord.trim());
          wordStart = currentIndex + 1;
        }

        currentIndex++;
        mustCorrectError = false;
      } else {
        // Wrong character
        characters[currentIndex].status = 'incorrect';
        characters[currentIndex].hasRedHighlight = true;
        errorPositions.add(currentIndex);
        mustCorrectError = true;
      }
      characters = characters;
    }
  }

  function resetGame() {
    currentIndex = 0;
    mustCorrectError = false;
    errorPositions.clear();
    wordStart = 0;
    currentWord = '';
    characters = sampleText.split('').map(char => ({
      char,
      status: 'pending',
      hasRedHighlight: false
    }));
  }

  let containerRef;
</script>


<div class="min-h-screen bg-gray-900 text-gray-200 p-8">

  <button on:click={startGame} class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition">Start New Game</button>
  <div 
    bind:this={containerRef}
    class="max-w-2xl mx-auto p-6 bg-gray-800 rounded-lg shadow-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
    tabindex="0"
    on:keydown={handleKeyPress}
  >

        <div class="text-2xl leading-relaxed font-mono">
          {#each characters as { char, status, hasRedHighlight }, i}
            <span 
              class="
                relative inline-block min-w-[1ch] 
                {status === 'correct' ? 'text-green-400' :
                status === 'incorrect' ? 'text-red-400' :
                status === 'corrected' ? 'text-green-400' :
                'text-gray-400'}
                {i === currentIndex ? 'border-b-2 border-blue-500' : ''}
                {(status === 'incorrect' || hasRedHighlight) ? 'error-highlight' : ''}
              "
            >
              {char === ' ' ? '\u00A0' : char}
            </span>
          {/each}
        </div>
        
        <div class="mt-8 space-y-4">
          <p class="text-gray-400">
            Progress: {currentIndex} / {characters.length} characters
          </p>
          <p class="text-gray-400">
            Mistakes: {errorPositions.size}
          </p>
          
          <button
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
            on:click={resetGame}
          >
            Reset
          </button>
        </div>
  </div>
</div>
<style>
  span {
    white-space: pre;
  }

  .error-highlight {
    outline: 1px solid rgb(248 113 113);
    outline-offset: -1px;
    border-radius: 2px;
  }

  .error-highlight:empty::before {
    content: ' ';
  }
</style>
