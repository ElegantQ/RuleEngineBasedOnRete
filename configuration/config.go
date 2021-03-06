package configuration

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/rules/common/model"
)

// RuleSessionDescriptor is a collection of rules to be loaded

type RuleActionDescriptor struct {
	Name       string               `json:"name"`
	IOMetadata *metadata.IOMetadata `json:"metadata"`
	Rules      []*RuleDescriptor    `json:"rules"`
}

type RuleSessionDescriptor struct {
	Rules []*RuleDescriptor `json:"rules"`
}

// RuleDescriptor defines a rule
type RuleDescriptor struct {
	Name        string
	Conditions  []*ConditionDescriptor
	ActionFunc  model.ActionFunction
	Priority    int
	Identifiers []string
}

// ConditionDescriptor defines a condition in a rule
type ConditionDescriptor struct {
	Name        string
	Identifiers []string
	Evaluator   model.ConditionEvaluator
	Expression  string
}

func (c *RuleDescriptor) UnmarshalJSON(d []byte) error {
	ser := &struct {
		Name         string                 `json:"name"`
		Conditions   []*ConditionDescriptor `json:"conditions"`
		ActionFuncId string                 `json:"actionFunction"`
		Priority     int                    `json:"priority"`
		Identifiers  []string               `json:"identifiers"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	c.Name = ser.Name
	c.Conditions = ser.Conditions
	c.ActionFunc = GetActionFunction(ser.ActionFuncId)
	c.Priority = ser.Priority
	c.Identifiers = ser.Identifiers

	return nil
}

func (c *RuleDescriptor) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"name\":" + "\"" + c.Name + "\",")
	if c.Identifiers != nil {
		buffer.WriteString("\"identifiers\":[")
		for _, id := range c.Identifiers {
			buffer.WriteString("\"" + id + "\",")
		}
		buffer.Truncate(buffer.Len() - 1)
		buffer.WriteString("],")
	}

	buffer.WriteString("\"conditions\":[")
	for _, condition := range c.Conditions {
		jsonCondition, err := condition.MarshalJSON()
		if err == nil {
			buffer.WriteString(string(jsonCondition) + ",")
		}
	}
	buffer.Truncate(buffer.Len() - 1)
	buffer.WriteString("],")

	actionFunctionID := GetActionFunctionID(c.ActionFunc)
	buffer.WriteString("\"actionFunction\":\"" + actionFunctionID + "\",")
	buffer.WriteString("\"priority\":" + strconv.Itoa(c.Priority) + "}")

	return buffer.Bytes(), nil
}

func (c *ConditionDescriptor) UnmarshalJSON(d []byte) error {
	ser := &struct {
		Name        string   `json:"name"`
		Identifiers []string `json:"identifiers"`
		EvaluatorId string   `json:"evaluator"`
		Expression  string   `json:"expression"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	c.Name = ser.Name
	c.Identifiers = ser.Identifiers
	c.Evaluator = GetConditionEvaluator(ser.EvaluatorId)
	c.Expression = ser.Expression

	return nil
}

func (c *ConditionDescriptor) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"name\":" + "\"" + c.Name + "\",")
	if c.Identifiers != nil {
		buffer.WriteString("\"identifiers\":[")
		for _, id := range c.Identifiers {
			buffer.WriteString("\"" + id + "\",")
		}
		buffer.Truncate(buffer.Len() - 1)
		buffer.WriteString("],")
	}

	conditionEvaluatorID := GetConditionEvaluatorID(c.Evaluator)
	buffer.WriteString("\"evaluator\":\"" + conditionEvaluatorID + "\",")
	buffer.WriteString("\"expression\":\"" + c.Expression + "\"}")

	return buffer.Bytes(), nil
}

//metadata support
type DefinitionConfig struct {
	Name     string               `json:"name"`
	Metadata *metadata.IOMetadata `json:"metadata"`
	Rules    []*RuleDescriptor    `json:"rules"`
}
