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

// mapCtrlInfo -
type mapCtrlInfo struct {
	Mapctrl map[int32]*CtrlInfo `yaml:"mapctrl"`
}

func newMapCtrlInfo() *mapCtrlInfo {
	return &mapCtrlInfo{
		Mapctrl: make(map[int32]*CtrlInfo),
	}
}

func loadMapCtrlInfo(filename string) (*mapCtrlInfo, error) {
	buf, err := loadFile(filename)
	if err != nil {
		return nil, err
	}

	mapctrlinfo := newMapCtrlInfo()
	err1 := yaml.Unmarshal(buf, mapctrlinfo)
	if err1 != nil {
		return nil, err1
	}

	return mapctrlinfo, nil
}

func (m *mapCtrlInfo) save(filename string) error {
	d, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	ioutil.WriteFile(filename, d, 0755)

	return nil
}

func (m *mapCtrlInfo) clear() {
	if len(m.Mapctrl) == 0 {
		return
	}

	for k := range m.Mapctrl {
		delete(m.Mapctrl, k)
	}

	m.Mapctrl = nil
	m.Mapctrl = make(map[int32]*CtrlInfo)
}

func (m *mapCtrlInfo) hasCtrl(ctrlid int32) bool {
	if _, ok := m.Mapctrl[ctrlid]; ok {
		return true
	}

	return false
}

func (m *mapCtrlInfo) addCtrl(ctrlid int32, command string) {
	m.Mapctrl[ctrlid] = &CtrlInfo{Command: command}
}

func (m *mapCtrlInfo) setCtrlResult(ctrlid int32, result string) {
	m.Mapctrl[ctrlid].Result = result
}
