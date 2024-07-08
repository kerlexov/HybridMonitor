<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
        />

        <q-toolbar-title>HybridMonitor</q-toolbar-title>

        <div>
          <q-btn stretch flat to="/login" v-if="!isAuthenticated">Login</q-btn>
          <q-btn stretch flat @click="logout" v-else>Logout</q-btn>
        </div>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered>
      <q-list>
        <q-item-label
          header
        >
          Essential Links
        </q-item-label>

        <EssentialLink
          v-for="link in essentialLinks"
          :key="link.title"
          v-bind="link"
        />
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>


<script>

import { defineComponent, ref } from "vue";
import { mapGetters } from 'vuex'
import EssentialLink from 'components/EssentialLink.vue'

const essentialLinks = [
  {
    title: 'Dashboard',
    icon: 'home',
    link: "/#/dashboard"
  },
  {
    title: 'Redfish',
    caption: 'Redfish',
    icon: '',
    link: '/#/redfish'
  },
  {
    title: 'vSphere',
    caption: '',
    icon: '',
    link: '/#/vsphere'
  },
]

export default defineComponent({
  name: "MainLayout",

  components: {
    EssentialLink
  },
  setup() {
    const leftDrawerOpen = ref(false);

    return {
      essentialLinks,
      leftDrawerOpen,
      toggleLeftDrawer() {
        leftDrawerOpen.value = !leftDrawerOpen.value;
      },
    };
  },
  methods: {
    logout () {
      this.$router.push('/')
    }
  },
  computed: {
    ...mapGetters('auth', ['isAuthenticated'])
  }
});
</script>
