package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"sigs.k8s.io/multi-tenancy/benchmarks/kubectl-mtb/pkg/benchmark"
	benchmarksuite "sigs.k8s.io/multi-tenancy/benchmarks/kubectl-mtb/pkg/benchmark_suite"
)

// DefaultReporter collects all the test summaries
type DefaultReporter struct {
	testSummaries []*TestSummary
}

// NewDefaultReporter returns the pointer of DefaultReporter
func NewDefaultReporter() *DefaultReporter {
	return &DefaultReporter{}
}

// SuiteWillBegin prints banner and total benchmarks to be run
func (r *DefaultReporter) SuiteWillBegin(suiteSummary *SuiteSummary) {
	writer.PrintBanner(writer.Colorize(boldStyle, "%s", suiteSummary.Suite.Title), "=")
	writer.Println(0, writer.Colorize(lightGrayColor, "Running %d of %d as %s in namespace %s",
		suiteSummary.NumberOfTotalTests, suiteSummary.Suite.Totals(), suiteSummary.User, suiteSummary.Namespace))
}

// TestWillRun prints each test status
func (r *DefaultReporter) TestWillRun(testSummary *TestSummary) {
	if testSummary.Validation {
		writer.Print(0, writer.Colorize(cyanColor, "[%s] [%s] ", testSummary.Benchmark.ID, testSummary.Benchmark.Category))
		writer.Println(0, testSummary.Benchmark.Title)
		writer.Println(0, writer.Colorize(grayColor, "%s", testSummary.Benchmark.Description))
		if testSummary.Test {
			passed := "Passed " + tick
			writer.Println(0, writer.Colorize(greenColor, passed))
		} else {
			failed := "Failed " + cross
			writer.Println(0, writer.Colorize(redColor, failed))
			writer.Print(0, writer.Colorize(lilac, "Remediation: "))
			writer.Println(0, writer.Colorize(lightGrayColor, testSummary.Benchmark.Remediation))

		}
		writer.PrintBanner(writer.Colorize(grayColor, "Completed in %v", testSummary.RunTime), "-")
		return
	}
	preRunfmt := writer.Colorize(magentaColor, "[PreRun-Validation Error]")
	errormsg := writer.Colorize(redColor, testSummary.ValidationError.Error())
	bannerText := fmt.Sprintf("%s [%s] %s: %s %s", preRunfmt, testSummary.Benchmark.ID, testSummary.Benchmark.Title, errormsg, cross)
	writer.PrintBanner(bannerText, "-")
	r.testSummaries = append(r.testSummaries, testSummary)
}

// SuiteDidEnd prints end result summary of benchmark suite
func (r *DefaultReporter) SuiteDidEnd(suiteSummary *SuiteSummary) {
	writer.Print(0, writer.Colorize(greenColor, "%d Passed | ", suiteSummary.NumberOfPassedTests))
	writer.Print(0, writer.Colorize(redColor, "%d Failed | ", suiteSummary.NumberOfFailedTests))
	writer.Print(0, writer.Colorize(yellowColor, "%d Skipped | ", suiteSummary.NumberOfSkippedTests))
	writer.Print(0, writer.Colorize(magentaColor, "%d Errors | ", suiteSummary.NumberOfFailedValidations))
	writer.PrintNewLine()
	writer.PrintBanner(writer.Colorize(grayColor, "Completed in %v", suiteSummary.RunTime), "=")

	printScoreCard(benchmarksuite.SortedBenchmarks)

	//if sonobuoyReport != "" {  // Can have flag
	sonobuoyReport := "/tmp/mtlb.json"
	f, _ := os.Create(sonobuoyReport)
	err := writeSonobuoyReport(f, benchmarksuite.SortedBenchmarks)
	if err != nil {
		fmt.Println("Error", err)
	}
	//}

}

// FullSummary prints end result of all the tests at one place.
func printScoreCard(benchmarks []*benchmark.Benchmark) {
	data := [][]string{}
	counter := 0

	for _, b := range benchmarks {
		counter++
		var status string

		switch b.Status {
		case "Error":
			status = writer.Colorize(magentaColor, "Error")
		case "Pass":
			status = writer.Colorize(greenColor, "Passed")
		case "Fail":
			status = writer.Colorize(redColor, "Failed")
		default:
			status = writer.Colorize(yellowColor, "Skipped")
		}

		testName := b.Title
		result := []string{strconv.Itoa(counter), b.ID, testName, status}
		data = append(data, result)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"No.", "ID", "Test", "Result"})
	table.SetAutoWrapText(false)

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
	writer.PrintBanner(writer.Colorize(defaultStyle, ""), "=")
}

type SonobuoyResultsItem struct {
	Name     string                `json:"name" yaml:"name"`
	Status   string                `json:"status" yaml:"status,omitempty"`
	Metadata map[string]string     `json:"meta,omitempty" yaml:"meta,omitempty"`
	Details  map[string]string     `json:"details,omitempty" yaml:"details,omitempty"`
	Items    []SonobuoyResultsItem `json:"items,omitempty" yaml:"items,omitempty"`
}

func createSonobuoyResultsForResult(benchmarks []*benchmark.Benchmark) []SonobuoyResultsItem {
	var results = []SonobuoyResultsItem{}

	for _, test := range benchmarks {
		resultItem := SonobuoyResultsItem{
			Name:   test.Title,
			Status: test.Status,
			Details: map[string]string{
				"description": test.Description,
				"category":    test.Category,
				"remediation": test.Remediation,
			},
		}
		results = append(results, resultItem)
	}
	return results
}

func writeSonobuoyReport(w io.Writer, benchmarks []*benchmark.Benchmark) error {
	sonobuoyResults := SonobuoyResultsItem{
		Name:   "mtb-benchmark",
		Status: "complete",
	}
	//for _, result := range *r {
	results := createSonobuoyResultsForResult(benchmarks)
	sonobuoyResults.Items = append(sonobuoyResults.Items, results...)
	//}
	j, err := json.Marshal(sonobuoyResults)
	if err != nil {
		return err
	}
	_, err = w.Write(j)
	return err
}
