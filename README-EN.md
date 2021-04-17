# tumbleweed workflow engine

[简体中文](./README.md) | English

## Introduction

Workflow engine based on bpmn 2.0, split users, forms, etc., and focus on the flow of the process itself

Use with workflow designer
> https://github.com/lzw5399/tumbleweed-designer

## Project positioning

`Tumbleweed` is a workflow independent microservice to provide RESTful services (subsequent to consider changing to grpc-gateway form). Separate dependencies such as [user/role] and [form] to support multi-tenancy

Processing for [user/role]
- Designate a tenant for the system connected to the workflow engine
- Then synchronize the id and name of the user/role that needs to use the workflow in the system to the database of the workflow engine
- The workflow engine itself only cares about the unique identification of the user/role for approval, etc.

Processing for [Form]
- The workflow engine itself does not save any form structure and data
- If there are some gateway conditions in the circulation that need to use form data, assign the judgment field in the form to variable, and then use the variable in the condition expression to judge

## Supported bpmn elements

- Event
   - StartEvent
   - EndEvent
- Activity
   - UserTask
   - ScriptTask
- Gateway
   - ExclusiveGateway
   - ParallelGateway
   - InclusiveGateway
- SequenceFlow

## Technology Architecture

golang + echo + gorm + postgres

### External component dependencies

No other dependencies except the database (postgres)
-mysql needs to replace the driver of gorm and make a small modification to some sql

### data access

gorm + postgres

Use the wf schema of the specified database. A database can be used alone or integrated into an existing business database

Support automatic migration (configurable on or off)

## Supported features

- All the basic functions that the above bpmn element should support
- Countersign
- Referral for approval
- Approval time limit Natural day/working day
   - Timeout consequences (automatically pass/reject or no action)
- WebHook
- Multi-tenancy
- Process link for display

## quality assurance

- Integration Testing