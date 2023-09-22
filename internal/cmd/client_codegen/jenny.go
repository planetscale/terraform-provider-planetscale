package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/go-openapi/spec"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func genEmbedStruct(file *jen.File, typename, embedname string) error {
	file.Type().Id(typename).Struct(jen.Id(embedname))
	return nil
}

func genParamStruct(defns spec.Definitions, file *jen.File, typename string, body *spec.Schema) error {
	toField := func(item spec.OrderSchemaItem) (jen.Code, error) {
		fieldName := snakeToCamel(item.Name)
		f := jen.Id(fieldName)
		isOptional := !slices.Contains(body.Required, item.Name)
		if isOptional {
			f = f.Op("*")
		}
		switch {
		case item.Type.Contains("string"):
			f = f.String()
		case item.Type.Contains("number"):
			f = f.Float64()
		case item.Type.Contains("boolean"):
			f = f.Bool()
		case item.Type.Contains("array"):
			itemTypename := ""
			switch {
			case item.Items.Schema.Type.Contains("string"):
				itemTypename += "string"
			case item.Items.Schema.Type.Contains("number"):
				itemTypename += "float64"
			case item.Items.Schema.Type.Contains("boolean"):
				itemTypename += "bool"
			case item.Items.Schema.Type.Contains("object"):
				itemTypename = typename + "_" + fieldName + "Item"

				if err := genParamStruct(defns, file, itemTypename, item.Items.Schema); err != nil {
					return nil, fmt.Errorf("generating child item type: %w", err)
				}
			case item.Items.Schema.Type.Contains("array"):
				return nil, fmt.Errorf("arrays of array aren't supported")
			case item.Items.Schema.Ref.GetURL() != nil || item.Items.Schema.Ref.GetURL().Fragment != "":
				fragment := item.Items.Schema.Ref.GetURL().Fragment
				defnName := strings.TrimPrefix(fragment, "/definitions/")
				_, ok := defns[defnName]
				if !ok {
					return nil, fmt.Errorf("no definition with name %q exists in the openapi spec", defnName)
				}
				itemTypename = snakeToCamel(defnName)
			}
			f = f.Id("[]" + itemTypename)

		case item.Type.Contains("object"):
			itemTypename := typename + "_" + fieldName
			if err := genParamStruct(defns, file, itemTypename, &item.Schema); err != nil {
				return nil, fmt.Errorf("generating child item type: %w", err)
			}
			f = f.Id(itemTypename)
		default:
			// perhaps it's a ref?
			if item.Ref.GetURL() == nil || item.Ref.GetURL().Fragment == "" {
				return nil, fmt.Errorf("unhandled item type %v", item.Type)
			}
			fragment := item.Ref.GetURL().Fragment
			defnName := strings.TrimPrefix(fragment, "/definitions/")
			_, ok := defns[defnName]
			if !ok {
				return nil, fmt.Errorf("no definition with name %q exists in the openapi spec", defnName)
			}
			defnTypeName := snakeToCamel(defnName)
			f = f.Id(defnTypeName)
		}
		jsonTag := item.Name
		if isOptional {
			jsonTag += ",omitempty"
		}
		f = f.Tag(map[string]string{
			"json":  jsonTag,
			"tfsdk": item.Name,
		})
		return f, nil
	}

	var fields []jen.Code
	if body != nil {
		for _, item := range body.Properties.ToOrderedSchemaItems() {
			f, err := toField(item)
			if err != nil {
				return fmt.Errorf("looking at item %q: %w", item.Name, err)
			}
			fields = append(fields, f)
		}
	}
	file.Type().Id(typename).Struct(fields...)
	return nil
}

func genErrRespParamStruct(defns spec.Definitions, file *jen.File, typename string, body *spec.Schema) error {
	toField := func(item spec.OrderSchemaItem) (jen.Code, error) {
		fieldName := snakeToCamel(item.Name)
		f := jen.Id(fieldName)
		isOptional := !slices.Contains(body.Required, item.Name)
		if isOptional {
			f = f.Op("*")
		}
		switch {
		case item.Type.Contains("string"):
			f = f.String()
		case item.Type.Contains("number"):
			f = f.Float64()
		case item.Type.Contains("boolean"):
			f = f.Bool()
		case item.Type.Contains("array"):
			itemTypename := ""
			switch {
			case item.Items.Schema.Type.Contains("string"):
				itemTypename += "string"
			case item.Items.Schema.Type.Contains("number"):
				itemTypename += "float64"
			case item.Items.Schema.Type.Contains("boolean"):
				itemTypename += "bool"
			case item.Items.Schema.Type.Contains("object"):
				itemTypename = typename + "_" + fieldName + "Item"

				if err := genParamStruct(defns, file, itemTypename, item.Items.Schema); err != nil {
					return nil, fmt.Errorf("generating child item type: %w", err)
				}
			case item.Items.Schema.Type.Contains("array"):
				return nil, fmt.Errorf("arrays of array aren't supported")
			}
			f = f.Id("[]" + itemTypename)

		case item.Type.Contains("object"):
			itemTypename := typename + "_" + fieldName
			if err := genParamStruct(defns, file, itemTypename, &item.Schema); err != nil {
				return nil, fmt.Errorf("generating child item type: %w", err)
			}
			f = f.Id(itemTypename)
		default:
			return nil, fmt.Errorf("unhandled item type %v", item.Type)
		}
		jsonTag := item.Name
		if isOptional {
			jsonTag += ",omitempty"
		}
		f = f.Tag(map[string]string{
			"json":  jsonTag,
			"tfsdk": item.Name,
		})
		return f, nil
	}

	fields := []jen.Code{
		jen.Op("*").Id("ErrorResponse"),
	}

	if body != nil {
		for _, item := range body.Properties.ToOrderedSchemaItems() {
			f, err := toField(item)
			if err != nil {
				return fmt.Errorf("looking at item %q: %w", item.Name, err)
			}
			fields = append(fields, f)
		}
	}
	file.Type().Id(typename).Struct(fields...)
	return nil
}

func genClientStruct(
	file *jen.File,
	spec *spec.Swagger,
) error {
	file.Type().Id("Client").Struct(
		jen.Id("httpCl").Op("*").Qual("net/http", "Client"),
		jen.Id("baseURL").Op("*").Qual("net/url", "URL"),
	)

	file.Func().Id("NewClient").Params(
		jen.Id("httpCl").Op("*").Qual("net/http", "Client"),
		jen.Id("baseURL").Op("*").Qual("net/url", "URL"),
	).Parens(jen.Op("*").Id("Client")).BlockFunc(func(g *jen.Group) {
		g.If(jen.Id("baseURL").Op("==").Nil()).Block(
			jen.Id("baseURL").Op("=").Op("&").Qual("net/url", "URL").Values(
				jen.Id("Scheme").Op(":").Lit("https"),
				jen.Id("Host").Op(":").Lit(spec.Host),
				jen.Id("Path").Op(":").Lit(spec.BasePath),
			),
		)
		g.If(jen.Op("!").Qual("strings", "HasSuffix").Call(jen.Id("baseURL").Dot("Path"), jen.Lit("/"))).Block(
			jen.Id("baseURL").Dot("Path").Op("=").Id("baseURL").Dot("Path").Op("+").Lit("/"),
		)

		g.Return(jen.Op("&").Id("Client").Values(
			jen.Id("httpCl").Op(":").Id("httpCl"),
			jen.Id("baseURL").Op(":").Id("baseURL"),
		))
	})
	return nil
}

func genErrorStruct(
	file *jen.File,
	spec *spec.Swagger,
) error {
	file.Type().Id("ErrorResponse").Struct(
		jen.Id("Code").Id("string").Tag(map[string]string{"json": "code"}),
		jen.Id("Message").Id("string").Tag(map[string]string{"json": "message"}),
	)

	file.Func().Params(
		jen.Id("err").Op("*").Id("ErrorResponse"),
	).Id("Error").Params().Parens(jen.String()).BlockFunc(func(g *jen.Group) {
		g.Return(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("error %s: %s"),
				jen.Id("err").Dot("Code"),
				jen.Id("err").Dot("Message"),
			),
		)
	})

	return nil
}

func genClientCall(
	file *jen.File,
	path, verb string,
	clientCallFuncName string,
	pathArgs []spec.Parameter, queryArgs []spec.Parameter,
	reqBodyTypeName string,
	responseTypeNames map[int]string,
) error {

	args := []jen.Code{jen.Id("ctx").Qual("context", "Context")}

	path = strings.TrimPrefix(path, "/")

	pathBuilderArg, pathArgs := pathInterpolator(path, pathArgs)
	for _, pathArg := range pathArgs {
		argName := lowerSnakeToCamel(pathArg.Name)

		argF := jen.Id(argName)
		switch pathArg.Type {
		case "string":
			argF = argF.String()
		case "number":
			argF = argF.Float64()
		default:
			return fmt.Errorf("unhandled pathArg type %v", pathArg.Type)
		}

		args = append(args, argF)
	}
	if reqBodyTypeName != "" {
		args = append(args, jen.Id("req").Id(reqBodyTypeName))
	}
	for _, queryArg := range queryArgs {
		argName := lowerSnakeToCamel(queryArg.Name)
		argF := jen.Id(argName)
		if !queryArg.Required {
			argF = argF.Op("*")
		}
		switch queryArg.Type {
		case "string":
			argF = argF.String()
		case "number":
			argF = argF.Int()
		default:
			return fmt.Errorf("unhandled queryArg type %v", queryArg.Type)
		}
		args = append(args, argF)
	}

	var returnVals []jen.Code
	var returnNames []jen.Code
	codes := maps.Keys(responseTypeNames)
	sort.Ints(codes)
	for _, code := range codes {
		if code >= 400 {
			continue
		}
		returnValName := "res" + strconv.Itoa(code)
		returnValTypeName := responseTypeNames[code]
		returnF := jen.Id(returnValName).Op("*").Id(returnValTypeName)
		returnVals = append(returnVals, returnF)
		returnNames = append(returnNames, jen.Id(returnValName))
	}
	returnVals = append(returnVals, jen.Id("err").Id("error"))
	returnNames = append(returnNames, jen.Id("err"))

	rcvrName := "cl"
	rcvrType := "Client"
	file.Func().Params(
		jen.Id(rcvrName).Op("*").Id(rcvrType),
	).Id(clientCallFuncName).Params(args...).Parens(
		jen.List(returnVals...),
	).BlockFunc(func(g *jen.Group) {
		g.Id("u").Op(":=").Id("cl").Dot("baseURL").Dot("ResolveReference").Call(
			jen.Op("&").Qual("net/url", "URL").Values(jen.Id("Path").Op(":").Add(pathBuilderArg...)),
		)
		if len(queryArgs) > 0 {
			g.Id("q").Op(":=").Id("u").Dot("Query").Call()
			for _, queryArg := range queryArgs {
				argName := lowerSnakeToCamel(queryArg.Name)
				queryVal := jen.Id(argName)
				if !queryArg.Required {
					queryVal = jen.Op("*").Add(queryVal)
				}
				switch queryArg.Type {
				case "string":
					// nothing to do
				case "number":
					queryVal = jen.Qual("strconv", "Itoa").Call(queryVal)
				default:
					panic("should have been handled earlier")
				}
				if !queryArg.Required {
					g.If(jen.Id(argName).Op("!=").Nil()).Block(
						jen.Id("q").Dot("Set").Call(
							jen.Lit(queryArg.Name),
							queryVal,
						),
					)
				} else {
					g.Id("q").Dot("Set").Call(
						jen.Lit(queryArg.Name),
						queryVal,
					)
				}
			}
			g.Id("u").Dot("RawQuery").Op("=").Id("q").Dot("Encode").Call()
		}

		var bodyStmt *jen.Statement
		if reqBodyTypeName == "" {
			bodyStmt = jen.Nil()
		} else {
			bodyStmt = jen.Id("body")
			g.Id("body").Op(":=").Qual("bytes", "NewBuffer").Call(jen.Nil())
			g.If(
				jen.Id("err").Op("=").Qual("encoding/json", "NewEncoder").Call(jen.Id("body")).Dot("Encode").Call(jen.Id("req")),
				jen.Id("err").Op("!=").Nil(),
			).Block(
				jen.Return(returnNames...),
			)
		}

		g.List(jen.Id("r"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
			jen.Id("ctx"),
			jen.Lit(verb),
			jen.Id("u").Dot("String").Call(),
			bodyStmt,
		)
		g.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(returnNames...),
		)
		g.Id("r").Dot("Header").Dot("Set").Call(jen.Lit("Content-Type"), jen.Lit("application/json"))
		g.Id("r").Dot("Header").Dot("Set").Call(jen.Lit("Accept"), jen.Lit("application/json"))

		g.List(jen.Id("res"), jen.Id("err")).Op(":=").Id("cl").Dot("httpCl").Dot("Do").Call(jen.Id("r"))
		g.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(returnNames...),
		)
		g.Defer().Id("res").Dot("Body").Dot("Close").Call()

		g.Switch(jen.Id("res").Dot("StatusCode")).BlockFunc(func(g *jen.Group) {

			for _, code := range codes {
				returnValName := "res" + strconv.Itoa(code)
				returnValTypeName := responseTypeNames[code]
				if code < 400 {

					g.Case(jen.Lit(code)).Block(
						jen.Id(returnValName).Op("=").New(jen.Id(returnValTypeName)),
						jen.Id("err").Op("=").Qual("encoding/json", "NewDecoder").Call(jen.Id("res").Dot("Body")).Dot("Decode").Call(jen.Op("&").Id(returnValName)),
					)
				} else {
					g.Case(jen.Lit(code)).Block(
						jen.Id(returnValName).Op(":=").New(jen.Id(returnValTypeName)),
						jen.Id("err").Op("=").Qual("encoding/json", "NewDecoder").Call(jen.Id("res").Dot("Body")).Dot("Decode").Call(jen.Op("&").Id(returnValName)),
						jen.If(jen.Id("err").Op("==").Nil()).Block(
							jen.Id("err").Op("=").Id(returnValName),
						),
					)
				}
			}
			g.Default().Block(
				jen.Var().Id("errBody").Op("*").Id("ErrorResponse"),
				jen.Id("_").Op("=").Qual("encoding/json", "NewDecoder").Call(jen.Id("res").Dot("Body")).Dot("Decode").Call(jen.Op("&").Id("errBody")),
				jen.If(jen.Id("errBody").Op("!=").Nil()).Block(
					jen.Id("err").Op("=").Id("errBody"),
				).Else().Block(
					jen.Id("err").Op("=").Qual("fmt", "Errorf").Call(jen.Lit("unexpected status code %d"), jen.Id("res").Dot("StatusCode")),
				),
			)
		})
		g.If(jen.Qual("errors", "Is").Call(jen.Id("err"), jen.Qual("io", "EOF"))).Block(jen.Id("err").Op("=").Nil())
		g.Return(returnNames...)
	})

	return nil
}

func pathInterpolator(path string, pathArgs []spec.Parameter) ([]jen.Code, []spec.Parameter) {
	type interpolateArg struct {
		StringLiteral *string
		VariableName  *string
	}
	args := []interpolateArg{
		{StringLiteral: &path},
	}
	for _, pathArg := range pathArgs {
		key := "{" + pathArg.Name + "}"
		argName := lowerSnakeToCamel(pathArg.Name)

		for i, arg := range args {
			if arg.StringLiteral == nil {
				continue
			}
			lit := *arg.StringLiteral
			idx := strings.Index(lit, key)
			if idx < 0 {
				continue
			}
			prePathPart := lit[:idx]

			inserts := []interpolateArg{
				{VariableName: &argName},
			}
			if len(lit) > idx+len(key) {
				postPathPart := lit[idx+len(key):]
				inserts = append(inserts, interpolateArg{StringLiteral: &postPathPart})
			}
			arg.StringLiteral = &prePathPart
			args[i] = arg
			if len(args) > i+1 {
				args = slices.Insert(args, i+1, inserts...)
			} else {
				args = append(args, inserts...)
			}
		}
	}
	var (
		out     []jen.Code
		ordered []spec.Parameter
	)
	for i, arg := range args {
		if i != 0 {
			out = append(out, jen.Op("+"))
		}
		switch {
		case arg.StringLiteral != nil:
			out = append(out, jen.Lit(*arg.StringLiteral))
		case arg.VariableName != nil:
			out = append(out, jen.Id(*arg.VariableName))
			// reorder the path args in the way they appear in the URL
			// so that they can be returned in a sane order for the
			// func argument names
			for _, pathArg := range pathArgs {
				argName := lowerSnakeToCamel(pathArg.Name)
				if argName == *arg.VariableName {
					ordered = append(ordered, pathArg)
				}
			}

		}
	}
	return out, ordered
}
