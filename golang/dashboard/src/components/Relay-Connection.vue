<template>
  <div>
    <v-alert
      v-if="isConnected"
      dense
      color="teal lighten-4"
      class="alert"
    >
      <div class="alert__content">
        <v-icon class="mr-3">check_circle</v-icon>Connected to the Relay
      </div>
    </v-alert>

    <v-alert
      v-else
      dense
      color="orange lighten-3"
      class="alert"
    >
      <div class="alert__content">
        <v-icon class="mr-3">error</v-icon>Not connected to the Relay. Try refreshing page or hit "Update" button.
      </div>
    </v-alert>
  </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'RelayConnection',

  computed: mapState({
    isConnected: (state) => state.relay.connected
  })
}
</script>

<style scoped>
.alert {
  top: 64px;
  color: #000;
}
.alert__content {
  max-width: 1140px;
  margin: auto;
}
</style>
