package basedef

// VERSION - jarviscore version
const VERSION = "0.7.60"

// TimeFailedServAddr - if the servaddr is failed, we won't try to connect it within this time
const TimeFailedServAddr = 60

// NumsReconnectFailedServAddr - if the servaddr is failed, we will retry this number before it is marked as failed
const NumsReconnectFailedServAddr = 3
