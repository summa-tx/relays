<template>
  <div class="display-mins">
    <span v-if="minsAgo > 1" class="white--text">{{ minsAgo }} min<span v-if="minsAgo > 1">s</span> ago</span>
    <span v-else-if="minsAgo >= 0" class="white--text">
      Less than 1 min ago
    </span>
    <span v-else-if="!minsAgo"></span>
  </div>
</template>

<script>
import { getMinsAgo } from '@/utils/utils'

export default {
  name: 'DisplayMins',

  props: {
    timestamp: {
      required: true,
      type: [Date, String],
      default: ''
    }
  },

  computed: {
    minsAgo () {
      return getMinsAgo(this.timestamp)
    }
  }
}
</script>
