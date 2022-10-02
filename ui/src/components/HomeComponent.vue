<template>
  <div class="HomeComponent">
    <h1>{{ msg }}</h1>
    <input type="button" value="Connexion" class="btn btn-danger" />
    <input type="button" value="Deconnexion" class="btn btn-danger" />
  </div>
</template>

<script>
export default {
    name: 'HomeComponent',
    props: {
        msg: String,
    },
}
</script>
