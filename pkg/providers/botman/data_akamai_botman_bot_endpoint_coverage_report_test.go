package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataBotEndpointCoverageReport(t *testing.T) {
	t.Run("DataBotEndpointCoverageReport", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := botman.GetBotEndpointCoverageReportResponse{
			Operations: []map[string]interface{}{
				{"operationId": "b85e3eaa-d334-466d-857e-33308ce416be", "testKey": "testValue1"},
				{"operationId": "69acad64-7459-4c1d-9bad-672600150127", "testKey": "testValue2"},
				{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				{"operationId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey": "testValue4"},
				{"operationId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey": "testValue5"},
			},
		}
		expectedJSON := `
{
	"operations": [
		{"operationId":"b85e3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"operationId":"69acad64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"operationId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"operationId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`
		mockedBotmanClient.On("GetBotEndpointCoverageReport",
			mock.Anything,
			botman.GetBotEndpointCoverageReportRequest{},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataBotEndpointCoverageReport/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_bot_endpoint_coverage_report.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
	t.Run("DataBotEndpointCoverageReport filter by id", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := botman.GetBotEndpointCoverageReportResponse{
			Operations: []map[string]interface{}{
				{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			},
		}
		expectedJSON := `
{
	"operations":[
		{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}
	]
}`
		mockedBotmanClient.On("GetBotEndpointCoverageReport",
			mock.Anything,
			botman.GetBotEndpointCoverageReportRequest{OperationID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataBotEndpointCoverageReport/filter_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_bot_endpoint_coverage_report.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
	t.Run("DataBotEndpointCoverageReport by config", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := botman.GetBotEndpointCoverageReportResponse{
			Operations: []map[string]interface{}{
				{"operationId": "b85e3eaa-d334-466d-857e-33308ce416be", "testKey": "testValue1"},
				{"operationId": "69acad64-7459-4c1d-9bad-672600150127", "testKey": "testValue2"},
				{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				{"operationId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey": "testValue4"},
				{"operationId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey": "testValue5"},
			},
		}
		expectedJSON := `
{
	"operations": [
		{"operationId":"b85e3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"operationId":"69acad64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"operationId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"operationId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`
		mockedBotmanClient.On("GetBotEndpointCoverageReport",
			mock.Anything,
			botman.GetBotEndpointCoverageReportRequest{ConfigID: 43253, Version: 15},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataBotEndpointCoverageReport/with_config.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_bot_endpoint_coverage_report.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
	t.Run("DataBotEndpointCoverageReport by config filter by id", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := botman.GetBotEndpointCoverageReportResponse{
			Operations: []map[string]interface{}{
				{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			},
		}
		expectedJSON := `
{
	"operations":[
		{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}
	]
}`
		mockedBotmanClient.On("GetBotEndpointCoverageReport",
			mock.Anything,
			botman.GetBotEndpointCoverageReportRequest{ConfigID: 43253, Version: 15, OperationID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataBotEndpointCoverageReport/with_config_filter_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_bot_endpoint_coverage_report.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
