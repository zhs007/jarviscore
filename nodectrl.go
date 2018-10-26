package jarviscore

import (
	"bytes"

	"go.uber.org/zap"

	ankadatabase "github.com/zhs007/ankadb/database"
	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
)

const (
	ctrldbAddr       = "addr"
	ctrldbPublicKey  = "pubkey"
	ctrldbCtrlPrefix = "ctrl:"
)

// nodeCtrlInfo -
type nodeCtrlInfo struct {
	srcAddr string
	pubKey  []byte
	ctrldb  ankadatabase.Database
	mapCtrl map[int64]*pb.CtrlDataInDB
}

func newNodeCtrlInfo(addr string) *nodeCtrlInfo {
	return nil
	// db, err := ankadatabase.NewAnkaLDB(getRealPath("ctrl-"+addr), 16, 16)
	// if err != nil {
	// 	jarviserr.ErrorLog("newNodeCtrlInfo:NewAnkaLDB", err)
	// 	return nil
	// }

	// err = db.Put([]byte("addr"), []byte(addr))
	// if err != nil {
	// 	jarviserr.ErrorLog("newNodeCtrlInfo:saveAddr", err)
	// 	return nil
	// }

	// nci := &nodeCtrlInfo{
	// 	srcAddr: addr,
	// 	ctrldb:  db,
	// 	mapCtrl: make(map[int64]*pb.CtrlDataInDB),
	// }

	// err = nci.onInit()
	// if err != nil {
	// 	jarviserr.ErrorLog("newNodeCtrlInfo:saveAddr", err)
	// 	return nil
	// }

	// return nci
}

func (nci *nodeCtrlInfo) onInit() error {
	addr, err := nci.ctrldb.Get([]byte(ctrldbAddr))
	if err != nil {
		err = nci.ctrldb.Put([]byte(ctrldbAddr), []byte(nci.srcAddr))
	} else if bytes.Compare(addr, []byte(nci.srcAddr)) != 1 {
		return ErrInvalidAddr
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
			return ErrInvalidPublishKey
		}

		return nil
	}

	err := nci.ctrldb.Put([]byte(ctrldbPublicKey), pubKey)
	if err != nil {
		jarvisbase.Error("newNodeCtrlInfo:saveAddr", zap.Error(err))

		return err
	}

	return nil
}

func (nci *nodeCtrlInfo) hasCtrl(ctrlid int64) bool {
	if _, ok := nci.mapCtrl[ctrlid]; ok {
		return true
	}

	ok, err := nci.ctrldb.Has([]byte(ctrldbCtrlPrefix + string(ctrlid)))
	if err != nil {
		return false
	}

	return ok
}

func (nci *nodeCtrlInfo) addCtrl(ctrlid int64, ctrltype pb.CTRLTYPE, command []byte, forwordAddr string, forwordNums int32) error {
	// if _, ok := nci.mapCtrl[ctrlid]; ok {
	// 	return ErrExistCtrlID
	// }

	// ctrlindb := &pb.CtrlDataInDB{
	// 	Ctrlid:      ctrlid,
	// 	CtrlType:    ctrltype,
	// 	ForwordAddr: forwordAddr,
	// 	Command:     command,
	// 	ForwordNums: forwordNums,
	// }

	// data, err := proto.Marshal(ctrlindb)
	// if err != nil {
	// 	return jarviserr.NewError(pb.CODE_CTRLDATAINDB_ENCODE_FAIL)
	// }

	// err = nci.ctrldb.Put([]byte(ctrldbCtrlPrefix+string(ctrlid)), data)
	// if err != nil {
	// 	return jarviserr.NewError(pb.CODE_CTRLDB_SAVE_CTRLDATA_FAIL)
	// }

	// nci.mapCtrl[ctrlid] = ctrlindb

	return nil
}

func (nci *nodeCtrlInfo) setCtrlResult(ctrlid int64, result []byte) error {
	// val, ok := nci.mapCtrl[ctrlid]
	// if !ok {
	// 	return jarviserr.NewError(pb.CODE_CTRLDB_NOT_EXIST_CTRLID)
	// }

	// val.Result = result

	// data, err := proto.Marshal(val)
	// if err != nil {
	// 	return jarviserr.NewError(pb.CODE_CTRLDATAINDB_ENCODE_FAIL)
	// }

	// err = nci.ctrldb.Put([]byte(ctrldbCtrlPrefix+string(ctrlid)), data)
	// if err != nil {
	// 	return jarviserr.NewError(pb.CODE_CTRLDB_SAVE_CTRLDATA_FAIL)
	// }

	return nil
}
