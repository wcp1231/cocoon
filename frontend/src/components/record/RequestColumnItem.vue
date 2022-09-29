<template>
  <div class="request-column">
    <dubbo-brief-request-column v-if="protocol === 'Dubbo'" :request="request" />
    <http-brief-request-column v-else-if="protocol === 'HTTP'" :request="request" />
    <redis-brief-request-column v-else-if="protocol === 'Redis'" :request="request" />
    <mongo-brief-request-column v-else-if="protocol === 'Mongo'" :request="request" />
    <mysql-brief-request-column v-else-if="protocol === 'Mysql'" :request="request" />
    <template v-else>
      {{ request }}
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import DubboBriefRequestColumn from "@/components/record/dubbo/DubboBriefRequestColumn.vue";
import HttpBriefRequestColumn from "@/components/record/http/HttpBriefRequestColumn.vue";
import MongoBriefRequestColumn from "@/components/record/mongo/MongoBriefRequestColumn.vue";
import MysqlBriefRequestColumn from "@/components/record/mysql/MysqlBriefRequestColumn.vue";
import RedisBriefRequestColumn from "@/components/record/redis/RedisBriefRequestColumn.vue";


export default defineComponent({
  name: "RequestColumnItem",
  props: {
    protocol: String,
    request: Object,
  },
  data() {
    return {
    };
  },
  computed: {
  },
  methods: {
  },
  components: {
    DubboBriefRequestColumn,
    MysqlBriefRequestColumn,
    HttpBriefRequestColumn,
    MongoBriefRequestColumn,
    RedisBriefRequestColumn,
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
.request-column-tag {
  margin-right:1rem;

  color: #000000;
  font-size: 0.75rem;
  font-weight: 500;
  padding: 0.25rem 0.4rem;
  border: solid 1px black;
  border-radius: 4px;
}
</style>
