<template>
  <Accordion :multiple="true" :activeIndex="[0]">
    <AccordionTab header="Request">
      <div class="header-panel">
        <div>
          <span class="header-key">{{ request.meta["CMD"] }}</span> <span class="header-val">{{ request.meta["KEY"] }}</span>
        </div>
      </div>
      <template v-if="request.body">
        <Divider type="dashed" />
        <div>
          {{ request.body }}
        </div>
      </template>
    </AccordionTab>
    <AccordionTab header="Response">
      <span class="overflow-wrap-item">{{ response.body }}</span>
    </AccordionTab>
  </Accordion>
</template>

<script lang="js">
import { defineComponent } from "vue";
import Accordion from 'primevue/accordion';
import AccordionTab from 'primevue/accordiontab';
import Divider from 'primevue/divider';

export default defineComponent({
  name: "RedisRecordDetailPanel",
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
    Divider,
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
