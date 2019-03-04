# jarviscore
[![Build Status](https://travis-ci.org/zhs007/jarviscore.svg?branch=master)](https://travis-ci.org/zhs007/jarviscore)

Jarvis是一个共享网络的私人助理，他的名字取自漫威里钢铁侠的私人助理。  

Jarvis初衷是用于管理多台服务器，合理分配每个节点的算力。  
Jarvis可以用于小型公司或团体的服务器管理，可以配合GitLab做CI，进行分布式的运算，大数据处理，可视化，以及知识库管理等。  

Jarvis提供一个chatbot，供大家控制网络服务。

jarviscore是Jarvis网络内核的golang实现。

### 版本更新

0.7
- 完善MessageID
- 完善Event
- 增加了Request回调
- 支持大于4MB的消息传递
- 进一步完善信任体系
- 加入group，构建新的信任组
- 新增可视化数据服务Viewer
- 新增独立的速记服务Note
- 新增NLP服务器节点
- 新增通用爬虫节点

0.6
- 签名校验规则简化
- 增加event通知接口
- 增加script file类型的message
- 初步实现分布式运算网络
- Telegram下的Jarvis已经可以支持各种远程指令调用、文件传输等

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