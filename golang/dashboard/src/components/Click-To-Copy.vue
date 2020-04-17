<template>
  <div class="ctc">
    <v-tooltip
      v-model="copied"
      nudge-top="8"
      bottom
      class="ctc__copy"
    >
      <template>
        <v-btn icon color="teal">
          <v-icon @click="handleCopy" id="ctc__icon">
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
    },
    copyText: {
      type: String,
      default: 'Copied'
    }
  },

  data: () => ({
    copied: false
  }),

  methods: {
    handleCopy () {
      const self = this

      if (navigator.clipboard) {
        navigator.clipboard.writeText(this.copyValue).then(function () {
          self.copyText = 'Copied'
          self.copied = true
        }, function () {
          self.copyText = 'Error: Not Copied'
          self.copied = false
        })
      } else {
        self.copyText = 'Error: Not Copied'
        self.copied = false
      }

      setTimeout(() => {
        self.copied = false
      }, 2000)
    }
  }
}
</script>

<style>
.ctc {
  position: relative;
  bottom: 5px;
  margin-left: 10px;
}

/* TODO: BAD. only doing this for now to get it working. fix later. */
#ctc__icon {
  font-size: 16px;
}
</style>
