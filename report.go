package main

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type PluginReport struct {
	Name       string
	Status     string
	ActiveInst int
	Value      string
}

func main() {
	var reports []PluginReport
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".zip" {
			report := analyzePlugin(path)
			reports = append(reports, report)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error while walking through files:", err)
		return
	}

	reportFile, err := os.Create("plugin_report.txt")
	if err != nil {
		fmt.Println("Error while creating report.txt:", err)
		return
	}
	defer reportFile.Close()

	for _, report := range reports {
		fmt.Fprintf(reportFile, "Name: %s\nStatus: %s\nActive Installs: %d\nValue: %s\n\n", report.Name, report.Status, report.ActiveInst, report.Value)
	}
}

func analyzePlugin(filePath string) PluginReport {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return PluginReport{Name: filePath, Status: "Invalid (could not open)", ActiveInst: 0, Value: "N/A"}
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "readme.txt") {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			content, err := ioutil.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}

			return parseReadme(filePath, string(content))
		}
	}

	return PluginReport{Name: filePath, Status: "No Readme", ActiveInst: 0, Value: "N/A"}
}

func parseReadme(fileName, content string) PluginReport {
	lines := strings.Split(content, "\n")
	activeInstalls := 0
	status := "Accepted"
	value := "Low"

	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "active installations") {
			activeInstalls = parseActiveInstalls(line)
		}
	}

	if activeInstalls > 5000000 {
		value = "High"
	} else if activeInstalls > 1000000 {
		value = "Medium"
	}

	return PluginReport{
		Name:       strings.TrimSuffix(fileName, ".zip"),
		Status:     status,
		ActiveInst: activeInstalls,
		Value:      value,
	}
}

func parseActiveInstalls(line string) int {
	line = strings.ToLower(line)
	if strings.Contains(line, "million") {
		parts := strings.Split(line, "million")
		if len(parts) > 0 {
			value := strings.TrimSpace(parts[0])
			if value == "5" {
				return 5000000
			}
		}
	} else if strings.Contains(line, "thousand") {
		parts := strings.Split(line, "thousand")
		if len(parts) > 0 {
			value := strings.TrimSpace(parts[0])
			if value == "500" {
				return 500000
			}
		}
	}
	return 0
}
