<template>
  <div class="q-pa-md">
    <q-btn label="Add" color="primary" @click="addVsphereHostModel = true" />
    <q-dialog v-model="addVsphereHostModel">
      <q-card style="min-width: 350px">
        <q-form
          @submit="onSubmit"
          @reset="onReset"
          class="q-gutter-md"
        >
          <q-card-section>
            <div class="text-h6">Host</div>
          </q-card-section>

          <q-card-section class="q-pt-none">
            <q-input dense filled v-model="host" type="text" autofocus
                     lazy-rules
                     :rules="[ val => val && val.length > 0 || 'Please enter host address']"/>
          </q-card-section>

          <q-card-section>
            <div class="text-h6">Username</div>
          </q-card-section>

          <q-card-section class="q-pt-none">
            <q-input dense filled v-model="username" type="text"
                     lazy-rules
                     :rules="[ val => val && val.length > 0 || 'Please enter username']"/>
          </q-card-section>

          <q-card-section>
            <div class="text-h6">Password</div>
          </q-card-section>

          <q-card-section class="q-pt-none">
            <q-input dense filled v-model="password" type="password"
                     lazy-rules
                     :rules="[ val => val && val.length > 0 || 'Please enter password']"/>
          </q-card-section>

          <q-card-actions align="right" class="text-primary">
            <q-btn flat label="Cancel" type="reset" v-model="cancel" v-close-popup />
            <q-btn flat label="Add vSphere host" type="submit"/>
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>
    <q-table
      title="vSphere agents"
      :rows="rows"
      :columns="columns"
      row-key="name"
    />
  </div>
</template>

<script>
import { ref } from 'vue'
import { useQuasar } from 'quasar'
import { mapActions } from 'vuex'
import Api from "../../api/Api";

const columns = [
  { name: 'host', align: 'center', label: 'Host', field: 'host', sortable: true, field: row => row.host, format: val => `${val}`},
  { name: 'user', align: 'center', label: 'Username', field: 'user', sortable: true },
  { name: 'session', align: 'center', label: 'Session', field: 'session' },
  { name: 'time', align: 'center', label: 'Started time', field: 'start', sortable: true }
]
const rows = [{}]

export default {
  name: "Vsphere",
    data () {
      return {
        host: '',
        username: '',
        password: '',
        $q : useQuasar()
      }
    },
    methods: {
      ...mapActions('auth', ['addVsphereHost']),
      async onSubmit () {
        const ok = await this.addVsphereHost(this.host,this.username,this.password)
        if (ok) {
          this.q.notify({
            color: 'green-4',
            textColor: 'white',
            icon: 'cloud_done',
            message: 'Submitted'
          })
        }else{
          this.q.notify({
            color: 'red-5',
            textColor: 'white',
            icon: 'warning',
            message: 'Error add Redfish host'
          })
        }
      },
      getAgents() {
        Api.get('/vsphere/agents').then(resp =>
          this.databasesToTableFormat(resp.data)).catch(err => notifications.handleError(err));
      },
      databasesToTableFormat(response) {
        response.forEach(e => {
          this.data.push({
            host: e.host,
            user: e.user,
            session: e.session,
            time: e.time,
          });
        });
      }
    },
    mounted() {
      this.getAgents();
    },
    onReset () {
      this.host.value = null
      this.username.value = null
      this.password.value = null
    },
    setup () {
      const host = ref(null)
      const username = ref(null)
      const password = ref(null)
      const addVsphereHostModel = ref(null)
      return {
        columns,
        rows,
        addVsphereHostModel,
      }
    }
  }

</script>

<style scoped>

</style>
