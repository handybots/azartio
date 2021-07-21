package database

func (db *DB) LastVote() (vote Vote, _ error) {
	vote.db = db
	const query = `
		SELECT * FROM votes 
		WHERE done=false AND message_id!='' 
		ORDER BY created_at DESC LIMIT 1`
	return vote, db.Get(&vote, query)
}
