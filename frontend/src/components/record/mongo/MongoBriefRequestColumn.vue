<template>
  <span class="request-column-tag">{{ request.payload["MONGO_OP_TYPE"] }}</span>
  <span class="request-column-info">{{ `${JSON.stringify(request.payload["MONGO_MESSAGE"])}` }}</span>
</template>

<script>
export default {
  name: "MongoBriefRequestColumn",
  props: {
    request: Object,
  },
}
</script>

<style scoped>

</style>