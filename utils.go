package jarviscore

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/peer"

	"github.com/golang/protobuf/proto"

	jarvisbase "github.com/zhs007/jarviscore/base"
	"github.com/zhs007/jarviscore/basedef"
	coredbpb "github.com/zhs007/jarviscore/coredb/proto"
	pb "github.com/zhs007/jarviscore/proto"
)

// GetMD5String - md5 buf and return string
func GetMD5String(buf []byte) string {
	return fmt.Sprintf("%x", md5.Sum(buf))
}

// MD5File - md5 file and return string
func MD5File(fn string) (string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return "", nil
	}

	return fmt.Sprintf("%x", md5hash.Sum(nil)), nil
}

// GetRealFilename - get filename
func GetRealFilename(fn string) string {
	arr := strings.Split(fn, "/")
	return arr[len(arr)-1]
}

// AbsInt64 - abs int64
func AbsInt64(num int64) int64 {
	if num < 0 {
		return -num
	}

	return num
}

// IsTimeOut - is JarvisMsg timeout
func IsTimeOut(msg *pb.JarvisMsg) bool {
	ct := time.Now().Unix()
	if AbsInt64(ct-msg.CurTime) < 5*60 {
		return false
	}

	return true
}

// IsValidNodeName - check node name
func IsValidNodeName(nodename string) bool {
	if len(nodename) == 0 {
		return false
	}

	for i, v := range nodename {
		if i == 0 && ((v >= '0' && v <= '9') || v == '_') {
			return false
		}

		if !((v >= '0' && v <= '9') || (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || v == '_') {
			return false
		}
	}

	if nodename[len(nodename)-1] == '_' {
		return false
	}

	return true
}

// StoreLocalFile - store filedata to local file systems
func StoreLocalFile(file *pb.FileData) error {
	f, err := os.Create(file.Filename)
	if err != nil {
		jarvisbase.Warn("StoreLocalFile:os.Create err", zap.Error(err))

		return err
	}

	defer f.Close()

	f.Write(file.File)
	f.Close()

	return nil
}

// StoreLocalFileEx - store filedata array to local file systems
func StoreLocalFileEx(files []*pb.FileData) error {
	f, err := os.Create(files[0].Filename)
	if err != nil {
		jarvisbase.Warn("StoreLocalFileEx:os.Create err", zap.Error(err))

		return err
	}

	defer f.Close()

	for i := 0; i < len(files); i++ {
		f.Write(files[i].File)
	}

	f.Close()

	return nil
}

// GetNodeBaseInfo - get nodebaseinfo from nodeinfo
func GetNodeBaseInfo(node *coredbpb.NodeInfo) *pb.NodeBaseInfo {
	return &pb.NodeBaseInfo{
		ServAddr:        node.ServAddr,
		Addr:            node.Addr,
		Name:            node.Name,
		NodeTypeVersion: node.NodeTypeVersion,
		NodeType:        node.NodeType,
		CoreVersion:     node.CoreVersion,
	}
}

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

// GetFileLength - get file length
func GetFileLength(fn string) (int64, error) {
	fileInfo, err := os.Stat(fn)
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

// DeepCopy - deep copy
func DeepCopy(src proto.Message, dest proto.Message) error {
	bret, err := proto.Marshal(src)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(bret, dest)
	if err != nil {
		return err
	}

	return nil
}

// CountMD5String - count MD5 string
func CountMD5String(lst []*pb.FileData) (int64, string, error) {
	totallen := int64(0)
	md5hash := md5.New()
	for i := 0; i < len(lst); i++ {
		cl, err := md5hash.Write(lst[i].File)
		if err != nil {
			return 0, "", err
		}

		totallen = totallen + int64(cl)
	}

	return totallen, fmt.Sprintf("%x", md5hash.Sum(nil)), nil
}

// FuncOnFileData - on proc filedata
type FuncOnFileData func(fd *pb.FileData, isend bool) error

// ProcFileDataWithBuff - proc filedata
func ProcFileDataWithBuff(buf []byte, onfunc FuncOnFileData) error {
	if onfunc == nil {
		return ErrNoFuncOnFileData
	}

	fl := len(buf)

	if fl <= basedef.BigFileLength {

		md5str := GetMD5String(buf)

		onfunc(&pb.FileData{
			File:          buf,
			Ft:            pb.FileType_FT_BINARY,
			Start:         0,
			Length:        int64(fl),
			TotalLength:   int64(fl),
			FileMD5String: md5str,
			Md5String:     md5str,
		}, true)

	} else {
		curstart := int64(0)
		curlength := int64(basedef.BigFileLength)

		for curstart < int64(fl) {
			if curstart+curlength >= int64(fl) {

				onfunc(&pb.FileData{
					File:          buf[curstart:fl],
					Ft:            pb.FileType_FT_BINARY,
					Start:         curstart,
					Length:        int64(int64(fl) - curstart),
					TotalLength:   int64(fl),
					FileMD5String: GetMD5String(buf),
					Md5String:     GetMD5String(buf[curstart:fl]),
				}, true)

				return nil

			}

			onfunc(&pb.FileData{
				File:        buf[curstart:(curstart + curlength)],
				Ft:          pb.FileType_FT_BINARY,
				Start:       curstart,
				Length:      int64(curlength),
				TotalLength: int64(fl),
				Md5String:   GetMD5String(buf[curstart:(curstart + curlength)]),
			}, false)

			curstart = curstart + curlength
		}
	}

	return nil
}

// ProcFileData - proc filedata
func ProcFileData(fn string, onfunc FuncOnFileData) error {
	if onfunc == nil {
		return ErrNoFuncOnFileData
	}

	fdata, err := os.Open(fn)
	if err != nil {
		return err
	}

	defer fdata.Close()

	fs, err := fdata.Stat()
	if err != nil {
		return err
	}

	fl := fs.Size()

	if fl <= basedef.BigFileLength {

		buf := make([]byte, fl)

		rn, err := fdata.Read(buf)
		if err != nil {
			return err
		}

		if int64(rn) != fl {
			return ErrLoadFileReadSize
		}

		totalmd5 := GetMD5String(buf)

		onfunc(&pb.FileData{
			File:          buf,
			Ft:            pb.FileType_FT_BINARY,
			Start:         0,
			Length:        int64(rn),
			TotalLength:   fl,
			FileMD5String: totalmd5,
			Md5String:     totalmd5,
		}, true)

	} else {
		curstart := int64(0)
		curlength := int64(basedef.BigFileLength)
		// buf := make([]byte, curlength)

		for curstart < int64(fl) {
			if curstart+curlength >= int64(fl) {

				buf := make([]byte, int64(fl)-curstart)

				rn, err := fdata.Read(buf)
				if err != nil {
					return err
				}

				if curstart+int64(rn) != fl {
					return ErrLoadFileReadSize
				}

				totalmd5, err := MD5File(fn)
				if err != nil {
					return err
				}

				onfunc(&pb.FileData{
					File:          buf[0:rn],
					Ft:            pb.FileType_FT_BINARY,
					Start:         curstart,
					Length:        int64(rn),
					TotalLength:   fl,
					FileMD5String: totalmd5,
					Md5String:     GetMD5String(buf[0:rn]),
				}, true)

				return nil

			}

			buf := make([]byte, curlength)

			rn, err := fdata.Read(buf)
			if err != nil {
				return err
			}

			if int64(rn) != curlength {
				return ErrLoadFileReadSize
			}

			onfunc(&pb.FileData{
				File:        buf[0:rn],
				Ft:          pb.FileType_FT_BINARY,
				Start:       curstart,
				Length:      int64(rn),
				TotalLength: fl,
				Md5String:   GetMD5String(buf[0:rn]),
			}, false)

			curstart = curstart + curlength
		}
	}

	return nil
}

// IsValidNodeAddr - is valid nodeaddr
func IsValidNodeAddr(addr string) bool {
	return len(addr) == 34
}

// AppendString - append string
func AppendString(strs ...string) string {
	var buffer bytes.Buffer

	for _, str := range strs {
		if len(str) > 0 {
			buffer.WriteString(str)
		}
	}

	return buffer.String()
}

// GetCurPath - get current path
func GetCurPath() string {
	file, _ := exec.LookPath(os.Args[0])
	p := filepath.Dir(file)
	return p
}

// GetAbsPath - get absolutely patg
func GetAbsPath(p string) string {
	fp := path.Join(GetCurPath(), p)
	return path.Clean(fp)
}

// RunCommand - run command
func RunCommand(ctx context.Context, jnode JarvisNode, cmdname string,
	cmd string) (string, string, error) {

	md5str := GetMD5String([]byte(cmd))
	tm := time.Now()
	timestr := tm.Format("2006-01-02_15:04:05")

	nc := exec.Command("sh", "-c", cmd)

	var outfn string
	var errfn string

	if jnode != nil {
		logpath := jnode.GetConfig().Log.LogPath

		outfn = path.Join(logpath,
			fmt.Sprintf("cmd.out.%v.%v.log", md5str, timestr))
		errfn = path.Join(logpath,
			fmt.Sprintf("cmd.err.%v.%v.log", md5str, timestr))

		jt := &pb.JarvisTask{
			Name:     cmdname,
			TaskType: pb.TASKTYPE_NORMAL,
			CurTime:  tm.Unix(),
			LogFiles: make(map[string]string),
		}

		jt.LogFiles["log"] = outfn
		jt.LogFiles["err"] = outfn

		jnode.GetCoreDB().AddTask(ctx, jt)
	} else {
		logpath := "./"

		outfn = path.Join(logpath,
			fmt.Sprintf("cmd.out.%v.%v.log", md5str, timestr))
		errfn = path.Join(logpath,
			fmt.Sprintf("cmd.err.%v.%v.log", md5str, timestr))
	}

	var errStdout, errStderr error
	stdoutIn, _ := nc.StdoutPipe()
	stderrIn, _ := nc.StderrPipe()

	cmdstd, err := NewCMDStdOutErr(outfn, errfn)
	if err != nil {
		jarvisbase.Warn("RunCommand:NewCMDStdOutErr",
			zap.Error(err))

		return "", "", err
	}

	go func() {
		errStdout = cmdstd.countStdOut(stdoutIn)
	}()

	go func() {
		errStderr = cmdstd.countStdErr(stderrIn)
	}()

	nc.Start()

	err = nc.Wait()
	if err != nil {
		return "", "", err
	}

	if errStdout != nil {
		jarvisbase.Warn("RunCommand:errStdout",
			zap.Error(errStdout))

		// return "", "", errStdout
	}

	if errStderr != nil {
		jarvisbase.Warn("RunCommand:errStderr",
			zap.Error(errStderr))

		// return "", "", errStderr
	}

	outbuf, err := cmdstd.GetStdOut()
	if err != nil {
		jarvisbase.Warn("RunCommand:GetStdOut",
			zap.Error(err))

		return "", "", err
	}

	errbuf, err := cmdstd.GetStdErr()
	if err != nil {
		jarvisbase.Warn("RunCommand:GetStdErr",
			zap.Error(err))

		return "", "", err
	}

	return string(outbuf), string(errbuf), nil
}
