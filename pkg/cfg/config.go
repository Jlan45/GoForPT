package cfg

import (
	"gopkg.in/yaml.v3"
	"os"
)

type SMTPConfig struct {
	Host     string // SMTP 服务器地址
	Port     int    // SMTP 服务器端口
	Username string // SMTP 用户名
	Password string // SMTP 密码
}
type CacheConfig struct {
	RegTime      int
	LoginTime    int
	AnnounceTime int
}
type DatabaseConfig struct {
	Driver   string // 数据库驱动
	Host     string // 数据库地址
	Port     int    // 数据库端口
	Username string // 数据库用户名
	Password string // 数据库密码
	Name     string // 数据库名称
}
type SiteConfig struct {
	Title       string // 网站名称
	Host        string // 网站domain+port
	SSL         bool   //是否开启https，不影响本程序功能，只影响生成的torrent中的announce
	Description string // 网站简介（首页展示，用markdown）
	TorrentPath string // 种子文件存放路径
	StaticPath  string // 静态文件存放路径
}
type Config struct {
	SMTP     SMTPConfig
	Cache    CacheConfig
	Database DatabaseConfig
	Site     SiteConfig
}

var Cfg Config

func LoadConfig() error {
	// 读取 YAML 文件
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}
	// 解析 YAML 文件
	err = yaml.Unmarshal(data, &Cfg)
	if err != nil {
		return err
	}
	return nil
}
