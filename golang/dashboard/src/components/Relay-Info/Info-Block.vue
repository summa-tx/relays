<template>
  <div class="d-flex flex-no-wrap white--text">
    <v-tooltip top nudge-bottom="5">
      <template v-slot:activator="{ on }">
        <div class="item" v-on="on"> {{ label }}: </div>
      </template>
      <span> {{ tip }} </span>
    </v-tooltip>

    <slot></slot>
  </div>
</template>

<script>
export default {
  name: 'InfoBlock',

  props: {
    label: {
      type: String
    },

    tip: {
      type: String
    }
  }
}
</script>

<style scoped>
.item {
  width: 150px;
  font-weight: bold;
}
</style>
