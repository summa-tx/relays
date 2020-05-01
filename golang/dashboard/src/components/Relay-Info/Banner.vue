<template>
  <v-card
    class="banner"
    tile
    color="teal"
    dark
  >
    <v-layout
      class="banner__title"
      row
      justify-space-between
      align-content-center
    >
      <v-layout column>
        <h2>
          <span>{{ chainNet | capitalize }}</span>
          Relay - [ {{  netType }} ]
        </h2>
        <v-layout>
          <p class="mr-2">Last Connected:</p>
          <Display-Mins :timestamp="lastCommsRelay" />
        </v-layout>
      </v-layout>
      <v-layout column>
        <h2>
          Bitcoin - [ main ]
        </h2>
        <External-Info/>
      </v-layout>
    </v-layout>
  </v-card>
</template>

<script>
import { mapState } from 'vuex'
import config from '@/config'

export default {
  name: "RelayInfoBanner",

  components: {
    ExternalInfo: () => import(/* webpackChunkName: 'External-Info' */ '@/components/External-Info'),
    DisplayMins: () => import(/* webpackChunkName: 'Display-Mins' */ '../Display-Mins')
  },

  data: () => ({
    netType: config.netType,
    chainNet: config.chainNet
  }),

  computed: {
    ...mapState({
      lastCommsRelay: state => state.relay.lastComms,
    })
  },

  filters: {
    capitalize (str) {
      return str.charAt(0).toUpperCase() + str.slice(1)
    }
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
