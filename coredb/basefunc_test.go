package coredb

import (
	"testing"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/zhs007/ankadb"
	"github.com/zhs007/jarviscore/crypto"
)

func TestBaseFunc(t *testing.T) {
	//------------------------------------------------------------------------
	// initial ankaDB

	cfg := ankadb.NewConfig()

	cfg.PathDBRoot = "../test/coredbbasefunc"
	cfg.ListDB = append(cfg.ListDB, ankadb.DBConfig{
		Name:   "coredb",
		Engine: "leveldb",
		PathDB: "coredb",
	})

	dblogic, err := ankadb.NewBaseDBLogic(graphql.SchemaConfig{
		Query:    typeQuery,
		Mutation: typeMutation,
	})
	if err != nil {
		t.Fatalf("TestBaseFunc NewBaseDBLogic %v", err)

		return
	}

	ankaDB, err := ankadb.NewAnkaDB(cfg, dblogic)
	if ankaDB == nil {
		t.Fatalf("TestBaseFunc NewAnkaDB %v", err)

		return
	}

	//------------------------------------------------------------------------
	// newPrivateData

	privKey := jarviscrypto.GenerateKey()

	pd, err := newPrivateData(ankaDB, jarviscrypto.Base58Encode(privKey.ToPrivateBytes()), jarviscrypto.Base58Encode(privKey.ToPublicBytes()), privKey.ToAddress(), time.Now().Unix())
	if err != nil {
		t.Fatalf("TestBaseFunc newPrivateData %v", err)

		return
	}

	if pd.PriKey != nil || pd.PubKey != nil || pd.StrPriKey != "" {
		t.Fatalf("TestBaseFunc getPrivateKey return PriKey or PubKey or StrPriKey is non-nil")

		return
	}

	if pd.StrPubKey != jarviscrypto.Base58Encode(privKey.ToPublicBytes()) {
		t.Fatalf("TestBaseFunc newPrivateData return StrPubKey fail! (%v - %v)", pd.StrPubKey, jarviscrypto.Base58Encode(privKey.ToPublicBytes()))

		return
	}

	//------------------------------------------------------------------------
	// getPrivateKey

	pd, err = getPrivateKey(ankaDB)
	if err != nil {
		t.Fatalf("TestBaseFunc getPrivateKey %v", err)

		return
	}

	if pd.PriKey != nil || pd.PubKey != nil {
		t.Fatalf("TestBaseFunc getPrivateKey return PriKey or PubKey is non-nil")

		return
	}

	if pd.StrPubKey != jarviscrypto.Base58Encode(privKey.ToPublicBytes()) {
		t.Fatalf("TestBaseFunc getPrivateKey return StrPubKey fail! (%v - %v)", pd.StrPubKey, jarviscrypto.Base58Encode(privKey.ToPublicBytes()))

		return
	}

	if pd.StrPriKey != jarviscrypto.Base58Encode(privKey.ToPrivateBytes()) {
		t.Fatalf("TestBaseFunc getPrivateKey return StrPriKey fail! (%v - %v)", pd.StrPriKey, jarviscrypto.Base58Encode(privKey.ToPublicBytes()))

		return
	}

	//------------------------------------------------------------------------
	// getPrivateData

	pd, err = getPrivateData(ankaDB)
	if err != nil {
		t.Fatalf("TestBaseFunc getPrivateData %v", err)

		return
	}

	if pd.PriKey != nil || pd.PubKey != nil || pd.StrPriKey != "" {
		t.Fatalf("TestBaseFunc getPrivateData return PriKey or PubKey or StrPriKey is non-nil")

		return
	}

	if pd.StrPubKey != jarviscrypto.Base58Encode(privKey.ToPublicBytes()) {
		t.Fatalf("TestBaseFunc getPrivateData return StrPubKey fail! (%v - %v)", pd.StrPubKey, jarviscrypto.Base58Encode(privKey.ToPublicBytes()))

		return
	}

	//------------------------------------------------------------------------
	// updPrivateData

	pd, err = updPrivateData(ankaDB, int64(100))
	if err != nil {
		t.Fatalf("TestBaseFunc updPrivateData %v", err)

		return
	}

	if pd.PriKey != nil || pd.PubKey != nil || pd.StrPriKey != "" {
		t.Fatalf("TestBaseFunc updPrivateData return PriKey or PubKey or StrPriKey is non-nil")

		return
	}

	if pd.StrPubKey != jarviscrypto.Base58Encode(privKey.ToPublicBytes()) {
		t.Fatalf("TestBaseFunc updPrivateData return StrPubKey fail! (%v - %v)", pd.StrPubKey, jarviscrypto.Base58Encode(privKey.ToPublicBytes()))

		return
	}

	t.Logf("TestBaseFunc is OK")
}
