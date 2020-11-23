package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"testing"
)

// testPostgresExistence will verify whether the tags are added to the vpc by msg-vpc module
func testPostgresExistence(t *testing.T, test TestCase) {
	logger.Logf(t, "Test Method Name: %s", WhereAmI())
	workingDir := test.TerraformDir
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	region := fmt.Sprintf("%s", terraformOptions.Vars["region"])
	projectId := fmt.Sprintf("%s", terraformOptions.Vars["project"])
	dbName := fmt.Sprintf("%s", terraformOptions.Vars["db_name"])

	instanceNameFromOutput := terraform.Output(t, terraformOptions, "master_instance_name")
	ipAddressesFromOutput := terraform.Output(t, terraformOptions, "master_ip_addresses")
	privateIPFromOutput := terraform.Output(t, terraformOptions, "master_private_ip")

	assert.Contains(t, ipAddressesFromOutput, "PRIVATE", "IP Addresses output has to contain 'PRIVATE'")
	assert.Contains(t, ipAddressesFromOutput, privateIPFromOutput, "IP Addresses output has to contain 'private_ip' from output")

	dbNameFromOutput := terraform.Output(t, terraformOptions, "db_name")
	proxyConnectionFromOutput := terraform.Output(t, terraformOptions, "master_proxy_connection")

	expectedDBConn := fmt.Sprintf("%s:%s:%s", projectId, region, instanceNameFromOutput)

	assert.Equal(t, dbName, dbNameFromOutput)
	assert.Equal(t, expectedDBConn, proxyConnectionFromOutput)
}
