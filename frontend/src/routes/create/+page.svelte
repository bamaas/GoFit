<script lang="ts">
    import type { PageData } from "./$types.js";
    import CreateForm from "./create-form.svelte";
    import * as Card from "$lib/components/ui/card/index.js";
    export let data: PageData;
</script>


<Card.Root>
    <Card.Header>
      <Card.Title>Add check-in</Card.Title>
      <Card.Description>Try to daily check-in to get the best results</Card.Description>
    </Card.Header>
    <Card.Content>
        <CreateForm data={data.form} />
    </Card.Content>
</Card.Root>
