package manifest

import (
	"fmt"
	"os"
	"strings"
)

// PackageManifest, MadGet.xml içindeki <application><info> alanlarından türetilir.
// Registry’de paket adı: önce package_name, yoksa name.
type PackageManifest struct {
	Name        string
	PackageName string
	Version     string
	Description string
}

// RegistryName publish ve resolve için kullanılan benzersiz paket adı.
func (m PackageManifest) RegistryName() string {
	if s := strings.TrimSpace(m.PackageName); s != "" {
		return s
	}
	return strings.TrimSpace(m.Name)
}

// Load MadGet.xml dosyasını okur (kök öğe <application>).
func Load(path string) (PackageManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PackageManifest{}, err
	}
	return LoadFromBytes(data)
}

// LoadFromBytes MadGet.xml içeriğinden PackageManifest üretir.
func LoadFromBytes(data []byte) (PackageManifest, error) {
	var m PackageManifest
	app, err := UnmarshalApplication(data)
	if err != nil {
		return m, fmt.Errorf("MadGet.xml parse: %w", err)
	}
	info := app.Info
	m.Name = strings.TrimSpace(info.Name)
	m.PackageName = strings.TrimSpace(info.PackageName)
	m.Version = strings.TrimSpace(info.Version)
	m.Description = strings.TrimSpace(info.Description)
	if m.RegistryName() == "" {
		return m, fmt.Errorf("MadGet.xml: name veya package_name zorunlu")
	}
	if m.Version == "" {
		return m, fmt.Errorf("MadGet.xml: version zorunlu")
	}
	return m, nil
}

