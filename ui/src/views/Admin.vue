<template>
  <div class="Admin">
    <admin-component></admin-component>
  </div>
</template>

<script>
// @ is an alias to /src
import AdminComponent from '@/components/AdminComponent.vue'

export default {
    components: { AdminComponent },
    name: 'Admin',
}
</script>
