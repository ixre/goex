package report

var _ IDataExportProvider = new(ExcelProvider)

type ExcelProvider struct {
	csv IDataExportProvider
}

func NewExcelProvider() IDataExportProvider {
	return &ExcelProvider{
		csv: NewCsvProvider(),
	}
}

func (e *ExcelProvider) Export(rows []map[string]interface{}, keys []string, alias []string) (binary []byte) {
	return e.csv.Export(rows, keys, alias)
}
