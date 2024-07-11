<script lang="ts">
    import MoveUpIcon from "lucide-svelte/icons/move-up";
    import MoveDownIcon from "lucide-svelte/icons/move-down";
    import MinusIcon from "lucide-svelte/icons/minus";
    import { profileStore } from "$lib/stores/profile";
	import type { Unsubscriber } from "svelte/store";
	import { onDestroy } from "svelte";

    export let weightDiff: number;
    let moveDownIconColor: string = "white";
    let moveUpIconColor: string = "white";

    const unsubProfileStore: Unsubscriber= profileStore.subscribe((userProfile) => {
        if (!userProfile) return;
        const goal = userProfile.goal;
        if (goal == "cut"){
            moveUpIconColor = "red";
            moveDownIconColor = "green";
        } else if (goal == "bulk"){
            moveUpIconColor = "green";
            moveDownIconColor = "red";
        } else if (goal == "maintain"){
            moveUpIconColor = "red";
            moveDownIconColor = "red";
        }
    });

    onDestroy(() => {
        unsubProfileStore();
    });

</script>

<style>
    :global(.red) {
        color: #7f1d1d;
    }

    :global(.green) {
        color: rgba(23, 104, 51, 0.84);
    }
</style>

<div class="flex items-center justify-center">
    {#if weightDiff > 0}
        <MoveUpIcon class="h-4 w-4 {moveUpIconColor}"/>
    {:else if weightDiff < 0}
        <MoveDownIcon class="h-4 w-4 {moveDownIconColor}"/>
    {:else}
        <MinusIcon class="h-4 w-4"/>
    {/if}
    {#if weightDiff != 0}
        <span>{Math.abs(weightDiff).toFixed(1)} kg</span>
    {/if}
</div>
