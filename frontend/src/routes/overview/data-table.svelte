<script lang="ts">
    import { createTable, Render, Subscribe } from "svelte-headless-table";
    import * as Table from "$lib/components/ui/table";
    import { onMount } from "svelte";
    import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
    import { apiData, checkIns } from "./store";

    onMount(async () => {
    fetch(`${PUBLIC_BACKEND_BASE_URL}/v1/check-ins`)
    .then(response => response.json())
    .then(data => {
        apiData.set(data);
    }).catch(error => {
        // TODO: error handling
        console.log(error);
        return [];
    });
    });

    const table = createTable(checkIns);

    const columns = table.createColumns([
        table.column({
            accessor: "created_at",
            header: "Date",
        }),
        table.column({
            accessor: "weight",
            header: "Weight",
            cell: ({ value }) => {return value + " kg"},
        })
    ]);

    const { headerRows, pageRows, tableAttrs, tableBodyAttrs } =
    table.createViewModel(columns);
</script>

<div class="rounded-md border">
    <Table.Root {...$tableAttrs}>
        <Table.Header>
        {#each $headerRows as headerRow}
            <Subscribe rowAttrs={headerRow.attrs()}>
            <Table.Row>
                {#each headerRow.cells as cell (cell.id)}
                <Subscribe attrs={cell.attrs()} let:attrs props={cell.props()}>
                    <Table.Head {...attrs}>
                        {#if cell.id === "weight"}
                            <div class="text-right">
                            <Render of={cell.render()} />
                            </div>
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
            <Table.Row {...rowAttrs}>
                {#each row.cells as cell (cell.id)}
                <Subscribe attrs={cell.attrs()} let:attrs>
                    <Table.Cell {...attrs}>
                        {#if cell.id === "weight"}
                            <div class="text-right">
                                <Render of={cell.render()} />
                            </div>
                        {:else}
                            <Render of={cell.render()} />
                        {/if}
                    </Table.Cell>
                </Subscribe>
                {/each}
            </Table.Row>
            </Subscribe>
        {/each}
        </Table.Body>
    </Table.Root>
</div>