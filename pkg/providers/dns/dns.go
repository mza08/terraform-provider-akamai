//go:build all || dns
// +build all dns

package dns

import "github.com/akamai/terraform-provider-akamai/v3/pkg/providers/registry"

func init() {
	registry.RegisterProvider(Subprovider())
}
