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

        <v-btn @click="updateAll">UPDATE</v-btn>

      </v-layout>
    </v-app-bar>

    <Relay-Connection />

    <v-content>
      <v-container
        v-resize="onResize"
        class="relay"
      >
        <Relay-Info />
        <ExternalInfo />
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
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
    RelayInfo: () => import(/* webpackChunkName: 'Relay-Info' */ './components/Relay-Info'),
    RelayConnection: () => import(/* webpackChunkName: 'Relay-Connection' */ './components/Relay-Connection'),
    ExternalInfo: () => import(/* webpackChunkName: 'External-Info' */ './components/External-Info')
  },

  mounted () {
    // Get relay info - best know digest(BKD), last common ancestor(LCA)
    this.getRelayInfo()
    this.getExternalInfo()
    // Get external info and set it in the store, start polling
    setInterval(this.getExternalInfo, 120000)
    setInterval(this.getRelayInfo, 60000)
  },

  methods: {
    getRelayInfo () {
      console.log('Getting relay info')
      this.$store.dispatch('relay/getBKD')
      this.$store.dispatch('relay/getLCA')
    },

    getExternalInfo () {
      this.$store.dispatch('info/getExternalInfo')
    },

    updateAll () {
      this.getRelayInfo()
      this.getExternalInfo()
    },

    onResize () {
      this.windowWidth = window.innerWidth
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

.relay {
  max-width: 1264px;
  padding: 60px;
}
</style>
