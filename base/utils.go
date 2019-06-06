package jarvisbase

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

// JSON - make json to field
func JSON(key string, obj interface{}) zap.Field {
	s, err := json.Marshal(obj)
	if err != nil {
		return zap.Error(err)
	}

	return zap.String(key, string(s))
}

// BuildLogSubFilename -
func BuildLogSubFilename(appName string, version string) string {
	tm := time.Now()
	str := tm.Format("2006-01-02_15:04:05")
	return fmt.Sprintf("%v.%v.%v", appName, version, str)
}

// BuildLogFilename -
func BuildLogFilename(logtype string, subname string) string {
	return fmt.Sprintf("%v.%v.log", subname, logtype)
}

// MD5Protobuf - md5 protobuf and return string
func MD5Protobuf(pb proto.Message) (string, error) {
	buf, err := proto.Marshal(pb)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(buf)), nil
}
