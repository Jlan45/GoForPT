package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Thread struct {
	gorm.Model
	Title    string         //标题按特定格式填写资源名称
	Content  string         //markdown写
	Tags     pq.StringArray `gorm:"type:text[]"`
	Torrents pq.StringArray `gorm:"type:text[]"` //只存种子的md5Sum（比infohash强）
	AuthorID uint
	Author   User `gorm:"foreignKey:AuthorID"`
}
type Post struct {
	gorm.Model
	Content  string
	ThreadID uint
	Thread   Thread `gorm:"foreignKey:ThreadID"`
	AuthorID uint
	Author   User `gorm:"foreignKey:AuthorID"`
}
