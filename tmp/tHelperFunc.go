package tmp

const HelperFuncTmp = `package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/0LuigiCode0/logger"
)

func ParseConfig() (*Config, error) {
	_, err := os.Stat(ConfigDir + ConfigFiel)
	if err != nil {
		return nil, fmt.Errorf(KeyErrorNotFound+": file: %v", ConfigDir+ConfigFiel)
	}
	file, err := os.Open(ConfigDir + ConfigFiel)
	if err != nil {
		return nil, fmt.Errorf(KeyErrorOpen+": file: %v", err)
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf(KeyErrorRead+": body: %v", err)
	}
	data := new(Config)
	err = json.Unmarshal(buf, data)
	if err != nil {
		return nil, fmt.Errorf(KeyErrorParse+": json: %v", err)
	}

	return data, err
}

func InitCtx() {
	Ctx, CloseCtx = context.WithCancel(context.Background())
}
func InitLogger() {
	Log = logger.InitLogger("")
}

func Dispatch(f interface{}, args ...interface{}) {
	ff := reflect.ValueOf(f)
	if ff.Kind() == reflect.Func {
		in := make([]reflect.Value, ff.Type().NumIn())
		for i, arg := range args {
			v := reflect.ValueOf(arg)
			if v.Type().ConvertibleTo(ff.Type().In(i)) {
				in[i] = v.Convert(ff.Type().In(i))
			} else {
				Log.Warningf("parameter: %v, expected %v got %v", i+1, ff.Type().In(i), v.Type())
				return
			}
		}

		Wg.Add(1)
		go func() {
			ff.Call(in)
			Wg.Done()
		}()
	}
}`
