<template>
  <span>{{text}}</span>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "TimespanColumnItem",
  props: {
    timespan: Number,
  },
  data() {
    return {};
  },
  computed: {
    text() {
      if (!this.timespan || this.timespan < 0) {
        return "-"
      }
      let time = this.timespan / 1e6;
      if (time < 1000) {
        return `${time.toFixed(1)}ms`;
      }
      time = time / 1e3;
      if (time < 120) {
        return `${time.toFixed(1)}s`;
      }
      let h = time / 3600;
      let m = (time % 3600) / 60;
      let s = time % 60;
      return `${h}:${m}:${s}`;
    }
  },
  methods: {
  },
  components: {
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
