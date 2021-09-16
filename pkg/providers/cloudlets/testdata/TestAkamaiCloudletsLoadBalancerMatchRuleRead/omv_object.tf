provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_application_load_balancer_match_rule" "test" {
  match_rules {
    name = "rule1"
    start = 10
    end = 10000
    match_url = "example.com"
    matches {
      match_type = "clientip"
      match_value = "127.0.0.1"
      object_match_value {
        type = "simple"
        value = "[\"fghi\"]"
      }
    }
    matches {
      case_sensitive = true
      match_type = "cookie"
      match_value = "cookie=cookievalue"
      object_match_value {
        type = "object"
        name = "abcde"
        name_case_sensitive = true
        name_has_wildcard = false
        options {
          value = "asfas"
          value_has_wildcard = false
          value_case_sensitive = true
          value_escaped = false
        }
      }
    }
    type = "albMatchRule"
  }
  match_rules {
    name = "rule2"
    type = "albMatchRule"
    id = 12333
    aka_rule_id = "abcd"
    forward_settings {
      origin_id = "1234"
    }
  }
}