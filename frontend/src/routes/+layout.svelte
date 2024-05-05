<script lang="ts">
	import { Toaster } from "$lib/components/ui/sonner";
	import '../app.pcss';
	import HeaderComponent from "./header.svelte";
	import ToolbarComponent from "./toolbar.svelte";
	import { page } from '$app/stores';
	import { onMount } from "svelte";

	// function getHeaderTitle(): string {
	// 	if ($page.data.headerTitle != undefined) {
	// 		return $page.data.headerTitle
	// 	}
	// 	return getHeaderTitleFromPathname($page.url.pathname)
	// }

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

</script>

<svelte:head>
	<title>{$page.data.title} • GoFit</title>
</svelte:head>

<div class="flex flex-col h-screen">
	<HeaderComponent title={headerTitle || ""}/>
	<Toaster position="top-right" />
	<div class="mb-auto">
		<slot />
	</div>
	{#if renderToolbar}
		<ToolbarComponent />
	{/if}
</div>
