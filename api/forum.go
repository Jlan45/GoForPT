package api

import (
	"GoForPT/model"
	"GoForPT/pkg/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ThreadData struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	//Author string `json:"author"` auther直接用session拿
	Torrent []string `json:"torrent"`
	Tag     []string `json:"tag"`
}
type ThreadListData struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	AuthorID uint   `json:"author_id"`
}

func PostThread(c *gin.Context) {
	postdata := ThreadData{}
	c.Bind(&postdata)
	//校验postdata，如果没有torrent说明是纯论坛
	if postdata.Title == "" || postdata.Content == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Invalid request",
		})
	}
	if len(postdata.Torrent) == 0 {
		thread := model.Thread{
			Title:    postdata.Title,
			Content:  postdata.Content,
			AuthorID: c.GetUint("UserID"),
			Tags:     postdata.Tag,
			Torrents: nil,
		}
		database.DB.Create(&thread)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "New thread",
		})
	} else {
		thread := model.Thread{
			Title:    postdata.Title,
			Content:  postdata.Content,
			AuthorID: c.GetUint("UserID"),
			Tags:     postdata.Tag,
			Torrents: postdata.Torrent,
		}
		database.DB.Create(&thread)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "New thread",
		})
	}
	return

}
func GetThread(c *gin.Context) {
	//获取thread
	threadid := c.Param("id")
	thread := model.Thread{}
	database.DB.Where("id = ?", threadid).First(&thread)
	if thread.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Thread not found",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"title":   thread.Title,
			"content": thread.Content,
			"torrent": thread.Torrents,
			"tag":     thread.Tags,
			"author":  thread.AuthorID,
		},
	})
}
func ListThread(c *gin.Context) {
	l := c.Query("limit")
	listnum, err := strconv.Atoi(l)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Invalid limit",
		})
		return
	}
	//获取thread
	threadListData := []ThreadListData{}
	database.DB.Model(&model.Thread{}).Order("created_at desc").Limit(listnum).Find(&threadListData)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": threadListData,
	})
	return

}
