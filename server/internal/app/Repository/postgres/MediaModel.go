package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"path"
)

type MediaModel struct {
	// db InstanceInterface
	db *sql.DB
}

type Media struct {
	Id   uint64
	Name string
}

type MediaRecord struct {
	Id   uint64
	Name sql.NullString
}

type AddedMedias struct {
	Exist []Media
	Added []Media
	Size  int
}

func (i AddedMedias) GetItem(n int) Media {
	if existLen := len(i.Exist); existLen > 0 && existLen >= n {
		return i.Exist[n]
	}

	if addedLen := len(i.Added); addedLen > 0 && addedLen >= n {
		return i.Added[n]
	}

	return Media{}
}

func (i AddedMedias) GetIDs() []uint64 {
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

func (i AddedMedias) GetMedias() []string {
	dmap := make(map[string]any, i.Size)
	rels := make([]string, 0, i.Size)

	for _, item := range i.Added {
		if _, exist := dmap[item.Name]; !exist {
			dmap[item.Name] = nil
			rels = append(rels, item.Name)
		}
	}

	for _, item := range i.Exist {
		if _, exist := dmap[item.Name]; !exist {
			dmap[item.Name] = nil
			rels = append(rels, item.Name)
		}
	}

	return rels
}

func (mod *MediaModel) Total() (total int, err error) {
	err = mod.db.
		QueryRow(`
			SELECT COUNT(media_id)
			FROM media`).
		Scan(&total)

	if errors.Is(err, sql.ErrNoRows) {
		total = 0
		err = nil
	}

	return
}

func (mod *MediaModel) Exist(Id uint64) (exist bool, err error) {
	var rec sql.NullInt64

	query := `
	SELECT media_id
	FROM media
	WHERE media_id = $1
	LIMIT 1`

	err = mod.db.
		QueryRow(query, Id).
		Scan(&rec)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		return
	}

	if err == nil && rec.Valid {
		exist = uint64(rec.Int64) == Id
	}

	return
}

func (mod *MediaModel) HasDuplicateName(media string) (exist bool, err error) {
	var rec sql.NullString

	query := `
	SELECT media_name
	FROM media
	WHERE media_name = $1
	LIMIT 1`

	err = mod.db.
		QueryRow(query, media).
		Scan(&rec)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		return
	}

	if err == nil && rec.Valid {
		exist = rec.String == media
	}

	return
}

func (mod *MediaModel) GetIdByName(media string) (Id uint64, err error) {
	const op = "MediaModel.GetIdByName"
	var rec sql.NullInt64

	query := `
	SELECT media_id
	FROM media
	WHERE media_name = $1
	LIMIT 1`

	err = mod.db.
		QueryRow(query, media).
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

func (mod *MediaModel) Insert(media string) (lastInsertId uint64, err error) {
	const op = "MediaModel.Insert"

	query := `INSERT INTO media(media_name) values($1) RETURNING media_id`
	err = mod.db.QueryRow(query, media).
		Scan(&lastInsertId)

	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
	}

	return
}

func (mod *MediaModel) GetOrInsertItems(medias []string) (list AddedMedias, err error) {
	var id uint64

	size := len(medias)
	dmap := make(map[uint64]any, size)

	list.Exist = make([]Media, 0, size)
	list.Added = make([]Media, 0, size)

	for _, med := range medias {
		filename := path.Base(med)
		id, err = mod.GetIdByName(filename)

		if err != nil {
			return
		}

		if id > 0 {

			if _, exist := dmap[id]; !exist {
				dmap[id] = nil

				list.Exist = append(list.Exist, Media{
					Id:   id,
					Name: filename,
				})
			}

		} else {

			id, err = mod.Insert(filename)
			if err != nil {
				return
			}

			if _, exist := dmap[id]; !exist {
				dmap[id] = nil

				list.Added = append(list.Added, Media{
					Id:   id,
					Name: filename,
				})
			}

		}
	}

	list.Size = len(list.Exist) + len(list.Added)
	return
}

func (mod *MediaModel) Delete(mediaId uint64) (int64, error) {
	query := `
		DELETE FROM media
		WHERE media_id = $1`

	result, err := mod.db.Exec(query, mediaId)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (mod *MediaModel) SetMedia(id uint64, name string) (int64, error) {
	query := `
		UPDATE media SET media_name = $1
		WHERE media_id = $2`

	result, err := mod.db.Exec(query, name, id)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
