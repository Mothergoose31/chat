func (db *database) getUser(nick string) (Userid, bool) {

	stmt := db.getStatement("getUser", `
		SELECT
			u.userId,
			IF(IFNULL(f.featureId, 0) >= 1, 1, 0) AS protected
		FROM dfl_users AS u
		LEFT JOIN dfl_users_features AS f ON (
			f.userId = u.userId AND
			featureId = (SELECT featureId FROM dfl_features WHERE featureName IN("protected", "admin") LIMIT 1)
		)
		WHERE u.username = ?
	`)
	db.Lock()
	defer stmt.Close()
	defer db.Unlock()

	var uid int32
	var protected bool
	err := stmt.QueryRow(nick).Scan(&uid, &protected)
	if err != nil {
		D("error looking up", nick, err)
		return 0, false
	}
	return Userid(uid), protected
}
