<template>
  <div class="FormJiraComponent">
    <form v-on:submit.prevent>
      <label for="GET-name">Key :</label>
      <input v-model="jiraKey" id="form1" type="text" name="name" />
      <input
        v-on:click="Submit"
        type="button"
        value="OK"
        class="btn btn-outline-primary btn-sm"
      />
    </form>
  </div>
</template>

<script>
import Api from '@/services/ServiceJira'

const api = new Api()
export default {
    name: 'FormJiraComponent',
    data() {
        return {
            jiraKey: '',
            r: '',
        }
    },
    methods: {
        async Submit() {
            if (this.jiraKey !== '') {
                api.GetJiraIssue(this.jiraKey).then((result) => {
                    this.r = result
                    this.$emit('onSubmitRes', this.r)
                    this.$emit('onSubmitKey', this.jiraKey)
                })
            }
        },
    },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}

</style>
