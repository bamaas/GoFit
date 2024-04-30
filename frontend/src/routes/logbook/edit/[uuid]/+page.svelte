<script lang="ts">
    import BackButton from "../../../../lib/components/gofit/buttons/PreviousPageButton.svelte"
    import DeleteButton from "./delete-button.svelte"
    import * as Card from "$lib/components/ui/card/index.js";
    import Form from "$lib/components/gofit/forms/LogbookForm.svelte";
    import type { PageData } from "./$types.js";

    export let data: PageData;

</script>

<div class="container items-center py-4 max-w-screen-2xl">
    <div class="pb-3 pt-4 justify-between flex">
        <BackButton />
        <DeleteButton authToken={data.authToken}/>
    </div>
    <Card.Root>
        <Card.Header>
        <Card.Title>Your check-in</Card.Title>
        <Card.Description>Every new day is another chance to change your life</Card.Description>
        </Card.Header>
        <Card.Content>
            <Form data={data.form} authToken={data.authToken}/>
        </Card.Content>
    </Card.Root>
</div>