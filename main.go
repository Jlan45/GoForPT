package main

import (
	"GoForPT/api"
	"GoForPT/pkg/cfg"
	"GoForPT/pkg/database"
	"GoForPT/pkg/ptcaches"
	"log"
	"os"
)

func main() {
	err := cfg.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 检查种子文件和静态文件夹是否存在，不存在则创建
	if _, err := os.Stat(cfg.Cfg.Site.TorrentPath); os.IsNotExist(err) {
		err = os.Mkdir(cfg.Cfg.Site.TorrentPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create torrent folder: %v", err)
			return
		}
	}
	if _, err := os.Stat(cfg.Cfg.Site.StaticPath); os.IsNotExist(err) {
		err = os.Mkdir(cfg.Cfg.Site.StaticPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create static folder: %v", err)
			return
		}
	}

	ptcaches.InitEmailCache()
	ptcaches.InitAnnounceCache()
	ptcaches.InitPeerCache()
	err = database.InitDB()
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	database.InitDBData()

	api.NewServer().Run(":8080")
}
func test() {

}
