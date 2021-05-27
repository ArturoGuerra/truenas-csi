package mounter

/*
Was going to use k8s.io/pkg/util/mount but it's missing in the 1.17 release don't know if that is intended,
Hopefully I understood mount propagation correctly lul
*/

import (

	// Filesystem stuff
	"os"

	"github.com/sirupsen/logrus"
	mountutils "k8s.io/mount-utils"
	utilsexec "k8s.io/utils/exec"
)

type (
	Mounter interface {
		//BindMount(source, target, fstype, opts string) error
		//FormatAndMount(source, target, fstype, opts string) error
		//Unmount(target string) error
		PathExists(target string) (bool, error)
		MakeDir(target string) error
		mountutils.Interface
		utilsexec.Interface
	}

	mounter struct {
		Logger *logrus.Logger
		mountutils.SafeFormatAndMount
		utilsexec.Interface
	}
)

func New(logger *logrus.Logger) Mounter {
	mount := mountutils.SafeFormatAndMount{
		Interface: mountutils.New(""),
		Exec:      utilsexec.New(),
	}
	return &mounter{
		logger,
		mount,
		mount.Exec,
	}
}

/* Mounts directory NOTE: Is used inside FormatAndMount */
/*func (m *mounter) BindMount(source, target, fstype, opts string) error {
	return gofsutil.BindMount(ctx.Background(), source, target, fstype, opts)
}

/* Formats and/or Mounts a device to directory
func (m *mounter) FormatAndMount(source, target, fstype string, opts string) error {
	return gofsutil.FormatAndMount(ctx.Background(), source, target, fstype, opts)
}

 Unmounts directory
func (m *mounter) Unmount(target string) error {
	return gofsutil.Unmount(ctx.Background(), target)
}*/

/* Check if directory exists */
func (m *mounter) PathExists(target string) (bool, error) {
	_, err := os.Stat(target)
	if os.IsNotExist(err) {
		m.Logger.Infof("Target: (%s) does not exist", target)
		return false, nil
	} else {
		m.Logger.Infof("Target: (%s) exists", target)
		return true, err
	}
}

func (m *mounter) MakeDir(target string) error {
	return os.Mkdir(target, 0775)
}
