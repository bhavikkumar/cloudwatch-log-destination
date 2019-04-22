data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    effect = "Allow"

    actions = [
      "sts:AssumeRole"
    ]

    principals {
      type = "Service"

      identifiers = [
        "lambda.amazonaws.com"
      ]
    }
  }
}

data "aws_iam_policy_document" "lambda_write_logs" {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]

    resources = [
      "${aws_cloudwatch_log_group.lambda.arn}"
    ]
  }
}

data "aws_iam_policy_document" "subscription_policy" {
  statement {
    effect = "Allow"

    actions = [
      "logs:PutSubscriptionFilter"
    ]

    resources = [
      "arn:aws:logs:*:*:*"
    ]
  }
}

resource "aws_iam_role" "lambda" {
  name               = "CloudWatchSubscriptionFilterLambda"
  description        = "Used by CloudWatch Destination Lambda"
  assume_role_policy = "${data.aws_iam_policy_document.lambda_assume_role.json}"
  tags               = "${var.tags}"
}

resource "aws_iam_role_policy" "lambda_write_logs" {
  name   = "CloudwatchLogWritePermissions"
  role   = "${aws_iam_role.lambda.name}"
  policy = "${data.aws_iam_policy_document.lambda_write_logs.json}"
}

resource "aws_iam_role_policy" "lambda_subscription_filter_policy" {
  name   = "AllowPutSubscriptionFilterPolicy"
  role   = "${aws_iam_role.lambda.name}"
  policy = "${data.aws_iam_policy_document.subscription_policy.json}"
}

resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${aws_lambda_function.lambda.function_name}"
  retention_in_days = "${var.log_retention_period}"
  kms_key_id        = "${var.kms_key_arn}"
  tags              = "${var.tags}"
}

resource "aws_cloudwatch_log_subscription_filter" "lambda" {
  name            = "default_log_destination"
  log_group_name  = "${aws_cloudwatch_log_group.lambda.name}"
  filter_pattern  = ""
  destination_arn = "${var.log_destination_arn}"
  distribution    = "ByLogStream"
}

resource "aws_lambda_function" "lambda" {
  function_name = "CloudWatchLogDestination"
  description   = "Sets the default cloudwatch log subscription filters"
  role          = "${aws_iam_role.lambda.arn}"
  handler       = "main"
  runtime       = "go1.x"
  memory_size   = 128
  kms_key_arn   = "${var.kms_key_arn}"
  filename      = "cloudwatch-log-destination${var.lambda_version}.zip"

  environment {
    variables = {
      DESTINATION_ARN = "${var.log_destination_arn}"
    }
  }
  tags = "${var.tags}"
}

resource "aws_cloudwatch_event_rule" "subscription_filter" {
  name        = "LogSubscriptionFilterModifications"
  description = "Captures when log groups are created or the subscription filters are modified"
  tags        = "${var.tags}"

  event_pattern = <<PATTERN
{
  "source": [
    "aws.logs"
  ],
  "detail-type": [
    "AWS API Call via CloudTrail"
  ],
  "detail": {
    "eventSource": [
      "logs.amazonaws.com"
    ],
    "eventName": [
      "CreateLogGroup",
      "PutSubscriptionFilter",
      "DeleteSubscriptionFilter"
    ]
  }
}
PATTERN
}

resource "aws_cloudwatch_event_target" "subscription_lambda" {
  rule      = "${aws_cloudwatch_event_rule.subscription_filter.name}"
  arn       = "${aws_lambda_function.lambda.arn}"
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowSubscriptionFilterLambdaExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.lambda.function_name}"
  principal     = "events.amazonaws.com"
  source_arn    = "${aws_cloudwatch_event_rule.subscription_filter.arn}"
}
