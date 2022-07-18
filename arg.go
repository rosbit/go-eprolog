package pl

import (
	"github.com/ichiban/prolog/engine"
	"reflect"
	"fmt"
	"bytes"
)

type PlArg interface {
	ToTerm() engine.Term
}

var (
	_ PlBool  = PlBool(false)
	_ PlInt   = PlInt(0)
	_ PlFloat = PlFloat(0.0)
	_ PlString= PlString("")
	_ PlVar   = PlVar("")
	_ PlStrTerm = PlStrTerm("[]")
	_ *PlList = &PlList{}
	_ *PlRecord = &PlRecord{}
)

func makePlTerm(v interface{}) (engine.Term, error) {
	if v == nil {
		return engine.Atom("[]"), nil
	}

	switch vv := v.(type) {
	case int, int8, int16, int32, int64,
	     uint,uint8,uint16,uint32,uint64:
		return makeInt(v), nil
	case string:
		return PlString(vv).ToTerm(), nil
	case bool:
		return PlBool(vv).ToTerm(), nil
	case float64:
		return PlFloat(vv).ToTerm(), nil
	case float32:
		return PlFloat(float64(vv)).ToTerm(), nil
	case PlVar:
		return vv.ToTerm(), nil
	case PlStrTerm:
		return vv.ToTerm(), nil
	case Record:
		r, err := newRecord(vv)
		if err != nil {
			return engine.Atom("[]"), err
		}
		return r.ToTerm(), nil
	default:
	}

	vv := reflect.ValueOf(v)
	switch vv.Kind() {
	case reflect.Slice:
		t := vv.Type()
		if t.Elem().Kind() == reflect.Uint8 {
			return PlString(string(v.([]byte))).ToTerm(), nil
		}
		fallthrough
	case reflect.Array:
		if plL, err := newPlList(v); err != nil {
			return engine.Atom("[]"), err
		} else {
			return plL.ToTerm(), nil
		}
	case reflect.Ptr:
		switch vv.Elem().Kind() {
		case reflect.Array, reflect.Struct:
			if plL, err := newPlList(v); err != nil {
				return engine.Atom("[]"), err
			} else {
				return plL.ToTerm(), nil
			}
		}
		return makePlTerm(vv.Elem().Interface())
	/*
	case reflect.Map:
	case reflect.Struct:
	case reflect.Func:
	*/
	default:
		return nil, fmt.Errorf("unsupported type %v", vv.Kind())
	}
}

func makeInt(i interface{}) engine.Term {
	switch i.(type) {
	case int,int8,int16,int32,int64:
		return PlInt(reflect.ValueOf(i).Int()).ToTerm()
	case uint8,uint16,uint32:
		return PlInt(int64(reflect.ValueOf(i).Uint())).ToTerm()
	case uint,uint64:
		return PlFloat(float64(reflect.ValueOf(i).Uint())).ToTerm()
	default:
		return engine.Integer(0)
	}
}

func fromPlTerm(plVal engine.Term) (goVal interface{}, err error) {
	if plVal == nil {
		return
	}
	switch v := plVal.(type) {
	case engine.Atom:
		s := string(v)
		l := len(s)
		if l <= 2 {
			goVal = s
			return
		}
		if s[0] == '\'' && s[l-1] == '\'' {
			goVal = s[1:l-1]
			return
		}
		goVal = s
		return
	case engine.Integer:
		goVal = int64(v)
		return
	case engine.Float:
		goVal = float64(v)
		return
	case *engine.Compound:
		if string(v.Functor) == "." {
			args := v.Args
			res, e := fromPlCons(args)
			if e != nil {
				err = e
				return
			}
			goVal = res
			return
		}
		b := &bytes.Buffer{}
		err = plVal.WriteTerm(b, &engine.WriteOptions{Quoted:true}, engine.NewEnv())
		if err != nil {
			goVal = b.String()
		}
		return
	// case term.VariableType:
	// case term.ErrorType:
	default:
		err = fmt.Errorf("unsupported type")
		return
	}
}

func fromPlCons(cons []engine.Term) (res []interface{}, err error) {
	head, e := fromPlTerm(cons[0])
	if e != nil {
		err = e
		return
	}
	res = append(res, head)
	for cons[1] != nil {
		if tail, ok := cons[1].(*engine.Compound); ok {
			if string(tail.Functor) == "." {
				cons = tail.Args
				head, e = fromPlTerm(cons[0])
				if e != nil {
					err = e
					break
				}
				res = append(res, head)
			}
		} else {
			break
		}
	}
	return
}

// var
type PlVar string
func (v PlVar) ToTerm() engine.Term {
	return engine.Variable(string(v))
}

// bool
type PlBool bool
func (b PlBool) ToTerm() engine.Term {
	if b {
		return engine.Atom("true")
	}
	return engine.Atom("false")
}

// int
type PlInt int64
func (i PlInt) ToTerm() engine.Term {
	return engine.Integer(int64(i))
}

// float
type PlFloat float64
func (f PlFloat) ToTerm() engine.Term {
	return engine.Float(float64(f))
}

// string
type PlString string
func (s PlString) ToTerm() engine.Term {
	return engine.Atom(string(s))
}

// string -> term
type PlStrTerm string
func (s PlStrTerm) ToTerm() engine.Term {
	return engine.Atom(string(s))
}

// list
type PlList struct {
	pa engine.Term
}
func newPlList(a interface{}) (plL *PlList, err error) {
	if a == nil {
		plL = &PlList{engine.Atom("[]")}
		return
	}
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		l := v.Len()
		pa := make([]engine.Term, l)
		for i:=0; i<l; i++ {
			e := v.Index(i).Interface()
			ev, err := makePlTerm(e)
			if err != nil {
				pa[i] = engine.Atom("[]")
			} else {
				pa[i] = ev
			}
		}
		plL = &PlList{engine.List(pa...)}
		return
	/*
	case reflect.Map:
	case reflect.Struct:
	*/
	case reflect.Ptr:
		return newPlList(v.Elem().Interface())
	default:
		err = fmt.Errorf("slice or list expected")
		return
	}
}
func (a *PlList) ToTerm() engine.Term {
	return a.pa
}

// record
type PlRecord struct {
	r Record
}
func newRecord(rec Record) (*PlRecord, error) {
	if len(rec.TableName()) == 0 {
		return nil, fmt.Errorf("table name expected for record")
	}
	return &PlRecord{r: rec}, nil
}
func (r *PlRecord) ToTerm() engine.Term {
	fields := r.r.FieldValues()
	ts := make([]engine.Term, len(fields))
	for i, f := range fields {
		ts[i], _ = makePlTerm(f)
	}
	return &engine.Compound{
		Functor: engine.Atom(r.r.TableName()),
		Args: ts,
	}
}
