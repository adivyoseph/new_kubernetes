/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generators

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/kubernetes/cmd/libs/go2idl/client-gen/generators/normalization"
	"k8s.io/kubernetes/pkg/api/unversioned"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	pluralExceptions := map[string]string{
		"Endpoints": "Endpoints",
	}
	return namer.NameSystems{
		"public":             namer.NewPublicNamer(0),
		"private":            namer.NewPrivateNamer(0),
		"raw":                namer.NewRawNamer("", nil),
		"publicPlural":       namer.NewPublicPluralNamer(pluralExceptions),
		"allLowercasePlural": namer.NewAllLowercasePluralNamer(pluralExceptions),
		"lowercaseSingular":  &lowercaseSingularNamer{},
	}
}

type lowercaseSingularNamer struct{}

func (n *lowercaseSingularNamer) Name(t *types.Type) string {
	return strings.ToLower(t.Name.Name)
}

// DefaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
func DefaultNameSystem() string {
	return "public"
}

func generatedBy(cmdArgs string) string {
	if len(cmdArgs) != 0 {
		return fmt.Sprintf("\n// This file was automatically generated by lister-gen with arguments: %s\n\n", cmdArgs)
	}
	return fmt.Sprintf("\n// This file was automatically generated by lister-gen with the default arguments.\n\n")
}

// Packages makes the client package definition.
func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		glog.Fatalf("Failed loading boilerplate: %v", err)
	}

	var cmdArgs string
	pflag.VisitAll(func(f *pflag.Flag) {
		if !f.Changed || f.Name == "verify-only" {
			return
		}
		cmdArgs += fmt.Sprintf("--%s=%s ", f.Name, f.Value)
	})

	generatedBy := generatedBy(cmdArgs)
	boilerplate = append(boilerplate, []byte(generatedBy)...)

	var groupVersions []unversioned.GroupVersion
	gvToTypes := map[unversioned.GroupVersion][]*types.Type{}
	for _, inputDir := range arguments.InputDirs {
		p := context.Universe.Package(inputDir)
		gv := unversioned.GroupVersion{Group: p.Name}
		groupVersions = append(groupVersions, gv)
		for _, t := range p.Types {
			// filter out types which dont have genclient=true.
			if extractBoolTagOrDie("genclient", t.SecondClosestCommentLines) == false {
				continue
			}
			if _, found := gvToTypes[gv]; !found {
				gvToTypes[gv] = []*types.Type{}
			}
			gvToTypes[gv] = append(gvToTypes[gv], t)
		}
	}

	var packageList generator.Packages

	orderer := namer.Orderer{Namer: namer.NewPrivateNamer(0)}
	for i := range groupVersions {
		gv := groupVersions[i]

		normalizedGroupName := normalization.BeforeFirstDot(gv.Group)
		if normalizedGroupName == "api" {
			normalizedGroupName = "core"
		}

		packageList = append(packageList, &generator.DefaultPackage{
			PackageName: filepath.Base(arguments.OutputPackagePath),
			PackagePath: arguments.OutputPackagePath,
			HeaderText:  boilerplate,
			GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
				return []generator.Generator{
					&genListers{
						DefaultGen: generator.DefaultGen{
							OptionalName: arguments.OutputFileBaseName + "." + normalizedGroupName,
						},
						outputPackage: arguments.OutputPackagePath,
						groupVersion:  gv,
						types:         orderer.OrderTypes(gvToTypes[gv]),
						imports:       generator.NewImportTracker(),
					},
				}
			},
			FilterFunc: func(c *generator.Context, t *types.Type) bool {
				// piggy-back on types that are tagged for client-gen
				return extractBoolTagOrDie("genclient", t.SecondClosestCommentLines) == true
			},
		})
	}

	return packageList
}

// genListers produces a file of listers for a group.
type genListers struct {
	generator.DefaultGen
	outputPackage string
	groupVersion  unversioned.GroupVersion
	// types in this group
	types   []*types.Type
	imports namer.ImportTracker
	// if true, we need to import k8s.io/kubernetes/pkg/api
	hasNonNamespaced bool
}

var _ generator.Generator = &genListers{}

func (g *genListers) Filter(c *generator.Context, t *types.Type) bool {
	return t == g.types[0]
}

func (g *genListers) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

func (g *genListers) Imports(c *generator.Context) (imports []string) {
	if g.hasNonNamespaced {
		objectMeta := c.Universe.Type(types.Name{Package: "k8s.io/kubernetes/pkg/api", Name: "ObjectMeta"})
		g.imports.AddType(objectMeta)
	}
	imports = append(imports, g.imports.ImportLines()...)
	imports = append(imports, "k8s.io/kubernetes/pkg/api/errors")
	imports = append(imports, "k8s.io/kubernetes/pkg/labels")
	// for Indexer
	imports = append(imports, "k8s.io/kubernetes/pkg/client/cache")
	return
}

func (g *genListers) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	groupName := g.groupVersion.Group
	if g.groupVersion.Group == "core" {
		groupName = "api"
	}

	for _, t := range g.types {
		glog.V(5).Infof("processing type %v", t)
		m := map[string]interface{}{
			"group": groupName,
			"type":  t,
		}

		namespaced := !extractBoolTagOrDie("nonNamespaced", t.SecondClosestCommentLines)
		if namespaced {
			sw.Do(typeListerInterface, m)
		} else {
			g.hasNonNamespaced = true
			sw.Do(typeListerInterface_NonNamespaced, m)
		}

		sw.Do(typeListerStruct, m)
		sw.Do(typeListerConstructor, m)
		sw.Do(typeLister_List, m)

		if namespaced {
			sw.Do(typeLister_NamespaceLister, m)
			sw.Do(namespaceListerInterface, m)
			sw.Do(namespaceListerStruct, m)
			sw.Do(namespaceLister_List, m)
			sw.Do(namespaceLister_Get, m)
		} else {
			sw.Do(typeLister_NonNamespacedGet, m)
		}
	}

	return sw.Error()
}

var typeListerInterface = `
// $.type|public$Lister helps list $.type|publicPlural$.
type $.type|public$Lister interface {
	// List lists all $.type|publicPlural$ in the indexer.
	List(selector labels.Selector) (ret []*$.type|raw$, err error)
	// $.type|publicPlural$ returns an object that can list and get $.type|publicPlural$.
	$.type|publicPlural$(namespace string) $.type|public$NamespaceLister
}
`

var typeListerInterface_NonNamespaced = `
// $.type|public$Lister helps list $.type|publicPlural$.
type $.type|public$Lister interface {
	// List lists all $.type|publicPlural$ in the indexer.
	List(selector labels.Selector) (ret []*$.type|raw$, err error)
	// Get retrieves the $.type|public$ from the index for a given name.
	Get(name string) (*$.type|raw$, error)
}
`

var typeListerStruct = `
// $.type|private$Lister implements the $.type|public$Lister interface.
type $.type|private$Lister struct {
	indexer cache.Indexer
}
`

var typeListerConstructor = `
// New$.type|public$Lister returns a new $.type|public$Lister.
func New$.type|public$Lister(indexer cache.Indexer) $.type|public$Lister {
	return &$.type|private$Lister{indexer: indexer}
}
`

var typeLister_List = `
// List lists all $.type|publicPlural$ in the indexer.
func (s *$.type|private$Lister) List(selector labels.Selector) (ret []*$.type|raw$, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*$.type|raw$))
	})
	return ret, err
}
`

var typeLister_NamespaceLister = `
// $.type|publicPlural$ returns an object that can list and get $.type|publicPlural$.
func (s *$.type|private$Lister) $.type|publicPlural$(namespace string) $.type|public$NamespaceLister {
	return $.type|private$NamespaceLister{indexer: s.indexer, namespace: namespace}
}
`

var typeLister_NonNamespacedGet = `
// Get retrieves the $.type|public$ from the index for a given name.
func (s *$.type|private$Lister) Get(name string) (*$.type|raw$, error) {
  key := &$.type|raw${ObjectMeta: api.ObjectMeta{Name: name}}
  obj, exists, err := s.indexer.Get(key)
  if err != nil {
    return nil, err
  }
  if !exists {
    return nil, errors.NewNotFound($.group$.Resource("$.type|lowercaseSingular$"), name)
  }
  return obj.(*$.type|raw$), nil
}
`

var namespaceListerInterface = `
// $.type|public$NamespaceLister helps list and get $.type|publicPlural$.
type $.type|public$NamespaceLister interface {
	// List lists all $.type|publicPlural$ in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*$.type|raw$, err error)
	// Get retrieves the $.type|public$ from the indexer for a given namespace and name.
	Get(name string) (*$.type|raw$, error)
}
`

var namespaceListerStruct = `
// $.type|private$NamespaceLister implements the $.type|public$NamespaceLister
// interface.
type $.type|private$NamespaceLister struct {
	indexer cache.Indexer
	namespace string
}
`

var namespaceLister_List = `
// List lists all $.type|publicPlural$ in the indexer for a given namespace.
func (s $.type|private$NamespaceLister) List(selector labels.Selector) (ret []*$.type|raw$, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*$.type|raw$))
	})
	return ret, err
}
`

var namespaceLister_Get = `
// Get retrieves the $.type|public$ from the indexer for a given namespace and name.
func (s $.type|private$NamespaceLister) Get(name string) (*$.type|raw$, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound($.group$.Resource("$.type|lowercaseSingular$"), name)
	}
	return obj.(*$.type|raw$), nil
}
`
