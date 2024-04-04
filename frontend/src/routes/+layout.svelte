<script>
	import { Toaster } from "$lib/components/ui/sonner";
	import '../app.pcss';
	import HeaderComponent from "./header.svelte";
</script>
<HeaderComponent />
<Toaster />

<slot />
