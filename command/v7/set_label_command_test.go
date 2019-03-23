package v7_test

import (
	"errors"
	"regexp"

	"code.cloudfoundry.org/cli/command/commandfakes"
	"code.cloudfoundry.org/cli/command/flag"
	. "code.cloudfoundry.org/cli/command/v7"
	"code.cloudfoundry.org/cli/util/configv3"
	"code.cloudfoundry.org/cli/util/ui"
	uuid "github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = FDescribe("set-label command", func() {
	var (
		cmd             SetLabelCommand
		fakeConfig      *commandfakes.FakeConfig
		fakeSharedActor *commandfakes.FakeSharedActor
		testUI          *ui.UI

		executeErr error
	)

	BeforeEach(func() {
		testUI = ui.NewTestUI(nil, NewBuffer(), NewBuffer())
		fakeConfig = new(commandfakes.FakeConfig)
		fakeSharedActor = new(commandfakes.FakeSharedActor)
		cmd = SetLabelCommand{
			UI:          testUI,
			Config:      fakeConfig,
			SharedActor: fakeSharedActor,
		}
	})

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	It("doesn't error", func() {
		Expect(executeErr).ToNot(HaveOccurred())
	})

	It("checks that the user is logged in and targeted to an org and space", func() {
		Expect(fakeSharedActor.CheckTargetCallCount()).To(Equal(1))
		checkOrg, checkSpace := fakeSharedActor.CheckTargetArgsForCall(0)
		Expect(checkOrg).To(BeTrue())
		Expect(checkSpace).To(BeTrue())
	})

	When("checking target succeeds", func() {
		var appName string

		BeforeEach(func() {
			fakeSharedActor.CheckTargetReturns(nil)
		})

		When("fetching current user's name succeeds", func() {
			BeforeEach(func() {
				fakeConfig.CurrentUserNameReturns("some-user", nil)
				fakeConfig.TargetedOrganizationReturns(configv3.Organization{Name: "fake-org"})
				fakeConfig.TargetedSpaceReturns(configv3.Space{Name: "fake-space"})

				u, _ := uuid.NewV4()
				appName = u.String()
				cmd.RequiredArgs = flag.SetLabelArgs{
					ResourceType: "app",
					ResourceName: appName,
				}
			})

			It("sets the provided labels on the app", func() {
				Expect(fakeActor.UpdateLabels
			})

			It("displays a message", func() {
				Expect(executeErr).ToNot(HaveOccurred())

				Expect(fakeSharedActor.CheckTargetCallCount()).To(Equal(1))
				Expect(fakeConfig.TargetedOrganizationCallCount()).To(Equal(1))
				Expect(fakeConfig.TargetedSpaceCallCount()).To(Equal(1))
				Expect(fakeConfig.CurrentUserNameCallCount()).To(Equal(1))

				Expect(testUI.Out).To(Say(regexp.QuoteMeta(`Setting label(s) for app %s in org fake-org / space fake-space as some-user...`), appName))
				Expect(testUI.Out).To(Say("OK"))
			})
		})

		When("fetching the current user's name fails", func() {
			BeforeEach(func() {
				fakeConfig.CurrentUserNameReturns("some-user", errors.New("boom"))
			})

			It("returns an error", func() {
				Expect(executeErr).To(MatchError("boom"))
			})
		})
	})

	When("checking targeted org and space fails", func() {
		BeforeEach(func() {
			fakeSharedActor.CheckTargetReturns(errors.New("nope"))
		})

		It("returns an error", func() {
			Expect(executeErr).To(MatchError("nope"))
		})
	})
})
