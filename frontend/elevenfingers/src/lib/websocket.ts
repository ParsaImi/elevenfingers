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
        case 'startGame':
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
