package models

import "time"

// Noun represents a Greek noun with all declined forms
type Noun struct {
	ID           int64     `db:"id"`
	English      string    `db:"english"`
	Gender       string    `db:"gender"`
	NominativeSg string    `db:"nominative_sg"`
	GenitiveSg   string    `db:"genitive_sg"`
	AccusativeSg string    `db:"accusative_sg"`
	NominativePl string    `db:"nominative_pl"`
	GenitivePl   string    `db:"genitive_pl"`
	AccusativePl string    `db:"accusative_pl"`
	NomSgArticle string    `db:"nom_sg_article"`
	GenSgArticle string    `db:"gen_sg_article"`
	AccSgArticle string    `db:"acc_sg_article"`
	NomPlArticle string    `db:"nom_pl_article"`
	GenPlArticle string    `db:"gen_pl_article"`
	AccPlArticle string    `db:"acc_pl_article"`
	CreatedAt    time.Time `db:"created_at"`
}
