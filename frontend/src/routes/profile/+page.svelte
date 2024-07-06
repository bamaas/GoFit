<script lang="ts">
	import { goto } from "$app/navigation";
	import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
	import { Button } from "$lib/components/ui/button";
    import * as Card from "$lib/components/ui/card/index.js";
	import { deleteCookie } from "$lib/functions/cookie";
	import { request } from "$lib/functions/request";
    import { Separator } from "$lib/components/ui/separator/index.js";
	import { authenticated } from "$lib/stores/auth";
    import ProfileForm from "./(components)/profile-form.svelte";
    import type { PageData } from "./$types.js";
    export let data: PageData;

    function logout(){
        request(`${PUBLIC_BACKEND_BASE_URL}/v1/tokens/retract-all`, {method: 'DELETE'}, false)
        .then( (response: Response) => {
            if (response.ok){
                deleteCookie("token")
                authenticated.set(false);
                goto("/login");
            }
        });
    }

</script>

<div class="container items-center py-6 mt-4 max-w-screen-2xl">
    <Card.Root>
        <Card.Header>
        <Card.Title>Your profile</Card.Title>
        <Card.Description>Hello, how are you doing?</Card.Description>
        </Card.Header>
        <Card.Content>
            <ProfileForm data={data.form} />
            <Separator orientation="horizontal" class="mt-5 mb-5"/>
            <Button on:click={logout} variant="outline" class="w-full">Logout</Button>
        </Card.Content>
    </Card.Root>
</div>