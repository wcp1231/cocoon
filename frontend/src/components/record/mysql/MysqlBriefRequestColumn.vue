<template>
  <span class="request-column-tag">{{ request.payload["MYSQL_OP_TYPE"] }}</span>
  <span class="request-column-info">{{ request.payload["MYSQL_QUERY"] }}</span>
</template>

<script>
export default {
  name: "MysqlBriefRequestColumn",
  props: {
    request: Object,
  },
}
</script>

<style scoped>

</style>