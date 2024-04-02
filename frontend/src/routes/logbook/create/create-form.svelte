<script lang="ts">
  import * as Form from "$lib/components/ui/form";
  import { Input } from "$lib/components/ui/input";
  import { formSchema, type FormSchema } from "./schema";
  import {
    type SuperValidated,
    type Infer,
    superForm,
  } from "sveltekit-superforms";
  import { zodClient } from "sveltekit-superforms/adapters";
  import { PUBLIC_BACKEND_BASE_URL } from '$env/static/public';
  import { goto } from '$app/navigation';
  import { dateProxy } from "sveltekit-superforms";
  import { toast } from "svelte-sonner";
 
  export let data: SuperValidated<Infer<FormSchema>>;

  let uuid = "";

  type CheckIn = {
    id: string;
    date: string;
    weight: number;
  };

  function postCheckIn(data: checkIn){
    fetch(`${PUBLIC_BACKEND_BASE_URL}/v1/check-ins`, 
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            "date": data.date,
            "weight": data.weight,
          })
        })
        .then(response => {
          if (response.ok) {
            toast.success("Check-in added.");
            goto("/")
          } else {
            toast.error("Oops! Something went wrong.");
          }
        })
        .catch(error => {
          toast.error("Oops! Something went wrong.");
          console.log(error);
        }
    );
  }

  function updateCheckIn(data: checkIn){
    fetch(`${PUBLIC_BACKEND_BASE_URL}/v1/check-ins`, 
        {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            "uuid": data.uuid,
            "date": data.date,
            "weight": data.weight,
          })
        })
        .then(response => {
          if (response.ok) {
            toast.success("Check-in updated.");
            goto("/")
          }
        })
        .catch(error => {
          toast.error("Oops! Something went wrong.");
          console.log(error);
        }
    );
  }

  const form = superForm(data, {
    SPA: true,
    validators: zodClient(formSchema),
    onUpdate: async ({ form }) => {
      if (form.valid) {
        const data: CheckIn = {
          uuid: form.data.uuid,
          date: form.data.date,
          weight: form.data.weight,
        }
        if (form.data.uuid == "") {
          postCheckIn(data);
        } else {
          updateCheckIn(data);
        }
      }
    }

  });
 
  const proxyDate = dateProxy(form, 'date', { format: 'date' }); 

  const { form: formData, enhance } = form;

</script>
 
<form method="POST" use:enhance>
  <Form.Field {form} name="date" class="pb-5">
    <Form.Control let:attrs>
      <Form.Label>Date</Form.Label>
      <Input {...attrs} bind:value={$proxyDate} />
      {#if data.data.uuid == ""}
        <Form.Description>When did you measure?</Form.Description>
      {:else}
        <Form.Description>This is the date you measured</Form.Description>
      {/if}
      <Form.FieldErrors />
    </Form.Control>
  </Form.Field>
  <Form.Field {form} name="weight" class="pb-5">
      <Form.Control let:attrs>
          <Form.Label>Weight</Form.Label>
          {#if data.data.uuid == ""}
            <Input {...attrs} bind:value={$formData.weight} autofocus/>
          {:else}
            <Input {...attrs} bind:value={$formData.weight}/>
          {/if}
      </Form.Control>
      {#if data.data.uuid == ""}
        <Form.Description>How much did you weight?</Form.Description>
      {:else}
        <Form.Description>That's how much you weighed that day</Form.Description>
      {/if}
      <Form.FieldErrors />
  </Form.Field>
  <Form.Button class="w-full">
    {#if data.data.uuid == ""}
      Add
    {:else}
      Update
    {/if}
  </Form.Button>
</form>