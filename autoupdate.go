package jarviscore

import (
	"bytes"
	"os/exec"
	"text/template"

	jarvisbase "github.com/zhs007/jarviscore/base"
	"go.uber.org/zap"
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
		return "", string(out), err
	}

	return b.String(), string(out), nil
}

// RestartNodeParam - the parameter for restart node
type RestartNodeParam struct {
}

// restartNode - restart node
func restartNode(params *RestartNodeParam, scriptUpd string) (string, string, error) {
	tpl, err := template.New("autoupd.restart").Parse(scriptUpd)
	if err != nil {
		return "", "", err
	}

	var b bytes.Buffer
	tpl.Execute(&b, params)

	jarvisbase.Info("restartNode script", zap.String("script", b.String()))

	out, err := exec.Command("sh", "-c", b.String()).CombinedOutput()
	if err != nil {
		return "", string(out), err
	}

	return b.String(), string(out), nil
}
