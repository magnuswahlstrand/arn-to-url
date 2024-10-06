package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/integrii/flaggy"
	"github.com/magnuswahlstrand/arn-to-url/awsurl"
	"log"
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
	roleMappings       []string
}

func configuration() Configuration {
	c := Configuration{
		// Default values:
		openInBrowser:      false,
		ignoreErrors:       false,
		accessPortalDomain: "",
		roleMappings:       []string{},
	}

	cliVersion := fmt.Sprintf("1.%s.0", VERSION)
	description := fmt.Sprintf("Reads AWS ARNs from standard input and resolves them to their corresponding AWS console URLs. Outputs can be directed to standard output or opened directly in a web browser.\n\nVersion: %s", cliVersion)
	flaggy.SetDescription(description)
	flaggy.SetVersion(cliVersion)
	flaggy.Bool(&c.openInBrowser, "w", "web", "Open URL(s) in the default browser (default standard out)")
	flaggy.Bool(&c.ignoreErrors, "e", "ignore-errors", "Ignore errors. Only opens or prints successfully resolved URLs")
	flaggy.String(&c.accessPortalDomain, "d", "domain", "Access portal domain. E.g. 'magnus' for magnus.awsapps.com/start")
	flaggy.StringSlice(&c.roleMappings, "r", "roles", "Comma separated list of account to IAM role to assume. E.g. 12345:admin,54321:dev. NOTE: Only used if access portal domain is set")
	flaggy.Parse()
	return c
}

func main() {
	c := configuration()
	resolver, err := awsurl.NewResolver(c.accessPortalDomain, c.roleMappings)
	if err != nil {
		log.Fatal(err)
	}

	// Chain resolver and handler together, for easier common error handling
	resolveAndHandle := func(arn string) error {
		url, err := resolver.FromArn2(arn)
		if err != nil {
			return err
		}

		if c.openInBrowser {
			return openInBrowser(url)
		}
		fmt.Println(url)
		return nil
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := resolveAndHandle(scanner.Text()); err != nil {
			if c.ignoreErrors {
				continue
			} else {
				log.Fatal(err)
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
