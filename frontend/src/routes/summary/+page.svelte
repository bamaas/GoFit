<script lang="ts">
    import * as Card from "$lib/components/ui/card/index.js";
    import TrendingDownIcon from "lucide-svelte/icons/trending-down";
    import AwardIcon from "lucide-svelte/icons/award";
    import RocketIcon from "lucide-svelte/icons/rocket";
    import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
	import { onMount } from "svelte";
	import { request } from "$lib/functions/request.js";

    let apiDataAverageWeightThisWeek: Promise<any> = new Promise(() => {});
	let apiDataWeightDifferenceThisWeek: Promise<any> = new Promise(() => {});
	let apiDataAllTimeWeightDifference: Promise<any> = new Promise(() => {});

	function getMonday(d: Date) {
		d = new Date(d);
		var day = d.getDay(),
			diff = d.getDate() - day + (day == 0 ? -6 : 1); // adjust when day is sunday
		return new Date(d.setDate(diff));
	}

	let today: string = new Date().toISOString().split("T")[0]
	let lastMonday: string = getMonday(new Date()).toISOString().split("T")[0]

    onMount(() => {
		apiDataAverageWeightThisWeek = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-average?start_time=${lastMonday}&end_time=${today}`)
		apiDataWeightDifferenceThisWeek = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-difference?start_time=${lastMonday}`)
		apiDataAllTimeWeightDifference = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-difference`)
    });

</script>

<style>
	.green {
		color: rgba(23, 104, 51, 0.84);
	}

	.red {
		color: #7f1d1d;
	}
</style>

<div class="container max-w-screen-2xl items-center py-14">
	<div class="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
		<a href="/summary/average">
			<Card.Root>
				<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
					<Card.Title class="text-sm font-medium">Average weight this week</Card.Title>
					<RocketIcon class="text-muted-foreground h-4 w-4" />
				</Card.Header>
				<Card.Content>
					{#await apiDataAverageWeightThisWeek then data}
						<div class="green text-2xl font-bold">
							{Number(data.weight_average).toFixed(1)} kg
						</div>
						<p class="text-muted-foreground text-xs">Keep it up!</p>
					{/await}
				</Card.Content>
			</Card.Root>
		</a>
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Lost this week</Card.Title>
				<TrendingDownIcon class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
				{#await apiDataWeightDifferenceThisWeek then data}
					<div class="red text-2xl font-bold">
						{Number(data.weight_difference).toFixed(1)} kg
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
					<!-- <div class="text-2xl font-bold">{Math.abs(data.weight_difference)} kg</div> -->
					<div class="text-2xl font-bold">{data.weight_difference} kg</div>
					<p class="text-muted-foreground text-xs">Amazing!</p>
				{/await}
			</Card.Content>
		</Card.Root>
	</div>
</div>