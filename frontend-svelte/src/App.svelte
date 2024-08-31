<script>
	import { onMount } from 'svelte';
	import axios from 'axios';
  
	let djName = '';
	let result = null;
	let error = null;
	let loading = false;
  
	async function handleSubmit() {
	  loading = true;
	  try {
		const response = await axios.post('http://localhost:3001/check-dj', { djName });
		result = response.data;
		error = null;
	  } catch (err) {
		error = 'An error occurred. Please try again.';
		result = null;
	  } finally {
		loading = false;
	  }
	}
  </script>
  
  <main class="min-h-screen bg-black text-gray-300 flex items-center justify-center px-4 py-12">
	<div class="max-w-2xl w-full">
	  <h1 class="text-4xl font-bold mb-12 text-center text-white tracking-widest">BERGHAIN DJ CHECKER</h1>
	  <form on:submit|preventDefault={handleSubmit} class="flex mb-8">
		<input 
		  type="text" 
		  bind:value={djName} 
		  placeholder="Enter DJ name" 
		  required 
		  class="flex-grow p-3 bg-gray-900 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-gray-700 border-b-2 border-gray-700"
		/>
		<button type="submit" disabled={loading} class="px-6 py-3 bg-gray-800 text-white font-semibold hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-700 disabled:bg-gray-900 disabled:text-gray-600 disabled:cursor-not-allowed transition duration-150 ease-in-out">
		  {loading ? 'Checking...' : 'Check'}
		</button>
	  </form>
	  
	  {#if error}
		<p class="text-red-500 mt-4">{error}</p>
	  {/if}
	  
	  {#if result}
		<div class="mt-8 text-gray-300">
		  <p class="text-xl mb-4">{result.message}</p>
		  {#if result.details}
			<p class="mb-2 text-gray-400">{result.details}</p>
		  {/if}
		  {#if result.suggestion && result.suggestion !== djName}
			<p class="italic mb-4 text-gray-400">Did you mean: {result.suggestion}?</p>
		  {/if}
		  {#if result.performances && result.performances.length > 0}
			<h2 class="text-2xl font-semibold mt-12 mb-6 text-white">Performance History</h2>
			<div class="overflow-x-auto">
			  <table class="w-full text-left">
				<thead>
				  <tr class="border-b border-gray-700">
					<th class="p-3 text-gray-400">Date</th>
					<th class="p-3 text-gray-400">Time</th>
					<th class="p-3 text-gray-400">Floor</th>
					<th class="p-3 text-gray-400">Label</th>
					<th class="p-3 text-gray-400">Closing Set</th>
				  </tr>
				</thead>
				<tbody>
				  {#each result.performances as perf}
					<tr class="border-b border-gray-800">
					  <td class="p-3">{perf.date}</td>
					  <td class="p-3">{perf.time}</td>
					  <td class="p-3">{perf.floor}</td>
					  <td class="p-3">{perf.label || 'N/A'}</td>
					  <td class="p-3">{perf.closing ? 'Yes' : 'No'}</td>
					</tr>
				  {/each}
				</tbody>
			  </table>
			</div>
		  {/if}
		</div>
	  {/if}
	</div>
  </main>
  
  <style global>
	@import url('https://fonts.googleapis.com/css2?family=Space+Mono&display=swap');
  
	body {
	  font-family: 'Space Mono', monospace;
	  background-color: black;
	}
  </style>