package basedef

// VERSION - jarviscore version
const VERSION = "0.7.143"

// TimeFailedServAddr - if the servaddr is failed, we won't try to connect it within this time
const TimeFailedServAddr = 60

// NumsReconnectFailedServAddr - if the servaddr is failed, we will retry this number before it is marked as failed
const NumsReconnectFailedServAddr = 3

// LastTimeDeprecated - This is a constant about timestampDeprecated
var LastTimeDeprecated = []int64{30, 5 * 60, 30 * 60, 60 * 60, 2 * 60 * 60, 4 * 60 * 60, 8 * 60 * 60, 16 * 60 * 60, 24 * 60 * 60}

// BigFileLength -if file length >= BigFileLength, the file is bigfile
const BigFileLength = 4*1024*1024 - 1024

// TimeReconnect - If the LastConnectTime time is within 30s, we will not repeat the connection.
const TimeReconnect = 30
