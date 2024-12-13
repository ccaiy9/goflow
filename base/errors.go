package base

type ErrFlow struct {
	FlowId int8   `json:"flow_id"`
	FName  string `json:"f_name"`
}

type ErrTask struct {
	TaskId int8   `json:"task_id"`
	TName  string `json:"t_name"`
}

type ErrorInfo struct {
	*ErrFlow  `json:"err_flow"`
	*ErrTask  `json:"err_task"`
	Exception string `json:"exception"`
}

type ErrExecute struct {
	Errors []ErrorInfo `json:"errors"`
}

type ErrRollback struct {
	Errors []ErrorInfo `json:"errors"`
}

type FlowError struct {
	//ErrFlow
	//ErrTask
	//
	//ErrExecute   error
	//ErrRollbacks []error
	ErrExecute  `json:"err_execute"`
	ErrRollback `json:"err_rollback"`
}

func (i *FlowError) SetExecuteErr(f, t int8, fName, tName string, err error) {
	ei := ErrorInfo{
		ErrFlow: &ErrFlow{
			FlowId: f,
			FName:  fName,
		},
		ErrTask: &ErrTask{
			TaskId: t,
			TName:  tName,
		},
		Exception: err.Error(),
	}
	i.ErrExecute.Errors = append(i.ErrExecute.Errors, ei)
}

func (i *FlowError) SetRollbackErr(f, t int8, fName, tName string, err error) {
	ei := ErrorInfo{
		ErrFlow: &ErrFlow{
			FlowId: f,
			FName:  fName,
		},
		ErrTask: &ErrTask{
			TaskId: t,
			TName:  tName,
		},
		Exception: err.Error(),
	}

	i.ErrRollback.Errors = append(i.ErrRollback.Errors, ei)
}

func (i *FlowError) GetExecuteErr() *ErrExecute {
	//return fmt.Sprintf("%+v", i.ErrExecute)
	//reqJSON, _ := json.Marshal(i.ErrExecute)
	//return string(reqJSON)
	return &i.ErrExecute
}

func (i *FlowError) GetRollbackErr() *ErrRollback {
	//reqJSON, _ := json.Marshal(i.ErrRollback)
	//return string(reqJSON)
	return &i.ErrRollback
}
