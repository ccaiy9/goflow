package __precheck_flow

import "goflow/base"

func RegisterPreCheckFlow(cts base.CtxStorage) *base.Flow {
	flow := base.NewFlow("PreCheckFlow")
	flow.SubmitTasks(registerPreCheckConnectTask(cts), registerPreCheckDiskTask(cts))

	return flow
}
