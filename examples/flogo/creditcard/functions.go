package main

import (
	"context"
	"fmt"

	"github.com/project-flogo/rules/configuration"

	"github.com/project-flogo/rules/common/model"
)

//add this sample file to your flogo project
func init() {

	// rule UserData
	configuration.RegisterConditionEvaluator("cNewUser", cNewUser)
	configuration.RegisterActionFunction("aNewUser", aNewUser)

	// rule NewUser
	configuration.RegisterConditionEvaluator("cBadUser", cBadUser)
	configuration.RegisterActionFunction("aBadUser", aBadUser)

	// rule NewUserApprove
	configuration.RegisterConditionEvaluator("cUserIdMatch", cUserIdMatch)
	configuration.RegisterConditionEvaluator("cUserCreditScore", cUserCreditScore)
	configuration.RegisterActionFunction("aApproveWithLowerLimit", aApproveWithLowerLimit)

	// // rule NewUserReject
	configuration.RegisterConditionEvaluator("cUserIdMatch", cUserIdMatch)
	configuration.RegisterConditionEvaluator("cUserLowCreditScore", cUserLowCreditScore)
	configuration.RegisterActionFunction("aUserReject", aUserReject)

	// // rule NewUserApprove1
	configuration.RegisterConditionEvaluator("cUserIdMatch", cUserIdMatch)
	configuration.RegisterConditionEvaluator("cUserHighCreditScore", cUserHighCreditScore)
	configuration.RegisterActionFunction("aApproveWithHigherLimit", aApproveWithHigherLimit)
}

func cNewUser(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	newaccount := tuples["NewAccount"]
	if newaccount != nil {
		address, _ := newaccount.GetString("Address")
		age, _ := newaccount.GetInt("Age")
		income, _ := newaccount.GetInt("Income")
		if address != "" && age >= 18 && income >= 10000 && age <= 44 {
			return true
		}
	}
	return false
}

func aNewUser(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	newaccount := tuples["NewAccount"]
	Id, _ := newaccount.GetInt("Id")
	name, _ := newaccount.GetString("Name")
	address, _ := newaccount.GetString("Address")
	age, _ := newaccount.GetInt("Age")
	income, _ := newaccount.GetInt("Income")
	userInfo, _ := model.NewTupleWithKeyValues("UserAccount", Id)
	userInfo.SetString(ctx, "Name", name)
	userInfo.SetString(ctx, "Addresss", address)
	userInfo.SetInt(ctx, "Age", age)
	userInfo.SetInt(ctx, "Income", income)
	fmt.Println(userInfo)
	rs.Assert(ctx, userInfo)
	fmt.Println("User information received")

}

func cBadUser(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	newaccount := tuples["NewAccount"]
	if newaccount != nil {
		address, _ := newaccount.GetString("Address")
		age, _ := newaccount.GetInt("Age")
		income, _ := newaccount.GetInt("Income")
		if address == "" || age < 18 || income < 10000 || age >= 45 {
			return true
		}
	}
	return false
}

func aBadUser(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	fmt.Println("Applicant is not eligible to apply for creditcard")
}

func cUserIdMatch(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	userInfo := tuples["UserAccount"]
	updateScore := tuples["UpdateCreditScore"]
	if userInfo != nil || updateScore != nil {
		userId, _ := userInfo.GetInt("Id")
		newUserId, _ := updateScore.GetInt("Id")
		if userId == newUserId {
			fmt.Println("Userid match found")
			return true
		}
	}
	return false
}

func cUserCreditScore(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	updateScore := tuples["UpdateCreditScore"]
	if updateScore != nil {
		CreditScore, _ := updateScore.GetInt("creditScore")
		if CreditScore >= 750 && CreditScore < 820 {
			fmt.Println("cUserCreditScore")
			return true
		}
	}
	return false
}

func cUserLowCreditScore(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	updateScore := tuples["UpdateCreditScore"]
	if updateScore != nil {
		CreditScore, _ := updateScore.GetInt("creditScore")
		if CreditScore < 750 {
			fmt.Println("cUserLowCreditScore")
			return true
		}
	}
	return false
}

func cUserHighCreditScore(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	updateScore := tuples["UpdateCreditScore"]
	if updateScore != nil {
		CreditScore, _ := updateScore.GetInt("creditScore")
		if CreditScore >= 820 && CreditScore <= 900 {
			fmt.Println("cUserHighCreditScore")
			return true
		}
	}
	return false
}

func aApproveWithLowerLimit(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	userInfo := tuples["UserAccount"]
	updateScore := tuples["UpdateCreditScore"]
	CreditScore, _ := updateScore.GetInt("creditScore")
	income, _ := userInfo.GetInt("Income")
	var limit = 2 * income
	userInfoMutable := userInfo.(model.MutableTuple)
	userInfoMutable.SetInt(ctx, "creditScore", CreditScore)
	userInfoMutable.SetString(ctx, "appStatus", "Approved")
	userInfoMutable.SetInt(ctx, "approvedLimit", limit)
	fmt.Println(userInfo)
}

func aApproveWithHigherLimit(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	userInfo := tuples["UserAccount"]
	updateScore := tuples["UpdateCreditScore"]
	CreditScore, _ := updateScore.GetInt("creditScore")
	income, _ := userInfo.GetInt("Income")
	var limit = 3 * income
	userInfoMutable := userInfo.(model.MutableTuple)
	userInfoMutable.SetInt(ctx, "creditScore", CreditScore)
	userInfoMutable.SetString(ctx, "appStatus", "Approved")
	userInfoMutable.SetInt(ctx, "approvedLimit", limit)
	fmt.Println(userInfo)
}

func aUserReject(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	userInfo := tuples["UserAccount"]
	updateScore := tuples["UpdateCreditScore"]
	CreditScore, _ := updateScore.GetInt("creditScore")
	userInfoMutable := userInfo.(model.MutableTuple)
	userInfoMutable.SetInt(ctx, "creditScore", CreditScore)
	userInfoMutable.SetString(ctx, "appStatus", "Rejected")
	userInfoMutable.SetInt(ctx, "approvedLimit", 0)
	fmt.Println(userInfo)
}
