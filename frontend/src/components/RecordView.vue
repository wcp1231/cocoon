<template>
  <div>
    <div class="content-section introduction">
      <div class="feature-intro">
        <h1 class="title">Records View <small>网络请求记录</small></h1>
      </div>
    </div>
    <card class="record-card" style="height: calc(100vh - 130px)">
      <template #content>
<!--        <div style="width: 100%;"> &lt;!&ndash; FIXME  virtual 不会自动加载新数据 &ndash;&gt;-->
          <DataTable ref="datatable"
                     v-model:value="records"
                     v-model:expandedRows="expandedRows"
                     :scrollable="true"
                     scrollHeight="800px"
                     class="record-table">
            <Column :expander="true" headerStyle="flex: 0 0 45px;" bodyStyle="flex: 0 0 45px;" headerClass="expander-column" bodyClass="expander-column"/>
            <Column field="id" header="ID" headerStyle="flex: 0 0 60px;" bodyStyle="flex: 0 0 60px;"></Column>
            <Column field="protocol" header="Proto" filterField="protocol" headerStyle="flex: 0 0 80px;" bodyStyle="flex: 0 0 80px;">
              <template #body="slotProps">
                <tag :class="`protocol-tag protocol-${slotProps.data.protocol}-tag`">{{ slotProps.data.protocol }}</tag>
              </template>
            </Column>
            <Column field="request" header="Request" headerStyle="flex: 5;" bodyStyle="overflow-wrap: anywhere; flex: 5;">
              <template #body="slotProps">
                <request-column-item :protocol="slotProps.data.protocol" :request="slotProps.data.request" />
              </template>
            </Column>
            <Column field="response" header="Response" headerStyle="flex: 1;" bodyStyle="flex: 1;">
              <template #body="slotProps">
                <response-column-item :protocol="slotProps.data.protocol" :request="slotProps.data.request" :response="slotProps.data.response" />
              </template>
            </Column>
            <Column field="timespan" header="Time" headerStyle="flex: 0 0 100px;" bodyStyle="flex: 0 0 100px;">
              <template #body="slotProps">
                <timespan-column-item :timespan="slotProps.data.timespan" />
              </template>
            </Column>
            <template #expansion="slotProps">
              <div class="row-expansion-panel">
                <http-record-detail-panel v-if="slotProps.data.protocol === 'HTTP'" :record="slotProps.data" />
                <redis-record-detail-panel v-if="slotProps.data.protocol === 'Redis'" :record="slotProps.data" />
                <dubbo-record-detail-panel v-if="slotProps.data.protocol === 'Dubbo'" :record="slotProps.data" />
                <mongo-record-detail-panel v-if="slotProps.data.protocol === 'Mongo'" :record="slotProps.data" />
                <mysql-record-detail-panel v-if="slotProps.data.protocol === 'Mysql'" :record="slotProps.data" />
              </div>
            </template>
          </DataTable>
<!--        </div>-->
      </template>
    </card>
  </div>
</template>

<script lang="js">
import { defineComponent } from "vue";
import Card from "primevue/card";
import Column from "primevue/column";
import DataTable from "primevue/datatable";
import Tag from "primevue/tag";
//@ts-ignore
import store from "@/store/index";
import RequestColumnItem from "@/components/record/RequestColumnItem.vue";
import ResponseColumnItem from "@/components/record/ResponseColumnItem.vue";
import TimespanColumnItem from "@/components/record/TimespanColumnItem.vue";
import HttpRecordDetailPanel from "@/components/record/HttpRecordDetailPanel.vue";
import RedisRecordDetailPanel from "@/components/record/RedisRecordDetailPanel.vue";
import DubboRecordDetailPanel from "@/components/record/DubboRecordDetailPanel.vue";
import MongoRecordDetailPanel from "@/components/record/MongoRecordDetailPanel.vue";
import MysqlRecordDetailPanel from "@/components/record/MysqlRecordDetailPanel.vue";

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
    recordLen() {
      return store.state.records.length;
    }
  },
  watch: {
    recordLen() {
      // FIXME
      // let datatable = this.$refs.datatable;
      // if (datatable) {
      //   datatable.$refs.table.__vueParentComponent.proxy.init();
      // }
      console.log(`watch. ${this.records.length}`)
    }
  },
  methods: {
  },
  components: {
    Card,
    Column,
    DataTable,
    Tag,
    RequestColumnItem,
    ResponseColumnItem,
    TimespanColumnItem,
    HttpRecordDetailPanel,
    RedisRecordDetailPanel,
    DubboRecordDetailPanel,
    MongoRecordDetailPanel,
    MysqlRecordDetailPanel,
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.content-section.introduction .feature-intro h1 {
  margin-top: 0;
}
.protocol-tag {
  background: rgba(0, 0, 0, 0.2);
  color: #ffffff;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0.25rem 0.4rem;
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
.protocol-tag.protocol-Mongo-tag {
  background-color: rgb(17, 97, 73);
}
.protocol-tag.protocol-Mysql-tag {
  background-color: #3E6E93;
}
</style>

<style>
.title small {
  font-size: .5em;
  color: #6c757d
}
.record-card.p-card .p-card-body {
  padding-top: 0;
}
.record-table .p-datatable-table .expander-column {
  padding: 5px;
}
.row-expansion-panel .p-accordion.p-component a.p-accordion-header-link {
  padding: 8px 10px;
}
.record-table.p-datatable .p-datatable-thead > tr > th {
  padding: 0.5rem 0;
}
.record-table.p-datatable .p-datatable-tbody > tr > td {
  padding: 0.5rem 0;
}
</style>
