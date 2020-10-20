provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  contract_id = "ctr_0"
  group_id    = "grp_0"
  product_id  = "prd_0"

  name = "renamed property"
}
