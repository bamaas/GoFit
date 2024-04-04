<script lang="ts">
	import { Toaster } from "$lib/components/ui/sonner";
	import '../app.pcss';
	import HeaderComponent from "./header.svelte";
	import { page } from '$app/stores';

	let navigationTitle = getNavigationTitleFromPathname($page.url.pathname);

	function getNavigationTitleFromPathname(pathname: string) {
		let title: string = $page.url.pathname.split('/')[1];
		// Capitalize first letter
		title =  title.charAt(0).toUpperCase() + title.slice(1); 
		return title
	}
</script>

<svelte:head>
	<title>{$page.data.title} • GoFit</title>
</svelte:head>

<HeaderComponent title={navigationTitle}/>
<Toaster />

<slot />
