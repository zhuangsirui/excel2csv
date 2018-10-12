package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
	"github.com/urfave/cli"
)

var (
	output      string
	trim        bool
	trimFloat   bool
	withBom     bool
	convertBool bool
)

func main() {
	app := cli.NewApp()
	app.Name = "excel2csv"
	app.Usage = "convert excel each sheets to a single csv"
	app.UsageText = "excel2csv [--output DIR] [--trim] [--trim-float] [--with-bom] file [file...]"
	app.Version = "0.0.9"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "output, o",
			Value:       ".",
			Usage:       "target directory for output csv",
			Destination: &output,
		},
		cli.BoolFlag{
			Name:        "trim",
			Usage:       "trim value",
			Destination: &trim,
		},
		cli.BoolFlag{
			Name:        "trim-float",
			Usage:       "try to parse string like 1.10000000000001 to 1.1",
			Destination: &trimFloat,
		},
		cli.BoolFlag{
			Name:        "convert-bool",
			Usage:       "covert 0 1 to \"true\" \"false\"",
			Destination: &convertBool,
		},
		cli.BoolFlag{
			Name:        "with-bom",
			Usage:       "add UTF-8 BOM to csv file",
			Destination: &withBom,
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
	f, err := os.OpenFile(csvPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer f.Close()
	log.Printf("convert %s into %s", sheet.Name, csvPath)
	if withBom {
		_, err := f.Write(bomBytes)
		if err != nil {
			return err
		}
	}
	w := csv.NewWriter(f)
	for _, row := range sheet.Rows {
		var records []string
		for _, cell := range row.Cells {
			record := cell.String()
			if convertBool && cell.Type() == xlsx.CellTypeBool {
				record = boolStringToCharacter(cell.String())
			}
			if trimFloat {
				record = roundFloat(record)
			}
			if trim {
				record = strings.TrimSpace(record)
			}
			records = append(records, record)
		}
		if err := w.Write(records); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}
