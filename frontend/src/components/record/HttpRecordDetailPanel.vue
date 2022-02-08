<template>
  <Accordion :multiple="true" :activeIndex="[0, 1]">
    <AccordionTab header="Request">
      <div>
        <div class="header-panel">
          <div>
            <span class="header-key">Method: </span><span class="header-val">{{ request.meta["METHOD"] }}</span>
          </div>
          <div>
            <span class="header-key">URL: </span><span class="header-val overflow-wrap-item">{{ request.meta['HOST'] + request.meta['URL'] }}</span>
          </div>
        </div>
      </div>
      <Divider type="dashed">
        <b>Request Header</b>
      </Divider>
      <div class="header-panel">
        <div v-for="(v, k) in request.header" v-bind:key="k">
          <span class="header-key">{{ k }}: </span>
          <span class="header-val overflow-wrap-item">{{ v }}</span>
        </div>
      </div>
      <template v-if="request.body">
        <Divider type="dashed">
          <b>Request Body</b>
        </Divider>
        <div>
          {{ request.body }}
        </div>
      </template>
    </AccordionTab>
    <AccordionTab header="Response Header">
      <div class="header-panel">
        <div v-for="(v, k) in response.header" v-bind:key="k">
          <span class="header-key">{{ k }}: </span>
          <span class="header-val overflow-wrap-item">{{ v }}</span>
        </div>
      </div>
    </AccordionTab>
    <AccordionTab header="Response Body">
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
  name: "HttpRecordDetailPanel",
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
