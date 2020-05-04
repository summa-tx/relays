<template>
  <v-container fluid class="relay-info-wrapper">
      <v-row justify="center">
        <v-card class="mt-10 mx-5">
          <div class="d-flex flex-no-wrap justify-space-between source-header">
            <div>
              <v-card-title class="headline mt-2 ml-2">
                Cosmos Relay
              </v-card-title>

              <v-card-subtitle class="source-header__subtitle white--text ml-2">
                <div class="d-flex flex-no-wrap">
                  <div class="info__item">Last Connected:</div>
                  <Display-Mins :timestamp="relay.lastComms" />
                </div>

                <div class="d-flex flex-no-wrap">
                  <v-tooltip top nudge-bottom="5">
                    <template v-slot:activator="{ on }">
                      <div class="info__item" v-on="on">Block Difference:</div>
                    </template>
                    <span>The difference in blocks between the relay and the explorer</span>
                  </v-tooltip>
                  <div>{{ bkd.height - currentBlock.height }}</div>
                </div>
              </v-card-subtitle>
            </div>

            <div class="network-type">Localnet</div>
          </div>

          <v-divider/>

          <v-col column class="relay-info">
            <Relay-Info-Item-Title
              title="Best Known Digest:"
              tip="The digest of the best block, updated approximately every 5 blocks"
            ></Relay-Info-Item-Title>

            <Relay-Info-Item
              :height="bkd.height"
              :hash="bkd.hash"
              :time="bkd.time"
              :updated="bkd.updatedAt"
            ></Relay-Info-Item>
          </v-col>

          <v-divider/>

          <v-col column class="relay-info">
            <Relay-Info-Item-Title
              title="Last Reorg Common Ancestor:"
              tip="The latest ancestral block of both the current best known digest and the previous best known digest"
            ></Relay-Info-Item-Title>

            <Relay-Info-Item
              :height="lca.height"
              :hash="lca.hash"
              :time="lca.time"
              :updated="lca.updatedAt"
            ></Relay-Info-Item>
          </v-col>
        </v-card>

<!-- Repeat Section -->
        <v-card class="mt-10 mx-5">
          <div class="d-flex flex-no-wrap justify-space-between source-header">
            <div class="d-block">
              <v-card-title class="headline mt-2 ml-2">
                Bitcoin Explorer <span class="external-source">{{ source }}</span>
              </v-card-title>

              <v-card-subtitle class="source-header__subtitle white--text ml-2">
                <div class="d-flex flex-no-wrap">
                  <div class="info__item">Last Connected:</div>
                  <Display-Mins :timestamp="external.lastComms" />
                </div>

                <div class="d-flex flex-no-wrap">
                  <div class="info__item">Block Changed:</div>
                  <Display-Mins :timestamp="currentBlock.time" />
                </div>
              </v-card-subtitle>
            </div>

            <div class="network-type">Mainnet</div>
          </div>

          <v-divider></v-divider>

          <v-layout column class="relay-info">
            <Relay-Info-Item-Title
              title="Current Block:"
              tip="The most recent block from the explorer"
            ></Relay-Info-Item-Title>

            <Relay-Info-Item
              :height="currentBlock.height"
              :hash="currentBlock.hash"
              :time="currentBlock.time"
              :updated="currentBlock.updatedAt"
            ></Relay-Info-Item>
          </v-layout>

        </v-card>
      </v-row>
  </v-container>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'RelayInfo',

  components: {
    // RelayInfoBanner: () => import(/* webpackChunkName: 'Relay-Info-Banner' */ './Banner'),
    RelayInfoItemTitle: () => import(/* webpackCHunkName: 'Relay-Info-Item-Title' */ './Item-Title'),
    RelayInfoItem: () => import(/* webpackChunkName: 'Relay-Info-Item' */ './Item'),
    DisplayMins: () => import(/* webpackChunkName: 'Display-Mins' */ '@/components/Display-Mins')
  },

  computed: {
    ...mapState({
      currentBlock: state => state.info.currentBlock,
      relay: state => state.relay,
      bkd: state => state.relay.bkd,
      lca: state => state.relay.lca,
      external: state => state.info,
      source: state => state.info.source
    })
  }
}
</script>

<style scoped>
.relay-info-wrapper {
  margin-top: 40px;
}

.relay-info {
  padding: 20px;
}

.info__item {
  width: 150px;
  font-weight: bold;
}

.source-header {
  background-color: teal;
  color: white;
  min-height: 135px;
}

.source-header__subtitle {
  margin-top: 0px;
}

.external-source {
  font-size: 0.65em;
  position: relative;
  top: 4px;
  left: 15px;
}

.network-type {
  border: 1px solid rgba(255,255,255,0.8);
  border-radius: 2px;
  color: rgba(255,255,255,0.9);
  height: 2.3em;
  text-align: center;
  margin: 20px;
  padding: 5px 20px;
  text-transform: uppercase;
  font-size: 0.94em;
  font-weight: 500;
  letter-spacing: 1px;
}
</style>
