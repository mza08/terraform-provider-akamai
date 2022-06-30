# minimum basic info B
resource "akamai_iam_user" "test" {
  first_name = "John"
  last_name  = "Smith"
  email      = "jsmith@example.com"
  country    = "country"
  phone      = "(111) 111-1111"
  enable_tfa = false
  lock       = true

  auth_grants_json = "[{\"groupId\":0,\"groupName\":\"group\",\"roleDescription\":\"\",\"roleName\":\"\"}]"
}
