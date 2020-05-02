package core

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/dekoch/gouniversal/module/instafollowbot/global"
	"github.com/dekoch/gouniversal/module/instafollowbot/instauser"
	"github.com/dekoch/gouniversal/shared/api/instaclient"
	"github.com/dekoch/gouniversal/shared/api/instafollow"
	"github.com/dekoch/gouniversal/shared/api/instamedia"
	"github.com/dekoch/gouniversal/shared/api/instashareddata"
	"github.com/dekoch/gouniversal/shared/api/instatag"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/sqlite3"
)

type jobRet struct {
	err error
	cnt int
}

var (
	chanFollowFinished   = make(chan jobRet)
	chanUnfollowFinished = make(chan jobRet)
)

func LoadConfig() {

	var dbconn sqlite3.SQLite

	err := dbconn.Open(global.Config.DBFile)
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

	go job()

	intvl := global.Config.GetCheckInterval()

	if intvl == time.Duration(-1*time.Minute) ||
		(time.Since(global.Config.GetUnfollowTime()) > intvl && intvl > time.Duration(0)) {

		go func() {
			chanUnfollowFinished <- unfollow(global.Config.GetUnfollowCount())
		}()

		return
	}

	if intvl == time.Duration(-1*time.Minute) ||
		(time.Since(global.Config.GetFollowTime()) > intvl && intvl > time.Duration(0)) {

		go func() {
			chanFollowFinished <- follow(global.Config.GetFollowCount())
		}()
	}
}

func Exit() {

}

func job() {

	var unfollowFollow bool

	intvlCheck := global.Config.GetCheckInterval()
	tCheck := time.NewTimer(intvlCheck)

	for {

		select {
		case <-tCheck.C:
			tCheck.Stop()

			if intvlCheck > 0 {
				if unfollowFollow {

					go func() {
						chanFollowFinished <- follow(global.Config.GetFollowCount())
					}()
				} else {

					go func() {
						chanUnfollowFinished <- unfollow(global.Config.GetUnfollowCount())
					}()
				}
			} else {
				intvlCheck = global.Config.GetCheckInterval()
				if intvlCheck > 0 {
					tCheck.Reset(intvlCheck)
				} else {
					tCheck.Reset(time.Duration(10) * time.Second)
				}
			}

		case ret := <-chanFollowFinished:
			if ret.err != nil {
				console.Log(ret.err, "InstaFollowBot follow")
			}

			unfollowFollow = false
			tCheck.Reset(intvlCheck)

			if ret.cnt == 0 {
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

func follow(cnt int) jobRet {

	var ret jobRet

	for i := 0; i <= 2; i++ {

		switch i {
		case 0:
			ret.err = global.Config.LoadConfig()

		case 1:
			var retTag jobRet
			tags := global.Config.GetTags()

			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(tags), func(i, j int) { tags[i], tags[j] = tags[j], tags[i] })

			for _, tag := range tags {

				if ret.cnt >= cnt {
					continue
				}

				retTag = checkTag(tag, cnt-ret.cnt)
				if ret.err != nil {
					return ret
				}

				ret.cnt += retTag.cnt
			}

		case 2:
			global.Config.SetFollowTime(time.Now())
			ret.err = global.Config.SaveConfig()
		}

		if ret.err != nil {
			return ret
		}
	}

	return ret
}

func checkTag(tagname string, first int) jobRet {

	var (
		ret    jobRet
		dbconn sqlite3.SQLite
		ic     instaclient.InstaClient
		co     instaclient.Cookies
		users  []instauser.InstaUser
	)

	for i := 0; i <= 5; i++ {

		switch i {
		case 0:
			co = global.Config.GetCookies()
			ret.err = ic.SetCookies(co)

		case 1:
			users, ret.err = getUserFromTag(tagname, first*2, &ic)

		case 2:
			ret.err = dbconn.Open(global.Config.GetDBFile())

		case 3:
			defer dbconn.Close()

		case 4:
			var (
				exists bool
				resp   instashareddata.Response
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

				users[i].UserName, ret.err = getUserName(users[i].UserID, &ic)
				if ret.err != nil {
					return ret
				}

				resp, ret.err = instashareddata.GetResponse(users[i].UserName, &ic)
				if ret.err != nil {
					console.Log(ret.err, "")
					continue
				}

				if len(resp.Config.CsrfToken) > 0 {

					co.CsrfToken = resp.Config.CsrfToken

					ret.err = ic.SetCookies(co)
					if ret.err != nil {
						return ret
					}
				}

				if len(resp.EntryData.ProfilePage) == 0 {
					continue
				}

				if resp.EntryData.ProfilePage[0].GraphQL.User.IsBusinessAccount == true {
					continue
				}

				if resp.EntryData.ProfilePage[0].GraphQL.User.IsVerified == true {
					continue
				}

				if resp.EntryData.ProfilePage[0].GraphQL.User.IsPrivate == true {
					continue
				}

				if resp.EntryData.ProfilePage[0].GraphQL.User.EdgeFollowedBy.Count > 5000 {
					continue
				}

				ret.err = followUser(users[i], tagname, &ic, &dbconn)
				if ret.err != nil {
					return ret
				}

				ret.cnt++
			}

		case 5:
			global.Config.SetCookies(co)
		}

		if ret.err != nil {
			ret.err = errors.New("checkTag() " + strconv.Itoa(i) + " " + ret.err.Error())
			return ret
		}
	}

	return ret
}

func getUserFromTag(tagname string, first int, ic *instaclient.InstaClient) ([]instauser.InstaUser, error) {

	var (
		err  error
		ret  []instauser.InstaUser
		hash string
		resp instatag.Response
	)

	for i := 0; i <= 2; i++ {

		switch i {
		case 0:
			hash, err = global.Config.TagQueryHash.GetHash()

		case 1:
			resp, err = instatag.GetResponse(tagname, hash, first, "", ic)

		case 2:
			co := global.Config.GetCookies()

			var n instauser.InstaUser
			n.AccountID = co.DsUserID

			for _, edge := range resp.Data.Hashtag.EhtMedia.Edges {

				n.UserID = edge.Node.Owner.ID
				ret = append(ret, n)
			}
		}

		if err != nil {
			return ret, err
		}
	}

	return ret, nil
}

func getUserName(userid string, ic *instaclient.InstaClient) (string, error) {

	var (
		err  error
		hash string
		resp instamedia.Response
	)

	for i := 0; i <= 2; i++ {

		switch i {
		case 0:
			hash, err = global.Config.MediaQueryHash.GetHash()

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

func followUser(user instauser.InstaUser, tagname string, ic *instaclient.InstaClient, dbconn *sqlite3.SQLite) error {

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
			err = instafollow.Follow(user.UserID, user.UserName, ic)

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

func unfollow(cnt int) jobRet {

	var (
		ret    jobRet
		dbconn sqlite3.SQLite
		ic     instaclient.InstaClient
		co     instaclient.Cookies
		users  []instauser.InstaUser
	)

	for i := 0; i <= 6; i++ {

		switch i {
		case 0:
			ret.err = global.Config.LoadConfig()

		case 1:
			ret.err = dbconn.Open(global.Config.GetDBFile())

		case 2:
			defer dbconn.Close()

		case 3:
			fromDate := time.Now().AddDate(-999, 0, 0)
			toDate := time.Now().Add(-global.Config.GetUnfollowAfter())

			users, ret.err = instauser.GetUsersFollowingBetween(fromDate, toDate, &dbconn)

			if len(users) == 0 {
				return ret
			}

		case 4:
			co = global.Config.GetCookies()
			ret.err = ic.SetCookies(co)

		case 5:
			var (
				resp     instashareddata.Response
				userName string
			)

			for i := range users {

				if ret.cnt >= cnt {
					continue
				}

				userName, _ = getUserName(users[i].UserID, &ic)
				if len(userName) > 0 {
					users[i].UserName = userName
				}

				resp, _ = instashareddata.GetResponse(users[i].UserName, &ic)
				if len(resp.Config.CsrfToken) > 0 {

					co.CsrfToken = resp.Config.CsrfToken

					ret.err = ic.SetCookies(co)
					if ret.err != nil {
						return ret
					}
				}

				ret.err = unfollowUser(users[i], &ic, &dbconn)
				if ret.err != nil {
					return ret
				}

				ret.cnt++
			}

		case 6:
			global.Config.SetCookies(co)
			global.Config.SetUnfollowTime(time.Now())
			ret.err = global.Config.SaveConfig()
		}

		if ret.err != nil {
			ret.err = errors.New("unfollow() " + strconv.Itoa(i) + " " + ret.err.Error())
			return ret
		}
	}

	return ret
}

func unfollowUser(user instauser.InstaUser, ic *instaclient.InstaClient, dbconn *sqlite3.SQLite) error {

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
			err = instafollow.Unfollow(user.UserID, user.UserName, ic)

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
