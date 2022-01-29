<template>
  <Panel :header="'Mock #' + mock.id" :toggleable="true" class="mock-panel">
    <template #icons>
      <Button icon="pi pi-trash" class="p-button-danger p-button-text mock-delete-btn" @click="tryDelete($event, mock.id)" />
    </template>
    <h3 class="mock-section-header">Request</h3>
    <div v-if="mock.request.cmd">
      <field-matcher-item field="CMD" :matcher="mock.request.cmd" />
    </div>
    <div v-if="mock.request.key">
      <field-matcher-item field="KEY" :matcher="mock.request.key" />
    </div>
    <h3 class="mock-section-header">Response</h3>
    <div v-if="mock.response.type">
      <p>Type <Tag :value="mock.response.type"></Tag></p>
    </div>
    <div>
      <p>Data</p>
      <Textarea :value="resp" disabled rows="10" cols="50" />
    </div>
    <ConfirmPopup />
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
    resp() {
      let data = this.mock.response[this.mock.response.type]
      if (typeof(data) === "string") {
        return data;
      }
      return JSON.stringify(data, null, 2);
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
