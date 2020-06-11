package coreconfig

import (
	"sync"
	"time"

	"github.com/dekoch/gouniversal/shared/api/instaclient"
	"github.com/dekoch/gouniversal/shared/hashstor"
)

type CoreConfig struct {
	CoreUUID              string
	CheckInterv           int // minutes (0=disabled)
	FollowCnt             int
	UnfollowCnt           int
	UnfollowAfter         int // hours
	Tags                  []string
	AddNewTags            bool
	FollowBusinessAccount bool
	FollowVerified        bool
	FollowPrivate         bool
	MaxFollow             int // (0=disabled)
	MaxFollowedBy         int // (0=disabled)
	Cookies               instaclient.Cookies
	TagQueryHash          hashstor.HashStor
	MediaQueryHash        hashstor.HashStor
	FollowTime            time.Time
	UnfollowTime          time.Time
}

var mut sync.RWMutex

func (hc *CoreConfig) LoadDefaults() {

	hc.CheckInterv = 25
	hc.FollowCnt = 15
	hc.UnfollowCnt = 15
	hc.UnfollowAfter = 2

	hc.Tags = append(hc.Tags, "")
	hc.AddNewTags = false

	hc.FollowBusinessAccount = false
	hc.FollowVerified = false
	hc.FollowPrivate = false
	hc.MaxFollow = 0
	hc.MaxFollowedBy = 1000

	hc.TagQueryHash.Add("")
	hc.MediaQueryHash.Add("")

	hc.FollowTime = time.Now().AddDate(-1, 0, 0)
	hc.UnfollowTime = time.Now().AddDate(-1, 0, 0)
}

func (hc *CoreConfig) LoadConfig() error {

	hc.TagQueryHash.Init()
	hc.MediaQueryHash.Init()

	return nil
}

func (hc *CoreConfig) SetCoreUUID(uid string) {

	mut.Lock()
	defer mut.Unlock()

	hc.CoreUUID = uid
}

func (hc *CoreConfig) GetCoreUUID() string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.CoreUUID
}

func (hc *CoreConfig) SetCheckInterval(minutes int) {

	mut.Lock()
	defer mut.Unlock()

	hc.CheckInterv = minutes
}

func (hc *CoreConfig) GetCheckInterval() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.CheckInterv) * time.Minute
}

func (hc *CoreConfig) SetFollowCount(cnt int) {

	mut.Lock()
	defer mut.Unlock()

	hc.FollowCnt = cnt
}

func (hc *CoreConfig) GetFollowCount() int {

	mut.RLock()
	defer mut.RUnlock()

	return hc.FollowCnt
}

func (hc *CoreConfig) SetUnfollowCount(cnt int) {

	mut.Lock()
	defer mut.Unlock()

	hc.UnfollowCnt = cnt
}

func (hc *CoreConfig) GetUnfollowCount() int {

	mut.RLock()
	defer mut.RUnlock()

	return hc.UnfollowCnt
}

func (hc *CoreConfig) SetUnfollowAfter(hours int) {

	mut.Lock()
	defer mut.Unlock()

	hc.UnfollowAfter = hours
}

func (hc *CoreConfig) GetUnfollowAfter() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.UnfollowAfter) * time.Hour
}

func (hc *CoreConfig) SetTags(tags []string) {

	mut.Lock()
	defer mut.Unlock()

	hc.Tags = tags
}

func (hc *CoreConfig) GetTags() []string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Tags
}

func (hc *CoreConfig) SetAddNewTags(b bool) {

	mut.Lock()
	defer mut.Unlock()

	hc.AddNewTags = b
}

func (hc *CoreConfig) GetAddNewTags() bool {

	mut.RLock()
	defer mut.RUnlock()

	return hc.AddNewTags
}

func (hc *CoreConfig) AddTags(tags []string) {

	mut.Lock()
	defer mut.Unlock()

	for i := range tags {

		found := false

		for ii := range hc.Tags {

			if tags[i] == hc.Tags[ii] {
				found = true
			}
		}

		if found == false {
			hc.Tags = append(hc.Tags, tags[i])
		}
	}
}

func (hc *CoreConfig) SetFollowBusinessAccount(b bool) {

	mut.Lock()
	defer mut.Unlock()

	hc.FollowBusinessAccount = b
}

func (hc *CoreConfig) GetFollowBusinessAccount() bool {

	mut.RLock()
	defer mut.RUnlock()

	return hc.FollowBusinessAccount
}

func (hc *CoreConfig) SetFollowVerified(b bool) {

	mut.Lock()
	defer mut.Unlock()

	hc.FollowVerified = b
}

func (hc *CoreConfig) GetFollowVerified() bool {

	mut.RLock()
	defer mut.RUnlock()

	return hc.FollowVerified
}

func (hc *CoreConfig) SetFollowPrivate(b bool) {

	mut.Lock()
	defer mut.Unlock()

	hc.FollowPrivate = b
}

func (hc *CoreConfig) GetFollowPrivate() bool {

	mut.RLock()
	defer mut.RUnlock()

	return hc.FollowPrivate
}

func (hc *CoreConfig) SetMaxFollow(n int) {

	mut.Lock()
	defer mut.Unlock()

	hc.MaxFollow = n
}

func (hc *CoreConfig) GetMaxFollow() int {

	mut.RLock()
	defer mut.RUnlock()

	return hc.MaxFollow
}

func (hc *CoreConfig) SetMaxFollowedBy(n int) {

	mut.Lock()
	defer mut.Unlock()

	hc.MaxFollowedBy = n
}

func (hc *CoreConfig) GetMaxFollowedBy() int {

	mut.RLock()
	defer mut.RUnlock()

	return hc.MaxFollowedBy
}

func (hc *CoreConfig) SetCookies(co instaclient.Cookies) {

	mut.Lock()
	defer mut.Unlock()

	hc.Cookies = co
}

func (hc *CoreConfig) GetCookies() instaclient.Cookies {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Cookies
}

func (hc *CoreConfig) SetFollowTime(t time.Time) {

	mut.Lock()
	defer mut.Unlock()

	hc.FollowTime = t
}

func (hc *CoreConfig) GetFollowTime() time.Time {

	mut.RLock()
	defer mut.RUnlock()

	return hc.FollowTime
}

func (hc *CoreConfig) SetUnfollowTime(t time.Time) {

	mut.Lock()
	defer mut.Unlock()

	hc.UnfollowTime = t
}

func (hc *CoreConfig) GetUnfollowTime() time.Time {

	mut.RLock()
	defer mut.RUnlock()

	return hc.UnfollowTime
}
