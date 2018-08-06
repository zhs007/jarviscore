package jarviscore

import (
	"gopkg.in/yaml.v2"
)

// peeraddrinfo -
type peeraddrinfo struct {
	PeerAddr []string `yaml:"peeraddr"`
}

func loadPeerAddr(filename string) (*peeraddrinfo, error) {
	buf, err := loadFile(filename)
	if err != nil {
		return nil, err
	}

	pa := peeraddrinfo{}
	err1 := yaml.Unmarshal(buf, &pa)
	if err1 != nil {
		return nil, err1
	}

	return &pa, nil
}
