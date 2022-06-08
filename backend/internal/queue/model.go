package queue

type Job struct {
	ID         int64  `db:"id"`
	Action     string `db:"action"`
	EntityType string `db:"entityType"`
	Payload    string `db:"payload"`
	Attempts   int64  `db:"attempts"`
	Result     string
	Error      string
}

type RespCount struct {
	Messages int   `json:"messages"`
	IDJwts   []int `json:"idJwts,omitempty"`
}
