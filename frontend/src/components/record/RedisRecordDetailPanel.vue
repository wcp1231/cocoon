<template>
  <Accordion :multiple="true" :activeIndex="[0]">
    <AccordionTab header="Request">
      <div class="header-panel">
        <div>
          <span class="header-key">{{ request.payload['REDIS_CMD'] }}</span> <span class="header-val">{{ request.payload["REDIS_KEY"] }}</span>
        </div>
      </div>
      <template v-if="request.payload['REDIS_REQ_CMD']">
        <Divider type="dashed" />
        <div>
          {{ request.payload['REDIS_REQ_CMD'] }}
        </div>
      </template>
    </AccordionTab>
    <AccordionTab header="Response">
      <span class="overflow-wrap-item">{{ response.payload['REDIS_RESP_OBJ'] }}</span>
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
