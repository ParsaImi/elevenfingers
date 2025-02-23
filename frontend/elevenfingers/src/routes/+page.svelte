<script lang="ts">
    import { isAuthenticated } from '$lib/stores/auth';
</script>

<div class="text-center">
    <h1 class="text-white text-4xl font-bold mb-8">Welcome to Typing Competition</h1>
    
    {#if $isAuthenticated}
        <a 
            href="/game" 
            class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600"
        >
            Start Playing
        </a>
    {:else}
        <a 
            href="/login" 
            class="bg-green-500 text-white px-6 py-3 rounded-lg hover:bg-green-600"
        >
            Login to Play
        </a>
    {/if}
</div>
