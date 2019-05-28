package jarvisbase

import (
	"encoding/json"
	"fmt"
	"time"

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

// BuildLogSubFileName -
func BuildLogSubFileName(appName string, version string) string {
	tm := time.Unix(curtime, 0)
	str := tm.Format("2006-01-02_15:04:05")
	return fmt.Sprintf("%v.%v.%v", appName, version, str)
}
