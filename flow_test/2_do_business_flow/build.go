package __do_business_flow

import "goflow/base"

func RegisterPreCheckDoBusinessFlow(cts base.CtxStorage) *base.Flow {
	flow := base.NewFlow("RegisterPreCheckDoBusinessFlow")
	flow.SubmitTasks(registerInstallToolsTask(cts), registerSendPkgTask(cts))

	return flow
}
