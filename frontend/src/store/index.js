import { createStore } from "vuex";
import API from "../remote/api";

function getProtocol(request) {
  let protocol = "UNKNOWN";
  if (!request) {
    return protocol;
  }
  protocol = request.meta["PROTOCOL"] || protocol;
  return protocol;
}

export default createStore({
  state: {
    mocks: {},
    recordMap: new Map(),
    records: [],
    socket: {
      // 连接状态
      isConnected: false,
      // 消息内容
      message: "",
      // 重新连接错误
      reconnectError: true,
      // 心跳消息发送时间
      heartBeatInterval: 50000,
      // 心跳定时器
      heartBeatTimer: 0
    }
  },
  mutations: {// 连接打开
    SOCKET_ONOPEN(state, event) {
      //main.config.globalProperties.$socket = event.currentTarget;
      state.socket.isConnected = true;
      // 连接成功时启动定时发送心跳消息，避免被服务器断开连接
      // state.socket.heartBeatTimer = setInterval(() => {
      //   const message = "心跳消息";
      //   state.socket.isConnected &&
      //   main.config.globalProperties.$socket.sendObj({
      //     code: 200,
      //     msg: message
      //   });
      // }, state.socket.heartBeatInterval);
      console.log("已经连接");
    },
    // 连接关闭
    SOCKET_ONCLOSE(state, event) {
      state.socket.isConnected = false;
      // 连接关闭时停掉心跳消息
      clearInterval(state.socket.heartBeatTimer);
      state.socket.heartBeatTimer = 0;
      console.log("连接已断开: " + new Date());
      console.log(event);
    },
    // 发生错误
    SOCKET_ONERROR(state, event) {
      console.error(state, event);
    },
    // 收到服务端发送的消息
    SOCKET_ONMESSAGE(state, message) {
      console.log("On message.")
      console.log(message)
      state.socket.message = message;
      let record = JSON.parse(message.data)
      if (record.isRequest) {
        this.commit("on_request_records", record);
      } else {
        this.commit("on_response_records", record);
      }
    },
    // 自动重连
    SOCKET_RECONNECT(state, count) {
      console.info("消息系统重连中...", state, count);
    },
    // 重连错误
    SOCKET_RECONNECT_ERROR(state) {
      state.socket.reconnectError = true;
    },
    on_request_records(state, request) {
      let record = {
        id: request.id,
        protocol: getProtocol(request),
        request: request,
        response: {},
        timespan: -1,
      };
      state.records.push(record);
      state.recordMap.set(request.id, record);
    },
    on_response_records(state, response) {
      let record = state.recordMap.get(response.id);
      record.response = response;
      record.timespan = response.captureTime - record.request.captureTime;
    },
    update_mocks(state, mocks) {
      state.mocks = mocks;
    },
  },
  actions: {
    refresh_mocks(context) {
      API.fetchMocks().then((resp) => {
        context.commit("update_mocks", resp.data);
      });
    },
  },
});
