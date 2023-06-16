package actions

import (
	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/chunhui2001/go-starter/core/googleapi"
	"github.com/chunhui2001/go-starter/core/utils"
)

// 列出指定文件的所有权限
func GoogleDocListAllPermissionsRouter(c *gin.Context) {
	fileId := c.Query("fileId")
	permissions, err := googleapi.AllPermissions(fileId)
	c.JSON(http.StatusOK, (&R{Data: permissions, Error: err}).IfErr(400))
}

// 创建一个excel表格
func GoogleDocCreateSheetRouter(c *gin.Context) {
	sheetTitle := c.Query("sheetTitle")
	sheetId, err := googleapi.CreateSheet(sheetTitle)
	c.JSON(http.StatusOK, (&R{Data: sheetId, Error: err}).IfErr(400))
}

// 分享文件
func GoogleDocShardWithReaderRouter(c *gin.Context) {
	fileId := c.Query("fileId")
	emailAddresses := c.Query("emailAddresses")
	p, err := googleapi.ShardWithReader(fileId, emailAddresses)
	c.JSON(http.StatusOK, (&R{Data: p, Error: err}).IfErr(400))
}

// 导入csv
func GoogleDocImportCsvRouter(c *gin.Context) {
	spreadsheetId := c.Query("spreadsheetId")
	csvFilePath := c.Query("csvFilePath")
	_range := c.Query("_range")
	count, err := googleapi.ImportCsv(spreadsheetId, _range, csvFilePath, ",")
	c.JSON(http.StatusOK, (&R{Data: count, Error: err}).IfErr(400))
}

// 导入csv
func GoogleDocClearSheetRouter(c *gin.Context) {
	spreadsheetId := c.Query("spreadsheetId")
	err := googleapi.ClearSheet(spreadsheetId)
	c.JSON(http.StatusOK, (&R{Data: true, Error: err}).IfErr(400))
}

// 导入csv
func GoogleDocCsvReaderRouter(c *gin.Context) {

	csvFilePath := c.Query("csvFilePath")
	windowSize := c.Query("windowSize")

	stream := func(rows [][]string, err error) {
		logger.Infof(`Count=%d, Error=%v`, len(rows), err)
	}

	csvReader := googleapi.CsvReader{
		HasHeader: true,
		FilePath:  csvFilePath,
		Stream:    stream,
	}

	csvReader.Read(utils.StrToInt(windowSize))

	c.JSON(http.StatusOK, (&R{Data: csvReader.TotalCount}).Ok())

}
