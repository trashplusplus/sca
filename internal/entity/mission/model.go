package mission

import (
	"sca/internal/entity/target"
)

type Mission struct {
	Id       int             `json:"id`
	CatId    int             `json:"cat_id"omitempty`
	Targets  []target.Target `json: "targets"`
	Complete bool            `json:"complete"`
}
