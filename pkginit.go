package jarviscore

var mgrConn *connMgr

// var mgrCtrl *ctrlMgr

func init() {
	mgrConn = &connMgr{}

	// mgrCtrl = &ctrlMgr{mapCtrl: make(map[string](Ctrl))}
	// mgrCtrl.Reg(CtrlTypeShell, &CtrlShell{})
	// mgrCtrl.Reg(CtrlTypeScriptFile, &CtrlScriptFile{})
	// mgrCtrl.Reg(CtrlTypeScriptFile2, &CtrlScriptFile2{})
}
