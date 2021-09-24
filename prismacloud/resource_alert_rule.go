package prismacloud

import (
	"log"

	pc "github.com/paloaltonetworks/prisma-cloud-go"
	"github.com/paloaltonetworks/prisma-cloud-go/alert/rule"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlertRule() *schema.Resource {
	return &schema.Resource{
		Create: createAlertRule,
		Read:   readAlertRule,
		Update: updateAlertRule,
		Delete: deleteAlertRule,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"policy_scan_config_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy scan config ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule/Scan name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled",
				Default:     true,
			},
			"scan_all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Scan all policies",
			},
			"policies": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of specific policies to scan",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"policy_labels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Policy labels",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"excluded_policies": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of policies to exclude from scan",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_modified_on": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Last modified on",
			},
			"last_modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified by",
			},
			"allow_auto_remediate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow auto-remediation",
			},
			"delay_notification_ms": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Delay notifications by the specified milliseconds",
			},
			"notify_on_open": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Include open alerts in notification",
				Default:     true,
			},
			"notify_on_snoozed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Include snoozed alerts in notification",
			},
			"notify_on_dismissed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Include dismissed alerts in notification",
			},
			"notify_on_resolved": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Include resolved alerts in notification",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Owner",
			},
			"notification_channels": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of notification channels",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"open_alerts_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Open alerts count",
			},
			"read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Read only",
			},
			"target": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "Model for the target filter",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_groups": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "List of account groups",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"excluded_accounts": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "List of excluded accounts",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"regions": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "List of regions",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of TargetTag objects",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource tag target",
									},
									"values": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of values for resource tag key",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"notification_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of data for notifications to third-party tools",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alert rule notification config ID",
						},
						"frequency": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Frequency",
							ValidateFunc: validation.StringInSlice(
								[]string{
									rule.FrequencyAsItHappens,
									rule.FrequencyDaily,
									rule.FrequencyWeekly,
									rule.FrequencyMonthly,
								},
								false,
							),
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Scan enabled",
						},
						"recipients": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "List of unique email addresses to notify",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"detailed_report": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Provide CSV detailed report",
						},
						"with_compression": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Compress detailed report",
						},
						"include_remediation": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Include remediation in detailed report",
						},
						"last_updated": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Last updated",
						},
						"last_sent_ts": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time of last notification in milliseconds",
						},
						"config_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Config type",
							ValidateFunc: validation.StringInSlice(
								[]string{
									rule.TypeEmail,
									rule.TypeSlack,
									rule.TypeSplunk,
									rule.TypeAmazonSqs,
									rule.TypeJira,
									rule.TypeMicrosoftTeams,
									rule.TypeWebhook,
									rule.TypeAwsSecurityHub,
									rule.TypeGoogleCscc,
									rule.TypeServiceNow,
									rule.TypePagerDuty,
									rule.TypeDemisto,
									rule.TypeAwsS3,
								},
								false,
							),
						},
						"template_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Template ID",
						},
						"timezone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timezone ID",
						},
						"day_of_month": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Day of month",
						},
						"r_rule_schedule": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "R rule schedule",
						},
						"frequency_from_r_rule": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Frequency from R rule",
						},
						"hour_of_day": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Hour of day",
						},
						"days_of_week": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Days of week",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Day",
									},
									"offset": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Offset",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func parseAlertRule(d *schema.ResourceData, id string) rule.Rule {
	tgt := ResourceDataInterfaceMap(d, "target")

	accountGroups := []string{}
	if tg := tgt["account_groups"]; tg != nil {
		accountGroups = SetToStringSlice(tg.(*schema.Set))
	}

	var excludedAccounts []string
	if ea := tgt["excluded_accounts"]; ea != nil {
		excludedAccounts = SetToStringSlice(ea.(*schema.Set))
	}

	var regions []string
	if r := tgt["regions"]; r != nil {
		regions = SetToStringSlice(r.(*schema.Set))
	}

	var tags []rule.Tag
	if tt := tgt["tags"]; tt != nil && len(tt.([]interface{})) > 0 {
		tlist := tt.([]interface{})
		tags = make([]rule.Tag, 0, len(tlist))
		for _, item := range tlist {
			tmap := item.(map[string]interface{})
			tags = append(tags, rule.Tag{
				Key:    tmap["key"].(string),
				Values: ListToStringSlice(tmap["values"].([]interface{})),
			})
		}
	}

	ans := rule.Rule{
		PolicyScanConfigId: id,
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Enabled:            d.Get("enabled").(bool),
		ScanAll:            d.Get("scan_all").(bool),
		Policies:           SetToStringSlice(d.Get("policies").(*schema.Set)),
		PolicyLabels:       SetToStringSlice(d.Get("policy_labels").(*schema.Set)),
		ExcludedPolicies:   SetToStringSlice(d.Get("excluded_policies").(*schema.Set)),
		Target: rule.Target{
			AccountGroups:    accountGroups,
			ExcludedAccounts: excludedAccounts,
			Regions:          regions,
			Tags:             tags,
		},
		AllowAutoRemediate:  d.Get("allow_auto_remediate").(bool),
		DelayNotificationMs: d.Get("delay_notification_ms").(int),
		NotifyOnOpen:        d.Get("notify_on_open").(bool),
		NotifyOnSnoozed:     d.Get("notify_on_snoozed").(bool),
		NotifyOnDismissed:   d.Get("notify_on_dismissed").(bool),
		NotifyOnResolved:    d.Get("notify_on_resolved").(bool),
	}
	if ans.Policies == nil {
		ans.Policies = []string{}
	}
	if ans.PolicyLabels == nil {
		ans.PolicyLabels = []string{}
	}
	if ans.ExcludedPolicies == nil {
		ans.ExcludedPolicies = []string{}
	}

	ncl := d.Get("notification_config").([]interface{})
	if len(ncl) > 0 {
		ans.NotificationConfig = make([]rule.NotificationConfig, 0, len(ncl))
		for _, i := range ncl {
			nc := i.(map[string]interface{})
			var days []rule.Day

			if dl := nc["days_of_week"]; dl != nil && len(dl.([]interface{})) != 0 {
				dil := dl.([]interface{})
				days = make([]rule.Day, 0, len(dil))
				for _, j := range dil {
					di := j.(map[string]interface{})
					days = append(days, rule.Day{
						Day:    di["day"].(string),
						Offset: di["offset"].(int),
					})
				}
			}

			ans.NotificationConfig = append(ans.NotificationConfig, rule.NotificationConfig{
				Frequency:          nc["frequency"].(string),
				Enabled:            nc["enabled"].(bool),
				Recipients:         SetToStringSlice(nc["recipients"].(*schema.Set)),
				DetailedReport:     nc["detailed_report"].(bool),
				WithCompression:    nc["with_compression"].(bool),
				IncludeRemediation: nc["include_remediation"].(bool),
				Type:               nc["config_type"].(string),
				TemplateId:         nc["template_id"].(string),
				RruleSchedule:      nc["r_rule_schedule"].(string),
			})

		}
	}

	return ans
}

func saveAlertRule(d *schema.ResourceData, o rule.Rule) {
	var err error

	d.Set("policy_scan_config_id", o.PolicyScanConfigId)
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("enabled", o.Enabled)
	d.Set("scan_all", o.ScanAll)
	if err = d.Set("policies", o.Policies); err != nil {
		log.Printf("[WARN] Error setting 'policies' for %q: %s", d.Id(), err)
	}
	if err = d.Set("policy_labels", o.PolicyLabels); err != nil {
		log.Printf("[WARN] Error setting 'policy_labels' for %q: %s", d.Id(), err)
	}
	if err = d.Set("excluded_policies", o.ExcludedPolicies); err != nil {
		log.Printf("[WARN] Error setting 'excluded_policies' for %q: %s", d.Id(), err)
	}
	d.Set("last_modified_on", o.LastModifiedOn)
	d.Set("last_modified_by", o.LastModifiedBy)
	d.Set("allow_auto_remediate", o.AllowAutoRemediate)
	d.Set("delay_notification_ms", o.DelayNotificationMs)
	d.Set("notify_on_open", o.NotifyOnOpen)
	d.Set("notify_on_snoozed", o.NotifyOnSnoozed)
	d.Set("notify_on_dismissed", o.NotifyOnDismissed)
	d.Set("notify_on_resolved", o.NotifyOnResolved)
	d.Set("owner", o.Owner)
	if err = d.Set("notification_channels", o.NotificationChannels); err != nil {
		log.Printf("[WARN] Error setting 'notification_channels' for %q: %s", d.Id(), err)
	}
	d.Set("open_alerts_count", o.OpenAlertsCount)
	d.Set("read_only", o.ReadOnly)

	tgt := map[string]interface{}{
		"account_groups":    o.Target.AccountGroups,
		"excluded_accounts": o.Target.ExcludedAccounts,
		"regions":           o.Target.Regions,
	}
	if len(o.Target.Tags) == 0 {
		tgt["tags"] = nil
	} else {
		tl := make([]interface{}, 0, len(o.Target.Tags))
		for _, tag := range o.Target.Tags {
			tl = append(tl, map[string]interface{}{
				"key":    tag.Key,
				"values": tag.Values,
			})
		}
		tgt["tags"] = tl
	}
	if err = d.Set("target", []interface{}{tgt}); err != nil {
		log.Printf("[WARN] Error setting 'target' for %q: %s", d.Id(), err)
	}

	if len(o.NotificationConfig) == 0 {
		d.Set("notification_config", nil)
	} else {
		ncList := make([]interface{}, 0, len(o.NotificationConfig))

		for _, nc := range o.NotificationConfig {
			var days []interface{}
			if len(nc.DaysOfWeek) != 0 {
				days = make([]interface{}, 0, len(nc.DaysOfWeek))
				for _, dow := range nc.DaysOfWeek {
					days = append(days, map[string]interface{}{
						"day":    dow.Day,
						"offset": dow.Offset,
					})
				}
			}

			ncList = append(ncList, map[string]interface{}{
				"config_id":             nc.Id,
				"frequency":             nc.Frequency,
				"enabled":               nc.Enabled,
				"recipients":            nc.Recipients,
				"detailed_report":       nc.DetailedReport,
				"with_compression":      nc.WithCompression,
				"include_remediation":   nc.IncludeRemediation,
				"last_updated":          nc.LastUpdated,
				"last_sent_ts":          nc.LastSentTs,
				"config_type":           nc.Type,
				"template_id":           nc.TemplateId,
				"timezone_id":           nc.TimezoneId,
				"day_of_month":          nc.DayOfMonth,
				"r_rule_schedule":       nc.RruleSchedule,
				"frequency_from_r_rule": nc.FrequencyFromRrule,
				"hour_of_day":           nc.HourOfDay,
				"days_of_week":          days,
			})
		}

		if err = d.Set("notification_config", ncList); err != nil {
			log.Printf("[WARN] Error setting 'notification_config' for %q: %s", d.Id(), err)
		}
	}
}

func createAlertRule(d *schema.ResourceData, meta interface{}) error {
	var err error
	client := meta.(*pc.Client)
	o := parseAlertRule(d, "")

	if err = rule.Create(client, o); err != nil {
		return err
	}

	PollApiUntilSuccess(func() error {
		_, err := rule.Identify(client, o.Name)
		return err
	})

	id, err := rule.Identify(client, o.Name)
	if err != nil {
		return err
	}

	PollApiUntilSuccess(func() error {
		_, err := rule.Get(client, id)
		return err
	})

	d.SetId(id)
	return readAlertRule(d, meta)
}

func readAlertRule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pc.Client)
	id := d.Id()

	o, err := rule.Get(client, id)
	if err != nil {
		if err == pc.ObjectNotFoundError {
			d.SetId("")
			return nil
		}
		return err
	}

	saveAlertRule(d, o)
	return nil
}

func updateAlertRule(d *schema.ResourceData, meta interface{}) error {
	var err error
	client := meta.(*pc.Client)
	id := d.Id()
	o := parseAlertRule(d, id)

	if err = rule.Update(client, o); err != nil {
		if err == pc.ObjectNotFoundError {
			d.SetId("")
			return nil
		}
		return err
	}

	return readAlertRule(d, meta)
}

func deleteAlertRule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pc.Client)
	id := d.Id()

	if err := rule.Delete(client, id); err != nil {
		if err != pc.ObjectNotFoundError {
			return err
		}
	}

	d.SetId("")
	return nil
}
