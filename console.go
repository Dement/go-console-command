package command

import (
	"os"
	"reflect"
	"errors"
	"log"
	"unicode"
	"strconv"
)

var funcs = map[string]interface{} {}

func Run()  {
	args := os.Args
	if len(args) < 2 {
		var list string
		for key, _ := range funcs {
			list = list + key + "\n"
		}
		err := errors.New("No command specified.\nAvailable commands:\n" + list + "\n")
		log.Fatal(err)
	}

	command := args[1]

	var param []string

	if len(args) > 2 {
		param = append(args[2:])
	}

	fun := reflect.ValueOf(call)
	in := make([]reflect.Value, len(param) + 2)
	in[0] = reflect.ValueOf(funcs)
	in[1] = reflect.ValueOf(command)
	for k, param := range param {
		if isInt(param) {
			item, _ := strconv.Atoi(param)
			in[k+2] = reflect.ValueOf(item)
		} else if param == "true" || param == "false" {
			item, _ := strconv.ParseBool(param)
			in[k+2] = reflect.ValueOf(item)
		} else {
			in[k+2] = reflect.ValueOf(param)
		}
	}
	fun.Call(in)
}

func AddCommand(command string, function interface{})  {
	funcs[command] = function
}

func call(m map[string]interface{}, name string, params ... interface{}) (result []reflect.Value, err error) {
	if m[name] == nil {
		err = errors.New("Command " + name + " not registered.")
		log.Fatal(err)
	}

	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		log.Fatal(err)
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}