# AnkaDB Development Log

### 2019-03-16

今天觉得需要添加文件系统，有个比较简单的方案。  

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
