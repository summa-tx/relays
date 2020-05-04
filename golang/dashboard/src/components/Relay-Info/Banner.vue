<template>
  <v-card
    class="mx-auto banner"
    color="teal"
    dark
  >
    <v-card-title>
      {{ chainNet }} Relay
    </v-card-title>
    <v-row
      justify="space-between"
    >
      <v-col cols="4">
        <button>{{  netType }}</button>
      </v-col>
    </v-row>

    <v-row>
      <p class="mr-2">Last Connected:</p>
      <Display-Mins :timestamp="lastCommsRelay" />
    </v-row>

  </v-card>
</template>

<script>
import { mapState } from 'vuex'
import config from '@/config'

export default {
  name: 'RelayInfoBanner',

  components: {
    DisplayMins: () => import(/* webpackChunkName: 'Display-Mins' */ '../Display-Mins')
  },

  data: () => ({
    netType: config.netType
  }),

  computed: {
    ...mapState({
      lastCommsRelay: state => state.relay.lastComms,
    }),

    chainNet () {
      return this.capitalize(config.chainNet)
    }
  },

  methods: {
    capitalize (str) {
      return str.charAt(0).toUpperCase() + str.slice(1)
    }
  },

  filters: {
    capitalize (str) {
      return str.charAt(0).toUpperCase() + str.slice(1)
    }
  }
}
</script>

<style scoped>

</style>
