// Lib to handle any disk related operations, including mounting .dmg files, etc

package macosutils

import (
	"github.com/groob/plist"
	"log"
	"os/exec"
)

type SystemEntities struct {
	Disks []DiskEntries `plist:"system-entities"`
}

type DiskEntries struct {
	Disk *string `plist:"mount-point"`
}

func MountDmg(dmgPath string) string {
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

	mount := MountPoint(mountinfo)
	return mount
}

func UnmountDmg(dmgPath string) bool {
	args := []string{"detach", dmgPath}
	_, err := exec.Command("/usr/bin/hdiutil", args...).Output()
	if err != nil {
		log.Fatal(err)
		return false
	} else {
		return true
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
