<template>
  <div class="AdminComponent">
    <h3>Commands:</h3>
    <form v-on:submit.prevent>
      <label for="GET-name">Key :</label>
      <input v-model="keymongo" id="form1" type="text" name="name" />
      <input
        v-on:click="
          index = 1;
          viewCommands();
        "
        type="button"
        value="DELETE"
        class="btn btn-outline-primary btn-sm"
      />
    </form>
    <h3></h3>
    <input
      v-on:click="
        index = 2;
        viewCommands();
      "
      type="button"
      value="Reset the Week"
      class="btn btn-danger"
    />
    <input
      v-on:click="
        index = 3;
        viewCommands();
      "
      type="button"
      value="update history"
      class="btn btn-danger"
    />
    <h3></h3>
    <h4 v-if="index != 0">{{ this.confirm }}</h4>
    <input
      v-if="index != 0"
      v-on:click="Commands"
      type="button"
      value="YES"
      class="btn btn-danger"
    />
    <input
      v-if="index != 0"
      v-on:click="index = 0"
      type="button"
      value="NO"
      class="btn btn-danger"
    />

    <h2 v-if="IssueWeekly != null">Issues reported this week :</h2>
    <b-table striped hover :items="IssueWeekly"> </b-table>
    <h2 v-if="history != null">Issues reported history :</h2>
    <b-table
      striped
      hover
      :items="history"
      v-if="history != null"
      :fields="fields"
    ></b-table>
  </div>
</template>

<script>
import Api from '@/services/ServiceJira'

const api = new Api()

export default {
    name: 'AdminComponent',
    data() {
        return {
            history: null,
            index: 0,
            keymongo: '',
            confirm: '',
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
                    sortable: true,
                },
            ],
        }
    },
    mounted() {
        this.GetAll()
        this.Issuehistory()
    },
    methods: {
        Commands() {
            switch (this.index) {
            case 1:
                this.IssueDelete()
                break
            case 2:
                this.mongoReset()
                break
            case 3:
                this.historyUpdate()
                break
            default:
                this.index = 0
                break
            }
            this.index = 0
        },
        viewCommands() {
            switch (this.index) {
            case 1:
                this.confirm = `Delete issue: ${this.keymongo} ?`
                break
            case 2:
                this.confirm = 'Reset the Database Reporting Weekly ? '
                break
            case 3:
                this.confirm = 'Update the database history ?'
                break
            default:
                break
            }
        },
        IssueDelete() {
            api.Delete(this.keymongo)
        },
        async mongoReset() {
            api.Reset().then(() => {
                api.GetAllIssue().then((result) => {
                    this.IssueWeekly = result.issues
                })
            })
        },
        historyUpdate() {
            api.Update().then(() => {
                api.History().then((result) => {
                    this.history = result.issues
                })
            })
        },
        async Issuehistory() {
            api.History().then((result) => {
                this.history = result.issues
            })
        },
        async GetAll() {
            api.GetAllIssue().then((result) => {
                this.IssueWeekly = result.issues
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
