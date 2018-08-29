package jarviscore

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// CtrlInfo -
type CtrlInfo struct {
	Command string `yaml:"command"`
	Result  string `yaml:"result"`
}

// nodeCtrlInfo -
type nodeCtrlInfo struct {
	Mapctrl map[int32]*CtrlInfo `yaml:"mapctrl"`
}

func newNodeCtrlInfo() *nodeCtrlInfo {
	return &nodeCtrlInfo{
		Mapctrl: make(map[int32]*CtrlInfo),
	}
}

func loadNodeCtrlInfo(filename string) (*nodeCtrlInfo, error) {
	buf, err := loadFile(filename)
	if err != nil {
		return nil, err
	}

	nodectrlinfo := newNodeCtrlInfo()
	err1 := yaml.Unmarshal(buf, nodectrlinfo)
	if err1 != nil {
		return nil, err1
	}

	return nodectrlinfo, nil
}

func (m *nodeCtrlInfo) save(filename string) error {
	d, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	ioutil.WriteFile(filename, d, 0755)

	return nil
}

func (m *nodeCtrlInfo) clear() {
	if len(m.Mapctrl) == 0 {
		return
	}

	for k := range m.Mapctrl {
		delete(m.Mapctrl, k)
	}

	m.Mapctrl = nil
	m.Mapctrl = make(map[int32]*CtrlInfo)
}

func (m *nodeCtrlInfo) hasCtrl(ctrlid int32) bool {
	if _, ok := m.Mapctrl[ctrlid]; ok {
		return true
	}

	return false
}

func (m *nodeCtrlInfo) addCtrl(ctrlid int32, command string) {
	m.Mapctrl[ctrlid] = &CtrlInfo{Command: command}
}

func (m *nodeCtrlInfo) setCtrlResult(ctrlid int32, result string) {
	m.Mapctrl[ctrlid].Result = result
}
