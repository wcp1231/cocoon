<template>
  <div>
    <div class="content-section introduction">
      <div class="feature-intro">
        <h1>Mocks View</h1>
        <p>可以查看和修改 Mock 设置</p>
      </div>
    </div>
    <card class="mock-card">
      <template #content>
        <Button label="添加 Mock" class="p-button-sm add-mock-btn" @click="showCreateMockDialog"/>
        <TabView :activeIndex="activeIndex" class="mock-tab-view">
          <TabPanel header="HTTP">
            <http-mock-item v-for="mock in mocks.HTTP" :mock="mock" v-bind:key="mock.id" />
          </TabPanel>
          <TabPanel header="Redis">
            <redis-mock-item v-for="mock in mocks.Redis" :mock="mock" v-bind:key="mock.id" />
          </TabPanel>
        </TabView>
      </template>
    </card>
    <Dialog :header="dialogHeader" v-model:visible="dialogDisplay" @close="closeMockDialog">
      <!-- TODO 后续再写 form 表单 -->
      <h5 style="margin-top: 0;">Json Config</h5>
      <Textarea v-model="createMockData" rows="30" cols="80" />

      <template #footer>
        <Button label="Cancel" icon="pi pi-times" class="p-button-text" @click="closeMockDialog" />
        <Button label="Submit" icon="pi pi-check" @click="submitMock"/>
      </template>
    </Dialog>
    <Toast />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Button from "primevue/button";
import Card from "primevue/card";
import Dialog from "primevue/dialog";
import TabView from "primevue/tabview";
import TabPanel from "primevue/tabpanel";
import Textarea from "primevue/textarea";
//@ts-ignore
import store from "@/store/index";
import HttpMockItem from "@/components/mock/HttpMockItem.vue";
import RedisMockItem from "@/components/mock/RedisMockItem.vue";
import API from "@/remote/api";

export default defineComponent({
  name: "MockView",
  props: {},
  data() {
    return {
      activeIndex: 0,
      dialogHeader: "",
      dialogDisplay: false,
      createMockData: "",
    };
  },
  computed: {
    mocks() {
      return store.state.mocks;
    }
  },
  methods: {
    showCreateMockDialog() {
      this.dialogHeader = "添加 Mock";
      this.createMockData = "";
      this.dialogDisplay = true;
    },
    closeMockDialog() {
      this.dialogDisplay = false;
    },
    submitMock() {
      let data = JSON.parse(this.createMockData);
      API.createMocks(data).then(() => {
        this.$toast.add({severity:'info', summary:'Success', detail:'添加 Mock 成功', life: 3000});
        store.dispatch("refresh_mocks");
        this.closeMockDialog();
      });
    }
  },
  mounted() {
    store.dispatch("refresh_mocks");
  },
  components: {
    Button,
    Card,
    Dialog,
    TabView,
    TabPanel,
    Textarea,
    HttpMockItem,
    RedisMockItem,
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.content-section.introduction .feature-intro h1 {
  margin-top: 0;
}

.add-mock-btn {
  position: absolute;
  right: 48px;
  z-index: 1;
}
</style>

<style>
.mock-card.p-card .p-card-body {
  padding-top: 0;
}
</style>