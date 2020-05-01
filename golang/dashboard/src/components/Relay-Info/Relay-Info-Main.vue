<template>
  <v-card class="mt-5">
    <Relay-Info-Banner />

    <v-layout column class="relay-info">
      <Relay-Info-Item-Title
        title="Current Block:"
        tip="The most recent block from the external source"
      ></Relay-Info-Item-Title>

      <Relay-Info-Item
        :height="currentBlock.height"
        :hash="currentBlock.hash"
        :time="currentBlock.time"
        :updated="currentBlock.updatedAt"
      ></Relay-Info-Item>
    </v-layout>

    <v-divider/>

    <v-layout column class="relay-info">
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
    </v-layout>

    <v-divider/>

    <v-layout column class="relay-info">
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
    </v-layout>
  </v-card>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'RelayInfo',

  components: {
    RelayInfoBanner: () => import(/* webpackChunkName: 'Relay-Info-Banner' */ './Banner'),
    RelayInfoItemTitle: () => import(/* webpackCHunkName: 'Relay-Info-Item-Title' */ './Item-Title'),
    RelayInfoItem: () => import(/* webpackChunkName: 'Relay-Info-Item' */ './Item')
  },

  computed: {
    ...mapState({
      currentBlock: state => state.info.currentBlock,
      bkd: state => state.info.bkd,
      lca: state => state.info.lca
    })
  }
}
</script>

<style scoped>
.relay-info {
  margin-top: 30px;
  padding: 20px;
}
</style>
