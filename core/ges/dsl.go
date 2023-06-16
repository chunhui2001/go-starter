package ges

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// {{$x := join ", " "hello" "world"}}
func funcMaps() template.FuncMap {
	return template.FuncMap{
		"join":       utils.Join, // {{- join "\",\"" .name .name  }} , ["{{- join "\",\"" .name .name  }}"]
		"JoinString": utils.JoinString,
		"ifnull":     utils.IfNull, // {{- ifnull .nullVal 99 }}
	}
}

type YamlFileItem struct {
	Name      string
	Path      string
	Content   string
	Templates map[string]interface{}
}

var (
	YamlMapBlocks map[string]*YamlFileItem = make(map[string]*YamlFileItem)
	mylog         *logrus.Entry
	FuncMaps      = funcMaps()
	PrettyPrint   bool
	DisablePrint  bool
)

func InitDSL(folder string, prettyPrint bool, disablePrint bool, log *logrus.Entry) {

	mylog = log
	dslFolder := filepath.Join(utils.RootDir(), folder)
	PrettyPrint = prettyPrint
	DisablePrint = disablePrint

	if ok, _ := utils.FileExists(dslFolder); ok {

		files, err := ioutil.ReadDir(dslFolder)

		if err != nil {
			mylog.Errorf("Els-InitDSL-Folder-Not-Exists: DslFolder=%s", dslFolder)
		} else {

			for _, file := range files {

				if file.IsDir() {
					continue
				}

				if filepath.Ext(file.Name()) != ".yaml" {
					continue
				}

				b, err := utils.ReadFile(filepath.Join(dslFolder, file.Name())) // just pass the file name

				if err != nil {
					mylog.Errorf("Els-InitDSL-File-Read-Error: Name=%s", file.Name())
					continue
				}

				templateMapEntry := make(map[string]interface{})

				if err := yaml.Unmarshal(b, &templateMapEntry); err != nil {
					mylog.Errorf("Els-InitDSL-File-Read-Error: Name=%s", file.Name())
					continue
				}

				item := &YamlFileItem{
					Name:      file.Name(),
					Path:      dslFolder,
					Content:   string(b),
					Templates: templateMapEntry,
				}

				mylog.Infof("Els-InitDSL-File-Loaded: Path=%s/%s, Size=%s", dslFolder, file.Name(), utils.HumanFileSize(float64(file.Size())))

				YamlMapBlocks[item.Name] = item

			}

			if len(YamlMapBlocks) == 0 {
				mylog.Warnf("Els-InitDSL-Folder-Empty: DslFolder=%s/*.yaml", dslFolder)
			}

		}

	} else {
		mylog.Errorf("Els-InitDSL-Folder-Not-Exists: DslFolder=%s", dslFolder)
	}

}

func DSLQuery(filename string, tplname string, templateData map[string]interface{}) (string, error) {

	if YamlMapBlocks[filename] == nil {
		mylog.Errorf("Els-Template-File-Not-Exists: filename=%s, tplname=%s", filename, tplname)
		return "", nil
	}

	if len(YamlMapBlocks[filename].Templates) == 0 || YamlMapBlocks[filename].Templates[tplname] == nil {
		errorMessage := fmt.Sprintf(`Els-Template-Blocks-Entry-Not-Found: filename=%s, tplname=%s`, filename, tplname)
		mylog.Errorf(errorMessage)
		return "", errors.New(errorMessage)
	}

	tmpl := utils.ToJsonString(YamlMapBlocks[filename].Templates[tplname])
	t, err := template.New(filename).Funcs(FuncMaps).Parse(tmpl)

	if err != nil {
		mylog.Errorf(`Els-Parse-DslHtml-Template-Error: tmpl=%s, ErrorMessage=%s`, tmpl, err.Error())
		return tmpl, err
	}

	// mapData := utils.MapOf("key", "val")

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, templateData); err != nil {
		mylog.Errorf(`Els-Exccute-Parse-Template-Error: Tmpl=%s, Data=%s, ErrorMessage=%s`, tmpl, utils.ToJsonString(tplname), err.Error())
		return tmpl, err
	}

	return prettyprint(filename, tplname, ([]byte)(html.UnescapeString(tpl.String())))

}

func prettyprint(filename string, tplname string, b []byte) (string, error) {

	if !PrettyPrint {
		if !DisablePrint {
			mylog.Infof(`DSLQuery: filename=%s, tplname=%s, rawjson=%s`, filename, tplname, b)
		}
		return string(b), nil
	}

	var out bytes.Buffer
	err := json.Indent(&out, b, "", "")

	if err != nil {
		return string(b), err
	}

	prettyjson := regexp.MustCompile(`\r?\n`).ReplaceAllString(out.String(), "")
	mylog.Infof(`DSLQuery: filename=%s, tplname=%s, prettyjson=%s`, filename, tplname, prettyjson)

	return prettyjson, nil

}
