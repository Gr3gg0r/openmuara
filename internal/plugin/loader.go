package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadedPlugin binds a parsed GatewayConfig to its source directory.
type LoadedPlugin struct {
	Name   string
	Dir    string
	Config GatewayConfig
}

var pluginNameRe = regexp.MustCompile(`^[a-z0-9-]+$`)

// LoadBuiltin discovers built-in plugins under the given directories.
func LoadBuiltin(dirs ...string) ([]*LoadedPlugin, error) {
	if len(dirs) == 0 {
		dirs = []string{"plugins"}
	}
	return loadFromDirs(dirs)
}

// LoadLocal discovers user-local plugin overrides.
func LoadLocal(dir string) ([]*LoadedPlugin, error) {
	if dir == "" {
		dir = ".muara/plugins"
	}
	return loadFromDirs([]string{dir})
}

func loadFromDirs(dirs []string) ([]*LoadedPlugin, error) {
	var out []*LoadedPlugin
	for _, d := range dirs {
		entries, err := os.ReadDir(d)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("read plugin dir %q: %w", d, err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			name := entry.Name()
			if !pluginNameRe.MatchString(name) {
				return nil, fmt.Errorf("invalid plugin directory name %q", name)
			}
			pdir := filepath.Join(d, name)
			if err := assertWithinRoot(d, pdir); err != nil {
				return nil, err
			}
			cfg, err := loadGatewayYAML(pdir)
			if err != nil {
				return nil, fmt.Errorf("load plugin %q: %w", name, err)
			}
			out = append(out, &LoadedPlugin{Name: name, Dir: pdir, Config: cfg})
		}
	}
	return out, nil
}

func loadGatewayYAML(dir string) (GatewayConfig, error) {
	path := filepath.Join(dir, "gateway.yml")
	// #nosec G304 -- intentional plugin manifest load from validated directory
	data, err := os.ReadFile(path)
	if err != nil {
		return GatewayConfig{}, fmt.Errorf("read %q: %w", path, err)
	}
	var cfg GatewayConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return GatewayConfig{}, fmt.Errorf("parse %q: %w", path, err)
	}
	return cfg, nil
}

func assertWithinRoot(root, target string) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve root %q: %w", root, err)
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("resolve target %q: %w", target, err)
	}
	prefix := absRoot + string(filepath.Separator)
	if absRoot == absTarget || strings.HasPrefix(absTarget, prefix) {
		return nil
	}
	return fmt.Errorf("plugin directory %q is outside root %q", target, root)
}
