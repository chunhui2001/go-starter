package ges

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chunhui2001/go-starter/core/utils"
	"gopkg.in/yaml.v3"
)

type YamlFileItem struct {
	Name      string
	Path      string
	Content   string
	Templates map[string]interface{}
}

var (
	YamlMapBlocks map[string]*YamlFileItem = make(map[string]*YamlFileItem)
)

func InitDSL() {

	dslFolder := filepath.Join(utils.RootDir(), esConf.DslFolder)

	if ok, _ := utils.FileExists(dslFolder); ok {

		files, err := ioutil.ReadDir(dslFolder)

		if err != nil {
			logger.Errorf("Els-InitDSL-Folder-Not-Exists: DslFolder=%s", dslFolder)
		} else {

			for _, file := range files {

				if file.IsDir() {
					continue
				}

				if filepath.Ext(file.Name()) != ".yaml" {
					continue
				}

				b, err := os.ReadFile(filepath.Join(dslFolder, file.Name())) // just pass the file name

				if err != nil {
					logger.Errorf("Els-InitDSL-File-Read-Error: Name=%s", file.Name())
					continue
				}

				templateMapEntry := make(map[string]interface{})

				if err := yaml.Unmarshal(b, &templateMapEntry); err != nil {
					logger.Errorf("Els-InitDSL-File-Read-Error: Name=%s", file.Name())
					continue
				}

				item := &YamlFileItem{
					Name:      file.Name(),
					Path:      dslFolder,
					Content:   string(b),
					Templates: templateMapEntry,
				}

				logger.Infof("Els-InitDSL-File-Loaded: Path=%s/%s, Size=%s", dslFolder, file.Name(), utils.HumanFileSize(float64(file.Size())))

				YamlMapBlocks[item.Name] = item

			}

			if len(YamlMapBlocks) == 0 {
				logger.Warnf("Els-InitDSL-Folder-Empty: DslFolder=%s/*.yml", dslFolder)
			}

		}

	} else {
		logger.Errorf("Els-InitDSL-Folder-Not-Exists: DslFolder=%s", dslFolder)
	}

}

func DSLQuery(filename string, tplname string, templateData map[string]interface{}) (string, error) {

	if YamlMapBlocks[filename] == nil {
		logger.Errorf("Els-InitDSL-Template-File-Not-Exists: filename=%s, tplname=%s", filename, tplname)
		return "", nil
	}

	if len(YamlMapBlocks[filename].Templates) == 0 || YamlMapBlocks[filename].Templates[tplname] == nil {
		logger.Errorf("Els-InitDSL-Template-Blocks-Entry-Not-Found: filename=%s, tplname=%s", filename, tplname)
		return "", nil
	}

	tmpl := utils.ToJsonString(YamlMapBlocks[filename].Templates[tplname])
	t, err := template.New(filename).Parse(tmpl)

	if err != nil {
		logger.Errorf(`Els-Parse-DslHtml-Template-Error: tmpl=%s, ErrorMessage=%s`, tmpl, err.Error())
		return "", err
	}

	// mapData := utils.MapOf("key", "val")

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, templateData); err != nil {
		logger.Errorf(`Els-Exccute-Parse-Template-Error: Tmpl=%s, Data=%s, ErrorMessage=%s`, tmpl, utils.ToJsonString(tplname), err.Error())
		return "", err
	}

	return tpl.String(), nil

}
