# jarviscore
[![Build Status](https://travis-ci.org/zhs007/jarviscore.svg?branch=master)](https://travis-ci.org/zhs007/jarviscore)

Jarvis is an open source personal assistant. His name was inspired by the personal assistant of Marvel Iron Man.

Jarvis was originally used to manage multiple servers, making it easier to manage multiple servers.  
Jarvis can be used for server management in small companies or groups. It can be used with GitLab for CI, distributed computing, big data processing, visualization, and knowledge base management.

Jarvis includes a chatbot for everyone to control web services.

Jarviscore is the golang implementation of the Jarvis network kernel.

### Update

#### **v0.7**
- Improved MessageID
- Improved Event module to add more event notifications
- Added Request callback interface
- Support for message delivery greater than 4MB
- Improved trust system
- Added group function to build a new trust group
- Added Visual Data Service Viewer
- Added separate note service Note
- Added NLP server node
- Added universal crawler node

#### **v0.6**
- Simplified signature verification rules
- Add event notification interface
- Added remote call support for scripts
- Can work on dozens of servers
- Jarvis for Telegram can support various remote command calls, file transfers, etc.

#### **v0.5**
- Refactored network kernel
- Added support for goroutine pool
- Unified server and client transaction processing
- Implement some test cases

#### **v0.3**
- Upgraded network protocol (not backward compatible)
- Opened chatbot for telegram
- Support chat query for telegram

#### **v0.2**
- Data storage switched to AnkaDB
- Added GraphQL data query
- Added GraphIQL support

#### **v0.1**
- Complete basic functions and establish a network
- BTC encryption algorithm
- leveldb data storage