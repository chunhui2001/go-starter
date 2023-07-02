package googleapi

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/file/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GoogleAPIConf struct {
	Enable          bool     `mapstructure:"GOOGLE_API_ENABLE"`
	CredentialsFile string   `mapstructure:"GOOGLE_API_CREDENTIALS_FILE"`
	TokenFile       string   `mapstructure:"GOOGLE_API_TOKEN_FILE"`
	Scopes          []string `mapstructure:"GOOGLE_API_SCOPES"`
}

var (
	logger           *logrus.Entry
	SCOPES           = []string{}
	CREDENTIALS_FILE = ""
	TOKEN_FILE       = ""
	Client           *http.Client
	FILE_SERVICE     *file.Service
	SHEET_SERVICE    *sheets.Service
	DRIVE_SERVICE    *drive.Service
)

func Init(conf *GoogleAPIConf, log *logrus.Entry) {

	logger = log
	CREDENTIALS_FILE = filepath.Join(utils.RootDir(), conf.CredentialsFile)
	TOKEN_FILE = filepath.Join(utils.RootDir(), conf.TokenFile)
	SCOPES = conf.Scopes

	b, err := os.ReadFile(CREDENTIALS_FILE)

	if err != nil {
		logger.Errorf("GoogleApi-Unable-to-read-client-secret-file: %v", err)
		return
	}

	config, err := google.ConfigFromJSON(b, conf.Scopes...)

	if err != nil {
		logger.Errorf("GoogleApi-Unable-to-parse-client-secret-file-to-config: %v", err)
		return
	}

	Client = getClient(config)

	SHEET_SERVICE, err = sheets.NewService(context.Background(), option.WithHTTPClient(Client))

	if err != nil {
		logger.Errorf("GoogleApi-Unable-to-retrieve-Sheets-client: %v", err)
		return
	}

	FILE_SERVICE, err = file.NewService(context.Background(), option.WithHTTPClient(Client))

	if err != nil {
		logger.Errorf("GoogleApi-Unable-to-retrieve-File-client: %v", err)
		return
	}

	DRIVE_SERVICE, err = drive.NewService(context.Background(), option.WithHTTPClient(Client))

	if err != nil {
		logger.Errorf("GoogleApi-Unable-to-Drive-Http-client: %v", err)
		return
	}

	logger.Infof("GoogleApi-Client-init-CredentialsFile: File=%s", strings.Replace(CREDENTIALS_FILE, utils.RootDir(), "", -1))
	logger.Infof("GoogleApi-Client-init-TokenFile: File=%s", strings.Replace(TOKEN_FILE, utils.RootDir(), "", -1))
	logger.Infof("GoogleApi-Client-init-Scope: Value=%v", SCOPES)

}

// spreadsheetId: 1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms
// readRange: Class Data!A2:E
// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms
func ReadSheet(spreadsheetId string, readRange string) ([][]interface{}, error) {

	// Prints the names and majors of students in a sample spreadsheet:
	resp, err := SHEET_SERVICE.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()

	if err != nil {
		logger.Errorf("GoogleApi-Unable-Read-Sheet: spreadsheetId=%s, readRange=%s, ErrorMessage=%v",
			spreadsheetId, readRange, err)
		return nil, err
	}

	return resp.Values, nil

	// if len(resp.Values) == 0 {
	// 	panic("No Sheet data found.")
	// } else {
	// 	fmt.Println("Name, Major:")
	// 	for _, row := range resp.Values {
	// 		// Print columns A and E, which correspond to indices 0 and 4.
	// 		fmt.Printf("%s, %s\n", row[0], row[4])
	// 	}
	// }

}

// https://docs.google.com/spreadsheets/u/0/
// https://docs.google.com/spreadsheets/d/${SpreadsheetId}
func CreateSheet(title string) (string, error) {

	rb := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title:    title,
			TimeZone: "UTC",
		},
	}

	resp, err := SHEET_SERVICE.Spreadsheets.Create(rb).Context(context.Background()).Do()

	if err != nil {
		logger.Errorf("GoogleApi-Created-Sheet-Failed: %v", err)
		return "", err
	}

	// fmt.Printf("%#v\n", resp)

	return resp.SpreadsheetId, nil

}

// Class Data!A2:E
// writeRange := "A1" // or "sheet1:A1" if you have a different sheet
func WriteToSpreadsheet(spreadsheetId string, writeRange string, values *[][]interface{}) error {

	var vr sheets.ValueRange

	vr.Values = append(vr.Values, *values...)

	var theWriteRange string = writeRange

	if writeRange == "" {
		theWriteRange = "A1"
	}

	_, err := SHEET_SERVICE.Spreadsheets.Values.Update(spreadsheetId, theWriteRange, &vr).ValueInputOption("RAW").Do()

	if err != nil {
		logger.Errorf("GoogleApi-WriteToSpreadsheet-Error: %v", err)
	}

	return err

}

func ClearSheet(spreadsheetId string) error {

	// 清空表格内容
	clearRequest := &sheets.BatchClearValuesRequest{}

	_, err := SHEET_SERVICE.Spreadsheets.Values.BatchClear(spreadsheetId, clearRequest).Do()

	return err

}

// 使用 sed 命令来打印文件的指定行数
// sed -n '1,5p' 文件路径, 这个命令将打印文件 file.txt 中的第 1 行到第 5 行之间的内容。
func ImportCsv(spreadsheetId string, sheetName string, csvFilePath string, separator string) (int, error) {

	fi, err := os.Stat(csvFilePath)

	if err != nil {
		logger.Errorf("GoogleApi-ImportCsv-StatError: %v", err)
		return 0, err
	}

	start0 := time.Now()

	// 读取CSV文件数据
	csvFile, err := os.Open(csvFilePath)

	if err != nil {
		logger.Errorf("GoogleApi-ImportCsv-OpenError: %v", err)
		return 0, err
	}

	defer csvFile.Close()

	logger.Debugf("GoogleApi-ImportCsv-打开文件: 耗时=%s, csvFilePath=%s", time.Since(start0), csvFilePath)

	start1 := time.Now()

	reader := csv.NewReader(csvFile)
	reader.Comma = []rune(separator)[0] // 设置分隔符
	csvData, err := reader.ReadAll()

	if err != nil {
		logger.Fatalf("无法读取CSV文件：%v", err)
		return 0, err
	}

	logger.Debugf("GoogleApi-ImportCsv-读取数据: 耗时=%s, Size=%s, csvFilePath=%s", time.Since(start1), utils.HumanFileSizeInt64(fi.Size()), csvFilePath)

	start2 := time.Now()

	// 转换CSV数据格式
	var values [][]interface{}

	for _, row := range csvData {

		var valueRow []interface{}

		for _, cell := range row {
			valueRow = append(valueRow, cell)
		}

		values = append(values, valueRow)

	}

	logger.Debugf("GoogleApi-ImportCsv-转换数据: 耗时=%s, csvFilePath=%s", time.Since(start2), csvFilePath)

	// 构建请求体
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	start3 := time.Now()

	// 执行导入请求
	_, err = SHEET_SERVICE.Spreadsheets.Values.Update(spreadsheetId, sheetName, valueRange).ValueInputOption("RAW").Do()

	logger.Debugf("GoogleApi-ImportCsv-更新文件: 耗时=%s, 数量=%d, spreadsheetId=%s", time.Since(start3), len(values), spreadsheetId)

	if err != nil {
		logger.Errorf("GoogleApi-ImportCsv-导入异常：spreadsheetId=%s, csvFilePath=%s, Error=%v", spreadsheetId, csvFilePath, err)
		return 0, err
	}

	return len(values), nil

}

// https://developers.google.com/drive/api/v2/reference/revisions/list
func AllRevisions(fileId string) ([]*drive.Revision, error) {

	r, err := DRIVE_SERVICE.Revisions.List(fileId).Do()

	if err != nil {
		logger.Errorf("GoogleApi-AllRevisions-Error: %v", err)
		return nil, err
	}

	return r.Items, nil

}

//	{
//	    "pinned": true,
//	    "publishAuto": true,
//	    "published": true,
//	    "publishedOutsideDomain": true
//	}
func PatchRevision(fileId string, revisionId string, revision *drive.Revision) error {

	r := revision

	if revision == nil {
		r = &drive.Revision{Pinned: true}
	}

	_, err := DRIVE_SERVICE.Revisions.Patch(fileId, revisionId, r).Do()

	if err != nil {
		logger.Errorf("GoogleApi-PatchRevision-Error: %v", err)
		return err
	}

	return nil

}

// AllPermissions fetches all permissions for a given file
func AllPermissions(fileId string) ([]*drive.Permission, error) {

	r, err := DRIVE_SERVICE.Permissions.List(fileId).Do()

	if err != nil {
		logger.Errorf("GoogleApi-AllPermissions-Error: %v", err)
		return nil, err
	}

	for _, p := range r.Items {
		p.Etag = utils.DecodeJsonString(p.Etag)
	}

	return r.Items, nil

}

// 设置文件权限
// https://developers.google.com/drive/api/guides/manage-sharing
// type — The type identifies the scope of the permission (user, group, domain, or anyone).
//        A permission with type=user applies to a specific user whereas a permission with type=domain applies to everyone in a specific domain.
// role — The role field identifies the operations that the type can perform.
//        For example, a permission with type=user and role=reader grants a specific user read-only access to the file or folder.
//        Or, a permission with type=domain and role=commenter lets everyone in the domain add comments to a file.
//        For a complete list of roles and the operations permitted by each, refer to Roles.

// InsertPermission adds a permission to the given file with value type and role
func InsertPermission(fileId string, value string, permType string, role string) error {

	p := &drive.Permission{Type: permType, Role: role}

	if value != "" {
		p = &drive.Permission{Value: value, Type: permType, Role: role}
	}

	_, err := DRIVE_SERVICE.Permissions.Insert(fileId, p).Do()

	if err != nil {
		logger.Errorf("GoogleApi-InsertPermission-Error: %v", err)
		return err
	}

	return nil

}

// 多个email地址用逗号分隔
func ShardWithReader(fileId string, userEmails string) (*drive.Permission, error) {

	permission := &drive.Permission{
		Type:         "user",
		Role:         "reader", // 为目标用户授予 "reader" 角色，以允许他们查看和复制文件，但不能修改原始文件。
		EmailAddress: userEmails,
		Value:        userEmails,
	}

	p, err := DRIVE_SERVICE.Permissions.Insert(fileId, permission).Do()

	// createdPermission, err := DRIVE_SERVICE.Permissions.Create(fileId, permission).Do()

	if err != nil {
		logger.Errorf("GoogleApi-ShardWithReader-Error: fileId=%s, %v", fileId, err)
		return nil, err
	}

	// 获取分享链接
	// shareURL := createdPermission.Link

	p.Etag = utils.DecodeJsonString(p.Etag)

	return p, nil

}

type RefreshToken struct {
	AccessToken string        `json:"access_token,omitempty"`
	ExpiresIn   time.Duration `json:"expires_in,omitempty"`
	Scope       string        `json:"scope,omitempty"`
	TokenType   string        `json:"token_type,omitempty"`
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {

	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first time.
	tok, err := tokenFromFile()

	if err != nil {
		tok = getTokenFromWeb(config)
		if tok == nil {
			logger.Errorf("GoogleApi-Could-not-Get-token: CredentialsFile=%s", CREDENTIALS_FILE)
			return nil
		}
		saveToken(tok)
	}

	if tok.Expiry.Before(time.Now()) {
		logger.Infof("GoogleApi-need-to-renew-new-access-token: %s", tok.Expiry)
		tok = RenewToken(config, tok)
	}

	return config.Client(context.Background(), tok)

}

func RenewToken(config *oauth2.Config, tok *oauth2.Token) *oauth2.Token {

	urlValue := url.Values{"client_id": {config.ClientID}, "client_secret": {config.ClientSecret}, "refresh_token": {tok.RefreshToken}, "grant_type": {"refresh_token"}}

	resp, err := http.PostForm("https://www.googleapis.com/oauth2/v3/token", urlValue)

	if err != nil {
		logger.Errorf("GoogleApi-Error-when-renew-token: %v", err)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		logger.Errorf("GoogleApi-Error-when-renew-token: %v", err)
		return nil
	}

	var refresh_token RefreshToken

	json.Unmarshal([]byte(body), &refresh_token)

	// logger.Infof("GoogleApi-Refresh-Token-Successful: ClientId=%s, NewToken=%+v", config.ClientID, refresh_token)
	logger.Infof("GoogleApi-Refresh-Token-Successful: ClientId=%s, ExpiresIn=%s, TokenType=%s", config.ClientID, refresh_token.ExpiresIn, refresh_token.TokenType)

	then := time.Now()
	then = then.Add(time.Duration(refresh_token.ExpiresIn) * time.Second)

	tok.Expiry = then
	tok.AccessToken = refresh_token.AccessToken

	saveToken(tok)

	return tok

}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {

	// authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	logger.Infof("GoogleApi-Go-to-the-following-link-in-your-browser-then-type-the-authorization-code: \n%v\n", authURL)

	var authCode string

	if _, err := fmt.Scan(&authCode); err != nil {
		logger.Errorf("GoogleApi-Unable-to-read-authorization-code: %v", err)
		return nil
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		logger.Errorf("GoogleApi-Unable-to-retrieve-token-from-web: %v", err)
		return nil
	}

	return tok
}

// Retrieves a token from a local file.
func tokenFromFile() (*oauth2.Token, error) {

	f, err := os.Open(TOKEN_FILE)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err

}

// Saves a token to a file path.
func saveToken(token *oauth2.Token) {

	f, err := os.OpenFile(TOKEN_FILE, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		logger.Errorf("GoogleApi-Unable-to-cache-oauth-token: %v", err)
		return
	}

	defer f.Close()

	json.NewEncoder(f).Encode(token)

	logger.Infof("GoogleApi-Saving-credential-file-to: %v", TOKEN_FILE)

}
