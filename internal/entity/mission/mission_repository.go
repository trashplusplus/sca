package mission

import (
	"database/sql"
	"errors"
	"fmt"
	"sca/internal/entity/target"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(m *Mission) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRow(
		`INSERT INTO missions (cat_id, complete) VALUES ($1, $2) RETURNING id`,
		nil,
		m.Complete,
	).Scan(&m.Id)
	if err != nil {
		return err
	}

	for i := range m.Targets {
		t := &m.Targets[i]
		err = tx.QueryRow(
			`INSERT INTO targets (mission_id, name, country, notes, complete)
			 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			m.Id, t.Name, t.Country, t.Notes, t.Complete,
		).Scan(&t.Id)
		if err != nil {
			return err
		}
	}

	//підтверджуємо транзакці
	return tx.Commit()
}

func (r *Repository) GetMissionById(missionID int) (*Mission, error) {
	//left join щоб повернути місію навіть якщо targets немає
	rows, err := r.db.Query(`
		SELECT 
			m.id, m.cat_id, m.complete,
			t.id, t.name, t.country, t.notes, t.complete
		FROM missions m
		LEFT JOIN targets t ON t.mission_id = m.id
		WHERE m.id = $1
	`, missionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mission *Mission
	targets := []target.Target{}

	for rows.Next() {
		var mID int
		var catID sql.NullInt64
		var mComplete bool
		var tID sql.NullInt64
		var tName, tCountry, tNotes sql.NullString
		var tComplete sql.NullBool

		if err := rows.Scan(&mID, &catID, &mComplete, &tID, &tName, &tCountry, &tNotes, &tComplete); err != nil {
			return nil, err
		}

		if mission == nil {
			mission = &Mission{
				Id:       mID,
				Complete: mComplete,
				Targets:  []target.Target{},
			}
			if catID.Valid {
				id := int(catID.Int64)
				mission.CatId = id
			}
		}

		if tID.Valid {
			target := target.Target{
				Id:       int(tID.Int64),
				Name:     tName.String,
				Country:  tCountry.String,
				Notes:    tNotes.String,
				Complete: tComplete.Bool,
			}
			targets = append(targets, target)
		}
	}

	if mission == nil {
		return nil, sql.ErrNoRows
	}

	mission.Targets = targets
	return mission, nil
}

func (r *Repository) List() ([]Mission, error) {

	rows, err := r.db.Query(`
		SELECT 
			m.id, m.cat_id, m.complete,
			t.id, t.name, t.country, t.notes, t.complete
		FROM missions m
		LEFT JOIN targets t ON t.mission_id = m.id
		ORDER BY m.id, t.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	missionMap := make(map[int]*Mission)
	for rows.Next() {
		var mID int
		var catID sql.NullInt64
		var mComplete bool
		var tID sql.NullInt64
		var tName, tCountry, tNotes sql.NullString
		var tComplete sql.NullBool

		if err := rows.Scan(&mID, &catID, &mComplete, &tID, &tName, &tCountry, &tNotes, &tComplete); err != nil {
			return nil, err
		}

		m, exists := missionMap[mID]
		if !exists {
			m = &Mission{
				Id:       mID,
				Complete: mComplete,
			}
			if catID.Valid {
				id := int(catID.Int64)
				m.CatId = id
			}
			m.Targets = []target.Target{}
			missionMap[mID] = m
		}

		if tID.Valid {
			target := target.Target{
				Id:       int(tID.Int64),
				Name:     tName.String,
				Country:  tCountry.String,
				Notes:    tNotes.String,
				Complete: tComplete.Bool,
			}
			m.Targets = append(m.Targets, target)
		}
	}

	var missions []Mission
	for _, m := range missionMap {
		missions = append(missions, *m)
	}

	return missions, nil
}

func (r *Repository) AssignCat(missionID int, catID int) error {
	res, err := r.db.Exec(`UPDATE missions SET cat_id=$1 WHERE id=$2`, catID, missionID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("mission not found")
	}
	return nil
}

func (r *Repository) UpdateMissionComplete(missionID int, complete bool) error {
	_, err := r.db.Exec(`UPDATE missions SET complete=$1 WHERE id=$2`, complete, missionID)
	return err
}

func (r *Repository) DeleteMission(missionID int) error {

	var catID sql.NullInt64
	row := r.db.QueryRow(`SELECT cat_id FROM missions WHERE id=$1`, missionID)
	if err := row.Scan(&catID); err != nil {
		return err
	}
	if catID.Valid {
		return errors.New("cannot delete mission assigned to cat")
	}

	_, err := r.db.Exec(`DELETE FROM missions WHERE id=$1`, missionID)
	return err
}

func (r *Repository) DeleteTarget(targetID int) error {
	var complete bool
	row := r.db.QueryRow(`SELECT complete FROM targets WHERE id=$1`, targetID)
	if err := row.Scan(&complete); err != nil {
		return err
	}
	if complete {
		return errors.New("cannot delete completed target")
	}

	_, err := r.db.Exec(`DELETE FROM targets WHERE id=$1`, targetID)
	return err
}

func (r *Repository) CreateTarget(missionID int, t target.Target) error {
	//статус місії
	var complete bool
	err := r.db.QueryRow(`SELECT complete FROM missions WHERE id = $1`, missionID).Scan(&complete)
	if err != nil {
		return fmt.Errorf("failed to check mission: %w", err)
	}

	if complete {
		return fmt.Errorf("cannot add target: mission %d is already completed", missionID)
	}

	//якщо місія не завершилась то додаю таргет
	err = r.db.QueryRow(`
		INSERT INTO targets (mission_id, name, country, notes, complete)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, missionID, t.Name, t.Country, t.Notes, t.Complete).Scan(&t.Id)

	if err != nil {
		return fmt.Errorf("failed to insert target: %w", err)
	}

	return nil
}
