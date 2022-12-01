provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_javascript_injection" "test" {
  config_id            = 43253
  security_policy_id   = "AAAA_81230"
  javascript_injection = <<-EOF
{
  "testKey": "testValue3"
}
EOF
}