package main

import (
	"github.com/rosbit/go-eprolog"
	"os"
	"fmt"
	"unicode"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <prolog-file> <predict>[ <args>...]\n", os.Args[0])
		return
	}

	// 1. init prolog
	ctx := pl.NewProlog()

	plFile, predict := os.Args[1], os.Args[2]
	// 2. consult prolog script file
	if err := ctx.LoadFile(plFile); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// 3. prepare arguments and variables
	args := make([]interface{}, len(os.Args) - 3)
	for i,j:=3,0; i<len(os.Args); i,j = i+1,j+1 {
		arg := []rune(os.Args[i])
		if len(arg) > 0 && (arg[0] == '_' || unicode.IsUpper(arg[0])) {
			args[j] = pl.PlVar(os.Args[i])
		} else {
			args[j] = pl.PlStrTerm(os.Args[i])
		}
	}

	// 4. query the goal with arguments and variables
	solutions, ok, err := ctx.Query(predict, args...)

	// 5. check the result
	//  5.1 error checking
	if err != nil {
		fmt.Printf("failed to query %s: %v\n", predict, err)
		return
	}
	//  5.2 proving checking with result `false`
	if !ok {
		fmt.Printf("false\n")
		return
	}
	//  5.3 proving checking with result `true`
	if solutions == nil {
		fmt.Printf("true\n")
		return
	}

	//  5.4 result set processing
	i := 0
	for sol := range solutions {
		i += 1
		fmt.Printf("solution #%d: %#v\n", i, sol)
	}
}

