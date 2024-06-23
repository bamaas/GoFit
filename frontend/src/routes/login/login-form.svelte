<script lang="ts">
	  import { goto } from '$app/navigation';
    import { PUBLIC_BACKEND_BASE_URL } from '$env/static/public';
    import * as Form from "$lib/components/ui/form";
    import { Input } from "$lib/components/ui/input";
	  import { toast } from 'svelte-sonner';
    import { formSchema, type FormSchema } from "./schema";
    import {
      type SuperValidated,
      type Infer,
      superForm,
    } from "sveltekit-superforms";
    import { zodClient } from "sveltekit-superforms/adapters";
    import LoaderCircleIcon from "lucide-svelte/icons/loader-circle";
    import { authenticated } from "$lib/stores/auth"

    let buttonDisabled: boolean = true;
   
    export let data: SuperValidated<Infer<FormSchema>>;

    function validateForm(){
      try {
        formSchema.parse($formData)
        buttonDisabled = false;
      } catch (e: any) {
        buttonDisabled = true;
      }
    }

    function showErrors(form: SuperValidated<Infer<FormSchema>>): void {
      form.errors.email = ["Invalid credentials"];
      form.errors.password = ["Invalid credentials"];
      toast.error("Invalid credentails.", {
        description: "Oops, was it a typo?",
        cancel: {label: "X", onClick: () => {}}
      });
    }
   
    const form = superForm(data, {
      validators: zodClient(formSchema),
      resetForm: false,
      SPA: true,
      delayMs: 1000,
      onChange: () => {validateForm()},
      onUpdate: async ({form}) => {
        try {
          const response = await fetch(`${PUBLIC_BACKEND_BASE_URL}/v1/tokens/authentication`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(form.data),
          });
          const r = await response.json();
          const expires = new Date(r.authentication_token.expiry).toUTCString();
          const token = r.authentication_token.token;
          document.cookie = `token=${token}; expires=${expires};path=/`;
          authenticated.set(true);
          goto("/logbook");
        } catch (e) {
          showErrors(form);
        }
      }
    });
   
    const { form: formData, enhance, delayed } = form;
</script>

<style>
  :global(.spinner) {
		animation: spinner-frames 3s infinite linear;
	}
  @keyframes spinner-frames {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
</style>
   
<form method="POST" use:enhance>
    <Form.Field {form} name="email" class="pb-5">
        <Form.Control let:attrs>
          <Input placeholder="username" {...attrs} bind:value={$formData.email} />
        </Form.Control>
        <Form.FieldErrors />
    </Form.Field>
    <Form.Field {form} name="password" class="pb-5 text-center">
      <Form.Control let:attrs>
        <Input placeholder="Password" type="password" {...attrs} bind:value={$formData.password} />
      </Form.Control>
      <Form.FieldErrors />
  </Form.Field>
  <Form.Button class="w-full" disabled={buttonDisabled || $delayed}>
    {#if $delayed}
      <LoaderCircleIcon class="spinner"/>
    {:else}
      Login
    {/if}
  </Form.Button>
</form>
