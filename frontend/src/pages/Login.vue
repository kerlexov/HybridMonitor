<template>
  <div class="fullscreen text-white text-center q-pa-md flex flex-center">
  <q-card-section>
          <q-form class="q-gutter-md" @submit.prevent="submitForm">
            <q-input label="Username" v-model="username">
            </q-input>
            <q-input label="Password" type="password" v-model="password">
            </q-input>
            <div>
              <q-btn class="full-width" color="primary" label="Login" type="submit" rounded></q-btn>
            </div>
          </q-form>
        </q-card-section>
  </div>
</template>

<script>
import { useQuasar } from 'quasar'
import api from 'src/api/Api'

let $q
export default {
  name: 'login',
  data () {
    return {
        username: 'user',
        password: ''
    }
  },
  methods: {
    async submitForm () {
      if (!this.username || !this.password) {
        $q.notify({
          type: 'negative',
          message: 'Please enter username and password'
        })
      } else {
        try {
          let request = {
              username: this.username,
              password: this.password
          };
          api.post('/users/login', request).then(resp => this.saveJwtToCookie(resp)).catch(err => this.handlerError(err));

          const toPath = this.$route.query.to || '/dashboard'
          $q.notify({
            type: 'positive',
            message: 'Login successful'
          })
          this.$router.push(toPath)
        } catch (err) {
          console.log(err)
          if (err.response.data.detail) {
            $q.notify({
              type: 'negative',
              message: err.response.data.detail
            })
          }
        $q.notify({
          type: 'negative',
          message: 'Invalid credentials'
        })
        }
      }
    },

  },
  mounted () {
    $q = useQuasar()
  }
}
</script>

<style scoped>
.wave {
  position: fixed;
  height: 100%;
  left: 0;
  bottom: 0;
  z-index: -1;
}
</style>
