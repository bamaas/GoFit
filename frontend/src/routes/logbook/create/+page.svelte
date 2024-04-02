<script lang="ts">
    import type { PageData } from "./$types.js";
    import CreateForm from "./create-form.svelte";
    import * as Card from "$lib/components/ui/card/index.js";
    import BackButton from "../back-button.svelte"
    import HeaderComponent from "../../header.svelte"

    export let data: PageData;
</script>

<HeaderComponent />
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
            <CreateForm data={data.form} />
        </Card.Content>
    </Card.Root>
</div>