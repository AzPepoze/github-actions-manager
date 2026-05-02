package service

import (
	"os/exec"
	"strings"
)

func ServiceCmd(installPath, action string) *exec.Cmd {
	cmd := exec.Command("sudo", "./svc.sh", action)
	cmd.Dir = installPath
	return cmd
}

func runServiceCmd(installPath, action string) error {
	return ServiceCmd(installPath, action).Run()
}

func InstallService(installPath string) error {
	return runServiceCmd(installPath, "install")
}

func UninstallService(installPath string) error {
	return runServiceCmd(installPath, "uninstall")
}

func StartService(installPath string) error {
	return runServiceCmd(installPath, "start")
}

func StopService(installPath string) error {
	return runServiceCmd(installPath, "stop")
}

func GetServiceStatus(installPath string) (installed bool, running bool, err error) {
	cmd := exec.Command("sudo", "./svc.sh", "status")
	cmd.Dir = installPath
	out, _ := cmd.CombinedOutput()
	output := string(out)

	if strings.Contains(output, "active (running)") {
		return true, true, nil
	}
	if strings.Contains(output, "inactive (dead)") || strings.Contains(output, "Stopped") {
		return true, false, nil
	}
	if strings.Contains(output, "not installed") || strings.Contains(output, "No such file") {
		return false, false, nil
	}

	return strings.Contains(output, "status"), false, nil
}
