<script lang="ts">
    import * as Card from "$lib/components/ui/card/index.js";
    import TrendingDownIcon from "lucide-svelte/icons/trending-down";
    import AwardIcon from "lucide-svelte/icons/award";
    import RocketIcon from "lucide-svelte/icons/rocket";
    import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
	import { onMount } from "svelte";
	import { request } from "$lib/functions/request.js";
	import { goal } from "$lib/stores/profile.js";

    let apiDataAverageWeightThisWeek: Promise<any> = new Promise(() => {});
	let apiDataWeightDifferenceThisWeek: Promise<any> = new Promise(() => {});
	let apiDataAllTimeWeightDifference: Promise<any> = new Promise(() => {});

	let today: string = new Date().toISOString().split("T")[0]
	let lastMonday: string = getMonday(new Date()).toISOString().split("T")[0]

	function getMonday(d: Date) {
		d = new Date(d);
		var day = d.getDay(),
			diff = d.getDate() - day + (day == 0 ? -6 : 1); // adjust when day is sunday
		return new Date(d.setDate(diff));
	}

	function getColorClass(goal: string, weight: number): string | undefined{
		// gain
		if (goal == "gain" && weight >= 0){
			return "green"
		} else if (goal == "gain") {
			return "red"
		// lose
		} else if (goal == "lose" && weight < 0) {
			return "green"
		} else if (goal == "lose") {
			return "red"
		// maintain
		} else if (goal == "maintain" && weight == 0) {
			return "green"
		} else if (goal == "maintain") {
			return "red"
		}
	}

	function goalMet(goal: string, weight: number): boolean | undefined {
		if (goal == "gain" && weight > 0) {
			return true
		} else if (goal == "gain") {
			return false
		} else if (goal == "lose" && weight < 0) {
			return true
		} else if (goal == "lose") {
			return false
		} else if (goal == "maintain" && weight == 0) {
			return true
		} else if (goal == "maintain") {
			return false
		}
	}

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
		<a href="/summary/weight-average">
			<Card.Root>
				<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
					<Card.Title class="text-sm font-medium">Average weight this week</Card.Title>
					<RocketIcon class="text-muted-foreground h-4 w-4" />
				</Card.Header>
				<Card.Content>
					{#await apiDataAverageWeightThisWeek then data}
						<div class="text-2xl font-bold">
							{Math.abs(Number(data.weight_average)).toFixed(1)} kg
						</div>
						<p class="text-muted-foreground text-xs">Keep it up!</p>
					{/await}
				</Card.Content>
			</Card.Root>
		</a>
		<Card.Root>
			{#await apiDataWeightDifferenceThisWeek then data}
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">
					{#if $goal == "gain"}
						Gained this week
					{:else if $goal == "lose"}
						Lost this week
					{/if}
				</Card.Title>
				<TrendingDownIcon class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
					<div class="text-2xl font-bold {getColorClass($goal, data.weight_difference)}">
						{#if data.weight_difference > 0}+{/if}{Number(data.weight_difference).toFixed(1)} kg
					</div>
					<p class="text-muted-foreground text-xs">
						{#if goalMet($goal, data.weight_difference)}
							Good work!
						{:else}
							Time to pull yourself together...
						{/if}
					</p>
			</Card.Content>
			{/await}
		</Card.Root>
		<Card.Root>
			{#await apiDataAllTimeWeightDifference then data}
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">
					{#if $goal == "gain"}
						Gained all time
					{:else if $goal == "lose"}
						Lost all time
					{/if}
				</Card.Title>
				<AwardIcon class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
					<div class="text-2xl font-bold {getColorClass($goal, data.weight_difference)}">
						{#if data.weight_difference > 0}+{/if}{Number(data.weight_difference).toFixed(1)} kg
					</div>
					<p class="text-muted-foreground text-xs">
						{#if goalMet($goal, data.weight_difference)}
							Amazing!
						{:else}
							Regain focus...
						{/if}
					</p>
			</Card.Content>
			{/await}
		</Card.Root>
	</div>
</div>