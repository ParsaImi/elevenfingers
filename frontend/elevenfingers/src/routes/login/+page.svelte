<script lang="ts">
    import { user, isAuthenticated } from '$lib/stores/auth';
    import { goto } from '$app/navigation';
    
    let username = '';
    let password = '';
    let error = '';

    async function handleLogin() {
        const data = {
              grant_type: 'password',
              username: username,
              password: password,
              scope: '',
              client_id: 'string',
              client_secret: 'string'
            };
        try {
            const response = await fetch('http://127.0.0.1:8000/auth/login', {
                method: "POST",
                headers: {
                    'accept': 'application/json',
                    'Content-Type': 'application/x-www-form-urlencoded'
                },
                body: new URLSearchParams(data).toString()
            });

             if (response.ok) {
                const data = await response.json();
                $user = data.access_token;
                console.log($user)
                $isAuthenticated = true;
                goto('/game');
            } else {
                error = 'Invalid credentials';
            }
        } catch (err) {
            error = 'Login failed';
        }
    }
</script>

<div class="max-w-md mx-auto bg-white p-8 rounded-lg shadow-md">
    <h2 class="text-2xl font-bold mb-6 text-center">Login</h2>
    
    {#if error}
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
        </div>
    {/if}
    
    <form on:submit|preventDefault={handleLogin} class="space-y-4">
        <div>
            <label for="username" class="block text-gray-700 mb-2">Username</label>
            <input
                id="username"
                type="text"
                bind:value={username}
                class="w-full p-2 border rounded"
                required
            />
        </div>
        
        <div>
            <label for="password" class="block text-gray-700 mb-2">Password</label>
            <input
                id="password"
                type="password"
                bind:value={password}
                class="w-full p-2 border rounded"
                required
            />
        </div>
        
        <button 
            type="submit"
            class="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
        >
            Login
        </button>
    </form>
</div>
