package __do_business_flow

import (
	"fmt"
	"goflow/base"
)

func registerSendPkgTask(cts base.CtxStorage) base.AtomTask {
	return base.NewTask("registerSendPkgTask", sendPkgExecute, sendPkgRollback)
}

func sendPkgExecute(cts base.CtxStorage) error {
	fmt.Println("---------- flow-2-task-2 sendPkgExecute ------------------")
	return nil
}

func sendPkgRollback(cts base.CtxStorage) error {
	fmt.Println("---------- flow-2-task-2 sendPkgRollback ------------------")
	return nil
}
