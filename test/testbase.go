package test

// testbase is the base of the test. This file controls what needs to be run before test suite or a test case and
// keeps tabs on the test suite, test case, test method counters.
import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/test-structure"
	"testing"
)

// suitesCount gives the total number of test suites that has run in the test.
// testCaseCount gives the total number of test cases that has run in the test
// methodsCount gives the total number of test methods that has run in the test.
// methodsPassedCount gives the total number of test methods that has passed in the test.
// methodsFailedCount gives the total number of test methods that has failed in the test.
// testCaseCounterInSuites gives the total number of test cases that has run in the test suite.
// methodsCounterInTestCases gives the total number of test methods that has run in the test case.
var counters = struct {
	suitesCount               int
	testCaseCount             int
	methodsCount              int
	methodsPassedCount        int
	methodsFailedCount        int
	testCaseCounterInSuites   int
	methodsCounterInTestCases int
}{
	suitesCount:               0,
	testCaseCount:             0,
	methodsCount:              0,
	methodsPassedCount:        0,
	methodsFailedCount:        0,
	testCaseCounterInSuites:   0,
	methodsCounterInTestCases: 0,
}

// TestCase struct is where you define your test cases
type TestCase struct {
	TestCaseDescription  string
	TerraformDir         string
	CopyFileName         string
	RunGlobalTestMethods bool
	RunLocalTestMethods  bool
	Vars                 map[string]interface{}
	LocalTestMethods     []func(t *testing.T, test TestCase)
}

var GlobalTestMethods = []func(t *testing.T, test TestCase){}

var errorStatus error

// TestPassed is an indicator to determine whether the test has passed or failed.
var TestPassed = true

// BeforeSuite function will hold tasks to be done before a test suite gets executed.
func BeforeSuite(t *testing.T, suiteDescription string) {
	fmt.Println("")
	fmt.Println("")
	counters.testCaseCounterInSuites = 0
	counters.suitesCount++
	logger.Logf(t, "START TEST SUITE %d: - %s", counters.suitesCount, suiteDescription)
}

// AfterSuite function will hold tasks to be done at the end of a test suite.
func AfterSuite(t *testing.T) {
	logger.Logf(t, "END TEST SUITE %d", counters.suitesCount)
}

func beforeTestCase(t *testing.T, test TestCase) {
	test_structure.RunTestStage(t, "deploy", func() {
		counters.testCaseCount++
		counters.testCaseCounterInSuites++
		logger.Logf(t, "START TEST CASE %d: %s", counters.testCaseCounterInSuites, test.TestCaseDescription)
		logger.Log(t, "Before Test Case: Deploy Terraform and load its contents to be used in the test cases")
		terraformOptions := CreateTerraformOptions(t, test)
		test_structure.SaveTerraformOptions(t, test.TerraformDir, terraformOptions)
		_, err := terraform.InitAndApplyE(t, terraformOptions)
		if err != nil {
			logger.Logf(t, "The test have failed in the deploy stage. The test will end now.")
			CheckPanic(t, err, fmt.Sprintf("Failure trace: %s", WhereAmI()), fmt.Sprintf("The test have failed in the deploy stage. The test will end now"))
		}
		fmt.Println("")
	})
}

func executeTestMethods(t *testing.T, test TestCase) {
	counters.methodsCounterInTestCases = 0
	test_structure.RunTestStage(t, "verify", func() {
		if test.RunLocalTestMethods {
			for _, fn := range test.LocalTestMethods {
				beforeTestMethod(t)
				fn(t, test)
				afterTestMethod(t)
			}
		}
		if test.RunGlobalTestMethods {
			for _, fn := range GlobalTestMethods {
				beforeTestMethod(t)
				fn(t, test)
				afterTestMethod(t)
			}
		}
	})
}

func afterTestCase(t *testing.T, test TestCase) {
	test_structure.RunTestStage(t, "destroy", func() {
		fmt.Println("")
		logger.Logf(t, "After Test Case: destroy terraform and all the resources it created.")
		terraformOptions := test_structure.LoadTerraformOptions(t, test.TerraformDir)
		_, err := terraform.DestroyE(t, terraformOptions)
		CheckPanic(t, err, fmt.Sprintf("Failure trace: %s", WhereAmI()), fmt.Sprintf("The test have failed in the destroy stage. The test will end now"))
		logger.Logf(t, "END TEST CASE %d: ", counters.testCaseCounterInSuites)
	})
}

func beforeTestMethod(t *testing.T) {
	errorStatus = nil
	counters.methodsCount++
	counters.methodsCounterInTestCases++
	logger.Logf(t, "** START TEST METHOD %d: ", counters.methodsCounterInTestCases)
}

func afterTestMethod(t *testing.T) {
	if errorStatus != nil {
		counters.methodsFailedCount++
	} else {
		counters.methodsPassedCount++
	}
	logger.Logf(t, "** END TEST METHOD %d : %s", counters.methodsCounterInTestCases, testMethodStatus(errorStatus))
	fmt.Println("======================================================================================")
}

// Testflow function will execute the test cases listed under the suites in 'suites_test' go test file.
func Testflow(t *testing.T, test TestCase) {
	defer afterTestCase(t, test)
	beforeTestCase(t, test)
	executeTestMethods(t, test)
	fmt.Println("")
}

// TestStatus function will report the status of the total test suites, total test cases, total test methods (pass & fail) and the test status.
func TestStatus() {
	fmt.Println("")
	fmt.Println("=============================================================================================")
	fmt.Printf("TOTAL TEST SUITES : %d\n", counters.suitesCount)
	fmt.Printf("TOTAL TEST CASES : %d\n", counters.testCaseCount)
	fmt.Printf("TOTAL TEST METHODS : %d\n", counters.methodsCount)
	fmt.Printf("TOTAL TEST METHODS PASS : %d\n", counters.methodsPassedCount)
	fmt.Printf("TOTAL TEST METHODS FAIL : %d\n", counters.methodsFailedCount)
	fmt.Println("")
	if TestPassed {
		fmt.Println("TEST STATUS : PASS")
	} else {
		fmt.Println("TEST STATUS : FAIL")
	}
	fmt.Println("=============================================================================================")
	fmt.Println("")
}

// CreateTerraformOptions will update the terraformOptions for terraform apply in the testbase
func CreateTerraformOptions(t *testing.T, test TestCase) *terraform.Options {
	terraformOptions := &terraform.Options{
		TerraformDir: test.TerraformDir,
		Vars:         test.Vars,
		EnvVars: map[string]string{
			"GOOGLE_PROJECT": fmt.Sprintf("%v", test.Vars["project"]),
		},
	}
	return terraformOptions
}
