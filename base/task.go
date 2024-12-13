package base

import "fmt"

type AtomTask interface {
	Name() string
	Id() int8

	Execute(CtxStorage) error
	Rollback(CtxStorage) error

	SetErrors(isExecute bool, f, t int8, fName, tName string, err error)
	GetErrors() (*ErrExecute, *ErrRollback)

	SetFId(fid int8)
	SetTId(tid int8)

	GetFId() int8
	GetTId() int8
}

type TaskFunc func(CtxStorage) error
type CtxStorage map[string]interface{}

func NewTask(name string, execute TaskFunc, rollback TaskFunc) AtomTask {

	if name == "" || execute == nil || rollback == nil {
		fmt.Errorf("new task param failed")
		return nil
	}

	return &atomTaskBase{name: name, exec: execute, rb: rollback, errors: &FlowError{}}
}

type atomTaskBase struct {
	name   string
	tId    int8
	fId    int8
	exec   TaskFunc
	rb     TaskFunc
	errors *FlowError
}

func (a *atomTaskBase) Name() string {
	return a.name
}

func (a *atomTaskBase) Id() int8 {
	return a.tId
}

func (a *atomTaskBase) Execute(cts CtxStorage) error {
	return a.exec(cts)
}

func (a *atomTaskBase) Rollback(cts CtxStorage) error {
	return a.rb(cts)
}

func (a *atomTaskBase) SetErrors(isExecute bool, f, t int8, fName, tName string, err error) {

	if isExecute {
		a.errors.SetExecuteErr(f, t, fName, tName, err)
	} else {
		a.errors.SetRollbackErr(f, t, fName, tName, err)
	}
}

func (a *atomTaskBase) GetErrors() (*ErrExecute, *ErrRollback) {
	return a.errors.GetExecuteErr(), a.errors.GetRollbackErr()
}

func (a *atomTaskBase) SetFId(fid int8) {
	a.fId = fid
}

func (a *atomTaskBase) SetTId(tid int8) {
	a.tId = tid
}

func (a *atomTaskBase) GetFId() int8 {
	return a.fId
}

func (a *atomTaskBase) GetTId() int8 {
	return a.tId
}
