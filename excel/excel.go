package excel

import (
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

// 通用Excel导出
// sheetName sheet 名，默认是 Sheet1
// outFile 输出文件
// header 表头字段
// rows Excel表每行数据
func exportNomarl(sheetName, outFile string, header []interface{}, rows [][]interface{}) error {
	if len(header) == 0 {
		return errors.New("请设置表头")
	}
	if len(rows) == 0 {
		return errors.New("Excel内容不能为空")
	}
	outf := excelize.NewFile()
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	sheet := outf.GetSheetName(0)
	outf.SetSheetName(sheet, sheetName)
	outf.NewSheet(sheetName)
	sheetWriter, err := outf.NewStreamWriter(sheetName)
	if err != nil {
		return fmt.Errorf("新建sheet 失败:%s", err.Error())
	}
	rowIndex := 1
	cell, _ := excelize.CoordinatesToCellName(1, rowIndex)
	err = sheetWriter.SetRow(cell, header)
	if err != nil {
		return fmt.Errorf("设置表头 失败:%s", err.Error())
	}
	rowIndex++
	for _, row := range rows {
		cell, _ = excelize.CoordinatesToCellName(1, rowIndex)
		err = sheetWriter.SetRow(cell, row)
		if err != nil {
			return fmt.Errorf("写入Excel行数据 失败:%s", err.Error())
		}
		rowIndex++
	}
	err = sheetWriter.Flush()
	if err != nil {
		return fmt.Errorf("写入Excel数据 失败:%s", err.Error())
	}
	err = outf.SaveAs(outFile)
	if err != nil {
		return fmt.Errorf("保存Excel 失败:%s", err.Error())
	}
	return nil
}
