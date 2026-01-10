package databases

type MySQLConnector struct{}

func (c MySQLConnector) Connect(dsn string) error {
	_ = dsn
	return nil
}
