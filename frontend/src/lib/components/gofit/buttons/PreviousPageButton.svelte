<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import { ChevronLeft } from "lucide-svelte"
    import { afterNavigate } from '$app/navigation';
    import { base } from '$app/paths'
	import { onMount } from "svelte";
    import { page } from '$app/stores';

    let previousPage : string = base ;

    afterNavigate(({from}) => {
        previousPage = from?.url.pathname || $page.url.pathname.substring(0, $page.url.pathname.lastIndexOf('/'))
    })

</script>
   
<a href={previousPage}>
    <Button variant="outline">
        <ChevronLeft class="size-4"/>
    </Button>
</a>