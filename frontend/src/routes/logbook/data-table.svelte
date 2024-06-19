<script lang="ts">
    import { createTable, Render, Subscribe, createRender} from "svelte-headless-table";
    import * as Table from "$lib/components/ui/table";
    import { onMount } from "svelte";
    import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
    import { apiData, checkIns, type CheckIn, type Metadata, type ApiResponse } from "./store";
    import { Button } from "$lib/components/ui/button/index.js";
	import { ArrowLeft, ArrowRight, XIcon } from "lucide-svelte";
	import { toast } from "svelte-sonner";
	import { goto } from "$app/navigation";
    import { Skeleton } from "$lib/components/ui/skeleton/index.js";
    import { page } from '$app/stores';
    import WeightDifference from "./weight-diff-column.svelte";
    import { request } from "$lib/functions/request";
    import moment from "moment";

    let pageNumber: number = Number($page.url.searchParams.get('page') || 1);
    let pageSize: number = 30;
    $: pageNumber;
    $: hasNextPage = false;
    $: hasPrevPage = false;

    let lastPage: number = 1;
    $: lastPage;

    $: recordsPresent = false;
    $: showPageNav = false;

    let promise: Promise<ApiResponse> = new Promise(() => {});

    onMount(() => {

        // Load dummy data
        let data: CheckIn[] = []
        for (let i = 0; i < 18; i++) {
            data.push({
                uuid: String(i),
                datetime: 1713034819,
                weight: 70,
                notes: "This is a note.",
                moving_average: 5,
                weight_difference: 0.5,
            });
        }
        let metadata: Metadata = {
            total_records: 30,
            current_page: 1,
            page_size: pageSize,
            first_page: 1,
            last_page: 1,
        }
        let dummyData = {
            data: data,
            metadata: metadata
        }

        apiData.set(dummyData);

        fetchData(pageNumber, dateRangeFilter);

        apiData.subscribe((data) => {

            // Check if records are present
            if (data.data.length == 0) {
                recordsPresent = false;
            } else {
                recordsPresent = true;
            }

            // show page navigation if there are more than 1 page
            if (data.metadata.last_page > 1) {
                showPageNav = true;
            } else {
                showPageNav = false;
            }

            // Check if there are more pages
            if (data.metadata.current_page < data.metadata.last_page) {
                hasNextPage = true;
            } else {
                hasNextPage = false;
            }
            if (data.metadata.current_page > data.metadata.first_page) {
                hasPrevPage = true;
            } else {
                hasPrevPage = false;
            }
            lastPage = data.metadata.last_page;
        });
    });

    function fetchData(pageNumber: number, dateFilter: DateRange | undefined): void {

        let uri: string = `/v1/check-ins?page=${pageNumber}&page_size=${pageSize}`
        if (dateFilter != undefined) {
            // setQueryParam('start_time', dateFilter?.start?.toString());
            // setQueryParam('end_time', dateFilter?.end?.toString());
            uri = uri + `&start_time=${dateFilter?.start?.toString()}&end_time=${dateFilter?.end?.toString()}`;
        }

        promise = (async () => {
            return await request(`${PUBLIC_BACKEND_BASE_URL}${uri}`);
        })();
        promise.then(response => {
            apiData.set(response);
        }).catch(() => {
            toast.error("Failed fetching data from server.", {description: "Oops!", cancel: { label: "X" }});
        });
    }

    const table = createTable(checkIns, {});

    const columns = table.createColumns([
        table.column({
            accessor: "datetime",
            header: "Date",
            cell: ({ value }) => {
                let date = new Date(value*1000).toISOString().split('T')[0];
                return moment(date).format("Do MMM YY");
            },
        }),
        table.column({
            accessor: "weight",
            header: "Weight",
            cell: ({ value }) => {return value.toFixed(1) + " kg"},
        }),
        table.column({
            accessor: "weight_difference",
            header: "Difference",
            cell: ({ value }) => createRender(WeightDifference, { weightDiff: value }),
        }),
        table.column({
            id: "movingAvg",
            accessor: "moving_average",
            header: "M. Avg",
            cell: ({ value }) => {return value.toFixed(1) + " kg"},
        }),
    ]);

    function handleClick(value: string){
        if (value != undefined) {
            goto(`/logbook/edit/${value}`);
        }
    }

    const { headerRows, pageRows, tableAttrs, tableBodyAttrs } = table.createViewModel(columns);

    function gotoNextPage(): void {
        pageNumber = pageNumber + 1;
        setQueryParam('page', String(pageNumber));
        fetchData(pageNumber, dateRangeFilter);
    }

    function goToPreviousPage(): void {
        pageNumber = pageNumber - 1;
        setQueryParam('page', String(pageNumber));
        fetchData(pageNumber, dateRangeFilter);
    }

    function setQueryParam(key: string, value: string): void {
        let query = new URLSearchParams($page.url.searchParams.toString());
        query.set(key, value);
        goto(`?${query.toString()}`);
    }

    // Date range picker TODO: put in seperate sub component
    import CalendarIcon from "lucide-svelte/icons/calendar";
    import type { DateRange } from "bits-ui";
    import {
      DateFormatter,
      getLocalTimeZone,
      today,
	  CalendarDate

    } from "@internationalized/date";
    import { cn } from "$lib/utils.js";
    import { RangeCalendar } from "$lib/components/ui/range-calendar/index.js";
    import * as Popover from "$lib/components/ui/popover/index.js";
	import { set } from "zod";

    let rangeCalendarOpen: boolean = false;
    const df = new DateFormatter("en-US", {
      dateStyle: "medium"
    });
    let dateRangeFilter: DateRange | undefined = undefined;
    $: dateRangeFilter: fetchData(pageNumber, dateRangeFilter);
    let startTime: number = Number($page.url.searchParams.get('start_time') || undefined);
    let endTime: number = Number($page.url.searchParams.get('end_time') || undefined);

</script>

<div class="">
    <div class="gap-2 mb-4 mt-6">
        <Popover.Root openFocus bind:open={rangeCalendarOpen} closeOnEscape closeOnOutsideClick>
            <div class="flex">
            <Popover.Trigger asChild let:builder>
            <Button
                variant="outline"
                class={cn(
                "w-screen justify-start text-left font-normal",
                !dateRangeFilter && "text-muted-foreground"
                )}
                builders={[builder]}
            >
                <CalendarIcon class="mr-2 h-4 w-4" />
                {#if dateRangeFilter && dateRangeFilter.start}
                {#if dateRangeFilter.end}
                    {df.format(dateRangeFilter.start.toDate(getLocalTimeZone()))} - {df.format(
                    dateRangeFilter.end.toDate(getLocalTimeZone())
                    )}
                {:else}
                    {df.format(dateRangeFilter.start.toDate(getLocalTimeZone()))}
                {/if}
                {:else}
                Pick a date
                {/if}
            </Button>
            </Popover.Trigger>
            {#if dateRangeFilter && dateRangeFilter.start && dateRangeFilter.end}
                <Button variant="outline" class="ml-4" size="default" on:click={() => dateRangeFilter = undefined}>
                    <XIcon class="h-3 w-3" />
                </Button>
            {/if}
            </div>
            <Popover.Content class="w-auto p-0" align="start">
            <RangeCalendar
                bind:value={dateRangeFilter}
                initialFocus
                numberOfMonths={2}
                weekStartsOn={1}
                placeholder={dateRangeFilter?.start}
                minValue={new CalendarDate(1900, 1, 1)}
                maxValue={today(getLocalTimeZone())}
                onValueChange={(v) => {
                    if (v.start && v.end) {
                        rangeCalendarOpen = false;
                    }
                }}
            />
            </Popover.Content>
        </Popover.Root>
    </div>
    <!-- <Button variant="outline" size="default" on:click={clearDateRangeFilter}>
        <XIcon class="size-3"/>
    </Button> -->
</div>
  {#if recordsPresent == true}
    <div> 
        <div class="rounded-md border">
            <Table.Root {...$tableAttrs}>
                <Table.Header>
                {#each $headerRows as headerRow}
                    <Subscribe rowAttrs={headerRow.attrs()}>
                    <Table.Row>
                        {#each headerRow.cells as cell (cell.id)}
                        <Subscribe attrs={cell.attrs()} let:attrs props={cell.props()} let:props >
                            <Table.Head {...attrs}>
                                {#if cell.id =="notes"}
                                    <div class="invisible lg:visible">
                                        <Render of={cell.render()} />
                                    </div>
                                {:else}
                                    <div class="text-center">
                                        <Render of={cell.render()} />
                                    </div>
                                {/if}
                            </Table.Head>
                        </Subscribe>
                        {/each}
                    </Table.Row>
                    </Subscribe>
                {/each}
                </Table.Header>
                <Table.Body {...$tableBodyAttrs}>
                {#each $pageRows as row (row.id)}
                    <Subscribe rowAttrs={row.attrs()} let:rowAttrs>
                    <Table.Row {...rowAttrs} on:click={() => handleClick(row.original['uuid'])} class="cursor-pointer">
                        {#each row.cells as cell (cell.id)}
                        <Subscribe attrs={cell.attrs()} let:attrs>
                            <Table.Cell {...attrs}>
                                {#await promise}
                                    <Skeleton class="h-4 w-full" />
                                {:then}
                                    {#if cell.id =="notes"}
                                        <div class="invisible lg:visible">
                                            <Render of={cell.render()} />
                                        </div>
                                    {:else}
                                        <div class="text-center">
                                            <Render of={cell.render()} />
                                        </div>
                                    {/if}
                                {:catch}
                                    <Skeleton class="h-4 w-full" />
                                {/await}
                            </Table.Cell>
                        </Subscribe>
                        {/each}
                    </Table.Row>
                    </Subscribe>
                {/each}
                </Table.Body>
            </Table.Root>
        </div>
        <!-- svelte-ignore empty-block -->
        {#await promise}
        {:then}
        {#if showPageNav == true}
            <div class="flex items-center justify-between space-x-4 py-4">
                <Button variant="outline" size="lg" disabled={!hasPrevPage} on:click={goToPreviousPage}>
                    <ArrowLeft class="size-4"/>
                </Button>
                <span class="text-sm text-muted-foreground">
                    Page {pageNumber} of {lastPage}
                </span>
                <Button variant="outline" size="lg" disabled={!hasNextPage} on:click={gotoNextPage}>
                    <ArrowRight class="size-4"/>
                </Button>
            </div>
        {/if}
        {/await}
    </div>
{:else}
    <!-- No records present -->
    <div class="text-center items-center justify-center align-middle mt-36">
        <h1 class="text-2xl font-semibold tracking-tight">Let's get started</h1>
        <p class="text-sm text-muted-foreground mt-2">Add your first check-in to get started.</p>
        <a href="/logbook/create">
            <Button class="mt-8">Add check-in</Button>
        </a>
    </div>
{/if}
