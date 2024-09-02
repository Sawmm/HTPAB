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
		const response = await axios.post('http://localhost:3001/check-dj', { djName }, {
		  headers: {
			'Content-Type': 'application/json',
			// You might need to add more headers here, depending on your server requirements
		  },
		//   withCredentials: true, // Include this if your server expects credentials
		});
		result = response.data;
		error = null;
		console.log('Success:', result);
	  } catch (err) {
		console.error('Error details:', err);
		error = `An error occurred: ${err.message}. ${err.response ? `Status: ${err.response.status}, Data: ${JSON.stringify(err.response.data)}` : ''}`;
		result = null;
	  } finally {
		loading = false;
	  }
	}

	onMount(() => {
	  // Test CORS preflight
	  fetch('http://localhost:3001/check-dj', {
		method: 'OPTIONS',
		headers: {
		  'Origin': window.location.origin,
		  'Access-Control-Request-Method': 'POST',
		  'Access-Control-Request-Headers': 'Content-Type',
		},
	  }).then(response => {
		console.log('CORS preflight response:', response);
	  }).catch(error => {
		console.error('CORS preflight error:', error);
	  });
	});
  </script>
  
  <main class="min-h-screen flex items-center justify-center px-4 py-12">
	<div class="max-w-4xl w-full bg-gray-800 rounded-lg shadow-xl p-8">
	  <h1 class="text-5xl font-bold mb-12 text-center text-white tracking-widest">BERGHAIN DJ CHECKER</h1>
	  <form on:submit|preventDefault={handleSubmit} class="flex mb-8">
		<input 
		  type="text" 
		  bind:value={djName} 
		  placeholder="Enter DJ name" 
		  required 
		  class="flex-grow p-4 bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 border-b-2 border-gray-600 rounded-l-lg"
		/>
		<button type="submit" disabled={loading} class="px-8 py-4 bg-purple-600 text-white font-semibold hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500 disabled:bg-gray-600 disabled:text-gray-400 disabled:cursor-not-allowed transition duration-150 ease-in-out rounded-r-lg">
		  {loading ? 'Checking...' : 'Check'}
		</button>
	  </form>
	  
	  {#if error}
		<p class="text-red-400 mt-4 p-4 bg-red-900 bg-opacity-50 rounded">{error}</p>
	  {/if}
	  
	  {#if result}
		<div class="mt-8">
		  <p class="text-2xl mb-4 text-purple-300">{result.message}</p>
		  {#if result.details}
			<p class="mb-6 text-gray-400">{result.details}</p>
		  {/if}
		  {#if result.performances && result.performances.length > 0}
			<h2 class="text-3xl font-semibold mt-12 mb-6 text-white">Performance History</h2>
			<div class="overflow-x-auto bg-gray-700 rounded-lg">
			  <table class="w-full text-left">
				<thead>
				  <tr class="border-b border-gray-600">
					<th class="p-4 text-purple-300">Date</th>
					<th class="p-4 text-purple-300">Time</th>
					<th class="p-4 text-purple-300">Floor</th>
					<th class="p-4 text-purple-300">Label</th>
					<th class="p-4 text-purple-300">Closing Set</th>
				  </tr>
				</thead>
				<tbody>
				  {#each result.performances as perf, i}
					<tr class="border-b border-gray-600 {i % 2 === 0 ? 'bg-gray-750' : 'bg-gray-700'}">
					  <td class="p-4">{perf.Date}</td>
					  <td class="p-4">{perf.Time}</td>
					  <td class="p-4">{perf.Floor}</td>
					  <td class="p-4">{perf.Label || 'N/A'}</td>
					  <td class="p-4">{perf.Closing ? 'Yes' : 'No'}</td>
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