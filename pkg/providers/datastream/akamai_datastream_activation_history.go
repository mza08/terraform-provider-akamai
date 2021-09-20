package datastream

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/datastream"
)

func dataAkamaiDatastreamActivationHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAkamaiDatastreamActivationHistoryRead,
		Schema: map[string]*schema.Schema{
			"stream_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"activations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stream_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"stream_version_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func populateSchemaFieldsWithActivationHistory(ac []datastream.ActivationHistoryEntry, d *schema.ResourceData, streamID int) error {

	var activations []map[string]interface{}
	for _, a := range ac {
		v := map[string]interface{}{
			"stream_id":         a.StreamID,
			"stream_version_id": a.StreamVersionID,
			"created_by":        a.CreatedBy,
			"created_date":      a.CreatedDate,
			"is_active":         a.IsActive,
		}
		activations = append(activations, v)
	}

	fields := map[string]interface{}{
		"stream_id":   streamID,
		"activations": activations,
	}

	err := tools.SetAttrs(d, fields)
	if err != nil {
		return fmt.Errorf("could not set schema attributes: %s", err)
	}

	return nil
}

func dataAkamaiDatastreamActivationHistoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("DataStream", "dataAkamaiDatastreamActivationHistoryRead")
	client := inst.Client(meta)

	streamID, err := tools.GetIntValue("stream_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Debug("Getting activation history")
	activationHistory, err := client.GetActivationHistory(ctx, datastream.GetActivationHistoryRequest{
		StreamID: int64(streamID),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	err = populateSchemaFieldsWithActivationHistory(activationHistory, d, streamID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", streamID))

	return nil
}
