<template>
  <div class="relay-updates">
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
  </div>
</template>

<script>
import { mapState } from 'vuex'
import { getMinsAgo } from '@/utils/utils'

export default {
  name: 'RelayHealthCheck',

  computed: {
    ...mapState({
      lastComms: state => state.info.lastComms,
      currentBlock: state => state.info.currentBlock,
      source: state => state.info.source,
      verifiedAt: state => state.info.minsAgo.currentBlockVerified,
      lastCommsExternal: state => state.info.minsAgo.sourceHealthCheck,
      lastCommsRelay: state => state.info.minsAgo.relayHealthCheck
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
      const currentBlockVerified = this.currentBlock.verifiedAt ? getMinsAgo(this.currentBlock.verifiedAt) : null
      const relayHealthCheck = this.lastComms.relay ? getMinsAgo(this.lastComms.relay) : null
      const sourceHealthCheck = this.lastComms.external ? getMinsAgo(this.lastComms.external) : null

      this.$store.dispatch('info/setMinsAgo', {
        currentBlockVerified,
        relayHealthCheck,
        sourceHealthCheck
      })
    }
  }
}
</script>

<style>
.relay-updates {
  margin-top: 20px;
  /* max-width: 500px; */
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
