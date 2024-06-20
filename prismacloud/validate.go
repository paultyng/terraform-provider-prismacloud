package prismacloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"gopkg.in/yaml.v2"
)

// Validates build policy YAML and the policy.Rule.Children object
func ValidateMetadata(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	m := v.(map[string]interface{})

	if val, ok := m["code"]; ok {
		log.Printf("[DEBUG] YAML for build policy: \n %v", val)
		if s, err := validateYAMLString(val); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Invalid YAML string",
				Detail:        fmt.Sprintf("\"code\" should be a valid YAML/JSON string. Got: \n%v", s),
				AttributePath: path,
			})
		}
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Bad map key",
			Detail:        fmt.Sprintf("Build policy definition should be assigned to the \"code\" key"),
			AttributePath: path,
		})
	}
	return diags
}

func validateYAMLString(yamlString interface{}) (string, error) {
	var y interface{}
	if yamlString == nil || yamlString.(string) == "" {
		return "", nil
	}
	s := yamlString.(string)
	err := yaml.Unmarshal([]byte(s), &y)
	return s, err
}
