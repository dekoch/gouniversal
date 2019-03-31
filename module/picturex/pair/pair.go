package pair

import (
	"errors"
	"image"
	"image/color"
	"time"

	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/file"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

const (
	MISSING int = 1 + iota
	LOCKED
	UNLOCKED
	VISIBLE
)

type Picture struct {
	user  string
	name  string
	state int
}

type Pair struct {
	uuid      string
	timestamp time.Time
	first     Picture
	second    Picture
}

var (
	pathSource      string
	pathDestination string
	pathStatic      string
)

func New(user string) Pair {

	var ret Pair

	u := uuid.Must(uuid.NewRandom())
	ret.uuid = u.String()

	ret.timestamp = time.Now()

	ret.first.user = user
	ret.first.state = MISSING

	ret.second.state = MISSING

	return ret
}

func SetSourcePath(p string) {
	pathSource = p
}

func SetDestinationPath(p string) {
	pathDestination = p
}

func SetStaticPath(p string) {
	pathStatic = p
}

func (pa *Pair) GetUUID() string {

	return pa.uuid
}

func (pa *Pair) GetTimeStamp() time.Time {

	return pa.timestamp
}

func (pa *Pair) Delete() error {

	if len(pa.first.name) > 0 {
		err := file.Remove(pathDestination + pa.first.name)
		if err != nil {
			return err
		}
	}

	if len(pa.second.name) > 0 {
		return file.Remove(pathDestination + pa.second.name)
	}

	return nil
}

func (pa *Pair) IsFirstUser(user string) bool {

	if user == pa.first.user {
		return true
	}

	return false
}

func (pa *Pair) IsSecondUser(user string) bool {

	if user == pa.second.user {
		return true
	}

	return false
}

func (pa *Pair) SetSecondUser(user string) error {

	if len(pa.second.user) > 0 {

		if pa.second.user != user {
			return errors.New("user is already set")
		}
	}

	pa.second.user = user

	return nil
}

func (pa *Pair) SetPicture(user, name string) error {

	if pa.first.user == user {

		if pa.first.state == UNLOCKED ||
			pa.first.state == VISIBLE {
			return nil
		}

		pa.first.name = name
		pa.first.state = LOCKED

		pa.timestamp = time.Now()

		return save(pa.first.state, pa.first.name)
	}

	if pa.second.state == UNLOCKED ||
		pa.second.state == VISIBLE {
		return nil
	}

	pa.second.user = user
	pa.second.name = name
	pa.second.state = LOCKED

	pa.timestamp = time.Now()

	return save(pa.second.state, pa.second.name)
}

func (pa *Pair) GetFirstPicture(user string) (string, error) {

	pa.timestamp = time.Now()

	if pa.first.state == MISSING {
		return "", nil
	}

	if pa.first.user != user &&
		pa.second.user != user {
		return "", errors.New("invalid user")
	}

	return pa.first.name, nil
}

func (pa *Pair) GetSecondPicture(user string) (string, error) {

	pa.timestamp = time.Now()

	if pa.second.state == MISSING {
		return "", nil
	}

	if pa.first.user != user &&
		pa.second.user != user {
		return "", errors.New("invalid user")
	}

	return pa.second.name, nil
}

func (pa *Pair) UnlockPicture(user string) error {

	if pa.first.user == user {

		if pa.second.state == LOCKED {

			pa.second.state = UNLOCKED

			if pa.first.state == UNLOCKED &&
				pa.second.state == UNLOCKED {

				pa.first.state = VISIBLE
				pa.second.state = VISIBLE

				err := save(pa.first.state, pa.first.name)
				if err != nil {
					return err
				}
			}

			return save(pa.second.state, pa.second.name)
		}
	} else {

		if pa.first.state == LOCKED {

			pa.first.state = UNLOCKED

			if pa.first.state == UNLOCKED &&
				pa.second.state == UNLOCKED {

				pa.first.state = VISIBLE
				pa.second.state = VISIBLE

				err := save(pa.second.state, pa.second.name)
				if err != nil {
					return err
				}
			}

			return save(pa.first.state, pa.first.name)
		}
	}

	return nil
}

func lock(pic image.Image, icon string) image.Image {

	// Create a blurred version of the image.
	img := imaging.Blur(pic, 10)

	ret := imaging.New(pic.Bounds().Max.X, pic.Bounds().Max.Y, color.NRGBA{0, 0, 0, 0})
	ret = imaging.Paste(ret, img, image.Pt(0, 0))

	imgLock, err := imaging.Open(pathStatic + icon)
	if err == nil {
		nrgbaLock := imaging.New(imgLock.Bounds().Max.X, imgLock.Bounds().Max.Y, color.NRGBA{0, 0, 0, 0})
		nrgbaLock = imaging.Paste(nrgbaLock, imgLock, image.Pt(0, 0))
		nrgbaLock = imaging.Resize(nrgbaLock, pic.Bounds().Max.X/6, 0, imaging.Lanczos)

		ret = imaging.PasteCenter(ret, nrgbaLock)
	}

	return ret
}

func save(state int, name string) error {

	img, err := imaging.Open(pathSource + name)
	if err != nil {
		return err
	}

	switch state {
	default:
		img = lock(img, "lock-solid.jpg")

	case UNLOCKED:
		img = lock(img, "lock-open-solid.jpg")

	case VISIBLE:
		//img = img
	}

	err = functions.CreateDir(pathDestination)
	if err != nil {
		return err
	}

	return imaging.Save(img, pathDestination+name)
}
