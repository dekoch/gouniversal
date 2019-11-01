package userconfig

type UserConfig struct {
	User    string
	InstaID []string
}

func (hc *UserConfig) LoadDefaults(user string) {

	hc.User = user
}

func (hc *UserConfig) AddID(instaid string) {

	for i := range hc.InstaID {

		if hc.InstaID[i] == instaid {
			return
		}
	}

	hc.InstaID = append(hc.InstaID, instaid)
}

func (hc *UserConfig) GetIDList() []string {

	return hc.InstaID
}
