<script lang="ts">
	import { Toaster } from "$lib/components/ui/sonner";
	import '../app.pcss';
	import HeaderComponent from "./header.svelte";
	import { page } from '$app/stores';

	function getNavigationTitleFromPathname(pathname: string) {
		let title: string = pathname.split('/')[1];
		// Capitalize first letter
		title =  title.charAt(0).toUpperCase() + title.slice(1); 
		if (title === "") {
			title = "GoFit"
		}
		return title
	}

	$: navigationTitle = getNavigationTitleFromPathname($page.url.pathname)

	// to be removed
	import { Button } from "$lib/components/ui/button/index.js";
	import { PlusIcon, NotebookTextIcon, BarChart3Icon } from "lucide-svelte"
</script>

<svelte:head>
	<title>{$page.data.title} • GoFit</title>
</svelte:head>



<div class="flex flex-col h-screen justify-between">
	<HeaderComponent title={navigationTitle || ""}/>
	<Toaster />
	<slot />
	<footer class="flex-none w-full border-b h-14 bg-gray-900 sticky bottom-0">
		<div class="flex container h-14 max-w-screen-2xl items-center justify-between">
			<a href="/logbook">
				<Button variant=ghost>
				  <NotebookTextIcon class="size-4"/>
				</Button>
			</a>
			<a href="/logbook/create">
				<Button>
				  <PlusIcon class="size-4"/>
				</Button>
			</a>
			<a href="/summary">
				<Button variant=ghost>
				  <BarChart3Icon class="size-4"/>
				</Button>
			</a>
		</div>
	</footer>
</div>
