<script lang="ts">
    import BackButton from "../../back-button.svelte"
    import DeleteButton from "./delete-button.svelte"
    import * as Card from "$lib/components/ui/card/index.js";
    import CreateForm from "../../create/create-form.svelte";
    import type { PageData } from "./$types.js";
    export let data: PageData;
	import { onMount } from "svelte";
	import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
	import { toast } from "svelte-sonner";
    import { page } from '$app/stores';
	import { superForm, superValidate } from "sveltekit-superforms";
	import { zod } from "sveltekit-superforms/adapters";
	import { formSchema } from "../../create/schema";
    import HeaderComponent from "../../../header.svelte"
    
    let uuid: string = "";
    $: uuid = $page.params.uuid;

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