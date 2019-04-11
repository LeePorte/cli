package v6

import (
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/translatableerror"
)

type FindUserCommand struct {
	RequiredArgs    flag.Username `positional-args:"no"`
	usage           interface{}   `usage:"CF_NAME find-user USERNAME"`
	relatedCommands interface{}   `related_commands:"org-users"`
}

func (FindUserCommand) Setup(config command.Config, ui command.UI) error {
	return nil
}

func (FindUserCommand) Execute(args []string) error {
	return translatableerror.UnrefactoredCommandError{}
}
