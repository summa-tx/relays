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
          <h3 class="relay__info__title">Current Block:</h3>
          <v-flex class="relay__info__info" row>
            <p><b>Height:</b> {{ height }}</p>
            <Click-To-Copy :copy-value="height"/>
          </v-flex>
          <v-flex class="relay__info__info" row>
            <Click-To-Copy :copy-value="currentBlock.hash"/>
          </v-flex>
          <v-flex class="relay__info__info" row>
            <p><b>Verified:</b> {{ currentBlock.verifiedAt || 'Unverified' }}</p>
          </v-flex>
        </v-layout>

        <v-divider/>

        <v-layout class="relay__info__line" column>
          <h3 class="relay__info__title">Best Known Digest:</h3>
          <v-flex class="relay__info__info" row>
            <p v-if="windowWidth < 800">{{ relay.bkd | crop }}</p>
            <p v-else>{{ relay.bkd }}</p>
            <Click-To-Copy :copy-value="relay.bkd"/>
          </v-flex>
        </v-layout>

        <v-divider/>

        <v-layout class="relay__info__line" column>
          <h3 class="relay__info__title">Common Ancestor of Last Reorg:</h3>
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
        <p v-if="!lastCommsRelay">Health check not completed</p>
        <p v-else-if="lastCommsRelay <= 1">Less than 1 minute ago</p>
        <p v-else>{{ lastCommsRelay }} minutes ago</p>
      </v-layout>
      <v-layout>
        <p class="relay__updates__title">Source:</p>
        <p>{{ source }}</p>
      </v-layout>
      <v-layout>
        <p class="relay__updates__title">Source Health Check:</p>
        <p v-if="!lastCommsExternal">Source health check not completed</p>
        <p v-else-if="lastCommsExternal <= 1">Less than 1 minute ago</p>
        <p v-else>{{ lastCommsExternal }} minutes ago</p>
      </v-layout>
      <v-layout>
        <p class="relay__updates__title">Source Block Changed:</p>
        <p v-if="!verifiedAt">Unknown</p>
        <p v-else-if="verifiedAt < 1">Less than 1 minute ago</p>
        <p v-else>{{ verifiedAt }} minutes ago</p>
      </v-layout>
      <v-layout>
        <p class="relay__updates__title">Source Height:</p>
        <p>{{ height || 'Unknown' }}</p>
      </v-layout>
    </v-card>

    <v-card>
      <v-btn @click="getBKD">Get BKD</v-btn>
      <v-btn @click="getLCA">Get LCA</v-btn>
      <v-btn @click="getHeight">Get Height</v-btn>
      <v-btn @click="getCurrentHeight">Get Current Height</v-btn>
    </v-card>
  </v-container>
</template>

<script>
import { mapState } from 'vuex'
import { getMinsAgo } from '@/utils/utils'
import { debugButtons } from '@/config'

export default {
  name: 'relay',

  components: {
    ClickToCopy: () => import(/* webpackChunkName: 'Click-To-Copy' */ './Click-To-Copy'),
    NetType: () => import(/* webpackChunkName: 'Net-Type' */ './Net-Type')
  },

  data: () => ({
    windowWidth: Number,
    debugButtons
  }),

  mounted () {
    this.onResize()
    console.log('Last relay communication (Relay Health Check): ', this.lastCommsRelay)
    console.log('Last external communication (Source Health Check): ', this.lastCommsExternal)
    console.log('Source block changed: ', this.verifiedAt)
    console.log('Source height: ', this.height)
    setTimeout(() => {
      console.log('Source Health Check: ', this.lastCommsExternal)
    }, 3000)
  },

  methods: {
    onResize () {
      this.windowWidth = window.innerWidth
    },

    getBKD () {
      this.$store.dispatch('relay/getBKD')
    },

    getLCA () {
      this.$store.dispatch('relay/getLCA')
    },

    getCurrentHeight () {
      this.$store.dispatch('relay/verifyHeight', this.currentBlock.hash)
    },

    getHeight () {
      this.$store.dispatch('relay/verifyHeight', '0000000000000346624ca7ac1dbbc16c5ffd3fa388ce9bdb4627264d117014dc')
    }
  },

  computed: {
    ...mapState({
      lastComms: state => state.info.lastComms,
      currentBlock: state => state.info.currentBlock,
      height: state => state.info.currentBlock.height,
      relay: state => state.info.relay,
      source: state => state.info.source
    }),

    verifiedAt () {
      return getMinsAgo(this.currentBlock.verifiedAt)
    },

    lastCommsExternal () {
      return getMinsAgo(this.lastComms.external)
    },

    lastCommsRelay () {
      return getMinsAgo(this.lastComms.relay)
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
  width: 350px;
  min-width: 350px;
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
