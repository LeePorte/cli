package user

import (
	"fmt"

	"code.cloudfoundry.org/cli/cf/api"
	"code.cloudfoundry.org/cli/cf/commandregistry"
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/cf/flags"
	. "code.cloudfoundry.org/cli/cf/i18n"
	"code.cloudfoundry.org/cli/cf/requirements"
	"code.cloudfoundry.org/cli/cf/terminal"
)

type FindUser struct {
	ui       terminal.UI
	config   coreconfig.Reader
	userRepo api.UserRepository
}

func init() {
	commandregistry.Register(&FindUser{})
}

func (cmd *FindUser) MetaData() commandregistry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["f"] = &flags.BoolFlag{ShortName: "f", Usage: T("Force deletion without confirmation")}

	return commandregistry.CommandMetadata{
		Name:        "find-user",
		Description: T("Find a user"),
		Usage: []string{
			T("CF_NAME find-user USERNAME"),
		},
	}
}

func (cmd *FindUser) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) ([]requirements.Requirement, error) {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires an argument\n\n") + commandregistry.Commands.CommandUsage("find-user"))
		return nil, fmt.Errorf("Incorrect usage: %d arguments of %d required", len(fc.Args()), 1)
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}

	return reqs, nil
}

func (cmd *FindUser) SetDependency(deps commandregistry.Dependency, pluginCall bool) commandregistry.Command {
	cmd.ui = deps.UI
	cmd.config = deps.Config
	cmd.userRepo = deps.RepoLocator.GetUserRepository()
	return cmd
}

func (cmd *FindUser) Execute(c flags.FlagContext) error {
	username := c.Args()[0]


	cmd.ui.Say(T("Finding user {{.TargetUser}} as {{.CurrentUser}}...",
		map[string]interface{}{
			"TargetUser":  terminal.EntityNameColor(username),
			"CurrentUser": terminal.EntityNameColor(cmd.config.Username()),
		}))

	users, err := cmd.userRepo.FindAllByUsername(username)

	switch err.(type) {
	case nil:
		//if len(users) == 1 {
		//	return fmt.Errorf(T(
		//		"Found user {{.Username}} \nThe user exists.",
		//		map[string]interface{}{
		//			"Username": username,
		//		}))
		//}
		if len(users) > 1 {
			return fmt.Errorf(T(
				"Error finding user {{.Username}} \nThe user exists in multiple origins.",
				map[string]interface{}{
					"Username": username,
				}))
		}
	case *errors.ModelNotFoundError:
		cmd.ui.Ok()
		cmd.ui.Warn(T("User {{.TargetUser}} does not exist.", map[string]interface{}{"TargetUser": username}))
		return nil
	default:
		return err
	}


	cmd.ui.Ok()
	return nil
}
