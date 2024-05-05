<script lang="ts">
    import * as Card from "$lib/components/ui/card/index.js";
    import TrendingDownIcon from "lucide-svelte/icons/trending-down";
    import AwardIcon from "lucide-svelte/icons/award";
    import RocketIcon from "lucide-svelte/icons/rocket";
    import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
	import { onMount } from "svelte";
	import { request } from "$lib/functions/request.js";

    type StatsResponse = {
        stats: {
            weight_difference: {
                week_ago: number;
                ninety_days_ago: number;
                all_time: number;
            }
        }
    }

    let apiData: Promise<StatsResponse> = new Promise(() => {});

    function fetchData() {
        return request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats`)
    }

    onMount(() => {
        apiData = fetchData();
    });

</script>

<style>
    .green {
        color: rgba(23, 104, 51, 0.84)
    }

    .red {
        color: #7f1d1d;
    }
</style>

<div class="container items-center py-14 max-w-screen-2xl">
    <div class="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
        <Card.Root>
            <Card.Header
                class="flex flex-row items-center justify-between space-y-0 pb-2">
                <Card.Title class="text-sm font-medium">Lost in 7 days</Card.Title>
                <TrendingDownIcon class="h-4 w-4 text-muted-foreground" />
            </Card.Header>
            <Card.Content>
                {#await apiData then data}
                    <div class="text-2xl font-bold green">{Math.abs(data.stats.weight_difference.week_ago)} kg</div>
                    <p class="text-xs text-muted-foreground">Keep going!</p>
                {/await}
            </Card.Content>
        </Card.Root>
        <Card.Root>
            <Card.Header
                class="flex flex-row items-center justify-between space-y-0 pb-2">
                <Card.Title class="text-sm font-medium">Lost in 90 days</Card.Title>
                <RocketIcon class="h-4 w-4 text-muted-foreground" />
            </Card.Header>
            <Card.Content>
                {#await apiData then data}
                    <div class="text-2xl font-bold red">{Math.abs(data.stats.weight_difference["ninety_days_ago"])} kg</div>
                    <p class="text-xs text-muted-foreground">Good work!</p>
                {/await}
            </Card.Content>
        </Card.Root>
        <Card.Root>
            <Card.Header
                class="flex flex-row items-center justify-between space-y-0 pb-2">
                <Card.Title class="text-sm font-medium">Lost all time</Card.Title>
                <AwardIcon class="h-4 w-4 text-muted-foreground" />
            </Card.Header>
            <Card.Content>
                {#await apiData then data}
                    <div class="text-2xl font-bold">{Math.abs(data.stats.weight_difference.all_time)} kg</div>
                    <p class="text-xs text-muted-foreground">Amazing!</p>
                {/await}
            </Card.Content>
        </Card.Root>
    </div>
</div>