package decorator

import (
	"os"

	"github.com/mattn/go-isatty"
)

type Decorator interface {
	RepositoryID(repoID string) string
	GroupName(name string) string
	EnvironmentValue(value string) string
	EnvironmentLabel(label string) string
}

func New(colorSettings string) (Decorator, error) {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return NewNoDecorator(), nil
	}
	return NewColorDecorator(colorSettings)
}

func NewNoDecorator() Decorator {
	return &noDecorator{}
}
