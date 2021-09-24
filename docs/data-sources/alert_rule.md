---
page_title: "Prisma Cloud: prismacloud_alert_rule"
---

# prismacloud_alert_rule

Retrieve information on a specific alert rule.

## Example Usage

```hcl
data "prismacloud_alert_rule" "example" {
    name = "My Alert Rule"
}
```

## Argument Reference

You must specify at least one of the following:

* `policy_scan_config_id` - Policy scan config ID
* `name` - Rule/Scan name

## Attribute Reference

* `description` - Description
* `enabled` - (bool) Enabled
* `scan_all` - (bool) Scan all policies
* `policies` - List of specific policies to scan
* `policy_labels` - List of policy labels
* `excluded_policies` - List of policies to exclude from scan
* `last_modified_on` - (int) Last modified on
* `last_modified_by` - Last modified by
* `allow_auto_remediate` - (bool) Allow auto-remediation
* `delay_notification_ms` - (int) Delay notifications by the specified miliseconds
* `notify_on_open` - (bool) Include open alerts in notification
* `notify_on_snoozed` - (bool) Include snoozed alerts in notification
* `notify_on_dismissed` - (bool) Include dismissed alerts in notification
* `notify_on_resolved` - (bool) Include resolved alerts in notification
* `owner` - Owner
* `notification_channels` - List of notification channels
* `open_alerts_count` - (int) Open alerts count
* `read_only` - (bool) Read only
* `target` - Model for the target filter, as defined [below](#target)
* `notification_config` - List of data for notifications to third-party tools, as defined [below](#notification-config)

### Target

* `account_groups` - List of account groups
* `excluded_accounts` - List of excluded accounts
* `regions` - List of regions
* `tags` - List of TargetTag objects, as defined [below](#tags)

### Tags

* `key` - Resource tag target
* `values` - List of values for resource tag key

### Notification Config

* `config_id` - Alert rule notification config ID
* `frequency` - Frequency
* `enabled` - (bool) Scan enabled
* `recipients` - List of unique email addresses to notify (For email notifications), List of integration ids (For integrations without notification templates), or List of notification template ids (For integrations with notification templates)
* `detailed_report` - (bool) Provide CSV detailed report
* `with_compression` - (bool) Compress detailed report
* `include_remediation` - (bool) Include remediation in detailed report
* `last_updated` - (int) Last updated
* `last_sent_ts` - (int) Time of last notification in miliseconds
* `config_type` - Config type
* `template_id` - Template ID
* `timezone_id` - Timezone ID
* `day_of_month` - (int) Day of month
* `r_rule_schedule` - R rule schedule
* `frequency_from_r_rule` - Frequency from R rule
* `hour_of_day` - (int) Hour of day
* `days_of_week` - List of days of week, as defined [below](#days-of-week)

### Days Of Week

* `day` - Day
* `offset` - (int) Offset
