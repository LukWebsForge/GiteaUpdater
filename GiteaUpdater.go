package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {

	fmt.Println("The Gitea Updater")

	version := readVersion()

	if len(version) < 2 || len(version) > 10 {
		fmt.Println("Your version doesn't seem to match :(")
		return
	}

	// Stopping Gitea service via systemctl
	runCommand("systemctl", "stop", "gitea")

	// Downloading the Gitea executable
	path, err := DownloadGitea(version)
	if err != nil {
		fmt.Println("Maybe your wanted version doesn't exists?")
		handleError(err)
	}

	fmt.Println("Path: " + path)
	fmt.Println("Installing ...")

	// Allowing the executable to claim ports under 1024
	runCommand("chmod", "+x", path)
	runCommand("chown", "git:git", path)
	runCommand("setcap", "cap_net_bind_service=+ep", path)

	// Creating a symlink to the new version
	createSymlink(version)

	// Starting Gitea service via systemctl
	runCommand("systemctl", "start", "gitea")

	fmt.Println("Update was sucessful!")
}

func readVersion() (version string) {
	// Ask the user for a version
	fmt.Print("Version (1.2.3): ")

	// Read the input
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	// Removing the linefeed
	text = strings.Replace(text, "\n", "", -1)
	text = strings.TrimSpace(text)

	return text
}

func runCommand(command string, args ...string) {
	// Constructing the command
	cmd := exec.Command(command, args...)

	// Running it
	fmt.Println("'" + command + " " + strings.Join(args, " ") + "'")
	cmd.Run()
}

func createSymlink(version string) {
	// If there's a old system link, we'll delete it
	_, err := os.Lstat("gitea")
	if err == nil {
		os.Remove("gitea")
	}

	// Creating the new system link
	err = os.Symlink("gitea-"+version, "gitea")
	if err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	// Print the error
	fmt.Println("\nError =>")
	panic(err)
}
