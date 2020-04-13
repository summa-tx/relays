<template>
  <v-app>
    <v-app-bar app>
      <v-layout
        justify-space-between
        align-center
        class="nav__content"
      >
        <v-toolbar-title class="headline text-uppercase">
          <v-img
            alt="Summa"
            src="./assets/Summa-Logo.png"
            height="42"
            width="35"
          ></v-img>
        </v-toolbar-title>

        <v-spacer/>

        <div class="nav__title">
          <h1>Bitcoin Relay</h1>
        </div>

        <v-spacer/>

      </v-layout>
    </v-app-bar>

    <v-content>
      <Websocket-Test/>
      <Relay-Info/>
    </v-content>
  </v-app>
</template>

<script>
import RelayInfo from './components/Relay-Info'
import WebsocketTest from './components/WebsocketTest'

export default {
  name: 'OperatedRelayDashboard',

  metaInfo: {
    title: 'Operated Bitcoin Relay',
    meta: [
      { 'http-equiv': 'Content-Type', content: 'text/html; charset=utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { name: 'description', content: 'Operated Bitcoin Relay Dashboard' }
    ]
  },

  components: {
    RelayInfo,
    WebsocketTest
  },

  mounted () {
    // Get relay info - bkd, lrca
    this.getRelayInfo()
    this.getExternalInfo()
    // Get external info and set it in the store, start polling
    setInterval(this.getExternalInfo, 120000)
    setInterval(this.getRelayInfo, 60000)
    setInterval(this.updateNow, 65000)
  },

  methods: {
    getRelayInfo () {
      console.log('Getting relay info')
      this.$socket.emit('get_bkd')
      this.$socket.emit('get_lca')
    },

    updateNow () {
      this.$store.dispatch('info/updateNow')
    },

    getExternalInfo () {
      this.$store.dispatch('info/getExternalInfo')
    }
  }
}
</script>

<style>
.nav__content {
  max-width: 1200px;
  margin: auto;
}
.nav__title {
  font-weight: 500;
  font-size: 0.8em;
}
</style>
