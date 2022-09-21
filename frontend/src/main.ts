import { createApp } from "vue";
import App from "./App.vue";
import "primeicons/primeicons.css";
import "primevue/resources/themes/bootstrap4-light-blue/theme.css";
import "primevue/resources/primevue.min.css";
import PrimeVue from "primevue/config";
import ConfirmationService from "primevue/confirmationservice";
import toastservice from "primevue/toastservice";
import VueNativeSock from "vue-native-websocket-vue3";

// @ts-ignore
import store from "@/store/index";

const app = createApp(App);

const wsUrl = `ws://${location.host}/api/ws`
//const wsUrl = `ws://127.0.0.1:7070/api/ws`

app
  .use(PrimeVue)
  .use(ConfirmationService)
  .use(toastservice)
  .use(store)
  .use(VueNativeSock, wsUrl, {
    store: store
  });

app.mount("#app");
