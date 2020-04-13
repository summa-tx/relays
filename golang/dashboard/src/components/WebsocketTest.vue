<template>
  <div>
    <v-alert
      v-if="!isConnected"
      dense
      type="error"
      color="orange"
    >Not connected to server. Try refreshing page.</v-alert>
  </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'WebSocketConnection',

  computed: mapState({
    isConnected: (state) => state.isConnected
  })
}
</script>
