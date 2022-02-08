<template>
  <div class="request-column">
    <span>{{ protocol }}</span>
    <template v-if="protocol === 'HTTP'">
      <span class="request-column">{{ request.meta['METHOD'] }}</span>
      <span class="request-column">{{ request.meta['HOST'] + request.meta['URL'] }}</span>
    </template>
    <template v-else-if="protocol === 'Redis'">
      <span class="request-column">{{ request.meta['CMD'] }}</span>
      <span class="request-column">{{ request.meta['KEY'] }}</span>
    </template>
    <template v-else-if="protocol === 'Dubbo'">
      <span v-if="request.body === 'N'" class="request-column">Heart Beat</span>
      <span v-else class="request-column">{{ request.body }}</span>
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
.request-column {
  margin-left:1rem;
}
</style>
