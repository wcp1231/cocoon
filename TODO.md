# Todo List

- 支持 Mock
  - [x] 支持 Mock 热更新 
  - [ ] 支持延时模拟
- 支持流量比较
  - [ ] 区分入流量和出流量
  - [ ] 基本比较逻辑
  - [ ] 通过 Mock 保证外部依赖返回数据一致
- 协议支持
  - [ ] dubbo
  - [ ] mysql
  - [ ] kafka
- Web 管理功能
  - [x] 简单 Web 页面
  - [x] 简单 API 接口
- 测试
  - [ ] 单元测试
- 其他优化
  - [x] 支持根据端口直接标识协议的配置机制
  - [ ] 重新考虑 agent 和 server 的职责
  - [ ] HTTP Gzip 类型自动解压
  - [ ] Redis 支持 empty 和 nil 的区分