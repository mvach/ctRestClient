package csv

type CsvData interface {
	Records() [][]string
	Header() []string
}