<script lang="ts">
    import { CircleUserRoundIcon } from "lucide-svelte"

    export let title: string = "GoFit";
</script>
<header class="flex-none sticky top-0 w-full border-b backdrop-blur supports-[backdrop-filter]:bg-background/60">
    <div class="flex container h-14 max-w-screen-2xl items-center justify-between">
        <h1 class="croll-m-20 text-l font-semibold tracking-tight lg:text-l">{title}</h1>
        <div class="flex items-center gap-4 text-sm lg:gap-6">
            <CircleUserRoundIcon class="size-5 transition-colors hover:text-foreground/80 text-foreground/60"/>
        </div>
    </div>
</header>