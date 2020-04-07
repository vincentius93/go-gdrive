package go_gdrive

import (
	"errors"
	"fmt"
	"google.golang.org/api/sheets/v4"
)

type spreadSheet struct {
	spreadSheetService *sheets.Service
}
type SpreadSheetInterface interface {
	CreateNewSheet(FileTitle string)(*sheets.Spreadsheet,error)
	WriteDataRaw(SpreadSheetId string,Range string, value []interface{})(error)
	UpdateCellStyle(SpreadSheetId string, cellStyle []sheets.Request)(error)
}

func (s *spreadSheet)CreateNewSheet(FileTitle string)(*sheets.Spreadsheet,error){
	spreadSheetOptions := &sheets.Spreadsheet{
		Properties:
			&sheets.SpreadsheetProperties{
				Title:FileTitle,
			},
	}
	spreadSheet,err := s.spreadSheetService.Spreadsheets.Create(spreadSheetOptions).Do()
	return spreadSheet,err
}

func (s *spreadSheet)WriteDataRaw(SpreadSheetId string,Range string, value []interface{})(error){
	var valueRange sheets.ValueRange
	valueRange.Values = append(valueRange.Values, value)

	_,err := s.spreadSheetService.Spreadsheets.Values.Update(SpreadSheetId, Range, &valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		err = errors.New(fmt.Sprintf("nable to retrieve data from sheet: %v\n", err))
		return err
	}
	return err
}

func (s *spreadSheet)UpdateCellStyle(SpreadSheetId string, cellStyle []sheets.Request)(error){
	batchUpdateRequest := sheets.BatchUpdateSpreadsheetRequest{}
	for i:=0;i < len(cellStyle); i++{
		batchUpdateRequest.Requests= append(batchUpdateRequest.Requests,&cellStyle[i])
	}
	if len(cellStyle)== 0{
		return nil
	}
	_,err := s.spreadSheetService.Spreadsheets.BatchUpdate(SpreadSheetId,&batchUpdateRequest).Do()
	if err != nil {
		return err
	}
	s.spreadSheetService.Spreadsheets.BatchUpdate(SpreadSheetId,&batchUpdateRequest)

	return nil
}