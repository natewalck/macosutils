// Lib to handle any disk related operations, including mounting .dmg files, etc

package macosutils

import (
	"github.com/groob/plist"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"strings"
)

type DMG struct {
	MountPoint string
	Pkgs       []string
	Apps       []string
}

type SystemEntities struct {
	Disks []DiskEntries `plist:"system-entities"`
}

type DiskEntries struct {
	Disk *string `plist:"mount-point"`
}

func (d *DMG) Mount(dmgPath string) error {
	log.Printf("Mounting dmg located at %v\n", dmgPath)
	args := []string{"attach", dmgPath, "-mountRandom", "/tmp", "-nobrowse", "-plist"}
	out, err := exec.Command("/usr/bin/hdiutil", args...).Output()
	if err != nil {
		log.Fatal(err)
	}

	var mountinfo SystemEntities
	err = plist.Unmarshal(out, &mountinfo)

	d.MountPoint = MountPoint(mountinfo)
	log.Printf("DMG mounted at %v\n", d.MountPoint)
	return err
}

func (d *DMG) Unmount(dmgPath string) error {
	args := []string{"detach", dmgPath}
	_, err := exec.Command("/usr/bin/hdiutil", args...).Output()
	return err
}

func (d *DMG) GetInstallables() {
	files, _ := ioutil.ReadDir(d.MountPoint)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		} else if path.Ext(f.Name()) == ".app" {
			d.Apps = append(d.Apps, f.Name())
		} else if path.Ext(f.Name()) == ".pkg" {
			d.Pkgs = append(d.Pkgs, f.Name())
		}
	}
}

func MountPoint(mountinfo SystemEntities) string {
	for _, v := range mountinfo.Disks {
		if v.Disk != nil {
			return *v.Disk
		}
	}
	return ""
}
