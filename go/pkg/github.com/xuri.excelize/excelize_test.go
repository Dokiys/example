package xuri_excelize

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

// TestExcelizeHello 测试Excelize
func TestExcelizeHello(t *testing.T) {
	f := excelize.NewFile()
	// 创建一个工作表
	index := f.NewSheet("Sheet2")
	// 设置单元格的值
	f.SetCellValue("Sheet2", "A2", "Hello work!")
	f.SetCellValue("Sheet1", "B2", 100)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs(PathPrefix + "BookHello_out.xlsx"); err != nil {
		fmt.Println(err)
	}
}

// TestExcelizeReadExcel 测试读取文件
func TestExcelizeReadExcel(t *testing.T) {
	f, err := excelize.OpenFile(PathPrefix + "BookRead.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 获取工作表中指定单元格的值
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}

	f.SaveAs(PathPrefix + "BookRead_out.xlsx")
}

// TestExcelizeChart 添加图表
func TestExcelizeChart(t *testing.T) {
	f := excelize.NewFile()
	categories := map[string]string{"A2": "Small", "A3": "Normal", "A4": "Large", "B1": "Apple", "C1": "Orange", "D1": "Pear"}
	values := map[string]int{"B2": 2, "C2": 3, "D2": 3, "B3": 5, "C3": 2, "D3": 4, "B4": 6, "C4": 7, "D4": 8}

	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}
	for k, v := range values {
		f.SetCellValue("Sheet1", k, v)
	}
	if err := f.AddChart("Sheet1", "E1", `{
        "type": "col3DClustered",
        "series": [
        {
            "name": "Sheet1!$A$2",
            "categories": "Sheet1!$B$1:$D$1",
            "values": "Sheet1!$B$2:$D$2"
        },
        {
            "name": "Sheet1!$A$3",
            "categories": "Sheet1!$B$1:$D$1",
            "values": "Sheet1!$B$3:$D$3"
        },
        {
            "name": "Sheet1!$A$4",
            "categories": "Sheet1!$B$1:$D$1",
            "values": "Sheet1!$B$4:$D$4"
        }],
        "title":
        {
            "name": "Fruit 3D Clustered Column Chart"
        }
    }`); err != nil {
		fmt.Println(err)
		return
	}
	// 根据指定路径保存文件
	if err := f.SaveAs(PathPrefix + "BookChart_out.xlsx"); err != nil {
		fmt.Println(err)
	}
}

// TestExcelizeImage 插入图片
func TestExcelizeImage(t *testing.T) {
	f := excelize.NewFile()
	// 插入图片
	if err := f.AddPicture("Sheet1", "A2", "image.png", ""); err != nil {
		fmt.Println(err)
	}
	// 在工作表中插入图片，并设置图片的缩放比例
	if err := f.AddPicture("Sheet1", "E2", "image.jpeg", `{
        "x_scale": 0.5,
        "y_scale": 0.5
    }`); err != nil {
		fmt.Println(err)
	}
	// 在工作表中插入图片，并设置图片的打印属性
	if err := f.AddPicture("Sheet1", "G2", "image.gif", `{
        "x_offset": 15,
        "y_offset": 10,
        "print_obj": true,
        "lock_aspect_ratio": false,
        "locked": false
    }`); err != nil {
		fmt.Println(err)
	}
	// 保存文件
	if err := f.SaveAs(PathPrefix + "BookImage_out.xlsx"); err != nil {
		fmt.Println(err)
	}
}

// TestExcelizeOneMerge 测试左侧单个合并单元格的插入是否增长
func TestExcelizeOneMerge(t *testing.T) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	// 设置单元格的值
	f.SetCellValue(sheet, "A1", 100)
	f.SetCellValue(sheet, "B2", 1)
	f.MergeCell(sheet, "A1", "A2")
	f.RemoveRow(sheet, 1)

	// f.DuplicateRowTo(sheet, 2, 3)
	// f.RemoveRow(sheet, 2)

	// 根据指定路径保存文件
	f.SaveAs(PathPrefix + "BookOneMerge_out.xlsx")
}

// TestExcelizeExp 测试表达式
func TestExcelizeExp(t *testing.T) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	// 设置单元格的值
	f.SetCellValue(sheet, "A1", 100)
	f.SetCellValue(sheet, "A2", 1)
	f.SetCellFormula(sheet, "A3", "=A1+A2")

	// f.DuplicateRowTo(sheet, 2, 3)
	// f.RemoveRow(sheet, 2)

	// 根据指定路径保存文件
	f.SaveAs(PathPrefix + "BookExp_out.xlsx")
}

// 获取主题色
func getCellBgColor(f *excelize.File, sheet, axix string) string {
	styleID, err := f.GetCellStyle(sheet, axix)
	if err != nil {
		return err.Error()
	}
	fillID := *f.Styles.CellXfs.Xf[styleID].FillID
	fgColor := f.Styles.Fills.Fill[fillID].PatternFill.FgColor
	if fgColor.Theme != nil {
		children := f.Theme.ThemeElements.ClrScheme.Children
		if *fgColor.Theme < 4 {
			dklt := map[int]string{
				0: children[1].SysClr.LastClr,
				1: children[0].SysClr.LastClr,
				2: *children[3].SrgbClr.Val,
				3: *children[2].SrgbClr.Val,
			}
			return strings.TrimPrefix(
				excelize.ThemeColor(dklt[*fgColor.Theme], fgColor.Tint), "FF")
		}
		srgbClr := *children[*fgColor.Theme].SrgbClr.Val
		return strings.TrimPrefix(excelize.ThemeColor(srgbClr, fgColor.Tint), "FF")
	}
	return strings.TrimPrefix(fgColor.RGB, "FF")
}
