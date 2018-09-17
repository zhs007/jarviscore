package jarviscore

import (
	"github.com/golang/protobuf/proto"
	"github.com/zhs007/jarviscore/crypto"
	"github.com/zhs007/jarviscore/db"
	"github.com/zhs007/jarviscore/err"
	pb "github.com/zhs007/jarviscore/proto"
)

const (
	coredbMyPrivKey        = "myprivkey"
	coredbMyNodeInfoPrefix = "ni:"
)

type coreDB struct {
	db      jarvisdb.Database
	privKey *jarviscrypto.PrivateKey
}

func newCoreDB() (*coreDB, error) {
	db, err := jarvisdb.NewJarvisLDB(getRealPath("coredb"), 16, 16)
	if err != nil {
		jarviserr.ErrorLog("newCoreDB:NewJarvisLDB", err)

		return nil, jarviserr.NewError(pb.CODE_COREDB_OPEN_FAIL)
	}

	return &coreDB{
		db: db,
	}, nil
}

func (db *coreDB) savePrivateKey() error {
	privkey := &pb.PrivateKey{
		PriKey: db.privKey.ToPrivateBytes(),
	}

	data, err := proto.Marshal(privkey)
	if err != nil {
		return err
	}

	err = db.db.Put([]byte(coredbMyPrivKey), data)
	if err != nil {
		return err
	}

	return nil
}

func (db *coreDB) loadPrivateKey() error {
	dat, err := db.db.Get([]byte(coredbMyPrivKey))
	if err != nil {
		db.privKey = jarviscrypto.GenerateKey()

		return db.savePrivateKey()
	}

	pbprivkey := &pb.PrivateKey{}
	err = proto.Unmarshal(dat, pbprivkey)
	if err != nil {
		db.privKey = jarviscrypto.GenerateKey()

		return db.savePrivateKey()
	}

	privkey := jarviscrypto.NewPrivateKey()
	err = privkey.FromBytes(pbprivkey.PriKey)
	if err != nil {
		db.privKey = jarviscrypto.GenerateKey()

		return db.savePrivateKey()
	}

	db.privKey = privkey

	return nil
}

func (db *coreDB) foreachNode(oneach func([]byte, *pb.NodeInfoInDB), errisbreak bool) error {
	iter := db.db.NewIteratorWithPrefix([]byte(coredbMyNodeInfoPrefix))
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		// key := iter.Key()
		value := iter.Value()

		ni2db := &pb.NodeInfoInDB{}
		err := proto.Unmarshal(value, ni2db)
		if err != nil {
			oneach(iter.Key(), ni2db)
		} else if errisbreak {
			return jarviserr.NewError(pb.CODE_NODEINFOINDB_DECODE_FAIL)
		}
	}

	iter.Release()
	err := iter.Error()
	if err != nil {
		return jarviserr.NewError(pb.CODE_LEVELDB_ITER_ERROR)
	}

	return nil
}

func (db *coreDB) saveNode(cni *NodeInfo) error {
	ni := &pb.NodeInfo{
		Name:     cni.baseinfo.Name,
		ServAddr: cni.baseinfo.ServAddr,
		Addr:     cni.baseinfo.Addr,
		NodeType: cni.baseinfo.NodeType,
	}

	ni2db := &pb.NodeInfoInDB{
		NodeInfo:      ni,
		ConnectNums:   int32(cni.connectNums),
		ConnectedNums: int32(cni.connectedNums),
	}

	data, err := proto.Marshal(ni2db)
	if err != nil {
		return jarviserr.NewError(pb.CODE_NODEINFOINDB_ENCODE_FAIL)
	}

	err = db.db.Put(append([]byte(coredbMyNodeInfoPrefix), cni.baseinfo.Addr...), data)
	if err != nil {
		return jarviserr.NewError(pb.CODE_COREDB_NODEINFO_PUT_FAIL)
	}

	return nil
}
