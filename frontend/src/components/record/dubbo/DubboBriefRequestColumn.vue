<template>
  <span v-if="request.meta['HEARTBEAT']" class="request-column-tag">HeartBeat</span>
  <template v-else>
    <span class="request-column-tag">Invoke</span>
    <span class="request-column-info">{{ `${request.header['target']}#${request.header['method']}` }}</span>
  </template>
</template>

<script>
export default {
  name: "DubboBriefRequestColumn",
  props: {
    request: Object,
  },
}
</script>

<style scoped>

</style>