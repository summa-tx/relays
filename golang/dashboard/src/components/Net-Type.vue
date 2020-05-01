<template>
  <h4>
    <span>({{ chainNet }}: {{ netType }}) | (bitcoin: main)</span>
  </h4>
</template>

<script>
import config from '@/config'

export default {
  name: 'NetType',

  data: () => ({
    netType: config.netType,
    chainNet: config.chainNet
  }),
}
</script>
