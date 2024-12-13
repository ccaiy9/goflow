package base

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"os"
	"strconv"
	"strings"
)

const (
	p = "./flow.hint"
)

type failedScene struct {
	f int8
	t int8
}

// Flow line-flow
type Flow struct {
	name             string
	autoRollBack     bool
	rollbackInFailed bool

	fId   int8
	tasks []AtomTask
}

func NewFlow(name string) *Flow {
	if name == "" {
		fmt.Errorf("new task param failed")
		return nil
	}

	autoRollBack, err := g.Cfg().Get(context.Background(), "policy.auto_rollback")
	if err != nil {
		fmt.Errorf("get installer config failed: %+v", err)
		return nil
	}

	rollbackInScene, err := g.Cfg().Get(context.Background(), "policy.rollback_in_scene")
	if err != nil {
		fmt.Errorf("get installer config failed: %+v", err)
		return nil
	}

	return &Flow{
		name:             name,
		autoRollBack:     autoRollBack.Bool(),
		rollbackInFailed: rollbackInScene.Bool(),
	}
}

func (f *Flow) Name() string {
	return f.name
}

func (f *Flow) Id() int8 {
	return f.fId
}

func (f *Flow) SetFId(fid int8) {
	f.fId = fid
}

func (f *Flow) SetTId(fid int8) {
	//f.fId = fid
}

func (f *Flow) GetFId() int8 {
	return f.fId
}

func (f *Flow) GetTId() int8 {
	return 0
}

// SetErrors no use
func (f *Flow) SetErrors(isExecute bool, fid, tid int8, fName, tName string, err error) {
	return
}

func (f *Flow) GetErrors() (*ErrExecute, *ErrRollback) {
	return nil, nil
}

var gFs failedScene
var gExecuteOut bool
var gRollBackOut bool
var gRestMainFlow bool

func (f *Flow) Execute(cts CtxStorage) error {
	for _, task := range f.tasks {
		if err := task.Execute(cts); err != nil {

			if !gExecuteOut {
				gFs = failedScene{task.GetFId(), task.GetTId()}
				gExecuteOut = true
			} else {
				// in main-flow
				gRestMainFlow = true
			}

			if t, ok := task.(*atomTaskBase); ok {
				t.SetErrors(true, t.GetFId(), t.GetTId(), f.Name(), t.Name(), err)
			}

			if f.autoRollBack {
				if err := f.Rollback(cts); err != nil {
					fmt.Errorf("rollback failed failed: %+v", err)
					return nil
				}
			} else {

				if tk, ok := task.(*atomTaskBase); ok {
					if errS := f.syncFailedHint(tk); errS != nil {
						fmt.Errorf("syncFailedHint failed: %+v", errS)
					}
				}
			}

			return err
		}
	}

	return nil
}

func (f *Flow) Rollback(cts CtxStorage) error {

	var seq = int8(len(f.tasks))
	if gRestMainFlow {
		seq = gFs.f - 1
		gRestMainFlow = false
	}

	if !gRollBackOut {
		if f.GetFId() == 0 {
			if f.rollbackInFailed {
				seq = gFs.f
			} else {
				seq = gFs.f - 1
			}

		} else {

			if f.rollbackInFailed {
				seq = gFs.t
			} else {
				seq = gFs.t - 1
			}
		}

		gRollBackOut = true
	}

	for i := seq - 1; i >= 0; i-- {
		tk := f.tasks[i]
		err := tk.Rollback(cts)
		if err != nil {
			tk.SetErrors(false, f.Id(), tk.Id(), tk.Name(), tk.Name(), err)
		}
	}

	return nil
}

func (f *Flow) syncTaskIndex(ts *[]AtomTask, depth int8, fIndex *int8) {

	for index, t := range *ts {
		if depth == 0 {
			*fIndex++
			index = 0
		}

		if at, ok := t.(*atomTaskBase); ok {
			at.SetTId(int8(index + 1))
			at.SetFId(*fIndex)

		} else {
			t.(*Flow).SetFId(*fIndex)
			f.syncTaskIndex(&t.(*Flow).tasks, depth+1, fIndex)
		}
	}
}

func (f *Flow) SubmitTasks(ts ...AtomTask) {

	f.tasks = append(f.tasks, ts...)
	var fIndex int8 = 0

	f.syncTaskIndex(&f.tasks, 0, &fIndex)
}

func (f *Flow) PrintErrors() string {
	var flowErrors FlowError

	handler := func(t AtomTask) {
		execErr, robErr := t.GetErrors()
		if execErr != nil {
			flowErrors.ErrExecute.Errors = append(flowErrors.ErrExecute.Errors, execErr.Errors...)
		}
		if robErr != nil {
			flowErrors.ErrRollback.Errors = append(flowErrors.ErrRollback.Errors, robErr.Errors...)
		}
	}

	for _, task := range f.tasks {

		if at, ok := task.(*Flow); ok {
			for _, a := range at.tasks {
				handler(a)
			}
		} else {
			handler(task)
		}
	}

	reqJSON, _ := json.Marshal(flowErrors)
	return string(reqJSON)
}

func (f *Flow) syncFailedHint(tk *atomTaskBase) error {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		file, err := os.Create(p)
		if err != nil {
			fmt.Errorf("error creating file: %+v", err)
			return err
		}
		file.Close()
	}

	file, err := os.OpenFile(p, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Errorf("error opening file: %+v", err)
		return err
	}
	defer file.Close()

	s := fmt.Sprintf("%d-%d", tk.GetFId(), tk.GetTId())
	_, err = file.WriteString(s)
	if err != nil {
		fmt.Errorf("error writing to file: %+v", err)
		return err
	}

	return nil
}

func (f *Flow) GetFailedHint() (int8, int8, error) {
	r, err := os.Open(p)
	if err != nil {
		fmt.Errorf("open %s failed: %+v", p, err)
		return 0, 0, err
	}

	scanner := bufio.NewScanner(r)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}

	parts := strings.Split(line, "-")
	fid, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Errorf("trans fid failed: %+v", err)
		return 0, 0, err
	}
	tid, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Errorf("trans tid failed: %+v", err)
		return 0, 0, err
	}

	return int8(fid), int8(tid), nil
}

func (f *Flow) UpdateFailedScene(fid, tid int8) {
	gFs = failedScene{fid, tid}
}

var gManuRollBackOut bool

func (f *Flow) RollBackByManual(cts CtxStorage) error {

	var seq = gFs.f

	for i := seq - 1; i >= 0; i-- {
		tk := f.tasks[i]

		if !gManuRollBackOut {
			gManuRollBackOut = true

			if _, ok := tk.(*atomTaskBase); ok {
				gRollBackOut = true
				if !f.rollbackInFailed {
					continue
				}
			}
		}

		err := tk.Rollback(cts)
		if err != nil {
			tk.SetErrors(false, f.Id(), tk.Id(), tk.Name(), tk.Name(), err)
		}
	}

	return nil
}
