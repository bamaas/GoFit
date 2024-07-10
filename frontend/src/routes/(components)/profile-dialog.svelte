<script lang="ts">
    import { UserRoundIcon } from "lucide-svelte"
    import * as Avatar from "$lib/components/ui/avatar/index.js";
    import { authenticated } from "$lib/stores/auth";
    import * as Form from "$lib/components/ui/form";
    import LoaderCircleIcon from "lucide-svelte/icons/loader-circle";
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
    import { toast } from 'svelte-sonner';
    import { formSchema, type FormSchema } from "./schema";
    import { type SuperValidated, type Infer, superForm } from "sveltekit-superforms";
	import { zod } from "sveltekit-superforms/adapters";
    import { defaults } from 'sveltekit-superforms';
    import { Skeleton } from "$lib/components/ui/skeleton/index.js";
    import { profileStore } from "$lib/stores/profile";
    import type { UserProfile } from "$lib/stores/profile";
	import { fetchUserProfile } from "$lib/functions/profile";
	import { onMount } from "svelte";
    
    let profileDialogOpen: boolean = false;

    let user: Promise<UserProfile> = new Promise(() => {});

    function handleDialogStateChange(open: boolean){
        if (open) {
            fetchUserProfile();
        } else {
            buttonDisabled = true;
        }
    }

    onMount(() => {
        profileStore.subscribe((userProfile) => {
            if (!userProfile) return;
			user = new Promise((resolve) => {
				resolve(userProfile)
			});
            $formData.goal = userProfile.goal;
        });
    });

    function logout(){
        request(`${PUBLIC_BACKEND_BASE_URL}/v1/tokens/retract-all`, {method: 'DELETE'}, false)
        .then( (response: Response) => {
            if (response.ok){
                deleteCookie("token")
                authenticated.set(false);
                profileStore.set(undefined);
                goto("/login");
            }
        });
    }

    let changes: number = 0;
    let buttonDisabled: boolean = true;

    const { form: formData, enhance, delayed }  = superForm(defaults(zod(formSchema)), {
        SPA: true,
        validators: zod(formSchema),
        resetForm: false,
        delayMs: 1000,
        onChange: () => {
            changes++;
            if (changes === 1) return;
            try {
                formSchema.parse($formData)
                buttonDisabled = false;
                changes++;
            } catch (e: any) {
                buttonDisabled = true;
            }
        },
        onUpdate: async ({form}) => {
            try {
                const response = await request(`${PUBLIC_BACKEND_BASE_URL}/v1/users/me`, {
                    method: "PUT",
                    headers: {
                    "Content-Type": "application/json",
                    },
                    body: JSON.stringify(form.data),
                });
                const user: UserProfile = {
                    id: response.data.id,
                    email: response.data.email,
                    goal: response.data.goal,
                    activated: response.data.activated,
                    createdAt: response.data.created_at,
                }
                profileStore.set(user);
                profileDialogOpen = false;
                toast.info("Got ya!", {
                    description: "Profile saved.",
                    cancel: {label: "X", onClick: () => {}}
                });
            } catch (e) {
                showErrors(form);
            }
        }
    });

    function showErrors(form: SuperValidated<Infer<FormSchema>>): void {
        form.errors.goal = ["Invalid goal"];
        toast.error("Oops!", {
            description: "Something went wrong.",
            cancel: {label: "X", onClick: () => {}}
        });
    }

    $: selectedGoal = $formData.goal
    ? {
        label: capitalizeFirstLetter($formData.goal),
        value: $formData.goal,
      }
    : undefined;

    $: handleDialogStateChange(profileDialogOpen);

</script>

<Dialog.Root bind:open={profileDialogOpen}>
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
        <div class="grid grid-cols-4 items-center gap-4">
            <Label for="email" class="text-right">Email</Label>
            {#await user}
            <Skeleton class="h-10 col-span-3" />
            {:then user}
            <Input disabled id="email" value={user.email} class="col-span-3" />
            {/await}
        </div>
        <div class="grid grid-cols-4 items-center gap-4">
            <Label for="goal" class="text-right">Goal</Label>
            {#await user}
            <Skeleton class="h-10 col-span-3" />
            {:then user}
            <div class="col-span-3">
                <Select.Root
                selected={selectedGoal}
                onSelectedChange={(v) => {
                    v && ($formData.goal = v.value);
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
            {/await}
        </div>
    </div>
    <Dialog.Footer>
        <Button variant="outline" on:click={logout}>Logout</Button>
        <form method="POST" use:enhance>
            <Form.Button disabled={buttonDisabled || $delayed} class="sm:w-auto w-full mb-4 sm:mb-0">
                {#if $delayed}
                <LoaderCircleIcon class="spinner"/>
                {:else}
                Save profile
                {/if}
            </Form.Button>
        </form>
    </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>