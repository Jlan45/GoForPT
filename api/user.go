package api

import (
	"GoForPT/model"
	"GoForPT/pkg/cfg"
	"GoForPT/pkg/database"
	"GoForPT/pkg/email"
	"GoForPT/pkg/ptcaches"
	"GoForPT/pkg/tools"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

type RegisterData struct {
	Email    string `json:"email"`
	Password string `json:"password"` //数据库是passhash
	Username string `json:"username"`
	Verify   string `json:"verify"`
}

func UserRegister(c *gin.Context) {
	//前端哈希再传过来
	postdata := RegisterData{}
	err := c.BindJSON(&postdata)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Invalid request",
		})
		return
	}
	var count int64
	database.DB.Model(&model.User{}).Where("email = ?", postdata.Email).Count(&count)
	if count > 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Email already registered",
		})
		return
	}
	//检查验证码
	if verifyCode, found := ptcaches.EmailCache.Get(fmt.Sprintf("reg@%d", postdata.Email)); !found || verifyCode != postdata.Verify {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Verification code error",
		})
		return
	}
	//注册成功销毁验证码
	ptcaches.EmailCache.Delete(fmt.Sprintf("reg@%d", postdata.Email))
	var userGroup model.UserGroup
	database.DB.Where("name = ?", "Normal").First(&userGroup)
	//将密码进行md5哈希
	postdata.Password = tools.MD5([]byte(postdata.Password))
	//注册
	user := model.User{
		Email:       postdata.Email,
		Username:    postdata.Username,
		PassHash:    postdata.Password,
		Token:       tools.GenerateToken(),
		Uploaded:    0,
		Downloaded:  0,
		MagicPower:  0,
		UserGroupID: userGroup.ID,
	}

	database.DB.Create(&user)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "Register successfully",
	})
	return
}
func UserLogin(c *gin.Context) {
	session := sessions.Default(c)
	postdata := gin.H{}
	err := c.BindJSON(&postdata)
	if err != nil {
		return
	}
	var user model.User
	database.DB.Where("email = ? and pass_hash = ?", postdata["email"], tools.MD5([]byte(postdata["password"].(string)))).First(&user)
	if user.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Login failed",
		})
		return
	}
	session.Set("UserID", user.ID)
	session.Save()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "Login successfully",
	})
}
func SendRegisterOTP(c *gin.Context) {
	postdata := gin.H{}
	err := c.BindJSON(&postdata)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Invalid request",
		})
		return
	}
	usermail := postdata["email"].(string)
	if usermail == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Please send email",
		})
		return
	}
	//检查是否存在
	if _, found := ptcaches.EmailCache.Get(fmt.Sprintf("reg@%d", usermail)); found {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Please wait for 5 minutes before sending again",
		})
		return
	}
	//检查数据库是否存在用户
	var count int64
	database.DB.Model(&model.User{}).Where("email = ?", usermail).Count(&count)
	if count == 1 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Email already registered",
		})
		return
	}
	//生成6位随机数
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	err = email.SendEmail(usermail, fmt.Sprintf("[%s]注册验证码", cfg.Cfg.Site.Title), otp)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to send Email",
		})
		return
	}
	ptcaches.EmailCache.Set(fmt.Sprintf("reg@%d", usermail), otp, time.Duration(cfg.Cfg.Cache.RegTime)*time.Second)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "Email sent successfully",
	})
	return
}
func SendLoginOTP(c *gin.Context) {
	postdata := gin.H{}
	err := c.BindJSON(&postdata)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Invalid request",
		})
		return
	}
	usermail := postdata["email"].(string)
	if usermail == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Please send email",
		})
		return
	}
	//检查是否存在
	if _, found := ptcaches.EmailCache.Get(fmt.Sprintf("login@%d", usermail)); found {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Please wait for 5 minutes before sending again",
		})
		return
	}
	//检查是否注册
	var count int64
	database.DB.Model(&model.User{}).Where("email = ?", usermail).Count(&count)
	if count == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Please register first",
		})
		return
	}
	//生成6位随机数
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	err = email.SendEmail(usermail, fmt.Sprintf("[%s]登录验证码", cfg.Cfg.Site.Title), otp)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to send Email",
		})
		return
	}
	ptcaches.EmailCache.Set(fmt.Sprintf("login@%d", usermail), otp, time.Duration(cfg.Cfg.Cache.LoginTime)*time.Second)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "Email sent successfully",
	})
	return
}
func ListUserTorrents(c *gin.Context) {
	userID, _ := c.Get("UserID")
	var torrents []model.Torrent
	var torrentsData []gin.H
	database.DB.Where("OwnerID = ?", userID).Find(&torrents)
	for _, torrent := range torrents {
		torrentsData = append(torrentsData, gin.H{"MD5Sum": torrent.MD5Sum, "Name": torrent.Name})
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": torrentsData,
	})
}
