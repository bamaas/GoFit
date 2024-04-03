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

  import { Calendar } from "$lib/components/ui/calendar/index.js";
  import * as Popover from "$lib/components/ui/popover/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import CalendarIcon from "lucide-svelte/icons/calendar";
  import {
    DateFormatter,
    type DateValue,
    getLocalTimeZone,

	CalendarDate,

	today,

	parseDate



  } from "@internationalized/date";
  import { cn } from "$lib/utils.js";
  
  export let data: SuperValidated<Infer<FormSchema>>;

  let uuid: string = "";
  let submitButtonDisabled: boolean = true;

  type CheckIn = {
    uuid: string;
    date: string;
    weight: number;
  };

  function postCheckIn(data: CheckIn){
    if (data.date.split('T').length == 1) data.date = data.date + "T17:34:02.666Z";
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

  function updateCheckIn(data: CheckIn){
    if (data.date.split('T').length == 1) data.date = data.date + "T17:34:02.666Z";
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

  const form = superForm(data, {
    SPA: true,
    validators: zodClient(formSchema),
    onChange() {
      const d: CheckIn = {
        uuid: $formData.uuid,
        date: $formData.date,
        weight: $formData.weight,
      }
      try {
        formSchema.parse(d)
        submitButtonDisabled = false;
      } catch (error) {
        submitButtonDisabled = true;
      }
    },
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

  const { form: formData, enhance } = form;

  const df = new DateFormatter("en-US", {
    dateStyle: "long"
  });
 
  let dateValue: DateValue | undefined = undefined;
  $: dateValue = $formData.date ? parseDate($formData.date.split('T')[0]) : undefined;
  let datePlaceholder: DateValue = today(getLocalTimeZone());

  let popOverOpen: boolean = false;

</script>

<form method="POST" use:enhance>
  <Form.Field {form} name="date" class="pb-5">
    <Form.Control let:attrs>
      <Form.Label>Date</Form.Label>
      <!-- <Input {...attrs} bind:value={$formData.date} /> -->
      <Popover.Root bind:open={popOverOpen}>
        <Popover.Trigger asChild let:builder>
          <Button
            variant="outline"
            class={cn(
              "w-full justify-start text-left font-normal",
              !dateValue && "text-muted-foreground"
            )}
            builders={[builder]}
          >
            <CalendarIcon class="mr-2 h-4 w-4" />
            {dateValue ? df.format(dateValue.toDate(getLocalTimeZone())) : "Pick a date"}
          </Button>
        </Popover.Trigger>
        <Popover.Content class="w-auto p-0">
          <Calendar 
            bind:value={dateValue}
            bind:placeholder={datePlaceholder}
            minValue={new CalendarDate(1900, 1, 1)}
            maxValue={today(getLocalTimeZone())}
            initialFocus
            onValueChange={(v) => {
              if (v) {
                $formData.date = v.toString();
              } else {
                $formData.date = "";
              }
              popOverOpen = false;
            }}
            />
        </Popover.Content>
      </Popover.Root>
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
        <Form.Description>That's how much you weighted that day</Form.Description>
      {/if}
      <Form.FieldErrors />
  </Form.Field>
  <Form.Button class="w-full" disabled={submitButtonDisabled}>
    {#if data.data.uuid == ""}
      Add
    {:else}
      Update
    {/if}
  </Form.Button>
</form>