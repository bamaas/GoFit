<script lang="ts">
    import { UserRoundIcon } from "lucide-svelte"
    import * as Avatar from "$lib/components/ui/avatar/index.js";
    import { authenticated } from "$lib/stores/auth";
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import { Label } from "$lib/components/ui/label/index.js";
	import { request } from "$lib/functions/request";
	import { deleteCookie } from "$lib/functions/cookie";
	import { goto } from "$app/navigation";
	import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
    import { capitalizeFirstLetter } from "$lib/functions/string";
    import * as Select from "$lib/components/ui/select";

    let saveProfileChangesButtonDisabled: boolean = true;
    let profileDialogOpen: boolean = false;

    let iets: any = "";

    $: selectedGoal = iets ? {value: iets, label: capitalizeFirstLetter(iets)} : null;

    let user: Promise<any> = new Promise(() => {});

    function onOpen(){
        request(`${PUBLIC_BACKEND_BASE_URL}/v1/users/me`).then((response) => {
            user = response.data;
        })
    }

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

    function updateProfile(){
        request(`${PUBLIC_BACKEND_BASE_URL}/v1/users/me`, {
            method: "PUT",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
                goal: selectedGoal?.value
            }),
          }).then(() => {
              profileDialogOpen = false;
          });
    }


</script>

<Dialog.Root onOpenChange={onOpen} bind:open={profileDialogOpen}>
    <Dialog.Trigger on:click={() => (profileDialogOpen = true)}>
        {#if $authenticated}
        <Avatar.Root class="items-center h-8 w-8 border">
            <Avatar.Image src="https://avatars.githubusercontent.com/u/59253720?s=400&u=6f8fcb11a70d4fe4c3a6aa8a4a7c7f0d8772d4e0&v=4" alt="@shadcn" />
            <Avatar.Fallback>
                <UserRoundIcon class="size-4"/>
            </Avatar.Fallback>
        </Avatar.Root>
        {/if}
    </Dialog.Trigger>
    <Dialog.Content class="sm:max-w-[425px]">
    <Dialog.Header>
        <Dialog.Title>Your profile</Dialog.Title>
        <Dialog.Description>Hello, how are you doing?</Dialog.Description>
    </Dialog.Header>
    <div class="grid gap-4 py-4">
        {#await user then user}
        <div class="grid grid-cols-4 items-center gap-4">
            <Label for="email" class="text-right">Email</Label>
            <Input disabled id="email" value={user.email} class="col-span-3" />
        </div>
        <div class="grid grid-cols-4 items-center gap-4">
            <Label for="goal" class="text-right">Goal</Label>
            <div class="col-span-3">
                <Select.Root
                selected={selectedGoal}
                onSelectedChange={(v) => {
                v && (v.value);
                }}
                >
                <Select.Trigger class="w-full">
                    <Select.Value placeholder="Goal" />
                </Select.Trigger>
                <Select.Content>
                    <Select.Item value="cut">Cut</Select.Item>
                    <Select.Item value="bulk">Bulk</Select.Item>
                    <Select.Item value="maintain">Maintain</Select.Item>
                </Select.Content>
                </Select.Root>
            </div>
        </div>
        {/await}
    </div>
    <Dialog.Footer>
        <Button variant="outline" on:click={logout}>Logout</Button>
        <Button type="submit" on:click={updateProfile}>Save changes</Button>
    </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>