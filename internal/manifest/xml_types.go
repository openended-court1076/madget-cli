package manifest

import (
	"encoding/json"
	"encoding/xml"
)

// Application MadGet.xml kök belgesi (registry metadatası ve CLI info için).
type Application struct {
	XMLName      xml.Name     `xml:"application" json:"-"`
	Info         Info         `xml:"info" json:"info"`
	Protocols    Protocols    `xml:"protocols" json:"protocols,omitempty"`
	FilesHandler FilesHandler `xml:"files_handler" json:"files_handler,omitempty"`
}

type Info struct {
	XMLName     xml.Name `xml:"info" json:"-"`
	Name        string   `xml:"name,attr" json:"name,omitempty"`
	Version     string   `xml:"version,attr" json:"version,omitempty"`
	PackageName string   `xml:"package_name,attr" json:"package_name,omitempty"`
	License     string   `xml:"license,attr" json:"license,omitempty"`
	Categories  string   `xml:"categories,attr" json:"categories,omitempty"`
	Description string   `xml:"description" json:"description,omitempty"`
	Author      Author   `xml:"author" json:"author,omitempty"`
	Readme      string   `xml:"readme" json:"readme,omitempty"`
	Homepage    string   `xml:"homepage" json:"homepage,omitempty"`
	Permissions []string `xml:"permissions>permission" json:"permissions,omitempty"`
}

type Author struct {
	XMLName xml.Name `xml:"author" json:"-"`
	ID      string   `xml:"id,attr" json:"id,omitempty"`
}

type Protocols struct {
	XMLName  xml.Name   `xml:"protocols" json:"-"`
	Protocol []Protocol `xml:"protocol" json:"protocol,omitempty"`
}

type Protocol struct {
	XMLName xml.Name `xml:"protocol" json:"-"`
	Scheme  string   `xml:"schema,attr" json:"schema,omitempty"`
	Handler string   `xml:"handler,attr" json:"handler,omitempty"`
}

type FilesHandler struct {
	XMLName     xml.Name      `xml:"files_handler" json:"-"`
	FileHandler []FileHandler `xml:"file_handler" json:"file_handler,omitempty"`
}

type FileHandler struct {
	XMLName xml.Name `xml:"file_handler" json:"-"`
	Ext     string   `xml:"ext,attr" json:"ext,omitempty"`
	Handler string   `xml:"handler,attr" json:"handler,omitempty"`
}

// UnmarshalApplication MadGet.xml baytlarını ayrıştırır.
func UnmarshalApplication(data []byte) (Application, error) {
	var app Application
	if err := xml.Unmarshal(data, &app); err != nil {
		return app, err
	}
	return app, nil
}

// MetadataJSON registry’de saklanacak yapılandırılmış metadatayı üretir (JSON).
func (a Application) MetadataJSON() ([]byte, error) {
	return json.Marshal(a)
}
