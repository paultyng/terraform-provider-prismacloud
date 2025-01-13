package prismacloud

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"golang.org/x/net/context"
	"log"
	"strings"
	"time"

	pc "github.com/paloaltonetworks/prisma-cloud-go"
	"github.com/paloaltonetworks/prisma-cloud-go/cloud/account"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: createCloudAccount,
		ReadContext:   readCloudAccount,
		UpdateContext: updateCloudAccount,
		DeleteContext: deleteCloudAccount,

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

			// AWS type.
			account.TypeAws: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AWS account type",
				MaxItems:    1,
				ConflictsWith: []string{
					account.TypeAzure,
					account.TypeGcp,
					account.TypeAlibaba,
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
							Default:     "account",
							Description: "Account type - organization or account",
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
					},
				},
			},

			// Azure type.
			account.TypeAzure: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Azure account type",
				MaxItems:    1,
				ConflictsWith: []string{
					account.TypeAws,
					account.TypeGcp,
					account.TypeAlibaba,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Application ID registered with Active Directory",
						},
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Application ID key",
							Sensitive:   true,
						},
						"monitor_flow_logs": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Automatically ingest flow logs",
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
						"account_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "account",
							Description: "Account type - organization or account",
						},
						"protection_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "MONITOR",
							Description: "Monitor or Monitor and Protect",
							ForceNew:    true,
						},
					},
				},
			},

			// GCP type.
			account.TypeGcp: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "GCP account type",
				MaxItems:    1,
				ConflictsWith: []string{
					account.TypeAws,
					account.TypeAzure,
					account.TypeAlibaba,
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
							DiffSuppressFunc: gcpCredentialsMatch,
						},
						"account_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "account",
							Description: "Account type - organization or account",
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
					},
				},
			},

			// Alibaba type.
			account.TypeAlibaba: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Alibaba account type",
				MaxItems:    1,
				ConflictsWith: []string{
					account.TypeAws,
					account.TypeAzure,
					account.TypeGcp,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Alibaba account ID",
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
						"ram_arn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique identifier for an Alibaba RAM role resource",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the account is enabled",
							Default:     true,
						},
					},
				},
			},
		},
	}
}

func gcpCredentialsMatch(k, old, new string, d *schema.ResourceData) bool {
	var (
		err       error
		prev, cur account.GcpCredentials
	)

	if err = json.Unmarshal([]byte(old), &prev); err != nil {
		return false
	}

	if err = json.Unmarshal([]byte(new), &cur); err != nil {
		return false
	}

	return (prev.Type == cur.Type &&
		prev.ProjectId == cur.ProjectId &&
		prev.PrivateKeyId == cur.PrivateKeyId &&
		//prev.PrivateKey == cur.PrivateKey &&       //Commenting this comparison of privateKey to avoid diff on every terraform plan
		prev.ClientEmail == cur.ClientEmail &&
		prev.ClientId == cur.ClientId &&
		prev.AuthUri == cur.AuthUri &&
		prev.TokenUri == cur.TokenUri &&
		prev.ProviderCertUrl == cur.ProviderCertUrl &&
		prev.ClientCertUrl == cur.ClientCertUrl)
}

func parseCloudAccount(d *schema.ResourceData) (string, string, interface{}) {
	if x := ResourceDataInterfaceMap(d, account.TypeAws); len(x) != 0 {
		return account.TypeAws, x["name"].(string), account.Aws{
			AccountId:      x["account_id"].(string),
			Enabled:        x["enabled"].(bool),
			ExternalId:     x["external_id"].(string),
			GroupIds:       SetToStringSlice(x["group_ids"].(*schema.Set)),
			Name:           x["name"].(string),
			RoleArn:        x["role_arn"].(string),
			ProtectionMode: x["protection_mode"].(string),
			AccountType:    x["account_type"].(string),
		}
	} else if x := ResourceDataInterfaceMap(d, account.TypeAzure); len(x) != 0 {
		return account.TypeAzure, x["name"].(string), account.Azure{
			Account: account.CloudAccount{
				AccountId:      x["account_id"].(string),
				Enabled:        x["enabled"].(bool),
				GroupIds:       SetToStringSlice(x["group_ids"].(*schema.Set)),
				Name:           x["name"].(string),
				ProtectionMode: x["protection_mode"].(string),
				AccountType:    x["account_type"].(string),
			},
			ClientId:           x["client_id"].(string),
			Key:                x["key"].(string),
			MonitorFlowLogs:    x["monitor_flow_logs"].(bool),
			TenantId:           x["tenant_id"].(string),
			ServicePrincipalId: x["service_principal_id"].(string),
		}
	} else if x := ResourceDataInterfaceMap(d, account.TypeGcp); len(x) != 0 {
		var creds account.GcpCredentials
		_ = json.Unmarshal([]byte(x["credentials_json"].(string)), &creds)

		return account.TypeGcp, x["name"].(string), account.Gcp{
			Account: account.CloudAccount{
				AccountId:      x["account_id"].(string),
				Enabled:        x["enabled"].(bool),
				GroupIds:       SetToStringSlice(x["group_ids"].(*schema.Set)),
				Name:           x["name"].(string),
				ProtectionMode: x["protection_mode"].(string),
				AccountType:    x["account_type"].(string),
			},
			CompressionEnabled:     x["compression_enabled"].(bool),
			DataflowEnabledProject: x["dataflow_enabled_project"].(string),
			FlowLogStorageBucket:   x["flow_log_storage_bucket"].(string),
			Credentials:            creds,
		}
	} else if x := ResourceDataInterfaceMap(d, account.TypeAlibaba); len(x) != 0 {
		return account.TypeAlibaba, x["name"].(string), account.Alibaba{
			AccountId: x["account_id"].(string),
			GroupIds:  SetToStringSlice(x["group_ids"].(*schema.Set)),
			Name:      x["name"].(string),
			RamArn:    x["ram_arn"].(string),
			Enabled:   x["enabled"].(bool),
		}
	}

	return "", "", nil
}

func saveCloudAccount(d *schema.ResourceData, dest string, obj interface{}) {
	var val map[string]interface{}

	switch v := obj.(type) {
	case account.Aws:
		val = map[string]interface{}{
			"account_id":      v.AccountId,
			"enabled":         v.Enabled,
			"external_id":     v.ExternalId,
			"group_ids":       v.GroupIds,
			"name":            v.Name,
			"role_arn":        v.RoleArn,
			"protection_mode": v.ProtectionMode,
			"account_type":    v.AccountType,
		}
	case account.Azure:
		x := ResourceDataInterfaceMap(d, account.TypeAzure)
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
			"client_id":            v.ClientId,
			"key":                  key,
			"monitor_flow_logs":    v.MonitorFlowLogs,
			"tenant_id":            v.TenantId,
			"service_principal_id": v.ServicePrincipalId,
			"protection_mode":      v.Account.ProtectionMode,
			"account_type":         v.Account.AccountType,
		}
	case account.Gcp:
		b, _ := json.Marshal(v.Credentials)
		val = map[string]interface{}{
			"account_id":               v.Account.AccountId,
			"enabled":                  v.Account.Enabled,
			"group_ids":                v.Account.GroupIds,
			"name":                     v.Account.Name,
			"compression_enabled":      v.CompressionEnabled,
			"dataflow_enabled_project": v.DataflowEnabledProject,
			"flow_log_storage_bucket":  v.FlowLogStorageBucket,
			"credentials_json":         string(b),
			"protection_mode":          v.Account.ProtectionMode,
			"account_type":             v.Account.AccountType,
		}
	case account.Alibaba:
		val = map[string]interface{}{
			"account_id": v.AccountId,
			"group_ids":  v.GroupIds,
			"name":       v.Name,
			"ram_arn":    v.RamArn,
			"enabled":    v.Enabled,
		}
	}

	for _, key := range []string{account.TypeAws, account.TypeAzure, account.TypeGcp, account.TypeAlibaba} {
		if key != dest {
			d.Set(key, nil)
			continue
		}

		if err := d.Set(key, []interface{}{val}); err != nil {
			log.Printf("[WARN] Error setting %q field for %q: %s", key, d.Id(), err)
		}
	}
}

func createCloudAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)
	cloudType, name, obj := parseCloudAccount(d)

	if err := account.Create(client, obj); err != nil {
		if strings.Contains(err.Error(), "duplicate_cloud_account") {
			if err := account.Update(client, obj); err != nil {
				return diag.FromErr(err)
			}
		} else {
			return diag.FromErr(err)
		}
	}

	PollApiUntilSuccess(func() error {
		_, err := account.Identify(client, cloudType, name)
		return err
	})

	id, err := account.Identify(client, cloudType, name)
	if err != nil {
		return diag.FromErr(err)
	}

	PollApiUntilSuccess(func() error {
		_, err := account.Get(client, cloudType, id)
		return err
	})

	d.SetId(TwoStringsToId(cloudType, id))
	return readCloudAccount(ctx, d, meta)
}

func readCloudAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)
	cloudType, id := IdToTwoStrings(d.Id())

	obj, err := account.Get(client, cloudType, id)
	if err != nil {
		if err == pc.ObjectNotFoundError {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	saveCloudAccount(d, cloudType, obj)

	return nil
}

func updateCloudAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)

	_, _, obj := parseCloudAccount(d)

	if err := account.Update(client, obj); err != nil {
		return diag.FromErr(err)
	}

	return readCloudAccount(ctx, d, meta)
}

func deleteCloudAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*pc.Client)
	cloudType, id := IdToTwoStrings(d.Id())
	disable := d.Get("disable_on_destroy").(bool)

	if disable {
		switch cloudType {
		case account.TypeAws:
			cloudAccount, _ := account.Get(client, cloudType, id)
			cloudAccountAws := cloudAccount.(account.Aws)
			cloudAccountAws.Enabled = false
			if err := account.Update(client, cloudAccountAws); err != nil {
				return diag.FromErr(err)
			}
			return nil

		case account.TypeAzure:
			cloudAccount, _ := account.Get(client, cloudType, id)
			cloudAccountAzure := cloudAccount.(account.Azure)
			cloudAccountAzure.Account.Enabled = false
			if err := account.Update(client, cloudAccountAzure); err != nil {
				return diag.FromErr(err)
			}
			return nil

		case account.TypeGcp:
			cloudAccount, _ := account.Get(client, cloudType, id)
			cloudAccountGcp := cloudAccount.(account.Gcp)
			cloudAccountGcp.Account.Enabled = false
			if err := account.Update(client, cloudAccountGcp); err != nil {
				return diag.FromErr(err)
			}
			return nil

		case account.TypeAlibaba:
			cloudAccount, _ := account.Get(client, cloudType, id)
			cloudAccountAlibaba := cloudAccount.(account.Alibaba)
			cloudAccountAlibaba.Enabled = false
			if err := account.Update(client, cloudAccountAlibaba); err != nil {
				return diag.FromErr(err)
			}
			return nil
		}
	}

	err := account.Delete(client, cloudType, id)
	if err != nil {
		if err != pc.ObjectNotFoundError {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}
