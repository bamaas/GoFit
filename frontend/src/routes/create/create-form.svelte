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
      SPA: true,
      validators: zodClient(formSchema),
      onUpdate: async ({ form }) => {
        if (form.valid) {
          fetch("http://localhost:8080/v1/check-ins", 
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              "Weight": form.data.weight,
            })
          })
          .then(response => response.json())
          .then((data) => {
            console.log("success");
          }).catch(error => {
            console.log("error");
          });
        }
      }

    });
   
    const { form: formData, enhance } = form;
  </script>
   
  <form method="POST" use:enhance>
    <Form.Field {form} name="weight">
        <!-- <Form.Control let:attrs>
            <Form.Label>Date</Form.Label>
            <Input {...attrs} bind:value={$formData.date} />
        </Form.Control> -->
        <Form.Control let:attrs>
            <Form.Label>Weight</Form.Label>
            <Input {...attrs} bind:value={$formData.weight} />
        </Form.Control>
        <Form.Description>How much did you weight</Form.Description>
      <Form.FieldErrors />
    </Form.Field>
    <Form.Button>Submit</Form.Button>
  </form>