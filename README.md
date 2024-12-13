#### Usage

------

This is a easy task-flow framework written in golang, referencing the taskflow of OpenStack.
Less code, but still provides complete functions.

- Support automatic rollback
- Support manual rollback
- Support rollback in failed scene or upper layer task
- Support still be executed if the rollback-tasks failed

You can execute the main_test.go file to see the effect.

```golang

// ...
mf.SubmitTasks(
    helloTmp,
    __precheck_flow.RegisterPreCheckFlow(ctxStorage),
    __do_business_flow.RegisterPreCheckDoBusinessFlow(ctxStorage),
    goodbyeTmp,
)
// ...

/*
    1. failed in  goodbyeTmp and rollback-tasks failed in PreCheckConnectTask and PreCheckDiskTask
    config: auto rollback, rollback in failed scene
*/

---------- hello Execute ------------------
---------- flow-1-task-1 checkConnectExecute ------------------
---------- flow-1-task-2 checkDiskExecute ------------------
---------- flow-2-task-1 -installToolsExecute ------------------
---------- flow-2-task-2 sendPkgExecute ------------------
---------- goodbyeExec Execute ------------------
---------- goodbyeExec rollback ------------------
---------- flow-2-task-2 sendPkgRollback ------------------
---------- flow-2-task-1 installToolsRollback ------------------
---------- flow-1-task-2 checkDiskRollback ------------------
---------- flow-1-task-1 checkConnectRollback ------------------
---------- hello rollback ------------------
errmsg: {"err_execute":{"errors":[{"err_flow":{"flow_id":4,"f_name":"main-flow"},"err_task":{"task_id":1,"t_name":"goodbye-task"},"exception":"goodbyeExec Execute failed"}]},"err_rollback":{"errors":[{"err_flow":{"flow_id":2,"f_name":"PreCheckConnectTask"},"err_task":{"task_id":1,"t_name":"PreCheckConnectTask"},"exception":"flow-1-task-1 checkConnectRollback failed"},{"err_flow":{"flow_id":2,"f_name":"PreCheckDiskTask"},"err_task":{"task_id":2,"t_name":"PreCheckDiskTask"},"exception":"flow-1-task-2 checkDiskRollback failed"}]}}


/*
   1. failed in  goodbyeTmp and rollback-tasks failed in PreCheckConnectTask and PreCheckDiskTask
   config: auto rollback, rollback in upper layer task
*/
---------- hello Execute ------------------
---------- flow-1-task-1 checkConnectExecute ------------------
---------- flow-1-task-2 checkDiskExecute ------------------
---------- flow-2-task-1 -installToolsExecute ------------------
---------- flow-2-task-2 sendPkgExecute ------------------
---------- goodbyeExec Execute ------------------
---------- flow-2-task-2 sendPkgRollback ------------------
---------- flow-2-task-1 installToolsRollback ------------------
---------- flow-1-task-2 checkDiskRollback ------------------
---------- flow-1-task-1 checkConnectRollback ------------------
---------- hello rollback ------------------
errmsg: {"err_execute":{"errors":[{"err_flow":{"flow_id":4,"f_name":"main-flow"},"err_task":{"task_id":1,"t_name":"goodbye-task"},"exception":"goodbyeExec Execute failed"}]},"err_rollback":{"errors":[{"err_flow":{"flow_id":2,"f_name":"PreCheckConnectTask"},"err_task":{"task_id":1,"t_name":"PreCheckConnectTask"},"exception":"flow-1-task-1 checkConnectRollback failed"},{"err_flow":{"flow_id":2,"f_name":"PreCheckDiskTask"},"err_task":{"task_id":2,"t_name":"PreCheckDiskTask"},"exception":"flow-1-task-2 checkDiskRollback failed"}]}}

```