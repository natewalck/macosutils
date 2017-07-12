// Lib to handle any disk related operations, including mounting .dmg files, etc

package macosutils

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/groob/plist"
)

// DMG reprents a macOS DMG image
type DMG struct {
	dmgpath    string
	MountPoint string
	Pkgs       []string
	Apps       []string
	logger     *log.Logger
}

// SystemEntities contains an array of volumes
type SystemEntities struct {
	Disks []DiskEntries `plist:"system-entities"`
}

// DiskEntries contains all mount points for the dmg
type DiskEntries struct {
	Disk *string `plist:"mount-point"`
}

// DMGOption allows the passing of dmg mounting options that differ from default
type DMGOption func(*DMG)

// WithLogger allows you to pass a custom logger to the NewDMG function
func WithLogger(logger *log.Logger) DMGOption {
	return func(d *DMG) {
		d.logger = logger
	}
}

// NewDMG will create a new DMG object, with various utility functions
func NewDMG(path string, opts ...DMGOption) *DMG {
	d := &DMG{
		dmgpath: path,
		logger:  log.New(os.Stderr, "test", 1),
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Mount the DMG
func (d *DMG) Mount() error {
	log.Printf("Mounting dmg located at %v\n", d.dmgpath)
	args := []string{"attach", d.dmgpath, "-mountRandom", "/tmp", "-nobrowse", "-plist"}
	out, err := exec.Command("/usr/bin/hdiutil", args...).Output()
	if err != nil {
		log.Printf("Failed to mount dmg with error: %v", err)
	}

	var mountinfo SystemEntities
	err = plist.Unmarshal(out, &mountinfo)

	d.MountPoint = MountPoint(mountinfo)
	log.Printf("DMG mounted at %v\n", d.MountPoint)
	return err
}

// Unmount the DMG
func (d *DMG) Unmount(dmgPath string) error {
	args := []string{"detach", dmgPath}
	_, err := exec.Command("/usr/bin/hdiutil", args...).Output()
	return err
}

// GetInstallables will show all valid installer types inside the dmg
func (d *DMG) GetInstallables() {
	files, _ := ioutil.ReadDir(d.MountPoint)
	d.Apps = []string{}
	d.Pkgs = []string{}
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

// MountPoint returns the filepath where the dmg is mounted
func MountPoint(mountinfo SystemEntities) string {
	for _, v := range mountinfo.Disks {
		if v.Disk != nil {
			return *v.Disk
		}
	}
	return ""
}
