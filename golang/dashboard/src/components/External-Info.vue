<template>
  <v-card class="external-info">
    <h3>External Info</h3>
    <v-divider class="external-info__divider" />
    <v-layout>
      <p class="external-info__item">Source:</p>
      <p>{{ source }}</p>
    </v-layout>
    <v-layout>
      <p class="external-info__item">Health Check:</p>
      <p v-if="lastCommsExternal === null">Not completed</p>
      <p v-else-if="lastCommsExternal < 1">Less than 1 minute ago</p>
      <p v-else>{{ lastCommsExternal }} minute<span v-if="lastCommsExternal > 1">s</span> ago</p>
    </v-layout>
    <v-layout>
      <p class="external-info__item">Block Changed:</p>
      <p v-if="verifiedAt === null">Unknown</p>
      <p v-else-if="verifiedAt < 1">Less than 1 minute ago</p>
      <p v-else>{{ verifiedAt }} minute<span v-if="verifiedAt > 1">s</span> ago</p>
    </v-layout>
    <v-layout>
      <p class="external-info__item">Height:</p>
      <p>{{ currentBlock.height || 'Unknown' }}</p>
    </v-layout>
  </v-card>
</template>

<script>
import { mapState } from 'vuex'
import { getMinsAgo } from '@/utils/utils'

export default {
  name: 'ExternalInfo',

  computed: {
    ...mapState({
      lastComms: state => state.info.lastComms,
      currentBlock: state => state.info.currentBlock,
      source: state => state.info.source,
      verifiedAt: state => state.info.minsAgo.currentBlockVerified,
      lastCommsExternal: state => state.info.minsAgo.sourceHealthCheck
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
.external-info {
  margin-top: 40px;
  padding: 20px;
  max-width: 500px;
}

.external-info__divider {
  margin: 10px 0 20px 0;
}

.external-info__item {
  width: 200px;
  font-weight: 900;
}

@media (max-width: 800px) {
  .external-info {
    max-width: none;
  }
}
</style>
