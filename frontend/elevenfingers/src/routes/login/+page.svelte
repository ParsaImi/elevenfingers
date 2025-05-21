<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  
  let username = '';
  let password = '';
  let rememberMe = false;
  let isLoading = false;
  let error: string | null = null;
  let showPassword = false;
  
  // Check if user is already logged in
  onMount(() => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      // Redirect to the main page if user is already authenticated
      goto('/');
    }
  });
  
  async function handleSubmit() {
    // Reset error state
    error = null;
    
    // Validate form
    if (!username || !password) {
      error = 'Please enter both username and password';
      return;
    }
    
    // Set loading state
    isLoading = true;
    
    try {
      const response = await fetch('http://localhost:8000/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.message || 'Login failed');
      }
      
      // Store auth token in localStorage
      if (data.access_token) {
        localStorage.setItem('auth_token', data.access_token);
      }
      
      // Redirect to rooms page
      goto('/rooms');
    } catch (err) {
      console.error('Login error:', err);
      error = err.message || 'Failed to login. Please try again.';
    } finally {
      isLoading = false;
    }
  }
  
  function togglePasswordVisibility() {
    showPassword = !showPassword;
  }
  
  function goToRegister() {
    goto('/register');
  }
</script>

<svelte:head>
  <title>Login - Tenfinger Typing Game</title>
  <meta name="description" content="Login to your Tenfinger typing game account" />
</svelte:head>

<div class="container">
  <div class="login-container">
    <h2>Login to Your Account</h2>
    
    {#if error}
      <div class="error-message">
        {error}
      </div>
    {/if}
    
    <form on:submit|preventDefault={handleSubmit}>
      <div class="form-group">
        <label for="username">Username or Email</label>
        <input 
          type="text" 
          id="username" 
          bind:value={username} 
          placeholder="Enter your username or email"
          autocomplete="username"
          disabled={isLoading}
        />
      </div>
      
      <div class="form-group">
        <label for="password">Password</label>
        <div class="password-input-container">
          <input 
            type={showPassword ? 'text' : 'password'} 
            id="password" 
            bind:value={password} 
            placeholder="Enter your password"
            autocomplete="current-password"
            disabled={isLoading}
          />
          <button 
            type="button" 
            class="toggle-password" 
            on:click={togglePasswordVisibility}
            disabled={isLoading}
          >
            {showPassword ? 'Hide' : 'Show'}
          </button>
        </div>
      </div>
      
      <div class="form-options">
        <div class="remember-me">
          <input 
            type="checkbox" 
            id="rememberMe" 
            bind:checked={rememberMe}
            disabled={isLoading}
          />
          <label for="rememberMe">Remember me</label>
        </div>
        
        <a href="/forgot-password" class="forgot-password">Forgot password?</a>
      </div>
      
      <button 
        type="submit" 
        class="login-button"
        disabled={isLoading}
      >
        {isLoading ? 'Logging in...' : 'Login'}
      </button>
    </form>
    
    <div class="register-prompt">
      <p>Don't have an account? <a href="/register" on:click|preventDefault={goToRegister}>Register now</a></p>
    </div>
    
    <div class="back-to-home">
      <a href="/">Back to Home</a>
    </div>
  </div>
</div>

<style>
  .container {
    max-width: 1000px;
    margin: 0 auto;
    padding: 2rem;
    display: flex;
    justify-content: center;
    min-height: calc(100vh - 4rem);
  }
  
  .login-container {
    width: 100%;
    max-width: 500px;
    background-color: #f9f9f9;
    padding: 2.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  }
  
  h2 {
    font-size: 2rem;
    margin-bottom: 1.5rem;
    text-align: center;
    color: #333;
  }
  
  .form-group {
    margin-bottom: 1.5rem;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: #444;
  }
  
  input[type="text"],
  input[type="password"] {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
    transition: border-color 0.2s;
  }
  
  input[type="text"]:focus,
  input[type="password"]:focus {
    border-color: #4a56e2;
    outline: none;
    box-shadow: 0 0 0 2px rgba(74, 86, 226, 0.2);
  }
  
  .password-input-container {
    position: relative;
  }
  
  .toggle-password {
    position: absolute;
    right: 10px;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    color: #666;
    font-size: 0.85rem;
    cursor: pointer;
    padding: 0.25rem 0.5rem;
  }
  
  .toggle-password:hover {
    color: #4a56e2;
  }
  
  .form-options {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    font-size: 0.9rem;
  }
  
  .remember-me {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  
  .forgot-password {
    color: #4a56e2;
  }
  
  .login-button {
    width: 100%;
    padding: 0.875rem;
    background-color: #4a56e2;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s, transform 0.2s;
  }
  
  .login-button:hover {
    background-color: #3643cf;
  }
  
  .login-button:disabled {
    background-color: #a8aef0;
    cursor: not-allowed;
  }
  
  .error-message {
    background-color: #fff8f8;
    border: 1px solid #ffcdd2;
    color: #d32f2f;
    padding: 0.75rem 1rem;
    border-radius: 4px;
    margin-bottom: 1.5rem;
    font-size: 0.9rem;
  }
  
  .register-prompt {
    text-align: center;
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid #eee;
  }
  
  .register-prompt a {
    color: #4a56e2;
    font-weight: 500;
  }
  
  .back-to-home {
    text-align: center;
    margin-top: 1.5rem;
    font-size: 0.9rem;
  }
  
  .back-to-home a {
    color: #666;
  }
  
  .back-to-home a:hover {
    color: #4a56e2;
  }
</style>
