<script lang="ts">
	import { Toaster } from "$lib/components/ui/sonner";
	import '../app.pcss';
	import HeaderComponent from "./header.svelte";
	import ToolbarComponent from "./toolbar.svelte";
	import { page } from '$app/stores';

	function getHeaderTitleFromPathname(pathname: string): string {
		let title: string = pathname.split('/')[1];
		// Capitalize first letter
		title =  title.charAt(0).toUpperCase() + title.slice(1); 
		if (title === "") {
			title = "GoFit"
		}
		return title
	}

	$: headerTitle = getHeaderTitleFromPathname($page.url.pathname)
	$: renderToolbar = $page.data.renderToolbar;
	$: renderHeader = $page.data.renderHeader;

</script>

<svelte:head>
	<title>{$page.data.title} • GoFit</title>
</svelte:head>

<div class="flex flex-col h-screen">
	{#if renderHeader != false}
		<HeaderComponent title={headerTitle || ""}/>
	{/if}
	<Toaster position="top-right" />
	<div class="mb-auto">
		<slot />
	</div>
	{#if renderToolbar}
		<ToolbarComponent />
	{/if}
</div>
