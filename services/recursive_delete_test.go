package services_test

import (
	. "github.com/cloudfoundry/cf-acceptance-tests/helpers/services"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	"github.com/cloudfoundry/cf-acceptance-tests/helpers/assets"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Recursive Delete", func() {
	var broker ServiceBroker
	var orgName string

	BeforeEach(func() {
		broker = NewServiceBroker(generator.RandomName(), assets.NewAssets().ServiceBroker, context)
		broker.Push()
		broker.Configure()
		broker.Create()
		broker.PublicizePlans()

		orgName = generator.RandomName()
		spaceName := generator.RandomName()
		appName := generator.RandomName()
		instanceName := generator.RandomName()

		cf.AsUser(context.AdminUserContext(), DEFAULT_TIMEOUT, func() {
			createOrg := cf.Cf("create-org", orgName).Wait(DEFAULT_TIMEOUT)
			Expect(createOrg).To(Exit(0), "failed to create org")

			createSpace := cf.Cf("create-space", spaceName, "-o", orgName).Wait(DEFAULT_TIMEOUT)
			Expect(createSpace).To(Exit(0), "failed to create space")

			target := cf.Cf("target", "-o", orgName, "-s", spaceName).Wait(CF_PUSH_TIMEOUT)
			Expect(target).To(Exit(0), "failed targeting")

			createApp := cf.Cf("push", appName, "-p", assets.NewAssets().Dora).Wait(CF_PUSH_TIMEOUT)
			Expect(createApp).To(Exit(0), "failed creating app")

			createService := cf.Cf("create-service", broker.Service.Name, broker.SyncPlans[0].Name, instanceName).Wait(DEFAULT_TIMEOUT)
			Expect(createService).To(Exit(0), "failed creating service")
		})
	})

	AfterEach(func() {
		broker.Destroy()
	})

	It("deletes all apps and services in all spaces in an org", func() {
		cf.AsUser(context.AdminUserContext(), DEFAULT_TIMEOUT, func() {
			deleteOrg := cf.Cf("delete-org", orgName, "-f").Wait(DEFAULT_TIMEOUT)
			Expect(deleteOrg).To(Exit(0), "failed deleting org")
		})
		getOrg := cf.Cf("org", orgName).Wait(DEFAULT_TIMEOUT)
		Expect(getOrg).To(Exit(1), "org still exists")
	})
})
