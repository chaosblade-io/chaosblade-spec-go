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

package util

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

// CreateYamlFile converts the spec.Models to spec file
func CreateYamlFile(models *spec.Models, specFile string) error {
	file, err := os.OpenFile(specFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	return MarshalModelSpec(models, file)
}

// MarshalModelSpec marshals the spec.Models to bytes and output to writer
func MarshalModelSpec(models *spec.Models, writer io.Writer) error {
	bytes, err := yaml.Marshal(models)
	if err != nil {
		return err
	}
	writer.Write(bytes)
	return nil
}

// ParseSpecsToModel parses the yaml file to spec.Models and set the executor to the spec.Models
func ParseSpecsToModel(file string, executor spec.Executor) (*spec.Models, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	models := &spec.Models{}
	err = yaml.Unmarshal(bytes, models)
	if err != nil {
		return nil, err
	}
	for idx := range models.Models {
		models.Models[idx].ExpExecutor = executor
	}
	return models, nil
}

// ConvertSpecToModels converts the spec.ExpModelCommandSpec to spec.Models
func ConvertSpecToModels(commandSpec spec.ExpModelCommandSpec, prepare spec.ExpPrepareModel, scope string) *spec.Models {
	models := &spec.Models{
		Version: "v1",
		Kind:    "plugin",
		Models:  make([]spec.ExpCommandModel, 0),
	}

	model := spec.ExpCommandModel{
		ExpName:         commandSpec.Name(),
		ExpShortDesc:    commandSpec.ShortDesc(),
		ExpLongDesc:     commandSpec.LongDesc(),
		ExpActions:      make([]spec.ActionModel, 0),
		ExpSubTargets:   make([]string, 0),
		ExpPrepareModel: prepare,
		ExpScope:        scope,
	}
	for _, action := range commandSpec.Actions() {
		actionModel := spec.ActionModel{
			ActionName:      action.Name(),
			ActionAliases:   action.Aliases(),
			ActionShortDesc: action.ShortDesc(),
			ActionLongDesc:  action.LongDesc(),
			ActionExample:   action.Example(),
			ActionMatchers: func() []spec.ExpFlag {
				matchers := make([]spec.ExpFlag, 0)
				for _, m := range action.Matchers() {
					matchers = append(matchers, spec.ExpFlag{
						Name:                  m.FlagName(),
						Desc:                  m.FlagDesc(),
						NoArgs:                m.FlagNoArgs(),
						Required:              m.FlagRequired(),
						RequiredWhenDestroyed: m.FlagRequiredWhenDestroyed(),
					})
				}
				return matchers
			}(),
			ActionFlags: func() []spec.ExpFlag {
				flags := make([]spec.ExpFlag, 0)
				for _, m := range action.Flags() {
					flags = append(flags, spec.ExpFlag{
						Name:                  m.FlagName(),
						Desc:                  m.FlagDesc(),
						NoArgs:                m.FlagNoArgs(),
						Required:              m.FlagRequired(),
						RequiredWhenDestroyed: m.FlagRequiredWhenDestroyed(),
					})
				}
				for _, m := range commandSpec.Flags() {
					flags = append(flags, spec.ExpFlag{
						Name:                  m.FlagName(),
						Desc:                  m.FlagDesc(),
						NoArgs:                m.FlagNoArgs(),
						Required:              m.FlagRequired(),
						RequiredWhenDestroyed: m.FlagRequiredWhenDestroyed(),
					})
				}
				flags = append(flags,
					spec.ExpFlag{
						Name:                  "timeout",
						Desc:                  "set timeout for experiment",
						Required:              false,
						RequiredWhenDestroyed: false,
					},
					spec.ExpFlag{
						Name:     "async",
						Desc:     "whether to create asynchronously, default is false",
						Required: false,
					},
					spec.ExpFlag{
						Name:     "endpoint",
						Desc:     "the create result reporting address. It takes effect only when the async value is true and the value is not empty",
						Required: false,
					},
				)
				return flags
			}(),
			ActionPrograms:   action.Programs(),
			ActionCategories: action.Categories(),
		}
		model.ExpActions = append(model.ExpActions, actionModel)
	}
	models.Models = append(models.Models, model)
	return models
}

// AddModels adds the child model to parent
func AddModels(parent *spec.Models, child *spec.Models) {
	for idx, model := range parent.Models {
		for _, sub := range child.Models {
			model.ExpSubTargets = append(model.ExpSubTargets, sub.ExpName)
		}
		parent.Models[idx] = model
	}
}

// MergeModels
func MergeModels(models ...*spec.Models) *spec.Models {
	result := &spec.Models{
		Models: make([]spec.ExpCommandModel, 0),
	}
	for _, model := range models {
		result.Version = model.Version
		result.Kind = model.Kind
		result.Models = append(result.Models, model.Models...)
	}
	return result
}
