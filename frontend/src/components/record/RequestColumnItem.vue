<template>
  <div class="request-column">
    <tag :class="`protocol-tag protocol-${protocol}-tag`">{{ protocol }}</tag>
    <template v-if="protocol === 'HTTP'">
      <span class="request-column">
        <span class="cmd-tag">{{ request.meta['METHOD'] }}</span>
      </span>
      <span class="request-column">{{ request.meta['HOST'] + request.meta['URL'] }}</span>
    </template>
    <template v-else-if="protocol === 'Redis'">
      <span class="request-column">
        <span class="cmd-tag">{{ request.meta['CMD'] }}</span>
      </span>
      <span class="request-column">{{ request.meta['KEY'] }}</span>
    </template>
    <template v-else-if="protocol === 'Dubbo'">
      <span v-if="request.meta['HEARTBEAT']" class="request-column">
        <span class="cmd-tag">HeartBeat</span>
      </span>
      <span v-else class="request-column">{{ request.body }}</span>
    </template>
    <template v-else>
      {{ request }}
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Tag from 'primevue/tag';

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
    Tag,
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.request-column {
  margin-left:1rem;
}
.request-column .protocol-tag {
  background: rgba(0, 0, 0, 0.2);
  color: #ffffff;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0.25rem 0.4rem;
  border-radius: 4px;
}
.request-column .cmd-tag {
  color: #000000;
  font-size: 0.75rem;
  font-weight: 500;
  padding: 0.25rem 0.4rem;
  border: solid 1px black;
  border-radius: 4px;
}
.protocol-tag.protocol-HTTP-tag {
  background-color: rgb(13, 110, 253);
}
.protocol-tag.protocol-Redis-tag {
  background-color: #a51f17;
}
.protocol-tag.protocol-Dubbo-tag {
  background-color: #2ba3de;
}
.request-column .heart-beat-tag {
  background-color: rgb(108, 117, 125);
}
</style>
