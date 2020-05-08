<template>
  <v-container fluid class="info-wrapper">
      <v-row justify="center">
        <v-card class="mt-10 mx-5">
          <Info-Banner
            network-type="localnet"
            :connection-type="relay"
          >
            <template v-slot:title>
              Cosmos Relay
            </template>

            <template v-slot:block>
              <Info-Block
                label="Block Difference"
                tip="The difference in blocks between the relay and the explorer"
              >
                <div>{{ blockDifference }}</div>
              </Info-Block>
            </template>
          </Info-Banner>

          <v-divider/>

          <v-col column class="info-item">
            <Info-Item-Title
              title="Best Known Digest:"
              tip="The digest of the best block, updated approximately every 5 blocks"
            ></Info-Item-Title>

            <Info-Item
              :height="bkd.height"
              :hash="bkd.hash"
              :time="bkd.time"
            ></Info-Item>
          </v-col>

          <v-divider/>

          <v-col column class="info-item">
            <Info-Item-Title
              title="Last Reorg Common Ancestor:"
              tip="The latest ancestral block of both the current best known digest and the previous best known digest"
            ></Info-Item-Title>

            <Info-Item
              :height="lca.height"
              :hash="lca.hash"
              :time="lca.time"
            ></Info-Item>
          </v-col>
        </v-card>

<!-- External Source -->
        <v-card class="mt-10 mx-5">
          <Info-Banner
            network-type="mainnet"
            :connection-type="external"
          >
            <template v-slot:title>
              Bitcoin Explorer <span class="external-source">{{ source }}</span>
            </template>

            <template v-slot:block>
              <Info-Block
                label="Block Changed"
                tip="Time since the block changed"
              >
                <Display-Mins :timestamp="currentBlock.time" />
              </Info-Block>
            </template>
          </Info-Banner>

          <v-divider></v-divider>

          <v-layout column class="info-item">
            <Info-Item-Title
              title="Current Block:"
              tip="The most recent block from the explorer"
            ></Info-Item-Title>

            <Info-Item
              :height="currentBlock.height"
              :hash="currentBlock.hash"
              :time="currentBlock.time"
            ></Info-Item>
          </v-layout>

        </v-card>
      </v-row>
  </v-container>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'InfoMain',

  components: {
    DisplayMins: () => import(/* webpackChunkName: 'Display-Mins' */ '@/components/Display-Mins'),
    InfoBanner: () => import(/* webpackChunkName: 'Info-Banner' */ './Info-Banner'),
    InfoItemTitle: () => import(/* webpackCHunkName: 'Info-Item-Title' */ './Item-Title'),
    InfoItem: () => import(/* webpackChunkName: 'Info-Item' */ './Item'),
    InfoBlock: () => import(/* webpackChunkName: 'Info-Block' */ './Info-Block')
  },

  computed: {
    ...mapState({
      currentBlock: state => state.external.currentBlock,
      relay: state => state.relay,
      bkd: state => state.relay.bkd,
      lca: state => state.relay.lca,
      external: state => state.external,
      source: state => state.external.source
    }),

    blockDifference () {
      return this.relay.bkd.height - this.currentBlock.height
    }
  }
}
</script>

<style scoped>
.info-wrapper {
  margin-top: 40px;
}

.info-item {
  padding: 20px;
}

.external-source {
  font-size: 0.65em;
  position: relative;
  top: 4px;
  left: 15px;
}
</style>
