<script lang="ts">
    import { UserRoundIcon } from "lucide-svelte"
    import * as Avatar from "$lib/components/ui/avatar/index.js";

    export let title: string = "GoFit";
</script>
<header class="flex-none sticky top-0 w-full border-b backdrop-blur supports-[backdrop-filter]:bg-background/60">
    <div class="flex container h-14 max-w-screen-2xl items-center justify-between">
        <h1 class="croll-m-20 text-l font-semibold tracking-tight lg:text-l">{title}</h1>
        <div class="flex items-center gap-4 text-sm lg:gap-6">
            <a href="/profile">
                <Avatar.Root class="items-center h-8 w-8 border">
                    <Avatar.Image src="https://avatars.githubusercontent.com/u/59253720?s=400&u=6f8fcb11a70d4fe4c3a6aa8a4a7c7f0d8772d4e0&v=4" alt="@shadcn" />
                    <Avatar.Fallback>
                        <UserRoundIcon class="size-4"/>
                    </Avatar.Fallback>
                </Avatar.Root>
            </a>
        </div>
    </div>
</header>