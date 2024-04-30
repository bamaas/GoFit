<script lang="ts">
  import * as Card from "$lib/components/ui/card/index.js";
  import type { PageData } from "./$types.js";
  import LoginForm from "./login-form.svelte";
  export let data: PageData;
</script>
 
<div class="container items-center py-14 max-w-screen-2xl">
  <Card.Root>
    <Card.Header>
        <Card.Title>Let's get it on</Card.Title>
        <Card.Description>Demo account: hi@gofit.nl | gofit123</Card.Description>
    </Card.Header>
    <Card.Content>
      <LoginForm data={data.form} />
    </Card.Content>
  </Card.Root>
</div>
