<template>
  <div class="display-mins">
    <v-layout>
      <p v-if="minsAgo === null">Not completed</p>
      <p v-else-if="minsAgo < 1">Less than 1 minute ago</p>
      <p v-else>{{ minsAgo }} minute<span v-if="minsAgo > 1">s</span> ago</p>
    </v-layout>
  </div>
</template>

<script>
import { getMinsAgo } from '@/utils/utils'

export default {
  name: 'DisplayMins',

  props: {
    timestamp: {
      required: true,
      type: Date,
      default: null
    }
  },

  data: () => ({
    minsAgo: null
  }),

  mounted () {
    this.updateMinsAgo()
    setInterval(() => {
      this.updateMinsAgo()
    }, 10000)
  },

  methods: {
    updateMinsAgo () {
      this.minsAgo = this.timestamp ? getMinsAgo(this.timestamp) : null
    }
  }
}
</script>
