package generator

import (
	"go/build"
	"go/types"
	"path"
	"path/filepath"

	"github.com/go-courier/codegen"
	"github.com/go-courier/loaderx"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/loader"

	"github.com/go-courier/status_error"
)

func NewStatusErrorGenerator(program *loader.Program, rootPkgInfo *loader.PackageInfo) *StatusErrorGenerator {
	return &StatusErrorGenerator{
		pkgInfo:      rootPkgInfo,
		scanner:      NewStatusErrorScanner(program),
		statusErrors: map[string]*StatusError{},
	}
}

type StatusErrorGenerator struct {
	pkgInfo      *loader.PackageInfo
	scanner      *StatusErrorScanner
	statusErrors map[string]*StatusError
}

func (g *StatusErrorGenerator) Scan(names ...string) {
	pkgInfo := loaderx.NewPackageInfo(g.pkgInfo)

	for _, name := range names {
		typeName := pkgInfo.TypeName(name)
		g.statusErrors[name] = &StatusError{
			TypeName: typeName,
			Errors:   g.scanner.StatusError(typeName),
		}
	}
}

func (g *StatusErrorGenerator) Output(cwd string) {
	for _, statusErr := range g.statusErrors {
		p, _ := build.Import(statusErr.TypeName.Pkg().Path(), "", build.FindOnly)
		dir, _ := filepath.Rel(cwd, p.Dir)
		filename := codegen.GeneratedFileSuffix(path.Join(dir, codegen.LowerSnakeCase(statusErr.Name())+".go"))

		file := codegen.NewFile(statusErr.TypeName.Pkg().Name(), filename)
		statusErr.WriteToFile(file)

		if _, err := file.WriteFile(); err != nil {
			logrus.Printf("%s generated", file)
		}
	}
}

type StatusError struct {
	TypeName *types.TypeName
	Errors   []*status_error.StatusErr
}

func (s *StatusError) Name() string {
	return s.TypeName.Name()
}

func (s *StatusError) WriteToFile(file *codegen.File) {
	s.WriteMethodImplements(file)

	s.WriteMethodStatusErrAndError(file)
	s.WriteMethodStatus(file)
	s.WriteMethodCode(file)

	s.WriteMethodKey(file)
	s.WriteMethodMsg(file)
	s.WriteMethodCanBeTalkError(file)
}

func (s *StatusError) WriteMethodImplements(file *codegen.File) {
	tpe := codegen.Type(file.Use("github.com/go-courier/status_error", "StatusError"))

	file.WriteBlock(
		file.Expr("var _ ? = (*?)(nil)", codegen.Interface(tpe), codegen.Type(s.Name())),
	)
}

func (s *StatusError) WriteMethodStatusErrAndError(file *codegen.File) {
	tpe := codegen.Type(file.Use("github.com/go-courier/status_error", "StatusErr"))

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("StatusErr").
			Return(codegen.Var(codegen.Star(tpe))).Do(
			file.Expr(`return &?{
Key: v.Key(),
Code: v.Code(),
Msg: v.Msg(),
CanBeTalkError: v.CanBeTalkError(),
}`, tpe)),
	)

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("Error").
			Return(codegen.Var(codegen.String)).Do(
			file.Expr(`return v.StatusErr().Error()`)),
	)
}

func (s *StatusError) WriteMethodStatus(file *codegen.File) {
	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("Status").
			Return(codegen.Var(codegen.Int)).Do(
			file.Expr(`return ?(int(v))`, codegen.Id(file.Use("github.com/go-courier/status_error", "GetStatus"))),
		),
	)
}

func (s *StatusError) WriteMethodCode(file *codegen.File) {
	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("Code").
			Return(codegen.Var(codegen.Int)).Do(
			file.Expr(`if withServiceCode, ok := (interface{})(v).(?); ok {
	return withServiceCode.ServiceCode() + int(v)
}
return int(v)
`, codegen.Id(file.Use("github.com/go-courier/status_error", "StatusErrorWithServiceCode"))),
		),
	)
}

func (s *StatusError) WriteMethodKey(file *codegen.File) {
	clauses := make([]*codegen.SnippetClause, 0)

	for _, statusErr := range s.Errors {
		clauses = append(clauses, codegen.Clause(codegen.Id(statusErr.Key)).Do(
			codegen.Return(
				file.Val(statusErr.Key),
			),
		))
	}

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("Key").
			Return(codegen.Var(codegen.String)).Do(
			codegen.Switch(codegen.Id("v")).When(
				clauses...,
			),
			codegen.Return(file.Val("UNKNOWN")),
		),
	)
}

func (s *StatusError) WriteMethodMsg(file *codegen.File) {
	clauses := make([]*codegen.SnippetClause, 0)

	for _, statusErr := range s.Errors {
		clauses = append(clauses, codegen.Clause(codegen.Id(statusErr.Key)).Do(
			codegen.Return(
				file.Val(statusErr.Msg),
			),
		))
	}

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("Msg").
			Return(codegen.Var(codegen.String)).Do(
			codegen.Switch(codegen.Id("v")).When(
				clauses...,
			),
			codegen.Return(file.Val("-")),
		),
	)
}

func (s *StatusError) WriteMethodCanBeTalkError(file *codegen.File) {
	clauses := make([]*codegen.SnippetClause, 0)

	for _, statusErr := range s.Errors {
		clauses = append(clauses, codegen.Clause(codegen.Id(statusErr.Key)).Do(
			codegen.Return(
				file.Val(statusErr.CanBeTalkError),
			),
		))
	}

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("CanBeTalkError").
			Return(codegen.Var(codegen.Bool)).Do(
			codegen.Switch(codegen.Id("v")).When(
				clauses...,
			),
			codegen.Return(file.Val(false)),
		),
	)
}
