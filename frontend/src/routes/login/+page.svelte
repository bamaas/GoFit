<script lang="ts">
  import * as Card from "$lib/components/ui/card/index.js";
  import type { PageData } from "./$types.js";
  import LoginForm from "./login-form.svelte";
  export let data: PageData;
</script>
<div 
class="bg-cover bg-center bg-no-repeat bg-cover h-screen w-screen items-center justify-center flex flex-col"
style="background-image:linear-gradient(rgba(0, 0, 0, 0.55), rgba(0, 0, 0, 0.55)), url(https://images.unsplash.com/photo-1605296867304-46d5465a13f1?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D);"
>
    <Card.Root class="lg:w-1/3 w-11/12 bg-opacity-80">
      <Card.Header>
          <Card.Title class="text-center text-2xl tracking-tight">GoFit</Card.Title>
          <Card.Description class="text-center text-sm text-muted-foreground">Enter your credentials to login</Card.Description>
      </Card.Header>
      <Card.Content>
        <LoginForm data={data.form} />
      </Card.Content>
    </Card.Root>
</div>