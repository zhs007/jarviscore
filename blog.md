# JarvisCore Development Log

### 2019-05-05

异步以后，出现一个新问题，就是sendmsg的返回可能比服务器异步以后的返回慢。  
这个问题有点麻烦，不过目前来说，只可能是IGOTIT比END慢，其实大部分逻辑不影响的，先简单处理了一下，后面再来调整吧。

### 2019-05-04

如果jarvissh不是docker运行的话，以前是建议将日志调整为文件格式的，但发现那样的话，就没办法捕获panic了，前面试着捕获了一下panic，但那个方法需要对每个goroutine加处理，比较麻烦。  
今天加了新的处理，将stderr映射到一个文件里，这样应该就可以正确的处理panic日志了。

今天还将所有的请求调整为异步请求了，主要是简化逻辑层，然后就是保证了请求语义的一致性。  
系统级的协议，最好还是不走异步请求，包括连接和获取msgstate。

### 2019-05-02

接下来，有几个比较重要的调整，记录一下：

- chatbot的内核需要独立服务处理，而telegram、微信这些聊天app的对接服务放更外层，这样可以一个chatbot对应多个聊天app，且可以绑定多个账号关系，``Jarvis无处不在``。
- 文件系统需要提上日程了。
- 日志系统需要进一步简化。
- 分布式的配置系统，基于文件系统吧。
- 分布式的资源存储。

``procMsgResultMgr``可能的内存泄露问题，今天发现还有没彻底处理干净的。  
这个主要是为了能正确处理回调，所以需要保证结束消息被收到，且连接断开，但现在有个可能，就是中间可能会报错，如果是网络错误，重连就能好，但如果是远端重启或者别的什么原因，可能就永远也收不到结束消息了。  
这时是不是应该每个进程有个liveid之类的，如果发现liveid变了，就清理掉老的数据。  
或者这个replyid应该记下来，尽可能维护对。  
但replyid记下来，如果是进程结束，有些数据也不太可能维护对。  
还有个方案，就是隔一段时间，是不是能主动询问一下，这样也好知道远端节点是否还在处理，或者是给个进度条之类的，如果是shell，没有进度条，是不是能给个console日志条数，哪怕不能知道百分比的进度，但能够隔一段时间知道任务在推进，也是件好事吧。  
估计会按最后的方案来执行。

### 2019-04-29

端口约定：

- 7788: jarvissh docker环境
- 7700: jarvissh 宿主环境
- 6061: pprof
- 7051: jarviscrawlercore
- 7100: dtdataserv的jarvisnode端口
- 7101: dtdataserv的http端口

今天把dtdataserv调通了，dtdataserv和jarviscrawler其实就是2个完全不同的思路了，jarviscrawler是独立于Jarvis节点以外的服务，而dtdataserv是JarvisNode，2者各有各的好处吧。  
本质上，dtdataserv还是用到了jarviscrawler，所以实际上不是一个竞争关系。  
如果需要更多的使用JarvisNode特性，譬如主动信任机制（安全性）、统一的自动更新机制、节点消息订阅推送（还未完成），这些特性时，基于JarvisNode会简单很多。  
当然，其实还有个很重要的原因，就是现在JarvisNode没有非golang的版本，其它语言的，目前还是独立开发的好，后面再用Golang套一层，更省事一些。

### 2019-04-28

关于``Jarvis``，首先，``Jarvis``不是一个``chatbot``，它是一个复杂的多层系统，最底层是一个分布式的运算网络，使用了BTC的加密验证算法，并通过信任系统保证本机算力不被滥用。  
然后，在这个分布式的运算网络上面，是一个应用层，应用层利用``Jarvis``网络算力来进行具体业务，现在被实际应用的部分主要有3个，一个是``chatops``，就是用来做服务器运维的，包括各种共有库私有库的更新，单元测试，发布等。  
然后是企业私有机器人助理，这部分包括打通各个独立的系统、数据监控、数据分析等。  
最后，是具体的bot应用，包括数据爬虫、对接各种API等，现在的新闻等，就属于这个部分。  
计划中，在底层网络和应用层之间，还应该有一层货币层，现在还没实现的，我的想法主要是为了给应用层一个清晰的货币激励，譬如我发布应用，我希望召集很多个节点一起来做点事，我需要有很方便的自动化的激励措施。  
这个货币层，会和现在的区块链不太一样，我不需要有一个通用货币，我只提供一个方便应用自行发布货币的功能，然后提供完善的市商节点，让很多货币能流通起来，这样单货币之间其实并不需要特别复杂的共识，交易效率会高很多。

### 2019-04-27

这几天一直在关注pprof，goroutine数量其实基本上维持在90多，内存的问题，开了2天才有点头绪，看起来应该是grpc接收stream数据时，某些情况下服务端goroutine没结束造成的。  

```
#	0x848870	github.com/zhs007/jarvissh/vendor/github.com/zhs007/jarviscore/proto.(*jarvisCoreServProcMsgClient).Recv+0x30	/go/src/github.com/zhs007/jarvissh/vendor/github.com/zhs007/jarviscore/proto/jarviscore.pb.go:1926
#	0x96ae0f	github.com/zhs007/jarvissh/vendor/github.com/zhs007/jarviscore.(*jarvisClient2)._sendMsg+0xa6f			/go/src/github.com/zhs007/jarvissh/vendor/github.com/zhs007/jarviscore/client2.go:316
#	0x9693af	github.com/zhs007/jarvissh/vendor/github.com/zhs007/jarviscore.(*clientTask).Run+0x8df				/go/src/github.com/zhs007/jarvissh/vendor/github.com/zhs007/jarviscore/client2.go:42
#	0x84bf76	github.com/zhs007/jarvissh/vendor/github.com/zhs007/jarviscore/base.(*l2routine).start+0x406			/go/src/github.com/zhs007/jarvissh/vendor/github.com/zhs007/jarviscore/base/l2routinepool.go:99
#	0x84dd70	github.com/zhs007/jarvissh/vendor/github.com/zhs007/jarviscore/base.(*l2routinePool).startRountine+0x90		/go/src/github.com/zhs007/jarvissh/vendor/github.com/zhs007/jarviscore/base/l2routinepool.go:344
```

感觉和chanEnd有关，这块本来实现就有点纠结，估计要找个时间再理一遍才好。  
后来发现没有goroutine卡主，所以应该不是卡recv了，而是msg没有被释放，因为前面有个回调的处理，如果error，而不是正常eof，是可能泄露的。

然后，切换到了go module，没有用dep了，有个小问题，就是目前没有发现go module能用非tag的方式精确定位版本，以前dep可以定位branch的，所以接下来可能一段时间tag会多一点。

再就是新开了一个特殊功能节点的项目，最初没打算加到jarvisnode节点里来的，本来是打算类似crawler server那样的方式，直接grpc服务，简单很多，后来考虑到权限各种配置，就干脆加进来好了，试试看这条路是否能走通。  
所以，会有些更新是专门为这种节点缺的接口准备的。

### 2019-04-24

前几天发现了内存泄露问题，今天加了pprof。

### 2019-04-02

节点的信任关系，暂时按简单的方案处理了，信任关系需要写在配置文件里。

### 2019-03-31

前几天一直在处理crawler，昨天开始准备做jarvisnote了，本来想找个全文搜索引擎来做的，最初选的是beleve，golang写的，但查issue的话，发现其实跟进很慢，而且效率非常低，于是又找了一圈，最后发现最靠谱的可能还是elasticsearch，基本确定后，突然又想到，其实我不需要这么复杂的全文搜索引擎，前期可以自己做个简单的，这样更可控一些。  
再说，后面处理对话时，还是要折腾这些的，于是就还是先分词建索引吧。

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
