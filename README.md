# tumbleweed(风滚草)工作流引擎

简体中文 | [English](./README-EN.md)

## 简介

基于bpmn 2.0的工作流引擎, 将用户、表单等拆分出去，专注于流程流转本身

配合工作流设计器使用
> https://github.com/lzw5399/tumbleweed-designer

## 项目定位

由于实际业务中要在主系统中集成一款可定制化的工作流组件，要求是专注于流程的设计和流转。于是就搞了一个工作流的独立微服务，提供RESTful服务(后续考虑改成grpc-gateway形式). 将【用户/角色】【表单】等依赖分离出去，支持多租户

针对【用户/角色】的处理
- 接入工作流引擎的系统指定一个租户
- 然后将本身系统中需要用到工作流的用户/角色 的id和name同步到工作流引擎的数据库中
- 工作流引擎本身只关心用户/角色的唯一标识进行审批等的判定

针对【表单】的处理
- 工作流引擎本身不保存任何表单的结构和数据
- 如果流转中有一些网关的条件需要用到表单的数据，将表单中的判断字段赋值给variable，然后条件表达式中使用该变量来判断

## 支持的bpmn元素

- 事件(Event)
   - 开始事件(StartEvent)
   - 结束事件(EndEvent)
- 活动(Activity)
   - 用户任务(UserTask)
   - 脚本任务(ScriptTask)
- 网关(Gateway)
   - 排他网关(ExclusiveGateway)
   - 并行网关(ParallelGateway)
   - 包容网关(InclusiveGateway)
- 顺序流(SequenceFlow)

## 技术架构

golang + echo + gorm + postgres

### 外部组件依赖

除了数据库(postgres)之外没有其他的依赖
- mysql需要替换gorm的驱动，以及对一些sql做小幅改造

### 数据访问

gorm + postgres

使用指定数据库的wf schema。可以单独使用一个数据库，也可以集成到已有业务的数据库中

支持自动迁移(可配置开启或关闭)

## 支持的功能

- 上述bpmn元素应该支持的所有基础功能
- 会签
- 转交审批
- 审批限时 自然日/工作日
   - 超时后果 (自动通过/拒绝 或者无操作)
- WebHook
- 多租户
- 展示用的流程链路

## 质量保证

- 集成测试