package user_test

import (
	"code.cloudfoundry.org/cli/cf/api/apifakes"
	"code.cloudfoundry.org/cli/cf/commandregistry"
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/cf/models"
	"code.cloudfoundry.org/cli/cf/requirements"
	"code.cloudfoundry.org/cli/cf/requirements/requirementsfakes"
	testcmd "code.cloudfoundry.org/cli/cf/util/testhelpers/commands"
	testconfig "code.cloudfoundry.org/cli/cf/util/testhelpers/configuration"
	testterm "code.cloudfoundry.org/cli/cf/util/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "code.cloudfoundry.org/cli/cf/util/testhelpers/matchers"
)

var _ = Describe("find-user command", func() {
	var (
		ui                  *testterm.FakeUI
		configRepo          coreconfig.Repository
		userRepo            *apifakes.FakeUserRepository
		requirementsFactory *requirementsfakes.FakeFactory
		deps                commandregistry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.UI = ui
		deps.Config = configRepo
		deps.RepoLocator = deps.RepoLocator.SetUserRepository(userRepo)
		commandregistry.Commands.SetCommand(commandregistry.Commands.FindCommand("find-user").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{Inputs: []string{"y"}}
		userRepo = new(apifakes.FakeUserRepository)
		requirementsFactory = new(requirementsfakes.FakeFactory)
		requirementsFactory.NewLoginRequirementReturns(requirements.Passing{})
		configRepo = testconfig.NewRepositoryWithDefaults()

		token, err := testconfig.EncodeAccessToken(coreconfig.TokenInfo{
			UserGUID: "admin-user-guid",
			Username: "admin-user",
		})
		Expect(err).ToNot(HaveOccurred())
		configRepo.SetAccessToken(token)
	})

	runCommand := func(args ...string) bool {
		return testcmd.RunCLICommand("delete-user", args, requirementsFactory, updateCommandDependency, false, ui)
	}

	Describe("requirements", func() {
		It("fails when not logged in", func() {
			requirementsFactory.NewLoginRequirementReturns(requirements.Failing{Message: "not logged in"})

			Expect(runCommand("my-user")).To(BeFalse())
		})

		It("fails with usage when no arguments are given", func() {
			runCommand()
			Expect(ui.Outputs()).To(ContainSubstrings(
				[]string{"Incorrect Usage", "Requires an argument"},
			))
		})
	})

	Context("when the given user exists", func() {
		BeforeEach(func() {
			userRepo.FindAllByUsernameReturns([]models.UserFields{{
				Username: "user-name",
				GUID:     "user-guid",
			}}, nil)
		})

		It("finds a user with the given name", func() {
			runCommand("user-name")

			Expect(ui.Outputs()).To(ContainSubstrings(
				[]string{"Found user", "user-name", "admin-user"},
				[]string{"OK"},
			))

			Expect(userRepo.FindAllByUsernameArgsForCall(0)).To(Equal("user-name"))
		})

	})

	Context("when multiple users exist", func() {
		BeforeEach(func() {
			userRepo.FindAllByUsernameReturns([]models.UserFields{
				{
					Username: "user-name1",
					GUID:     "user-guid1"},
				{
					Username: "user-name2",
					GUID:     "user-guid2"},
			}, nil)
		})

		It("returns a muliple users found error", func() {
			runCommand("user-name")

			Expect(ui.Outputs()).To(ContainSubstrings(
				[]string{"Error finding user"},
				[]string{"The user exists in multiple origins."},
			))
		})
	})
	Context("when the given user does not exist", func() {
		BeforeEach(func() {
			userRepo.FindAllByUsernameReturns(nil, errors.NewModelNotFoundError("User", ""))
		})

		It("prints a warning", func() {
			runCommand("-f", "user-name")

			Expect(ui.WarnOutputs).To(ContainSubstrings([]string{"user-name", "does not exist"}))

			Expect(userRepo.FindAllByUsernameArgsForCall(0)).To(Equal("user-name"))
		})
	})
})
