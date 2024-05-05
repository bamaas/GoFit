<script lang="ts">
	import { goto } from "$app/navigation";
    import BackButton from "$lib/components/gofit/buttons/PreviousPageButton.svelte"
	import { Button } from "$lib/components/ui/button";
    import * as Card from "$lib/components/ui/card/index.js";
	import { deleteCookie } from "$lib/functions/cookie";

    function logout(){
        deleteCookie("token")
        goto("/login");
    }

</script>

<div class="container items-center py-4 max-w-screen-2xl">
    <div class="pb-3 pt-4 justify-between flex">
        <BackButton />
    </div>
    <Card.Root>
        <Card.Header>
        <Card.Title>Your profile</Card.Title>
        <Card.Description>Hello mynamehere</Card.Description>
        </Card.Header>
        <Card.Content>
            <Button on:click={logout} variant="outline">Logout</Button>
        </Card.Content>
    </Card.Root>
</div>