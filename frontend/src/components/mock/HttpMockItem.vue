<template>
  <Panel :header="'Mock #' + mock.id" :toggleable="true" class="mock-panel">
    <template #icons>
      <Button icon="pi pi-trash" class="p-button-danger p-button-text mock-delete-btn" @click="tryDelete($event, mock.id)"/>
    </template>
    <h3 class="mock-section-header">Request</h3>
    <div v-if="mock.request.method">
      <field-matcher-item field="METHOD" :matcher="mock.request.method" />
    </div>
    <div v-if="mock.request.host">
      <field-matcher-item field="HOST" :matcher="mock.request.host" />
    </div>
    <div v-if="mock.request.url">
      <field-matcher-item field="URL" :matcher="mock.request.url" />
    </div>
    <div v-if="mock.request.header">
      <field-matcher-item field="header-" :matcher="mock.request.header" />
    </div>
    <h3 class="mock-section-header">Response</h3>
    <div v-if="mock.response.status">
      <p>Status <Tag :value="mock.response.status"></Tag></p>
    </div>
    <div v-if="mock.response.header">
      <p>Header</p>
      <Textarea :value="respHeaders" disabled rows="5" cols="50" />
    </div>
    <div v-if="mock.response.body">
      <p>Body</p>
      <Textarea :value="respBody" disabled rows="10" cols="50" />
    </div>
    <ConfirmPopup></ConfirmPopup>
    <Toast />
  </Panel>
</template>

<script>
import { defineComponent } from "vue";
import Panel from 'primevue/panel';
import Button from "primevue/button";
import Tag from 'primevue/tag';
import Textarea from 'primevue/textarea';
import Toast from 'primevue/toast';
import ConfirmPopup from 'primevue/confirmpopup';
import FieldMatcherItem from "@/components/mock/FieldMatcherItem.vue";
import store from "@/store/index";
import API from "@/remote/api";

export default defineComponent({
  name: "HttpMockItem",
  props: {
    mock: Object
  },
  data() {
    return {};
  },
  computed: {
    respHeaders() {
      return JSON.stringify(this.mock.response.header, null, 2);
    },
    respBody() {
      let body = JSON.parse(this.mock.response.body);
      return JSON.stringify(body, null, 2);
    }
  },
  methods: {
    tryDelete(event, id) {
      this.$confirm.require({
        target: event.currentTarget,
        message: '确定要删除这个 Mock 配置吗？',
        icon: 'pi pi-exclamation-triangle',
        accept: () => {
          this.deleteMock(id);
        }
      });
    },
    deleteMock(id) {
      API.deleteMock(id).then(() => {
        this.$toast.add({severity:'info', summary:'Success', detail:'删除 Mock 成功', life: 3000});
        store.dispatch("refresh_mocks");
      });
    }
  },
  components: {
    Button,
    Panel,
    Tag,
    Textarea,
    Toast,
    ConfirmPopup,
    FieldMatcherItem,
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.mock-panel .mock-section-header:first-child {
  margin-top: 0;
}
</style>
