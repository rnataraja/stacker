package stacker

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/openSUSE/umoci"
)

type BaseLayerOpts struct {
	Config StackerConfig
	Name   string
	Target string
	Layer  *Layer
	Cache  *BuildCache
	OCI    *umoci.Layout
}

func GetBaseLayer(o BaseLayerOpts) error {
	switch o.Layer.From.Type {
	case BuiltType:
		/* nothing to do assuming layers are imported in dependency order */
		return nil
	case TarType:
		return getTar(o)
	case OCIType:
		return fmt.Errorf("not implemented")
	case DockerType:
		return getDocker(o)
	default:
		return fmt.Errorf("unknown layer type: %v", o.Layer.From.Type)
	}
}

func getDocker(o BaseLayerOpts) error {
	out, err := exec.Command(
		"skopeo",
		"copy",
		fmt.Sprintf("oci:%s:%s", o.Config.OCIDir, o.Name, o.Layer.From.Url),
		o.Layer.From.Url).CombinedOutput()
	if err != nil {
		return fmt.Errorf("skopeo copy: %s: %s", err, string(out))
	}

	return o.OCI.Unpack(o.Name, path.Join(o.Config.RootFSDir, o.Name), nil)
}

func getTar(o BaseLayerOpts) error {
	tar, err := download(path.Join(o.Config.StackerDir, "layer-bases"), o.Layer.From.Url)
	if err != nil {
		return err
	}

	layerPath := path.Join(o.Config.RootFSDir, o.Target)
	if err := os.MkdirAll(layerPath, 0755); err != nil {
		return err
	}

	output, err := exec.Command("tar", "xf", tar, "-C", layerPath).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error: %s: %s", err, string(output))
	}

	return nil
}
