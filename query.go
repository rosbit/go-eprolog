package pl

import (
	"fmt"
	"bytes"
)

func (p *Prolog) Query(goal string, args ...interface{}) (it <-chan map[string]interface{}, ok bool, err error) {
	if len(goal) == 0 {
		err = fmt.Errorf("goal expected")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			if v, o := r.(error); o {
				err = v
				return
			}
			err = fmt.Errorf("%v", r)
			return
		}
	}()
	argc, strArgs, argv, vars, e := makeGoalArgs(args...)
	if e != nil {
		err = e
		return
	}

	return p.doQuery(goal, argc, strArgs, argv, vars)
}

func makeGoalArgs(args ...interface{}) (argc int, strArgs string, argv []interface{}, vars map[string]interface{}, err error) {
	argc = len(args)
	if argc == 0 {
		return
	}

	argv = make([]interface{}, argc)
	argvCount := 0
	argsBuf := &bytes.Buffer{}
	vars = make(map[string]interface{})
	for i, arg := range args {
		if i > 0 {
			argsBuf.WriteByte(',')
		}
		switch arg.(type) {
		case PlVar:
			plVar := arg.(PlVar)
			vName := string(plVar)
			if len(vName) == 0 {
				vName = fmt.Sprintf("_Var%d", i)
			}
			// vars[i] = vName
			// argv[i] = PlVar(vName).ToTerm()
			argsBuf.WriteString(vName)
			vars[vName] = nil
		case PlStrTerm:
			plStrTerm := arg.(PlStrTerm)
			argsBuf.WriteString(string(plStrTerm))
		default:
			t, e := makePlTerm(arg)
			if e != nil {
				err = e
				return
			}
			argsBuf.WriteByte('?')
			argv[argvCount] = t
			argvCount += 1
		}
	}
	argv = argv[:argvCount]
	strArgs = argsBuf.String()
	return
}

func (p *Prolog) doQuery(goal string, argc int, strArgs string, argv []interface{}, vars map[string]interface{}) (it <-chan map[string]interface{}, ok bool, err error) {
	realGoal := fmt.Sprintf("%s(%s).", goal, strArgs)
	sols, e := p.i.Query(realGoal, argv...)
	if e != nil {
		err = e
		return
	}

	if len(vars) == 0 {
		ok = sols.Next()
		sols.Close()
		return
	}

	ok = true
	// var binding
	res := make(chan map[string]interface{})
	go func() {
		for sols.Next() {
			m := map[string]interface{}{}
			if e := sols.Scan(m); e != nil {
				err = e
				break
			}
			res <- m
		}
		sols.Close()
		close(res)
	}()

	it = res
	return
}
