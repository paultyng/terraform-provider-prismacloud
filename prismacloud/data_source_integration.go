package prismacloud

import (
	pc "github.com/paloaltonetworks/prisma-cloud-go"
	"github.com/paloaltonetworks/prisma-cloud-go/integration"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIntegration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIntegrationRead,

		Schema: map[string]*schema.Schema{
			// Input.
			"integration_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Integration ID",
				AtLeastOneOf: []string{"integration_id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Name of the integration",
				AtLeastOneOf: []string{"integration_id", "name"},
			},

			// Output.
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description",
			},
			"integration_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Integration type",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enabled",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created by",
			},
			"created_ts": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Created on",
			},
			"last_modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified by",
			},
			"last_modified_ts": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Last modified timestamp",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status",
			},
			"valid": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Valid",
			},
			"reason": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Model for the integration status details",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_updated": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Last updated",
						},
						"error_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Error type",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message",
						},
						"details": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Model for message details",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Status code",
									},
									"subject": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subject",
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Internationalization key",
									},
								},
							},
						},
					},
				},
			},
			"integration_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Integration configuration, the values depend on the integration type",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"queue_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Queue URL you used when you configured Prisma Cloud in Amazon SQS",
						},
						"login": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "(Qualys/ServiceNow) Login",
						},
						"base_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Qualys Security Operations Center server API URL (without \"http(s)\")",
						},
						"password": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "(Qualys/ServiceNow) Password",
						},
						"host_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ServiceNow URL",
						},
						"tables": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Key/value pairs that identify the ServiceNow module tables with which to integrate (e.g. - incident, sn_si_incident, or em_event)",
							Elem: &schema.Schema{
								Type: schema.TypeBool,
							},
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ServiceNow release version",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Webhook URL or Splunk HTTP event collector URL",
						},
						"headers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Webhook headers",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Header name",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Header value",
									},
									"secure": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Secure",
									},
									"read_only": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Read only",
									},
								},
							},
						},
						"auth_token": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PagerDuty authentication token for the event collector",
						},
						"integration_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PagerDuty/Splunk integration key",
						},
						"source_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GCP Source ID for Google CSCC integration",
						},
						"org_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GCP Organization ID for Google CSCC integration",
						},
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS account ID",
						},
						"regions": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "AWS regions",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AWS region name",
									},
									"api_identifier": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AWS region code",
									},
									"cloud_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cloud Type",
									},
									"sdk_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "SDK ID",
									},
								},
							},
						},
						"s3_uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS S3 URI",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS region",
						},
						"role_arn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS role ARN",
						},
						"external_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS external ID",
						},
						"roll_up_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "File Roll Up Time",
						},
						"source_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source type for splunk integration",
						},
					},
				},
			},
		},
	}
}

func dataSourceIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pc.Client)

	var err error
	id := d.Get("integration_id").(string)

	if id == "" {
		name := d.Get("name").(string)
		id, err = integration.Identify(client, name)
		if err != nil {
			if err == pc.ObjectNotFoundError {
				d.SetId("")
				return nil
			}
			return err
		}
	}

	o, err := integration.Get(client, id)
	if err != nil {
		if err == pc.ObjectNotFoundError {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(o.Id)
	saveIntegration(d, o)

	return nil
}
