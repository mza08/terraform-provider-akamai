package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigurationCreate,
		ReadContext:   resourceConfigurationRead,
		UpdateContext: resourceConfigurationUpdate,
		DeleteContext: resourceConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contract_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"host_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_from_config_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"create_from_version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"config_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationCreate")
	logger.Debug("in resourceConfigurationCreate")

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	description, err := tools.GetStringValue("description", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := tools.GetIntValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	hostnameset, err := tools.GetSetValue("host_names", d)
	if err != nil {
		return diag.FromErr(err)
	}
	hostnames := make([]string, 0, len(hostnameset.List()))
	for _, h := range hostnameset.List() {
		hostnames = append(hostnames, h.(string))
	}
	createFromConfigID, err := tools.GetIntValue("create_from_config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createFromVersion, err := tools.GetIntValue("create_from_version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	if createFromVersion > 0 && createFromConfigID > 0 {
		createConfigurationClone := appsec.CreateConfigurationCloneRequest{
			Name:        name,
			Description: description,
			ContractID:  contractID,
			GroupID:     groupID,
			Hostnames:   hostnames,
		}
		createConfigurationClone.CreateFrom.ConfigID = createFromConfigID
		createConfigurationClone.CreateFrom.Version = createFromVersion

		response, err := client.CreateConfigurationClone(ctx, createConfigurationClone)
		if err != nil {
			logger.Errorf("calling 'createConfigurationClone': %s", err.Error())
			return diag.FromErr(err)
		}
		if err := d.Set("config_id", response.ConfigID); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}

		d.SetId(fmt.Sprintf("%d", response.ConfigID))

	} else {
		createConfiguration := appsec.CreateConfigurationRequest{
			Name:        name,
			Description: description,
			ContractID:  contractID,
			GroupID:     groupID,
			Hostnames:   hostnames,
		}

		response, err := client.CreateConfiguration(ctx, createConfiguration)
		if err != nil {
			logger.Errorf("calling 'createConfiguration': %s", err.Error())
			return diag.FromErr(err)
		}
		if err := d.Set("config_id", response.ConfigID); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
		d.SetId(fmt.Sprintf("%d", response.ConfigID))
	}

	return resourceConfigurationRead(ctx, d, m)
}

func resourceConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRead")
	logger.Debug("in resourceConfigurationRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getConfiguration := appsec.GetConfigurationRequest{
		ConfigID: configID,
	}

	configuration, err := client.GetConfiguration(ctx, getConfiguration)
	if err != nil {
		logger.Errorf("calling 'getConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}

	if err = d.Set("name", configuration.Name); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err = d.Set("description", configuration.Description); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err = d.Set("config_id", configuration.ID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	getSelectedHostnamesRequest := appsec.GetSelectedHostnamesRequest{
		ConfigID: configID,
		Version:  version,
	}

	selectedhostnames, err := client.GetSelectedHostnames(ctx, getSelectedHostnamesRequest)
	if err != nil {
		logger.Errorf("calling 'getSelectedHostname': %s", err.Error())
		return diag.FromErr(err)
	}
	selectedhostnameset := schema.Set{F: schema.HashString}
	for _, hostname := range selectedhostnames.HostnameList {
		selectedhostnameset.Add(hostname.Hostname)
	}

	if err = d.Set("host_names", selectedhostnameset.List()); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationUpdate")
	logger.Debug("in resourceConfigurationUpdate")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	description, err := tools.GetStringValue("description", d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateConfiguration := appsec.UpdateConfigurationRequest{
		ConfigID:    configID,
		Name:        name,
		Description: description,
	}

	_, err = client.UpdateConfiguration(ctx, updateConfiguration)
	if err != nil {
		logger.Errorf("calling 'updateConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}

	if d.HasChange("host_names") {
		hostnameset, err := tools.GetSetValue("host_names", d)
		if err != nil {
			return diag.FromErr(err)
		}
		hostnamelist := tools.SetToStringSlice(hostnameset)
		hostnames := make([]appsec.Hostname, 0, len(hostnamelist))
		for _, name := range hostnamelist {
			hostname := appsec.Hostname{Hostname: name}
			hostnames = append(hostnames, hostname)
		}

		version, err := getModifiableConfigVersion(ctx, configID, "configuration", m)
		if err != nil {
			return diag.FromErr(err)
		}
		updateSelectedHostnames := appsec.UpdateSelectedHostnamesRequest{
			ConfigID:     configID,
			Version:      version,
			HostnameList: hostnames,
		}

		_, err = client.UpdateSelectedHostnames(ctx, updateSelectedHostnames)
		if err != nil {
			logger.Errorf("calling 'UpdateSelectedHostnames': %s", err.Error())
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}

	return resourceConfigurationRead(ctx, d, m)
}

func resourceConfigurationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationDelete")
	logger.Debug("in resourceConfigurationDelete")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Check whether any versions of this config have ever been activated
	getConfigVersionsRequest := appsec.GetConfigurationVersionsRequest{
		ConfigID: configID,
	}

	configurationVersions, err := client.GetConfigurationVersions(ctx, getConfigVersionsRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, configVersion := range configurationVersions.VersionList {
		if configVersion.Production.Status != "Inactive" || configVersion.Staging.Status != "Inactive" {
			return diag.Errorf("cannot delete configuration '%s' as version %d has been active in staging or production",
				configurationVersions.ConfigName, configVersion.Version)
		}
	}

	removeConfiguration := appsec.RemoveConfigurationRequest{
		ConfigID: configID,
	}

	_, err = client.RemoveConfiguration(ctx, removeConfiguration)
	if err != nil {
		logger.Errorf("calling 'removeConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
