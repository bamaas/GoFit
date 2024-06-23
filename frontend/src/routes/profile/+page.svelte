<script lang="ts">
	import { goto } from "$app/navigation";
	import { Button } from "$lib/components/ui/button";
    import * as Card from "$lib/components/ui/card/index.js";
	import { deleteCookie } from "$lib/functions/cookie";
	import { authenticated } from "$lib/stores/auth";

    function logout(){
        deleteCookie("token")
        authenticated.set(false);
        goto("/login");
    }

</script>

<div class="container items-center py-6 mt-4 max-w-screen-2xl">
    <Card.Root>
        <Card.Header>
        <Card.Title>Your profile</Card.Title>
        <Card.Description>Hello, how are you doing?</Card.Description>
        </Card.Header>
        <Card.Content>
            <Button on:click={logout} variant="outline">Logout</Button>
        </Card.Content>
    </Card.Root>
</div>