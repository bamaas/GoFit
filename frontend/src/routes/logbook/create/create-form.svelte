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
  
  const form = superForm(data, {
    SPA: true,
    validators: zodClient(formSchema),
    onUpdate: async ({ form }) => {
      if (form.valid) {
        fetch(`${PUBLIC_BACKEND_BASE_URL}/v1/check-ins`, 
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            "date": form.data.date,
            "weight": form.data.weight,
          })
        })
        .then(response => {
          if (response.ok) {
            toast.success("Check-in added");
            goto("/")
          }
        })
        .catch(error => {
          toast.error("Oops! Something went wrong.");
          console.log(error);
        });
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
      <Form.Description>When did you measure?</Form.Description>
      <Form.FieldErrors />
    </Form.Control>
  </Form.Field>
  <Form.Field {form} name="weight" class="pb-5">
      <Form.Control let:attrs>
          <Form.Label>Weight</Form.Label>
          <Input {...attrs} bind:value={$formData.weight} />
      </Form.Control>
      <Form.Description>How much did you weight this morning?</Form.Description>
      <Form.FieldErrors />
  </Form.Field>
  <Form.Button class="w-full">Submit</Form.Button>
</form>