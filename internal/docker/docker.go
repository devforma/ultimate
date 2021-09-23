package docker

import "os/exec"

func StartContainer() {
	exec.Command("docker", "start")
}
