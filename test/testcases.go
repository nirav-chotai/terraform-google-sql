package test

// testcases go file has the following
//  1. test cases
//  2. default vars for the test cases
//  3. define the global functions
import (
	"testing"
)

// define the default vars for the test cases
func setDefaultVars(test TestCase) TestCase {
	test.Vars["region"] = "asia-northeast1"
	test.Vars["project"] = "my-gcp-project"
	test.Vars["name_prefix"] = "postgres-private"
	test.Vars["master_user_name"] = "admin"
	test.Vars["master_user_password"] = "password"
	GlobalTestMethods = globalTestMethods
	return test
}

// define the global functions
var globalTestMethods = []func(t *testing.T, test TestCase){}

// slice of testcases mapped to a name, which in turn is called from the suites_test go test file.
var testCases = map[string]TestCase{

	"Deploying SQL module": {
		TestCaseDescription:  "Deploying SQL module",
		TerraformDir:         "../examples/postgres-private-ip",
		RunGlobalTestMethods: false,
		RunLocalTestMethods:  true,
		Vars: map[string]interface{}{
			"db_name": "test_postgres",
		},
		LocalTestMethods: []func(t *testing.T, test TestCase){
			testPostgresExistence,
		},
	},
}
