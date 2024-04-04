<script lang="ts">
    import BackButton from "../../back-button.svelte"
    import DeleteButton from "./delete-button.svelte"
    import * as Card from "$lib/components/ui/card/index.js";
    import CreateForm from "../../create/create-form.svelte";
    import type { PageData } from "./$types.js";

    export let data: PageData;

</script>

<div class="container items-center py-4 max-w-screen-2xl">
    <div class="pb-3 pt-4 justify-between flex">
        <BackButton />
        <DeleteButton />
    </div>
    <Card.Root>
        <Card.Header>
        <Card.Title>Your check-in</Card.Title>
        <Card.Description>Every new day is another chance to change your life</Card.Description>
        </Card.Header>
        <Card.Content>
            <CreateForm data={data.form} />
        </Card.Content>
    </Card.Root>
</div>