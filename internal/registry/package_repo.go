package registry

type PackageVersion struct {
	Name         string
	Version      string
	Checksum     string
	Tarball      string
	Published    string
	ManifestXML  string
	MetadataJSON string
}
