package target

import (
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(missionID int, t *Target) error {
	err := r.db.QueryRow(
		`INSERT INTO targets (mission_id, name, country, notes, complete)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		missionID, t.Name, t.Country, t.Notes, t.Complete,
	).Scan(&t.Id)
	return err
}

func (r *Repository) Update(t *Target) error {
	// 1. перевіряємо стан таргета та місії
	var targetComplete, missionComplete bool
	query := `
		SELECT t.complete, m.complete
		FROM targets t
		JOIN missions m ON t.mission_id = m.id
		WHERE t.id = $1
	`
	row := r.db.QueryRow(query, t.Id)
	if err := row.Scan(&targetComplete, &missionComplete); err != nil {
		return err
	}

	// 2.якщо таргет або місія завершені — забороняємо оновлювати Notes
	if (targetComplete || missionComplete) && t.Notes != "" {
		return errors.New("cannot update notes: mission or target already completed")
	}

	// 3.оновлюємо всі інші поля
	_, err := r.db.Exec(
		`UPDATE targets SET name=$1, country=$2, notes=$3, complete=$4 WHERE id=$5`,
		t.Name, t.Country, t.Notes, t.Complete, t.Id,
	)
	return err
}

func (r *Repository) GetTargetById(id int) (*Target, error) {
	var t Target
	err := r.db.QueryRow(`SELECT id, name, country, notes, complete FROM targets WHERE id=$1`, id).Scan(&t.Id, &t.Name, &t.Country, &t.Notes, &t.Complete)
	return &t, err
}

func (r *Repository) DeleteTargetFromMission(missionID, targetID int) error {
	//спочатку перевіряємо чи target існує і чи він не completed
	var completed bool
	err := r.db.QueryRow(`
		SELECT completed 
		FROM targets 
		WHERE id = $1 AND mission_id = $2
	`, targetID, missionID).Scan(&completed)
	if err == sql.ErrNoRows {
		return fmt.Errorf("target not found in mission")
	}
	if err != nil {
		return err
	}

	if completed {
		return fmt.Errorf("cannot delete a completed target")
	}

	//видаляємо
	_, err = r.db.Exec(`
		DELETE FROM targets 
		WHERE id = $1 AND mission_id = $2
	`, targetID, missionID)
	if err != nil {
		return err
	}

	return nil
}
