# JarvisCore Development Log

### 2019-03-25

今天终于把昨天折腾一天的bug给找到了，最后发现是我为了省内存，在ProcFileData里循环内部有个buf复用，如果这个数据没有及时被复制出去，其实就会变化的！

然后，这几天还折腾了一下爬虫，图省事，直接用的pdf，估计接下来还是得分析dom才好，其实现在article这边需求没那么复杂。  

### 2019-03-23

今天把前几天blog里提到的END和WAITPUSH这些处理了。  
双向的异步请求回应都能正常关联上了。  
其实可能还有更复杂的request->reply->request->reply这样的结构，但按道理可以在逻辑层回避掉。  
先不考虑这么多。

### 2019-03-21

今天开始新一轮的重构了，支持了双向流。  

消息流分为3种，分别是单个消息、普通消息流 和 流式消息。  
普通消息流，就是这一组消息可以一起处理，所以会占用内存，等消息结束后一起处理。  
流式消息，目前想到的主要是文件传输，因为如果文件太大，这部分完全可以流式直接处理掉，没必要在内存里收完最后一起处理。  
今天其实做的是普通消息流，流式消息后面再加吧，短期内也没需求。  

因为有些耗时操作，譬如runscript这样的，不是收到request后直接reply的，以前希望通过msgid来把request和reply绑定在一起，特地增加了END来标识结束，今天加了WAITPUSH，想干脆把2种请求分开，如果是异步reply的，先告诉另一端，这样缓存处理简单一些。

### 2019-03-18

今天把大文件传输正式部署到线上了，今天的Jarvis更新就是用的它自己传输的文件来做的。  

接下来需要把文件系统彻底整理一下，现在在telegram chatbot这边，文件都是放ankadb里的，感觉有点不合适，特别是支持大文件以后。  

文件系统现在的想法，其实还是HASH文件名，避免文件重复，然后数据库存文件映射关系就好。

### 2019-03-17

今天找到``ErrInvalidMsgID``bug的原因了，因为现在其实不同节点直接交互有2种渠道，分别是响应请求的stream和发起请求的sendmsg，这2种渠道是无关的，接收顺序不被保障，所以当2边同时进行时，就可能``ErrInvalidMsgID``。  
暂时没有很好的解决，只是先将stream部分不处理msgid了，后面再来解决。  

然后就是发现今天``updnodes``的效率比以前要高很多，很快就能把节点更新完，不知道是不是前几天优化``L2RoutinePool``的原因。

### 2019-03-16

今天觉得需要添加文件系统，有个比较简单的方案。  
这个不重要。

### 2019-03-15

最近一直在关注输出，并解决了一批小buf，最新的版本要稳定很多。  
这几天还优化了日志，现在INFO级别的日志也变得可用了。

### 2019-03-10

我在昨天晚上终于发现一个测试用例可能报错的bug了，就是普通节点在请求数据时，可能还没有收到root节点的返回。  
今天修正这个bug，但又发现一个可能root节点未listen成功，普通节点就可能去连接，导致连接失败。  
其实这些bug都是测试用例的bug。

I fixed some bugs in the test case.  

### 2019-03-09

前面遇到的几个自动更新bug都被修正了。  
将requestnodes测试用例的写法切换到updnodes了。  
这几天还在大刀阔斧的重构，以前有些实现冗余了，会逐步删除掉。  

I fixed some bugs for the automatic update today.  
I rewrote the test case for requestnodes.  
I have started a new refactoring.  

### 2019-03-06

今天Jarvis的自动更新已经初步完成，有几台资源占用较高的服务器未能自动更新成功，原因待查。  
这几天发现前面的一些设计有些混乱，准备新一轮的重构。  
还有测试用例的写法，updnode将是一个新的范例。  

对Jarvis的定位更加清晰了，Jarvis是帮助普通人用算法进行决策。  

I completed the automatic update for Jarvis today, but several servers failed to update (their CPU usage is high), and I have to spend more time looking for the cause of the problem.    
I found some of the previous designs to be bad and prepared for a new round of refactoring.  
I have completed a new test case today, and updnode is a better way to write.  

I think Jarvis is used to help people make decisions with algorithms.

### 2019-03-02

今天将线上的jarvissh和jarvistelebot升级到0.7.10了，初步测试通过。  
jarvissh部署流程比以前简单了。  
今天升级全部服务器其实还比较简单，但接下来需要加入自动更新的功能。

Today I upgraded the jarvissh and jarvistelebot kernels to v0.7.10.  
The deployment of jarvis is simpler than before.  
Today, I didn't spend much time upgrading all the servers. It would be better if I could add the automatic update feature.
