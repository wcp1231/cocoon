<!--<template>-->
<!--  <div>-->
<!--    <template v-if="protocol === 'HTTP'">-->
<!--      <span style="margin-right: 1rem;">{{ response.meta['STATUS'] }}</span>-->
<!--    </template>-->
<!--    <template v-else>-->
<!--      {{ response }}-->
<!--    </template>-->
<!--  </div>-->
<!--</template>-->

<script lang="jsx">
import { defineComponent, computed } from "vue";
import Tag from "primevue/tag";

export default defineComponent({
  name: "ResponseColumnItem",
  props: {
    protocol: String,
    request: Object,
    response: Object
  },
  setup(props) {
    let status = computed(() => "UNKNOWN")
    const getValueFromMeta = (data, key) => {
      if (!data) {
        return null;
      }
      if (!data.meta) {
        return null;
      }
      return data.meta[key];
    };
    const getValueFromPayload = (data, key) => {
      if (!data) {
        return null;
      }
      if (!data.payload) {
        return null;
      }
      return data.payload[key];
    };
    const checkResponse = (resp, key) => {
      if (!resp) {
        return false;
      }
      if (!resp[key]) {
        return false;
      }
      return true;
    }

    // TODO 定义统一的异常格式，现在 HTTP、Redis、Mysql 等的错误字段都不一样

    if (props.protocol === 'HTTP') {
      status = computed(() => {
        return getValueFromPayload(props.response, "HTTP_STATUS") || "PENDING";
      });
    } else if (props.protocol === 'Redis') {
      status = computed(() => {
        let ok = checkResponse(props.response, "payload");
        return ok ? "OK" : "PENDING";
      });
    } else if (props.protocol === 'Dubbo') {
      status = computed(() => {
        let ok = checkResponse(props.response, "body");
        return ok ? "OK" : "PENDING";
      });
    } else if (props.protocol === 'Mongo') {
      status = computed(() => {
        let ok = getValueFromPayload(props.response, "MONGO_MESSAGE");
        return ok ? "OK": "PENDING";
      });
    } else if (props.protocol === 'Mysql') {
      status = computed(() => {
        let ok = checkResponse(props.response, "payload");
        let isStmtClose = props.request.payload['MYSQL_OP_TYPE'] === 'COM_STMT_CLOSE'
        if (isStmtClose) {
          return "OK"
        }
        return ok ? "OK": "PENDING";
      });
    }
    let mock = computed(() => getValueFromMeta(props.response, "MOCK") )
    if (mock.value) {
      return () => (
          <>
            <span class={`response-status-tag response-status-tag-${status.value}`}>{ status.value }</span>
            <Tag severity="info">MOCK</Tag>
          </>
      )
    }
    return () => (
        <>
          <span class={`response-status-tag response-status-tag-${status.value}`} >{ status.value }</span>
        </>
    );
  },
  components: {
    Tag,
  }
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.response-status-tag {
  margin-right:1rem;

  color: #000000;
  font-size: 0.75rem;
  font-weight: 500;
  padding: 0.25rem 0.4rem;
  border: solid 1px black;
  border-radius: 4px;
}
.response-status-tag.response-status-tag-OK {
  color: #28a745;
  border: solid 1px #28a745;
}
.response-status-tag.response-status-tag-200 {
  color: #28a745;
  border: solid 1px #28a745;
}
.response-status-tag.response-status-tag-300 {
  color: #6c757d;
  border: solid 1px #6c757d;
}
.response-status-tag.response-status-tag-400 {
  color: #ffc107;
  border: solid 1px #ffc107;
}
.response-status-tag.response-status-tag-500 {
  color: #dc3545;
  border: solid 1px #dc3545;
}
</style>
