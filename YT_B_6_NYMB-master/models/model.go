package models

/*

Models represent any row retrieved from a database query.
The standard http routes that should be in the following form:

read    → GET    /model
create  → POST   /model
update  → PUT    /model/id
destroy → DELETE /model/id

*/

// GetQuery takes a base query and inserts the given args
func GetQuery(base string, args map[string][]string) string {
	query := base
	ok := false
	if len(args) > 0 {
		query += " WHERE"
		for k, vals := range args {
			if vals[0] != "" {
				ok = true
				query += " " + k + "=" + vals[0] + " AND"
			}
		}
		query = query[0 : len(query)-4]
	}
	if !ok {
		return base
	}
	return query
}
