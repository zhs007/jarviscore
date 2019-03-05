package jarviscore

import (
	"bytes"
	"os/exec"
	"text/template"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/base"
)

// UpdateNodeParam - the parameter for update node
type UpdateNodeParam struct {
	NewVersion string
}

// updateNode - update node
func updateNode(params *UpdateNodeParam, scriptUpd string) (string, string, error) {
	tpl, err := template.New("autoupd").Parse(scriptUpd)
	if err != nil {
		return "", "", err
	}

	var b bytes.Buffer
	tpl.Execute(&b, params)

	jarvisbase.Info("updateNode script", zap.String("script", b.String()))

	out, err := exec.Command("sh", "-c", b.String()).CombinedOutput()
	if err != nil {
		return "", "", err
	}

	return b.String(), string(out), nil
}
