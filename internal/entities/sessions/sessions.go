package sessions

type Sessions struct {
	UserID           string `db:"user_id" json:"user_id"`
	RefreshTokenHash []byte `db:"refresh_token_hash" json:"refresh_token_hash"`
	UserAgent        string `db:"user_agent" json:"user_agent"`
	IPAddr           string `db:"ip_addr" json:"ip_addr"`
}
