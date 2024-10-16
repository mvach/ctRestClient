package main

import (
    "ctRestClient/app"
    "ctRestClient/config"
    "flag"
    "log"
    "os"
    "path/filepath"
)

func main() {
    var configPath string
    var outputDirectory string
    
    flag.StringVar(&configPath, "c", "config.yml", "the config file path")
    flag.StringVar(&outputDirectory, "o", getExecutablePath(), "the output directory")
    flag.Parse()
    

    config, err := config.LoadConfig(configPath)
    if err != nil {
        log.Fatalf("Failed to load config from path %s: %v", configPath, err)
    }

    

    err = app.NewInstancesProcessor(*config, outputDirectory, app.NewLogger()).Process(
        app.NewGroupExporter(),
        app.NewCSVFileWriter(),
    )
    if err != nil {
        log.Fatalf("Failed to process instances: %v", err)
    }
}

func getExecutablePath() string {
    exePath, err := os.Executable()
    if err != nil {
        log.Fatalf("Failed to get executable path: %v", err)
    }
    return filepath.Dir(exePath)   
}