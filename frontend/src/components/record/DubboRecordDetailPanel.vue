<template>
  <Accordion :multiple="true" :activeIndex="[0]">
    <AccordionTab header="Request">
      <div v-if="request.meta['HEARTBEAT']">Heartbeat</div>
      <div v-else class="header-panel">
        <div v-for="(v, k) in request.header" v-bind:key="k">
          <span class="header-key">{{ k }}: </span>
          <span class="header-val overflow-wrap-item">{{ v }}</span>
        </div>
      </div>
<!--      <div>-->
<!--        {{ request.body }}-->
<!--      </div>-->
    </AccordionTab>
    <AccordionTab header="Response">
      <div v-if="response.meta['HEARTBEAT']">Heartbeat</div>
      <div v-else class="header-panel">
        <div v-for="(v, k) in response.header" v-bind:key="k">
          <span class="header-key">{{ k }}: </span>
          <span class="header-val overflow-wrap-item">{{ v }}</span>
        </div>
      </div>
<!--      <span class="overflow-wrap-item">{{ response.body }}</span>-->
    </AccordionTab>
  </Accordion>
</template>

<script lang="js">
import { defineComponent } from "vue";
import Accordion from 'primevue/accordion';
import AccordionTab from 'primevue/accordiontab';

export default defineComponent({
  name: "DubboRecordDetailPanel",
  props: {
    record: Object
  },
  data() {
    return {};
  },
  computed: {
    request() {
      return this.record.request;
    },
    response() {
      return this.record.response;
    }
  },
  methods: {
  },
  components: {
    Accordion,
    AccordionTab,
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.header-panel .header-key {
  font-weight: bold;
}
.header-panel .header-val {
  margin-left: 1.25rem;
}
.overflow-wrap-item {
  overflow-wrap: anywhere;
}
</style>
