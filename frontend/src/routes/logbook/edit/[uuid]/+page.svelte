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
    
    let uuid = "";
    $: uuid = $page.params.uuid;

    // onMount(async () => {
    //     fetch(`${PUBLIC_BACKEND_BASE_URL}/v1/check-ins/${uuid}`)
    //     .then(response => response.json())
    //     .then(data => {
    //         let d = data;
    //         console.log(d)
    //         form: superValidate(data, zod(formSchema));
    //     }).catch(error => {
    //         toast.error("Oops! Failed fetching data from server.");
    //         console.log(error);
    //         return [];
    //     });
    //     // return {form: await superValidate(d, zod(formSchema))}
    // });

</script>

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