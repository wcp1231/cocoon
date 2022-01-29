import { createStore } from "vuex";
import API from "../remote/api";

export default createStore({
  state: {
    mocks: {},
    records: [],
  },
  mutations: {
    update_mocks(state, mocks) {
      state.mocks = mocks;
    },
  },
  actions: {
    refresh_mocks(context) {
      API.fetchMocks().then((resp) => {
        console.log(resp.data);
        context.commit("update_mocks", resp.data);
      });
    },
  },
});
