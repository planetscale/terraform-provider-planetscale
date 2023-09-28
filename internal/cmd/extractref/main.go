package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/pkg/browser"
	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/exp/slog"
)

func main() {
	cfgFilepath := flag.String("cfg", "../../../openapi/extract-ref-cfg.json", "")
	specFilepath := flag.String("spec", "../../../openapi/openapi-spec.json", "")
	flag.Parse()

	if err := realMain(*cfgFilepath, *specFilepath); err != nil {
		slog.Error("failed", "err", err.Error())
	}
}

type ExtractConfig struct {
	Extractions []ExtractRule `json:"extractions"`
}

type ExtractRule struct {
	Path     string `json:"path"`
	Method   string `json:"method"`
	Response int    `json:"responses"`
	Prop     string `json:"prop"`

	BecomeRef string `json:"become_ref"`
}

func readCfg(filepath string) (*ExtractConfig, error) {
	var cfg ExtractConfig
	cfgRaw, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return &cfg, json.Unmarshal(cfgRaw, &cfg)
}

func readSpec(filepath string) (*spec.Swagger, error) {
	var (
		file *os.File
		err  error
	)
	if filepath == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(filepath)
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var spec spec.Swagger
	if err := json.NewDecoder(file).Decode(&spec); err != nil {
		return nil, fmt.Errorf("decoding JSOn: %w", err)
	}
	return &spec, nil
}

func realMain(cfgFile, specFile string) error {
	slog.Info("loading cfg")
	cfg, err := readCfg(cfgFile)
	if err != nil {
		return fmt.Errorf("loading cfg: %w", err)
	}
	slog.Info("loading spec")

	spec, err := readSpec(specFile)
	if err != nil {
		return fmt.Errorf("loading spec: %w", err)
	}

	slog.Info("loaded openapi spec", "spec.id", spec.ID)

	for _, extraction := range cfg.Extractions {
		slog.Info("applying extraction rule", "path", extraction.Path)
		p, ok := spec.Paths.Paths[extraction.Path]
		if !ok {
			return fmt.Errorf("path doesn't exist in openapi spec: %q", extraction.Path)
		}
		if err := handlePath(spec, extraction, p); err != nil {
			return fmt.Errorf("handling rule for path %q: %v", extraction.Path, err)
		}
	}
	slog.Info("encoding modified spec")

	return json.NewEncoder(os.Stdout).Encode(spec)
}

func handlePath(doc *spec.Swagger, rule ExtractRule, path spec.PathItem) error {
	ref, err := spec.NewRef("#/definitions/" + rule.BecomeRef)
	if err != nil {
		return fmt.Errorf("invalid `become_ref` rule: %v", err)
	}

	var op *spec.Operation
	switch strings.ToUpper(rule.Method) {
	case "GET":
		op = path.Get
	case "PUT":
		op = path.Put
	case "POST":
		op = path.Post
	case "DELETE":
		op = path.Delete
	case "OPTIONS":
		op = path.Options
	case "HEAD":
		op = path.Head
	case "PATCH":
		op = path.Patch
	default:
		return fmt.Errorf("unsupported method %q", rule.Method)
	}
	if op == nil {
		return fmt.Errorf("no definition for method %q", rule.Method)
	}

	resp, ok := op.Responses.StatusCodeResponses[rule.Response]
	if !ok {
		return fmt.Errorf("response doesn't support code %d", rule.Response)
	}
	if resp.Schema == nil {
		return fmt.Errorf("response at this path has no schema")
	}

	pathParts := strings.Split(rule.Prop, ".")
	if len(pathParts) == 1 && pathParts[0] == "" {
		pathParts = nil
	}

	tgt, err := resolvePath(pathParts, resp.Schema, func(path string, parent, schema *spec.Schema) {
		// replace the tgt schema with the ref
		switch {
		case parent.Type.Contains("array"):
			parent.Items.Schema = spec.RefSchema(ref.String())
		case parent.Type.Contains("object"):
			if path == "" {
				// the root schema itself is changed
				desc := resp.Schema.Description
				nullable := resp.Schema.Nullable
				resp.Schema = spec.RefSchema(ref.String())
				resp.Schema.Description = desc
				resp.Schema.Nullable = nullable
			} else {
				parent.Properties[path] = *spec.RefSchema(ref.String())
			}
		default:
			panic(fmt.Sprintf("unhandled case: %#v", parent.Type))
		}
	})
	if err != nil {
		return fmt.Errorf("resolving prop at path %q: %v", rule.Prop, err)
	}
	op.Responses.StatusCodeResponses[rule.Response] = resp
	existingDef, ok := doc.Definitions[rule.BecomeRef]
	if !ok {
		doc.Definitions[rule.BecomeRef] = *tgt
	} else {
		oldDef, err := json.MarshalIndent(existingDef, "", "   ")
		if err != nil {
			return fmt.Errorf("encoding existing def %q: %v", rule.BecomeRef, err)
		}
		newDef, err := json.MarshalIndent(*tgt, "", "   ")
		if err != nil {
			return fmt.Errorf("encoding new def %q: %v", rule.BecomeRef, err)
		}
		if !bytes.Equal(oldDef, newDef) {
			slog.Error("old definition", "def", string(oldDef))
			slog.Error("new definition", "def", string(newDef))
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(string(oldDef), string(newDef), false)
			if err := browser.OpenReader(bytes.NewBufferString(dmp.DiffPrettyHtml(diffs))); err != nil {
				panic(err)
			}
			return fmt.Errorf("duplicate reference to %q, using a non-equal schema definition", rule.BecomeRef)
		}
	}

	return nil
}

func resolvePath(pathParts []string, schema *spec.Schema, atTarget func(path string, parent, schema *spec.Schema)) (*spec.Schema, error) {
	return resolvePathRecurse(pathParts, "", schema, schema, atTarget)
}

func resolvePathRecurse(
	pathParts []string,
	atPath string,
	parent,
	schema *spec.Schema,
	atTarget func(path string, parent, schema *spec.Schema),
) (*spec.Schema, error) {
	if schema.Type.Contains("array") {
		if schema.Items.Schema == nil {
			return nil, fmt.Errorf("path is an array and its `items` schema isn't unitary")
		}
		parent = schema
		schema = schema.Items.Schema
	}

	if len(pathParts) == 0 {
		atTarget(atPath, parent, schema)
		return schema, nil
	}
	currentPath := pathParts[0]

	var (
		nextSchema spec.Schema
		ok         bool
	)

	nextSchema, ok = schema.Properties[currentPath]
	if !ok {
		return nil, fmt.Errorf("path %q doesn't exist", currentPath)
	}
	return resolvePathRecurse(pathParts[1:], currentPath, schema, &nextSchema, atTarget)
}
