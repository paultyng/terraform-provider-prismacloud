---
page_title: "Prisma Cloud: prismacloud_alert_rule"
---

# prismacloud_alert_rule

Manage an alert rule.

## Example Usage

```hcl
resource "prismacloud_alert_rule" "example" {
    name = "My Alert Rule"
    description = "Made by Terraform"
}
```

## Argument Reference

* `name` - (Required) Rule/Scan name
* `description` - Description
* `enabled` - (bool) Enabled (default: `true`)
* `scan_all` - (bool) Scan all policies
* `policies` - List of specific policies to scan
* `policy_labels` - List of policy labels
* `excluded_policies` - List of policies to exclude from scan
* `allow_auto_remediate` - (bool) Allow auto-remediation
* `delay_notification_ms` - (int) Delay notifications by the specified miliseconds
* `notify_on_open` - (bool) Include open alerts in notification (default: `true`)
* `notify_on_snoozed` - (bool) Include snoozed alerts in notification
* `notify_on_dismissed` - (bool) Include dismissed alerts in notification
* `notify_on_resolved` - (bool) Include resolved alerts in notification
* `target` - (Required) Model for the target filter, as defined [below](#target)
* `notification_config` - List of data for notifications to third-party tools, as defined [below](#notification-config)

### Target

There should be one and only one target block:

* `account_groups` - (Required) List of account groups
* `excluded_accounts` - List of excluded accounts
* `regions` - List of regions
* `tags` - List of tag models, as defined [below](#tags)

### Tags

* `key` - Resource tag target
* `values` - List of values for resource tag key

### Notification Config

* `frequency` - Frequency.  Valid values are `as_it_happens`, `daily`, `weekly`, or `monthly`.
* `enabled` - (bool) Scan enabled
* `recipients` - List of unique email addresses to notify (For email notifications), List of integration ids (For integrations without notification templates), or List of notification template ids (For integrations with notification templates)
* `detailed_report` - (bool) Provide CSV detailed report
* `with_compression` - (bool) Compress detailed report
* `include_remediation` - (bool) Include remediation in detailed report
* `config_type` - Config type.  Valid values are `email`, `slack`, `splunk`, `amazon_sqs`, `jira`, `microsoft_teams`, `webhook`, `aws_security_hub`, `google_cscc`, `service_now`, `pager_duty`, `aws_s3`, or `demisto`
* `template_id` - Template ID
* `r_rule_schedule` - R rule schedule

## Attribute Reference

* `policy_scan_config_id` - Policy scan config ID
* `last_modified_on` - (int) Last modified on
* `last_modified_by` - Last modified by
* `owner` - Owner
* `notification_channels` - List of notification channels
* `open_alerts_count` - (int) Open alerts count
* `read_only` - (bool) Read only

In each `notification_config` section, the following attributes are available:

* `config_id` - Alert rule notification config ID
* `last_updated` - (int) Last updated
* `last_sent_ts` - (int) Time of last notification in miliseconds
* `timezone_id` - Timezone ID
* `day_of_month` - (int) Day of month
* `frequency_from_r_rule` - Frequency from R rule
* `hour_of_day` - (int) Hour of day
* `days_of_week` - List of days of week, as defined [below](#days-of-week)

### Days Of Week

* `day` - Day
* `offset` - (int) Offset



## Import

Resources can be imported using the policy scan config ID:

```
$ terraform import prismacloud_alert_rule.example 11111111-2222-3333-4444-555555555555
```
