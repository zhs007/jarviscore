package jarviscore

import (
	"context"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc/peer"
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
func IsMyServAddr(destaddr string, myaddr string) bool {
	if destaddr == myaddr {
		return true
	}

	dh, dp, err := net.SplitHostPort(destaddr)
	if err != nil {
		return false
	}

	sh, sp, err := net.SplitHostPort(myaddr)
	if err != nil {
		return false
	}

	if dp == sp && (dh == sh || dh == "127.0.0.1") {
		return true
	}

	return false
}

// IsValidServAddr - is the server address is valid?
func IsValidServAddr(servaddr string) bool {
	arr := strings.Split(servaddr, ":")
	if len(arr) != 2 {
		return false
	}

	if arr[0] == "" || arr[1] == "" {
		return false
	}

	_, _, err := net.SplitHostPort(servaddr)
	if err != nil {
		return false
	}

	return true
}

// IsLocalHostAddr - is localhost
func IsLocalHostAddr(servaddr string) bool {
	arr := strings.Split(servaddr, ":")
	if len(arr) != 2 {
		return false
	}

	if arr[0] == "" {
		return false
	}

	if arr[0] == "127.0.0.1" || arr[0] == "localhost" {
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
