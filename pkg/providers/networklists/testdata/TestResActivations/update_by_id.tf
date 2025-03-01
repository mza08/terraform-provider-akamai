provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_networklist_activations" "test" {
  name                = "Network list test"
  network             = "STAGING"
  notes               = "TEST Notes updated"
  notification_emails = ["user@example.com"]
}

