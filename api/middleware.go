package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
)

func CheckLogin(c *gin.Context) {
	//通过session校验是否登录
	session := sessions.Default(c)
	userid := session.Get("UserID")
	if userid == nil {
		session.Clear()
		session.Save()
		c.JSON(200, gin.H{"code": 401, "msg": "Please Login"})
		c.Abort()
	}
	c.Set("UserID", userid)
	c.Next()
}
func NoMoreThunder(c *gin.Context) {
	//干掉迅雷
	peerid := c.Query("peer_id")
	if strings.HasPrefix(peerid, "-XL") {
		c.JSON(200, gin.H{"code": 403, "msg": "No More Thunder"})
		c.Abort()
		return
	}
	c.Next()
}
