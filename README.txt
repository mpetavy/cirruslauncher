CirrusLauncher.exe -username czmadmin -password czmAdmin2008 -patientId 1234 -issuerOfPatientId issuerABC

Windows compile:

Install GO (assume C:\GO)
Set env variable GOROOT to the root directory of your GO installation (assume C:\GO)
Set env variable GOPATH to the root directory of your GO development workspace directory (assume D:\GODEV)

Docu: https:github.com/josephspurrier/goversioninfo
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go generate
CD to the ForumLauncher/main.go directory
go build -ldflags -H=windowsgui --> generate ForumLauncher.exe in the current directory 

MacOS compile:
Docu: https:medium.com/@mattholt/packaging-a-go-application-for-macos-f7084b00f6b5

