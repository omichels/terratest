// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2014-04-01/sql"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureSQLDBExample(t *testing.T) {
	t.Parallel()

	expectedResourceGroupName := fmt.Sprintf("sqldatabase-rg-%s", random.UniqueId())

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-sqldb-example",
		Vars: map[string]interface{}{
			"resource_group_name": expectedResourceGroupName,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	expectedSQLServerID := terraform.Output(t, terraformOptions, "sql_server_id")
	expectedSQLServerName := terraform.Output(t, terraformOptions, "sql_server_name")

	expectedSQLServerFullDomainName := terraform.Output(t, terraformOptions, "sql_server_full_domain_name")
	expectedSQLDBName := terraform.Output(t, terraformOptions, "sql_database_name")

	expectedSQLDBID := terraform.Output(t, terraformOptions, "sql_database_id")
	expectedSQLDBStatus := "Online"

	// website::tag::4:: Get the SQL server details and assert them against the terraform output
	actualSQLServerID := azure.GetSQLServerID(t, expectedResourceGroupName, expectedSQLServerName, "")
	actualSQLServerFullDomainName := azure.GetSQLServerFullDomainName(t, expectedResourceGroupName, expectedSQLServerName, "")
	actualSQLServerState := azure.GetSQLServerState(t, expectedResourceGroupName, expectedSQLServerName, "")

	assert.Equal(t, expectedSQLServerID, actualSQLServerID)
	assert.Equal(t, expectedSQLServerFullDomainName, actualSQLServerFullDomainName)
	assert.Equal(t, actualSQLServerState, sql.ServerStateReady)

	// website::tag::5:: Get the SQL server DB details and assert them against the terraform output
	actualSQLDBID := azure.GetDatabaseID(t, expectedResourceGroupName, expectedSQLServerName, expectedSQLDBName, "")
	actualSQLDBStatus := azure.GetDatabaseStatus(t, expectedResourceGroupName, expectedSQLServerName, expectedSQLDBName, "")

	assert.Equal(t, expectedSQLDBID, actualSQLDBID)
	assert.Equal(t, expectedSQLDBStatus, actualSQLDBStatus)

}
