<script lang="ts">
    import * as Card from "$lib/components/ui/card/index.js";
    import TrendingDownIcon from "lucide-svelte/icons/trending-down";
    import AwardIcon from "lucide-svelte/icons/award";
    import RocketIcon from "lucide-svelte/icons/rocket";
    import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
	import { onMount } from "svelte";
	import { request } from "$lib/functions/request.js";

    let apiDataWeekAgoWeightDifference: Promise<any> = new Promise(() => {});
	let apiDataNinetyDaysAgoWeightDifference: Promise<any> = new Promise(() => {});
	let apiDataAllTimeWeightDifference: Promise<any> = new Promise(() => {});

    onMount(() => {
        let today: string = new Date().toISOString().split("T")[0]
        let weekAgo: string = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split("T")[0]
		let ninetyDaysAgo: string = new Date(Date.now() - 90 * 24 * 60 * 60 * 1000).toISOString().split("T")[0]
		let allTime: string = "1990-01-01"
        apiDataWeekAgoWeightDifference = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-difference?start_time=${weekAgo}&end_time=${today}`)
		apiDataNinetyDaysAgoWeightDifference = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-difference?start_time=${ninetyDaysAgo}&end_time=${today}`)
		apiDataAllTimeWeightDifference = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-difference?end_time=${today}`)
    });

</script>

<div class="container max-w-screen-2xl items-center py-14">
	<div class="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Lost in 7 days</Card.Title>
				<TrendingDownIcon class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
				{#await apiDataWeekAgoWeightDifference then data}
					<div class="green text-2xl font-bold">
						{Math.abs(data.weight_difference)} kg
					</div>
					<p class="text-muted-foreground text-xs">Keep going!</p>
				{/await}
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Lost in 90 days</Card.Title>
				<RocketIcon class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
				{#await apiDataNinetyDaysAgoWeightDifference then data}
					<div class="red text-2xl font-bold">
						{Math.abs(data.weight_difference)} kg
					</div>
					<p class="text-muted-foreground text-xs">Good work!</p>
				{/await}
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Lost all time</Card.Title>
				<AwardIcon class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
				{#await apiDataAllTimeWeightDifference then data}
					<div class="text-2xl font-bold">{Math.abs(data.weight_difference)} kg</div>
					<p class="text-muted-foreground text-xs">Amazing!</p>
				{/await}
			</Card.Content>
		</Card.Root>
	</div>
</div>

<style>
	.green {
		color: rgba(23, 104, 51, 0.84);
	}

	.red {
		color: #7f1d1d;
	}
</style>
