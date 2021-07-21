package database

func (db *DB) Voter(voteID int, userID int64) (voter Voter, _ error) {
	const query = `SELECT * FROM voters WHERE vote_id=$1 AND user_id=$2`
	return voter, db.Get(&voter, query, voteID, userID)
}
