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
import Tag from "primevue/tag";

export default defineComponent({
  name: "ResponseColumnItem",
  props: {
    protocol: String,
    response: Object
  },
  setup(props) {
    let status = computed(() => "UNKNOWN")
    const getValueFromMeta = (data, key) => {
      if (!data) {
        return null;
      }
      if (!data.meta) {
        return null;
      }
      return data.meta[key];
    };
    const checkResponse = (resp, key) => {
      if (!resp) {
        return false;
      }
      if (!resp[key]) {
        return false;
      }
      return true;
    }
    let mock = getValueFromMeta(props.response, "MOCK");

    if (props.protocol === 'HTTP') {
      status = computed(() => {
        return getValueFromMeta(props.response, "STATUS") || "PENDING";
      });
    } else if (props.protocol === 'Redis') {
      status = computed(() => {
        let ok = checkResponse(props.response, "meta");
        return ok ? "OK" : "PENDING";
      });
    } else if (props.protocol === 'Dubbo') {
      status = computed(() => {
        let ok = checkResponse(props.response, "body");
        return ok ? "OK" : "PENDING";
      });
    } else if (props.protocol === 'Mongo') {
      status = computed(() => {
        let ok = checkResponse(props.response, "body");
        return ok ? "OK": "PENDING";
      });
    }
    if (mock) {
      return () => (
          <div><span class="response-status-tag">{ status.value }</span><Tag severity="info">MOCK</Tag></div>
      )
    }
    return () => (
        <div><span class="response-status-tag">{ status.value }</span></div>
    );
  },
  components: {
    Tag,
  }
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
.response-status-tag {
  margin-right:1rem;

  color: #000000;
  font-size: 0.75rem;
  font-weight: 500;
  padding: 0.25rem 0.4rem;
  border: solid 1px black;
  border-radius: 4px;
}
</style>
