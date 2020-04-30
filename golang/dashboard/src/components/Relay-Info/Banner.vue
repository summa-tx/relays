<template>
  <v-card
    class="banner"
    tile
    color="teal"
    dark
  >
    <v-layout column>
      <v-layout
        class="banner__title"
        row
        justify-space-between
        align-content-center
      >
        <v-layout column>
          <h2>Relay Info</h2>
          <v-layout>
            <p class="mr-2">Health Check:</p>
            <Display-Mins :timestamp="lastCommsRelay" />
          </v-layout>
        </v-layout>
        <Net-Type/>
      </v-layout>
    </v-layout>
  </v-card>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: "RelayInfoBanner",

  components: {
    NetType: () => import(/* webpackChunkName: 'Net-Type' */ '../Net-Type')
  },

  computed: {
    ...mapState({
      lastCommsRelay: state => state.info.lastComms.relay,
    })
  }
}
</script>

<style scoped>
.banner {
  padding: 20px;
}

.banner__title {
  margin: 0;
}
</style>
