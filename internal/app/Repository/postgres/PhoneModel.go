package postgres

import (
	"database/sql"
	"errors"
	"fmt"
)

type PhoneModel struct {
	// db InstanceInterface
	db *sql.DB
}

type Phone struct {
	Id     uint64
	Number string
}

type PhoneRecord struct {
	Id     uint64
	Number sql.NullString
}

type AddedPhones struct {
	Exist []Phone
	Added []Phone
	Size  int
}

func (i AddedPhones) GetItem(n int) Phone {
	if existLen := len(i.Exist); existLen > 0 && existLen >= n {
		return i.Exist[n]
	}

	if addedLen := len(i.Added); addedLen > 0 && addedLen >= n {
		return i.Added[n]
	}

	return Phone{}
}

func (i AddedPhones) GetIDs() []uint64 {
	dmap := make(map[uint64]any, i.Size)
	rels := make([]uint64, 0, i.Size)

	for _, item := range i.Added {
		if _, exist := dmap[item.Id]; !exist {
			dmap[item.Id] = nil
			rels = append(rels, item.Id)
		}
	}

	for _, item := range i.Exist {
		if _, exist := dmap[item.Id]; !exist {
			dmap[item.Id] = nil
			rels = append(rels, item.Id)
		}
	}

	return rels
}

func (i AddedPhones) GetNumbers() []string {
	dmap := make(map[string]any, i.Size)
	rels := make([]string, 0, i.Size)

	for _, item := range i.Added {
		if _, exist := dmap[item.Number]; !exist {
			dmap[item.Number] = nil
			rels = append(rels, item.Number)
		}
	}

	for _, item := range i.Exist {
		if _, exist := dmap[item.Number]; !exist {
			dmap[item.Number] = nil
			rels = append(rels, item.Number)
		}
	}

	return rels
}

func (mod *PhoneModel) Total() (total int, err error) {
	err = mod.db.
		QueryRow(`
			SELECT COUNT(phone_id)
			FROM phone`).
		Scan(&total)

	if errors.Is(err, sql.ErrNoRows) {
		total = 0
		err = nil
	}

	return
}

func (mod *PhoneModel) Exist(Id uint64) (exist bool, err error) {
	var val sql.NullInt64

	query := `
	SELECT phone_id
	FROM phone
	WHERE phone_id = $1
	LIMIT 1`

	err = mod.db.
		QueryRow(query, Id).
		Scan(&val)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		return
	}

	if err == nil && val.Valid {
		exist = uint64(val.Int64) == Id
	}

	return
}

func (mod *PhoneModel) HasDuplicateNumber(phone string) (exist bool, err error) {
	var rec sql.NullString

	query := `
	SELECT phone_number
	FROM phone
	WHERE phone_number = $1
	LIMIT 1`

	err = mod.db.
		QueryRow(query, phone).
		Scan(&rec)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		return
	}

	if err == nil && rec.Valid {
		exist = rec.String == phone
	}

	return
}

func (mod *PhoneModel) GetIdByNumber(phone string) (Id uint64, err error) {
	const op = "PhoneModel.GetIgByNumber"
	var rec sql.NullInt64

	query := `
	SELECT phone_id
	FROM phone
	WHERE phone_number = $1
	LIMIT 1`

	err = mod.db.
		QueryRow(query, phone).
		Scan(&rec)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		return
	}

	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
	}

	if err == nil && rec.Valid {
		Id = uint64(rec.Int64)
	}

	return
}

func (mod *PhoneModel) Insert(phone string) (lastInsertId uint64, err error) {
	const op = "PhoneModel.Insert"

	query := `INSERT INTO phone(phone_number) values($1) RETURNING phone_id`
	err = mod.db.QueryRow(query, phone).
		Scan(&lastInsertId)

	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
	}

	return
}

func (mod *PhoneModel) GetOrInsertItems(phones []string) (list AddedPhones, err error) {
	var id uint64

	size := len(phones)
	dmap := make(map[uint64]any, size)

	list.Exist = make([]Phone, 0, len(phones))
	list.Added = make([]Phone, 0, len(phones))

	for _, num := range phones {
		id, err = mod.GetIdByNumber(num)

		if err != nil {
			return
		}

		if id > 0 {

			if _, exist := dmap[id]; !exist {
				dmap[id] = nil

				list.Exist = append(list.Exist, Phone{
					Id:     id,
					Number: num,
				})
			}

		} else {

			id, err = mod.Insert(num)
			if err != nil {
				return
			}

			if _, exist := dmap[id]; !exist {
				dmap[id] = nil

				list.Added = append(list.Added, Phone{
					Id:     id,
					Number: num,
				})
			}

		}
	}

	list.Size = len(list.Exist) + len(list.Added)
	return
}

func (mod *PhoneModel) Delete(phoneId uint64) (int64, error) {
	query := `
		DELETE FROM phone
		WHERE phone_id = $1`

	result, err := mod.db.Exec(query, phoneId)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (mod *PhoneModel) SetPhone(id uint64, phone string) (int64, error) {
	query := `
		UPDATE phone SET phone_number = $1
		WHERE phone_id = $2`

	result, err := mod.db.Exec(query, phone, id)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
