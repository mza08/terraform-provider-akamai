//go:build all || gtm
// +build all gtm

package gtm

import "github.com/akamai/terraform-provider-akamai/v3/pkg/providers/registry"

func init() {
	registry.RegisterProvider(Subprovider())
}
