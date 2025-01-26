package tools

import (
	"GoForPT/model"
	"GoForPT/pkg/ptcaches"
	"bytes"
	"github.com/jackpal/bencode-go"
	"strconv"
)

type PeerResponse struct {
	Interval int              "interval"
	Peers    []model.UserConn "peers"
}

func GetSeedUserList(md5sum string) []uint {
	if _, found := ptcaches.AnnounceUserCache.Get(md5sum); !found {
		ptcaches.AnnounceUserCache.Set(md5sum, make([]uint, 0), 0)
	}
	peerList, _ := ptcaches.AnnounceUserCache.Get(md5sum)
	return peerList.([]uint)
}
func GenerateSeedTrackerResponse(md5sum string) []byte {
	userIDList, _ := ptcaches.AnnounceUserCache.Get(md5sum)
	peerList := PeerResponse{
		Interval: 600,
	}
	for _, userID := range userIDList.([]uint) {
		userConnList, existed := ptcaches.PeerCache.Get(strconv.Itoa(int(userID)))
		if existed {
			peerList.Peers = append(peerList.Peers, userConnList.([]model.UserConn)...)
		}
	}
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, peerList)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}
