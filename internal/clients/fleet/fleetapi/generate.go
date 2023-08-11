//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	fleetSchemaURLTmpl = "https://raw.githubusercontent.com/elastic/kibana/%s/x-pack/plugins/fleet/common/openapi/bundled.json"
)

type OpenAPISchema struct {
	Paths          map[string]*Path `json:"paths"`
	OpenAPIVersion string           `json:"openapi"`
	Tags           []any            `json:"tags,omitempty"`
	Servers        []any            `json:"servers,omitempty"`
	Components     Fields           `json:"components,omitempty"`
	Security       []any            `json:"security,omitempty"`
	Info           map[string]any   `json:"info"`
}

type Path struct {
	Parameters []Fields  `json:"parameters,omitempty"`
	Get        *Endpoint `json:"get,omitempty"`
	Post       *Endpoint `json:"post,omitempty"`
	Put        *Endpoint `json:"put,omitempty"`
	Delete     *Endpoint `json:"delete,omitempty"`
}

type Endpoint struct {
	Summary     string   `json:"summary,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Responses   Fields   `json:"responses,omitempty"`
	RequestBody Fields   `json:"requestBody,omitempty"`
	OperationID string   `json:"operationId,omitempty"`
	Parameters  []Fields `json:"parameters,omitempty"`
	Deprecated  bool     `json:"deprecated,omitempty"`
}

var includePaths = map[string][]string{
	"/agent_policies":                      {"post"},
	"/agent_policies/{agentPolicyId}":      {"get", "put"},
	"/agent_policies/delete":               {"post"},
	"/enrollment_api_keys":                 {"get"},
	"/fleet_server_hosts":                  {"post"},
	"/fleet_server_hosts/{itemId}":         {"get", "put", "delete"},
	"/outputs":                             {"post"},
	"/outputs/{outputId}":                  {"get", "put", "delete"},
	"/package_policies":                    {"post"},
	"/package_policies/{packagePolicyId}":  {"get", "put", "delete"},
	"/epm/packages/{pkgName}/{pkgVersion}": {"get", "put", "post", "delete"},
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: HTTP %v: %v", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func stringInSlice(value string, slice []string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}

// filterKbnXsrfParameter filters out an entry if it is a kbn_xsrf parameter.
// Returns a copy of the slice if it was modified, otherwise returns the original
// slice if no match was found.
func filterKbnXsrfParameter(parameters []Fields) []Fields {
	removeIndex := -1

	for i, param := range parameters {
		if ref, ok := param["$ref"].(string); ok && ref == "#/components/parameters/kbn_xsrf" {
			removeIndex = i
			break
		}
	}
	if removeIndex != -1 {
		ret := make([]Fields, 0)
		ret = append(ret, parameters[:removeIndex]...)
		return append(ret, parameters[removeIndex+1:]...)
	}

	return parameters
}

// transformSchemasInputsType transforms the "inputs" property on the
// "new_package_policy" component schema from an array to an object.
func transformSchemasInputsType(schema *OpenAPISchema) {
	inputs, ok := schema.Components.GetFields("schemas.new_package_policy.properties.inputs")
	if !ok {
		return
	}

	inputs.Set("items.properties.streams.type", "object")

	inputs.Set("type", "object")
	inputs.Move("items", "additionalProperties")

	// Drop package_policies from Agent Policy
	agentPolicy, _ := schema.Components.GetFields("schemas.agent_policy")
	agentPolicy.Delete("properties.package_policies")
}

// transformInlinePackageDefinitions relocates inline type definitions for the
// EPM endpoints to the dedicated schemas section of the OpenAPI schema. This needs
// to be done as there is a bug in the OpenAPI generator which causes types to
// be generated with invalid names:
// https://github.com/deepmap/oapi-codegen/issues/1121
func transformInlinePackageDefinitions(schema *OpenAPISchema) {
	epmPath, ok := schema.Paths["/epm/packages/{pkgName}/{pkgVersion}"]
	if !ok {
		panic("epm path not found")
	}

	// Get
	{
		props, ok := epmPath.Get.Responses.GetFields("200.content.application/json.schema.allOf.1.properties")
		if !ok {
			panic("properties not found")
		}

		// status needs to be moved to schemes and a ref inserted in its place.
		value, _ := props.Get("status")
		schema.Components.Set("schemas.package_status", value)
		props.Delete("status")
		props.Set("status.$ref", "#/components/schemas/package_status")
	}

	// Post
	{
		props, ok := epmPath.Post.Responses.GetFields("200.content.application/json.schema.properties")
		if !ok {
			panic("properties not found")
		}

		// _meta.properties.install_source
		value, _ := props.GetFields("_meta.properties.install_source")
		schema.Components.Set("schemas.package_install_source", value)
		props.Delete("_meta.properties.install_source")
		props.Set("_meta.properties.install_source.$ref", "#/components/schemas/package_install_source")

		// items.items.properties.type
		value, _ = props.GetFields("items.items.properties.type")
		schema.Components.Set("schemas.package_item_type", value)
		props.Delete("items.items.properties.type")
		props.Set("items.items.properties.type.$ref", "#/components/schemas/package_item_type")
	}

	// Put
	{
		props, ok := epmPath.Put.Responses.GetFields("200.content.application/json.schema.properties")
		if !ok {
			panic("properties not found")
		}

		// items.items.properties.type (definition already moved by Post)
		props.Delete("items.items.properties.type")
		props.Set("items.items.properties.type.$ref", "#/components/schemas/package_item_type")
	}

	// Delete
	{
		props, ok := epmPath.Delete.Responses.GetFields("200.content.application/json.schema.properties")
		if !ok {
			panic("properties not found")
		}

		// items.items.properties.type (definition already moved by Post)
		props.Delete("items.items.properties.type")
		props.Set("items.items.properties.type.$ref", "#/components/schemas/package_item_type")
	}

	// Move embedded objects (structs) to schemas so Go-types are generated.
	{
		// package_policy_request_input_stream
		field, _ := schema.Components.GetFields("schemas.package_policy_request.properties.inputs.additionalProperties.properties.streams")
		props, _ := field.Get("additionalProperties")
		schema.Components.Set("schemas.package_policy_request_input_stream", props)
		field.Delete("additionalProperties")
		field.Set("additionalProperties.$ref", "#/components/schemas/package_policy_request_input_stream")

		// package_policy_request_input
		field, _ = schema.Components.GetFields("schemas.package_policy_request.properties.inputs")
		props, _ = field.Get("additionalProperties")
		schema.Components.Set("schemas.package_policy_request_input", props)
		field.Delete("additionalProperties")
		field.Set("additionalProperties.$ref", "#/components/schemas/package_policy_request_input")

		// package_policy_package_info
		field, _ = schema.Components.GetFields("schemas.new_package_policy.properties")
		props, _ = field.Get("package")
		schema.Components.Set("schemas.package_policy_package_info", props)
		field.Delete("package")
		field.Set("package.$ref", "#/components/schemas/package_policy_package_info")

		// package_policy_input
		field, _ = schema.Components.GetFields("schemas.new_package_policy.properties.inputs")
		props, _ = field.Get("additionalProperties")
		schema.Components.Set("schemas.package_policy_input", props)
		field.Delete("additionalProperties")
		field.Set("additionalProperties.$ref", "#/components/schemas/package_policy_input")
	}
}

func main() {
	outFile := flag.String("o", "", "output file")
	inFile := flag.String("i", "", "input file")
	apiVersion := flag.String("v", "main", "api version")
	flag.Parse()

	if *outFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	var err error
	var schemaData []byte

	if *inFile != "" {
		schemaData, err = os.ReadFile(*inFile)
	} else {
		schemaData, err = downloadFile(fmt.Sprintf(fleetSchemaURLTmpl, *apiVersion))
	}
	if err != nil {
		log.Fatal(err)
	}

	var schema OpenAPISchema
	if err = json.Unmarshal(schemaData, &schema); err != nil {
		log.Fatal(err)
	}

	for path, pathInfo := range schema.Paths {
		// Remove paths not in filter list.
		if _, exists := includePaths[path]; !exists {
			delete(schema.Paths, path)
			continue
		}

		// Filter out kbn-xsrf parameter (already set by API client).
		pathInfo.Parameters = filterKbnXsrfParameter(pathInfo.Parameters)

		// Filter out endpoints not if filter list, filter out kbn-xsrf
		// parameter in endpoint (already set by API client).
		allowedMethods := includePaths[path]
		filterEndpointFn := func(endpoint *Endpoint, method string) *Endpoint {
			if endpoint == nil {
				return nil
			}
			if !stringInSlice(method, allowedMethods) {
				return nil
			}

			endpoint.Parameters = filterKbnXsrfParameter(endpoint.Parameters)

			return endpoint
		}
		pathInfo.Get = filterEndpointFn(pathInfo.Get, "get")
		pathInfo.Post = filterEndpointFn(pathInfo.Post, "post")
		pathInfo.Put = filterEndpointFn(pathInfo.Put, "put")
		pathInfo.Delete = filterEndpointFn(pathInfo.Delete, "delete")
	}

	transformSchemasInputsType(&schema)
	transformInlinePackageDefinitions(&schema)

	outData, err := json.MarshalIndent(&schema, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile(*outFile, outData, 0664); err != nil {
		log.Fatal(err)
	}
}

// Fields wraps map[string]any with convenience functions for interacting
// with nested map values.
type Fields map[string]any

// Get will get the value at 'key' as the first returned
// parameter. The second parameter is a bool indicating
// if 'key' exists.
func (f Fields) Get(key string) (any, bool) {
	indexSliceFn := func(slice []any, key string) (any, string, bool) {
		indexStr, subKeys, _ := strings.Cut(key, ".")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			log.Printf("Failed to parse slice index key %q: %v", indexStr, err)
			return nil, "", false
		}

		if index < 0 || index >= len(slice) {
			log.Printf("Slice index is out of bounds (%d, target slice len: %d)", index, len(slice))
			return nil, "", false
		}

		return slice[index], subKeys, true
	}

	rootKey, subKeys, split := strings.Cut(key, ".")
	if split {
		switch t := f[rootKey].(type) {
		case Fields:
			return t.Get(subKeys)
		case map[string]any:
			return Fields(t).Get(subKeys)
		case []any:
			slicedValue, postSliceKeys, ok := indexSliceFn(t, subKeys)
			if !ok {
				return nil, false
			}
			if m, isMap := slicedValue.(map[string]any); ok && isMap {
				return Fields(m).Get(postSliceKeys)
			}
			return slicedValue, true

		default:
			rootKey = key
		}
	}

	value, ok := f[rootKey]
	return value, ok
}

// GetFields is like Get, but converts the found value to Fields.
// If the key is not found or the type conversion fails, the
// second return value will be false.
func (f Fields) GetFields(key string) (Fields, bool) {
	value, ok := f.Get(key)
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case Fields:
		return t, true
	case map[string]any:
		return t, true
	}

	return nil, false
}

// Set will set key to the value of 'value'.
func (f Fields) Set(key string, value any) {
	rootKey, subKeys, split := strings.Cut(key, ".")
	if split {
		if v, ok := f[rootKey]; ok {
			switch t := v.(type) {
			case Fields:
				t.Set(subKeys, value)
			case map[string]any:
				Fields(t).Set(subKeys, value)
			}
		} else {
			subMap := Fields{}
			subMap.Set(subKeys, value)
			f[rootKey] = subMap
		}
	} else {
		f[rootKey] = value
	}
}

// Move will move the value from 'key' to 'target'. If 'key' does not
// exist, the operation is a no-op.
func (f Fields) Move(key, target string) {
	value, ok := f.Get(key)
	if !ok {
		return
	}

	f.Set(target, value)
	f.Delete(key)
}

// Delete will remove the key from the Fields. If key is nested,
// empty sub-keys will be removed as well.
func (f Fields) Delete(key string) {
	rootKey, subKeys, split := strings.Cut(key, ".")
	if split {
		if v, ok := f[rootKey]; ok {
			switch t := v.(type) {
			case Fields:
				t.Delete(subKeys)
			case map[string]any:
				Fields(t).Delete(subKeys)
			}
		}
	} else {
		delete(f, rootKey)
	}
}
