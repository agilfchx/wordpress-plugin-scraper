package main

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io"
        "net/http"
        "os"
        "path/filepath"
        "sync"
        "syscall"
        "time"
)

const (
        baseURL       = "https://api.wordpress.org/plugins/info/1.2/?action=query_plugins&request[page]=%d"
        minInstalls   = 1000
        maxInstalls   = 100000
        maxWorkers    = 5
        telegramToken = "TELEGRAM_TOKEN"
        chatID        = "TELEGRAM_CHAT_ID"

        intervalSend = 3 * time.Hour
        //intervalSend = 10 * time.Second // TEST ONLY
)

type Plugin struct {
        Slug           string `json:"slug"`
        Version        string `json:"version"`
        DownloadLink   string `json:"download_link"`
        ActiveInstalls int    `json:"active_installs"`
}

type PluginList struct {
        Plugins []Plugin `json:"plugins"`
}

func createDownloadedFolder() error {
        if _, err := os.Stat("Downloaded"); os.IsNotExist(err) {
                err := os.MkdirAll("Downloaded", os.ModePerm)
                if err != nil {
                        return fmt.Errorf("failed to create 'Downloaded' folder: %v", err)
                }
        }
        return nil
}

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

func getTotalStorage() (float64, error) {
        var stat syscall.Statfs_t
        wd, err := os.Getwd()
        if err != nil {
                return 0, err
        }

        if err := syscall.Statfs(wd, &stat); err != nil {
                return 0, err
        }

        total := float64(stat.Blocks) * float64(stat.Bsize) / (1024 * 1024 * 1024)
        return total, nil
}

func getFreeStorage() (float64, error) {
        var stat syscall.Statfs_t
        wd, err := os.Getwd()
        if err != nil {
                return 0, err
        }

        if err := syscall.Statfs(wd, &stat); err != nil {
                return 0, err
        }

        available := float64(stat.Bavail) * float64(stat.Bsize) / (1024 * 1024 * 1024)
        return available, nil
}

func getDownloadedFolderSize() (float64, error) {
        var totalSize float64

        downloadedFolder := "./Downloaded"
        err := filepath.Walk(downloadedFolder, func(path string, info os.FileInfo, err error) error {
                if err != nil {
                        return err
                }

                if !info.IsDir() {
                        totalSize += float64(info.Size()) / (1024 * 1024 * 1024)
                }
                return nil
        })

        if err != nil {
                return 0, err
        }

        return totalSize, nil
}


func sendNotification(pageNumber int, downloadCounter int) {
        freeStorage, err := getFreeStorage()
        if err != nil {
                fmt.Println("Error calculating free storage:", err)
                freeStorage = 0
        }

        downloadedStorage, err := getDownloadedFolderSize()
        if err != nil {
                fmt.Println("Error calculating downloaded folder size:", err)
                downloadedStorage = 0
        }


        message := fmt.Sprintf(
                "📢 Plugin Download Report\n"+
                        "- - - - - - - - - - - - - - - - - - - - - - - - - - -\n"+
                        "📄 Current Page: %d\n"+
                        "📥 Total Downloads: %d\n"+
                        "📊 Actv. Installation Range: %d - %d\n"+
                        "💾 Storage Usage: %.2f GB / %.2f GB",
                pageNumber, downloadCounter, minInstalls, maxInstalls, downloadedStorage, freeStorage)

        url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramToken)
        payload := map[string]string{
                "chat_id": chatID,
                "text":    message,
        }
        jsonData, _ := json.Marshal(payload)
        http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}


func downloadPlugin(plugin Plugin, wg *sync.WaitGroup, downloadCounter *int) {
        defer wg.Done()

        err := createDownloadedFolder()
        if err != nil {
                fmt.Println(err)
                return
        }

        resp, err := http.Get(plugin.DownloadLink)
        if err != nil {
                fmt.Printf("Failed to download %s: %v\n", plugin.Slug, err)
                return
        }
        defer resp.Body.Close()

        fileName := fmt.Sprintf("Downloaded/%s-%s.zip", plugin.Slug, plugin.Version)
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
        *downloadCounter++
        time.Sleep(1 * time.Second)
}

func startNotifier(downloadCounter *int, pageNumber *int) {
        for {
                time.Sleep(intervalSend)
                sendNotification(*pageNumber, *downloadCounter)
        }
}

func main() {
        pageNumber := 1
        var wg sync.WaitGroup
        jobs := make(chan Plugin, maxWorkers)
        downloadCounter := 0

        go startNotifier(&downloadCounter, &pageNumber)

        for i := 0; i < maxWorkers; i++ {
                go func() {
                        for plugin := range jobs {
                                downloadPlugin(plugin, &wg, &downloadCounter)
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
                        if plugin.ActiveInstalls >= minInstalls && plugin.ActiveInstalls <= maxInstalls {
                                wg.Add(1)
                                jobs <- plugin
                        }
                }
                pageNumber++
        }

        wg.Wait()
        close(jobs)
}