---
page_title: "Prisma Cloud: prismacloud_integration"
---

# prismacloud_integration

Manages an integration.

## Example Usage

```hcl
resource "prismacloud_integration" "example" {
    name = "SQS"
    integration_type = "amazon_sqs"
    description = "Made by Terraform"
    enabled = true
    integration_config {
        queue_url = "https://sqs.us-east-1.amazonaws.com/12345/url"
    }
}
```

## Argument Reference

* `name` - (Required) Name of the integration.
* `description` - Description.
* `integration_type` - (Required) Integration type (valid values can be found [here](https://api.docs.prismacloud.io/reference#integrations).
* `enabled` - (bool) Enabled.
* `integration_config` - (Required) Integration configuration, the values depend on the integration type, as defined [below](#integration-config).

### Integration Config

Refer to the [Prisma Cloud integration documentation](https://api.docs.prismacloud.io/reference#integration-configuration) if you need more information on a specific integration.

* `queue_url` - The Queue URL you used when you configured Prisma Cloud in Amazon SQS
* `login` - (Qualys/ServiceNow) Login.
* `base_url` - Qualys Security Operations Center server API URL (without "http(s)")
* `password` - (Qualys/ServiceNow) Password
* `host_url` - ServiceNow URL.
* `tables` - (Map of bools) Key/value pairs that identify the ServiceNow module tables with which to integrate (e.g. - incident, sn_si_incident, or em_event).
* `version` - ServiceNow release version.
* `url` - Webhook URL or Splunk HTTP event collector URL.
* `headers` - List of webhook headers, as defined [below](#headers).
* `auth_token` - PagerDuty/Splunk authentication token for the event collector.
* `integration_key` - PagerDuty integration key.
* `source_id` - GCP Source ID for Google CSCC integration.
* `org_id` - GCP Organization ID for Google CSCC integration.
* `account_id` - AWS account ID.
* `regions` - List of AWS regions, as defined [below](#regions).
* `s3_uri` - AWS S3 URI.
* `region` - AWS region.
* `role_arn` - AWS role ARN.
* `external_id` - AWS external ID.
* `roll_up_interval` - (int) File roll up time in minutes for AWS S3 integration. Valid values are `15`, `30`, `60`, or `180`.
* `source_type` - Source type for splunk integration.

### Headers

* `key` - (Required) Header name.
* `value` - (Required) Header value.
* `secure` - (bool) Secure.
* `read_only` - (bool) Read-only.

### Regions

* `name` - AWS region name.
* `api_identifier` - AWS region code.
* `cloud_type` - Cloud Type (default: `aws`).
* `sdk_id` - SDK ID.

## Attribute Reference

* `integration_id` - Integration ID.
* `status` - Status.
* `valid` - (bool) Valid.
* `created_by` - Created by.
* `created_ts` - (int) Created timestamp.
* `last_modified_by` - Last modified by.
* `last_modified_ts` - (int) Last modified timestamp.
* `reason` - Model for the integration status details, as defined [below](#reason).

### Reason

* `last_updated` - (int) Last updated.
* `error_type` - Error type.
* `message` - Message.
* `details` - Model for message details, as defined [below](#details).

### Details

* `status_code` - (int) Status code.
* `subject` - Subject.
* `message` - Internationalization key.
