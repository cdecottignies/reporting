<template>
  <div class="home">
    <h1>Welcome to your weekly Jira ticket reporting</h1>
    <home-component></home-component>
  </div>
</template>

<script>
// @ is an alias to /src
import HomeComponent from '@/components/HomeComponent.vue'

export default {
    name: 'Home',
    components: {
        HomeComponent,
    },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">

</style>
