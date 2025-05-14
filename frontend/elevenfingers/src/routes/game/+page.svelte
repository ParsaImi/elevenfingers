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
  let typingGameComponent: TypingGame;
  let selectedRoom: string = '';
  let isPersianRoom = false; // Flag for Persian room
  
  onMount(() => {
    // Get the selected room from localStorage
    selectedRoom = localStorage.getItem('selected_room') || 'test';
    
    // Check if it's the Persian room
    isPersianRoom = selectedRoom === 'room3';
    
    // Connect to WebSocket server
    ws = new WebSocket('ws://localhost:9000/ws'); // Replace with your actual WebSocket URL
    
    ws.onopen = () => {
      console.log('Connected to the server');
      // Join the selected room when the page loads
      joinGame(selectedRoom);
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
        case 'playerRank':
          // Handle player rank updates
          if (typingGameComponent && data.playerrank) {
            typingGameComponent.updatePlayerRanks(data.playerrank);
          }
          break;
        case 'endGame':
          // Handle game end signal
          if (typingGameComponent) {
            typingGameComponent.endGame(data);
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
  
  function joinGame(roomName: string = selectedRoom) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'join',
        content: { room: roomName }
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
  
  function handlePlayAgain() {
    // Return to waiting room view
    currentView = 'waiting';
    
    // Rejoin the game room
    joinGame();
  }
  
  function changeRoom() {
    // Navigate back to the rooms list
    goto('/rooms');
  }
</script>

<svelte:head>
  <title>Tenfinger Typing Game - {selectedRoom} Room</title>
  <meta name="description" content="Multiplayer typing game" />
</svelte:head>

<div class="container {isPersianRoom ? 'persian-container' : ''}">
  {#if currentView === 'waiting'}
    <div class="room-header">
      <h2>Room: {selectedRoom} {isPersianRoom ? '(Persian)' : ''}</h2>
      <button class="change-room-btn" on:click={changeRoom}>Change Room</button>
    </div>
    <WaitingRoom 
      players={roomData.players} 
      onReady={sendReady}
      isPersian={isPersianRoom}
    />
  {:else if currentView === 'game'}
    <TypingGame 
      gameData={gameData} 
      onWordComplete={sendWordComplete}
      bind:this={typingGameComponent}
      on:playAgain={handlePlayAgain}
      isPersian={isPersianRoom}
    />
  {/if}
</div>

<style>
  .container {
    max-width: 1000px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .persian-container {
    /* Specific styles for Persian text direction */
    direction: rtl;
    text-align: right;
    font-family: 'Tahoma', 'Arial', sans-serif;
  }
  
  .room-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #eee;
  }
  
  .room-header h2 {
    margin: 0;
    font-size: 1.75rem;
  }
  
  .change-room-btn {
    padding: 0.5rem 1rem;
    background-color: #f5f5f5;
    border: none;
    border-radius: 4px;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .change-room-btn:hover {
    background-color: #e0e0e0;
  }
</style>
