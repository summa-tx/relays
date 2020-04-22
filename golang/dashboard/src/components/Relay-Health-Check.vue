<template>
  <v-card class="relay-updates">
    <v-layout>
      <p class="relay-updates__title">Relay Health Check:</p>
      <p v-if="lastCommsRelay === null">Not completed</p>
      <p v-else-if="lastCommsRelay < 1">Less than 1 minute ago</p>
      <p v-else>{{ lastCommsRelay }} minute<span v-if="lastCommsRelay > 1">s</span> ago</p>
    </v-layout>
    <v-layout>
      <p class="relay-updates__title">Source:</p>
      <p>{{ source }}</p>
    </v-layout>
    <v-layout>
      <p class="relay-updates__title">Source Health Check:</p>
      <p v-if="lastCommsExternal === null">Not completed</p>
      <p v-else-if="lastCommsExternal < 1">Less than 1 minute ago</p>
      <p v-else>{{ lastCommsExternal }} minute<span v-if="lastCommsExternal > 1">s</span> ago</p>
    </v-layout>
    <v-layout>
      <p class="relay-updates__title">Source Block Changed:</p>
      <p v-if="verifiedAt === null">Unknown</p>
      <p v-else-if="verifiedAt < 1">Less than 1 minute ago</p>
      <p v-else>{{ verifiedAt }} minute<span v-if="verifiedAt > 1">s</span> ago</p>
    </v-layout>
    <v-layout>
      <p class="relay-updates__title">Source Height:</p>
      <p>{{ currentBlock.height || 'Unknown' }}</p>
    </v-layout>
  </v-card>
</template>

<script>
import { mapState } from 'vuex'
import { getMinsAgo } from '@/utils/utils'

export default {
  name: 'RelayHealthCheck',

  components: {
    ClickToCopy: () => import(/* webpackChunkName: 'Click-To-Copy' */ './Click-To-Copy'),
    NetType: () => import(/* webpackChunkName: 'Net-Type' */ './Net-Type')
  },

  data: () => ({
    verifiedAt: null,
    lastCommsExternal: null,
    lastCommsRelay: null
  }),

  computed: {
    ...mapState({
      lastComms: state => state.info.lastComms,
      currentBlock: state => state.info.currentBlock,
      source: state => state.info.source
    })
  },

  watch: {
    lastComms: {
      handler: function () {
        console.log('Updated info')
        this.healthCheckMins()
      },
      deep: true
    }
  },

  mounted () {
    // Calculate minutes for health check
    this.healthCheckMins()

    // Updates every minute
    setInterval(() => {
      this.healthCheckMins()
    }, 60000)
  },

  methods: {
    healthCheckMins () {
      this.verifiedAt = this.currentBlock.verifiedAt ? getMinsAgo(this.currentBlock.verifiedAt) : null
      this.lastCommsExternal = this.lastComms.external ? getMinsAgo(this.lastComms.external) : null
      this.lastCommsRelay = this.lastComms.relay ? getMinsAgo(this.lastComms.relay) : null
    }
  }
}
</script>

<style>
.relay-updates {
  margin-top: 50px;
  padding: 20px;
  max-width: 500px;
}

.relay-updates__title {
  width: 200px;
  font-weight: 900;
}

@media (max-width: 800px) {
  .relay-updates {
    max-width: none;
  }
}
</style>
