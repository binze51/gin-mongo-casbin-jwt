# gosvcpublic

mongo+pulsar+casbin+jwt+gin 基础 rest api 业务框架

## 特性

- 基于 gin 两层模块隔离封装，提供了丰富的中间件支持（用户认证、跨域、访问日志、权限验证）

- 基于 Casbin 的 RBAC 访问控制模型

- JWT 认证及黑名单

- 支持 mongo 文档结构强模式验证

- zap 日志器，支持 mongo 慢日志记录

- 支持 pulsar 消息队列

- 模块错误码分离

- 国密 加解密

- 适应 K8S 环境

- 单元测试
