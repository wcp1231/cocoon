<template>
  <div>
    <Sidebar
      :visible="true"
      :modal="false"
      :show-close-icon="false"
      class="layout-sidebar"
    >
      <h3>Cocoon</h3>
      <Menu :model="items">
        <template #item="{ item }">
          <a
            class="p-menuitem-link"
            role="menuitem"
            :tabindex="item.id"
            @click="menuClick(item.id)"
          >
            <span class="p-menuitem-icon" :class="item.icon"></span>
            <span class="p-menuitem-text">{{ item.label }}</span>
          </a>
        </template>
      </Menu>
    </Sidebar>
    <tab-view class="layout-content" :activeIndex="activeIndex">
      <tab-panel>
        Home
      </tab-panel>
      <tab-panel>
        <record-view />
      </tab-panel>
      <tab-panel>
        <mock-view />
      </tab-panel>
    </tab-view>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Menu from "primevue/menu";
import Sidebar from "primevue/sidebar";
import TabView from "primevue/tabview";
import TabPanel from "primevue/tabpanel";
import MockView from "@/components/MockView.vue";
import RecordView from "@/components/RecordView.vue";

export default defineComponent({
  name: "Layout",
  props: {},
  data() {
    return {
      activeIndex: 0,
      items: [
        {
          id: 0,
          label: "Home",
          icon: "pi pi-home",
        },
        {
          id: 1,
          label: "Records",
          icon: "pi pi-inbox",
        },
        {
          id: 2,
          label: "Mocks",
          icon: "pi pi-pencil",
        },
      ],
    };
  },
  methods: {
    menuClick(idx: number) {
      this.activeIndex = idx;
    },
  },
  components: {
    Menu,
    Sidebar,
    TabView,
    TabPanel,
    MockView,
    RecordView,
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
.layout-sidebar.p-sidebar-left {
  width: 200px;
}
.layout-sidebar .p-sidebar-header {
  display: none;
}
.layout-sidebar.p-sidebar .p-sidebar-content {
  padding: 0;
}

.layout-content {
  margin-left: 200px;
}
.layout-content > .p-tabview-nav-container {
  display: none;
}
.layout-content.p-tabview > .p-tabview-panels {
  height: 100vh;
  background-color: #efefef;
  overflow: scroll;
}
</style>
