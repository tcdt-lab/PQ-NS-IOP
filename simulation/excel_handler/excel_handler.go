package excel_handler

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func WriteToAnExcelFile(fileName string, rows [][]float64) {

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	for idx, row := range rows {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetSheetRow("Sheet1", cell, &row)
	}
	if err := f.SaveAs(fileName + ".xlsx"); err != nil {
		fmt.Println(err)
	}
}
