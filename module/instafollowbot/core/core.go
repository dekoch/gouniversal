package core

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/dekoch/gouniversal/module/instafollowbot/core/coreconfig"
	"github.com/dekoch/gouniversal/module/instafollowbot/global"
	"github.com/dekoch/gouniversal/module/instafollowbot/instauser"
	"github.com/dekoch/gouniversal/shared/api/instaclient"
	"github.com/dekoch/gouniversal/shared/api/instafollow"
	"github.com/dekoch/gouniversal/shared/api/instamedia"
	"github.com/dekoch/gouniversal/shared/api/instashareddata"
	"github.com/dekoch/gouniversal/shared/api/instatag"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/sqlite3"
)

type Core struct {
	Config      coreconfig.CoreConfig
	instaClient instaclient.InstaClient
}

type jobRet struct {
	err error
	cnt int
}

var (
	chanFollowFinished   = make(chan jobRet)
	chanUnfollowFinished = make(chan jobRet)
)

func (co *Core) LoadConfig(conf coreconfig.CoreConfig) {

	co.Config = conf

	var dbconn sqlite3.SQLite

	err := dbconn.Open(global.Config.GetDBFile())
	if err != nil {
		console.Log(err, "InstaFollowBot")
		return
	}

	defer dbconn.Close()

	err = instauser.LoadConfig(&dbconn)
	if err != nil {
		console.Log(err, "InstaFollowBot")
		return
	}

	err = co.instaClient.SetCookies(co.Config.GetCookies())
	if err != nil {
		console.Log(err, "InstaFollowBot")
		return
	}

	go co.job()
}

func (co *Core) Exit() {

}

func (co *Core) job() {

	var unfollowFollow bool

	intvlCheck := co.Config.GetCheckInterval()
	tCheck := time.NewTimer(1 * time.Second)

	for {

		select {
		case <-tCheck.C:
			tCheck.Stop()

			if intvlCheck > 0 {
				if unfollowFollow {

					if time.Since(co.Config.GetFollowTime()) > intvlCheck {

						go func() {
							chanFollowFinished <- co.follow()
						}()
					} else {
						tCheck.Reset(10 * time.Second)
					}
				} else {

					if time.Since(co.Config.GetUnfollowTime()) > intvlCheck {

						go func() {
							chanUnfollowFinished <- co.unfollow()
						}()
					} else {
						tCheck.Reset(10 * time.Second)
					}
				}
			} else {
				intvlCheck = co.Config.GetCheckInterval()
				if intvlCheck > 0 {
					tCheck.Reset(intvlCheck)
				} else {
					tCheck.Reset(10 * time.Second)
				}
			}

		case ret := <-chanFollowFinished:
			unfollowFollow = false
			tCheck.Reset(intvlCheck)

			if ret.err != nil {
				console.Log(ret.err, "InstaFollowBot follow")

			} else if ret.cnt == 0 {
				tCheck.Reset(1 * time.Second)
			}

		case ret := <-chanUnfollowFinished:

			tCheck.Reset(intvlCheck)

			if ret.err == nil {

				unfollowFollow = true

				if ret.cnt == 0 {
					tCheck.Reset(1 * time.Second)
				}
			} else {
				console.Log(ret.err, "InstaFollowBot unfollow")
			}
		}
	}
}

func (co *Core) follow() jobRet {

	var ret jobRet

	for i := 0; i <= 4; i++ {

		switch i {
		case 0:
			ret.err = global.Config.LoadConfig()

		case 1:
			co.Config, ret.err = global.Config.GetCoreConfig(co.Config.CoreUUID)

		case 2:
			var retTag jobRet
			tags := co.Config.GetTags()
			cnt := co.Config.GetFollowCount()

			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(tags), func(i, j int) { tags[i], tags[j] = tags[j], tags[i] })

			for _, tag := range tags {

				if ret.cnt >= cnt {
					continue
				}

				if functions.IsEmpty(tag) {
					continue
				}

				retTag = co.checkTag(tag, cnt-ret.cnt)
				if retTag.err != nil {
					ret.err = retTag.err
					return ret
				}

				ret.cnt += retTag.cnt
			}

		case 3:
			co.Config.SetFollowTime(time.Now())
			ret.err = global.Config.SetCoreConfig(co.Config)

		case 4:
			ret.err = global.Config.SaveConfig()
		}

		if ret.err != nil {
			return ret
		}
	}

	return ret
}

func (co *Core) checkTag(tagname string, first int) jobRet {

	var (
		ret    jobRet
		dbconn sqlite3.SQLite
		coo    instaclient.Cookies
		users  []instauser.InstaUser
	)

	for i := 0; i <= 5; i++ {

		switch i {
		case 0:
			users, ret.err = co.getUserFromTag(tagname, first*5, &co.instaClient)

		case 1:
			ret.err = dbconn.Open(global.Config.GetDBFile())

		case 2:
			defer dbconn.Close()

		case 3:
			var (
				exists bool
				resp   instashareddata.Response
				n      int
			)

			for i := range users {

				if ret.cnt >= first {
					continue
				}

				exists, ret.err = instauser.Exists(users[i].UserID, &dbconn)
				if ret.err != nil {
					return ret
				}

				if exists {
					continue
				}

				users[i].UserName, ret.err = co.getUserName(users[i].UserID, &co.instaClient)
				if ret.err != nil {
					return ret
				}

				resp, ret.err = instashareddata.GetResponse(users[i].UserName, &co.instaClient)
				if ret.err != nil {
					console.Log(ret.err, "")
					continue
				}

				if len(resp.EntryData.ProfilePage) == 0 {
					continue
				}

				if co.Config.GetFollowBusinessAccount() == false &&
					resp.EntryData.ProfilePage[0].GraphQL.User.IsBusinessAccount {
					continue
				}

				if co.Config.GetFollowVerified() == false &&
					resp.EntryData.ProfilePage[0].GraphQL.User.IsVerified {
					continue
				}

				if co.Config.GetFollowPrivate() == false &&
					resp.EntryData.ProfilePage[0].GraphQL.User.IsPrivate {
					continue
				}

				n = co.Config.GetMaxFollow()

				if resp.EntryData.ProfilePage[0].GraphQL.User.EdgeFollow.Count > n && n > 0 {
					continue
				}

				n = co.Config.GetMaxFollowedBy()

				if resp.EntryData.ProfilePage[0].GraphQL.User.EdgeFollowedBy.Count > n && n > 0 {
					continue
				}

				if len(resp.RolloutHash) > 0 {

					ret.err = co.followUser(users[i], tagname, resp.RolloutHash, &co.instaClient, &dbconn)
					if ret.err != nil {
						return ret
					}

					if co.Config.GetAddNewTags() {
						co.Config.AddTags(users[i].NewTags)
					}
				}

				ret.cnt++
			}

		case 4:
			coo, ret.err = co.instaClient.GetCookies()

		case 5:
			co.Config.SetCookies(coo)
		}

		if ret.err != nil {
			ret.err = errors.New("checkTag() " + strconv.Itoa(i) + " " + ret.err.Error())
			return ret
		}
	}

	return ret
}

func (co *Core) getUserFromTag(tagname string, first int, ic *instaclient.InstaClient) ([]instauser.InstaUser, error) {

	var (
		err  error
		ret  []instauser.InstaUser
		hash string
		resp instatag.Response
	)

	for i := 0; i <= 2; i++ {

		switch i {
		case 0:
			hash, err = co.Config.TagQueryHash.GetHash()

		case 1:
			resp, err = instatag.GetResponse(tagname, hash, first, "", ic)

		case 2:
			coo := co.Config.GetCookies()

			var n instauser.InstaUser
			n.AccountID = coo.DsUserID

			for _, edge := range resp.Data.Hashtag.EhtMedia.Edges {

				n.UserID = edge.Node.Owner.ID

				if len(edge.Node.EmtCaption.Edges) > 0 {
					n.NewTags, err = co.getTagsFromString(edge.Node.EmtCaption.Edges[0].Node.Text)
				}

				ret = append(ret, n)
			}
		}

		if err != nil {
			return ret, err
		}
	}

	return ret, nil
}

func (co *Core) getUserName(userid string, ic *instaclient.InstaClient) (string, error) {

	var (
		err  error
		hash string
		resp instamedia.Response
	)

	for i := 0; i <= 2; i++ {

		switch i {
		case 0:
			hash, err = co.Config.MediaQueryHash.GetHash()

		case 1:
			resp, err = instamedia.GetResponse(userid, hash, 10, "", ic)

		case 2:
			for _, edge := range resp.Data.User.EottMedia.Edges {

				if edge.Node.Owner.ID == userid {
					return edge.Node.Owner.UserName, nil
				}
			}
		}

		if err != nil {
			return "", err
		}
	}

	return "", errors.New("username not found")
}

func (co *Core) followUser(user instauser.InstaUser, tagname, xinstagramajax string, ic *instaclient.InstaClient, dbconn *sqlite3.SQLite) error {

	var err error

	for i := 0; i <= 4; i++ {

		switch i {
		case 0:
			dbconn.Tx, err = dbconn.DB.Begin()

		case 1:
			defer func() {
				if err != nil {
					dbconn.Tx.Rollback()
				}
			}()

		case 2:
			err = instafollow.Follow(user.UserID, user.UserName, xinstagramajax, ic)

			user.Tag = tagname
			user.Following = true
			user.Follow = time.Now()

		case 3:
			err = user.Save(dbconn.Tx)

		case 4:
			err = dbconn.Tx.Commit()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (co *Core) unfollow() jobRet {

	var (
		ret    jobRet
		dbconn sqlite3.SQLite
		coo    instaclient.Cookies
		users  []instauser.InstaUser
	)

	for i := 0; i <= 8; i++ {

		switch i {
		case 0:
			ret.err = global.Config.LoadConfig()

		case 1:
			co.Config, ret.err = global.Config.GetCoreConfig(co.Config.CoreUUID)

		case 2:
			ret.err = dbconn.Open(global.Config.GetDBFile())

		case 3:
			defer dbconn.Close()

		case 4:
			fromDate := time.Now().AddDate(-999, 0, 0)
			toDate := time.Now().Add(-co.Config.GetUnfollowAfter())

			users, ret.err = instauser.GetUsersFollowingBetween(fromDate, toDate, &dbconn)

			if len(users) == 0 {
				return ret
			}

		case 5:
			var (
				resp     instashareddata.Response
				userName string
			)

			cnt := co.Config.GetUnfollowCount()

			for i := range users {

				if ret.cnt >= cnt {
					continue
				}

				userName, _ = co.getUserName(users[i].UserID, &co.instaClient)
				if len(userName) > 0 {
					users[i].UserName = userName
				}

				resp, ret.err = instashareddata.GetResponse(users[i].UserName, &co.instaClient)
				if ret.err != nil {
					console.Log(ret.err, "")
					continue
				}

				if len(resp.RolloutHash) > 0 {

					ret.err = co.unfollowUser(users[i], resp.RolloutHash, &co.instaClient, &dbconn)
					if ret.err != nil {
						return ret
					}
				}

				ret.cnt++
			}

		case 6:
			coo, ret.err = co.instaClient.GetCookies()

		case 7:
			co.Config.SetCookies(coo)
			co.Config.SetUnfollowTime(time.Now())
			ret.err = global.Config.SetCoreConfig(co.Config)

		case 8:
			ret.err = global.Config.SaveConfig()
		}

		if ret.err != nil {
			ret.err = errors.New("unfollow() " + strconv.Itoa(i) + " " + ret.err.Error())
			return ret
		}
	}

	return ret
}

func (co *Core) unfollowUser(user instauser.InstaUser, xinstagramajax string, ic *instaclient.InstaClient, dbconn *sqlite3.SQLite) error {

	var err error

	for i := 0; i <= 4; i++ {

		switch i {
		case 0:
			dbconn.Tx, err = dbconn.DB.Begin()

		case 1:
			defer func() {
				if err != nil {
					dbconn.Tx.Rollback()
				}
			}()

		case 2:
			err = instafollow.Unfollow(user.UserID, user.UserName, xinstagramajax, ic)

			user.Following = false
			user.Unfollow = time.Now()

		case 3:
			err = user.Save(dbconn.Tx)

		case 4:
			err = dbconn.Tx.Commit()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (co *Core) getTagsFromString(s string) ([]string, error) {

	var (
		err error
		ret []string
	)

	s = strings.Replace(s, "\r", " ", -1)
	s = strings.Replace(s, "\n", " ", -1)

	splitSuffix := strings.SplitAfter(s, " ")

	for i := range splitSuffix {

		if strings.Contains(splitSuffix[i], "#") == false {
			continue
		}

		splitPrefix := strings.Split(splitSuffix[i], "#")

		for ii := range splitPrefix {

			tag := strings.TrimSpace(splitPrefix[ii])

			if len(tag) == 0 {
				continue
			}

			ret = append(ret, tag)
		}
	}

	return ret, err
}
