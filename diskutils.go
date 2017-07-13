// Lib to handle any disk related operations, including mounting .dmg files, etc

package macosutils

import (
	"bytes"
	"fmt"
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
	ImageInfo  BackingStoreInfo
	Pkgs       []string
	Apps       []string
	logger     *log.Logger
}

// BackingStoreInfo Stores all info about the disk image
type BackingStoreInfo struct {
	ChecksumType      string              `plist:"Checksum Type"`
	ChecksumValue     string              `plist:"Checksum Value"`
	ClassName         string              `plist:"Checksum Name"`
	Format            string              `plist:"Format"`
	FormatDescription string              `plist:"Format Description"`
	Properties        DiskImageProperties `plist:"Properties"`
}

// DiskImageProperties represents the properties of a disk image
type DiskImageProperties struct {
	Checksummed              bool `plist:"Checksummed"`
	Compressed               bool `plist:"Compressed"`
	Encrypted                bool `plist:"Encrypted"`
	KernelCompatible         bool `plist:"Kernel Compatible"`
	Partitioned              bool `plist:"Partitioned"`
	SoftwareLicenseAgreement bool `plist:"Software License Agreement"`
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
func NewDMG(path string, opts ...DMGOption) (*DMG, error) {
	d := &DMG{
		dmgpath: path,
		logger:  log.New(os.Stderr, "test", 1),
	}
	for _, opt := range opts {
		opt(d)
	}

	args := []string{"imageinfo", d.dmgpath, "-plist"}
	out, err := exec.Command("/usr/bin/hdiutil", args...).Output()
	data := bytes.Replace(out, []byte(`<integer>-1</integer>`), []byte(`<string>-1</string>`), -1)
	if err != nil {
		return nil, fmt.Errorf("Failed to get the info from dmg : %s", err)
	}
	var diskInfo BackingStoreInfo
	err = plist.Unmarshal(data, &diskInfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to read disk info: %s", err)
	}
	d.ImageInfo = diskInfo
	return d, nil
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
	if err != nil {
		log.Printf("Failed to read mount info: %v", err)
	}

	d.MountPoint = MountPoint(mountinfo)
	log.Printf("DMG mounted at %v\n", d.MountPoint)
	return err
}

// Unmount the DMG
func (d *DMG) Unmount(dmgPath string) error {
	args := []string{"detach", dmgPath}
	_, err := exec.Command("/usr/bin/hdiutil", args...).Output()
	if err != nil {
		log.Printf("Failed to unmount dmg: %v", err)
	}
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

// HasSLA returns true if the DMG has an SLA
func (d *DMG) HasSLA() bool {
	return d.ImageInfo.Properties.SoftwareLicenseAgreement
}

// IsWritable returns true if the DMG is in a writable format
func (d *DMG) IsWritable() bool {
	writableFormats := []string{"UDSB", "UDSP", "UDRW", "RdWr"}
	for _, format := range writableFormats {
		if d.ImageInfo.Format == format {
			return true
		}
	}
	return false
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
