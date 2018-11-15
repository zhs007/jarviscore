# jarviscore
[![Build Status](https://travis-ci.org/zhs007/jarviscore.svg?branch=master)](https://travis-ci.org/zhs007/jarviscore)

Jarvis是一个共享网络的私人助理，他的名字取自漫威里钢铁侠的私人助理。

Jarvis初衷是用于管理多台服务器，或者说是多台服务器有效的利用资源。

jarviscore是Jarvis网络节点内核的golang实现。

### 版本更新

0.6
- 签名校验规则简化
- 增加event通知接口
- 增加script file类型的message

0.5
- 网络协议重构，彻底抛弃前面的网络层实现
- msg处理采用goroutine pool方式
- 统一server端和client端的事务处理
- 实现部分test用例

0.3
- 网络协议升级
- 对接telegram chatbot
- 支持telegram聊天查询

0.2  
- 数据存储切换为AnkaDB
- 增加GraphQL数据查询
- 增加GraphIQL支持

0.1  
- 完成基本功能，建立网络
- BTC加密算法
- leveldb的数据存储