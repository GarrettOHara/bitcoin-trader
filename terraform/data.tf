data "aws_caller_identity" "current" {}

data "archive_file" "lambda" {
  type             = "zip"
  source_dir       = "lambda"
  output_file_mode = "0666"
  output_path      = "${path.module}/files/source_code.zip"
}
