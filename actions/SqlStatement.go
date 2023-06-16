package actions

import (
	"net/http"
	"strings"

	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/gin-gonic/gin"
)

type SqlStatement struct {
	Preparing  string `form:"preparing"`
	Parameters string `form:"parameters"`
}

func RawSql(c *gin.Context) {

	stmt := &SqlStatement{}

	if err := c.Bind(stmt); err != nil {
		c.JSON(http.StatusOK, (&R{Error: err}).Fail(400))
		return
	}

	sql := stmt.Preparing
	params := strings.Split(stmt.Parameters, ", ")

	var rawSql strings.Builder
	var placeholderIndex uint
	var regEx = `(?P<Val>[\w\W]+)(?P<Typ>\([a-zA-Z]+\))`

	match1 := utils.MatchesGroup(regEx, "BTC-USDT(String)")
	logger.Infof(`BTC-USDT(String): matched=%s`, utils.ToJsonString(match1))
	match2 := utils.MatchesGroup(regEx, "BTC-USDT(BigDecimal)")
	logger.Infof(`BTC-USDT(BigDecimal): matched=%s`, utils.ToJsonString(match2))
	match3 := utils.MatchesGroup(regEx, "BTC-USDT(Integer)")
	logger.Infof(`BTC-USDT(Integer): matched=%s`, utils.ToJsonString(match3))

	for _, ch := range sql {
		if ch == '?' {
			// BTC-USDT(String)
			// 0.001(Double)
			// 10(Integer)
			// 8000.000000000000000000(BigDecimal)
			currentValue := params[placeholderIndex]
			match1 := utils.MatchesGroup(regEx, currentValue)
			if match1["Typ"] == `(String)` {
				rawSql.WriteString("'" + match1["Val"] + "'")
			} else {
				rawSql.WriteString(match1["Val"])
			}
			placeholderIndex++
		} else {
			rawSql.WriteString(string(ch))
		}
	}

	logger.Infof(`paramsSize=%d`, len(params))
	logger.Infof(`placeholderCount=%d`, placeholderIndex)

	c.Data(200, "text/plain; charset=utf-8", []byte(rawSql.String()))

}
