<template>
  <v-container
    v-resize="onResize"
    class="relay"
  >
    <v-card class="relay__card">
      <v-toolbar color="teal" dark>
        <v-toolbar-title>Relay Info</v-toolbar-title>
        <v-spacer/>
        <Net-Type></Net-Type>

        <!-- <v-spacer></v-spacer>
        <v-btn icon>
          <v-icon>mdi-dots-vertical</v-icon>
        </v-btn> -->
      </v-toolbar>

      <div class="relay__info">
        <v-layout class="relay__info__line" column>
          <v-layout>
            <h3 class="relay__info__title">Current Block:</h3>
            <v-tooltip top nudge-bottom="5">
              <template v-slot:activator="{ on }">
                <v-icon v-on="on" size="20px">help</v-icon>
              </template>
              <span>Something here</span>
            </v-tooltip>
          </v-layout>

          <v-flex class="relay__info__info" row>
            <p>Height: {{ currentBlock.height }}</p>
            <Click-To-Copy :copy-value="height"/>
          </v-flex>

          <v-flex class="relay__info__info" row>
            <p>
              <span>Hash: </span>
              <span v-if="windowWidth < 800">{{ currentBlock.hash | crop }}</span>
              <span v-else>{{ currentBlock.hash }}</span>
            </p>
            <Click-To-Copy :copy-value="currentBlock.hash"/>
          </v-flex>

          <!-- <v-flex class="relay__info__info" row>
            TODO: Fix verifiedAt
            <p v-if="verifiedAt === null">Unverified</p>
            <p v-else-if="verifiedAt < 1">Verified: Less than 1 minute ago</p>
            <p v-else>Verified: {{ verifiedAt }} minute<span v-if="verifiedAt > 1">s</span> ago</p>
          </v-flex> -->
        </v-layout>

        <v-divider/>

        <v-layout class="relay__info__line" column>
          <v-layout>
            <h3 class="relay__info__title">Best Known Digest:</h3>
            <v-tooltip top nudge-bottom="5">
              <template v-slot:activator="{ on }">
                <v-icon v-on="on" size="20px">help</v-icon>
              </template>
              <span>Something here</span>
            </v-tooltip>
          </v-layout>

          <v-flex class="relay__info__info" row>
            <p v-if="windowWidth < 800">{{ relay.bkd | crop }}</p>
            <p v-else>{{ relay.bkd }}</p>
            <Click-To-Copy :copy-value="relay.bkd"/>
          </v-flex>
        </v-layout>

        <v-divider/>

        <v-layout class="relay__info__line" column>
          <v-layout>
            <h3 class="relay__info__title">Common Ancestor of Last Reorg:</h3>

            <v-tooltip top nudge-bottom="5">
              <template v-slot:activator="{ on }">
                <v-icon v-on="on" size="20px">help</v-icon>
              </template>
              <span>Something here</span>
            </v-tooltip>
          </v-layout>
          <v-flex class="relay__info__info" row>
            <p v-if="windowWidth < 800">{{ relay.lca | crop }}</p>
            <p v-else>{{ relay.lca }}</p>
            <Click-To-Copy :copy-value="relay.lca"/>
          </v-flex>
        </v-layout>
      </div>
    </v-card>

    <Relay-Health-Check />
  </v-container>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'Relay',

  components: {
    RelayHealthCheck: () => import(/* webpackChunkName: 'Relay-Health-Check' */ './Relay-Health-Check'),
    ClickToCopy: () => import(/* webpackChunkName: 'Click-To-Copy' */ './Click-To-Copy'),
    NetType: () => import(/* webpackChunkName: 'Net-Type' */ './Net-Type')
  },

  data: () => ({
    windowWidth: Number,
    verifiedAt: null
  }),

  computed: {
    ...mapState({
      currentBlock: state => state.info.currentBlock,
      relay: state => state.info.relay
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
    }
  }
}
</script>

<style scoped>
.relay {
  max-width: 1264px;
  padding: 0 60px;
}

.relay__card {
  margin-top: 50px;
}

.relay__info {
  padding: 20px;
}

.relay__info__title {
  margin-right: 7px;
  font-weight: 900;
}

.relay__info__line {
  margin-top: 30px;
}

.relay__info__info {
  font-weight: 400;
  margin-left: 0px;
}

@media (max-width: 800px) {
  .relay {
    padding: 0 20px;
  }
}
</style>
