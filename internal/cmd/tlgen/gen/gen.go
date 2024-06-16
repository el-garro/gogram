package gen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"

	"github.com/amarnathcjd/gogram/internal/cmd/tlgen/tlparser"
)

// некоторые названия (id, api, url etc.) нужно или капсом или никак (например нельзя написать Id, или
// Api, только ID API URL)
var capitalizePatterns = []string{
	"id",
	"api",
	"url",
	"p2p",
	"sha",
	"srp",
}

type Generator struct {
	schema *internalSchema
	outdir string

	PackageName   string
	PackageHeader string
}

func NewGenerator(tlschema *tlparser.Schema, licenseHeader, outdir string) (*Generator, error) {
	internalSchema, err := createInternalSchema(tlschema)
	if err != nil {
		return nil, errors.Wrap(err, "analyzing schema")
	}

	return &Generator{
		schema:        internalSchema,
		outdir:        outdir,
		PackageName:   "telegram",
		PackageHeader: licenseHeader + "\nCode generated by tlgen; DO NOT EDIT.",
	}, nil
}

func (g *Generator) Generate(d bool) error {
	err := g.generateFile(g.generateEnumDefinitions, filepath.Join(g.outdir, "enums_gen.go"), d)
	if err != nil {
		return fmt.Errorf("generate enums: %w", err)
	}

	err = g.generateFile(g.generateSpecificStructs, filepath.Join(g.outdir, "types_gen.go"), d)
	if err != nil {
		return fmt.Errorf("generate types: %w", err)
	}

	err = g.generateFile(g.generateInterfaces, filepath.Join(g.outdir, "interfaces_gen.go"), d)
	if err != nil {
		return fmt.Errorf("generate interfaces: %w", err)
	}

	err = g.generateFile(g.generateMethods, filepath.Join(g.outdir, "methods_gen.go"), d)
	if err != nil {
		return fmt.Errorf("generate methods: %w", err)
	}

	err = g.generateFile(g.generateInit, filepath.Join(g.outdir, "init_gen.go"), d)
	if err != nil {
		return fmt.Errorf("generate init: %w", err)
	}

	return nil
}

func (*Generator) generateFile(f func(file *jen.File, d bool), filename string, genDocs bool) error {
	file := jen.NewFile("telegram")
	file.HeaderComment("Code generated by TLParser; DO NOT EDIT. (c) @amarnathcjd")
	f(file, genDocs)

	buf := bytes.NewBuffer([]byte{})
	if err := file.Render(buf); err != nil {
		return err
	}

	return os.WriteFile(filename, buf.Bytes(), 0644)
}
