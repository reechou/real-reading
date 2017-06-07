package ext

import (
	"strconv"
	"sync"
)

type TulingUserId struct {
	sync.Mutex

	userid        int
	userMap       map[string]int
	userMapString map[string]string
}

func NewTulingUserId() *TulingUserId {
	return &TulingUserId{
		userMap:       make(map[string]int),
		userMapString: make(map[string]string),
	}
}

func (self *TulingUserId) GetUserId(user string) int {
	self.Lock()
	defer self.Unlock()

	userid := self.userMap[user]
	if userid == 0 {
		self.userid++
		self.userMap[user] = self.userid
		return self.userid
	}
	return userid
}

func (self *TulingUserId) GetUserIdString(user string) string {
	self.Lock()
	defer self.Unlock()

	userid := self.userMapString[user]
	if userid == "" {
		self.userid++
		userid = strconv.Itoa(self.userid)
		self.userMapString[user] = userid
	}
	return userid
}
