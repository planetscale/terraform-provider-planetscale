package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slog"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	specFilepath := flag.String("spec", "../../../openapi-spec.json", "")
	flag.Parse()

	if err := realMain(*specFilepath); err != nil {
		slog.Error("failed", "err", err)
	}
}

func realMain(specFilepath string) error {
	slog.Info("loading spec")
	doc, err := loads.Spec(specFilepath)
	if err != nil {
		return fmt.Errorf("loading spec: %w", err)
	}

	spec := doc.Spec()
	slog.Info("loaded openapi spec", "spec.id", spec.ID)

	f := jen.NewFile("planetscale")

	if err := genClientStruct(f, spec); err != nil {
		return fmt.Errorf("generating client struct: %w", err)
	}

	if err := genErrorStruct(f, spec); err != nil {
		return fmt.Errorf("generating error struct: %w", err)
	}

	for name, defn := range spec.Definitions {
		ll := slog.With("definition", name)
		typeName := snakeToCamel(name)
		ll = ll.With("type_name", typeName)
		ll.Info("generating type for definition")
		if err := genParamStruct(spec.Definitions, f, typeName, &defn); err != nil {
			return fmt.Errorf("generating type for definition %q: %w", name, err)
		}
	}

	paths := maps.Keys(spec.Paths.Paths)
	sort.Strings(paths)
	for _, path := range paths {
		ll := slog.With("path", path)
		pathItem := spec.Paths.Paths[path]
		props := pathItem.PathItemProps
		if err := handlePath(ll, spec.Definitions, f, path, props); err != nil {
			return fmt.Errorf("handling path %q: %w", path, err)
		}
	}

	return f.Render(os.Stdout)
}

func handlePath(ll *slog.Logger, defns spec.Definitions, f *jen.File, path string, props spec.PathItemProps) error {
	if props.Get != nil {
		if err := handleVerbPath(ll, defns, f, path, "GET", props.Get); err != nil {
			return fmt.Errorf("handling GET props: %w", err)
		}
	}
	if props.Put != nil {
		if err := handleVerbPath(ll, defns, f, path, "PUT", props.Put); err != nil {
			return fmt.Errorf("handling PUT props: %w", err)
		}
	}
	if props.Post != nil {
		if err := handleVerbPath(ll, defns, f, path, "POST", props.Post); err != nil {
			return fmt.Errorf("handling POST props: %w", err)
		}
	}
	if props.Delete != nil {
		if err := handleVerbPath(ll, defns, f, path, "DELETE", props.Delete); err != nil {
			return fmt.Errorf("handling DELETE props: %w", err)
		}
	}
	if props.Options != nil {
		if err := handleVerbPath(ll, defns, f, path, "OPTIONS", props.Options); err != nil {
			return fmt.Errorf("handling OPTIONS props: %w", err)
		}
	}
	if props.Head != nil {
		if err := handleVerbPath(ll, defns, f, path, "HEAD", props.Head); err != nil {
			return fmt.Errorf("handling HEAD props: %w", err)
		}
	}
	if props.Patch != nil {
		if err := handleVerbPath(ll, defns, f, path, "PATCH", props.Patch); err != nil {
			return fmt.Errorf("handling PATCH props: %w", err)
		}
	}

	return nil
}

func handleVerbPath(ll *slog.Logger, defns spec.Definitions, f *jen.File, path, verb string, operation *spec.Operation) error {
	ll.Info("looking at prop", "verb", verb)
	pathParams, queryParams, reqBody, err := splitParams(operation.Parameters)
	if err != nil {
		return fmt.Errorf("splitting params: %w", err)
	}

	var reqBodyStructName string
	if reqBody != nil {
		reqBodyStructName = kebabToCamel(removeFillerWords(operation.ID)) + "Req"
		if err := genParamStruct(defns, f, reqBodyStructName, reqBody.Schema); err != nil {
			return fmt.Errorf("generating call param struct: %w", err)
		}
	}

	responses := make(map[int]string)
	resCodes := maps.Keys(operation.Responses.StatusCodeResponses)
	successResponseTypes := 0
	for _, code := range resCodes {
		if code < 400 {
			successResponseTypes++
		}
	}
	for _, code := range resCodes {
		resBodyStructName := kebabToCamel(removeFillerWords(operation.ID)) + "Res" + strconv.Itoa(code)
		res := operation.Responses.StatusCodeResponses[code]
		if code < 400 {
			if successResponseTypes == 1 {
				resBodyStructName = kebabToCamel(removeFillerWords(operation.ID)) + "Res"
			}
			respSchema := res.ResponseProps.Schema

			if respSchema != nil && respSchema.Ref.GetURL() != nil && respSchema.Ref.GetURL().Fragment != "" {
				defnName := strings.TrimPrefix(respSchema.Ref.GetURL().Fragment, "/definitions/")
				_, ok := defns[defnName]
				if !ok {
					return fmt.Errorf("no definition with name %q exists in the openapi spec", defnName)
				}
				if err := genEmbedStruct(f, resBodyStructName, snakeToCamel(defnName)); err != nil {
					return fmt.Errorf("generating call response struct: %w", err)
				}

			} else {
				if err := genParamStruct(defns, f, resBodyStructName, res.Schema); err != nil {
					return fmt.Errorf("generating call response struct: %w", err)
				}
			}

		} else {
			if err := genErrRespParamStruct(defns, f, resBodyStructName, res.Schema); err != nil {
				return fmt.Errorf("generating call response struct: %w", err)
			}
		}
		responses[code] = resBodyStructName
	}

	clientCallFuncName := kebabToCamel(removeFillerWords(operation.ID))
	if err := genClientCall(f, path, verb, clientCallFuncName, pathParams, queryParams, reqBodyStructName, responses); err != nil {
		return fmt.Errorf("generating client call method: %w", err)
	}

	return nil
}

func removeFillerWords(name string) string {
	name = strings.ReplaceAll(name, "-an-", "-")
	name = strings.ReplaceAll(name, "-a-", "-")
	return name
}

func kebabToCamel(kebab string) string {
	var out strings.Builder
	for _, w := range strings.Split(kebab, "-") {
		out.WriteString(cases.Title(language.AmericanEnglish).String(w))
	}
	return out.String()
}

func snakeToCamel(snake string) string {
	var out strings.Builder
	for _, w := range strings.Split(snake, "_") {
		out.WriteString(cases.Title(language.AmericanEnglish).String(w))
	}
	return out.String()
}

func lowerSnakeToCamel(snake string) string {
	var out strings.Builder
	for i, w := range strings.Split(snake, "_") {
		if i == 0 {
			out.WriteString(w)
		} else {
			out.WriteString(cases.Title(language.AmericanEnglish).String(w))
		}
	}
	return out.String()
}

func splitParams(params []spec.Parameter) (path, query []spec.Parameter, body *spec.Parameter, err error) {
	for _, param := range params {
		switch param.In {
		case "path":
			path = append(path, param)
		case "query":
			query = append(query, param)
		case "body":
			if body != nil {
				return nil, nil, nil, fmt.Errorf("multiple bodies specified: %q", param.Name)
			}
			if param.Type != "object" && len(param.Schema.Properties) == 0 {
				return nil, nil, nil, fmt.Errorf("body should be an object: was a %q", param.Type)
			}
			body = &param
		default:
			return nil, nil, nil, fmt.Errorf("unhandled param.In: %q", param.In)
		}
	}
	return
}
