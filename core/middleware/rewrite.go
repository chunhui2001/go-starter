package middleware

import (
	"encoding/xml"
	"net/http"
	"path/filepath"

	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/gin-gonic/gin"
)

type UrlRewrite struct {
	Root  xml.Name `xml:"urlrewrite"`
	Rules []rule   `xml:"rule"`
}

type rule struct {
	From string `xml:"from"`
	To   string `xml:"to"`
}

var (
	logger      = config.Log
	UrlRewriter = UrlRewrite{}
)

// curl -sL https://git.io/fN4Pq | zek -e -p -c
// https://astaxie.gitbooks.io/build-web-application-with-golang/content/ja/07.1.html
// curl -L http://localhost:8080/google
func init() {

	rewirteXmlFilePath := filepath.Join(utils.RootDir(), "resources", "urlrewrite.xml")

	if ok, _ := utils.FileExists(rewirteXmlFilePath); ok {

		b, err := utils.ReadFile(rewirteXmlFilePath) // just pass the file name

		if err != nil {
			logger.Errorf("Go-Rewrite-Xml-Read-Error: FilePath=%s, ErrorMessage=%s", rewirteXmlFilePath, err.Error())
		}

		err = xml.Unmarshal(b, &UrlRewriter)

		if err != nil {
			logger.Errorf("Go-Rewrite-Xml-Parse-Error: FilePath=%s, ErrorMessage=%s", rewirteXmlFilePath, err.Error())
			return
		}

		logger.Infof(`Go-Rewrite-Xml-Loaded: FilePath=%s, Content=%s`, rewirteXmlFilePath, utils.ToJsonString(UrlRewriter))

	}

}

func Urlwriter() gin.HandlerFunc {

	return func(c *gin.Context) {

		currentPath := c.Request.URL.Path
		redirectTo := ""

		for _, rule := range UrlRewriter.Rules {
			if rule.From == currentPath {
				redirectTo = rule.To
				break
			}
		}

		if redirectTo != "" {
			logger.Infof(`Urlwriter: Path=%s, RedirectTo=%s`, currentPath, redirectTo)
			c.Redirect(http.StatusMovedPermanently, redirectTo)
			c.AbortWithStatus(301)
		}

	}

}
