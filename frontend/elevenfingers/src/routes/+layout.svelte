<script lang="ts">
  import '../app.css';
</script>

<div class="app">
  <header>
    <nav>
      <div class="logo">
        <a href="/">Tenfinger</a>
      </div>
      <div class="nav-links">
        <a href="/">Home</a>
        <a href="/game">Play</a>
      </div>
    </nav>
  </header>

  <main>
    <slot />
  </main>

  <footer>
    <p>© {new Date().getFullYear()} Tenfinger Typing Competition</p>
  </footer>
</div>

<style>
  .app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }
  
  header {
    background-color: #fff;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }
  
  nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem 2rem;
  }
  
  .logo a {
    font-size: 1.5rem;
    font-weight: bold;
    color: #4a56e2;
    text-decoration: none;
  }
  
  .nav-links {
    display: flex;
    gap: 2rem;
  }
  
  .nav-links a {
    color: #333;
    text-decoration: none;
    font-weight: 500;
  }
  
  .nav-links a:hover {
    color: #4a56e2;
  }
  
  main {
    flex: 1;
  }
  
  footer {
    background-color: #f5f5f5;
    text-align: center;
    padding: 1rem;
    margin-top: auto;
  }
</style>
