package services

import (
	"strings"
	"sync/atomic"
	"time"

	"im-server/commons/caches"
	"im-server/commons/tools"
)

type UserStatus struct {
	appkey string
	userId string
	// LastSyncTime        *int64
	// LastSendBoxSyncTime *int64
	LatestMsgTime *int64 // latest msg time
	// LatestSendMsgTime *int64
	TerminalNum  int
	OnlineStatus bool //online state

	isNtf bool //is ntf

	PushSwitch int32
}

var userOnlineStatusCache *caches.LruCache
var userLocks *tools.SegmentatedLocks

func init() {
	userOnlineStatusCache = caches.NewLruCacheWithReadTimeout(100000, func(key, value interface{}) {}, 10*time.Minute)
	userLocks = tools.NewSegmentatedLocks(512)
}

/*
record user's  status when sync msg
*/
func RecordUserOnlineStatus(appKey, userId string, onlineStatus bool, terminalNum int) {
	user := GetUserStatus(appKey, userId)
	key := getKey(appKey, userId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	user.OnlineStatus = onlineStatus
	user.TerminalNum = terminalNum
}

func (user *UserStatus) IsOnline() bool {
	return user.OnlineStatus
}

func (user *UserStatus) SetPushSwitch(pushSwitch int32) {
	atomic.StoreInt32(&user.PushSwitch, pushSwitch)
}

func (user *UserStatus) OpenPushSwitch() bool {
	return user.PushSwitch > 0
}

func (user *UserStatus) CheckNtfWithSwitch() bool {
	if !user.OnlineStatus || user.TerminalNum > 1 {
		return true
	}
	if user.isNtf {
		return true
	} else {
		key := getKey(user.appkey, user.userId)
		lock := userLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if user.isNtf {
			return true
		} else {
			ret := user.isNtf
			user.isNtf = true
			return ret
		}
	}
}
func (user *UserStatus) SetNtfStatus(isNtf bool) {
	key := getKey(user.appkey, user.userId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	user.isNtf = isNtf
}
func (user *UserStatus) SetLatestMsgTime(time int64) {
	key := getKey(user.appkey, user.userId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if user.LatestMsgTime == nil || *user.LatestMsgTime < time {
		user.LatestMsgTime = &time
	}
}

func GetUserStatus(appKey, userId string) *UserStatus {
	key := getKey(appKey, userId)
	if val, exist := userOnlineStatusCache.Get(key); exist {
		return val.(*UserStatus)
	} else {
		l := userLocks.GetLocks(appKey, userId)
		l.Lock()
		defer l.Unlock()
		if val, exist := userOnlineStatusCache.Get(key); exist {
			return val.(*UserStatus)
		} else {
			userInfo := initUserInfo(appKey, userId)
			userOnlineStatusCache.Add(key, userInfo)
			return userInfo
		}
	}
}

func RegenateSendTime(appkey, userId string, currentTime int64) int64 {
	user := GetUserStatus(appkey, userId)

	key := getKey(appkey, userId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	ret := currentTime
	if user.LatestMsgTime == nil || currentTime > *user.LatestMsgTime {
		user.LatestMsgTime = &currentTime
	} else {
		ret = *user.LatestMsgTime + 1
		user.LatestMsgTime = &ret
	}
	return ret
}

func getKey(appkey, userId string) string {
	return strings.Join([]string{appkey, userId}, "_")
}

func initUserInfo(appkey, userId string) *UserStatus {
	return &UserStatus{
		appkey:       appkey,
		userId:       userId,
		OnlineStatus: true,
	}
}
