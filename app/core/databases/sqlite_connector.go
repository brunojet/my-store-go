package databases

type SQLiteConnector struct{}

func (c SQLiteConnector) Connect(path string) error {
	_ = path
	return nil
}
