<template>
  <v-layout class="item-title">
    <h3 class="item-title__title">{{ title }}</h3>
    <v-tooltip top nudge-bottom="5">
      <template v-slot:activator="{ on }">
        <v-icon v-on="on" size="20px">help</v-icon>
      </template>
      <span>{{ tip }}</span>
    </v-tooltip>
  </v-layout>
</template>

<script>
export default {
  name: 'InfoItemTitle',

  props: {
    title: {
      required: true,
      type: String,
      default: ''
    },
    tip: {
      required: true,
      type: String,
      default: ''
    }
  }
}
</script>

<style scoped>
.item-title {
  margin-bottom: 10px;
}

.item-title__title {
  margin-right: 7px;
  font-weight: 900;
}
</style>
