<script lang="ts">
import { onMount } from "svelte";
import { apiData, entries } from '../store';
import * as Table from "$lib/components/ui/table";

onMount(async () => {
  fetch("http://localhost:8080/v1/check-ins")
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
        <Table.Head class="w-[100px]">Date</Table.Head>
        <Table.Head class="text-right">Weight</Table.Head>
      </Table.Row>  
	</Table.Header>
	<Table.Body>
    {#each $entries as entry}
	  <Table.Row>
      <Table.Cell class="w-[100px]">{entry.created_at}</Table.Cell>
      <Table.Cell class="text-right">{entry.weight}</Table.Cell>
	  </Table.Row>
    {/each}
	</Table.Body>
  </Table.Root>
