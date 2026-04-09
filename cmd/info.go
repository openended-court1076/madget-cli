package cmd

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Application struct {
	XMLName      xml.Name     `xml:"application"`
	Info         Info         `xml:"info"`
	Protocols    Protocols    `xml:"protocols"`
	FilesHandler FilesHandler `xml:"files_handler"`
}

type Info struct {
	XMLName     xml.Name `xml:"info"`
	Name        string   `xml:"name,attr"`
	Version     string   `xml:"version,attr"`
	PackageName string   `xml:"package_name,attr"`
	License     string   `xml:"license,attr"`
	Categories  string   `xml:"categories,attr"`

	Description string   `xml:"description"`
	Author      Author   `xml:"author"`
	Readme      string   `xml:"readme"`
	Homepage    string   `xml:"homepage"`
	Permissions []string `xml:"permissions>permission"`
}

type Author struct {
	XMLName xml.Name `xml:"author"`
	ID      string   `xml:"id,attr"`
}

type Protocols struct {
	XMLName  xml.Name   `xml:"protocols"`
	Protocol []Protocol `xml:"protocol"`
}

type Protocol struct {
	XMLName xml.Name `xml:"protocol"`
	Scheme  string   `xml:"schema,attr"`
	Handler string   `xml:"handler,attr"`
}

type FilesHandler struct {
	XMLName     xml.Name      `xml:"files_handler"`
	FileHandler []FileHandler `xml:"file_handler"`
}

type FileHandler struct {
	XMLName xml.Name `xml:"file_handler"`
	Ext     string   `xml:"ext,attr"`
	Handler string   `xml:"handler,attr"`
}

var infoCmd = &cobra.Command{
	Use:   "info [package]",
	Short: "MadGet Package Information",
	Long:  `MadGet provides information about Mad packages.`,
	Run: func(cmd *cobra.Command, args []string) {
		wdPath, wdErr := os.Getwd()
		if wdErr != nil {
			fmt.Println("Error getting current working directory:", wdErr)
			return
		}

		_x, fileErr := os.Stat(wdPath + "/MadGet.xml")

		if fileErr != nil {
			errorMessage("Error: MadGet.xml file not found in the current directory. Please run this command in a valid Mad package directory.")
			return
		} else {
			PackageData := iniGet(wdPath + "/" + _x.Name())

			color.New(color.FgHiGreen).Println("Current working directory: ", wdPath)
			color.New(color.FgHiGreen).Println("Current read MadGet.xml file path: ", wdPath+"/MadGet.xml")
			fmt.Println()

			fgColor := color.New(color.FgHiBlue).PrintlnFunc()
			packageName := strings.Replace(PackageData.Info.PackageName, " ", "", -1)

			fgColor("Package Name: ", clean(PackageData.Info.Name), "("+packageName+")")
			fgColor("Version: ", clean(PackageData.Info.Version))
			fgColor("Description:", clean(PackageData.Info.Description))
			fgColor("Categories: ", clean(PackageData.Info.Categories))
			fgColor("Permissions:", strings.Join(PackageData.Info.Permissions, ", "))
			fgColor("Author ID Or Name: ", clean(PackageData.Info.Author.ID))
			fgColor("Homepage: ", clean(PackageData.Info.Homepage))
			color.New(color.FgHiYellow).Println("License: ", clean(PackageData.Info.License))

			fmt.Println()
			showProtocols(PackageData.Protocols, fgColor)
			fmt.Println()
			readmeContentView(wdPath)
			fmt.Println()
		}
	},
}

func errorMessage(text string, a ...any) {
	color.New(color.FgHiRed).Println(text, a)
	fmt.Println()
}

func showProtocols(protocols Protocols, fgColor func(a ...interface{})) {
	color.New(color.FgHiCyan).Println("------- PROTOCOLS -------")
	for _, p := range protocols.Protocol {
		fgColor("Schema: ", clean(p.Scheme), " - Handler: ", clean(p.Handler))
	}
	color.New(color.FgHiCyan).Println("------- PROTOCOLS -------")
}

func readmeContentView(wdPath string) {
	color.New(color.FgHiMagenta).Println("------- README -------")
	readmeFile := "README.md"

	readmeContent, readmeErr := readmeContent(wdPath + "/" + readmeFile)
	if readmeErr != nil {
		errorMessage("Error reading README file: " + readmeErr.Error())
		return
	}

	fmt.Println()
	fmt.Println(string(readmeContent))
	fmt.Println()
	color.New(color.FgHiMagenta).Println("------- README -------")
}

func iniGet(path string) Application {
	xmlFile, err := os.ReadFile(path)
	if err != nil {
		errorMessage("Error reading XML file: " + err.Error())
	}

	var appManifest Application
	err = xml.Unmarshal(xmlFile, &appManifest)
	if err != nil {
		errorMessage("Error parsing XML file: " + err.Error())
	}

	return appManifest
}

func readmeContent(path string) (string, error) {
	readmePath := path
	readmeContent, readmeErr := os.ReadFile(readmePath)
	if readmeErr != nil {
		errorMessage("Error reading README file: " + readmeErr.Error())
		return "", readmeErr
	}

	return string(readmeContent), nil
}

func clean(s string) string {
	return strings.Trim(s, "\" ")
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
