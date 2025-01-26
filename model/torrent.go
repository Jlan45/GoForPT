package model

import "gorm.io/gorm"

type EventType uint8

const (
	Started EventType = iota
	Stopped
	Completed
)

type UserConn struct {
	IP   string
	Port int
}
type Torrent struct {
	//入库内容
	gorm.Model
	Hash    []byte //种子文件的InfoHash
	Name    string
	MD5Sum  string //以TorrentsFile为起始，末尾加.torrent后缀作为文件种子哈希
	OwnerID uint
	Owner   User `gorm:"foreignKey:OwnerID"`
}
type TorrentFileData struct {
	Announce string //tracker地址
}
type AnnounceRequest struct {
	Token      string   `form:"token"`
	InfoHash   string   `form:"info_hash"`
	PeerID     string   `form:"peer_id"`
	Port       string   `form:"port"`
	Uploaded   string   `form:"uploaded"`
	Downloaded string   `form:"downloaded"`
	Left       string   `form:"left"`
	Event      string   `form:"event"`
	Numwant    string   `form:"numwant"`
	IPv6       []string `form:"ipv6"`
	IP         []string `form:"ip"`
}
