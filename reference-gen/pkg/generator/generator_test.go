package generator

import (
	"bytes"
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-git/go-git/v5/utils/diff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	testDataPackage = "github.com/oauth2-proxy/tools/reference-gen/pkg/generator/testdata"
)

//go:embed testdata/*.md
var testOutputs embed.FS

var _ = Describe("Generator", func() {
	type generatorTableInput struct {
		requestedTypes         []string
		headerFileName         string
		expectedOutputFileName string
	}

	DescribeTable("should generate the expected output", func(in generatorTableInput) {
		By("Creating an output file")
		outputFile, err := ioutil.TempFile("", "oauth2-proxy-reference-generator-suite-")
		Expect(err).ToNot(HaveOccurred())

		outputFileName := outputFile.Name()
		Expect(outputFile.Close()).To(Succeed())

		By("Constructing the generator")
		gen, err := NewGenerator(testDataPackage, in.requestedTypes, in.headerFileName, outputFileName, "")
		Expect(err).ToNot(HaveOccurred())

		By("Running the generator")
		Expect(gen.Run()).To(Succeed())

		By("Loading the output")
		output, err := os.ReadFile(outputFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Loading the expected output")
		expectedOutput, err := testOutputs.ReadFile(in.expectedOutputFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Comparing the outputs")
		diffs := diff.Do(string(expectedOutput), string(output))
		if len(diffs) > 1 {
			// A single diff means the two files are equal, only fail if there is more than one diff.
			fmt.Printf("\nUnexpected diff:\n\n%s\n", prettyPrintDiff(diffs))
			Fail("Unexpected diff in generated output")
		}
	},
		Entry("With the full test structure, pulls in references for all substructs", generatorTableInput{
			requestedTypes:         []string{"MyTestStruct"},
			expectedOutputFileName: "testdata/fullMyTestStruct.md",
		}),
		Entry("With only a sub test structure", generatorTableInput{
			requestedTypes:         []string{"SomeSubStruct"},
			expectedOutputFileName: "testdata/someSubStructOnly.md",
		}),
		Entry("With a header file specified, should prefix the generated content", generatorTableInput{
			requestedTypes:         []string{"SomeSubStruct"},
			expectedOutputFileName: "testdata/someSubStructWithHeader.md",
			headerFileName:         "testdata/header.md",
		}),
		Entry("With two unrelated structs", generatorTableInput{
			requestedTypes:         []string{"SomeSubStruct", "AnEmbeddedStruct"},
			expectedOutputFileName: "testdata/unrelatedStructs.md",
			headerFileName:         "testdata/header.md",
		}),
	)
})

// prettyPrintDiff prints the diff for the file out as if it were a git diff.
func prettyPrintDiff(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString("\x1b[32m")
			printDiffLines(&buff, "+ ", text)
			_, _ = buff.WriteString("\x1b[0m")
		case diffmatchpatch.DiffDelete:
			_, _ = buff.WriteString("\x1b[31m")
			printDiffLines(&buff, "- ", text)
			_, _ = buff.WriteString("\x1b[0m")
		case diffmatchpatch.DiffEqual:
			printDiffLines(&buff, "  ", text)
		}
	}

	return buff.String()
}

// printDiffLines prints each line in the diff as a separate line with the given prefix.
func printDiffLines(buff *bytes.Buffer, prefix, in string) {
	in = strings.TrimSuffix(in, "\n")
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		_, _ = buff.WriteString(prefix + line + "\n")
	}
}
