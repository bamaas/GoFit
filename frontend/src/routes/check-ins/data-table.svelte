<script lang="ts">
    import { createTable, Render, Subscribe } from "svelte-headless-table";
    import { readable } from "svelte/store";
    import * as Table from "$lib/components/ui/table";
    import DataTableCreateButton from "./data-table-create-button.svelte";
    import { onMount } from "svelte";
    import { writable, derived } from 'svelte/store';
    import { PUBLIC_BASE_URL } from "$env/static/public";

    onMount(async () => {
    fetch(`${PUBLIC_BASE_URL}/v1/check-ins`)
    .then(response => response.json())
    .then(data => {
        console.log(data);
        apiData.set(data);
    }).catch(error => {
        console.log(error);
        return [];
    });
    });

    // type CheckIn = {
    //   id: string;
    //   created_at: string;
    //   weight: number;
    // };

    export const apiData = writable([]);

    export const entries = derived(apiData, ($apiData) => {
    if ($apiData.length > 0){
        return $apiData;
    }
    return [];
    });

    const table = createTable(apiData);

    const columns = table.createColumns([
        // table.column({
        // accessor: "id",
        // header: "ID",
        // }),
        table.column({
        accessor: "created_at",
        header: "Date",
        }),
        table.column({
        accessor: "weight",
        header: "Weight",
        cell: ({ value }) => {return value + " kg"},
        }),
        // table.column({
        // accessor: ({ id }) => id,
        // header: "",
        // }),
    ]);

    const { headerRows, pageRows, tableAttrs, tableBodyAttrs } =
    table.createViewModel(columns);
</script>

<div>
    <div class="py-4">
        <DataTableCreateButton/>
    </div>
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
</div>