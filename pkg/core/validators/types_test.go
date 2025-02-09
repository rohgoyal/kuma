package validators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Kong/kuma/pkg/core/validators"
)

var _ = Describe("Validation Error", func() {
	It("should construct errors", func() {
		// given
		err := validators.ValidationError{}

		// when
		err.AddViolation("name", "invalid name")

		// and
		addressErr := validators.ValidationError{}
		addressErr.AddViolation("street", "invalid format")
		err.AddError("address", addressErr)

		// then
		Expect(err.HasViolations()).To(BeTrue())
		Expect(validators.IsValidationError(&err)).To(BeTrue())
		Expect(err.OrNil()).To(MatchError("name: invalid name; address.street: invalid format"))
	})

	It("should convert to nil error when there are no violations", func() {
		// given
		validationErr := validators.ValidationError{}

		// when
		err := validationErr.OrNil()

		Expect(err).To(BeNil())
	})

	Describe("Append()", func() {
		It("should add a given error to the end of the list", func() {
			// given
			err := validators.ValidationError{}
			err1 := validators.ValidationError{}
			err1.AddViolationAt(validators.RootedAt("sources"), "unknown error")
			err2 := validators.ValidationError{}
			err2.AddViolationAt(validators.RootedAt("destinations"), "yet another error")

			By("adding the first error")
			// when
			err.Add(err1)
			// then
			Expect(err).To(Equal(validators.ValidationError{
				Violations: []validators.Violation{
					{Field: "sources", Message: "unknown error"},
				},
			}))

			By("adding the second error")
			// when
			err.Add(err2)
			// then
			Expect(err).To(Equal(validators.ValidationError{
				Violations: []validators.Violation{
					{Field: "sources", Message: "unknown error"},
					{Field: "destinations", Message: "yet another error"},
				},
			}))
		})
	})

	Describe("AddViolationAt()", func() {
		It("should accept nil PathBuilder", func() {
			// given
			err := validators.ValidationError{}
			// when
			err.AddViolationAt(nil, "unknown error")
			// then
			Expect(err).To(Equal(validators.ValidationError{
				Violations: []validators.Violation{
					{Field: "", Message: "unknown error"},
				},
			}))
		})

		It("should accept non-nil PathBuilder", func() {
			// given
			err := validators.ValidationError{}
			path := validators.RootedAt("sources").Index(0).Field("match").Key("service")
			// when
			err.AddViolationAt(path, "unknown error")
			// and
			Expect(err).To(Equal(validators.ValidationError{
				Violations: []validators.Violation{
					{Field: `sources[0].match["service"]`, Message: "unknown error"},
				},
			}))
		})
	})
})

var _ = Describe("PathBuilder", func() {
	It("should produce empty path by default", func() {
		Expect(validators.PathBuilder{}.String()).To(Equal(""))
	})

	It("should produce valid root path", func() {
		Expect(validators.RootedAt("spec").String()).To(Equal("spec"))
	})

	It("should produce valid field path", func() {
		Expect(validators.RootedAt("spec").Field("sources").String()).To(Equal("spec.sources"))
	})

	It("should produce valid array index", func() {
		Expect(validators.RootedAt("spec").Field("sources").Index(0).String()).To(Equal("spec.sources[0]"))
	})

	It("should produce valid array index", func() {
		Expect(validators.RootedAt("spec").Field("sources").Index(0).Field("match").Key("service").String()).To(Equal(`spec.sources[0].match["service"]`))
	})
})
