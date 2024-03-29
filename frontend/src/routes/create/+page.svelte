<script lang="ts">
    import type { PageData } from "./$types.js";
    import CreateForm from "./create-form.svelte";
    export let data: PageData;
</script>
   
<CreateForm data={data.form} />