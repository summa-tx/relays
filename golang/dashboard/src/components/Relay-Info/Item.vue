<template>
  <v-layout v-resize="onResize" class="info-item" column>
    <v-flex class="info-item__data" row>
      <p>Height: {{ height }}</p>
      <Click-To-Copy :copy-value="height"/>
    </v-flex>

    <v-flex class="info-item__data" row>
      <p>
        <span>Hash: </span>
        <span v-if="windowWidth < 800">{{ hash | crop }}</span>
        <span v-else>{{ hash }}</span>
      </p>
      <Click-To-Copy :copy-value="hash"/>
    </v-flex>

    <v-flex class="info-item__data" row>
      <p>
        <span>Timestamp: </span>
        <span>{{ time | formatTime }}</span>
      </p>
    </v-flex>

    <v-flex class="info-item__data" row>
      <p class="mr-2">Verified:</p>
      <Display-Mins :timestamp="verified" />
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  name: 'RelayInfoItem',

  props: {
    height: {
      required: true,
      type: Number,
      default: undefined
    },
    hash: {
      required: true,
      type: String,
      default: ''
    },
    time: {
      required: true,
      type: String
    },
    verified: {
      required: true,
      type: [Date, String]
    }
  },

  components: {
    ClickToCopy: () => import(/* webpackChunkName: 'Click-To-Copy' */ '../Click-To-Copy'),
    DisplayMins: () => import(/* webpackChunkName: 'Display-Mins' */ '../Display-Mins'),
  },

  mounted () {
    this.onResize()
  },

  data: () => ({
    windowWidth: Number,
  }),

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

    formatTime (dateString) {
      const d = new Date(dateString)
      const date = d.toUTCString().slice(4, d.length)
      return `${date}`
    }
  }
}
</script>

<style scoped>
.info-item__title {
  margin-right: 7px;
  font-weight: 900;
}

.info-item__data {
  font-weight: 400;
  margin-left: 0px;
}
</style>
