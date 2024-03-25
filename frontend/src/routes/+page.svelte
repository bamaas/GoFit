<script lang="ts">
import { onMount } from "svelte";
import { apiData, entries } from '../store';
import * as Table from "$lib/components/ui/table";

onMount(async () => {
  fetch("https://gofit-api.azurewebsites.net/v1/check-ins")
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

<Table.Root>
	<Table.Header>
	  <Table.Row>
      <Table.Head class="w-[100px]">ID</Table.Head>
      <Table.Head class="text-right">Weight</Table.Head>
	  </Table.Row>
	</Table.Header>
	<Table.Body>
    {#each $entries as entry}
	  <Table.Row>
      <Table.Cell class="font-medium">{entry.id}</Table.Cell>
      <Table.Cell class="text-right">{entry.weight}</Table.Cell>
	  </Table.Row>
    {/each}
	</Table.Body>
  </Table.Root>
