<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  
  let email = '';
  let username = '';
  let password = '';
  let confirmPassword = '';
  let isLoading = false;
  let error: string | null = null;
  let successMessage: string | null = null;
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
    // Reset states
    error = null;
    successMessage = null;
    
    // Validate form
    if (!email || !username || !password) {
      error = 'Please fill in all required fields';
      return;
    }
    
    // Validate email format
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      error = 'Please enter a valid email address';
      return;
    }
    
    // Validate password match
    if (password !== confirmPassword) {
      error = 'Passwords do not match';
      return;
    }
    
    // Set loading state
    isLoading = true;
    
    try {
      const response = await fetch('http://localhost:8000/auth/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email, username, password })
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.message || 'Registration failed');
      }
      
      // Show success message
      successMessage = 'Account created successfully! Redirecting to login...';
      
      // Redirect to login after a short delay
      setTimeout(() => {
        goto('/login');
      }, 2000);
      
    } catch (err) {
      console.error('Registration error:', err);
      error = err.message || 'Failed to create account. Please try again.';
    } finally {
      isLoading = false;
    }
  }
  
  function togglePasswordVisibility() {
    showPassword = !showPassword;
  }
  
  function goToLogin() {
    goto('/login');
  }
</script>

<svelte:head>
  <title>Sign Up - Tenfinger Typing Game</title>
  <meta name="description" content="Create your Tenfinger typing game account" />
</svelte:head>

<div class="container">
  <div class="signup-container">
    <h2>Create an Account</h2>
    
    {#if error}
      <div class="error-message">
        {error}
      </div>
    {/if}
    
    {#if successMessage}
      <div class="success-message">
        {successMessage}
      </div>
    {/if}
    
    <form on:submit|preventDefault={handleSubmit}>
      <div class="form-group">
        <label for="email">Email <span class="required">*</span></label>
        <input 
          type="email" 
          id="email" 
          bind:value={email} 
          placeholder="Enter your email address"
          autocomplete="email"
          disabled={isLoading}
          required
        />
      </div>
      
      <div class="form-group">
        <label for="username">Username <span class="required">*</span></label>
        <input 
          type="text" 
          id="username" 
          bind:value={username} 
          placeholder="Choose a username"
          autocomplete="username"
          disabled={isLoading}
          required
        />
      </div>
      
      <div class="form-group">
        <label for="password">Password <span class="required">*</span></label>
        <div class="password-input-container">
          <input 
            type={showPassword ? 'text' : 'password'} 
            id="password" 
            bind:value={password} 
            placeholder="Create a password"
            autocomplete="new-password"
            disabled={isLoading}
            required
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
      
      <div class="form-group">
        <label for="confirmPassword">Confirm Password <span class="required">*</span></label>
        <div class="password-input-container">
          <input 
            type={showPassword ? 'text' : 'password'} 
            id="confirmPassword" 
            bind:value={confirmPassword} 
            placeholder="Confirm your password"
            autocomplete="new-password"
            disabled={isLoading}
            required
          />
        </div>
      </div>
      
      <div class="terms-privacy">
        <p>By creating an account, you agree to our <a href="/terms">Terms of Service</a> and <a href="/privacy">Privacy Policy</a>.</p>
      </div>
      
      <button 
        type="submit" 
        class="signup-button"
        disabled={isLoading}
      >
        {isLoading ? 'Creating Account...' : 'Create Account'}
      </button>
    </form>
    
    <div class="login-prompt">
      <p>Already have an account? <a href="/login" on:click|preventDefault={goToLogin}>Login</a></p>
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
  
  .signup-container {
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
  
  .required {
    color: #d32f2f;
  }
  
  input[type="text"],
  input[type="email"],
  input[type="password"] {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
    transition: border-color 0.2s;
  }
  
  input[type="text"]:focus,
  input[type="email"]:focus,
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
  
  .terms-privacy {
    margin-bottom: 1.5rem;
    font-size: 0.85rem;
    color: #666;
    text-align: center;
  }
  
  .terms-privacy a {
    color: #4a56e2;
  }
  
  .signup-button {
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
  
  .signup-button:hover {
    background-color: #3643cf;
  }
  
  .signup-button:disabled {
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
  
  .success-message {
    background-color: #e8f5e9;
    border: 1px solid #a5d6a7;
    color: #2e7d32;
    padding: 0.75rem 1rem;
    border-radius: 4px;
    margin-bottom: 1.5rem;
    font-size: 0.9rem;
  }
  
  .login-prompt {
    text-align: center;
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid #eee;
  }
  
  .login-prompt a {
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
