<script lang="ts">
    import MoveUpIcon from "lucide-svelte/icons/move-up";
    import MoveDownIcon from "lucide-svelte/icons/move-down";
    import MinusIcon from "lucide-svelte/icons/minus";
    export let weightDiff: number;
</script>
<div class="flex flex-row-reverse">
    {#if weightDiff != 0}
        <span>{Math.abs(weightDiff)} kg</span>
    {/if}
    {#if weightDiff > 0}
        <MoveUpIcon class="h-4 w-4"/>
    {:else if weightDiff < 0}
        <MoveDownIcon class="h-4 w-4"/>
    {:else}
        <MinusIcon class="h-4 w-4"/>
    {/if}
</div>
