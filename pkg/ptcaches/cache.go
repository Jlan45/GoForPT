package ptcaches

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	EmailCache        *cache.Cache
	AnnounceUserCache *cache.Cache //结构是key为种子md5sum，value为当前种子的userid list，然后再去下面的Cache里面取对应peer的IPLIST，也就是说这个cache指向的是当前种子有多少用户在做种以及下载
	PeerCache         *cache.Cache
)

func InitEmailCache() error {
	EmailCache = cache.New(5*time.Minute, 1*time.Minute)
	return nil
}
func InitAnnounceCache() error {
	AnnounceUserCache = cache.New(24*time.Hour, 10*time.Minute)
	return nil
}
func InitPeerCache() error {
	PeerCache = cache.New(24*time.Hour, 10*time.Minute)
	return nil
}
