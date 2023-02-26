resource "aws_ssm_parameter" "param1" {
  name  = "string-test"
  type  = "String"
  value = "string-test-1"
}

resource "aws_ssm_parameter" "param2" {
  name  = "string-list-test"
  type  = "StringList"
  value = "stringlist-test-1"
}

resource "aws_ssm_parameter" "param3" {
  name  = "secure-string-test"
  type  = "SecureString"
  value = "securestring-test"
}
