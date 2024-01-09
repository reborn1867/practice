package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type SortedStatsList []stats

type stats struct {
	name            string
	referenceCount  int
	downloadedCount int
}

func (s SortedStatsList) Len() int {
	return len(s)
}
func (s SortedStatsList) Less(i, j int) bool {
	return s[i].referenceCount > s[j].referenceCount
}
func (s SortedStatsList) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

func main() {
	f, err := excelize.OpenFile("stats.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get value from cell by given worksheet name and cell reference.
	cell, err := f.GetCellValue("Sheet2", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet2")
	if err != nil {
		fmt.Println(err)
		return
	}
	newRows := [][]string{}
	for _, row := range rows {
		if strings.Contains(row[0], ";") {
			names := strings.Split(row[0], ";")
			for _, name := range names {
				if name != "" {
					newRows = append(newRows, append([]string{name}, row[1:]...))
				}
			}
		} else {
			newRows = append(newRows, row)
		}
	}

	m := map[string]stats{}
	for _, row := range newRows {
		newRCount, _ := strconv.Atoi(row[1])
		newDCount, _ := strconv.Atoi(row[2])
		s, ok := m[row[0]]
		if ok {
			s.referenceCount += newRCount
			s.downloadedCount += newDCount
			m[row[0]] = s
		} else {
			m[row[0]] = stats{
				name:            row[0],
				referenceCount:  newRCount,
				downloadedCount: newDCount,
			}
		}
	}

	statsList := []stats{}
	for _, stats := range m {
		statsList = append(statsList, stats)
	}

	sort.Sort(SortedStatsList(statsList))

	sortedRows := [][]string{}
	for _, stats := range statsList {
		sortedRows = append(sortedRows, []string{stats.name, fmt.Sprint(stats.referenceCount), fmt.Sprint(stats.downloadedCount)})
	}

	for _, row := range sortedRows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}

	updatedFile := excelize.NewFile()
	defer func() {
		if err := updatedFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	updatedFile.SetCellValue("Sheet1", "A1", "姓名")
	updatedFile.SetCellValue("Sheet1", "B1", "被引")
	updatedFile.SetCellValue("Sheet1", "C1", "下载")
	for i, row := range sortedRows {
		updatedFile.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), row[0])
		updatedFile.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), row[1])
		updatedFile.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), row[2])
	}

	if err := updatedFile.SaveAs("sorted.xlsx"); err != nil {
		fmt.Println(err)
	}
}
