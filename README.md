
# gouniversal
Start writing your **new Go Application** with gouniversal
and use a prepared user/group administration, navigation and alarm system.

HTML front end based on **Bootstrap v4**.

### Files and Folders
```
gouniversal/
├── data/ (Folder to ship with binary)
│   ├── config/
│   │   ├── group
│   │   ├── openespm/
│   │   ├── program
│   │   ├── ui
│   │   └── user
│   ├── lang/
│   │   ├── openespm/
│   │   │   └── en
│   │   └── program/
│   │       └── en
│   └── ui/ (HTML Templates/Stylesheets/Images/...)
│       ├── openespm/ (add a new folder for every Application/Module)
│       │   └── 1.0/ (use a subfolder for every version)
│       ├── program/
│       │   └── 1.0/
│       └── static/
│           └── 1.0/
├── gouniversal.go
├── modules/ (add your Applications and Modules here)
│   ├── modules.go (register your Application/Module)
│   └── openespm/ (first Application/Template/Test)
│       └── openespm.go
├── program/ (Just!!! user/group administration, navigation, alarms.
│              Please add other things as module)
├── README.md
└── shared/ (shared functions for program and modules)
```
# Applications
## openESPM
open-source Manager for IoT Devices
 - SimpleSwitch v1.0

Device Firmware: [https://github.com/dekoch/GOopenESPM](https://github.com/dekoch/GOopenESPM)

## Fileshare
Up/Download, Manage and Share your Files.

## Homepage
Auto generate Menu entries for your html Files.
Demo: [https://dekoch.net](https://dekoch.net)

## Console
Live Console output with HTML5 SSE.