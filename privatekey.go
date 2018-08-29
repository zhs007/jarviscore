package jarviscore

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// privatekey - privatekey
type privatekey struct {
	Token string `yaml:"token"`
}

const privateKeyFilename = "prikey.yaml"

func loadPrivateKeyFile() (*privatekey, error) {
	buf, err := loadFile(getRealPath(privateKeyFilename))
	if err != nil {
		return nil, err
	}

	prikey := privatekey{}
	err1 := yaml.Unmarshal(buf, &prikey)
	if err1 != nil {
		return nil, err1
	}

	return &prikey, nil
}

func savePrivateKeyFile(prikey *privatekey) error {
	d, err := yaml.Marshal(prikey)
	if err != nil {
		return err
	}

	ioutil.WriteFile(getRealPath(privateKeyFilename), d, 0755)

	return nil
}
