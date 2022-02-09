<template>
  <div class="request-column">
    <template v-if="protocol === 'HTTP'">
      <span class="request-column-tag">{{ request.meta['METHOD'] }}</span>
      <span class="request-column-info">{{ request.meta['HOST'] + request.meta['URL'] }}</span>
    </template>
    <template v-else-if="protocol === 'Redis'">
      <span class="request-column-tag">{{ request.meta['CMD'] }}</span>
      <span class="request-column-info">{{ request.meta['KEY'] }}</span>
    </template>
    <template v-else-if="protocol === 'Dubbo'">
      <span v-if="request.meta['HEARTBEAT']" class="request-column-tag">HeartBeat</span>
      <template v-else>
        <span class="request-column-tag">Invoke</span>
        <span class="request-column-info">{{ `${request.header['target']}#${request.header['method']}` }}</span>
      </template>
    </template>
    <template v-else-if="protocol === 'Mongo'">
      <span class="request-column-tag">{{ request.header["op_type"] }}</span>
      <span class="request-column-info">{{ `${request.header['collection']} ${request.header["query"]}` }}</span>
    </template>
    <template v-else>
      {{ request }}
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

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
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
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
