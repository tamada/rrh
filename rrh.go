package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/internal"
	"github.com/tamada/rrh/lib"
)

func executeInternalCommand(commands map[string]cli.CommandFactory, args []string) (int, error) {
	var c = cli.NewCLI("rrh", lib.VERSION)
	c.Name = "rrh"
	c.Args = args
	c.Autocomplete = true
	c.Commands = commands
	return c.Run()
}

func executeCommand(path string, args []string) (int, error) {
	var cmd = exec.Command(path, args...)
	var output, err = cmd.Output()
	if err != nil {
		return 4, err
	}
	fmt.Print(string(output))
	return 0, nil
}

func findExecutableFromPathEnv(commandName string) string {
	var pathEnv = os.Getenv("PATH")
	for _, env := range strings.Split(pathEnv, ":") {
		if findExecutableFromDir(env, commandName) {
			return filepath.Join(env, commandName)
		}
	}
	return ""
}

func findExecutableFromDir(dir, commandName string) bool {
	var path = filepath.Join(dir, commandName)
	var finfo, err = os.Stat(path)
	if err != nil {
		return false
	}
	if finfo.Mode().IsRegular() && (finfo.Mode().Perm()&0777) == 0755 {
		return true
	}
	return false
}

type rrhOptions struct {
	help       bool
	version    bool
	configPath string
}

func parseOptions(args []string, opts *rrhOptions) []string {
	var configPathFlag = false
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			switch arg {
			case "-h", "--help":
				opts.help = true
			case "-v", "--version":
				opts.version = true
			case "-c", "--config-file":
				configPathFlag = true
			}
		} else if configPathFlag {
			opts.configPath = arg
			configPathFlag = false
		} else {
			return args[i:]
		}
	}
	return []string{}
}

func (opts *rrhOptions) printHelpOrVersion(args []string) (int, error) {
	if opts.version {
		var com, _ = internal.VersionCommandFactory()
		com.Run([]string{})
	}
	if opts.help || (!opts.version && len(args) == 0) {
		var com, _ = internal.HelpCommandFactory()
		com.Run([]string{})
	}
	return 0, nil
}

func executeExternalCommand(args []string) (int, error) {
	var commandName = fmt.Sprintf("rrh-%s", args[0])
	var executablePath = findExecutableFromPathEnv(commandName)
	if executablePath == "" {
		return 3, fmt.Errorf("%s: command not found", args[0])
	}
	return executeCommand(executablePath, args[1:])
}

func (opts *rrhOptions) updateConfigPath() {
	if opts.configPath != "" {
		os.Setenv(lib.RrhConfigPath, opts.configPath)
	}
}

func goMain(args []string) (int, error) {
	var commands = internal.BuildCommandFactoryMap()
	var opts = new(rrhOptions)
	var newArgs = parseOptions(args[1:], opts)
	if len(newArgs) == 0 || opts.help || opts.version {
		return opts.printHelpOrVersion(newArgs)
	}
	opts.updateConfigPath()
	if commands[newArgs[0]] != nil {
		return executeInternalCommand(commands, newArgs)
	}
	return executeExternalCommand(newArgs)
}

func main() {
	var exitStatus, err = goMain(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(internal.GenerateDefaultHelp())
	}
	os.Exit(exitStatus)
}
