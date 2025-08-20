package cat

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(c *Cat) error {
	return r.db.QueryRow(
		`INSERT INTO cats (name, experience, breed, salary) VALUES ($1, $2, $3, $4) RETURNING id`,
		c.Name, c.Experience, c.Breed, c.Salary,
	).Scan(&c.Id)
}

func (r *Repository) GetByID(id int) (*Cat, error) {
	var c Cat
	err := r.db.QueryRow(
		`SELECT id, name, experience, breed, salary FROM cats WHERE id=$1`,
		id,
	).Scan(&c.Id, &c.Name, &c.Experience, &c.Breed, &c.Salary)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) GetAll() ([]Cat, error) {
	rows, err := r.db.Query(`SELECT id, name, experience, breed, salary FROM cats`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []Cat
	for rows.Next() {
		var c Cat
		if err := rows.Scan(&c.Id, &c.Name, &c.Experience, &c.Breed, &c.Salary); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

func (r *Repository) Update(c *Cat) error {
	_, err := r.db.Exec(
		`UPDATE cats SET name=$1, experience=$2, breed=$3, salary=$4 WHERE id=$5`,
		c.Name, c.Experience, c.Breed, c.Salary, c.Id,
	)
	return err
}

func (r *Repository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM cats WHERE id=$1`, id)
	return err
}
