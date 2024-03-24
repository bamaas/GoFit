<script>
import { onMount } from "svelte";
import { apiData, entries } from '../store';

onMount(async () => {
  fetch("https://gofit-api.azurewebsites.net/entries/")
  .then(response => response.json())
  .then(data => {
		console.log(data);
    apiData.set(data);
  }).catch(error => {
    console.log(error);
    return [];
  });
});
</script>

<h1>Welcome to GoFit</h1>
<ul>

{#each $entries as entry}
		<li>{entry.weight}</li>
	{/each}
</ul>
