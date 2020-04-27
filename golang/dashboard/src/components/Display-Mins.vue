<template>
  <div class="display-mins">
    <v-layout>
      <p v-if="minsAgo === null">Not completed</p>
      <p v-else-if="minsAgo < 1">Less than 1 minute ago</p>
      <p v-else>{{ minsAgo }} minute<span v-if="lastCommsExternal > 1">s</span> ago</p>
    </v-layout>
  </div>
</template>

<script>
export default {
  name: 'DisplayMins',

  props: {
    minsAgo: {
      required: true,
      type: Number
    }
  }
}
</script>
