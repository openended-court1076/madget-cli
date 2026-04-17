package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mehmetalidsy/madget-cli/internal/manifest"
	"github.com/spf13/cobra"
)

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

func showProtocols(protocols manifest.Protocols, fgColor func(a ...interface{})) {
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

func iniGet(path string) manifest.Application {
	xmlFile, err := os.ReadFile(path)
	if err != nil {
		errorMessage("Error reading XML file: " + err.Error())
		return manifest.Application{}
	}

	appManifest, err := manifest.UnmarshalApplication(xmlFile)
	if err != nil {
		errorMessage("Error parsing XML file: " + err.Error())
		return manifest.Application{}
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
