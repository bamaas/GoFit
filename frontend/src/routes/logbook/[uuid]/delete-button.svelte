<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import { Trash2Icon } from "lucide-svelte"
    import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
    import { page } from '$app/stores';
	import { toast } from "svelte-sonner";
	import { goto } from "$app/navigation";

    let uuid = "";
    $: uuid = $page.params.uuid;

    function deleteCheckIn(){
        fetch(`${PUBLIC_BACKEND_BASE_URL}/v1/check-ins/${uuid}`, 
        {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
          },
        })
        .then(response => {
          if (response.ok) {
            toast.success("Check-in deleted");
            goto("/logbook")
          }
        })
        .catch(error => {
          toast.error("Oops! Something went wrong.");
          console.log(error);
        });
      }

</script>
   
<Button variant="destructive" on:click={deleteCheckIn}>
    <Trash2Icon class="size-3.5"/>
</Button>