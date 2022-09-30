# Cocoon 茧

Cocoon 源自《柯南：贝克街的亡灵》中的现实模拟设备「茧」。

CoCoon 的主要功能是接管目标服务的所有进出流量，并对这些流量进行解析、记录甚至 Mock。
使得目标服务完全处于虚拟环境中，从而可以不依赖外部环境地进行开发、测试等工作。

## Mock 结果数据

Cocoon 支持根据请求数据的部分字段进行比较，对命中对请求 Mock 返回结果。

字段比较对方式当前支持两种：
- `equals` 字符串匹配
- `regex` 正则比较

同时，当前支持 HTTP 和 Redis 的结果 Mock。

### HTTP Mock

HTTP 请求支持 `Host`、`Method` 和 `URL` 三个字段的判断。比如

```json
{
  "request": {
    "host": { "equals": "http://httpbin.org" },
    "method": { "equals": "post" },
    "url": { "regex": ".*/get/.*" }
  },
  "response": {
    "status": "200",
    "header": {
      "Cache-Control": "no-cache, no-store, max-age=0, must-revalidate",
      "Content-Type": "application/json;charset=UTF-8"
    },
    "body": "[{\"some\":\"mock data\"}]"
  }
}
```

### Redis Mock

Redis 请求支持 `cmd` 和 `key` 两个字段的判断。比如

```json
[
  {
    "request": {
      "cmd": { "equals": "get" },
      "key": { "equals": "string_value" }
    },
    "response": {
      "type": "string",
      "value": "some data"
    }
  },
  {
    "request": {
      "cmd": { "equals": "get" },
      "key": { "equals": "integer_value" }
    },
    "response": {
      "type": "integer",
      "value": "100"
    }
  },
  {
    "request": {
      "cmd": { "equals": "get" },
      "key": { "equals": "response_error" }
    },
    "response": {
      "type": "error",
      "value": "some error"
    }
  },
  {
    "request": {
      "cmd": { "equals": "get" },
      "key": { "equals": "null_value_data" }
    },
    "response": {
      "type": "null",
      "value": ""
    }
  },
  {
    "request": {
      "cmd": { "equals": "zrange" },
      "key": { "equals": "array_response" }
    },
    "response": {
      "type": "array",
      "array": [
        { "type": "string", "value": "mem_13" },
        { "type": "string", "value": "100" }
      ]
    }
  }
]
```

### Mongo Mock

待支持

### Mysql Mock

待支持