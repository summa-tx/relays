<template>
  <div class="d-flex flex-no-wrap justify-space-between banner-header">
    <div>
      <v-card-title class="headline mt-2 ml-2">
        <slot name="title"></slot>
      </v-card-title>

      <v-card-subtitle class="banner-header__subtitle ml-2">
        <div class="d-flex flex-no-wrap">
          <div class="banner-header__subtitle__label">
            Last Connected:
          </div>
          <Display-Mins :timestamp="connectionType.lastComms" />
        </div>

        <slot name="block"></slot>
      </v-card-subtitle>
    </div>

    <div class="network-type">{{ networkType }}</div>
  </div>
</template>

<script>
export default {
  name: 'InfoBanner',

  components: {
    DisplayMins: () => import(/* webpackChunkName: 'Display-Mins' */ '@/components/Display-Mins')
  },

  props: {
    networkType: {
      type: String
    },

    connectionType: {
      type: Object
    }
  }
}
</script>

<style scoped>
.banner-header {
  background-color: teal;
  color: white;
  min-height: 135px;
}

.banner-header__subtitle {
  margin-top: 0px;
}

.banner-header__subtitle__label {
  width: 150px;
  font-weight: bold;
  color: white;
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
