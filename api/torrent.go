package api

import (
	"GoForPT/model"
	"GoForPT/pkg/cfg"
	"GoForPT/pkg/database"
	"GoForPT/pkg/ptcaches"
	"GoForPT/pkg/tools"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackpal/bencode-go"
	"os"
	"strconv"
	"sync"
	"time"
)

var announceLock sync.Mutex

func GetTorrentFile(c *gin.Context) {
	torrentid := c.Param("id")
	userID, _ := c.Get("UserID")
	var user model.User
	if err := database.DB.Where("id = ?", userID).Take(&user).Error; err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "User not found",
		})
		return
	}
	if torrentid == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Invalid request",
		})
		return
	}
	torrentfile := model.Torrent{}
	//从数据库查询种子文件，用hash查
	database.DB.Where("md5_sum = ?", []byte(torrentid)).Take(&torrentfile)
	if torrentfile.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Torrent not found",
		})
		return
	}
	//读取种子文件，修改announce参数
	file, err := os.ReadFile(fmt.Sprintf("%s/%s.torrent", cfg.Cfg.Site.TorrentPath, torrentfile.MD5Sum))
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Torrent not found",
		})
		return
	}
	torr, _ := bencode.Decode(bytes.NewReader(file))
	torr.(map[string]interface{})["announce"] = fmt.Sprintf("http://%s/announce?token=%s", cfg.Cfg.Site.Host, user.Token)
	torr.(map[string]interface{})["announce-list"] = nil
	var torrentFileBuffer bytes.Buffer
	err = bencode.Marshal(&torrentFileBuffer, torr)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Invaild torrent file",
		})
		return
	}
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", torrentfile.Name))
	c.Writer.Header().Set("Content-Type", "application/x-bittorrent")
	c.Writer.WriteHeader(200)
	c.Writer.Write(torrentFileBuffer.Bytes())
	return
	//修改announce参数
	// 将种子文件内容返回给用户
	//从数据库查询用户token，然后修改torrentfile的
}
func UploadTorrentFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "Failed to upload file",
		})
		return
	}
	content, err := file.Open()
	defer content.Close()

	torrentData, err := bencode.Decode(content)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Invaild torrent file",
		})
		return
	}
	//
	torrentInfo := torrentData.(map[string]interface{})["info"].(map[string]interface{})
	//判断是否为私有种子，其实设定announce之后没有必要了就
	//if private, ok := torrentInfo["private"].(int); !ok || private != 1 {
	//	c.JSON(200, gin.H{
	//		"code": 400,
	//		"msg":  "Invaild torrent file",
	//	})
	//	return
	//}

	var buf bytes.Buffer
	err = bencode.Marshal(&buf, torrentInfo)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Invaild torrent file",
		})
		return
	}
	infoHash := tools.SHA1(buf.Bytes())
	md5sum := tools.MD5(infoHash[:])
	// Save the file to disk
	if err := c.SaveUploadedFile(file, fmt.Sprintf("%s/%s.torrent", cfg.Cfg.Site.TorrentPath, md5sum)); err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Failed to save file",
		})
		return
	}
	torrentFile := model.Torrent{
		Hash:    infoHash[:],
		Name:    file.Filename,
		MD5Sum:  md5sum,
		OwnerID: c.GetUint("UserID"),
	}
	if err := database.DB.Create(&torrentFile).Error; err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Failed to save file",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "File uploaded successfully",
	})
	return
}
func Announce(c *gin.Context) {
	req := model.AnnounceRequest{
		Token:      c.Query("token"),     //必要参数
		InfoHash:   c.Query("info_hash"), //必要参数
		PeerID:     c.Query("peer_id"),   //必要参数
		Port:       c.Query("port"),      //必要参数
		Uploaded:   c.Query("uploaded"),
		Downloaded: c.Query("downloaded"),
		Left:       c.Query("left"),
		Event:      c.Query("event"),
		Numwant:    c.Query("numwant"),
		IPv6:       c.QueryArray("ipv6"),
		IP:         c.QueryArray("ip"),
	}
	//检测必要参数是否存在
	if req.Token == "" || req.InfoHash == "" || req.PeerID == "" || req.Port == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Missing params",
		})
		return
	}

	//通过token查询用户
	var u model.User
	database.DB.Where("token = ?", req.Token).First(&u)
	if u.ID == 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "User not found",
		})
		return
	}
	//查询种子是否存在
	var torrent model.Torrent
	if err := database.DB.Where("hash = ?", []byte(req.InfoHash)).Take(&torrent).Error; err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "Torrent not found",
		})
		return
	}
	portnum, err := strconv.Atoi(req.Port)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Invalid port",
		})
		return
	}
	userConnectionList := make([]model.UserConn, 0)
	userConnectionList = append(userConnectionList, model.UserConn{
		IP:   c.RemoteIP(),
		Port: portnum,
	})
	for _, ipv6 := range req.IPv6 {
		userConnectionList = append(userConnectionList, model.UserConn{
			IP:   ipv6,
			Port: portnum,
		})
	}
	for _, ip := range req.IP {
		userConnectionList = append(userConnectionList, model.UserConn{
			IP:   ip,
			Port: portnum,
		})
	}
	//都有就可以开始处理announce响应了，根据两种不同event操作iplist
	//加锁
	announceLock.Lock()
	defer announceLock.Unlock()
	switch req.Event {
	case "started", "completed", "":
		uidList := tools.GetSeedUserList(torrent.MD5Sum)
		//判断有没有自己的userid
		hasPeerID := false
		for _, uid := range uidList {
			if uid == u.ID {
				hasPeerID = true
				break
			}
		}
		if !hasPeerID {
			uidList = append(uidList, u.ID)
		}
		ptcaches.AnnounceUserCache.Set(torrent.MD5Sum, uidList, 30*time.Minute)
		ptcaches.PeerCache.Set(strconv.Itoa(int(u.ID)), userConnectionList, 30*time.Minute)
		c.Writer.WriteHeader(200)
		_, err := c.Writer.Write(tools.GenerateSeedTrackerResponse(torrent.MD5Sum))
		if err != nil {
			return
		}
		return
	case "stopped":
		//删除对应peer的缓存，其实直接从announce那个list里买呢删除就好了
		ptcaches.PeerCache.Delete(strconv.Itoa(int(u.ID)))
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "Stopped torrent",
		})
	}
}
