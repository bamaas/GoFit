<script lang="ts">
    import type { PageData } from "./$types.js";
    import Form from "$lib/components/gofit/forms/LogbookForm.svelte";
    import * as Card from "$lib/components/ui/card/index.js";
    import BackButton from "../../../lib/components/gofit/buttons/PreviousPageButton.svelte"

    export let data: PageData;
</script>

<div class="container items-center py-4 max-w-screen-2xl">
    <div class="pb-3 pt-4">
        <BackButton />
    </div>
    <Card.Root>
        <Card.Header>
        <Card.Title>Add check-in</Card.Title>
        <Card.Description>Try to check-in daily to get the best results</Card.Description>
        </Card.Header>
        <Card.Content>
            <Form data={data.form} />
        </Card.Content>
    </Card.Root>
</div>