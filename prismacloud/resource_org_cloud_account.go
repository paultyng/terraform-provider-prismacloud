package prismacloud

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"golang.org/x/net/context"
	"log"
	"strings"
	"time"

	pc "github.com/paloaltonetworks/prisma-cloud-go"
	"github.com/paloaltonetworks/prisma-cloud-go/cloud/account/org"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOrgCloudAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: createOrgCloudAccount,
		ReadContext:   readOrgCloudAccount,
		UpdateContext: updateOrgCloudAccount,
		DeleteContext: deleteOrgCloudAccount,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"disable_on_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "to disable cloud account instead of deleting on calling destroy",
				Default:     false,
			},

			org.TypeAwsOrg: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AWS account type",
				MaxItems:    1,
				ConflictsWith: []string{
					org.TypeGcpOrg,
					org.TypeAzureOrg,
					org.TypeOci,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AWS account ID",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the account is enabled",
							Default:     true,
						},
						"external_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "AWS account external ID",
							Sensitive:   true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return true
							},
						},
						"group_ids": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "List of account IDs to which you are assigning this account",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to be used for the account on the Prisma Cloud platform (must be unique)",
						},
						"role_arn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique identifier for an AWS resource (ARN)",
						},
						"account_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "organization",
							Description: "Account type - organization or account",
							ValidateFunc: validation.StringInSlice(
								[]string{
									"organization",
									"account",
								},
								true,
							),
						},
						"member_role_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AWS Member account role name",
						},
						"member_external_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "AWS Member account role's external ID",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return true
							},
						},
						"member_role_status": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "true = The member role created using stack set exists in all the member accounts. All the Org accounts will be added.\nfalse = Only the master account will be added.",
						},
						"protection_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "MONITOR",
							Description: "Monitor or Monitor and Protect",
							ValidateFunc: validation.StringInSlice(
								[]string{
									"MONITOR",
									"MONITOR_AND_PROTECT",
								},
								false,
							),
						},
						"hierarchy_selection": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "List of hierarchy selection. Each item has resource id, display name, node type and selection type",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource ID. Valid values are AWS OU ID, AWS account ID, or AWS Organization ID. Note you must escape any double quotes in the resource ID with a backslash.",
									},
									"display_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Display name for AWS OU, AWS account, or AWS organization",
									},
									"node_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Valid values: OU, ACCOUNT, ORG",
										ValidateFunc: validation.StringInSlice(
											[]string{
												"OU",
												"ACCOUNT",
												"ORG",
											},
											false,
										),
									},
									"selection_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Selection type. Valid values: INCLUDE to include the specified resource to onboard, EXCLUDE to exclude the specified resource and onboard the rest, ALL to onboard all resources in the organization.",
										ValidateFunc: validation.StringInSlice(
											[]string{
												"INCLUDE",
												"EXCLUDE",
												"ALL",
											},
											false,
										),
									},
								},
							},
						},
					},
				},
			},

			org.TypeAzureOrg: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Azure account type",
				MaxItems:    1,
				ConflictsWith: []string{
					org.TypeGcpOrg,
					org.TypeAwsOrg,
					org.TypeOci,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Application ID registered with Active Directory",
						},
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Azure account ID",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the account is enabled",
							Default:     true,
						},
						"group_ids": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "List of account IDs to which you are assigning this account",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to be used for the account on the Prisma Cloud platform (must be unique)",
						},
						"account_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "tenant",
							Description: "Account type - account or tenant",
							ValidateFunc: validation.StringInSlice(
								[]string{
									"tenant",
									"account",
								},
								true,
							),
						},
						"protection_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "MONITOR",
							Description: "Monitor or Monitor and Protect",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Active Directory ID associated with Azure",
						},
						"service_principal_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique ID of the service principal object associated with the Prisma Cloud application that you create",
						},
						"monitor_flow_logs": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Automatically ingest flow logs",
						},
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Application ID key",
							Sensitive:   true,
						},
						"root_sync_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Azure tenant has children. Must be set to true when azure tenant is onboarded with children",
						},
						"hierarchy_selection": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "List of subscriptions and/or management groups to onboard",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource ID. Management group ID or subscription ID.\nNote you must escape any double quotes in the resource ID with a backslash.",
									},
									"display_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Display name for management group or subscription",
									},
									"node_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Node type. Valid values: SUBSCRIPTION, TENANT, MANAGEMENT_GROUP",
										ValidateFunc: validation.StringInSlice(
											[]string{
												"SUBSCRIPTION",
												"TENANT",
												"MANAGEMENT_GROUP",
											},
											false,
										),
									},
									"selection_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Selection type. Valid values: INCLUDE to include the specified resource to onboard, EXCLUDE to exclude the specified resource and onboard the rest, ALL to onboard all resources in the tenant.",
										ValidateFunc: validation.StringInSlice(
											[]string{
												"INCLUDE",
												"EXCLUDE",
												"ALL",
											},
											false,
										),
									},
								},
							},
						},
					},
				},
			},

			org.TypeGcpOrg: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Gcp account type",
				MaxItems:    1,
				ConflictsWith: []string{
					org.TypeAwsOrg,
					org.TypeAzureOrg,
					org.TypeOci,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "GCP project ID",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the account is enabled",
							Default:     true,
						},
						"group_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "List of account IDs to which you are assigning this account",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to be used for the account on the Prisma Cloud platform (must be unique)",
						},
						"compression_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable flow log compression",
						},
						"dataflow_enabled_project": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "GCP project for flow log compression",
						},
						"flow_log_storage_bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "GCP flow logs storage bucket",
						},
						// Use a json string until this feature is added:
						// https://github.com/hashicorp/terraform-plugin-sdk/issues/248
						"credentials_json": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Content of the JSON credentials file",
							Sensitive:        true,
							DiffSuppressFunc: gcpOrgCredentialsMatch,
						},
						"account_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "organization",
							Description: "Account type - organization or account",
							ValidateFunc: validation.StringInSlice(
								[]string{
									"organization",
									"account",
								},
								true,
							),
						},
						"protection_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "MONITOR",
							Description: "Protection Mode - Monitor or Monitor and Protect",
							ValidateFunc: validation.StringInSlice(
								[]string{
									"MONITOR",
									"MONITOR_AND_PROTECT",
								},
								true,
							),
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if strings.ToLower(old) == strings.ToLower(new) {
									return true
								}
								return false
							},
						},
						"organization_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "GCP organization name",
						},
						"account_group_creation_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "MANUAL",
							Description: "Cloud account group creation mode. Valid values - MANUAL, AUTO or RECURSIVE",
							ValidateFunc: validation.StringInSlice(
								[]string{
									"MANUAL",
									"AUTO",
									"RECURSIVE",
								},
								false,
							),
						},
						"hierarchy_selection": {
							Type:        schema.TypeSet,
							Computed:    true,
							Optional:    true,
							Description: "List of hierarchy selection. Each item has resource id, display name, node type and selection type",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource ID. For folders, format is folders/{folder ID}. For projects, format is {project number}. For orgs, format is organizations/{org ID}",
									},
									"display_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Display name for folder, project, or organization",
									},
									"node_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Node type. Valid values - FOLDER, PROJECT, ORG",
										ValidateFunc: validation.StringInSlice(
											[]string{
												"FOLDER",
												"PROJECT",
												"ORG",
											},
											false,
										),
									},
									"selection_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Selection type. Valid values: INCLUDE, EXCLUDE, ALL. If hierarchySelection.nodeType is PROJECT or FOLDER, then a valid value is either INCLUDE or EXCLUDE",
										ValidateFunc: validation.StringInSlice(
											[]string{
												"INCLUDE",
												"EXCLUDE",
												"ALL",
											},
											false,
										),
									},
								},
							},
						},
					},
				},
			},

			org.TypeOci: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Oci account type",
				MaxItems:    1,
				ConflictsWith: []string{
					org.TypeAwsOrg,
					org.TypeAzureOrg,
					org.TypeGcpOrg,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Org account ID",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether or not the account is enabled",
						},
						"group_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "OCI identity group name that you define. Can be an existing group",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to be used for the account on the Prisma Cloud platform (must be unique)",
						},
						"account_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Account type - account or tenant",
							ValidateFunc: validation.StringInSlice(
								[]string{
									"tenant",
									"account",
								},
								true,
							),
						},
						"default_account_group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Account group id for this account",
						},
						"home_region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "OCI tenancy home region",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "OCI identity policy name that you define. Can be an existing policy that has the right policy statements",
						},
						"user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "OCI identity user name that you define. Can be an existing user that has the right privileges",
						},
						"user_ocid": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "OCI identity user Ocid that you define. Can be an existing user that has the right privileges",
						},
					},
				},
			},
		},
	}
}

func gcpOrgCredentialsMatch(k, old, new string, d *schema.ResourceData) bool {
	var (
		err       error
		prev, cur org.GcpOrgCredentials
	)

	if err = json.Unmarshal([]byte(old), &prev); err != nil {
		return false
	}

	if err = json.Unmarshal([]byte(new), &cur); err != nil {
		return false
	}
	if cur.Type == "service_account" {
		return (prev.Type == cur.Type &&
			prev.ProjectId == cur.ProjectId &&
			prev.PrivateKeyId == cur.PrivateKeyId &&
			prev.ClientEmail == cur.ClientEmail &&
			prev.ClientId == cur.ClientId &&
			prev.AuthUri == cur.AuthUri &&
			prev.TokenUri == cur.TokenUri &&
			prev.ProviderCertUrl == cur.ProviderCertUrl &&
			prev.ClientCertUrl == cur.ClientCertUrl)
	} else {
		return (prev.Type == cur.Type &&
			prev.Audience == cur.Audience &&
			prev.ServiceAccountImpersonationUrl == cur.ServiceAccountImpersonationUrl &&
			prev.SubjectTokenType == cur.SubjectTokenType &&
			prev.TokenUrl == cur.TokenUrl &&
			prev.CredentialSource == cur.CredentialSource)
	}

}

func parseOrgCloudAccount(d *schema.ResourceData) (string, string, interface{}) {
	if x := ResourceDataInterfaceMap(d, org.TypeAwsOrg); len(x) != 0 {
		ans := org.AwsOrg{
			AccountId:        x["account_id"].(string),
			Enabled:          x["enabled"].(bool),
			ExternalId:       x["external_id"].(string),
			GroupIds:         SetToStringSlice(x["group_ids"].(*schema.Set)),
			Name:             x["name"].(string),
			RoleArn:          x["role_arn"].(string),
			ProtectionMode:   x["protection_mode"].(string),
			AccountType:      x["account_type"].(string),
			MemberRoleName:   x["member_role_name"].(string),
			MemberExternalId: x["member_external_id"].(string),
			MemberRoleStatus: x["member_role_status"].(bool),
		}
		hsl := x["hierarchy_selection"].(*schema.Set).List()
		ans.HierarchySelection = make([]org.HierarchySelection, 0, len(hsl))
		for _, hsi := range hsl {
			hs := hsi.(map[string]interface{})
			ans.HierarchySelection = append(ans.HierarchySelection, org.HierarchySelection{
				ResourceId:    hs["resource_id"].(string),
				DisplayName:   hs["display_name"].(string),
				SelectionType: hs["selection_type"].(string),
				NodeType:      hs["node_type"].(string),
			})
		}
		return org.TypeAwsOrg, x["name"].(string), ans
	} else if x := ResourceDataInterfaceMap(d, org.TypeOci); len(x) != 0 {
		return org.TypeOci, x["name"].(string), org.Oci{
			AccountId:             x["account_id"].(string),
			Enabled:               x["enabled"].(bool),
			Name:                  x["name"].(string),
			AccountType:           x["account_type"].(string),
			DefaultAccountGroupId: x["default_account_group_id"].(string),
			GroupName:             x["group_name"].(string),
			HomeRegion:            x["home_region"].(string),
			PolicyName:            x["policy_name"].(string),
			UserName:              x["user_name"].(string),
			UserOcid:              x["user_ocid"].(string),
		}
	} else if x := ResourceDataInterfaceMap(d, org.TypeAzureOrg); len(x) != 0 {
		ans := org.AzureOrg{
			Account: org.AzureCloudAccount{
				AccountId:      x["account_id"].(string),
				Enabled:        x["enabled"].(bool),
				Name:           x["name"].(string),
				ProtectionMode: x["protection_mode"].(string),
				AccountType:    x["account_type"].(string),
				GroupIds:       SetToStringSlice(x["group_ids"].(*schema.Set)),
			},
			ClientId:           x["client_id"].(string),
			TenantId:           x["tenant_id"].(string),
			ServicePrincipalId: x["service_principal_id"].(string),
			MonitorFlowLogs:    x["monitor_flow_logs"].(bool),
			Key:                x["key"].(string),
			RootSyncEnabled:    x["root_sync_enabled"].(bool),
		}
		hsl := x["hierarchy_selection"].(*schema.Set).List()
		ans.HierarchySelection = make([]org.HierarchySelection, 0, len(hsl))
		for _, hsi := range hsl {
			hs := hsi.(map[string]interface{})
			ans.HierarchySelection = append(ans.HierarchySelection, org.HierarchySelection{
				ResourceId:    hs["resource_id"].(string),
				DisplayName:   hs["display_name"].(string),
				SelectionType: hs["selection_type"].(string),
				NodeType:      hs["node_type"].(string),
			})
		}
		return org.TypeAzureOrg, x["name"].(string), ans
	} else if x := ResourceDataInterfaceMap(d, org.TypeGcpOrg); len(x) != 0 {
		var creds org.GcpOrgCredentials
		_ = json.Unmarshal([]byte(x["credentials_json"].(string)), &creds)

		ans := org.GcpOrg{
			Account: org.GcpCloudAccount{
				AccountId:      x["account_id"].(string),
				Enabled:        x["enabled"].(bool),
				ProjectId:      creds.ProjectId,
				Name:           x["name"].(string),
				ProtectionMode: x["protection_mode"].(string),
				AccountType:    x["account_type"].(string),
				GroupIds:       SetToStringSlice(x["group_ids"].(*schema.Set)),
			},
			CompressionEnabled:       x["compression_enabled"].(bool),
			DataflowEnabledProject:   x["dataflow_enabled_project"].(string),
			FlowLogStorageBucket:     x["flow_log_storage_bucket"].(string),
			Credentials:              creds,
			OrganizationName:         x["organization_name"].(string),
			AccountGroupCreationMode: x["account_group_creation_mode"].(string),
		}
		hsl := x["hierarchy_selection"].(*schema.Set).List()
		ans.HierarchySelection = make([]org.HierarchySelection, 0, len(hsl))
		for _, hsi := range hsl {
			hs := hsi.(map[string]interface{})
			ans.HierarchySelection = append(ans.HierarchySelection, org.HierarchySelection{
				ResourceId:    hs["resource_id"].(string),
				DisplayName:   hs["display_name"].(string),
				SelectionType: hs["selection_type"].(string),
				NodeType:      hs["node_type"].(string),
			})
		}
		return org.TypeGcpOrg, x["name"].(string), ans
	}
	return "", "", nil
}

func saveOrgCloudAccount(d *schema.ResourceData, dest string, obj interface{}) {
	var val map[string]interface{}

	switch v := obj.(type) {
	case org.AwsOrg:
		val = map[string]interface{}{
			"account_id":         v.AccountId,
			"enabled":            v.Enabled,
			"external_id":        v.ExternalId,
			"group_ids":          v.GroupIds,
			"name":               v.Name,
			"role_arn":           v.RoleArn,
			"protection_mode":    v.ProtectionMode,
			"account_type":       v.AccountType,
			"member_role_name":   v.MemberRoleName,
			"member_external_id": v.MemberExternalId,
			"member_role_status": v.MemberRoleStatus,
		}
		if len(v.HierarchySelection) == 0 {
			val["hierarchy_selection"] = nil
		} else {
			hsList := make([]interface{}, 0, len(v.HierarchySelection))
			for _, hs := range v.HierarchySelection {
				hsList = append(hsList, map[string]interface{}{
					"resource_id":    hs.ResourceId,
					"display_name":   hs.DisplayName,
					"node_type":      hs.NodeType,
					"selection_type": hs.SelectionType,
				})
			}
			val["hierarchy_selection"] = hsList
		}
		x := ResourceDataInterfaceMap(d, org.TypeAwsOrg)
		if x["hierarchy_selection"] == nil {
			val["hierarchy_selection"] = nil
		} else {
			hsl := x["hierarchy_selection"].(*schema.Set).List()
			hsList := make([]interface{}, 0, len(hsl))
			for _, hsi := range hsl {
				hs := hsi.(map[string]interface{})
				hsList = append(hsList, map[string]interface{}{
					"resource_id":    hs["resource_id"].(string),
					"display_name":   hs["display_name"].(string),
					"selection_type": hs["selection_type"].(string),
					"node_type":      hs["node_type"].(string),
				})
			}
			val["hierarchy_selection"] = hsList
		}
	case org.Oci:
		val = map[string]interface{}{
			"account_id":               v.AccountId,
			"enabled":                  v.Enabled,
			"name":                     v.Name,
			"account_type":             v.AccountType,
			"default_account_group_id": v.DefaultAccountGroupId,
			"group_name":               v.GroupName,
			"home_region":              v.HomeRegion,
			"policy_name":              v.PolicyName,
			"user_name":                v.UserName,
			"user_ocid":                v.UserOcid,
		}
	case org.AzureOrg:
		x := ResourceDataInterfaceMap(d, org.TypeAzureOrg)
		var key string
		if x["key"] == nil {
			key = v.Key
		} else {
			key = x["key"].(string)
		}
		val = map[string]interface{}{
			"account_id":           v.Account.AccountId,
			"enabled":              v.Account.Enabled,
			"group_ids":            v.Account.GroupIds,
			"name":                 v.Account.Name,
			"protection_mode":      v.Account.ProtectionMode,
			"account_type":         v.Account.AccountType,
			"client_id":            v.ClientId,
			"tenant_id":            v.TenantId,
			"service_principal_id": v.ServicePrincipalId,
			"monitor_flow_logs":    v.MonitorFlowLogs,
			"key":                  key,
			"root_sync_enabled":    v.RootSyncEnabled,
		}
		if len(v.HierarchySelection) == 0 {
			val["hierarchy_selection"] = nil
		} else {
			hsList := make([]interface{}, 0, len(v.HierarchySelection))
			for _, hs := range v.HierarchySelection {
				hsList = append(hsList, map[string]interface{}{
					"resource_id":    hs.ResourceId,
					"display_name":   hs.DisplayName,
					"node_type":      hs.NodeType,
					"selection_type": hs.SelectionType,
				})
			}
			val["hierarchy_selection"] = hsList
		}
		if x["hierarchy_selection"] == nil {
			val["hierarchy_selection"] = nil
		} else {
			hsl := x["hierarchy_selection"].(*schema.Set).List()
			hsList := make([]interface{}, 0, len(hsl))
			for _, hsi := range hsl {
				hs := hsi.(map[string]interface{})
				hsList = append(hsList, map[string]interface{}{
					"resource_id":    hs["resource_id"].(string),
					"display_name":   hs["display_name"].(string),
					"selection_type": hs["selection_type"].(string),
					"node_type":      hs["node_type"].(string),
				})
			}
			val["hierarchy_selection"] = hsList
		}
	case org.GcpOrg:
		b, _ := json.Marshal(v.Credentials)
		val = map[string]interface{}{
			"account_id":                  v.Account.AccountId,
			"name":                        v.Account.Name,
			"enabled":                     v.Account.Enabled,
			"group_ids":                   v.Account.GroupIds,
			"compression_enabled":         v.CompressionEnabled,
			"dataflow_enabled_project":    v.DataflowEnabledProject,
			"flow_log_storage_bucket":     v.FlowLogStorageBucket,
			"credentials_json":            string(b),
			"protection_mode":             v.Account.ProtectionMode,
			"account_type":                v.Account.AccountType,
			"organization_name":           v.OrganizationName,
			"account_group_creation_mode": v.AccountGroupCreationMode,
		}
		if len(v.HierarchySelection) == 0 {
			val["hierarchy_selection"] = nil
		} else {
			hsList := make([]interface{}, 0, len(v.HierarchySelection))
			for _, hs := range v.HierarchySelection {
				hsList = append(hsList, map[string]interface{}{
					"resource_id":    hs.ResourceId,
					"display_name":   hs.DisplayName,
					"node_type":      hs.NodeType,
					"selection_type": hs.SelectionType,
				})
			}
			val["hierarchy_selection"] = hsList
		}
		x := ResourceDataInterfaceMap(d, org.TypeGcpOrg)
		if x["hierarchy_selection"] == nil {
			val["hierarchy_selection"] = nil
		} else {
			hsl := x["hierarchy_selection"].(*schema.Set).List()
			hsList := make([]interface{}, 0, len(hsl))
			for _, hsi := range hsl {
				hs := hsi.(map[string]interface{})
				hsList = append(hsList, map[string]interface{}{
					"resource_id":    hs["resource_id"].(string),
					"display_name":   hs["display_name"].(string),
					"selection_type": hs["selection_type"].(string),
					"node_type":      hs["node_type"].(string),
				})
			}
			val["hierarchy_selection"] = hsList
		}
	}

	for _, key := range []string{org.TypeAwsOrg, org.TypeGcpOrg, org.TypeAzureOrg, org.TypeOci} {
		if key != dest {
			d.Set(key, nil)
			continue
		}

		if err := d.Set(key, []interface{}{val}); err != nil {
			log.Printf("[WARN] Error setting %q field for %q: %s", key, d.Id(), err)
		}
	}
}

func createOrgCloudAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)
	cloudType, name, obj := parseOrgCloudAccount(d)
	if err := org.Create(client, obj); err != nil {
		if strings.Contains(err.Error(), "duplicate_cloud_account") {
			if err := org.Update(client, obj); err != nil {
				return diag.FromErr(err)
			}
		} else {
			return diag.FromErr(err)
		}
	}
	PollApiUntilSuccess(func() error {
		_, err := org.Identify(client, cloudType, name)
		return err
	})

	id, err := org.Identify(client, cloudType, name)
	if err != nil {
		return diag.FromErr(err)
	}

	PollApiUntilSuccess(func() error {
		_, err := org.Get(client, cloudType, id)
		return err
	})

	d.SetId(TwoStringsToId(cloudType, id))
	return readOrgCloudAccount(ctx, d, meta)
}

func readOrgCloudAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)
	cloudType, id := IdToTwoStrings(d.Id())

	obj, err := org.Get(client, cloudType, id)
	if err != nil {
		if err == pc.ObjectNotFoundError {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	saveOrgCloudAccount(d, cloudType, obj)

	return nil
}

func updateOrgCloudAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)
	_, _, obj := parseOrgCloudAccount(d)

	if err := org.Update(client, obj); err != nil {
		return diag.FromErr(err)
	}

	return readOrgCloudAccount(ctx, d, meta)
}

func deleteOrgCloudAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)
	cloudType, id := IdToTwoStrings(d.Id())
	disable := d.Get("disable_on_destroy").(bool)

	if disable {
		switch cloudType {
		case org.TypeAwsOrg:
			cloudAccount, _ := org.Get(client, cloudType, id)
			cloudAccountAws := cloudAccount.(org.AwsOrg)
			cloudAccountAws.Enabled = false
			if err := org.Update(client, cloudAccountAws); err != nil {
				return diag.FromErr(err)
			}
			return nil

		case org.TypeAzureOrg:
			cloudAccount, _ := org.Get(client, cloudType, id)
			cloudAccountAzure := cloudAccount.(org.AzureOrg)
			cloudAccountAzure.Account.Enabled = false
			if err := org.Update(client, cloudAccountAzure); err != nil {
				return diag.FromErr(err)
			}
			return nil

		case org.TypeGcpOrg:
			cloudAccount, _ := org.Get(client, cloudType, id)
			cloudAccountGcp := cloudAccount.(org.GcpOrg)
			cloudAccountGcp.Account.Enabled = false
			if err := org.Update(client, cloudAccountGcp); err != nil {
				return diag.FromErr(err)
			}
			return nil
		case org.TypeOci:
			cloudAccount, _ := org.Get(client, cloudType, id)
			cloudAccountGcp := cloudAccount.(org.Oci)
			cloudAccountGcp.Enabled = false
			if err := org.Update(client, cloudAccountGcp); err != nil {
				return diag.FromErr(err)
			}
			return nil
		}
	}

	err := org.Delete(client, cloudType, id)
	if err != nil {
		if err != pc.ObjectNotFoundError {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}
