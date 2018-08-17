/*
Copyright 2017 The Kubernetes Authors.

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

package plugins

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	// ErrIncompletePlugin indicates plugin is incomplete.
	ErrIncompletePlugin = fmt.Errorf("incomplete plugin descriptor: name, shortDesc and command fields are required")
	// ErrInvalidPluginName indicates plugin name is invalid.
	ErrInvalidPluginName = fmt.Errorf("plugin name can't contain spaces")
	// ErrIncompleteFlag indicates flag is incomplete.
	ErrIncompleteFlag = fmt.Errorf("incomplete flag descriptor: name and desc fields are required")
	// ErrInvalidFlagName indicates flag name is invalid.
	ErrInvalidFlagName = fmt.Errorf("flag name can't contain spaces")
	// ErrInvalidFlagShorthand indicates flag shorthand is invalid.
	ErrInvalidFlagShorthand = fmt.Errorf("flag shorthand must be only one letter")
	// ErrInvalidFlagType indicates flag type is invalid.
	ErrInvalidFlagType = fmt.Errorf("invalid flag type, must be one of the following: string(default), bool")
	// ErrInvalidFlagDefaultValue indicates flag default value is invalid.
	ErrInvalidFlagDefaultValue = fmt.Errorf("invalid flag default value: must be parseable as the flag type")
)

// Plugin is the representation of a CLI extension (plugin).
type Plugin struct {
	Description
	Source
	Context RunningContext `json:"-"`
}

// Description holds everything needed to register a
// plugin as a command. Usually comes from a descriptor file.
type Description struct {
	Name      string  `json:"name"`
	ShortDesc string  `json:"shortDesc"`
	LongDesc  string  `json:"longDesc,omitempty"`
	Example   string  `json:"example,omitempty"`
	Command   string  `json:"command"`
	Flags     []Flag  `json:"flags,omitempty"`
	Tree      Plugins `json:"tree,omitempty"`
}

// Source holds the location of a given plugin in the filesystem.
type Source struct {
	Dir            string `json:"-"`
	DescriptorName string `json:"-"`
}

// Validate validates plugin data.
func (p Plugin) Validate() error {
	if len(p.Name) == 0 || len(p.ShortDesc) == 0 || (len(p.Command) == 0 && len(p.Tree) == 0) {
		return ErrIncompletePlugin
	}
	if strings.Contains(p.Name, " ") {
		return ErrInvalidPluginName
	}
	for _, flag := range p.Flags {
		if err := flag.Validate(); err != nil {
			return err
		}
	}
	for _, child := range p.Tree {
		if err := child.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// IsValid returns true if plugin data is valid.
func (p Plugin) IsValid() bool {
	return p.Validate() == nil
}

// Plugins is a list of plugins.
type Plugins []*Plugin

// Flag describes a single flag supported by a given plugin.
type Flag struct {
	Name      string   `json:"name"`
	Type      FlagType `json:"type,omitempty"`
	Shorthand string   `json:"shorthand,omitempty"`
	Desc      string   `json:"desc"`
	DefValue  string   `json:"defValue,omitempty"`
}

// FlagType specifies the datatype of the flag.
type FlagType string

const (
	// BoolFlagType represents boolean valued flags
	// The presence of the flag indicates the value is true
	// For example: --debug is equivalent to --debug=true
	BoolFlagType FlagType = "bool"
	// StringFlagType represents string valued flags.
	// This is the default flag type when no value is specified.
	StringFlagType FlagType = "string"
)

// Validate validates flag data.
func (f *Flag) Validate() error {
	// the set of supported flag types
	var flagTypes = map[FlagType]struct{}{
		BoolFlagType:   {},
		StringFlagType: {},
	}

	if len(f.Name) == 0 || len(f.Desc) == 0 {
		return ErrIncompleteFlag
	}
	if strings.Contains(f.Name, " ") {
		return ErrInvalidFlagName
	}

	// Default to a string flag when not specified to preserve backwards compatibility.
	if f.Type == "" {
		f.Type = StringFlagType
	}

	if _, ok := flagTypes[f.Type]; !ok {
		return ErrInvalidFlagType
	}

	if f.Type == BoolFlagType {
		if f.DefValue == "" {
			f.DefValue = "false"
		}

		if _, err := strconv.ParseBool(f.DefValue); err != nil {
			return ErrInvalidFlagDefaultValue
		}
	}

	return f.ValidateShorthand()
}

// ValidateShorthand validates flag shorthand data.
func (f *Flag) ValidateShorthand() error {
	length := len(f.Shorthand)
	if length == 0 || (length == 1 && unicode.IsLetter(rune(f.Shorthand[0]))) {
		return nil
	}
	return ErrInvalidFlagShorthand
}

// Shorthanded returns true if flag shorthand data is valid.
func (f *Flag) Shorthanded() bool {
	return f.ValidateShorthand() == nil
}
