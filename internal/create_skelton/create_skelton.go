package create_skelton

import (
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
)

func CreateSnippets(ps *tfjson.ProviderSchema) (map[string][]string, map[string][]string) {
	resni := map[string][]string{}
	dsni := map[string][]string{}
	for k, v := range ps.ResourceSchemas {
		snipList := []string{}
		startSnip := fmt.Sprintf("snippet r_%s \"terraform %s snip\"", k, k)
		endSnip := fmt.Sprintf("endsnippet")
		ret := CreateSource("resource", k, v)
		snipList = append(snipList, startSnip)
		snipList = append(snipList, ret...)
		snipList = append(snipList, endSnip)
		resni[k] = snipList
	}
	for k, v := range ps.DataSourceSchemas {
		snipList := []string{}
		startSnip := fmt.Sprintf("snippet d_%s \"terraform %s snip\"", k, k)
		endSnip := fmt.Sprintf("endsnippet")
		ret := CreateSource("data", k, v)
		snipList = append(snipList, startSnip)
		snipList = append(snipList, ret...)
		snipList = append(snipList, endSnip)
		dsni[k] = snipList
	}
	return resni, dsni
}

func CreateSource(sourceType string, keyName string, schema *tfjson.Schema) []string {
	out := CreateSchemaBlockLine(schema.Block)
	startDoc := fmt.Sprintf("%s %s ${$1:Usage} {", sourceType, keyName)
	endDoc := "}"
	ret := []string{}
	ret = append(ret, startDoc)
	ret = append(ret, AddTabStrings(out)...)
	ret = append(ret, endDoc)
	return ret
}

//non nested TODO: composit type
func CreateLine(keyName string, schemaAttr *tfjson.SchemaAttribute) []string {
	friendlyType := schemaAttr.AttributeType.FriendlyName()
	return []string{
		fmt.Sprintf("%s = ${%s}", keyName, friendlyType),
	}
}

//nested
func CreateNestedLine(keyName string, nestedSchemaAttr *tfjson.SchemaBlockType) []string {
	nestingMode := nestedSchemaAttr.NestingMode
	attrs := CreateSchemaBlockLine(nestedSchemaAttr.Block)
	comment := fmt.Sprintf("## NestingMode: %s, MinItems: %d, MaxItemd %d", nestingMode, nestedSchemaAttr.MinItems, nestedSchemaAttr.MaxItems)
	startDoc := ""
	endDoc := ""
	switch nestingMode {
	case "single", "list", "set":
		startDoc = fmt.Sprintf("%s {", keyName)
		endDoc = "}"
	case "map":
		panic("unexpected nestingMode [single, list, set, ] is neededd")
	default:
		panic("unexpected nestingMode [single, list, set, ] is neededd")
	}
	ret := []string{}
	ret = append(ret, comment)
	ret = append(ret, startDoc)
	for _, k := range attrs {
		ret = append(ret, fmt.Sprintf("\t%s", k))
	}
	ret = append(ret, endDoc)
	return ret
}

func CreateSchemaBlockLine(schemaBlock *tfjson.SchemaBlock) []string {
	a, o := CreateAttrbutesLine(schemaBlock.Attributes)
	blockNote := []string{}
	originBlock := schemaBlock.NestedBlocks
	for k, bt := range originBlock {
		blockNote = append(blockNote, CreateNestedLine(k, bt)...)
	}
	r := append(a, blockNote...)
	r = append(r, o...)
	return r
}

func CreateAttrbutesLine(am map[string]*tfjson.SchemaAttribute) ([]string, []string) {
	required := []string{}
	optional := []string{}
	output := []string{}
	optional = append(optional, "## optional")
	output = append(output, "## output")

	for key, sb := range am {
		if sb.Required {
			required = append(required, CreateLine(key, sb)...)
		} else if sb.Optional {
			optional = append(optional, CreateLine(key, sb)...)
		} else {
			output = append(output, CreateLine(key, sb)...)
		}
	}
	optional = CommentOutStrings(optional)
	output = CommentOutStrings(output)
	return append(required, optional...), output
}

func CommentOutStrings(input []string) []string {
	o := []string{}
	for _, in := range input {
		o = append(o, fmt.Sprintf("## %s", in))
	}
	return o
}

func AddTabStrings(input []string) []string {
	o := []string{}
	for _, in := range input {
		o = append(o, fmt.Sprintf("\t %s", in))
	}
	return o
}
