<script lang="ts">
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
    import * as Select from "$lib/components/ui/select";
	  import { request } from '$lib/functions/request';
    import { capitalizeFirstLetter } from "$lib/functions/string";

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
      form.errors.goal = ["Invalid goal"];
      toast.error("Oops!", {
        description: "Something went wrong.",
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
          await request(`${PUBLIC_BACKEND_BASE_URL}/v1/users/me`, {
            method: "PUT",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(form.data),
          });
        } catch (e) {
          showErrors(form);
        }
      }
    });
  
    const { form: formData, enhance, delayed } = form;

    $: selectedGoal = $formData.goal
    ? {
        label: capitalizeFirstLetter($formData.goal),
        value: $formData.goal,
      }
    : undefined;

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
  <Form.Field {form} name="goal" class="pb-5">
      <Form.Control let:attrs>
        <Form.Label>Goal</Form.Label>
        <Select.Root
        selected={selectedGoal}
        onSelectedChange={(v) => {
          v && ($formData.goal = v.value);
        }}
        >
          <Select.Trigger class="w-full mb-3">
            <Select.Value placeholder="Goal" />
          </Select.Trigger>
          <Select.Content>
            <Select.Item value="cut">Cut</Select.Item>
            <Select.Item value="bulk">Bulk</Select.Item>
            <Select.Item value="maintain">Maintain</Select.Item>
          </Select.Content>
        </Select.Root>
        <Input style="display: none;" {...attrs} bind:value={$formData.goal} />
      </Form.Control>
      <Form.Description>What is the goal you try to archive?</Form.Description>
      <Form.FieldErrors />
  </Form.Field>
<Form.Button class="w-full" disabled={buttonDisabled || $delayed}>
  {#if $delayed}
    <LoaderCircleIcon class="spinner"/>
  {:else}
    Save profile
  {/if}
</Form.Button>
</form>
