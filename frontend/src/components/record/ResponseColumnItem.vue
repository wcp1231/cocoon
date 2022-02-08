<!--<template>-->
<!--  <div>-->
<!--    <template v-if="protocol === 'HTTP'">-->
<!--      <span style="margin-right: 1rem;">{{ response.meta['STATUS'] }}</span>-->
<!--    </template>-->
<!--    <template v-else>-->
<!--      {{ response }}-->
<!--    </template>-->
<!--  </div>-->
<!--</template>-->

<script lang="jsx">
import { defineComponent, computed } from "vue";

export default defineComponent({
  name: "ResponseColumnItem",
  props: {
    protocol: String,
    response: Object
  },
  setup(props) {
    let status = computed(() => "UNKNOWN")
    if (props.protocol === 'HTTP') {
      status = computed(() => {
        let status = "PENDING";
        if (!props.response) {
          return status;
        }
        if (!props.response.meta) {
          return status;
        }
        return props.response.meta["STATUS"] || status;
      });
    } else if (props.protocol === 'Redis') {
      status = computed(() => {
        let status = "PENDING";
        if (!props.response) {
          return status;
        }
        if (!props.response.meta) {
          return status;
        }
        return "OK";
      });
    } else if (props.protocol === 'Dubbo') {
      status = computed(() => {
        let status = "PENDING";
        if (!props.response) {
          return status;
        }
        if (!props.response.body) {
          return status;
        }
        return "OK";
      });
    }
    return () => (
        <div><span>{ status.value }</span></div>
    );
  },
});
</script>

<!--<script lang="ts">-->
<!--import { defineComponent } from "vue";-->

<!--export default defineComponent({-->
<!--  name: "ResponseColumnItem",-->
<!--  props: {-->
<!--    response: Object,-->
<!--  },-->
<!--  data() {-->
<!--    return {-->
<!--    };-->
<!--  },-->
<!--  computed: {-->
<!--    protocol() {-->
<!--      let protocol = 'UNKNOWN';-->
<!--      if (!this.response) {-->
<!--        return protocol;-->
<!--      }-->
<!--      protocol = this.response.meta['PROTOCOL'] || protocol;-->
<!--      return protocol;-->
<!--    }-->
<!--  },-->
<!--  methods: {-->
<!--  },-->
<!--  components: {-->
<!--  },-->
<!--});-->
<!--</script>-->

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
