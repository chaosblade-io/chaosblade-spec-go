/*
 * Copyright 1999-2019 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package spec

import (
	"fmt"
	"strings"
)

// ExpModelCommandSpec defines the command interface for the experimental plugin
type ExpModelCommandSpec interface {
	// Name returns the target name
	Name() string

	// Scope returns the experiment scope
	Scope() string

	// ShortDesc returns short description for the command
	ShortDesc() string

	// LongDesc returns full description for the command
	LongDesc() string

	// Actions returns the list of actions supported by the command
	Actions() []ExpActionCommandSpec

	// Flags returns the command flags
	Flags() []ExpFlagSpec

	// SetFlags
	SetFlags(flags []ExpFlagSpec)
}

// ExpActionCommandSpec defines the action command interface for the experimental plugin
type ExpActionCommandSpec interface {
	// Name returns the action name
	Name() string

	// Aliases returns command alias names
	Aliases() []string

	// ShortDesc returns short description for the action
	ShortDesc() string

	// LongDesc returns full description for the action
	LongDesc() string

	// SetLongDesc
	SetLongDesc(longDesc string)

	// Matchers returns the list of matchers supported by the action
	Matchers() []ExpFlagSpec

	// Flags returns the list of flags supported by the action
	Flags() []ExpFlagSpec

	//Example returns command example
	Example() string

	//Example returns command example
	SetExample(example string)

	// ExpExecutor returns the action command ExpExecutor
	Executor() Executor

	// SetExecutor
	SetExecutor(executor Executor)

	// Programs executed
	Programs() []string

	// Scenario categories
	Categories() []string

	// SetCategories
	SetCategories(categories []string)
}

type ExpFlagSpec interface {
	// FlagName returns the flag FlagName
	FlagName() string
	// FlagDesc returns the flag description
	FlagDesc() string
	// FlagNoArgs returns true if the flag is bool type
	FlagNoArgs() bool
	// FlagRequired returns true if the flag is necessary when creating experiment
	FlagRequired() bool
	// FlagRequiredWhenDestroyed returns true if the flag is necessary when destroying experiment
	FlagRequiredWhenDestroyed() bool
	// 	FlagDefault return the flag Defaule
	FlagDefault() string
}

// ExpFlag defines the action flag
type ExpFlag struct {
	// Name returns the flag FlagName
	Name string `yaml:"name"`

	// Desc returns the flag description
	Desc string `yaml:"desc"`

	// NoArgs means no arguments
	NoArgs bool `yaml:"noArgs"`

	// Required means necessary or not
	Required bool `yaml:"required"`
	// RequiredWhenDestroyed is true if the flag is necessary when destroying experiment
	RequiredWhenDestroyed bool `yaml:"requiredWhenDestroyed"`

	// default value
	Default string `yaml:"default,omitempty"`
}

func (f *ExpFlag) FlagName() string {
	return f.Name
}

func (f *ExpFlag) FlagDesc() string {
	return f.Desc
}

func (f *ExpFlag) FlagNoArgs() bool {
	return f.NoArgs
}

func (f *ExpFlag) FlagRequired() bool {
	return f.Required
}

func (f *ExpFlag) FlagRequiredWhenDestroyed() bool {
	return f.RequiredWhenDestroyed
}

func (f *ExpFlag) FlagDefault() string {
	return f.Default
}

// BaseExpModelCommandSpec defines the common struct of the implementation of ExpModelCommandSpec
type BaseExpModelCommandSpec struct {
	ExpScope   string
	ExpActions []ExpActionCommandSpec
	ExpFlags   []ExpFlagSpec
}

// Scope default value is "" means localhost
func (b *BaseExpModelCommandSpec) Scope() string {
	return ""
}

func (b *BaseExpModelCommandSpec) Actions() []ExpActionCommandSpec {
	return b.ExpActions
}

func (b *BaseExpModelCommandSpec) Flags() []ExpFlagSpec {
	return b.ExpFlags
}

func (b *BaseExpModelCommandSpec) SetFlags(flags []ExpFlagSpec) {
	b.ExpFlags = flags
}

// BaseExpActionCommandSpec defines the common struct of the implementation of ExpActionCommandSpec
type BaseExpActionCommandSpec struct {
	ActionMatchers   []ExpFlagSpec
	ActionFlags      []ExpFlagSpec
	ActionExecutor   Executor
	ActionLongDesc   string
	ActionExample    string
	ActionPrograms   []string
	ActionCategories []string
}

func (b *BaseExpActionCommandSpec) Matchers() []ExpFlagSpec {
	return b.ActionMatchers
}

func (b *BaseExpActionCommandSpec) Flags() []ExpFlagSpec {
	return b.ActionFlags
}

func (b *BaseExpActionCommandSpec) Executor() Executor {
	return b.ActionExecutor
}

func (b *BaseExpActionCommandSpec) SetExecutor(executor Executor) {
	b.ActionExecutor = executor
}

func (b *BaseExpActionCommandSpec) SetLongDesc(longDesc string) {
	b.ActionLongDesc = longDesc
}

func (b *BaseExpActionCommandSpec) SetExample(example string) {
	b.ActionExample = example
}

func (b *BaseExpActionCommandSpec) Example() string {
	return b.ActionExample
}

func (b *BaseExpActionCommandSpec) Programs() []string {
	return b.ActionPrograms
}

func (b *BaseExpActionCommandSpec) Categories() []string {
	return b.ActionCategories
}

func (b *BaseExpActionCommandSpec) SetCategories(categories []string) {
	b.ActionCategories = categories
}

// ActionModel for yaml file
type ActionModel struct {
	ActionName       string    `yaml:"action"`
	ActionAliases    []string  `yaml:"aliases,flow,omitempty"`
	ActionShortDesc  string    `yaml:"shortDesc"`
	ActionLongDesc   string    `yaml:"longDesc"`
	ActionMatchers   []ExpFlag `yaml:"matchers,omitempty"`
	ActionFlags      []ExpFlag `yaml:"flags,omitempty"`
	ActionExample    string    `yaml:"example"`
	executor         Executor
	ActionPrograms   []string `yaml:"programs,omitempty"`
	ActionCategories []string `yaml:"categories,omitempty"`
}

func (am *ActionModel) Programs() []string {
	return am.ActionPrograms
}

func (am *ActionModel) SetExample(example string) {
	am.ActionExample = example
}

func (am *ActionModel) Example() string {
	return am.ActionExample
}

func (am *ActionModel) SetExecutor(executor Executor) {
	am.executor = executor
}

func (am *ActionModel) Executor() Executor {
	return am.executor
}

func (am *ActionModel) Name() string {
	return am.ActionName
}

func (am *ActionModel) Aliases() []string {
	return am.ActionAliases
}

func (am *ActionModel) ShortDesc() string {
	return am.ActionShortDesc
}

func (am *ActionModel) SetLongDesc(longDesc string) {
	am.ActionLongDesc = longDesc
}

func (am *ActionModel) LongDesc() string {
	return am.ActionLongDesc
}

func (am *ActionModel) Matchers() []ExpFlagSpec {
	flags := make([]ExpFlagSpec, 0)
	for idx := range am.ActionMatchers {
		flags = append(flags, &am.ActionMatchers[idx])
	}
	return flags
}

func (am *ActionModel) Flags() []ExpFlagSpec {
	flags := make([]ExpFlagSpec, 0)
	for idx := range am.ActionFlags {
		flags = append(flags, &am.ActionFlags[idx])
	}
	return flags
}

func (am *ActionModel) Categories() []string {
	return am.ActionCategories
}

func (am *ActionModel) SetCategories(categories []string) {
	am.ActionCategories = categories
}

type ExpPrepareModel struct {
	PrepareType     string    `yaml:"type"`
	PrepareFlags    []ExpFlag `yaml:"flags"`
	PrepareRequired bool      `yaml:"required"`
}

type ExpCommandModel struct {
	ExpName         string          `yaml:"target"`
	ExpShortDesc    string          `yaml:"shortDesc"`
	ExpLongDesc     string          `yaml:"longDesc"`
	ExpActions      []ActionModel   `yaml:"actions"`
	ExpExecutor     Executor        `yaml:"-"`
	ExpFlags        []ExpFlag       `yaml:"flags,omitempty"`
	ExpScope        string          `yaml:"scope"`
	ExpPrepareModel ExpPrepareModel `yaml:"prepare,omitempty"`
	ExpSubTargets   []string        `yaml:"subTargets,flow,omitempty"`
}

func (ecm *ExpCommandModel) Scope() string {
	return ecm.ExpScope
}

func (ecm *ExpCommandModel) Name() string {
	return ecm.ExpName
}

func (ecm *ExpCommandModel) ShortDesc() string {
	return ecm.ExpShortDesc
}

func (ecm *ExpCommandModel) LongDesc() string {
	return ecm.ExpLongDesc
}

func (ecm *ExpCommandModel) Actions() []ExpActionCommandSpec {
	specs := make([]ExpActionCommandSpec, 0)
	for idx := range ecm.ExpActions {
		if ecm.ExpExecutor != nil {
			ecm.ExpActions[idx].executor = ecm.ExpExecutor
		}
		specs = append(specs, &ecm.ExpActions[idx])
	}
	return specs
}

func (ecm *ExpCommandModel) Flags() []ExpFlagSpec {
	flags := make([]ExpFlagSpec, 0)
	for idx := range ecm.ExpFlags {
		flags = append(flags, &ecm.ExpFlags[idx])
	}
	return flags
}

func (ecm *ExpCommandModel) SetFlags(flags []ExpFlagSpec) {
	expFlags := make([]ExpFlag, 0)
	for idx := range flags {
		expFlags = append(expFlags, *flags[idx].(*ExpFlag))
	}
	ecm.ExpFlags = expFlags
}

type Models struct {
	Version string            `yaml:"version"`
	Kind    string            `yaml:"kind"`
	Models  []ExpCommandModel `yaml:"items"`
}

type Empty struct{}

// ConvertExpMatchersToString returns the flag arguments for cli
func ConvertExpMatchersToString(expModel *ExpModel, createExcludeKeyFunc func() map[string]Empty) string {
	matchers := ""
	excludeKeys := createExcludeKeyFunc()
	flags := expModel.ActionFlags
	if flags != nil && len(flags) > 0 {
		for name, value := range flags {
			// exclude unsupported key in blade
			if _, ok := excludeKeys[name]; ok {
				continue
			}
			if value == "" {
				continue
			}
			if strings.Contains(value, " ") {
				value = strings.ReplaceAll(value, " ", "@@##")
			}
			matchers = fmt.Sprintf(`%s --%s=%s`, matchers, name, value)
		}
	}
	return matchers
}

// ConvertCommandsToExpModel returns the ExpModel by action, target and flags
func ConvertCommandsToExpModel(action, target, rules string) *ExpModel {
	model := &ExpModel{
		Target:      target,
		ActionName:  action,
		ActionFlags: make(map[string]string, 0),
	}
	flags := strings.Split(rules, " ")
	for _, flag := range flags {
		keyAndValue := strings.SplitN(flag, "=", 2)
		if len(keyAndValue) != 2 {
			continue
		}
		key := keyAndValue[0][2:]
		model.ActionFlags[key] = strings.ReplaceAll(keyAndValue[1], "@@##", " ")
	}
	return model
}

// AddFlagsToModelSpec
func AddFlagsToModelSpec(flagsFunc func() []ExpFlagSpec, expSpecs ...ExpModelCommandSpec) {
	flagSpecs := flagsFunc()
	for _, expSpec := range expSpecs {
		flags := expSpec.Flags()
		if flags == nil {
			flags = make([]ExpFlagSpec, 0)
		}
		flags = append(flags, flagSpecs...)
		expSpec.SetFlags(flags)
	}
}

// AddExecutorToModelSpec
func AddExecutorToModelSpec(executor Executor, expSpecs ...ExpModelCommandSpec) {
	for _, expSpec := range expSpecs {
		actions := expSpec.Actions()
		if actions == nil {
			continue
		}
		for _, action := range actions {
			action.SetExecutor(executor)
		}
	}
}
