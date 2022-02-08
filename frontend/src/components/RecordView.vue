<template>
  <div>
    <div class="content-section introduction">
      <div class="feature-intro">
        <h1>Record View</h1>
        <p>可以查看网络记录</p>
      </div>
    </div>
    <card class="record-card">
      <template #content>
        <div style="width: 100%">
          <DataTable :value="records" v-model:expandedRows="expandedRows">
            <Column :expander="true" headerStyle="width: 3rem" />
            <Column field="id" header="ID" headerStyle="width: 2rem"></Column>
            <Column field="request" header="Request" bodyStyle="overflow-wrap: anywhere; max-width: 50vh">
              <template #body="slotProps">
                <request-column-item :protocol="slotProps.data.protocol" :request="slotProps.data.request" />
              </template>
            </Column>
            <Column field="response" header="Response" >
              <template #body="slotProps">
                <response-column-item :protocol="slotProps.data.protocol" :response="slotProps.data.response" />
              </template>
            </Column>
            <Column field="timespan" header="Time" >
              <template #body="slotProps">
                <timespan-column-item :timespan="slotProps.data.timespan" />
              </template>
            </Column>
            <template #expansion="slotProps">
              <http-record-detail-panel v-if="slotProps.data.protocol === 'HTTP'" :record="slotProps.data" />
              <redis-record-detail-panel v-if="slotProps.data.protocol === 'Redis'" :record="slotProps.data" />
              <dubbo-record-detail-panel v-if="slotProps.data.protocol === 'Dubbo'" :record="slotProps.data" />
            </template>
          </DataTable>
        </div>
      </template>
    </card>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Card from "primevue/card";
import Column from "primevue/column";
import DataTable from "primevue/datatable";
//@ts-ignore
import store from "@/store/index";
import RequestColumnItem from "@/components/record/RequestColumnItem.vue";
import ResponseColumnItem from "@/components/record/ResponseColumnItem.vue";
import TimespanColumnItem from "@/components/record/TimespanColumnItem.vue";
import HttpRecordDetailPanel from "@/components/record/HttpRecordDetailPanel.vue";
import RedisRecordDetailPanel from "@/components/record/RedisRecordDetailPanel.vue";
import DubboRecordDetailPanel from "@/components/record/DubboRecordDetailPanel.vue";

export default defineComponent({
  name: "RecordView",
  props: {},
  data() {
    return {
      expandedRows: []
    };
  },
  computed: {
    records() {
      return store.state.records;
    },
  },
  methods: {},
  components: {
    Card,
    Column,
    DataTable,
    RequestColumnItem,
    ResponseColumnItem,
    TimespanColumnItem,
    HttpRecordDetailPanel,
    RedisRecordDetailPanel,
    DubboRecordDetailPanel,
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.content-section.introduction .feature-intro h1 {
  margin-top: 0;
}
</style>

<style>
.record-card.p-card .p-card-body {
  padding-top: 0;
}
</style>