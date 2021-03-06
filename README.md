
# gouniversal
Start writing your **new [Go](https://golang.org) Application** with gouniversal
and use a prepared user/group administration, navigation and alarm system.

HTML front end based on **Bootstrap v4**.

### Files and Folders
```
gouniversal/
├── build/ (build options and scripts)
├── data/ (Folder to ship with binary)
│   ├── config/
│   │   ├── openespm/
│   │   ├── group
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
├── module/ (add your Applications and Modules here)
│   ├── module.go (register your Application/Module)
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

Device Firmware: [github.com/dekoch/GOopenESPM](https://github.com/dekoch/GOopenESPM)

## Fileshare
Up/Download, Manage and Share your Files.

## Homepage
Auto generate Menu entries for your html Files.
Demo: [dekoch.net](https://dekoch.net)

## MeshFileSync
Sync files over a mesh network

## MediaDownloader
Download Link Creator

## GasPrice
Track Gas Prices

## Console
Live Console output with HTML5 SSE.

## LogViewer
View gouniversal log files.

# Modules
## Mesh
decentralised connectivity control

## HeatingMath
goroutine demo

## Mark
Benchmark module

# used packages

[github.com/adrianmo/go-nmea](https://github.com/adrianmo/go-nmea)

[github.com/asaskevich/govalidator](https://github.com/asaskevich/govalidator)

[github.com/goburrow/modbus](https://github.com/goburrow/modbus)

[github.com/google/uuid](https://github.com/google/uuid)

[github.com/gorilla/securecookie](https://github.com/gorilla/securecookie)

[github.com/gorilla/sessions](https://github.com/gorilla/sessions)

[github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid](https://github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid)

[github.com/tarm/serial](https://github.com/tarm/serial)

[github.com/twpayne/go-gpx](https://github.com/twpayne/go-gpx)