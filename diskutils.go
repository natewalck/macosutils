// Lib to handle any disk related operations, including mounting .dmg files, etc

package macosutils

import (
	"fmt"
	"github.com/groob/plist"
	"log"
	"os/exec"
)

type DMG struct {
	MountPoint string
}

type SystemEntities struct {
	Disks []DiskEntries `plist:"system-entities"`
}

type DiskEntries struct {
	Disk *string `plist:"mount-point"`
}

func (d *DMG) Mount(dmgPath string) {
	fmt.Printf("Mounting dmg located at %v\n", dmgPath)
	args := []string{"attach", dmgPath, "-mountRandom", "/tmp", "-nobrowse", "-plist"}
	out, err := exec.Command("/usr/bin/hdiutil", args...).Output()
	if err != nil {
		log.Fatal(err)
	}

	var mountinfo SystemEntities
	err = plist.Unmarshal(out, &mountinfo)
	if err != nil {
		log.Fatal(err)
	}

	d.MountPoint = MountPoint(mountinfo)
	fmt.Printf("DMG mounted at %v\n", d.MountPoint)
}

func (d *DMG) Unmount(dmgPath string) bool {
	args := []string{"detach", dmgPath}
	_, err := exec.Command("/usr/bin/hdiutil", args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func MountPoint(mountinfo SystemEntities) string {
	for _, v := range mountinfo.Disks {
		if v.Disk != nil {
			return *v.Disk
		}
	}
	return ""
}
