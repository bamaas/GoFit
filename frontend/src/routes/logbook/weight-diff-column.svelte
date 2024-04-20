<script lang="ts">
    import MoveUpIcon from "lucide-svelte/icons/move-up";
    import MoveDownIcon from "lucide-svelte/icons/move-down";
    import MinusIcon from "lucide-svelte/icons/minus";
    export let weightDiff: number;
</script>
<div class="flex items-center justify-center">
    {#if weightDiff > 0}
        <MoveUpIcon style="color: #7f1d1d;" class="h-4 w-4"/>
    {:else if weightDiff < 0}
        <MoveDownIcon style="color: #1d7f3f;" class="h-4 w-4"/>
    {:else}
        <MinusIcon class="h-4 w-4"/>
    {/if}
    {#if weightDiff != 0}
        <span>{Math.abs(weightDiff).toFixed(1)} kg</span>
    {/if}
</div>
