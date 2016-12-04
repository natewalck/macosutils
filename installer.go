// Lib to handle installing mac pkgs and apps

package macosutils

import (
	"log"
	"os/exec"
	"path"
)

func InstallApp(appPath string, appName string) error {
	args := []string{"--noqtn", path.Join(appPath, appName), path.Join("/Applications/", appName)}
	_, err := exec.Command("/usr/bin/ditto", args...).Output()
	if err != nil {
		log.Printf("Failed to install:", appName)
	}
	return err
}

func InstallPkg(pkgPath string, pkgName string) error {
	args := []string{"-pkg", path.Join(pkgPath, pkgName), "-tgt", "/"}
	_, err := exec.Command("/usr/sbin/installer", args...).Output()
	if err != nil {
		log.Printf("Failed to install:", pkgName)
	}
	return err
}
