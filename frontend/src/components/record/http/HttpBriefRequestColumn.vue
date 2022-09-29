<template>
  <span class="request-column-tag">{{ request.payload['HTTP_REQUEST']['Method'] }}</span>
  <span class="request-column-info">{{ request.payload['HTTP_REQUEST']['Host'] +request.payload['HTTP_REQUEST']['URL'] }}</span>
</template>

<script>
export default {
  name: "HttpBriefRequestColumn",
  props: {
    request: Object,
  },
}
</script>

<style scoped>

</style>