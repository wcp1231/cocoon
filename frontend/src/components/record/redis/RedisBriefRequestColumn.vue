<template>
  <span class="request-column-tag">{{ request.payload['REDIS_CMD'] }}</span>
  <span class="request-column-info">{{ request.payload['REDIS_KEY'] }}</span>
</template>

<script>
export default {
  name: "RedisBriefRequestColumn",
  props: {
    request: Object,
  },
}
</script>

<style scoped>

</style>