package jarviscore

import (
	"bytes"

	"github.com/golang/protobuf/proto"
	"github.com/zhs007/jarviscore/db"
	"github.com/zhs007/jarviscore/err"
	pb "github.com/zhs007/jarviscore/proto"
)

const (
	ctrldbAddr      = "addr"
	ctrldbPublicKey = "pubkey"
	ctrldbCtrl      = "ctrl:"
)

// nodeCtrlInfo -
type nodeCtrlInfo struct {
	srcAddr string
	pubKey  []byte
	ctrldb  jarvisdb.Database
	mapCtrl map[int64]*pb.CtrlDataInDB
}

func newNodeCtrlInfo(addr string) *nodeCtrlInfo {
	db, err := jarvisdb.NewJarvisLDB(getRealPath("ctrl"+addr), 16, 16)
	if err != nil {
		jarviserr.ErrorLog("newNodeCtrlInfo:NewJarvisLDB", err)
		return nil
	}

	err = db.Put([]byte("addr"), []byte(addr))
	if err != nil {
		jarviserr.ErrorLog("newNodeCtrlInfo:saveAddr", err)
		return nil
	}

	nci := &nodeCtrlInfo{
		srcAddr: addr,
		ctrldb:  db,
		mapCtrl: make(map[int64]*pb.CtrlDataInDB),
	}

	err = nci.onInit()
	if err != nil {
		jarviserr.ErrorLog("newNodeCtrlInfo:saveAddr", err)
		return nil
	}

	return nci
}

// func loadNodeCtrlInfo(filename string) (*nodeCtrlInfo, error) {
// 	buf, err := loadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}

// 	nodectrlinfo := newNodeCtrlInfo()
// 	err1 := yaml.Unmarshal(buf, nodectrlinfo)
// 	if err1 != nil {
// 		return nil, err1
// 	}

// 	return nodectrlinfo, nil
// }

// func (m *nodeCtrlInfo) save(filename string) error {
// 	d, err := yaml.Marshal(m)
// 	if err != nil {
// 		return err
// 	}

// 	ioutil.WriteFile(filename, d, 0755)

// 	return nil
// }

// func (m *nodeCtrlInfo) clear() {
// 	if len(m.Mapctrl) == 0 {
// 		return
// 	}

// 	for k := range m.Mapctrl {
// 		delete(m.Mapctrl, k)
// 	}

// 	m.Mapctrl = nil
// 	m.Mapctrl = make(map[int32]*CtrlInfo)
// }

func (nci *nodeCtrlInfo) onInit() error {
	addr, err := nci.ctrldb.Get([]byte(ctrldbAddr))
	if err != nil {
		err = nci.ctrldb.Put([]byte(ctrldbAddr), []byte(nci.srcAddr))
	} else if bytes.Compare(addr, []byte(nci.srcAddr)) != 1 {
		return jarviserr.NewError(pb.CODE_INVALID_ADDR)
	}

	pubkey, err := nci.ctrldb.Get([]byte(ctrldbPublicKey))
	if err != nil {
		return nil
	}

	nci.pubKey = pubkey

	return nil
}

func (nci *nodeCtrlInfo) setPublicKey(pubKey []byte) error {
	if len(nci.pubKey) != 0 {
		if bytes.Compare(nci.pubKey, pubKey) != 1 {
			return jarviserr.NewError(pb.CODE_INVALID_PUBLICKEY)
		}

		return nil
	}

	err := nci.ctrldb.Put([]byte(ctrldbPublicKey), pubKey)
	if err != nil {
		jarviserr.ErrorLog("newNodeCtrlInfo:saveAddr", err)
		return err
	}

	return nil
}

func (nci *nodeCtrlInfo) hasCtrl(ctrlid int64) bool {
	if _, ok := nci.mapCtrl[ctrlid]; ok {
		return true
	}

	ok, err := nci.ctrldb.Has([]byte("ctrl:" + string(ctrlid)))
	if err != nil {
		return false
	}

	return ok
}

func (nci *nodeCtrlInfo) addCtrl(ctrlid int64, ctrltype pb.CTRLTYPE, command []byte, forwordAddr string, forwordNums int32) error {
	if _, ok := nci.mapCtrl[ctrlid]; ok {
		return jarviserr.NewError(pb.CODE_CTRLDB_EXIST_CTRLID)
	}

	ctrlindb := &pb.CtrlDataInDB{
		Ctrlid:      ctrlid,
		CtrlType:    ctrltype,
		ForwordAddr: forwordAddr,
		Command:     command,
		ForwordNums: forwordNums,
	}

	data, err := proto.Marshal(ctrlindb)
	if err != nil {
		return jarviserr.NewError(pb.CODE_CTRLDATAINDB_ENCODE_FAIL)
	}

	err = nci.ctrldb.Put([]byte(ctrldbCtrl+string(ctrlid)), data)
	if err != nil {
		return jarviserr.NewError(pb.CODE_CTRLDB_SAVE_CTRLDATA_FAIL)
	}

	nci.mapCtrl[ctrlid] = ctrlindb

	return nil
}

func (nci *nodeCtrlInfo) setCtrlResult(ctrlid int64, result []byte) error {
	val, ok := nci.mapCtrl[ctrlid]
	if !ok {
		return jarviserr.NewError(pb.CODE_CTRLDB_NOT_EXIST_CTRLID)
	}

	val.Result = result

	data, err := proto.Marshal(val)
	if err != nil {
		return jarviserr.NewError(pb.CODE_CTRLDATAINDB_ENCODE_FAIL)
	}

	err = nci.ctrldb.Put([]byte(ctrldbCtrl+string(ctrlid)), data)
	if err != nil {
		return jarviserr.NewError(pb.CODE_CTRLDB_SAVE_CTRLDATA_FAIL)
	}

	return nil
}
