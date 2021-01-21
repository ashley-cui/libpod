package integration

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/containers/podman/v2/test/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Podman secret", func() {
	var (
		tempdir    string
		err        error
		podmanTest *PodmanTestIntegration
	)

	BeforeEach(func() {
		tempdir, err = CreateTempDirInTempDir()
		if err != nil {
			os.Exit(1)
		}
		podmanTest = PodmanTestCreate(tempdir)
		podmanTest.Setup()
		podmanTest.SeedImages()
	})

	AfterEach(func() {
		podmanTest.CleanupSecrets()
		f := CurrentGinkgoTestDescription()
		processTestResult(f)

	})

	It("podman secret create", func() {
		secretFilePath := filepath.Join(podmanTest.TempDir, "secret")
		err := ioutil.WriteFile(secretFilePath, []byte("mysecret"), 0755)
		Expect(err).To(BeNil())

		session := podmanTest.Podman([]string{"secret", "create", "a", secretFilePath})
		session.WaitWithDefaultTimeout()
		secrID := session.OutputToString()
		Expect(session.ExitCode()).To(Equal(0))

		inspect := podmanTest.Podman([]string{"secret", "inspect", "--format", "{{.ID}}", secrID})
		inspect.WaitWithDefaultTimeout()
		Expect(inspect.ExitCode()).To(Equal(0))
		Expect(inspect.OutputToString()).To(Equal(secrID))
	})
	It("podman secret create bad name should fail", func() {
		secretFilePath := filepath.Join(podmanTest.TempDir, "secret")
		err := ioutil.WriteFile(secretFilePath, []byte("mysecret"), 0755)
		Expect(err).To(BeNil())

		session := podmanTest.Podman([]string{"secret", "create", "?!", secretFilePath})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Not(Equal(0)))
	})
	It("podman secret inspect", func() {
		secretFilePath := filepath.Join(podmanTest.TempDir, "secret")
		err := ioutil.WriteFile(secretFilePath, []byte("mysecret"), 0755)
		Expect(err).To(BeNil())

		session := podmanTest.Podman([]string{"secret", "create", "a", secretFilePath})
		session.WaitWithDefaultTimeout()
		secrID := session.OutputToString()
		Expect(session.ExitCode()).To(Equal(0))

		inspect := podmanTest.Podman([]string{"secret", "inspect", secrID})
		inspect.WaitWithDefaultTimeout()
		Expect(inspect.ExitCode()).To(Equal(0))
		Expect(inspect.IsJSONOutputValid()).To(BeTrue())
	})

	It("podman secret inspect with --format", func() {
		secretFilePath := filepath.Join(podmanTest.TempDir, "secret")
		err := ioutil.WriteFile(secretFilePath, []byte("mysecret"), 0755)
		Expect(err).To(BeNil())

		session := podmanTest.Podman([]string{"secret", "create", "a", secretFilePath})
		session.WaitWithDefaultTimeout()
		secrID := session.OutputToString()
		Expect(session.ExitCode()).To(Equal(0))

		inspect := podmanTest.Podman([]string{"secret", "inspect", "--format", "{{.ID}}", secrID})
		inspect.WaitWithDefaultTimeout()
		Expect(inspect.ExitCode()).To(Equal(0))
		Expect(inspect.OutputToString()).To(Equal(secrID))
	})
	It("podman secret inspect multiple secrets", func() {
		secretFilePath := filepath.Join(podmanTest.TempDir, "secret")
		err := ioutil.WriteFile(secretFilePath, []byte("mysecret"), 0755)
		Expect(err).To(BeNil())

		session := podmanTest.Podman([]string{"secret", "create", "a", secretFilePath})
		session.WaitWithDefaultTimeout()
		secrID := session.OutputToString()
		Expect(session.ExitCode()).To(Equal(0))

		session2 := podmanTest.Podman([]string{"secret", "create", "b", secretFilePath})
		session2.WaitWithDefaultTimeout()
		secrID2 := session2.OutputToString()
		Expect(session2.ExitCode()).To(Equal(0))

		inspect := podmanTest.Podman([]string{"secret", "inspect", secrID, secrID2})
		inspect.WaitWithDefaultTimeout()
		Expect(inspect.ExitCode()).To(Equal(0))
		Expect(inspect.IsJSONOutputValid()).To(BeTrue())
	})

	It("podman secret inspect bogus", func() {
		secretFilePath := filepath.Join(podmanTest.TempDir, "secret")
		err := ioutil.WriteFile(secretFilePath, []byte("mysecret"), 0755)
		Expect(err).To(BeNil())

		inspect := podmanTest.Podman([]string{"secret", "inspect", "bogus"})
		inspect.WaitWithDefaultTimeout()
		Expect(inspect.ExitCode()).To(Not(Equal(0)))

	})

	It("podman secret ls", func() {
		secretFilePath := filepath.Join(podmanTest.TempDir, "secret")
		err := ioutil.WriteFile(secretFilePath, []byte("mysecret"), 0755)
		Expect(err).To(BeNil())

		session := podmanTest.Podman([]string{"secret", "create", "a", secretFilePath})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		list := podmanTest.Podman([]string{"secret", "ls"})
		list.WaitWithDefaultTimeout()
		Expect(list.ExitCode()).To(Equal(0))
		Expect(len(list.OutputToStringArray())).To(Equal(2))

	})
	It("podman secret ls with Go template", func() {
		secretFilePath := filepath.Join(podmanTest.TempDir, "secret")
		err := ioutil.WriteFile(secretFilePath, []byte("mysecret"), 0755)
		Expect(err).To(BeNil())

		session := podmanTest.Podman([]string{"secret", "create", "a", secretFilePath})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		list := podmanTest.Podman([]string{"secret", "ls", "--format", "table {{.Name}}"})
		list.WaitWithDefaultTimeout()

		Expect(list.ExitCode()).To(Equal(0))
		Expect(len(list.OutputToStringArray())).To(Equal(2), list.OutputToString())
	})
	It("podman secret rm", func() {
		secretFilePath := filepath.Join(podmanTest.TempDir, "secret")
		err := ioutil.WriteFile(secretFilePath, []byte("mysecret"), 0755)
		Expect(err).To(BeNil())

		session := podmanTest.Podman([]string{"secret", "create", "a", secretFilePath})
		session.WaitWithDefaultTimeout()
		secrID := session.OutputToString()
		Expect(session.ExitCode()).To(Equal(0))

		removed := podmanTest.Podman([]string{"secret", "rm", "a"})
		removed.WaitWithDefaultTimeout()
		Expect(removed.ExitCode()).To(Equal(0))
		Expect(removed.OutputToString()).To(Equal(secrID))

		session = podmanTest.Podman([]string{"secret", "ls"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(len(session.OutputToStringArray())).To(Equal(1))
	})
	It("podman secret rm --all", func() {
		secretFilePath := filepath.Join(podmanTest.TempDir, "secret")
		err := ioutil.WriteFile(secretFilePath, []byte("mysecret"), 0755)
		Expect(err).To(BeNil())

		session := podmanTest.Podman([]string{"secret", "create", "a", secretFilePath})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.Podman([]string{"secret", "create", "b", secretFilePath})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		removed := podmanTest.Podman([]string{"secret", "rm", "-a"})
		removed.WaitWithDefaultTimeout()
		Expect(removed.ExitCode()).To(Equal(0))

		session = podmanTest.Podman([]string{"secret", "ls"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(len(session.OutputToStringArray())).To(Equal(1))
	})

	It("podman run container with secret", func() {
		secretFilePath := filepath.Join(podmanTest.TempDir, "secret")
		err := ioutil.WriteFile(secretFilePath, []byte("somesecretdata"), 0755)
		Expect(err).To(BeNil())

		session := podmanTest.Podman([]string{"secret", "create", "mysecret", secretFilePath})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		removed := podmanTest.Podman([]string{"secret", "rm", "-a"})
		removed.WaitWithDefaultTimeout()
		Expect(removed.ExitCode()).To(Equal(0))

		session = podmanTest.Podman([]string{"secret", "ls"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(len(session.OutputToStringArray())).To(Equal(1))
	})

})
