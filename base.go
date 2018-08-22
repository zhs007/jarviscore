package jarviscore

import (
	"os"

	"github.com/zhs007/jarviscore/errcode"
)

func loadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileinfo.Size()
	buffer := make([]byte, fileSize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	if int64(bytesread) != fileSize {
		return nil, newError(jarviserrcode.FILEREADSIZEINVALID)
	}

	return buffer, nil
}

// // GetDNSPulicIP -
// func GetDNSPulicIP() string {
// 	conn, err := net.Dial("udp", "8.8.8.8:53")
// 	if err != nil {
// 		return ""
// 	}
// 	defer conn.Close()
// 	localAddr := conn.LocalAddr().String()
// 	idx := strings.LastIndex(localAddr, ":")
// 	return localAddr[0:idx]
// }

// // GetHTTPPulicIP -
// func GetHTTPPulicIP() string {
// 	resp, err := http.Get("http://ifconfig.me")
// 	if err != nil {
// 		return ""
// 	}
// 	defer resp.Body.Close()
// 	// content, _ := ioutil.ReadAll(resp.Body)
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(resp.Body)
// 	//s := buf.String()
// 	return string(buf.String())
// }
