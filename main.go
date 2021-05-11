package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/0LuigiCode0/go-gen/tmp"

	"github.com/0LuigiCode0/logger"
)

type Config struct {
	ModuleName string                     `json:"module_name"`
	GoVersion  float32                    `json:"go_version"`
	DBS        map[string]tmp.DBType      `json:"dbs"`
	Handlers   map[string]tmp.HandlerType `json:"handlers"`
	WorkDir    string                     `json:"work_dir"`
}

var log = logger.InitLogger("")
var fileConfig string

func main() {
	flag.StringVar(&fileConfig, "file", "", "generation setup file")
	flag.Parse()
	conf, err := parseConfig(fileConfig)
	if err != nil {
		log.Fatalf("cannot parse tmp.Config: %v", err)
	}
	if len(conf.Handlers) == 0 {
		log.Fatal("handlers si null")
	}
	if conf.GoVersion <= 0 {
		log.Fatal("version si null")
	}
	if conf.ModuleName == "" {
		log.Fatal("module si null")
	}

	if err = conf.bMain(); err != nil {
		log.Fatalf("cannot create %v: %v", tmp.FileMain, err)
	}
	if err = conf.bServer(); err != nil {
		log.Fatalf("cannot create %v: %v", tmp.FileServer, err)
	}
	if len(conf.DBS) > 0 {
		if err = conf.bDatabase(); err != nil {
			log.Fatalf("cannot create %v: %v", tmp.FileDatabase, err)
		}
		if err = conf.bStore(); err != nil {
			log.Fatalf("cannot create %v: %v", tmp.DirStore, err)
		}
	}
	if err = conf.bHub(); err != nil {
		log.Fatalf("cannot create %v: %v", tmp.DirHub, err)
	}
	if err = conf.bHandlers(); err != nil {
		log.Fatalf("cannot create %v: %v", tmp.DirHandlers, err)
	}
	if err = conf.bHelper(); err != nil {
		log.Fatalf("cannot create %v: %v", tmp.DirHelper, err)
	}
	if err = conf.bUtils(); err != nil {
		log.Fatalf("cannot create %v: %v", "utils", err)
	}
}

//ParseConfig парсит конфиг
func parseConfig(configName string) (*Config, error) {
	_, err := os.Stat(configName)
	if err != nil {
		return nil, fmt.Errorf("file not found: %v", configName)
	}
	file, err := os.Open(configName)
	if err != nil {
		return nil, fmt.Errorf("cannot open file : %v", err)
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read body is invalid : %v", err)
	}
	data := new(Config)
	err = json.Unmarshal(buf, data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal is invalid : %v", err)
	}
	data.ModuleName = strings.TrimSpace(data.ModuleName)
	data.WorkDir = strings.TrimSpace(data.WorkDir)
	return data, err
}

func (c *Config) isOneTCP() bool {
	for _, v := range c.Handlers {
		if v == tmp.TCP {
			return true
		}
	}
	return false
}
func (c *Config) isOneMQTT() bool {
	for _, v := range c.Handlers {
		if v == tmp.MQTT {
			return true
		}
	}
	return false
}
func (c *Config) isOneWS() bool {
	for _, v := range c.Handlers {
		if v == tmp.WS {
			return true
		}
	}
	return false
}
func (c *Config) isOnePostgres() bool {
	for _, v := range c.DBS {
		if v == tmp.Postgres {
			return true
		}
	}
	return false
}
func (c *Config) isOneMongo() bool {
	for _, v := range c.DBS {
		if v == tmp.Mongodb {
			return true
		}
	}
	return false
}

func (c *Config) bMain() error {
	pathDir := filepath.Join(c.WorkDir, tmp.DirCore, tmp.DirCmd)
	pathFile := filepath.Join(pathDir, tmp.FileMain)
	os.MkdirAll(pathDir, 0777)

	f, err := os.OpenFile(pathFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileMain, err)
	}
	defer f.Close()

	t, err := template.New("main").Parse(tmp.MainTmp)
	if err != nil {
		return err
	}
	if err = t.Execute(f, c); err != nil {
		return fmt.Errorf("file %v cannot write: %v", tmp.FileMain, err)
	}
	return nil
}

func (c *Config) bServer() error {
	pathDir := filepath.Join(c.WorkDir, tmp.DirCore)
	pathFile := filepath.Join(pathDir, tmp.FileServer)
	os.MkdirAll(pathDir, 0777)

	f, err := os.OpenFile(pathFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileServer, err)
	}
	defer f.Close()

	t, err := template.New("server").Parse(tmp.ServerTmp)
	if err != nil {
		return err
	}
	if err = t.Execute(f, c); err != nil {
		return fmt.Errorf("file %v cannot write: %v", tmp.FileServer, err)
	}
	return nil
}

func (c *Config) bDatabase() error {
	pathDir := filepath.Join(c.WorkDir, tmp.DirCore, tmp.DirDatabase)
	pathFile := filepath.Join(pathDir, tmp.FileDatabase)
	os.MkdirAll(pathDir, 0777)

	f, err := os.OpenFile(pathFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileDatabase, err)
	}
	defer f.Close()

	fmap := template.FuncMap{
		"title":         strings.Title,
		"isOnePostgres": c.isOnePostgres,
		"isOneMongo":    c.isOneMongo,
	}

	t, err := template.New("database").Funcs(fmap).Parse(tmp.DatabaseTmp)
	if err != nil {
		return err
	}
	if err = t.Execute(f, c); err != nil {
		return fmt.Errorf("file %v cannot write: %v", tmp.FileDatabase, err)
	}
	return nil
}

func (c *Config) bStore() error {
	pathDir := filepath.Join(c.WorkDir, tmp.DirStore)
	os.MkdirAll(pathDir, 0777)
	for i := range c.DBS {
		pathDirStore := filepath.Join(pathDir, i)
		pathFileStore := filepath.Join(pathDirStore, tmp.FileStore)
		os.MkdirAll(pathDirStore, 0777)

		f, err := os.OpenFile(pathFileStore, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
		if err != nil {
			return fmt.Errorf("file %v cannot open: %v", tmp.FileStore, err)
		}
		defer f.Close()

		t, err := template.New("store").Parse(tmp.StoreTmp)
		if err != nil {
			return err
		}
		if err = t.Execute(f, i); err != nil {
			return fmt.Errorf("file %v cannot write: %v", tmp.FileStore, err)
		}
	}
	return nil
}

func (c *Config) bHub() error {
	pathDir := filepath.Join(c.WorkDir, tmp.DirHub)
	pathDirHelper := filepath.Join(pathDir, tmp.DirHelper)
	pathFileHub := filepath.Join(pathDir, tmp.FileHub)
	pathFileHelper := filepath.Join(pathDirHelper, tmp.FileHelper)
	os.MkdirAll(pathDir, 0777)
	os.MkdirAll(pathDirHelper, 0777)

	fhub, err := os.OpenFile(pathFileHub, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileHub, err)
	}
	defer fhub.Close()
	fhelp, err := os.OpenFile(pathFileHelper, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileHelper, err)
	}
	defer fhelp.Close()

	fmap := template.FuncMap{
		"title":    strings.Title,
		"isOneTCP": c.isOneTCP,
	}

	t, err := template.New("hub").Funcs(fmap).Parse(tmp.HubTmp)
	if err != nil {
		return err
	}
	if err = t.Execute(fhub, c); err != nil {
		return fmt.Errorf("file %v cannot write: %v", tmp.FileHub, err)
	}
	t, err = template.New("helper").Funcs(fmap).Parse(tmp.HubHelperTmp)
	if err != nil {
		return err
	}
	if err = t.Execute(fhelp, c); err != nil {
		return fmt.Errorf("file %v cannot write: %v", tmp.FileHelper, err)
	}
	return nil
}

func (c *Config) bHandlers() error {
	pathDir := filepath.Join(c.WorkDir, tmp.DirHandlers)
	os.MkdirAll(pathDir, 0777)
	for i, v := range c.Handlers {
		pathDirHandler := filepath.Join(pathDir, i)
		pathDiHelper := filepath.Join(pathDirHandler, tmp.DirHelper)
		pathFileHandler := filepath.Join(pathDirHandler, tmp.FileHandler)
		pathFileMiddleware := filepath.Join(pathDirHandler, tmp.FileHubMiddleware)
		pathFileHelper := filepath.Join(pathDiHelper, tmp.FileHelper)
		os.MkdirAll(pathDirHandler, 0777)
		os.MkdirAll(pathDiHelper, 0777)

		fhand, err := os.OpenFile(pathFileHandler, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
		if err != nil {
			return fmt.Errorf("file %v cannot open: %v", tmp.FileHandler, err)
		}
		defer fhand.Close()
		fmidl, err := os.OpenFile(pathFileMiddleware, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
		if err != nil {
			return fmt.Errorf("file %v cannot open: %v", tmp.FileHubMiddleware, err)
		}
		defer fmidl.Close()
		fhelp, err := os.OpenFile(pathFileHelper, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
		if err != nil {
			return fmt.Errorf("file %v cannot open: %v", tmp.FileHelper, err)
		}
		defer fhelp.Close()

		var thub, tmiddl, thelp *template.Template
		switch v {
		case tmp.TCP:
			thub, err = template.New("handler").Parse(tmp.HandlerTCPTmp)
			if err != nil {
				return err
			}
			tmiddl, err = template.New("middleware").Parse(tmp.MiddlewareTCPTmp)
			if err != nil {
				return err
			}
		case tmp.MQTT:
			thub, err = template.New("handler").Parse(tmp.HandlerMQTTTmp)
			if err != nil {
				return err
			}
			tmiddl, err = template.New("middleware").Parse(tmp.MiddleWareMQTTTmp)
			if err != nil {
				return err
			}
		case tmp.WS:
			thub, err = template.New("handler").Parse(tmp.HandlerWSTmp)
			if err != nil {
				return err
			}
			tmiddl, err = template.New("middleware").Parse(tmp.MiddleWareWSTmp)
			if err != nil {
				return err
			}
		}
		thelp, err = template.New("helper").Parse(string(tmp.HandlerHelperTmp))
		if err != nil {
			return err
		}

		if err = thub.Execute(fhand, []interface{}{i, c.ModuleName}); err != nil {
			return fmt.Errorf("file %v cannot write: %v", tmp.FileHandler, err)
		}
		if err = tmiddl.Execute(fmidl, []interface{}{i, c.ModuleName}); err != nil {
			return fmt.Errorf("file %v cannot write: %v", tmp.FileHubMiddleware, err)
		}
		if err = thelp.Execute(fhelp, []interface{}{i, v}); err != nil {
			return fmt.Errorf("file %v cannot write: %v", tmp.FileHelper, err)
		}
	}
	return nil
}

func (c *Config) bHelper() error {
	pathDir := filepath.Join(c.WorkDir, tmp.DirHelper)
	pathFileFunc := filepath.Join(pathDir, tmp.FileFunctions)
	pathFileModel := filepath.Join(pathDir, tmp.FileModel)
	os.MkdirAll(pathDir, 0777)

	ff, err := os.OpenFile(pathFileFunc, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileFunctions, err)
	}
	defer ff.Close()
	fm, err := os.OpenFile(pathFileModel, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileModel, err)
	}
	defer fm.Close()

	if _, err = ff.WriteString(tmp.HelperFuncTmp); err != nil {
		return fmt.Errorf("file %v cannot write: %v", tmp.FileFunctions, err)
	}
	if _, err = fm.WriteString(tmp.HelperModelTmp); err != nil {
		return fmt.Errorf("file %v cannot write: %v", tmp.FileModel, err)
	}
	return nil
}

func (c *Config) bUtils() error {
	pathDir := filepath.Join(c.WorkDir, tmp.DirSource)
	pathDirConf := filepath.Join(pathDir, tmp.DirConfigs)
	pathDirUplo := filepath.Join(pathDir, tmp.DirUploads)
	pathFileConf := filepath.Join(pathDirConf, tmp.FileConfigServer)
	pathFileDockerfile := filepath.Join(c.WorkDir, tmp.FileDocker)
	pathFileComposeLocal := filepath.Join(c.WorkDir, tmp.FileComposeLocal)
	pathFileComposeBuild := filepath.Join(c.WorkDir, tmp.FileComposeBuild)
	pathFileMod := filepath.Join(c.WorkDir, tmp.FileMod)
	pathFileSum := filepath.Join(c.WorkDir, tmp.FileSum)
	os.MkdirAll(pathDir, 0777)
	os.MkdirAll(pathDirConf, 0777)
	os.MkdirAll(pathDirUplo, 0777)

	fc, err := os.OpenFile(pathFileConf, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileConfigServer, err)
	}
	defer fc.Close()
	fd, err := os.OpenFile(pathFileDockerfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileDocker, err)
	}
	defer fd.Close()
	fl, err := os.OpenFile(pathFileComposeLocal, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileComposeLocal, err)
	}
	defer fl.Close()
	fb, err := os.OpenFile(pathFileComposeBuild, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileComposeBuild, err)
	}
	defer fb.Close()
	fm, err := os.OpenFile(pathFileMod, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileMod, err)
	}
	defer fm.Close()
	fs, err := os.OpenFile(pathFileSum, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("file %v cannot open: %v", tmp.FileSum, err)
	}
	defer fs.Close()

	t, err := template.New("config").Parse(tmp.ConfigTmp)
	if err != nil {
		return err
	}
	if err = t.Execute(fc, c); err != nil {
		return fmt.Errorf("file %v cannot write: %v", tmp.FileConfigServer, err)
	}
	t, err = template.New("docker").Parse(tmp.DockerTmp)
	if err != nil {
		return err
	}
	if err = t.Execute(fd, c); err != nil {
		return fmt.Errorf("file %v cannot write: %v", tmp.FileDocker, err)
	}
	t, err = template.New("build").Parse(tmp.ComposeBuildTmp)
	if err != nil {
		return err
	}
	if err = t.Execute(fb, c); err != nil {
		return fmt.Errorf("file %v cannot write: %v", tmp.FileComposeBuild, err)
	}
	t, err = template.New("local").Parse(tmp.ComposeLocalTmp)
	if err != nil {
		return err
	}
	if err = t.Execute(fl, c); err != nil {
		return fmt.Errorf("file %v cannot write: %v", tmp.FileComposeLocal, err)
	}
	t, err = template.New("mod").Funcs(template.FuncMap{
		"isOneTCP":      c.isOneTCP,
		"isOneWS":       c.isOneWS,
		"isOneMQTT":     c.isOneMQTT,
		"isOnePostgres": c.isOnePostgres,
		"isOneMongo":    c.isOneMongo,
	}).Parse(tmp.ModTmp)
	if err != nil {
		return err
	}
	if err = t.Execute(fm, c); err != nil {
		return fmt.Errorf("file %v cannot write: %v", tmp.FileMod, err)
	}
	return nil
}
