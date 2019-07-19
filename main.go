// +build windows

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// cirruslauncher.exe -username czmadmin -password czmAdmin2008 -patientId 1234 -issuerOfPatientId issuerABC
// czmcirrus://server/appr?username=czmadmin&password=czmAdmin2008&patientId=1234&issuerOfPatientId=issuerABC
// czmcirrus://server/app?username=$$a8ea4f8bd53a4667&password=$$a8ea4fabd53a466712ab4a07&patientId=1234&issuerOfPatientId=issuerABC
//
// Windows compile:
// Docu: https://github.com/josephspurrier/goversioninfo
// go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
// go generate
// go install -ldflags -H=windowsgui
//
// MacOS compile:
// Docu: https://medium.com/@mattholt/packaging-a-go-application-for-macos-f7084b00f6b5

const (
	secret = "JHjh()z&)/hlLZ(jn(jnjnJHJ68JoUu7" // 32 Bytes
)

var (
	cirruspath string
	logFile    *os.File
)

func isWindows() bool {
	return strings.ToLower(runtime.GOOS) == "windows"
}

func decrypt(key, text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	ciphertext, err := hex.DecodeString(text)
	if err != nil {
		panic(err)
	}
	cfb := cipher.NewCFBEncrypter(block, []byte(key[:16]))
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)

	return string(plaintext)
}

// fileExists checks if a file/dir exists
func fileExists(filename string) bool {
	var b bool
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		b = false
	}

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		b = true
	}

	return b
}

func currentPath() string {
	return filepath.Dir(os.Args[0])
}

func info(t string, arg ...interface{}) {
	if len(arg) > 0 {
		t = fmt.Sprintf(t, arg...)
	}

	t = time.Now().Format("2006-01-02 15:04:05.000") + " INFO  " + t + "\n"
	fmt.Printf(t)

	if logFile != nil {
		logFile.WriteString(t)
	}
}

func fatal(e error) {
	t := e.Error()

	t = time.Now().Format("2006-01-02 15:04:05.000") + " ERROR " + t + "\n"
	fmt.Printf(t)

	if logFile != nil {
		logFile.WriteString(t)
	}
}

// initLog initializes the logging to console and log file
func initLog() {
	filename, err := os.Executable()
	if err != nil {
		filename = os.Args[0]
	}

	ext := filepath.Ext(filename)

	if len(ext) > 0 {
		filename = string(filename[:len(filename)-len(ext)])
	}

	filename += ".log"

	filename = filepath.Join(currentPath(), filepath.Base(filename))

	logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		usr, err := user.Current()
		if err != nil {
			fatal(err)
		}

		filename = filepath.Join(usr.HomeDir, filepath.Base(filename))

		logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			fatal(err)
		}
	}

	logFile.Close()

	if fileExists(filename) {
		fi, _ := os.Stat(filename)

		if fi.Size() > int64(10*1024*1024) {
			err := os.Remove(filename)
			if err != nil {
				fatal(err)
			}
		}
	}

	logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fatal(err)
	}

	info("---------- new invocation")
}

// initCirruspath initializes the cirruspath variable with the final "viewer" application filename
func initCirruspath() {
	info("")

	cirruspath = currentPath() + string(filepath.Separator) + "Czm.Sdoct.AnalysisEngine.exe"

	if !fileExists(cirruspath) {
		fatal(fmt.Errorf("no viewer executable found: %s", cirruspath))
	}
}

func main() {
	initLog()
	initCirruspath()

	defer func() {
		logFile.Close()
	}()

	appName, err := os.Executable()
	if err != nil {
		fatal(err)
	}

	appName = filepath.Base(appName)

	extLen := len(filepath.Ext(appName))
	if extLen > 0 {
		appName = appName[:len(appName)-extLen]
	}

	fmt.Printf("\n")
	fmt.Printf("%s - Launcher for Cirrus\n", strings.ToUpper(appName))
	fmt.Printf("\n")

	info("found cirrus: %s", cirruspath)

	usr, err := user.Current()
	if err != nil {
		fatal(err)
	}

	info("user home dir: %s", usr.HomeDir)

	var cmdLine string

	if len(os.Args) > 1 {
		cmdLine = strings.Join(os.Args[1:], " ")
	}

	if strings.HasPrefix(cmdLine, "czmcirrus://") {
		info("detected URL protocol launcher parameter: %s", cmdLine)

		u, err := url.Parse(cmdLine)
		if err != nil {
			fatal(err)
		}

		args := strings.Split(u.RawQuery, "&")
		cmdLine = ""

		for _, e := range args {
			e, err := url.QueryUnescape(e)
			if err != nil {
				fatal(err)
			}

			e = strings.Replace(e, "=", ":", 1)

			if len(cmdLine) > 0 {
				cmdLine += " "
			}

			cmdLine += e
		}

		info("converted args from url: %s", cmdLine)
	} else {
		cmdLine = ""

		for i := 1; i < len(os.Args); i++ {
			txt := ""
			if strings.Index(os.Args[i], ":") == -1 {
				txt = os.Args[i] + ":"
				if i+1 < len(os.Args) {
					txt += os.Args[i+1]
					i++
				}
			} else {
				txt = os.Args[i]
			}

			if strings.HasPrefix(txt, "-") {
				txt = txt[1:]
			}

			if len(cmdLine) > 0 {
				cmdLine += " "
			}

			cmdLine += txt
		}
	}

	info("cirrus parameter: %s", cmdLine)

	var args []string

	args = append(args, cirruspath)
	args = append(args, strings.Split(cmdLine, " ")...)

	pargs := make([]string, len(args))

	copy(pargs, args)

	for i := range pargs {
		if isWindows() {
			pargs[i] = "\"" + pargs[i] + "\""
		} else {
			pargs[i] = "'" + pargs[i] + "'"
		}
	}

	info("exec command: %s %s", pargs[0], strings.Join(pargs[1:], " "))

	for i := 1; i < len(args); i++ {
		arg := args[i]

		if strings.HasPrefix(args[i], "-") {
			args[i] = args[i][1:]
		}

		if (i+1) < len(args) && (strings.HasPrefix(arg, "username") || strings.HasPrefix(arg, "password")) {
			i++

			arg = args[i]

			if strings.HasPrefix(arg, "$$") {
				args[i] = decrypt(secret, arg[2:])
			}
		}
	}

	cmd := exec.Command(args[0], args[1:]...)

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}
}
