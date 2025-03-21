package dict

import (
	"fmt"
	"schoolonline/postgree"
)

func GetTgIdDirector(id_user int64) (res int64) {
	query := `
	SELECT u.tg_id FROM users AS u 
WHERE u.role = 'director'
AND u.school_id = (SELECT p.school_id FROM users AS p WHERE p.tg_id = $1)
`

	row := postgree.MainDBX.QueryRow(query, id_user)
	err := row.Scan(&res)
	if err != nil {

		fmt.Println(err)
		return
	}
	return
}
