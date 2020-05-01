<template>
  <v-layout column>
    <v-layout>
      <p class="external-info__item mr-2">Source:</p>
      <p>{{ source }}</p>
    </v-layout>
    <v-layout>
      <p class="external-info__item mr-2">Last Connected:</p>
      <Display-Mins :timestamp="lastComms.external" />
    </v-layout>
    <v-layout>
      <p class="external-info__item mr-2">Block Changed:</p>
      <Display-Mins :timestamp="currentBlock.time" />
    </v-layout>
    <v-layout>
      <p class="external-info__item mr-2">Height:</p>
      <p>{{ currentBlock.height || 'Unknown' }}</p>
    </v-layout>
  </v-layout>
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

<style scoped>
/* .external-info {
  margin-top: 40px;
  padding: 20px;
  max-width: 500px;
} */

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
