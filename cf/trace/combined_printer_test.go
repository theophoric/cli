package trace_test

import (
	. "github.com/theophoric/cf-cli/cf/trace"

	"github.com/theophoric/cf-cli/cf/trace/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CombinePrinters", func() {
	var printer1, printer2 *fakes.FakePrinter
	var printer Printer

	BeforeEach(func() {
		printer1 = &fakes.FakePrinter{}
		printer2 = &fakes.FakePrinter{}

		printer = CombinePrinters([]Printer{printer1, printer2})
	})

	It("returns a combined printer that Prints", func() {
		printer.Print("foo", "bar")

		Expect(printer1.PrintCallCount()).To(Equal(1))
		Expect(printer2.PrintCallCount()).To(Equal(1))

		expectedArgs := []interface{}{"foo", "bar"}

		Expect(printer1.PrintArgsForCall(0)).To(Equal(expectedArgs))
		Expect(printer2.PrintArgsForCall(0)).To(Equal(expectedArgs))
	})

	It("returns a combined printer that Printfs", func() {
		printer.Printf("format %s %s", "arg1", "arg2")

		Expect(printer1.PrintfCallCount()).To(Equal(1))
		Expect(printer2.PrintfCallCount()).To(Equal(1))

		expectedArgs := []interface{}{"arg1", "arg2"}

		fmt1, args1 := printer1.PrintfArgsForCall(0)
		fmt2, args2 := printer2.PrintfArgsForCall(0)

		Expect(fmt1).To(Equal("format %s %s"))
		Expect(fmt2).To(Equal("format %s %s"))

		Expect(args1).To(Equal(expectedArgs))
		Expect(args2).To(Equal(expectedArgs))
	})

	It("returns a combined printer that Printlns", func() {
		printer.Println("foo", "bar")

		Expect(printer1.PrintlnCallCount()).To(Equal(1))
		Expect(printer2.PrintlnCallCount()).To(Equal(1))

		expectedArgs := []interface{}{"foo", "bar"}

		Expect(printer1.PrintlnArgsForCall(0)).To(Equal(expectedArgs))
		Expect(printer2.PrintlnArgsForCall(0)).To(Equal(expectedArgs))
	})
})
