package test

// suites_test go test file is where the test suites are defined.
// The test cases from the 'testcases' go file are identified by their name and
// can be used in the order we want to run the tests. This approach gives more control on the test cases.
// Note:
//  1. The go test file must always end with a suffix '_test'
//  2. The test functions in the go test file must be prefixed with 'Test' to be executed.
import (
	"os"
	"testing"
)

// TestMain function allows us to control four aspects of test execution.
// 1. Setup - not used here, but ideally can call any function before the start of the test
// 2. Run the tests.
// 3. TearDown - calls the TestStatus function at the end of the test.
// 4. Exit behavior - exits the test.
func TestMain(m *testing.M) {
	resultCode := m.Run()
	// TearDown step
	TestStatus()
	// Exit the test
	os.Exit(resultCode)
}

func TestModuleFunctionsSuite(t *testing.T) {

	BeforeSuite(t, " Test the Cloud SQL module functions ")

	t.Run("Deploying SQL module", func(t *testing.T) {
		test := testCases["Deploying SQL module"]
		test = setDefaultVars(test)
		defer func() {
			if r := recover(); r != nil {
				PanicCall(t)
			}
		}()
		Testflow(t, test)
	})

	AfterSuite(t)
}
