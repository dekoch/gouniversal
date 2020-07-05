package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type BackupList struct {
	Menu    string
	Title   string
	Name    string
	Backup  string
	Restore string
}

type Backup struct {
	BackupTitle          string
	RestoreTitle         string
	Name                 string
	Address              string
	Rack                 string
	Slot                 string
	MaxBackupCnt         string
	UUID                 string
	BackupPLC            string
	BackupDB             string
	RestorePLC           string
	RestoreDB            string
	DBNo                 string
	DBByteLength         string
	SavedToDatabase      string
	RestoredFromDatabase string
}

type PLCList struct {
	Menu   string
	Title  string
	Name   string
	Edit   string
	AddPLC string
}

type PLCEdit struct {
	Title        string
	Apply        string
	Delete       string
	Name         string
	Address      string
	Rack         string
	Slot         string
	MaxBackupCnt string
	UUID         string
	DBNo         string
	DBByteLength string
	AddDB        string
}

type ScheduleList struct {
	Menu        string
	Title       string
	Name        string
	PLC         string
	Edit        string
	AddSchedule string
}

type ScheduleEdit struct {
	Title     string
	Apply     string
	Delete    string
	Name      string
	PLC       string
	ActiveDay string
	Day
	DBNo string
}

type Day struct {
	Monday    string
	Tuesday   string
	Wednesday string
	Thursday  string
	Friday    string
	Saturday  string
	Sunday    string
}

type Alert struct {
	Success string
	Info    string
	Warning string
	Error   string
}

type LangFile struct {
	Header       config.FileHeader
	BackupList   BackupList
	Backup       Backup
	PLCList      PLCList
	PLCEdit      PLCEdit
	ScheduleList ScheduleList
	ScheduleEdit ScheduleEdit
	Alert        Alert
}

func DefaultEn() LangFile {

	var l LangFile

	l.Header = config.BuildHeader("en", "LangS7Backup", 1.0, "Language File")

	l.BackupList.Menu = "S7Backup"
	l.BackupList.Title = "Backup/Restore"
	l.BackupList.Name = "Name"
	l.BackupList.Backup = "Backup"
	l.BackupList.Restore = "Restore"

	l.Backup.BackupTitle = "Backup"
	l.Backup.RestoreTitle = "Restore"
	l.Backup.Name = "Name"
	l.Backup.Address = "Address"
	l.Backup.Rack = "Rack"
	l.Backup.Slot = "Slot"
	l.Backup.MaxBackupCnt = "Max. Backup Cnt"
	l.Backup.UUID = "UUID"
	l.Backup.BackupPLC = "Backup PLC"
	l.Backup.BackupDB = "Backup DB"
	l.Backup.RestorePLC = "Restore PLC"
	l.Backup.RestoreDB = "Restore DB"
	l.Backup.DBNo = "DB No."
	l.Backup.DBByteLength = "DB Byte length"
	l.Backup.SavedToDatabase = "saved to database"
	l.Backup.RestoredFromDatabase = "restored from database"

	l.PLCList.Menu = "S7Backup"
	l.PLCList.Title = "Configure"
	l.PLCList.Name = "Name"
	l.PLCList.Edit = "Edit"
	l.PLCList.AddPLC = "Add PLC"

	l.PLCEdit.Title = "Edit PLC"
	l.PLCEdit.Apply = "Apply"
	l.PLCEdit.Delete = "Delete"
	l.PLCEdit.Name = "Name"
	l.PLCEdit.Address = "Address"
	l.PLCEdit.Rack = "Rack"
	l.PLCEdit.Slot = "Slot"
	l.PLCEdit.MaxBackupCnt = "Max. Backup Cnt"
	l.PLCEdit.UUID = "UUID"
	l.PLCEdit.DBNo = "DB No."
	l.PLCEdit.DBByteLength = "DB Byte length"
	l.PLCEdit.AddDB = "Add DB"

	l.ScheduleList.Menu = "S7Backup"
	l.ScheduleList.Title = "Schedule"
	l.ScheduleList.Name = "Name"
	l.ScheduleList.PLC = "PLC"
	l.ScheduleList.Edit = "Edit"
	l.ScheduleList.AddSchedule = "Add Schedule"

	l.ScheduleEdit.Title = "Edit Schedule"
	l.ScheduleEdit.Apply = "Apply"
	l.ScheduleEdit.Delete = "Delete"
	l.ScheduleEdit.Name = "Name"
	l.ScheduleEdit.PLC = "PLC"
	l.ScheduleEdit.ActiveDay = "Active Day"
	l.ScheduleEdit.Monday = "Monday"
	l.ScheduleEdit.Tuesday = "Tuesday"
	l.ScheduleEdit.Wednesday = "Wednesday"
	l.ScheduleEdit.Thursday = "Thursday"
	l.ScheduleEdit.Friday = "Friday"
	l.ScheduleEdit.Saturday = "Saturday"
	l.ScheduleEdit.Sunday = "Sunday"
	l.ScheduleEdit.DBNo = "DB No."

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
