package sessions

type Sessions struct {
	UserID           string `db:"user_id" json:"user_id"`
	RefreshTokenHash []byte `db:"refresh_token" json:"refresh_token"`
	UserAgent        string `db:"user_agent" json:"user_agent"`
	IssuedIP         string `db:"issued_ip" json:"issued_ip"`
}
