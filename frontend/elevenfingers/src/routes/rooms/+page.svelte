<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  
  let ws: WebSocket;
  let rooms: Record<string, Record<string, any>> = {};
  let loading = true;
  let error = null;
  
  onMount(() => {
    // Connect to WebSocket server
    ws = new WebSocket('wss://websocket.parsaimi.xyz/ws');
    
    ws.onopen = () => {
      console.log('Connected to the server');
      // Request the list of rooms
      requestRooms();
    };
    
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      
      switch(data.type) {
        case 'roomsStatus':
          rooms = data.rooms || {};
          loading = false;
          break;
        default:
          console.log('Received unknown message type:', data.type);
      }
    };
    
    ws.onerror = (err) => {
      console.error('WebSocket error:', err);
      error = 'Failed to connect to the server. Please try again later.';
      loading = false;
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
  
  function requestRooms() {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'roomsStatus'
      }));
    }
  }
  
  function refreshRooms() {
    loading = true;
    requestRooms();
  }
  
  function joinRoom(roomName: string) {
    // Store the room name to be used in the game page
    localStorage.setItem('selected_room', roomName);
    // Navigate to the game page
    goto('/game');
  }
  
  function createNewRoom() {
    // Open a modal or prompt for room name
    const roomName = prompt('Enter a new room name:');
    if (roomName && roomName.trim()) {
      localStorage.setItem('selected_room', roomName.trim());
      goto('/game');
    }
  }
</script>

<svelte:head>
  <title>Available Rooms - Tenfinger Typing Game</title>
  <meta name="description" content="Join a multiplayer typing game room" />
</svelte:head>

<div class="container">
  <div class="rooms-container">
    <h2>Available Rooms</h2>
    
    {#if loading}
      <div class="loading">
        <p>Loading available rooms...</p>
      </div>
    {:else if error}
      <div class="error">
        <p>{error}</p>
        <button on:click={refreshRooms}>Try Again</button>
      </div>
    {:else}
      <div class="room-controls">
        <button class="refresh-btn" on:click={refreshRooms}>
          Refresh Rooms
        </button>
        <button class="create-btn" on:click={createNewRoom}>
          Create New Room
        </button>
      </div>
      
      {#if Object.keys(rooms).length === 0}
        <div class="no-rooms">
          <p>No active rooms found. Create a new room to start playing!</p>
        </div>
      {:else}
        <div class="rooms-list">
          {#each Object.entries(rooms) as [roomName, players]}
            <div class="room-card">
              <div class="room-info">
                <h3>{roomName}</h3>
                <p>Players: {Object.keys(players).length}</p>
              </div>
              <button class="join-btn" on:click={() => joinRoom(roomName)}>
                Join Room
              </button>
            </div>
          {/each}
        </div>
      {/if}
    {/if}
  </div>
</div>

<style>
  .container {
    max-width: 1000px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .rooms-container {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  h2 {
    font-size: 2.5rem;
    margin-bottom: 2rem;
    text-align: center;
  }
  
  .loading, .error, .no-rooms {
    width: 100%;
    max-width: 600px;
    background-color: #f9f9f9;
    padding: 2rem;
    border-radius: 8px;
    text-align: center;
    margin-bottom: 2rem;
  }
  
  .error {
    background-color: #fff8f8;
    border: 1px solid #ffcdd2;
  }
  
  .loading p, .error p, .no-rooms p {
    margin-bottom: 1rem;
  }
  
  .room-controls {
    display: flex;
    gap: 1rem;
    margin-bottom: 2rem;
    width: 100%;
    max-width: 600px;
    justify-content: space-between;
  }
  
  .refresh-btn, .create-btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .refresh-btn {
    background-color: #f5f5f5;
    color: #333;
  }
  
  .refresh-btn:hover {
    background-color: #e0e0e0;
  }
  
  .create-btn {
    background-color: #4a56e2;
    color: white;
  }
  
  .create-btn:hover {
    background-color: #3643cf;
  }
  
  .rooms-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    width: 100%;
    max-width: 600px;
  }
  
  .room-card {
    background-color: #f9f9f9;
    border-radius: 8px;
    padding: 1.5rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s, box-shadow 0.2s;
  }
  
  .room-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  }
  
  .room-info h3 {
    margin: 0 0 0.5rem 0;
    font-size: 1.25rem;
  }
  
  .room-info p {
    margin: 0;
    color: #666;
  }
  
  .join-btn {
    padding: 0.5rem 1rem;
    background-color: #4a56e2;
    color: white;
    border: none;
    border-radius: 4px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .join-btn:hover {
    background-color: #3643cf;
  }
  
  button {
    cursor: pointer;
  }
</style>
