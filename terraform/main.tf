resource "aws_iam_role" "lambda" {
  name               = "bitcoin-trader-lambda"
  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "sts:AssumeRole",
            "Principal": {
                "Service": "lambda.amazonaws.com"
            }
        }
    ]
}
EOF
}

resource "aws_iam_role_policy" "lambda" {
  name   = var.name
  role   = aws_iam_role.lambda.id
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "logs:CreateLogGroup",
            "Resource": "arn:aws:logs:${var.aws_region}:${data.aws_caller_identity.current.id}:*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "arn:aws:logs:${var.aws_region}:${data.aws_caller_identity.current.id}:log-group:/aws/lambda/*"
            ]
        }
    ]
}
EOF
}

# data "archive_file" "archive" {
#   type        = "zip"
#   source_file = "../lambda/lambda_function.py"
#   output_path = "source_code.zip"
# }

resource "aws_lambda_function" "lambda" {
  # checkov:skip=CKV_AWS_50: No x-ray tracing
  # checkov:skip=CKV_AWS_116: No DLQ
  # checkov:skip=CKV_AWS_115: No concurrent executions
  # checkov:skip=CKV_AWS_117: Not inside VPC
  # checkov:skip=CKV_AWS_272: No code signing
  filename      = "source_code.zip"
  function_name = var.name
  role          = aws_iam_role.lambda.arn
  handler       = "lambda_handler"
  runtime       = "python3.9"

  source_code_hash = data.archive_file.lambda.output_base64sha256
}

resource "aws_iam_role" "event_bridge_scheduler" {
  name               = "${var.name}-event-bridge-scheduler-role"
  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "scheduler.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
EOF
}

resource "aws_iam_policy" "event_bridge_policy" {
  name        = "${var.name}-event-bridge-scheduler-policy"
  description = "Grants permission for Eventbridge Scheduler resource to trigger Step Function."

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "lambda:InvokeFunction"
            ],
            "Resource": [
                "${aws_lambda_function.lambda.arn}"
            ]
        }
    ]
}
EOF
}

resource "aws_scheduler_schedule_group" "schedule_group" {
  name = "${var.name}-group"
}

resource "aws_scheduler_schedule" "schedule" {
  # checkov:skip=CKV_AWS_297: Not using CMK
  name       = "${var.name}-schedule"
  group_name = aws_scheduler_schedule_group.schedule_group

  # Schedule expression translates to every day at 16:00 UTC / 08:00 PST
  schedule_expression = "cron(00 16 * * ? *)"

  flexible_time_window {
    mode = var.flexible_time_window
  }

  target {
    arn      = aws_lambda_function.lambda.arn
    role_arn = aws_iam_role.event_bridge_scheduler.arn
  }
}
