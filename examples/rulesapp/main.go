package main

import (
	"context"
	"fmt"
	"github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/mqtt"
	"github.com/project-flogo/rules/ruleapi"
	"sync"
)

func main() {
	//Load the tuple descriptor file (relative to GOPATH)
	tupleDescAbsFileNm := common.GetAbsPathForResource("src/github.com/project-flogo/rules/examples/rulesapp/rulesapp.json")
	tupleDescriptor := common.FileToString(tupleDescAbsFileNm)

	fmt.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	//First register the tuple descriptors
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	//Create a RuleSession
	rs, _ := ruleapi.GetOrCreateRuleSession("asession")

	//// check for name "Bob" in n1
	rule := ruleapi.NewRule("n1.name == Bob")
	rule.AddCondition("c1", []string{"n1"}, checkForBob, nil)
	rule.SetAction(checkForBobAction)
	rule.SetContext("This is a test of context")
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	// check for name "Bob" in n1, match the "name" field in n2,
	// in effect, fire the rule when name field in both tuples is "Bob"
	//rule2 := ruleapi.NewRule("n1.name == Bob && n1.name == n2.name")
	//rule2.AddCondition("c1", []string{"n1"}, checkForBob, nil)
	//rule2.AddCondition("c2", []string{"n1", "n2"}, checkSameNamesCondition, nil)
	//rule2.SetAction(checkSameNamesAction)
	//rs.AddRule(rule2)
	//fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	//Start the rule session
	rs.Start(nil)

	//Now assert a "n1" tuple
	fmt.Println("Asserting n1 tuple with name=Tom")
	t1, _ := model.NewTupleWithKeyValues("n1", "Tom")
	t1.SetString(context.TODO(), "name", "Tom")
	rs.Assert(context.TODO(), t1)

	//Now assert a "n1" tuple
	fmt.Println("Asserting n1 tuple with name=Bob")
	t2, _ := model.NewTupleWithKeyValues("n1", "Bob")
	t2.SetString(context.TODO(), "name", "Bob")
	rs.Assert(context.TODO(), t2)

	//Now assert a "n2" tuple
	fmt.Println("Asserting n2 tuple with name=Bob")
	t3, _ := model.NewTupleWithKeyValues("n2", "Bob")
	t3.SetString(nil, "name", "Bob")
	rs.Assert(context.TODO(), t3)

	//Retract tuples
	rs.Retract(context.TODO(), t1)
	rs.Retract(context.TODO(), t2)
	rs.Retract(context.TODO(), t3)

	//delete the rule
	//rs.DeleteRule(rule.GetName())

	//unregister the session, i.e; cleanup
	//rs.Unregister()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go mqtt.Init()
	wg.Wait()
}

func checkForBob(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	//This conditions filters on name="Bob"
	t1 := tuples["n1"]
	if t1 == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	name, _ := t1.GetString("name")
	return name == "Bob"
}

func checkForBobAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	fmt.Printf("Context is [%s]\n", ruleCtx)
	t1 := tuples["n1"]
	if t1 == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return
	}
}

func checkSameNamesCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	t1 := tuples["n1"]
	t2 := tuples["n2"]
	if t1 == nil || t2 == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return false
	}
	name1, _ := t1.GetString("name")
	name2, _ := t2.GetString("name")
	return name1 == name2
}

func checkSameNamesAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	t1 := tuples["n1"]
	t2 := tuples["n2"]
	if t1 == nil || t2 == nil {
		fmt.Println("Should not get nil tuples here in Action! This is an error")
		return
	}
	name1, _ := t1.GetString("name")
	name2, _ := t2.GetString("name")
	fmt.Printf("n1.name = [%s], n2.name = [%s]\n", name1, name2)
}
