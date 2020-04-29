<template>
  <v-card class="external-info">
    <h3>External Info</h3>
    <v-divider class="external-info__divider" />
    <v-layout>
      <p class="external-info__item mr-2">Source:</p>
      <p>{{ source }}</p>
    </v-layout>
    <v-layout>
      <p class="external-info__item mr-2">Health Check:</p>
      <Display-Mins :timestamp="lastComms.external" />
    </v-layout>
    <v-layout>
      <p class="external-info__item mr-2">Block Changed:</p>
      <Display-Mins :timestamp="currentBlock.verifiedAt" />
    </v-layout>
    <v-layout>
      <p class="external-info__item mr-2">Height:</p>
      <p>{{ currentBlock.height || 'Unknown' }}</p>
    </v-layout>
  </v-card>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'ExternalInfo',

  components: {
    DisplayMins: () => import(/* webpackChunkName: 'Display-Mins' */ './Display-Mins')
  },

  computed: {
    ...mapState({
      lastComms: state => state.info.lastComms,
      currentBlock: state => state.info.currentBlock,
      source: state => state.info.source
    })
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
