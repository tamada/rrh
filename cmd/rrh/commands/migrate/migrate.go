package migrate

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/common"
)

type Version int

const (
	V1 Version = iota + 1
	V2
	unknownVersion Version = -1
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migrate <major_version>",
		Version: rrh.VERSION,
		Short:   "migrate the rrh settings from the given version",
		Long: `migrate the rrh settings from the given version.
Available major versions are:
  1.x.x`,
		Args: cobra.MinimumNArgs(1),
		RunE: perform,
	}
	return cmd
}

func parseMajorVersion(str string) Version {
	if strings.HasPrefix(str, "1.") {
		return V1
	} else if strings.HasPrefix(str, "2.") {
		return V2
	}
	return unknownVersion
}

func expandPath(path string) string {
	if strings.Contains(path, "${HOME}") {
		home, _ := os.UserHomeDir()
		return strings.ReplaceAll(path, "${HOME}", home)
	}
	return path
}

func migrateFromVersion1Impl(c *cobra.Command, legacyPath, newPath string) error {
	if err := os.MkdirAll(newPath, 0755); err != nil {
		return err
	}
	finfos, err := ioutil.ReadDir(legacyPath)
	if err != nil {
		return err
	}
	errs := common.NewErrorList()
	for _, finfo := range finfos {
		err := copyFile(filepath.Join(legacyPath, finfo.Name()), filepath.Join(newPath, finfo.Name()))
		errs = errs.Append(err)
	}
	return errs.NilOrThis()
}

func copyFile(from, to string) error {
	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()
	dest, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	return err
}

func removeLegacySettingDir(c *cobra.Command, legacyPath string) error {
	c.Printf("Remove legacy settings in %s? [y/N] ", legacyPath)
	if isInputYes(c) {
		if err := os.RemoveAll(legacyPath); err != nil {
			return err
		}
		fmt.Printf("remove %s done\n", legacyPath)
	} else {
		c.Printf("do not remove %s\n", legacyPath)
	}
	return nil
}

func isInputYes(c *cobra.Command) bool {
	scanner := bufio.NewScanner(c.InOrStdin())
	scanner.Scan()
	text := strings.ToLower(scanner.Text())
	return text == "yes" || text == "y"
}

func migrateFromVersion1(c *cobra.Command) error {
	legacyPath := expandPath("${HOME}/.rrh")
	if rrh.IsExistDir(legacyPath) {
		newPath := expandPath("${HOME}/.config/rrh")
		if !rrh.IsExistDir(newPath) {
			if err := migrateFromVersion1Impl(c, legacyPath, newPath); err != nil {
				return err
			}
		} else {
			c.Println("Already migrate from version 1.x.x")
		}
		return removeLegacySettingDir(c, legacyPath)
	}
	c.Printf("%s: directory not found", legacyPath)
	return nil
}

func perform(c *cobra.Command, args []string) error {
	version := parseMajorVersion(args[0])
	switch version {
	case V1:
		return migrateFromVersion1(c)
	case V2:
		c.Printf("nothing to migrate from %s\n", args[0])
	case unknownVersion:
		return fmt.Errorf("%s: unknown version", args[0])
	}
	return nil
}
