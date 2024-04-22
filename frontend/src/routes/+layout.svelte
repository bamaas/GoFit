<script lang="ts">
	import { Toaster } from "$lib/components/ui/sonner";
	import '../app.pcss';
	import HeaderComponent from "./header.svelte";
	import ToolbarComponent from "./toolbar.svelte";
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
</script>

<svelte:head>
	<title>{$page.data.title} • GoFit</title>
</svelte:head>

<div class="flex flex-col h-screen">
	<HeaderComponent title={navigationTitle || ""}/>
	<Toaster />
	<div class="mb-auto">
		<slot />
	</div>
	<ToolbarComponent />
</div>
