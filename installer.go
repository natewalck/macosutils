// Lib to handle installing mac pkgs and apps

package macosutils

import (
	"io/ioutil"
	"os/exec"
	"path"
)

func InstallApp(appPath string) error {
	files, _ := ioutil.ReadDir(appPath)
	for _, f := range files {
		if path.Ext(f.Name()) == ".app" {
			args := []string{"--noqtn", path.Join(appPath, f.Name()), path.Join("/Applications/", f.Name())}
			_, err := exec.Command("/usr/bin/ditto", args...).Output()
			return err
		}
	}
	return nil
}
