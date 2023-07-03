package googleapi

import (
	"encoding/csv"
	"os"
)

type CsvReader struct {
	FilePath   string
	Separator  string
	HasHeader  bool
	HeaderRow  []string
	Stream     func([][]string, error)
	TotalCount int
}

func (r *CsvReader) Size() (int64, error) {

	fi, err := os.Stat(r.FilePath)

	if err != nil {
		return 0, err
	}

	return fi.Size(), nil

}

func (r *CsvReader) Read(windowSize int) *CsvReader {

	// 读取CSV文件数据
	csvFile, err := os.Open(r.FilePath)

	if err != nil {
		r.Stream(nil, err)
		return r
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	if r.Separator != "" {
		reader.Comma = []rune(r.Separator)[0] // 设置分隔符, 将 separator 字符串转换为 rune 类型的切片
	}

	if r.HasHeader {

		headerRow, err := reader.Read()

		if err != nil {
			r.Stream(nil, err)
			return r
		}

		r.HeaderRow = headerRow

	}

	for {

		rows, err := r.loopRows(reader, windowSize)

		if err != nil {
			r.Stream(nil, err)
			return r
		}

		if len(rows) > 0 {
			r.TotalCount += len(rows)
			r.Stream(rows, nil)
		} else {
			break
		}

	}

	return nil

}

func (r *CsvReader) loopRows(reader *csv.Reader, size int) ([][]string, error) {

	var rowCount = 0
	var rows [][]string

	// 读取记录
	for {

		if rowCount >= size {
			break
		}

		row, err := reader.Read()

		if err != nil {

			// 检查是否到达文件结尾
			if err.Error() == "EOF" {
				return rows, nil
			}

			// 处理其他错误
			return nil, err

		}

		// 处理记录
		rows = append(rows, row)
		rowCount++

	}

	return rows, nil

}
