package scheduleconfig

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ScheduleConfig struct {
	Backups []BackupSchedule
}

type BackupSchedule struct {
	UUID    string
	Name    string
	Day     [7]bool
	PLC     string
	DB      []string
	Backup  time.Time
	checked time.Time
}

var mut sync.RWMutex

func (sc *ScheduleConfig) Add(bs BackupSchedule) error {

	mut.Lock()
	defer mut.Unlock()

	for i := range sc.Backups {

		if sc.Backups[i].UUID == bs.UUID {
			sc.Backups[i] = bs
			return nil
		}
	}

	sc.Backups = append(sc.Backups, bs)
	return nil
}

func (sc *ScheduleConfig) Get(uid string) (BackupSchedule, error) {

	mut.RLock()
	defer mut.RUnlock()

	for i := range sc.Backups {

		if sc.Backups[i].UUID == uid {
			return sc.Backups[i], nil
		}
	}

	return BackupSchedule{}, errors.New("id not found")
}

func (sc *ScheduleConfig) GetList() []BackupSchedule {

	mut.RLock()
	defer mut.RUnlock()

	return sc.Backups
}

func (sc *ScheduleConfig) Delete(uid string) error {

	mut.Lock()
	defer mut.Unlock()

	var n []BackupSchedule

	for i := range sc.Backups {

		if sc.Backups[i].UUID != uid {
			n = append(n, sc.Backups[i])
		}
	}

	sc.Backups = n

	return nil
}

func (sc *ScheduleConfig) SetBackup(uid string, t time.Time) error {

	mut.Lock()
	defer mut.Unlock()

	for i := range sc.Backups {

		if sc.Backups[i].UUID == uid {
			sc.Backups[i].Backup = t
			return nil
		}
	}

	return errors.New("id not found")
}

func (sc *ScheduleConfig) SetChecked(uid string, t time.Time) error {

	mut.Lock()
	defer mut.Unlock()

	for i := range sc.Backups {

		if sc.Backups[i].UUID == uid {
			sc.Backups[i].checked = t
			return nil
		}
	}

	return errors.New("id not found")
}

func NewBackupSchedule() BackupSchedule {

	var ret BackupSchedule
	ret.UUID = uuid.Must(uuid.NewRandom()).String()
	ret.Name = ret.UUID

	for i := range ret.Day {
		ret.Day[i] = true
	}

	return ret
}

func (bs *BackupSchedule) AddDB(uid string) error {

	mut.Lock()
	defer mut.Unlock()

	for i := range bs.DB {

		if bs.DB[i] == uid {
			return nil
		}
	}

	bs.DB = append(bs.DB, uid)
	return nil
}

func (bs *BackupSchedule) SetChecked(t time.Time) {

	mut.Lock()
	defer mut.Unlock()

	bs.checked = t
}

func (bs *BackupSchedule) GetChecked() time.Time {

	mut.RLock()
	defer mut.RUnlock()

	return bs.checked
}
