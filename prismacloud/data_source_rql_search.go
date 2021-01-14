package prismacloud

import (
	"log"

	pc "github.com/paloaltonetworks/prisma-cloud-go"
	"github.com/paloaltonetworks/prisma-cloud-go/rql/search"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceRqlSearch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRqlSearchRead,

		Schema: map[string]*schema.Schema{
			// Input.
			"search_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The search type",
				Default:      "config",
				ValidateFunc: validation.StringInSlice([]string{"config", "network", "event"}, false),
			},
			"time_range": timeRangeSchema("data_source_rql_search"),
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The RQL search to perform",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Limit results",
				Default:     100,
			},
			"finalized": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set to true when you've finished iterating on a RQL search",
			},
			/*
			   "with_resource_json": {
			       Type: schema.TypeBool,
			       Optional: true,
			       Description: "Return back the resource information",
			   },
			*/

			// Output.
			"group_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Group by",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"search_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The search ID",
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cloud type",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The search name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description",
			},
			"config_data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of config data structs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL",
						},
					},
				},
			},
			"event_data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of event data structs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Region ID",
						},
						"region_api_identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region API identifier",
						},
					},
				},
			},
		},
	}
}

func dataSourceRqlSearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pc.Client)
	finalized := d.Get("finalized").(bool)
	query := d.Get("query").(string)
	limit := d.Get("limit").(int)
	tr := ParseTimeRange(ResourceDataInterfaceMap(d, "time_range"))

	if d.Id() != "" && finalized {
		return nil
	}

	switch d.Get("search_type").(string) {
	case "config":
		req := search.ConfigRequest{
			Query:     query,
			Limit:     limit,
			TimeRange: tr,
		}

		resp, err := search.ConfigSearch(client, req)
		if err != nil {
			return err
		}

		d.SetId(resp.Id)
		d.Set("search_id", resp.Id)
		if err = d.Set("group_by", resp.GroupBy); err != nil {
			log.Printf("[WARN] Error setting 'group_by' for %q: %s", d.Id(), err)
		}
		d.Set("cloud_type", resp.CloudType)
		d.Set("name", resp.Name)
		d.Set("description", resp.Description)
		d.Set("event_data", nil)

		if len(resp.Data.Items) == 0 {
			d.Set("config_data", nil)
		} else {
			list := make([]interface{}, 0, len(resp.Data.Items))
			for _, x := range resp.Data.Items {
				list = append(list, map[string]interface{}{
					"state_id": x.StateId,
					"name":     x.Name,
					"url":      x.Url,
				})
			}

			if err = d.Set("config_data", list); err != nil {
				log.Printf("[WARN] Error setting 'config_data' for %q: %s", d.Id(), err)
			}
		}
	case "network":
		req := search.NetworkRequest{
			Query:     query,
			Limit:     limit,
			TimeRange: tr,
		}

		resp, err := search.NetworkSearch(client, req)
		if err != nil {
			return err
		}

		d.SetId(resp.Id)
		d.Set("search_id", resp.Id)
		if err = d.Set("group_by", resp.GroupBy); err != nil {
			log.Printf("[WARN] Error setting 'group_by' for %q: %s", d.Id(), err)
		}
		d.Set("cloud_type", resp.CloudType)
		d.Set("name", resp.Name)
		d.Set("description", resp.Description)
		d.Set("config_data", nil)
		d.Set("event_data", nil)
	case "event":
		req := search.EventRequest{
			Query:     query,
			Limit:     limit,
			TimeRange: tr,
		}

		resp, err := search.EventSearch(client, req)
		if err != nil {
			return err
		}

		d.SetId(resp.Id)
		d.Set("search_id", resp.Id)
		if err = d.Set("group_by", resp.GroupBy); err != nil {
			log.Printf("[WARN] Error setting 'group_by' for %q: %s", d.Id(), err)
		}
		d.Set("cloud_type", resp.CloudType)
		d.Set("name", resp.Name)
		d.Set("description", resp.Description)
		d.Set("config_data", nil)

		if len(resp.Data.Items) == 0 {
			d.Set("event_data", nil)
		} else {
			list := make([]interface{}, 0, len(resp.Data.Items))
			for _, x := range resp.Data.Items {
				list = append(list, map[string]interface{}{
					"account":               x.Account,
					"region_id":             x.RegionId,
					"region_api_identifier": x.RegionApiIdentifier,
				})
			}

			if err = d.Set("event_data", list); err != nil {
				log.Printf("[WARN] Error setting 'event_data' for %q: %s", d.Id(), err)
			}
		}
	}

	return nil
}