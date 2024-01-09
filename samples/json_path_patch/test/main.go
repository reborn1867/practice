package main

import (
	"encoding/json"
	"fmt"

	"github.wdf.sap.corp/practice-learning/json_path_patch/patch"
)

type Tenant struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func main() {
	raw := `{
		"dataHubSpec": {
		  "auditlog": {
			"hpaConfig": {
			  "enabled": true,
			  "maxReplicas": 10,
			  "minReplicas": 2,
			  "unknown": "hahaha"
			},
			"tenants": [
				{
					"id": "a",
					"name": "test"
				}
			],
			"networkPolicies": {
			  "enabled": false
			}
		  }
		}
	  }
	`
	var m map[string]interface{}
	json.Unmarshal([]byte(raw), &m)

	// fmt.Printf("original object: %+v\n", m)
	if err := patch.PartialPatch(m, "dataHubSpec.auditlog.tenants.0", patch.MergePatchGenerator(func(doc json.RawMessage) (patch interface{}, err error) {
		tenant := Tenant{}
		if err := json.Unmarshal(doc, &tenant); err != nil {
			return nil, err
		}

		fmt.Printf("tenants: %+v\n", tenant)
		tenant.ID = "fake"

		fmt.Printf("\nnew tenants: %+v\n", tenant)
		return tenant, nil
	})); err != nil {
		fmt.Printf("failed to patch, err: %s\n", err)
	}
	// fmt.Printf("final object: %+v\n", m)
}
