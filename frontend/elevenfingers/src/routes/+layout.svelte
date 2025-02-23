<script lang="ts">
    import { isAuthenticated } from '$lib/stores/auth';
    import '../app.css';
</script>

<div class="min-h-screen bg-gray-100">
    <nav class="bg-black shadow-lg">
        <div class="max-w-6xl mx-auto px-4">
            <div class="flex justify-between">
                <div class="flex space-x-7">
                    <a href="/" class="flex items-center py-4">
                        <span class="font-semibold text-gray-500 text-lg">
                            Typing Competition
                        </span>
                    </a>
                </div>
                <div class="flex items-center space-x-3">
                    {#if $isAuthenticated}
                        <a href="/game" class="py-2 px-4 text-gray-500 hover:text-gray-700">Game</a>
                        <button 
                            on:click={() => ($isAuthenticated = false)}
                            class="py-2 px-4 bg-red-500 text-white rounded hover:bg-red-600"
                        >
                            Logout
                        </button>
                    {:else}
                        <a 
                            href="/login"
                            class="py-2 px-4 bg-blue-500 text-white rounded hover:bg-blue-600"
                        >
                            Login
                        </a>
                    {/if}
                </div>
            </div>
        </div>
    </nav>
    
    <main class="bg-black max-w-6xl mx-auto px-4 py-8">
        <slot />
    </main>
</div>

