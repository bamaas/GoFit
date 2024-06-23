<script lang="ts">
    import * as Card from "$lib/components/ui/card/index.js";
    import TrendingDownIcon from "lucide-svelte/icons/trending-down";
    import AwardIcon from "lucide-svelte/icons/award";
	import CalendarIcon from "lucide-svelte/icons/calendar";
    import RocketIcon from "lucide-svelte/icons/rocket";
    import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
	import { onMount } from "svelte";
	import { request } from "$lib/functions/request.js";
	import { goal } from "$lib/stores/profile.js";
	import { Chart, Svg, Axis, TooltipItem, Tooltip, Highlight, Spline } from 'layerchart';

	// Current average weight
    let apiDataAverageWeight: Promise<any> = new Promise(() => {});

	// Lost all time
	let apiDataAllTimeWeightDifference: Promise<any> = new Promise(() => {});

	// Gained/lost this week
	let apiDataAverageWeightThisWeek: Promise<any> = new Promise(() => {});
	let apiDataAverageWeightLastWeek: Promise<any> = new Promise(() => {});
	let apiDataAverageWeightDifference: Promise<any> = new Promise(() => {});
	let apiDataWeeklyAverageWeight: Promise<any> = new Promise(() => {});

	let today: string = new Date().toISOString().split("T")[0]
	let lastMonday: string = getMonday(new Date()).toISOString().split("T")[0]
	let sevenDaysAgo: string = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split("T")[0]
	let previousWeekModay: string = getMonday(new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)).toISOString().split("T")[0]
	let previousWeekSunday: string = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split("T")[0]

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
		apiDataAverageWeight = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-average?start_time=${sevenDaysAgo}&end_time=${today}`)

		apiDataAllTimeWeightDifference = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-difference`)

		apiDataAverageWeightThisWeek = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-average?start_time=${lastMonday}&end_time=${today}`)
		apiDataAverageWeightLastWeek = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-average?start_time=${previousWeekModay}&end_time=${previousWeekSunday}`)
		Promise.all([apiDataAverageWeightThisWeek, apiDataAverageWeightLastWeek]).then((results) => {
			apiDataAverageWeightDifference = new Promise((resolve) => {
				resolve(results[0].weight_average - results[1].weight_average)
			})
		})

		apiDataWeeklyAverageWeight = request(`${PUBLIC_BACKEND_BASE_URL}/v1/stats/weight-average-by-week`)
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

<div class="container max-w-screen-2xl items-center py-8">
	<div class="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Current average weight</Card.Title>
				<RocketIcon class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
				{#await apiDataAverageWeight then data}
					<div class="text-2xl font-bold">
						{Math.abs(Number(data.weight_average)).toFixed(1)} kg
					</div>
					<p class="text-muted-foreground text-xs">Keep it up!</p>
				{/await}
			</Card.Content>
		</Card.Root>
		<Card.Root>
			{#await apiDataAverageWeightDifference then data}
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
					<div class="text-2xl font-bold {getColorClass($goal, data)}">
						{#if data > 0}+{/if}{Number(data).toFixed(1)} kg
					</div>
					<p class="text-muted-foreground text-xs">
						{#if goalMet($goal, data)}
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
	{#await apiDataWeeklyAverageWeight then dateSeriesData}
    <Card.Root class="mt-3">
        <Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
            <Card.Title class="text-sm font-medium">Average weight per week</Card.Title>
            <CalendarIcon class="text-muted-foreground h-4 w-4" />
        </Card.Header>
        <Card.Content>
            <div class="h-[300px] p-2 rounded">
                <Chart
                  data={dateSeriesData}
                  x="week"
                  y="weight"
                  yDomain={[0, null]}
                  yNice
                  padding={{ left: 16, bottom: 24 }}
                  tooltip={{ mode: "bisect-x" }}
                >
                  <Svg>
                    <Axis placement="left" grid rule />
                    <Axis
                        placement="bottom"
                        format={(d) => d}
                        rule
                    />
                    <Spline class="stroke-2 stroke-primary" />
                    <Highlight points lines />
                  </Svg>
                  <Tooltip header={(data) => "Week: " + data.week} let:data>
                    <TooltipItem label="Weight" value={data.weight.toFixed(1)} />
                  </Tooltip>
                </Chart>
            </div>
        </Card.Content>
    </Card.Root>    
    {/await} 
</div>