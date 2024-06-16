<script lang="ts">
    import * as Card from "$lib/components/ui/card/index.js";
    import TrendingDownIcon from "lucide-svelte/icons/trending-down";
    import AwardIcon from "lucide-svelte/icons/award";
    import RocketIcon from "lucide-svelte/icons/rocket";
    import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
	import { onMount } from "svelte";
	import { request } from "$lib/functions/request.js";
	import BackButton from "$lib/components/gofit/buttons/PreviousPageButton.svelte"

    let apiDataAverageWeightThisWeek: Promise<any> = new Promise(() => {});
	let apiDataAverageWeightLastWeek: Promise<any> = new Promise(() => {});
	let averageWeightDifference: Promise<any> = new Promise(() => {});

	function getMonday(d: Date) {
		d = new Date(d);
		var day = d.getDay(),
			diff = d.getDate() - day + (day == 0 ? -6 : 1); // adjust when day is sunday
		return new Date(d.setDate(diff));
	}

	let today: string = new Date().toISOString().split("T")[0]
	let lastMonday: string = getMonday(new Date()).toISOString().split("T")[0]
	let previousWeekModay: string = getMonday(new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)).toISOString().split("T")[0]
	let previousWeekSunday: string = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split("T")[0]

    onMount(() => {
		apiDataAverageWeightThisWeek = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-average?start_time=${lastMonday}&end_time=${today}`)
		apiDataAverageWeightLastWeek = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-average?start_time=${previousWeekModay}&end_time=${previousWeekSunday}`)
		Promise.all([apiDataAverageWeightThisWeek, apiDataAverageWeightLastWeek]).then((results) => {
			averageWeightDifference = new Promise((resolve) => {
				resolve(results[0].weight_average - results[1].weight_average)
			})
		})
    });

</script>

<div class="container items-center py-4 max-w-screen-2xl">
	<div class="pb-3 pt-4 justify-between flex">
		<BackButton />
	</div>
	<div class="grid my-3 gap-3 md:grid-cols-2 lg:grid-cols-3">
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Average weight this week</Card.Title>
				<TrendingDownIcon class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
				{#await apiDataAverageWeightThisWeek then data}
					<div class="green text-2xl font-bold">
						{Math.abs(Number(data.weight_average)).toFixed(1)} kg
					</div>
					<p class="text-muted-foreground text-xs">Keep going!</p>
				{/await}
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Average weight last week</Card.Title>
				<RocketIcon class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
				{#await apiDataAverageWeightLastWeek then data}
					<div class="red text-2xl font-bold">
						{Math.abs(Number(data.weight_average)).toFixed(1)} kg
					</div>
					<p class="text-muted-foreground text-xs">Good work!</p>
				{/await}
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Average weight lost last week</Card.Title>
				<RocketIcon class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
				{#await averageWeightDifference then data}
					<div class="red text-2xl font-bold">
						{Math.abs(Number(data)).toFixed(1)} kg
					</div>
					<p class="text-muted-foreground text-xs">Good work!</p>
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
