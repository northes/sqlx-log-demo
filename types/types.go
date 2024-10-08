package types

type Table struct {
	ID   int    `db:"id" json:"id"`
	Text string `db:"text" json:"text"`
}
