<script lang="ts">
    import BackButton from "../back-button.svelte"
    import DeleteButton from "./delete-button.svelte"
    import * as Card from "$lib/components/ui/card/index.js";
    let uuid;
</script>

<div class="pb-3 pt-4 justify-between flex">
    <BackButton />
    <DeleteButton />
</div>
<Card.Root>
    <Card.Header>
      <Card.Title>Check-in</Card.Title>
      <Card.Description>Try to check-in daily to get the best results</Card.Description>
    </Card.Header>
    <Card.Content>
        Lorem ipsum
    </Card.Content>
</Card.Root>