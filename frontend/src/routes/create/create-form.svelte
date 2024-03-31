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
    

    import CalendarIcon from "lucide-svelte/icons/calendar";
    import {
      type DateValue,
      DateFormatter,
      getLocalTimeZone,
      parseDate,
      CalendarDate,
      today
    } from "@internationalized/date";
    import { cn } from "$lib/utils.js";
    import { Button, buttonVariants } from "$lib/components/ui/button/index.js";
    import { Calendar } from "$lib/components/ui/calendar/index.js";
    import * as Popover from "$lib/components/ui/popover/index.js";
    import { browser } from "$app/environment";
   
    export let data: SuperValidated<Infer<FormSchema>>;

    const df = new DateFormatter("en-US", {
      dateStyle: "long"
    });
    
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
   
    // const proxyDate = dateProxy(form, 'date', { format: 'date' }); 

    const { form: formData, enhance } = form;

    let value: DateValue | undefined;
    $: value = $formData.date ? parseDate($formData.date.toString()) : undefined;
    let placeholder: DateValue = today(getLocalTimeZone());
  </script>
   
  <form method="POST" use:enhance>
    <Form.Field {form} name="date" class="pb-5">
      <Form.Control let:attrs>
        <Form.Label>Date</Form.Label>
        <Popover.Root>
          <Popover.Trigger
            {...attrs}
            class={cn(
              buttonVariants({ variant: "outline" }),
              "w-full justify-start text-left font-normal",
              !value && "text-muted-foreground"
            )}
          >
            {value ? df.format(value.toDate(getLocalTimeZone())) : "Pick a date"}
            <CalendarIcon class="ml-auto h-4 w-4 opacity-50" />
          </Popover.Trigger>
          <Popover.Content class="w-auto p-0" side="top">
            <Calendar
              {value}
              bind:placeholder
              minValue={new CalendarDate(1900, 1, 1)}
              calendarLabel="Date of measurement"
              initialFocus
              onValueChange={(v) => {
                if (v) {
                  $formData.date = v.toString();
                } else {
                  $formData.date = "";
                }
              }}
            />
          </Popover.Content>
        </Popover.Root>
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