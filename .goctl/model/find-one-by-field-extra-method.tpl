
func (m *default{{.upperStartCamelObject}}Model) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", {{.primaryKeyLeft}}, primary)
}

func (m *default{{.upperStartCamelObject}}Model) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where {{.originalPrimaryField}} = ? limit 1", {{.lowerStartCamelObject}}Rows, m.table )
	return conn.QueryRow(v, query, primary)
}
