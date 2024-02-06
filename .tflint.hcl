plugin "terraform" {
  enabled = true
  preset  = "all"
}

rule "terraform_documented_variables" {
  enabled = false
}
