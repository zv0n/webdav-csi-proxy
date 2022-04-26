package webdav

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/zv0n/webdav-proxy/configuration"
	"k8s.io/utils/mount"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MountInput struct {
	URL        string
	Dir        string
	User       string
	Password   string
	ConfigName string
	TargetPath string
	UID        string
	GID        string
}

func Umount(targetPath string, configName string) error {
	notMnt, err := mount.New("").IsLikelyNotMountPoint(targetPath)

	if err != nil {
		if os.IsNotExist(err) {
			return status.Error(codes.NotFound, "Targetpath not found")
		} else {
			return status.Error(codes.Internal, err.Error())
		}
	}
	if notMnt {
		//return status.Error(codes.NotFound, "Volume not mounted")
		return nil
	}
	rcloneCmd := "rclone"
	configRemoveArgs := []string{}

	configRemoveArgs = append(
		configRemoveArgs,
		"config",
		"delete",
		configName,
	)

	err = exec.Command(rcloneCmd, configRemoveArgs...).Run()
	if err != nil {
		return fmt.Errorf("Failed to remove webdav configuration: %v, cmd: '%s %s'", err, rcloneCmd, configRemoveArgs)
	}

	err = mount.New("").Unmount(targetPath)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func Mount(input MountInput, config *configuration.Configuration) (err error) {
	mountCmd := "rclone"
	passwordArgs := []string{}
	configArgs := []string{}
	mountArgs := []string{}

	passwordArgs = append(
		passwordArgs,
		"obscure",
		input.Password,
	)
	obscuredPassword, err := exec.Command(mountCmd, passwordArgs...).Output()
	if err != nil {
		return fmt.Errorf("Failed to obscure password, aborting: %v", err)
	}

	configArgs = append(
		configArgs,
		"config",
		"create",
		input.ConfigName,
		"webdav",
		"url",
		input.URL,
		"vendor",
		"other",
		"user",
		input.User,
		"pass",
		string(obscuredPassword),
	)

	err = exec.Command(mountCmd, configArgs...).Run()
	if err != nil {
		return fmt.Errorf("Failed to create a webdav configuration: %v, cmd: '%s %s'", err, mountCmd, configArgs)
	}

	// create target, os.Mkdirall is noop if it exists
	err = os.MkdirAll(input.TargetPath, 0750)
	if err != nil {
		return err
	}

	mountArgs = append(
		mountArgs,
		"mount",
		input.ConfigName+":"+input.Dir,
		input.TargetPath,
		"--uid",
		input.UID,
		"--gid",
		input.GID,
		"--allow-other",
		"--daemon",
	)

	fmt.Printf("Command: %s %s", mountCmd, mountArgs)
	for _, x := range mountArgs {
		fmt.Printf("Arg: %s", x)
	}

	cmd := exec.Command(mountCmd, mountArgs...)
	// start command in its own process group so we can kill the parent without killing the mounts
	// this is POSIX specific, windows needs a different SysProcAttr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("mounting failed: %v cmd: '%s %s'",
			err, mountCmd, strings.Join(mountArgs, " "))
	}

	return nil
}
