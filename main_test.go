package main

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/test/gtest"
	"goflow/base"
	__precheck_flow "goflow/flow_test/1_precheck_flow"
	__do_business_flow "goflow/flow_test/2_do_business_flow"
	"testing"
)

func helloExec(cts base.CtxStorage) error {
	fmt.Println("---------- hello Execute ------------------")
	return nil
}

func helloRoll(cts base.CtxStorage) error {
	fmt.Println("---------- hello rollback ------------------")
	return nil
}

func goodbyeExec(cts base.CtxStorage) error {
	fmt.Println("---------- goodbyeExec Execute ------------------")
	return fmt.Errorf("goodbyeExec Execute failed")
}

func goodbyeRoll(cts base.CtxStorage) error {
	fmt.Println("---------- goodbyeExec rollback ------------------")
	return nil
}

func loadFlow(ctxStorage map[string]interface{}) *base.Flow {
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("./install.yaml")

	mf := base.NewFlow("main-flow")
	helloTmp := base.NewTask("hello-task", helloExec, helloRoll)
	goodbyeTmp := base.NewTask("goodbye-task", goodbyeExec, goodbyeRoll)

	mf.SubmitTasks(
		helloTmp,
		__precheck_flow.RegisterPreCheckFlow(ctxStorage),
		__do_business_flow.RegisterPreCheckDoBusinessFlow(ctxStorage),
		goodbyeTmp,
	)

	return mf
}

func TestExampleExecuteLineFlow(t *testing.T) {

	gtest.C(t, func(xt *gtest.T) {
		ctxStorage := make(map[string]interface{})
		mf := loadFlow(ctxStorage)
		mf.Execute(ctxStorage)
		fmt.Println(mf.PrintErrors())

	})
}

func TestExampleManuRollBackLineFlow(t *testing.T) {

	gtest.C(t, func(xt *gtest.T) {
		ctxStorage := make(map[string]interface{})
		mf := loadFlow(ctxStorage)

		fid, tid, _ := mf.GetFailedHint()
		mf.UpdateFailedScene(fid, tid)
		mf.RollBackByManual(ctxStorage)
		fmt.Println(mf.PrintErrors())
	})
}
