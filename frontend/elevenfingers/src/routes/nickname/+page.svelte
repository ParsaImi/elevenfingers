<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  
  let nickname = '';
  let error = '';
  let destination = '';
  
  onMount(() => {
    // Check if user already has a nickname
    const existingNickname = localStorage.getItem('user_nickname');
    
    if (existingNickname) {
      nickname = existingNickname;
    }
    
    // Get the destination from URL query parameter or default to /rooms
    const url = new URL(window.location.href);
    destination = url.searchParams.get('redirect') || '/rooms';
  });
  
  function handleSubmit() {
    if (!nickname || nickname.trim().length < 3) {
      error = 'Please enter a nickname (at least 3 characters)';
      return;
    }
    
    if (nickname.trim().length > 20) {
      error = 'Nickname must be 20 characters or less';
      return;
    }
    
    // Store nickname in localStorage
    localStorage.setItem('user_nickname', nickname.trim());
    
    // Navigate to destination
    goto(destination);
  }
</script>

<svelte:head>
  <title>Enter Your Nickname - Tenfinger Typing Game</title>
  <meta name="description" content="Set your nickname for the multiplayer typing game" />
</svelte:head>

<div class="container">
  <div class="nickname-form">
    <h1>Welcome to Tenfinger Typing Game</h1>
    <p class="subtitle">Enter a nickname to continue</p>
    
    <form on:submit|preventDefault={handleSubmit}>
      <div class="form-group">
        <label for="nickname">Your Nickname</label>
        <input 
          type="text" 
          id="nickname" 
          bind:value={nickname} 
          placeholder="Enter a nickname (3-20 characters)"
          class:error={error}
          autocomplete="off"
        />
        {#if error}
          <p class="error-message">{error}</p>
        {/if}
      </div>
      
      <button type="submit" class="continue-btn">Continue</button>
    </form>
  </div>
</div>

<style>
  .container {
    max-width: 600px;
    margin: 0 auto;
    padding: 2rem;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .nickname-form {
    width: 100%;
    background-color: #f9f9f9;
    padding: 2.5rem;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    text-align: center;
  }
  
  h1 {
    font-size: 2.25rem;
    margin-bottom: 0.75rem;
    color: #333;
  }
  
  .subtitle {
    font-size: 1.1rem;
    color: #666;
    margin-bottom: 2rem;
  }
  
  .form-group {
    text-align: left;
    margin-bottom: 1.5rem;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
  }
  
  input {
    width: 100%;
    padding: 0.875rem;
    font-size: 1rem;
    border: 2px solid #e0e0e0;
    border-radius: 4px;
    transition: border-color 0.2s;
  }
  
  input:focus {
    outline: none;
    border-color: #4a56e2;
  }
  
  input.error {
    border-color: #e74c3c;
  }
  
  .error-message {
    color: #e74c3c;
    margin-top: 0.5rem;
    font-size: 0.875rem;
  }
  
  .continue-btn {
    width: 100%;
    padding: 0.875rem;
    font-size: 1.125rem;
    font-weight: 600;
    color: white;
    background-color: #4a56e2;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s, transform 0.2s;
  }
  
  .continue-btn:hover {
    background-color: #3643cf;
    transform: translateY(-2px);
  }
</style>
