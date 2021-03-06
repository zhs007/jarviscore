package basedef

// VERSION - jarviscore version
const VERSION = "0.7.212"

// TimeFailedServAddr - if the servaddr is failed, we won't try to connect it within this time
const TimeFailedServAddr = 60

// NumsReconnectFailedServAddr - if the servaddr is failed, we will retry this number before it is marked as failed
const NumsReconnectFailedServAddr = 3

// LastTimeDeprecated - This is a constant about timestampDeprecated
var LastTimeDeprecated = []int64{
	30,
	5 * 60,
	30 * 60,
	60 * 60,
	2 * 60 * 60,
	4 * 60 * 60,
	8 * 60 * 60,
	16 * 60 * 60,
	24 * 60 * 60,
}

// BigFileLength -if file length >= BigFileLength, the file is bigfile
const BigFileLength = 4*1024*1024 - 1024

// TimeReconnect - If the LastConnectTime time is within 30s, we will not repeat the connection.
const TimeReconnect = 30

// BigMsgLength -if jarvismsg length >= BigMsgLength, the message is big message
const BigMsgLength = 4*1024*1024 - 1024

// TimeMsgState - Each time interval, it will check msgstate
const TimeMsgState = 60

// TimeClearEndMsgState - If msg has already been processed, its end state will remain for a while, avoiding concurrent bugs.
const TimeClearEndMsgState = 15 * 60

// BigLogFileLength -if command log file length >= BigLogFileLength, the message is big message
const BigLogFileLength = 1 * 1024 * 1024
