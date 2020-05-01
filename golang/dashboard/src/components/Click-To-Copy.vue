<template>
  <div class="ctc" @mouseleave="resetCopyText">
    <v-tooltip
      nudge-top="8"
      bottom
      class="ctc__copy"
      content-class="tooltip-text"
    >
      <template v-slot:activator="{ on }">
        <v-btn v-on="on" icon color="teal">
          <v-icon @click="handleCopy" size="16px">
            content_copy
          </v-icon>
        </v-btn>
      </template>
      <span>{{ copyText }}</span>
    </v-tooltip>
  </div>
</template>

<script>
export default {
  name: 'ClickToCopy',

  props: {
    copyValue: {
      required: true,
      type: [String, Number, undefined],
      default: ''
    }
  },

  data: () => ({
    copyText: 'Copy'
  }),

  methods: {
    handleCopy () {
      if (navigator.clipboard) {
        navigator.clipboard.writeText(this.copyValue).then(() => {
          this.copyText = 'Copied'
        }, function () {
          this.copyText = 'Error: Not Copied'
        })
      } else if (document.execCommand && document.body.createTextRange) {
        const range = document.body.createTextRange()
        range.moveToElementText(this.copyText)
        range.select()
        document.execCommand('copy')
      } else {
        this.copyText = 'Error: Not Copied'
      }
    },

    resetCopyText () {
      setTimeout(() => {
        this.copyText = 'Copy'
      }, 500)
    }
  }
}
</script>

<style scoped>
.ctc {
  position: relative;
  bottom: 5px;
  margin-left: 10px;
}

.tooltip-text {
  font-size: 12px;
  padding: 2px 10px;
}
</style>
