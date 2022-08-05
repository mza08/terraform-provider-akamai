package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAkamaiApiHostnameCoverage_data_basic(t *testing.T) {
	t.Run("match by ApiHostnameCoverage ID", func(t *testing.T) {
		client := &mockappsec{}

		getCoverageResponse := appsec.GetApiHostnameCoverageResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSApiHostnameCoverage/ApiHostnameCoverage.json"), &getCoverageResponse)

		client.On("GetApiHostnameCoverage",
			mock.Anything,
			appsec.GetApiHostnameCoverageRequest{ConfigID: 0, Version: 0},
		).Return(&getCoverageResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSApiHostnameCoverage/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_hostname_coverage.test", "id", "3644"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
