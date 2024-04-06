<script lang="ts">
    import { createTable, Render, Subscribe} from "svelte-headless-table";
    import * as Table from "$lib/components/ui/table";
    import { onMount } from "svelte";
    import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
    import { apiData, checkIns, type CheckIn } from "./store";
    import { addSortBy } from "svelte-headless-table/plugins";
    import { Button } from "$lib/components/ui/button/index.js";
	import { ArrowUpDown } from "lucide-svelte";
	import { toast } from "svelte-sonner";
	import { goto } from "$app/navigation";
    import { Skeleton } from "$lib/components/ui/skeleton/index.js";

    let promise: Promise<CheckIn[]> = new Promise(() => {});

    const table = createTable(checkIns, {
        sort: addSortBy(),
    });

    const columns = table.createColumns([
        table.column({
            accessor: "date",
            header: "Date",
            cell: ({ value }) => {return value.split("T")[0]},
        }),
        table.column({
            accessor: "weight",
            header: "Weight",
            cell: ({ value }) => {return value + " kg"},
        })
    ]);

    const { headerRows, pageRows, tableAttrs, tableBodyAttrs } =
    table.createViewModel(columns);

    function handleClick(value: string){
        if (value != undefined) {
            goto(`/logbook/edit/${value}`);
        }
    }

    onMount(() => {
        let dummyData: CheckIn[] = []
        for (let i = 0; i < 18; i++) {
            dummyData.push({
                uuid: String(i),
                date: "2021-10-01T00:00:00Z",
                weight: 70
            });
        }
        apiData.set(dummyData);

        promise = (async () => {
            const res = await fetch(`${PUBLIC_BACKEND_BASE_URL}/v1/check-ins`);
            const data = await res.json();
            return data; 
        })();

        promise.then(data => {
            apiData.set(data);
        }).catch(() => {
            toast.error("Oops! Failed fetching data from server.");
        });
    });

</script>

<div class="rounded-md border">
    <Table.Root {...$tableAttrs}>
        <Table.Header>
        {#each $headerRows as headerRow}
            <Subscribe rowAttrs={headerRow.attrs()}>
            <Table.Row>
                {#each headerRow.cells as cell (cell.id)}
                <Subscribe attrs={cell.attrs()} let:attrs props={cell.props()} let:props >
                    <Table.Head {...attrs}>
                        {#if cell.id === "weight"}
                            <div class="text-right">
                            <Render of={cell.render()} />
                            </div>
                        {:else if cell.id === "date"}
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
                            {#if cell.id === "weight"}
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