package create_skelton_test

import (
	"encoding/json"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	cs "github.com/mikiya771/terraform-snippetter/internal/create_skelton"
)

type SkeltonTestData struct {
	attrName  string
	inputJson string
	expected  string
}

func TestCreateLine(t *testing.T) {
	testSkelton := map[string]SkeltonTestData{
		"computedJson": {
			attrName: "name",
			inputJson: `
{
	"type": "string",
	"description_kind": "plain",
	"computed": true
}
			`,
			expected: "## name = \"${string}\"",
		},
	}
	for k, v := range testSkelton {
		t.Run(k, func(t *testing.T) {
			var sa tfjson.SchemaAttribute
			json.Unmarshal([]byte(v.inputJson), &sa)
			if v.expected != cs.CreateLine(v.attrName, &sa) {
				t.Errorf("Error: CreateLine(%s, %s) should return %s, but return %s", v.attrName, v.inputJson, v.expected, cs.CreateLine(v.attrName, &sa))
			}
		})
	}
}

func TestCreateNestedLine(t *testing.T) {
	testSkelton := map[string]SkeltonTestData{
		"computedJson": {
			attrName: "assume_role",
			inputJson: `
{
	"nesting_mode": "list",
	"block": {
		"attributes":{
			"policy_arns": {
				"type": [
					"set",
					"string"
				],
				"description": "Amazon Resource Names (ARNs) of IAM Policies describing further restricting permissions for the IAM Role being assumed.",
				"description_kind": "plain",
				"optional": true
			},
			"role_arn": {
				"type": "string",
				"description": "Amazon Resource Name of an IAM Role to assume prior to making API calls.",
				"description_kind": "plain",
				"optional": true
			},
			"tags": {
				"type": [
					"map",
					"string"
				],
				"description": "Assume role session tags.",
				"description_kind": "plain",
				"optional": true
			}
		},
		"description_kind": "plain"
	},
	"max_items": 1
}
			`,
			expected: "## name = \"${string}\"",
		},
	}
}
