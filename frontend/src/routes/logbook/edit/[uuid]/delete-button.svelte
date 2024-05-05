<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import { Trash2Icon } from "lucide-svelte"
    import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
    import { page } from '$app/stores';
	  import { toast } from "svelte-sonner";
	  import { goto } from "$app/navigation";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";
	  import { request } from "$lib/functions/request";

    let uuid = "";
    $: uuid = $page.params.uuid;

    function deleteCheckIn(){
        request(`${PUBLIC_BACKEND_BASE_URL}/v1/check-ins/${uuid}`, {method: 'DELETE'})
        .then( () => {
            toast.success("Check-in deleted.",{
              description: "Bye bye...",
              cancel: { label: "X" }
            });
            goto("/logbook")
        });
      }

</script>

<AlertDialog.Root>
    <AlertDialog.Trigger>
        <Button variant="destructive">
            <Trash2Icon class="size-3.5"/>
        </Button>
    </AlertDialog.Trigger>
    <AlertDialog.Content>
      <AlertDialog.Header>
        <AlertDialog.Title>Are you absolutely sure?</AlertDialog.Title>
        <AlertDialog.Description>
          This action cannot be undone. This will permanently delete your check-in from our servers.
        </AlertDialog.Description>
      </AlertDialog.Header>
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <AlertDialog.Action on:click={deleteCheckIn}>Delete</AlertDialog.Action>
      </AlertDialog.Footer>
    </AlertDialog.Content>
</AlertDialog.Root>
   
