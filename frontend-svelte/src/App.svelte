<script>
	import { onMount } from 'svelte';
	import axios from 'axios';
  
	let djName = '';
	let result = null;
	let error = null;
  
	async function handleSubmit() {
	  try {
		const response = await axios.post('http://localhost:3001/check-dj', { djName });
		result = response.data;
		error = null;
	  } catch (err) {
		error = 'An error occurred. Please try again.';
		result = null;
	  }
	}
  </script>
  
  <main>
	<h1>Berghain DJ Checker</h1>
	<form on:submit|preventDefault={handleSubmit}>
	  <input 
		type="text" 
		bind:value={djName} 
		placeholder="Enter DJ name" 
		required 
	  />
	  <button type="submit">Check</button>
	</form>
	
	{#if error}
	  <p class="error">{error}</p>
	{/if}
	
	{#if result}
	  <p class="result">{result.message}</p>
	  {#if result.details}
		<p>{result.details}</p>
	  {/if}
	  {#if result.suggestion && result.suggestion !== djName}
		<p class="suggestion">Did you mean: {result.suggestion}?</p>
	  {/if}
	  {#if result.performances && result.performances.length > 0}
		<h2>Performance History</h2>
		<table>
		  <thead>
			<tr>
			  <th>Date</th>
			  <th>Time</th>
			  <th>Floor</th>
			  <th>Label</th>
			  <th>Closing Set</th>
			</tr>
		  </thead>
		  <tbody>
			{#each result.performances as perf}
			  <tr>
				<td>{perf.date}</td>
				<td>{perf.time}</td>
				<td>{perf.floor}</td>
				<td>{perf.label || 'N/A'}</td>
				<td>{perf.closing ? 'Yes' : 'No'}</td>
			  </tr>
			{/each}
		  </tbody>
		</table>
	  {/if}
	{/if}
  </main>
  
  <style>
	main {
	  font-family: Arial, sans-serif;
	  max-width: 800px;
	  margin: 0 auto;
	  padding: 20px;
	}
  
	form {
	  display: flex;
	  margin-bottom: 20px;
	}
  
	input[type="text"] {
	  flex-grow: 1;
	  padding: 10px;
	  font-size: 16px;
	}
  
	button {
	  padding: 10px 20px;
	  font-size: 16px;
	  background-color: #4CAF50;
	  color: white;
	  border: none;
	  cursor: pointer;
	}
  
	.result {
	  font-size: 18px;
	  font-weight: bold;
	  margin-bottom: 10px;
	}
  
	.suggestion {
	  margin-top: 10px;
	  font-style: italic;
	}
  
	table {
	  width: 100%;
	  border-collapse: collapse;
	  margin-top: 20px;
	}
  
	th, td {
	  border: 1px solid #ddd;
	  padding: 8px;
	  text-align: left;
	}
  
	th {
	  background-color: #f2f2f2;
	}
  
	tr:nth-child(even) {
	  background-color: #f9f9f9;
	}
  
	.error {
	  color: red;
	}
  </style>