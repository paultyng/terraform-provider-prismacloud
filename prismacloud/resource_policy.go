package prismacloud

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"golang.org/x/net/context"

	pc "github.com/paloaltonetworks/prisma-cloud-go"
	"github.com/paloaltonetworks/prisma-cloud-go/cloud/account"
	"github.com/paloaltonetworks/prisma-cloud-go/policy"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/mapstructure"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: createPolicy,
		ReadContext:   readPolicy,
		UpdateContext: updatePolicy,
		DeleteContext: deletePolicy,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy type",
				ValidateFunc: validation.StringInSlice(
					[]string{
						policy.PolicyTypeConfig,
						policy.PolicyTypeAuditEvent,
						policy.PolicyTypeNetwork,
						policy.PolicyTypeIAM,
						policy.PolicyTypeAnomaly,
						policy.PolicyTypeData,
					},
					false,
				),
			},
			"policy_subtypes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Policy subtypes",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"system_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If policy is a system default policy or not",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description",
			},
			"severity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Severity",
				Default:     policy.SeverityLow,
				ValidateFunc: validation.StringInSlice(
					[]string{
						policy.SeverityLow,
						policy.SeverityMedium,
						policy.SeverityHigh,
						policy.SeverityCritical,
						policy.SeverityInformational,
					},
					false,
				),
			},
			"recommendation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remediation recommendation",
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cloud type (Required for config policies)",
				ValidateFunc: validation.StringInSlice(
					[]string{
						account.TypeAws,
						account.TypeAzure,
						account.TypeGcp,
						account.TypeAlibaba,
						"all",
					},
					false,
				),
				Default: "all",
			},
			"labels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Labels",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled",
				Default:     true,
			},
			"created_on": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Created on",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created by",
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
			"rule_last_modified_on": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule last modified on",
			},
			"overridden": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Overridden",
			},
			"deleted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Deleted",
				ForceNew:    true,
			},
			"restrict_alert_dismissal": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Restrict alert dismissal",
			},
			"open_alerts_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Open alerts count",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Owner",
			},
			"policy_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy mode",
			},
			"policy_category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy category",
			},
			"policy_class": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy class",
			},
			"remediable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is remediable or not",
			},
			"rule": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Model for rule",
				MaxItems:    1,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name",
						},
						"cloud_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cloud type",
						},
						"cloud_account": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cloud account",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource type",
						},
						"api_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "API name",
						},
						"resource_id_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource ID path",
						},
						"criteria": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Saved search ID that defines the rule criteria",
						},
						"data_criteria": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Criteria for DLP Rule",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"classification_result": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Data Profile name required for DLP rule criteria",
									},
									"exposure": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "File exposure",
										ValidateFunc: validation.StringInSlice(
											[]string{
												"private",
												"public",
												"conditional",
											},
											false,
										),
									},
									"extension": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "File extensions",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"parameters": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Parameters",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"rule_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of rule or RQL query",
							ValidateFunc: validation.StringInSlice(
								[]string{
									policy.RuleTypeConfig,
									policy.RuleTypeAuditEvent,
									policy.RuleTypeNetwork,
									policy.RuleTypeIAM,
									policy.RuleTypeAnomaly,
									policy.RuleTypeData,
									policy.RuleTypeNetworkConfig,
								},
								false,
							),
						},
						"children": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Children",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"criteria": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Criteria for build policy",
									},
									"metadata": {
										Type:        schema.TypeMap,
										Optional:    true,
										Computed:    true,
										Description: "YAML string for code build policy",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										ValidateDiagFunc: ValidateMetadata,
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Type of build policy",
										ValidateFunc: validation.StringInSlice(
											[]string{
												"tf",
												"cft",
												"k8s",
												"build",
											},
											false,
										),
									},
									"recommendation": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Recommendation",
									},
								},
							},
						},
					},
				},
			},
			"remediation": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Model for remediation",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Template type",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description",
						},
						"cli_script_template": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CLI script template",
						},
						"cli_script_json_schema_string": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CLI script JSON schema",
						},
						"actions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operation": {
								  		Type:     schema.TypeString,
								  		Optional: true,
									},
									"payload": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"compliance_metadata": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of compliance data. Each item has compliance standard, requirement, and/or section information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"standard_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance standard name",
						},
						"standard_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance standard description",
						},
						"requirement_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Requirement ID",
						},
						"requirement_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Requirement name",
						},
						"requirement_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Requirement description",
						},
						"section_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Section ID",
						},
						"section_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Section description",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy ID",
						},
						"compliance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Compliance Section UUID",
						},
						"section_label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Section label",
						},
						"custom_assigned": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Custom assigned",
						},
					},
				},
			},
		},
	}
}

func parsePolicy(d *schema.ResourceData, id string) policy.Policy {
	rspec := d.Get("rule").([]interface{})[0].(map[string]interface{})
	ps := d.Get("policy_subtypes")
	ans := policy.Policy{
		PolicyId:               id,
		Name:                   d.Get("name").(string),
		PolicyType:             d.Get("policy_type").(string),
		Description:            d.Get("description").(string),
		Severity:               d.Get("severity").(string),
		Recommendation:         d.Get("recommendation").(string),
		CloudType:              d.Get("cloud_type").(string),
		Labels:                 SetToStringSlice(d.Get("labels").(*schema.Set)),
		Enabled:                d.Get("enabled").(bool),
		Overridden:             d.Get("overridden").(bool),
		Deleted:                d.Get("deleted").(bool),
		RestrictAlertDismissal: d.Get("restrict_alert_dismissal").(bool),
		PolicySubTypes:         SetToStringSlice(ps.(*schema.Set)),
		Rule: policy.Rule{
			Name:           rspec["name"].(string),
			CloudType:      rspec["cloud_type"].(string),
			CloudAccount:   rspec["cloud_account"].(string),
			ResourceType:   rspec["resource_type"].(string),
			ApiName:        rspec["api_name"].(string),
			ResourceIdPath: rspec["resource_id_path"].(string),
			Type:           rspec["rule_type"].(string),
		},
	}
	if rspec["criteria"].(string) != "" {
		ans.Rule.Criteria = rspec["criteria"].(string)
	}

	dataCriteria := rspec["data_criteria"].([]interface{})
	if len(dataCriteria) > 0 {
		if data := dataCriteria[0].(map[string]interface{}); len(data) > 0 {
			ans.Rule.DataCriteria.ClassificationResult = data["classification_result"].(string)
			ans.Rule.DataCriteria.Exposure = data["exposure"].(string)
			ans.Rule.DataCriteria.Extension = SetToStringSlice(data["extension"].(*schema.Set))
		}
	}

	cld := rspec["children"].([]interface{})
	ans.Rule.Children = make([]policy.Children, 0, len(cld))
	for _, chi := range cld {
		cl := chi.(map[string]interface{})
		var md policy.Metadata
		err := mapstructure.Decode(cl["metadata"], &md)
		if err != nil {
			panic(err)
		}
		ans.Rule.Children = append(ans.Rule.Children, policy.Children{
			Criteria:       cl["criteria"].(string),
			Type:           cl["type"].(string),
			Recommendation: cl["recommendation"].(string),
			Metadata:       md,
		})
	}

	rem := d.Get("remediation").([]interface{})
	if len(rem) > 0 {
		if rems := rem[0].(map[string]interface{}); len(rems) > 0 {
			ans.Remediation.TemplateType = rems["template_type"].(string)
			ans.Remediation.Description = rems["description"].(string)
			ans.Remediation.CliScriptTemplate = rems["cli_script_template"].(string)
			if rems["actions"] != nil && len(rems["actions"].([]interface{})) > 0 {
				actions := make([]policy.Action, 0, len(rems["actions"].([]interface{})))
				for _, action := range rems["actions"].([]interface{}) {
					actions = append(actions, policy.Action{
						Operation:    action.(map[string]interface{})["operation"].(string),
						Payload: action.(map[string]interface{})["payload"].(string),
					})
				}
				ans.Remediation.Actions = actions
			}
			var csjs interface{}
			if err := json.Unmarshal([]byte(rems["cli_script_json_schema_string"].(string)), &csjs); err != nil {
				log.Printf("[WARN] Error unmarshalling 'cli_script_json_schema_string' for %q: %s", d.Id(), err)
			}
			ans.Remediation.CliScriptJsonSchema = csjs
		}
	}

	rsp := rspec["parameters"].(map[string]interface{})
	if len(rsp) > 0 {
		ans.Rule.Parameters = make(map[string]string)
		for key, val := range rsp {
			ans.Rule.Parameters[key] = val.(string)
		}
	}

	cms := d.Get("compliance_metadata").(*schema.Set).List()
	ans.ComplianceMetadata = make([]policy.ComplianceMetadata, 0, len(cms))
	for _, csmi := range cms {
		cmd := csmi.(map[string]interface{})
		ans.ComplianceMetadata = append(ans.ComplianceMetadata, policy.ComplianceMetadata{
			ComplianceId: cmd["compliance_id"].(string),
		})
	}

	return ans
}

func savePolicy(d *schema.ResourceData, obj policy.Policy) {
	d.Set("policy_id", obj.PolicyId)
	d.Set("name", obj.Name)
	d.Set("policy_type", obj.PolicyType)
	d.Set("system_default", obj.SystemDefault)
	d.Set("description", obj.Description)
	d.Set("severity", obj.Severity)
	d.Set("recommendation", obj.Recommendation)
	d.Set("cloud_type", obj.CloudType)
	d.Set("enabled", obj.Enabled)
	d.Set("created_on", obj.CreatedOn)
	d.Set("created_by", obj.CreatedBy)
	d.Set("last_modified_on", obj.LastModifiedOn)
	d.Set("last_modified_by", obj.LastModifiedBy)
	d.Set("rule_last_modified_on", obj.RuleLastModifiedOn)
	d.Set("overridden", obj.Overridden)
	d.Set("deleted", obj.Deleted)
	d.Set("restrict_alert_dismissal", obj.RestrictAlertDismissal)
	d.Set("open_alerts_count", obj.OpenAlertsCount)
	d.Set("owner", obj.Owner)
	d.Set("policy_mode", obj.PolicyMode)
	d.Set("policy_category", obj.PolicyCategory)
	d.Set("policy_class", obj.PolicyClass)
	d.Set("remediable", obj.Remediable)

	if err := d.Set("policy_subtypes", StringSliceToSet(obj.PolicySubTypes)); err != nil {
		log.Printf("[WARN] Error setting 'policy_subtypes' for %q: %s", d.Id(), err)
	}

	if err := d.Set("labels", StringSliceToSet(obj.Labels)); err != nil {
		log.Printf("[WARN] Error setting 'labels' for %q: %s", d.Id(), err)
	}

	// Rule.
	rv := map[string]interface{}{
		"name":             obj.Rule.Name,
		"cloud_type":       obj.Rule.CloudType,
		"cloud_account":    obj.Rule.CloudAccount,
		"resource_type":    obj.Rule.ResourceType,
		"api_name":         obj.Rule.ApiName,
		"resource_id_path": obj.Rule.ResourceIdPath,
		"rule_type":        obj.Rule.Type,
	}

	switch v := obj.Rule.Criteria.(type) {
	case string:
		x := ResourceDataInterfaceMap(d, "rule")
		if x["criteria"] == nil {
			rv["criteria"] = v
		} else {
			rv["criteria"] = x["criteria"]
		}
	case interface{}:
		b, err := json.Marshal(v)
		if err != nil {
			log.Printf("[WARN] Failed to marshal criteria for %q: %s", d.Id(), err)
		}
		rv["criteria"] = string(b)
	}

	if obj.Rule.DataCriteria.ClassificationResult == "" {
		rv["data_criteria"] = nil
	} else {
		data := map[string]interface{}{
			"classification_result": obj.Rule.DataCriteria.ClassificationResult,
			"exposure":              obj.Rule.DataCriteria.Exposure,
			"extension":             obj.Rule.DataCriteria.Extension,
		}
		rv["data_criteria"] = []interface{}{data}
	}

	cld := make([]interface{}, 0, len(obj.Rule.Children))
	for _, chi := range obj.Rule.Children {
		cld = append(cld, map[string]interface{}{
			"criteria":       chi.Criteria,
			"type":           chi.Type,
			"recommendation": chi.Recommendation,
			"metadata":       map[string]string{"code": chi.Metadata.Code},
		})
	}
	rv["children"] = cld

	pm := make(map[string]interface{})
	for k, v := range obj.Rule.Parameters {
		pm[k] = v
	}
	x := ResourceDataInterfaceMap(d, "rule")
	y := x["parameters"]
	if rec, ok := y.(map[string]interface{}); ok {
		for key, val := range rec {
			pm[key] = val
		}
	}
	rv["parameters"] = pm

	if err := d.Set("rule", []interface{}{rv}); err != nil {
		log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
	}

	// Remediation.
	if obj.Remediation.TemplateType == "" && obj.Remediation.Description == "" && obj.Remediation.CliScriptTemplate == "" && obj.Remediation.CliScriptJsonSchema == nil {
		d.Set("remediation", nil)
	} else {
		var csjs string
		if obj.Remediation.CliScriptJsonSchema != nil {
			b, err := json.Marshal(obj.Remediation.CliScriptJsonSchema)
			if err != nil {
				log.Printf("[WARN] Failed to marshal cli script json schema for %s: %s", d.Id(), err)
			}
			csjs = string(b)
		}
		var actions []map[string]string
		if obj.Remediation.Actions != nil  && len(obj.Remediation.Actions) > 0 {
			actions = make([]map[string]string, len(obj.Remediation.Actions))
			for i, action := range obj.Remediation.Actions {
				actions[i] = map[string]string{
					"operation": action.Operation,
					"payload":   action.Payload,
				}
			}
		}
		rem := map[string]interface{}{
			"template_type":                 obj.Remediation.TemplateType,
			"description":                   obj.Remediation.Description,
			"cli_script_template":           obj.Remediation.CliScriptTemplate,
			"actions":                       actions,
			"cli_script_json_schema_string": csjs,
		}
		if err := d.Set("remediation", []interface{}{rem} ); err != nil {
			log.Printf("[WARN] Error setting 'remediation' for %q: %s", d.Id(), err)
		}
	}

	if len(obj.ComplianceMetadata) == 0 {
		d.Set("compliance_metadata", nil)
		return
	}

	cmList := make([]interface{}, 0, len(obj.ComplianceMetadata))
	for _, cm := range obj.ComplianceMetadata {
		cmList = append(cmList, map[string]interface{}{
			"standard_name":           cm.StandardName,
			"standard_description":    cm.StandardDescription,
			"requirement_id":          cm.RequirementId,
			"requirement_name":        cm.RequirementName,
			"requirement_description": cm.RequirementDescription,
			"section_id":              cm.SectionId,
			"section_description":     cm.SectionDescription,
			"policy_id":               cm.PolicyId,
			"compliance_id":           cm.ComplianceId,
			"section_label":           cm.SectionLabel,
			"custom_assigned":         cm.CustomAssigned,
		})
	}
	if err := d.Set("compliance_metadata", cmList); err != nil {
		log.Printf("[WARN] Error setting 'compliance_metadata' for %q: %s", d.Id(), err)
	}
}

func validateBuildPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)
	obj := parsePolicy(d, "")

	if err := policy.ValidateBuildPolicy(client, obj); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func createPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)
	obj := parsePolicy(d, "")

	if diag := validateBuildPolicy(ctx, d, client); diag != nil {
		return diag
	}

	if err := policy.Create(client, obj); err != nil {
		return diag.FromErr(err)
	}

	PollApiUntilSuccess(func() error {
		_, err := policy.Identify(client, obj.Name)
		return err
	})

	id, err := policy.Identify(client, obj.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	PollApiUntilSuccess(func() error {
		_, err := policy.Get(client, id)
		return err
	})

	d.SetId(id)
	return readPolicy(ctx, d, meta)
}

func readPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)
	id := d.Id()

	obj, err := policy.Get(client, id)
	if err != nil {
		if err == pc.ObjectNotFoundError {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	savePolicy(d, obj)

	return nil
}

func updatePolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)
	id := d.Id()
	obj := parsePolicy(d, id)

	if err := policy.Update(client, obj); err != nil {
		return diag.FromErr(err)
	}

	return readPolicy(ctx, d, meta)
}

func deletePolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)
	id := d.Id()
	obj := parsePolicy(d, "")

	err := policy.Delete(client, id, obj)
	if err != nil {
		if err != pc.ObjectNotFoundError {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}
