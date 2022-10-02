<template>
  <div class="User">

    <h1>Reporting Jira ticket</h1>
    <form-jira-component @onSubmitRes='validateRes' @onSubmitKey='valueKey'></form-jira-component>
    <form-jira-reponse :res="Response" @onClickButton='valueNull' :rep='key'></form-jira-reponse>
  </div>
</template>

<script>
// @ is an alias to /src
import FormJiraComponent from '@/components/FormJiraComponent.vue'
import FormJiraReponse from '@/components/FormJiraReponse.vue'

export default {
    data() {
        return {
            Response: null,
            key: '',
        }
    },
    name: 'User',
    components: {
        FormJiraComponent,
        FormJiraReponse

    },
    methods: {
        valueKey(val) {
            this.key = val
        },
        validateRes(r) {
            this.Response = r
        },
        valueNull() {
            this.Response = null
        }
    }
}
</script>
