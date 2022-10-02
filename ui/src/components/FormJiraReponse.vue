<template>
  <div class="FormJiraResponse">
    <br />
    <h4>{{ res }}</h4>
    <h3 v-if="res != null">Add this issue ?</h3>
    <input
      v-if="res != null"
      v-on:click="AddYes"
      type="button"
      value="YES"
      class="btn btn-danger"
    />
    <input
      v-if="res != null"
      v-on:click="AddNo"
      type="button"
      value="NO"
      class="btn btn-danger"
    />
    <h2></h2>
    <br />
    <H2> Issues Reported this week :</H2>
    <b-table class="table table-striped " striped hover
    :items="IssueWeekly"
    :fields="fields">

    </b-table>
  </div>
</template>

<script>
import Api from '@/services/ServiceJira'

const api = new Api()
export default {
    name: 'FormJiraResponse',
    props: {
        res: Object,
        rep: String,
    },
    data() {
        return {
            IssueWeekly: null,
            fields: [
                {
                    key: 'project',
                    sortable: true,
                },
                {
                    key: 'key',
                    sortable: true,
                },
                {
                    key: 'type',
                    sortable: true,
                },
                {
                    key: 'desc',
                },
                {
                    key: 'assigned',
                    sortable: true,
                },
                {
                    key: 'date',
                },
            ],
        }
    },
    mounted() {
        this.GetAll()
    },
    methods: {
        AddNo() {
            this.$emit('onClickButton')
        },
        async GetAll() {
            api.GetAllIssue().then((result) => {
                this.IssueWeekly = result.issues
            })
        },
        async AddYes() {
            this.loading = true
            api.AddByKey(this.rep).then(() => {
                api.GetAllIssue().then((result) => {
                    this.IssueWeekly = result.issues
                    this.$emit('onClickButton')
                })
            })
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

.btn-danger {
  background-color: #4e00e4 !important;

  border-color: #4e00e4 #4e00e4 hsl(351, 68%, 1.5%);

  color: #fff !important;

  text-shadow: 0 -1px 0 rgba(0, 0, 0, 0.62);
}

.btn-danger:hover,
.btn-danger:focus,
.btn-danger:active,
.btn-danger.active,
.open .dropdown-toggle.btn-danger {
  background-color: #be00a6 !important;

  border-color: #9200be #9200be #be00a6;

  color: #fff !important;

  text-shadow: 0 -1px 0 rgba(0, 0, 0, 0.62);
}

.table {
  color: #4e00e4
}
</style>
