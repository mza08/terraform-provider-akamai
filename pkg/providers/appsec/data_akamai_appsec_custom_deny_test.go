package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiCustomDeny_data_basic(t *testing.T) {
	t.Run("match by CustomDeny ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getCustomDenyListResponse := appsec.GetCustomDenyListResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSCustomDeny/CustomDenyList.json"), &getCustomDenyListResponse)
		require.NoError(t, err)

		client.On("GetCustomDenyList",
			mock.Anything,
			appsec.GetCustomDenyListRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_54994"},
		).Return(&getCustomDenyListResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSCustomDeny/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_custom_deny.test", "custom_deny_id", "deny_custom_54994"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
