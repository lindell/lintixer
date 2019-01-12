package fixer

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// Fixer does contain information what should be fixed in the code and how
type Fixer struct {
	logger Logger
	fixers []NodeFixer
}

// New creates a new Fixer
func New(opts ...Option) *Fixer {
	f := &Fixer{
		logger: &nopLogger{},
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// Option is an option for the fixer
type Option func(*Fixer)

// NodeFixer looks at an ast node  and fixes it if something is wrong
type NodeFixer func(node ast.Node) (changeMade bool)

// Logger is the logger used
type Logger interface {
	Info(string)
}

// nopLogger is used when no other logger is specified
type nopLogger struct{}

func (n *nopLogger) Info(string) {}

// WithLogger sets a logger of the
func WithLogger(logger Logger) Option {
	return func(f *Fixer) { f.logger = logger }
}

// WithNodeFixers adds one Node
func WithNodeFixers(nf NodeFixer) Option {
	return func(f *Fixer) { f.fixers = append(f.fixers, nf) }
}

// Fix any files within a path
func (f *Fixer) Fix(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("file or directory does not exist")
		}
		return err
	}

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return f.fixFile(path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (f *Fixer) fixFile(path string) error {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	anyChanged := false
	for _, nodeFixer := range f.fixers {
		changed := nodeFixer(astFile)
		if changed {
			anyChanged = true
		}
	}

	if !anyChanged {
		return nil
	}

	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open to file: %s\n%s", path, err.Error())
	}
	defer file.Close()

	err = format.Node(file, fset, astFile)
	if err != nil {
		return fmt.Errorf("could not open to file: %s\n%s", path, err.Error())
	}

	f.logger.Info(fmt.Sprintf("changed file: %s", path))

	return nil
}
