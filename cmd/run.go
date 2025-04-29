package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/scenery/mediax/config"
	"github.com/scenery/mediax/database"
	"github.com/scenery/mediax/dataops"
	"github.com/scenery/mediax/routes"
	"github.com/scenery/mediax/version"
)

func Execute() {
	configPath := flag.String("config", "config.json", "指定配置文件路径")
	port := flag.Int("port", config.DefaultConfig.Server.Port, "指定启动端口")
	versionFlag := flag.Bool("version", false, "显示 mediaX 版本信息")
	importType := flag.String("import", "", "导入数据来源: bangumi 或 douban")
	filePath := flag.String("file", "", "导入文件的相对路径")
	downloadImage := flag.Bool("download-image", false, "可选，导入数据时是否下载图片")
	exportType := flag.String("export", "", "导出数据: all, book, movie, tv, anime, game")
	limit := flag.Int("limit", -1, "可选，指定导出的数量，默认导出所有数据")

	flag.Usage = func() {
		fmt.Println("Document: https://github.com/scenery/mediax/README.md")
		fmt.Println("\nAvailable parameters:")
		flag.PrintDefaults()
	}

	flag.Parse()

	usedFlags := map[string]bool{}
	flag.Visit(func(f *flag.Flag) {
		usedFlags[f.Name] = true
	})

	if *versionFlag {
		if len(usedFlags) > 1 {
			fmt.Println("Error: --version flag cannot be used with other parameters")
			os.Exit(1)
		}
		fmt.Println("mediaX", version.Version)
		os.Exit(0)
	}

	if usedFlags["config"] {
		if len(usedFlags) > 1 {
			fmt.Println("Error: Error: --config flag cannot be used with other parameters")
			os.Exit(1)
		}
		err := config.LoadConfig(*configPath)
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		startServer(config.App.Server.Port)
		return
	}

	config.App = config.DefaultConfig

	if *importType != "" {
		if *filePath == "" {
			fmt.Println("Error: File path (relative path) is required for data import")
			flag.Usage()
			os.Exit(1)
		}
		database.InitDB()
		err := dataops.ImportFromJSON(*importType, *filePath, *downloadImage)
		if err != nil {
			fmt.Println("Error during import:", err)
			os.Exit(1)
		}
		return
	}
	if *downloadImage && *importType == "" {
		fmt.Println("Error: --download-image is only supported during import")
		flag.Usage()
		os.Exit(1)
	}

	if *exportType != "" {
		database.InitDB()
		err := dataops.ExportToJSON(*exportType, *limit)
		if err != nil {
			fmt.Println("Error during export:", err)
			os.Exit(1)
		}
		return
	}

	startServer(*port)
}

func startServer(port int) {
	database.InitDB()
	routes.Init()

	address := fmt.Sprintf("%s:%d", config.App.Server.Address, port)
	fmt.Print(`
                    _ _      __  __
 _ __ ___   ___  __| (_) __ _\ \/ /
| '_   _ \ / _ \/ _  | |/ _  |\  / 
| | | | | |  __/ (_| | | (_| |/  \ 
|_| |_| |_|\___|\__,_|_|\__,_/_/\_\

`)
	fmt.Printf("mediaX %s is running at: http://%s\n", version.Version, address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("mediaX failed to start: %v", err)
	}
}
