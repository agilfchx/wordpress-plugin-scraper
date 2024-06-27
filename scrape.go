package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	baseURL      = "https://api.wordpress.org/plugins/info/1.2/?action=query_plugins&request[page]=%d"
	minInstalls  = 1000
	maxWorkers   = 5
	requestRetry = 3
)

// Plugin structure to hold plugin details from the API
type Plugin struct {
	Slug          string `json:"slug"`
	Version       string `json:"version"`
	DownloadLink  string `json:"download_link"`
	ActiveInstalls int   `json:"active_installs"`
}

// PluginList structure to hold the list of plugins fetched from the API
type PluginList struct {
	Plugins []Plugin `json:"plugins"`
}

// fetchPluginList fetches the list of plugins from the WordPress API
func fetchPluginList(pageNumber int) (PluginList, error) {
	var pluginList PluginList

	url := fmt.Sprintf(baseURL, pageNumber)
	resp, err := http.Get(url)
	if err != nil {
		return pluginList, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return pluginList, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&pluginList)
	return pluginList, err
}

// downloadPlugin handles the downloading of the plugin file
func downloadPlugin(plugin Plugin, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(plugin.DownloadLink)
	if err != nil {
		fmt.Printf("Failed to download %s: %v\n", plugin.Slug, err)
		return
	}
	defer resp.Body.Close()

	fileName := fmt.Sprintf("%s-%s.zip", plugin.Slug, plugin.Version)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Failed to write file %s: %v\n", fileName, err)
		return
	}

	fmt.Printf("Downloaded %s version %s\n", plugin.Slug, plugin.Version)
	time.Sleep(1 * time.Second)
}

func main() {
	pageNumber := 1
	var wg sync.WaitGroup
	jobs := make(chan Plugin, maxWorkers)

	// Worker pool to limit the concurrent downloads
	for i := 0; i < maxWorkers; i++ {
		go func() {
			for plugin := range jobs {
				downloadPlugin(plugin, &wg)
			}
		}()
	}

	for {
		pluginList, err := fetchPluginList(pageNumber)
		if err != nil {
			fmt.Printf("Failed to fetch plugin list: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if len(pluginList.Plugins) == 0 {
			break
		}

		for _, plugin := range pluginList.Plugins {
			if plugin.ActiveInstalls >= minInstalls {
				wg.Add(1)
				jobs <- plugin
			}
		}
		pageNumber++
	}

	wg.Wait()
	close(jobs)
}
