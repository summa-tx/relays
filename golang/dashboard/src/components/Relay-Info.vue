<template>
  <v-card v-resize="onResize" class="relay-info">
    <v-card
      class="relay-info__banner"
      tile
      color="teal"
      dark
    >
      <v-layout column>
        <v-layout
          class="relay-info__banner__title"
          row
          justify-space-between
          align-content-center
        >
          <v-layout column>
            <h2>Relay Info</h2>
            <v-layout>
              <p class="mr-2">Health Check:</p>
              <Display-Mins :timestamp="lastComms.relay" />
            </v-layout>
          </v-layout>
          <Net-Type/>
        </v-layout>
      </v-layout>
    </v-card>

    <div class="relay-info__info">
      <v-layout class="relay-info__info__line" column>
        <v-layout>
          <h3 class="relay-info__info__title">Current Block:</h3>
          <v-tooltip top nudge-bottom="5">
            <template v-slot:activator="{ on }">
              <v-icon v-on="on" size="20px">help</v-icon>
            </template>
            <span>The most recent block from the external source</span>
          </v-tooltip>
        </v-layout>

        <v-flex class="relay-info__info__data" row>
          <p>Height: {{ currentBlock.height }}</p>
          <Click-To-Copy :copy-value="currentBlock.height"/>
        </v-flex>

        <v-flex class="relay-info__info__data" row>
          <p>
            <span>Hash: </span>
            <span v-if="windowWidth < 800">{{ currentBlock.hash | crop }}</span>
            <span v-else>{{ currentBlock.hash }}</span>
          </p>
          <Click-To-Copy :copy-value="currentBlock.hash"/>
        </v-flex>

        <v-flex class="relay-info__info__data" row>
          <p>
            <span>Timestamp: </span>
            <span>{{ currentBlock.time | formatTime }}</span>
          </p>
        </v-flex>

        <v-flex class="relay-info__info__data" row>
          <p class="mr-2">Updated:</p>
          <Display-Mins :timestamp="currentBlock.updatedAt" />
        </v-flex>
      </v-layout>

      <v-divider/>

      <v-layout class="relay-info__info__line" column>
        <v-layout>
          <h3 class="relay-info__info__title">Best Known Digest:</h3>
          <v-tooltip top nudge-bottom="5">
            <template v-slot:activator="{ on }">
              <v-icon v-on="on" size="20px">help</v-icon>
            </template>
            <span>The digest of the best block, updated approximately every 5 blocks</span>
          </v-tooltip>
        </v-layout>

        <v-flex class="relay-info__info__data" row>
          <p>Height: {{ bkd.height }}</p>
          <Click-To-Copy :copy-value="bkd.height"/>
        </v-flex>

        <v-flex class="relay-info__info__data" row>
          <p>
            <span>Hash: </span>
            <span v-if="windowWidth < 800">{{ bkd.hash | crop }}</span>
            <span v-else>{{ bkd.hash }}</span>
          </p>
          <Click-To-Copy :copy-value="bkd.hash"/>
        </v-flex>

        <v-flex class="relay-info__info__data" row>
          <p class="mr-2">Verified:</p>
          <Display-Mins :timestamp="bkd.verifiedAt" />
        </v-flex>
      </v-layout>

      <v-divider/>

      <v-layout class="relay-info__info__line" column>
        <v-layout>
          <h3 class="relay-info__info__title">Common Ancestor of Last Reorg:</h3>

          <v-tooltip top nudge-bottom="5">
            <template v-slot:activator="{ on }">
              <v-icon v-on="on" size="20px">help</v-icon>
            </template>
            <span>The latest ancestral block of both the current best known digest and the previous best known digest</span>
          </v-tooltip>
        </v-layout>

        <v-flex class="relay-info__info__data" row>
          <p>Height: {{ lca.height }}</p>
          <Click-To-Copy :copy-value="lca.height"/>
        </v-flex>

        <v-flex class="relay-info__info__data" row>
          <p>
            <span>Hash: </span>
            <span v-if="windowWidth < 800">{{ lca.hash | crop }}</span>
            <span v-else>{{ lca.hash }}</span>
          </p>
          <Click-To-Copy :copy-value="lca.hash"/>
        </v-flex>

        <v-flex class="relay-info__info__data" row>
          <p class="mr-2">Verified:</p>
          <Display-Mins :timestamp="lca.verifiedAt" />
        </v-flex>
      </v-layout>
    </div>
  </v-card>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'Relay',

  components: {
    ClickToCopy: () => import(/* webpackChunkName: 'Click-To-Copy' */ './Click-To-Copy'),
    DisplayMins: () => import(/* webpackChunkName: 'Display-Mins' */ './Display-Mins'),
    NetType: () => import(/* webpackChunkName: 'Net-Type' */ './Net-Type')
  },

  data: () => ({
    windowWidth: Number,
  }),

  computed: {
    ...mapState({
      lastComms: state => state.info.lastComms,
      currentBlock: state => state.info.currentBlock,
      bkd: state => state.info.bkd,
      lca: state => state.info.lca
    })
  },

  mounted () {
    this.onResize()
  },

  methods: {
    onResize () {
      this.windowWidth = window.innerWidth
    }
  },

  filters: {
    crop (str) {
      var first = str.slice(0, 6)
      var last = str.slice((str.length - 6), str.length)
      return `${first} . . . ${last}`
    },

    formatTime (time) {
      const d = new Date(time)
      return `${d.toDateString()} ${d.toTimeString()}`
    }
  }
}
</script>

<style scoped>
.relay-info {
  margin-top: 20px;
}

.relay-info__banner {
  padding: 20px;
}

.relay-info__banner__title {
  margin: 0;
}

.relay-info__info {
  padding: 20px;
}

.relay-info__info__title {
  margin-right: 7px;
  font-weight: 900;
}

.relay-info__info__line {
  margin-top: 30px;
}

.relay-info__info__data {
  font-weight: 400;
  margin-left: 0px;
}
</style>
