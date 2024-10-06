package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/integrii/flaggy"
	"github.com/magnuswahlstrand/arn-to-url/awsurl"
	"os"
	"os/exec"
	"runtime"
)

//go:embed VERSION
var VERSION string

type Configuration struct {
	openInBrowser      bool
	ignoreErrors       bool
	accessPortalDomain string
	roleName           string
}

func configuration() Configuration {
	c := Configuration{
		// Default values:
		openInBrowser:      false,
		ignoreErrors:       false,
		accessPortalDomain: "",
		roleName:           "",
	}

	flaggy.SetVersion(fmt.Sprintf("1.%s.0", VERSION))
	flaggy.Bool(&c.openInBrowser, "w", "web", "Open URL(s) in the default browser")
	flaggy.Bool(&c.ignoreErrors, "e", "ignore-errors", "Ignore errors. Only opens or prints successfully resolved URLs")
	flaggy.String(&c.accessPortalDomain, "d", "domain", "Access portal domain. E.g. 'magnus' for magnus.awsapps.com/start")
	//flaggy.String(&c.roleName, "r", "role", "Access portal domain. E.g. 'magnus' for magnus.awsapps.com/start")
	flaggy.Parse()
	return c
}

func main() {
	c := configuration()
	resolver := awsurl.Resolver{
		AccessPortalDomain: c.accessPortalDomain,
		// TODO: This should be a mapping function instead
		// TODO: Not implemented yet
		RoleName: c.roleName,
	}
	// Set up action to be called when a URL is parsed
	action := func(url string) error {
		if c.openInBrowser {
			return openInBrowser(url)
		}
		fmt.Println(url)
		return nil
	}

	// Chain resolver and action together, for easier common error handling
	chain := func(arn string) error {
		url, err := resolver.FromArn2(arn)
		if err != nil {
			return err
		}

		return action(url)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := chain(scanner.Text()); err != nil {
			if c.ignoreErrors {
				continue
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from stdin:", err)
	}
}

// From https://stackoverflow.com/a/39324149
// opens the specified URL in the default browser of the user.
func openInBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
