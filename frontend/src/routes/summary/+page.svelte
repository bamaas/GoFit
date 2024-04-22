<!-- https://github.com/huntabyte/shadcn-svelte/blob/main/apps/www/src/lib/components/docs/dashboard/dashboard-page.svelte -->
<script lang="ts">
    import * as Card from "$lib/components/ui/card/index.js";
    import TrendingDownIcon from "lucide-svelte/icons/trending-down";
    import AwardIcon from "lucide-svelte/icons/award";
    import RocketIcon from "lucide-svelte/icons/rocket";

    export let data;
</script>
<div class="container items-center py-14 max-w-screen-2xl">
    <div class="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
        <Card.Root>
            <Card.Header
                class="flex flex-row items-center justify-between space-y-0 pb-2">
                <Card.Title class="text-sm font-medium">Lost in 7 days</Card.Title>
                <TrendingDownIcon class="h-4 w-4 text-muted-foreground" />
            </Card.Header>
            <Card.Content>
                <div style="color: rgba(23, 104, 51, 0.84);" class="text-2xl font-bold">{data.stats.stats.weight_diff.week_ago} kg</div>
                <p class="text-xs text-muted-foreground">Keep going!</p>
            </Card.Content>
        </Card.Root>
        <Card.Root>
            <Card.Header
                class="flex flex-row items-center justify-between space-y-0 pb-2">
                <Card.Title class="text-sm font-medium">Lost in 90 days</Card.Title>
                <RocketIcon class="h-4 w-4 text-muted-foreground" />
            </Card.Header>
            <Card.Content>
                <div style="color: #7f1d1d;" class="text-2xl font-bold">{data.stats.stats.weight_diff["90_days_ago"]} kg</div>
                <p class="text-xs text-muted-foreground">Good work!</p>
            </Card.Content>
        </Card.Root>
        <Card.Root>
            <Card.Header
                class="flex flex-row items-center justify-between space-y-0 pb-2">
                <Card.Title class="text-sm font-medium">Lost all time</Card.Title>
                <AwardIcon class="h-4 w-4 text-muted-foreground" />
            </Card.Header>
            <Card.Content>
                <div class="text-2xl font-bold">{data.stats.stats.weight_diff.all_time} kg</div>
                <p class="text-xs text-muted-foreground">Amazing!</p>
            </Card.Content>
        </Card.Root>
    </div>
</div>