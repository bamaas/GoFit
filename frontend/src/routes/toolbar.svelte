<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
	import { PlusIcon, NotebookTextIcon, BarChart3Icon } from "lucide-svelte"
    import { page } from '$app/stores';
</script>

<div class="flex-none w-full border-t h-14 bg-gray-900 sticky bottom-0">
    <div class="flex container h-14 max-w-screen-2xl items-center justify-between">
        <a href="/logbook">
            <Button variant=ghost class={$page.url.pathname.includes("logbook") ? 'bg-accent' : ''}>
              <NotebookTextIcon class="size-4"/>
            </Button>
        </a>
        <a href="/logbook/create">
            <Button>
              <PlusIcon class="size-4"/>
            </Button>
        </a>
        <a href="/summary">
            <Button variant=ghost class={$page.url.pathname.includes("summary") ? 'bg-accent' : ''}>
              <BarChart3Icon class="size-4"/>
            </Button>
        </a>
    </div>
</div>