package app

import (
	"ctRestClient/logger"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewKeepassCli", func() {
	It("returns error if file does not exist", func() {
		_, err := NewKeepassCli("/tmp/nonexistent.kdbx", "password", logger.NewLogger("/dev/null"))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("could not be found"))
	})

	It("returns error if path is a directory", func() {
		dir := GinkgoT().TempDir()
		_, err := NewKeepassCli(dir, "password", logger.NewLogger("/dev/null"))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not a regular file"))
	})

	It("returns cli if file exists", func() {
		file := GinkgoT().TempDir() + "/test.kdbx"
		f, err := os.Create(file)
		Expect(err).NotTo(HaveOccurred())
		f.Close()
		cli, err := NewKeepassCli(file, "password", logger.NewLogger("/dev/null"))
		Expect(err).NotTo(HaveOccurred())
		Expect(cli).NotTo(BeNil())
	})
})
