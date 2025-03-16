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
    ws = new WebSocket('ws://localhost:9000/ws'); // Replace with your actual WebSocket URL
    
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
