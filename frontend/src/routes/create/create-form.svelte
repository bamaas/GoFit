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
   
    export let data: SuperValidated<Infer<FormSchema>>;
   
    const form = superForm(data, {
      validators: zodClient(formSchema),
    });
   
    const { form: formData, enhance } = form;
  </script>
   
  <form method="POST" use:enhance>
    <Form.Field {form} name="weight">
        <Form.Control let:attrs>
            <Form.Label>Date</Form.Label>
            <Input {...attrs} bind:value={$formData.date} />
        </Form.Control>
        <Form.Control let:attrs>
            <Form.Label>Weight</Form.Label>
            <Input {...attrs} bind:value={$formData.weight} />
        </Form.Control>
      <Form.FieldErrors />
    </Form.Field>
    <Form.Button>Submit</Form.Button>
  </form>