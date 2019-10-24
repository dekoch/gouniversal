package printtable

import (
	"fmt"
	"strings"
)

type Table struct {
	colCnt   int
	colWidth []int
	rows     []row
}

type row struct {
	col []string
}

func (ta *Table) AddRow(columns []string) {

	if len(columns) > ta.colCnt {
		ta.colCnt = len(columns)
	}

	for i := range columns {
		columns[i] = strings.Trim(columns[i], " ")
	}

	var n row
	n.col = columns
	ta.rows = append(ta.rows, n)
}

func (ta *Table) Print() {

	ta.setColWidth()

	for i := range ta.rows {

		for ii := range ta.rows[i].col {

			s := ta.rows[i].col[ii]
			l := len(s)

			if l < ta.colWidth[ii] {

				for n := l; n < ta.colWidth[ii]; n++ {

					s += " "
				}
			}

			fmt.Print(s)
		}

		fmt.Println()
	}
}

func (ta *Table) Clear() {

	ta.rows = []row{}
}

func (ta *Table) setColWidth() {

	ta.colWidth = make([]int, ta.colCnt)

	for i := range ta.rows {

		for ii := range ta.rows[i].col {

			if ii > ta.colCnt {
				continue
			}

			n := len(ta.rows[i].col[ii])

			if ta.colWidth[ii] < n {
				ta.colWidth[ii] = n + 3
			}
		}
	}
}
