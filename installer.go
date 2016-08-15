// Lib to handle installing mac pkgs and apps

package macosutils

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path"
)

func InstallApp(appPath string) bool {
	files, _ := ioutil.ReadDir(appPath)
	result := false
	for _, f := range files {
		if path.Ext(f.Name()) == ".app" {
			args := []string{"--noqtn", path.Join(appPath, f.Name()), path.Join("/Applications/", f.Name())}
			_, err := exec.Command("/usr/bin/ditto", args...).Output()
			if err != nil {
				log.Fatal(err)
				result = false
			} else {
				result = true
			}
		}
	}
	return result
}
