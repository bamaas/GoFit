<script lang="ts">
	import { Toaster } from "$lib/components/ui/sonner";
	import '../app.pcss';
	import HeaderComponent from "./header.svelte";
	import ToolbarComponent from "./toolbar.svelte";
	import { page } from '$app/stores';
	import { fetchUserProfile } from "$lib/functions/profile.js";
	import { profileReadable } from "$lib/stores/profile.js";
	import type { Unsubscriber } from "svelte/store";
	import { onDestroy } from "svelte";

	const unsubPageStore: Unsubscriber = page.subscribe(page => {
		const paths: string[] = ["/", "/login"]
		if(paths.indexOf(page.url.pathname) === -1 && $profileReadable === undefined) {
			fetchUserProfile();
		}
	})

	onDestroy(() => {
		unsubPageStore();
	});

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
