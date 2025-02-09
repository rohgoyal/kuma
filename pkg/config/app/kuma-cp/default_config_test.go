package kuma_cp

import (
	"github.com/Kong/kuma/pkg/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strings"
)

var _ = Describe("Default config", func() {

	var backupEnvVars []string

	BeforeEach(func() {
		backupEnvVars = os.Environ()
		os.Clearenv()
	})
	AfterEach(func() {
		for _, envVar := range backupEnvVars {
			parts := strings.SplitN(envVar, "=", 2)
			Expect(os.Setenv(parts[0], parts[1])).To(Succeed())
		}
	})

	It("should be check against the kuma-cp.defaults.yaml file", func() {
		// given
		cfg := Config{}

		// when
		err := config.Load("kuma-cp.defaults.yaml", &cfg)

		// then
		Expect(err).ToNot(HaveOccurred())
		Expect(DefaultConfig()).To(Equal(cfg), "The default config generated by Kuma and kuma-cp.defaults.yaml config file are different. Please update the kuma-cp.defaults.yaml file.")
	})
})
