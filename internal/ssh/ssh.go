package ssh

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/tergel-sama/kssh/internal/models"
)

func RunSSH(host models.HostConfig) {
	args := []string{}
	if host.Key != "" {
		args = append(args, "-i", expandPath(host.Key))
	}
	if host.Port != 0 {
		args = append(args, "-p", fmt.Sprint(host.Port))
	}
	target := fmt.Sprintf("%s@%s", host.User, host.Hostname)
	args = append(args, target)
	fmt.Printf("🔐 Connecting to %s...\n", target)
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("SSH failed: %v", err)
	}
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		return filepath.Join(usr.HomeDir, path[2:])
	}
	return path
}
