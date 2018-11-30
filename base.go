package jarviscore

import (
	"context"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc/peer"
)

// VERSION - jarviscore version
const VERSION = "0.6.19"

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
		return nil, ErrLoadFileReadSize
	}

	return buffer, nil
}

func getGRPCClientIP(ctx context.Context) string {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	if pr.Addr == net.Addr(nil) {
		return ""
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	return addSlice[0]
}

// IsMyServAddr - check destaddr is same addr for me
func IsMyServAddr(destaddr string, srcaddr string) bool {
	dh, dp, err := net.SplitHostPort(destaddr)
	if err != nil {
		return false
	}

	sh, sp, err := net.SplitHostPort(srcaddr)
	if err != nil {
		return false
	}

	if dp == sp && (dh == sh || dh == "127.0.0.1") {
		return true
	}

	return false
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
