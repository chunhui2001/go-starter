package middleware

import (
	"encoding/xml"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

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

		logger.Infof(`Go-Rewrite-Xml-Loaded: FilePath=%s`, rewirteXmlFilePath)

	}

}

// ### Go by Example: Regular Expressions
// https://gobyexample.com/regular-expressions
// https://go.dev/play/
func Urlwriter() gin.HandlerFunc {

	return func(c *gin.Context) {

		currentPath := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		redirectTo := ""
		fromRegex := ""

		for _, rule := range UrlRewriter.Rules {

			fromRegex = rule.From
			r, _ := regexp.Compile(rule.From)
			match2 := r.MatchString(currentPath)

			if match2 {

				// ^/col/list/(\w+)/(\w+)\.html
				// /col/list/article/1.html
				// [[/col/list/article/1.html article 1]]
				allMatchs := r.FindAllStringSubmatch(currentPath, -1)[0]
				redirectTo = rule.To

				for i := range allMatchs {

					if i == 0 {
						continue
					}

					var replacer = strings.NewReplacer(
						"$"+strconv.Itoa(i), allMatchs[i],
					)

					redirectTo = replacer.Replace(redirectTo)

				}

				break

			}

		}

		if redirectTo != "" {

			if raw != "" {
				currentPath = currentPath + "?" + raw
				redirectTo = redirectTo + "?" + raw
			}

			logger.Infof(`Urlwriter: Regexp=%s, Path=%s, RedirectTo=%s`, fromRegex, currentPath, redirectTo)
			c.Redirect(http.StatusMovedPermanently, redirectTo)
			c.AbortWithStatus(301)

		}

	}

}
