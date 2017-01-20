package main

// #cgo LDFLAGS: -lX11 -lasound
// #include <X11/Xlib.h>
// #include "getvol.h"
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"
)

var dpy = C.XOpenDisplay(nil)

func getVolumePerc() int {
	return int(C.get_volume_perc())
}

func getBatteryPercentage(path string) (perc int, err error) {
	now := newFileErrReader()
	now.read(path + "/energy_now")
	now.read(path + "/charge_now")
	if now.err != nil {
		perc = -1
		return
	}

	full := newFileErrReader()
	full.read(path + "/energy_full")
	full.read(path + "/charge_full")
	if full.err != nil {
		perc = -1
		return
	}

	var enow, efull int
	fmt.Sscanf(now.String(), "%d", &enow)
	fmt.Sscanf(full.String(), "%d", &efull)
	perc = enow * 100 / efull
	return
}

func getBatteryStatus(path string) (string, error) {
	status := newFileErrReader()
	status.read(path + "/status")
	if status.err != nil {
		return "", status.err
	}

	switch strings.TrimSpace(status.String()) {
	case "Charging":
		return "+", nil
	case "Discharging":
		return "-", nil
	default:
		return "", nil
	}
}

func getLoadAverage(file string) (lavg string, err error) {
	loadavg, err := ioutil.ReadFile(file)
	if err != nil {
		return "Couldn't read loadavg", err
	}
	lavg = strings.Join(strings.Fields(string(loadavg))[:3], " ")
	return
}

func setStatus(s *C.char) {
	C.XStoreName(dpy, C.XDefaultRootWindow(dpy), s)
	C.XSync(dpy, 1)
}

func nowPlaying(addr string) (np string, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		np = "Couldn't connect to mpd."
		return
	}
	defer conn.Close()
	reply := make([]byte, 512)
	conn.Read(reply) // The mpd OK has to be read before we can actually do things.

	message := "status\n"
	conn.Write([]byte(message))
	conn.Read(reply)
	r := string(reply)
	arr := strings.Split(string(r), "\n")
	if arr[8] != "state: play" { //arr[8] is the state according to the mpd documentation
		status := strings.SplitN(arr[8], ": ", 2)
		np = fmt.Sprintf("mpd - [%s]", status[1]) //status[1] should now be stopped or paused
		return
	}

	message = "currentsong\n"
	conn.Write([]byte(message))
	conn.Read(reply)
	r = string(reply)
	arr = strings.Split(string(r), "\n")
	if len(arr) > 5 {
		var artist, title string
		for _, info := range arr {
			field := strings.SplitN(info, ":", 2)
			switch field[0] {
			case "Artist":
				artist = strings.TrimSpace(field[1])
			case "Title":
				title = strings.TrimSpace(field[1])
			default:
				//do nothing with the field
			}
		}
		np = artist + " - " + title
		return
	} else { //This is a nonfatal error.
		np = "Playlist is empty."
		return
	}
}

func formatStatus(format string, args ...interface{}) *C.char {
	status := fmt.Sprintf(format, args...)
	return C.CString(status)
}

func main() {
	if dpy == nil {
		log.Fatal("Can't open display")
	}
	battery := "/sys/class/power_supply/BAT0"

	for {
		t := time.Now().Format("Mon 02 15:04")
		b, err := getBatteryPercentage(battery)
		if err != nil {
			log.Println(err)
		}
		s, err := getBatteryStatus(battery)
		if err != nil {
			log.Println(err)
		}
		l, err := getLoadAverage("/proc/loadavg")
		if err != nil {
			log.Println(err)
		}
		m, err := nowPlaying("localhost:6600")
		if err != nil {
			log.Println(err)
		}
		vol := getVolumePerc()
		status := formatStatus("%s :: %d%% :: %s :: %s :: %s%d%%", m, vol, l, t, s, b)
		setStatus(status)
		time.Sleep(time.Second)
	}
}
