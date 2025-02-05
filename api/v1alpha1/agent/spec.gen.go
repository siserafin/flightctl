// Package v1alpha1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package v1alpha1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	externalRef0 "github.com/flightctl/flightctl/api/v1alpha1"
	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	pathPrefix := path.Dir(pathToFile)

	for rawPath, rawFunc := range externalRef0.PathToRawSpec(path.Join(pathPrefix, "../openapi.yaml")) {
		if _, ok := res[rawPath]; ok {
			// it is not possible to compare functions in golang, so always overwrite the old value
		}
		res[rawPath] = rawFunc
	}
	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
