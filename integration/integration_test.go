package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/dagger"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	spec.Run(t, "Integration", testIntegration, spec.Report(report.Terminal{}))
}

func testIntegration(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect  func(interface{}, ...interface{}) GomegaAssertion
		app     *dagger.App
		builder string
		err     error
	)

	it.Before(func() {
		Expect = NewWithT(t).Expect
		config, err := dagger.ParseConfig("config.json")
		Expect(err).NotTo(HaveOccurred())

		builder = config.Builder

	})

	it.After(func() {
		app.Destroy()
	})

	when("trying to build a .NET app using the metabuildpack", func() {
		it("successfully builds and runs an app", func() {
			app, err = dagger.NewPack(
				filepath.Join("testdata", "simple_3.1_source"),
				dagger.RandomImage(),
				dagger.SetBuilder(builder),
				dagger.NoPull(),
			).Build()
			Expect(err).NotTo(HaveOccurred())

			if builder == "bionic" {
				app.SetHealthCheck("stat /workspace", "2s", "15s")
			}

			Expect(app.Start()).To(Succeed())

			// we use the correct stack & buildpack

			// this is kind of awkwards...
			if builder == "cflinuxfs3" {
				Expect(app.BuildLogs()).NotTo(ContainSubstring("paketo-buildpacks/icu"))
			}
			body, _, err := app.HTTPGet("/")

			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(ContainSubstring("Hello World!"))
		})
	})

}
