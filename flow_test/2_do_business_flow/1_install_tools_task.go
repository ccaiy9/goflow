package __do_business_flow

import (
	"fmt"
	"goflow/base"
)

func registerInstallToolsTask(cts base.CtxStorage) base.AtomTask {
	return base.NewTask("registerInstallToolsTask", installToolsExecute, installToolsRollback)
}

func installToolsExecute(cts base.CtxStorage) error {
	fmt.Println("---------- flow-2-task-1 -installToolsExecute ------------------")
	return nil
}

func installToolsRollback(cts base.CtxStorage) error {
	fmt.Println("---------- flow-2-task-1 installToolsRollback ------------------")
	return nil
}
