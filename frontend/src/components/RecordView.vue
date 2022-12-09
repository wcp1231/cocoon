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
                     selectionMode="single"
                     v-model:selection="selectedRecord"
                     @rowSelect="onRowSelect"
                     :scrollable="true"
                     v-model:filters="filters"
                     filterDisplay="menu"
                     scrollHeight="800px"
                     class="record-table">
            <Column field="id" header="ID" :sortable="true" headerStyle="flex: 0 0 60px; padding-left: 0.5rem" bodyStyle="flex: 0 0 60px; padding-left: 0.5rem"></Column>
            <Column field="protocol" header="Proto" filterField="protocol"
                    :showFilterMatchModes="false" headerStyle="flex: 0 0 80px;" bodyStyle="flex: 0 0 80px;"
                    filterMenuClass="protocol-filter">
              <template #body="slotProps">
                <tag :class="`protocol-tag protocol-${slotProps.data.protocol}-tag`">{{ slotProps.data.protocol }}</tag>
              </template>
              <template #filter="{filterModel}">
                <Listbox v-model="filterModel.value" :multiple="true"
                         :options="protocols" optionLabel="val" optionValue="val"
                         class="protocol-filter-list">
                  <template #option="slotProps">
                    <tag :class="`protocol-tag protocol-${slotProps.option.val}-tag`">{{ slotProps.option.val }}</tag>
                  </template>
                </Listbox>
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
          </DataTable>
<!--        </div>-->
      </template>
    </card>
    <Sidebar
        class="detail-sidebar"
        v-model:visible="sidebarVisible"
        :dismissable="false"
        :modal="false"
        position="right">
      <template #header>
        <tag :class="`protocol-tag protocol-${selectedRecord.protocol}-tag`" style="margin-right: 2rem">{{ selectedRecord.protocol }}</tag>

        <response-column-item :protocol="selectedRecord.protocol" :request="selectedRecord.request" :response="selectedRecord.response" />
        <timespan-column-item :timespan="selectedRecord.timespan" style="margin-left: 2rem" />
      </template>
      <http-record-detail-panel v-if="selectedRecord.protocol === 'HTTP'" :record="selectedRecord" />
      <redis-record-detail-panel v-if="selectedRecord.protocol === 'Redis'" :record="selectedRecord" />
      <dubbo-record-detail-panel v-if="selectedRecord.protocol === 'Dubbo'" :record="selectedRecord" />
      <mongo-record-detail-panel v-if="selectedRecord.protocol === 'Mongo'" :record="selectedRecord" />
      <mysql-record-detail-panel v-if="selectedRecord.protocol === 'Mysql'" :record="selectedRecord" />
      {{ selectedRecord }}
    </Sidebar>
  </div>
</template>

<script lang="js">
import { defineComponent } from "vue";
import Card from "primevue/card";
import Column from "primevue/column";
import DataTable from "primevue/datatable";
import {FilterMatchMode} from 'primevue/api';
import Listbox from 'primevue/listbox';
import Sidebar from 'primevue/sidebar';
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
      filters: null,
      expandedRows: [],
      sidebarVisible: false,
      selectedRecord: {},
    };
  },
  computed: {
    records() {
      return store.state.records;
    },
    recordLen() {
      return store.state.records.length;
    },
    protocols() {
      return [...new Set(store.state.records.map(r => r.protocol))]
          .map((p, i) => { return { id: i, val: p }; });
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
  created() {
    this.initFilters();
  },
  methods: {
    onRowSelect() {
      this.sidebarVisible = true;
    },
    clearFilters() {
      this.initFilters();
    },
    initFilters() {
      this.filters = {
        protocol: { value: null, matchMode: FilterMatchMode.IN }
      };
    }
  },
  components: {
    Card,
    Column,
    DataTable,
    Listbox,
    Sidebar,
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
  padding: 0;
}
.record-card.p-card .p-card-content {
  padding: 0;
}
.record-table .p-datatable-table .expander-column {
  padding: 5px;
}
.record-table .p-datatable-table .p-datatable-tbody > tr.p-highlight {
  background: #EFF6FF;
  color: #1D4ED8;
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
.p-sidebar-right.detail-sidebar {
  width: 45%;
}
.p-sidebar-right.detail-sidebar .p-sidebar-header {
  justify-content: space-between;
}

.protocol-filter.p-column-filter-overlay-menu .p-column-filter-buttonbar {
  padding-top: 0;
}
.protocol-filter-list.p-listbox .p-listbox-list {
  padding: 0;
}
</style>
