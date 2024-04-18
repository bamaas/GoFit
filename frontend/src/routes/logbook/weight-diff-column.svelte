<script lang="ts">
    import MoveUpIcon from "lucide-svelte/icons/move-up";
    import MoveDownIcon from "lucide-svelte/icons/move-down";
    export let weightDiff: number;
</script>
<div class="flex flex-row-reverse">
    <span>{Math.abs(weightDiff)} kg</span>
    {#if weightDiff > 0}
        <MoveUpIcon class="h-4 w-4"/>
    {:else if weightDiff < 0}
        <MoveDownIcon class="h-4 w-4"/>
    {/if}
</div>
