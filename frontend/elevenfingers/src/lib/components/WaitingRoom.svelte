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

