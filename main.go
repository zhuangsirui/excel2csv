package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tealeg/xlsx"
	"github.com/urfave/cli"
)

var (
	output    string
	trimFloat bool
)

func main() {
	app := cli.NewApp()
	app.Name = "excel2csv"
	app.Usage = "convert excel each sheets to a single csv"
	app.Version = "0.0.2"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "output, o",
			Value:       ".",
			Usage:       "target directory for output csv",
			Destination: &output,
		},
		cli.BoolFlag{
			Name:        "trim-float",
			Usage:       "try to parse string like 1.10000000000001 to 1.1",
			Destination: &trimFloat,
		},
	}
	app.Action = convert
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func convert(c *cli.Context) error {
	for _, arg := range c.Args() {
		if err := convertExcelTo(arg); err != nil {
			return err
		}
	}
	return nil
}

func convertExcelTo(filePath string) error {
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("open xlsx file %s failed: %s", filePath, err))
	}
	var errors []error
	for _, sheet := range xlFile.Sheets {
		if err := convertSheetTo(sheet); err != nil {
			errors = append(errors, err)
		}
	}
	for _, err := range errors {
		fmt.Printf("convert %s has failed: %s", filePath, err)
	}
	return nil
}

func convertSheetTo(sheet *xlsx.Sheet) error {
	csvName := sheet.Name + ".csv"
	csvPath := filepath.Join(output, csvName)
	f, err := os.OpenFile(csvPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	log.Printf("convert %s into %s", sheet.Name, csvPath)
	w := csv.NewWriter(f)
	for _, row := range sheet.Rows {
		var record []string
		for _, cell := range row.Cells {
			if trimFloat {
				record = append(record, roundFloat(cell.String()))
			} else {
				record = append(record, cell.String())
			}
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}
