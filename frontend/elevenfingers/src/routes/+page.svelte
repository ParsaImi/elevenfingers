<script lang="ts">
  import { goto } from '$app/navigation';

  function startGame() {
    goto('/game');
  }
</script>

<svelte:head>
  <title>Tenfinger - Typing Competition</title>
  <meta name="description" content="Test your typing speed and accuracy" />
</svelte:head>

<div class="hero">
  <h1>Tenfinger Typing Competition</h1>
  <p>Challenge your friends and improve your typing skills</p>
  
  <div class="cta">
    <button on:click={startGame}>Start Game</button>
  </div>
</div>

<style>
  .hero {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 80vh;
    text-align: center;
    padding: 0 1rem;
  }
  
  h1 {
    font-size: 3rem;
    margin-bottom: 1rem;
    color: #333;
  }
  
  p {
    font-size: 1.2rem;
    margin-bottom: 2rem;
    color: #666;
  }
  
  .cta {
    margin-top: 2rem;
  }
  
  button {
    padding: 1rem 2rem;
    font-size: 1.2rem;
    font-weight: bold;
    color: white;
    background-color: #4a56e2;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  button:hover {
    background-color: #3643cf;
  }
</style>
