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
            <p>Height: {{ height }}</p>
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

          <v-flex class="relay__info__info" row>
            <p v-if="verifiedAt === null">Unverified</p>
            <p v-else-if="verifiedAt < 1">Verified: Less than 1 minute ago</p>
            <p v-else>Verified: {{ verifiedAt }} minute<span v-if="verifiedAt > 1">s</span> ago</p>
          </v-flex>
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

    <v-card class="relay__updates">
      <v-layout>
        <p class="relay__updates__title">Relay Health Check:</p>
        <p v-if="lastCommsRelay === null">Not completed</p>
        <p v-else-if="lastCommsRelay < 1">Less than 1 minute ago</p>
        <p v-else>{{ lastCommsRelay }} minute<span v-if="lastCommsRelay > 1">s</span> ago</p>
      </v-layout>
      <v-layout>
        <p class="relay__updates__title">Source:</p>
        <p>{{ source }}</p>
      </v-layout>
      <v-layout>
        <p class="relay__updates__title">Source Health Check:</p>
        <p v-if="lastCommsExternal === null">Not completed</p>
        <p v-else-if="lastCommsExternal < 1">Less than 1 minute ago</p>
        <p v-else>{{ lastCommsExternal }} minute<span v-if="lastCommsExternal > 1">s</span> ago</p>
      </v-layout>
      <v-layout>
        <p class="relay__updates__title">Source Block Changed:</p>
        <p v-if="verifiedAt === null">Unknown</p>
        <p v-else-if="verifiedAt < 1">Less than 1 minute ago</p>
        <p v-else>{{ verifiedAt }} minute<span v-if="verifiedAt > 1">s</span> ago</p>
      </v-layout>
      <v-layout>
        <p class="relay__updates__title">Source Height:</p>
        <p>{{ height || 'Unknown' }}</p>
      </v-layout>
    </v-card>

    <!-- <v-card>
      <v-btn @click="getBKD">Get BKD</v-btn>
      <v-btn @click="getLCA">Get LCA</v-btn>
      <v-btn @click="getCurrentHeight">Get Current Height</v-btn>
    </v-card> -->
  </v-container>
</template>

<script>
import { mapState } from 'vuex'
import { getMinsAgo } from '@/utils/utils'

export default {
  name: 'relay',

  components: {
    ClickToCopy: () => import(/* webpackChunkName: 'Click-To-Copy' */ './Click-To-Copy'),
    NetType: () => import(/* webpackChunkName: 'Net-Type' */ './Net-Type')
  },

  data: () => ({
    windowWidth: Number,
    verifiedAt: null,
    lastCommsExternal: null,
    lastCommsRelay: null
  }),

  mounted () {
    this.onResize()

    // Calculate minutes for health check
    this.healthCheckMins()

    // Updates every minute
    setInterval(() => {
      this.healthCheckMins()
    }, 60000)
  },

  methods: {
    onResize () {
      this.windowWidth = window.innerWidth
    },

    healthCheckMins () {
      this.verifiedAt = this.currentBlock.verifiedAt ? getMinsAgo(this.currentBlock.verifiedAt) : null
      this.lastCommsExternal = this.lastComms.external ? getMinsAgo(this.lastComms.external) : null
      this.lastCommsRelay = this.lastComms.relay ? getMinsAgo(this.lastComms.relay) : null
    }

    // getBKD () {
    //   this.$store.dispatch('relay/getBKD')
    // },

    // getLCA () {
    //   this.$store.dispatch('relay/getLCA')
    // },

    // getCurrentHeight () {
    //   this.$store.dispatch('relay/verifyHeight', this.currentBlock.hash)
    // }
  },

  computed: {
    ...mapState({
      lastComms: state => state.info.lastComms,
      currentBlock: state => state.info.currentBlock,
      height: state => state.info.currentBlock.height,
      relay: state => state.info.relay,
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

.relay__updates {
  margin-top: 50px;
  padding: 20px;
  max-width: 500px;
}

.relay__updates__title {
  width: 200px;
  font-weight: 900;
}

@media (max-width: 800px) {
  .relay {
    padding: 0 20px;
  }
  .relay__updates {
    max-width: none;
  }
}
</style>
