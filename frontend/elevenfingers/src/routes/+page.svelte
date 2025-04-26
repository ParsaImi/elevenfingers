<script lang="ts">
  import { goto } from '$app/navigation';
  
  function startGame() {
    goto('/rooms');
  }

  function gototest() {
      goto("/test")
  }
</script>

<svelte:head>
  <title>Tenfinger Typing Game</title>
  <meta name="description" content="Multiplayer typing game" />
</svelte:head>

<div class="container">
  <div class="hero">
    <h1>Tenfinger Typing Game</h1>
    <p class="subtitle">Compete with players from around the world in real-time typing challenges</p>
    
    <div class="cta-buttons">
      <button class="start-button" on:click={startGame}>Start Game</button>
      <a href="/login" class="login-link">Login / Register</a>
    </div>
  </div>
  
  <div class="features">
    <div class="feature">
      <h3>Real-Time Competition</h3>
      <p>Race against other players to see who can type the fastest with perfect accuracy</p>
    </div>
    
    <div class="feature">
      <h3>Improve Your Skills</h3>
      <p>Track your progress and see your typing speed and accuracy improve over time</p>
    </div>
    
    <div class="feature">
      <h3>Multiple Rooms</h3>
      <p>Join existing game rooms or create your own custom room to challenge friends</p>
    </div>
  </div>
</div>

<style>
  .container {
    max-width: 1000px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .hero {
    text-align: center;
    padding: 3rem 1rem;
    margin-bottom: 3rem;
  }
  
  h1 {
    font-size: 3.5rem;
    margin-bottom: 1rem;
    color: #333;
  }
  
  .subtitle {
    font-size: 1.25rem;
    color: #666;
    max-width: 600px;
    margin: 0 auto 2.5rem auto;
  }
  
  .cta-buttons {
    display: flex;
    justify-content: center;
    gap: 1.5rem;
    margin-bottom: 2rem;
    flex-wrap: wrap;
  }
  
  .start-button {
    padding: 0.875rem 2rem;
    font-size: 1.125rem;
    font-weight: 600;
    color: white;
    background-color: #4a56e2;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s, transform 0.2s;
  }
  
  .start-button:hover {
    background-color: #3643cf;
    transform: translateY(-2px);
  }
  
  .login-link {
    padding: 0.875rem 2rem;
    font-size: 1.125rem;
    font-weight: 600;
    color: #4a56e2;
    background-color: transparent;
    border: 2px solid #4a56e2;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s, color 0.2s;
    text-decoration: none;
    display: inline-block;
  }
  
  .login-link:hover {
    background-color: #4a56e2;
    color: white;
    text-decoration: none;
  }
  
  .features {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 2rem;
    margin-bottom: 3rem;
  }
  
  .feature {
    background-color: #f9f9f9;
    padding: 2rem;
    border-radius: 8px;
    text-align: center;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s;
  }
  
  .feature:hover {
    transform: translateY(-5px);
  }
  
  .feature h3 {
    margin-top: 0;
    margin-bottom: 1rem;
    color: #333;
  }
  
  .feature p {
    color: #666;
    margin: 0;
  }
</style>
