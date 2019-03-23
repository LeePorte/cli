package v7

import (
	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
)

type SetLabelCommand struct {
	RequiredArgs flag.SetLabelArgs `positional-args:"yes"`
	usage        interface{}       `usage:"cf set-label RESOURCE RESOURCE_NAME KEY=VALUE...\n\nEXAMPLES:\n   cf set-label app dora env=production\n\n RESOURCES:\n   APP\n\nSEE ALSO:\n   delete-label, labels"`

	UI          command.UI
	Config      command.Config
	SharedActor command.SharedActor
	// Actor SetHealthCheckActor
}

func (cmd *SetLabelCommand) Setup(config command.Config, ui command.UI) error {
	cmd.UI = ui
	cmd.Config = config
	cmd.SharedActor = sharedaction.NewActor(config)
	return nil
}

func (cmd SetLabelCommand) Execute(args []string) error {
	err := cmd.SharedActor.CheckTarget(true, true)
	if err != nil {
		return err
	}

	username, err := cmd.Config.CurrentUserName()
	if err != nil {
		return err
	}

	cmd.UI.DisplayTextWithFlavor(
		"Setting label(s) for {{.ResourceType}} {{.ResourceName}} in org {{.OrgName}} / space {{.SpaceName}} as {{.User}}...",
		map[string]interface{}{
			"ResourceType": cmd.RequiredArgs.ResourceType,
			"ResourceName": cmd.RequiredArgs.ResourceName,
			"OrgName":      cmd.Config.TargetedOrganization().Name,
			"SpaceName":    cmd.Config.TargetedSpace().Name,
			"User":         username,
		},
	)
	cmd.UI.DisplayOK()

	return nil
}
