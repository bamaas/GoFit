<script lang="ts">
    import { createTable, Render, Subscribe} from "svelte-headless-table";
    import * as Table from "$lib/components/ui/table";
    import { onMount } from "svelte";
    import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
    import { apiData, checkIns, type CheckIn, type Metadata, type ApiResponse } from "./store";
    import { addSortBy } from "svelte-headless-table/plugins";
    import { Button } from "$lib/components/ui/button/index.js";
	import { ArrowUpDown, ArrowLeft, ArrowRight } from "lucide-svelte";
	import { toast } from "svelte-sonner";
	import { goto } from "$app/navigation";
    import { Skeleton } from "$lib/components/ui/skeleton/index.js";
    import { page } from '$app/stores';
	import { readable } from "svelte/store";

    let pageNumber: number = Number($page.url.searchParams.get('page') || 1);
    let pageSize: number = 30;
    $: pageNumber;
    $: hasNextPage = false;
    $: hasPrevPage = false;

    let lastPage: number = 1;
    $: lastPage;

    let promise: Promise<ApiResponse> = new Promise(() => {});

    onMount(() => {

        // Load dummy data
        let data: CheckIn[] = []
        for (let i = 0; i < 18; i++) {
            data.push({
                uuid: String(i),
                datetime: 1713034819,
                weight: 70,
                moving_average: 5
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

        fetchData(pageNumber);

        apiData.subscribe((data) => {
            console.log(data.metadata.last_page)
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

    function fetchData(pageNumber: number): void{
        promise = (async () => {
            const res = await fetch(`${PUBLIC_BACKEND_BASE_URL}/v1/check-ins?page=${pageNumber}&page_size=${pageSize}`);
            const response = await res.json();
            return response; 
        })();
        promise.then(response => {
            apiData.set(response);
        }).catch(() => {
            toast.error("Oops! Failed fetching data from server.");
        });
    }

    const table = createTable(checkIns, {
        sort: addSortBy()
    });

    const columns = table.createColumns([
        table.column({
            accessor: "datetime",
            header: "Date",
            cell: ({ value }) => {return new Date(value*1000).toISOString().split('T')[0]},
        }),
        table.column({
            accessor: "weight",
            header: "Weight",
            cell: ({ value }) => {return value + " kg"},
        }),
        table.column({
            id: "movingAvg",
            accessor: "moving_average",
            header: "M. Avg",
            cell: ({ value }) => {return value.toFixed(1) + " kg"},
        })
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
        fetchData(pageNumber);
    }

    function goToPreviousPage(): void {
        pageNumber = pageNumber - 1;
        setQueryParam('page', String(pageNumber));
        fetchData(pageNumber);
    }

    function setQueryParam(key: string, value: string): void {
        let query = new URLSearchParams($page.url.searchParams.toString());
        query.set(key, value);
        goto(`?${query.toString()}`);
    }
</script>

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
                            {#if cell.id === "weight" || cell.id === "movingAvg"}
                                <div class="text-right">
                                    <Render of={cell.render()} />
                                </div>
                            {:else if cell.id === "datetime"}
                                <Button variant="ghost" on:click={props.sort.toggle}>
                                    <Render of={cell.render()} />
                                    <ArrowUpDown class={"ml-2 h-4 w-4"} />
                                </Button>
                            {:else}
                                <Render of={cell.render()} />
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
                                {#if cell.id === "weight" || cell.id === "movingAvg"}
                                    <div class="text-right">
                                        <Render of={cell.render()} />
                                    </div>
                                {:else}
                                    <Render of={cell.render()} />
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
</div>