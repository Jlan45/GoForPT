package api

import (
	"GoForPT/pkg/cfg"
	"GoForPT/pkg/tools"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	r := gin.Default()
	v1api := r.Group("/api/v1")
	store := cookie.NewStore([]byte(tools.GenerateToken()))
	v1api.Use(sessions.Sessions("goforpt", store))
	userAPI := v1api.Group("/user")
	userAPI.POST("/login", UserLogin)
	userAPI.POST("/register", UserRegister)
	userAPI.POST("/register/code", SendRegisterOTP)
	userAPI.POST("/login/code", SendLoginOTP)
	userAPI.GET("/list/torrents", CheckLogin, ListUserTorrents)

	r.GET("/announce", NoMoreThunder, Announce)

	torrentAPI := v1api.Group("/torrent")
	torrentAPI.GET("/file/:id", CheckLogin, GetTorrentFile)
	torrentAPI.POST("/upload", CheckLogin, UploadTorrentFile)

	forumAPI := v1api.Group("/forum")
	forumAPI.GET("/list/thread", CheckLogin, ListThread)
	forumAPI.POST("/thread/post", CheckLogin, PostThread)
	forumAPI.GET("/thread/:id", CheckLogin, GetThread)

	r.Static("/static/user", cfg.Cfg.Site.StaticPath)
	r.Static("/web", "web")

	r.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/web")
	})
	return r
}
