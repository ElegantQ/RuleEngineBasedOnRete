package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

//1 condition, 1 expression
func Test_1_Expr(t *testing.T) {

	actionCount := map[string]int{"count": 0}
	rs, _ := createRuleSession()
	rs.Start(nil)
	r1 := ruleapi.NewRule("r1")
	r1.AddExprCondition("c1", "$.t2.p2 > $.t1.p1", nil)
	r1.SetAction(a1)
	r1.SetContext(actionCount)

	rs.AddRule(r1)


	var ctx context.Context

	t1, _ := model.NewTupleWithKeyValues("t1", "t1")
	t1.SetInt(context.TODO(), "p1", 1)
	t1.SetDouble(context.TODO(), "p2", 1.3)
	t1.SetString(context.TODO(), "p3", "t3")

	ctx = context.WithValue(context.TODO(), TestKey{}, t)
	rs.Assert(ctx, t1)

	t2, _ := model.NewTupleWithKeyValues("t2", "t2")
	t2.SetInt(context.TODO(), "p1", 1)
	t2.SetDouble(context.TODO(), "p2", 1.0001)
	t2.SetString(context.TODO(), "p3", "t3")

	ctx = context.WithValue(context.TODO(), TestKey{}, t)
	rs.Assert(ctx, t2)
	rs.Unregister()
	count := actionCount["count"]
	if count != 1 {
		t.Errorf("expected [%d], got [%d]\n", 1, count)
	}
}

func a1(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t := ctx.Value(TestKey{}).(*testing.T)
	t.Logf("Test_1_Expr executed!")
	actionCount := ruleCtx.(map[string]int)
	count := actionCount["count"]
	actionCount["count"] = count + 1
}
