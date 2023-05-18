package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var version = "0.1.3"

//go:embed help.txt
var helpText string

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("")
		printHelp()
		os.Exit(1)
	}
}

func run() error {
	if _, err := exec.LookPath("git"); err != nil {
		return errors.New("cannot find git")
	}

	args := os.Args[1:]

	if len(args) == 0 || args[0] == "help" {
		return printHelp()
	}

	if args[0] == "version" {
		return printVersion()
	}

	id := args[0]
	if !isValidUserId(id) {
		return errors.New("invalid user ID")
	}

	if id == "default" {
		items := []string{"user.name", "user.email", "core.sshCommand"}
		for _, item := range items {
			if err := exec.Command("git", "config", "--unset", item).Run(); err != nil {
				return err
			}
		}

		return nil
	}

	args = args[1:]

	data := make(map[string]string)
	keys := []string{"name", "email", "cert"}

	for _, key := range keys {
		outputBytes, err := exec.Command("git", "config", "--global", fmt.Sprintf("git-as.%s.%s", id, key)).Output()
		if err != nil {
			return fmt.Errorf("entry for user ID '%s' not found or corrupted", id)
		}

		data[key] = strings.Trim(string(outputBytes), " \r\n")
	}

	if len(args) > 0 && args[0] == "clone" {
		if len(args) < 2 {
			return errors.New("repo argument not found")
		}

		args = append(args, "-c", fmt.Sprintf("user.name=%s", data["name"]), "-c", fmt.Sprintf("user.email=%s", data["email"]), "-c", fmt.Sprintf("core.sshCommand=ssh -i %s", data["cert"]))

		cmd := exec.Command("git", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	gitCommands := [][]string{
		{"config", "user.name", data["name"]},
		{"config", "user.email", data["email"]},
		{"config", "core.sshCommand", fmt.Sprintf("ssh -i %s", data["cert"])},
	}

	for _, arg := range gitCommands {
		if err := exec.Command("git", arg...).Run(); err != nil {
			return errors.New("error executing commands, manual cleanup of global .gitconfig file may be necessary")
		}
	}

	fmt.Printf("%s successfully activated.\n", id)

	return nil
}

func printHelp() error {
	fmt.Println(helpText)

	return nil
}

func printVersion() error {
	fmt.Printf("git-as v%s\n", version)

	return nil
}

func isValidUserId(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9-_]+$`).MatchString(s)
}
