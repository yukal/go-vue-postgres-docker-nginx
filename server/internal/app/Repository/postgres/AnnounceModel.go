package postgres

import (
	"database/sql"
	"errors"
	"time"
)

type AnnounceModel struct {
	// db InstanceInterface
	db *sql.DB
}

type AnnounceRecord struct {
	Id          uint64
	RegionId    uint8
	ImageId     uint64
	ImageName   sql.NullString
	Sex         sql.NullByte
	Age         sql.NullByte
	Height      sql.NullInt16
	Weight      sql.NullInt16
	Title       sql.NullString
	Description sql.NullString
	Link        sql.NullString
	CreatedAt   sql.NullString
}

type Announce struct {
	Id          uint64    `json:"announce_id"`
	RegionId    uint8     `json:"region_id"`
	ImageId     uint64    `json:"media_id"`
	ImageName   string    `json:"media_name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`

	Images []AnnounceAddition `json:"images"`
	Phones []AnnounceAddition `json:"phones"`
}

type AnnounceAddition struct {
	Id   uint64 `json:"id"`
	Item string `json:"item"`
}

type AnnounceAdditions struct {
	Medias []AnnounceAddition `json:"medias"`
	Phones []AnnounceAddition `json:"phones"`
}

func (mod *AnnounceModel) Total() (total uint64, err error) {
	err = mod.db.
		QueryRow(`
			SELECT COUNT(announce_id)
			FROM announce`).
		Scan(&total)

	if errors.Is(err, sql.ErrNoRows) {
		total = 0
		err = nil
	}

	return
}

func (mod *AnnounceModel) Exist(Id uint64) (exist bool, err error) {
	var rec sql.NullInt64

	query := `
	SELECT announce_id
	FROM announce
	WHERE announce_id = $1
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

func (mod *AnnounceModel) Get(offset, limit uint64) ([]Announce, error) {
	items := []Announce{}
	params := []any{}

	query := `
    SELECT
      announce_id,
      region_id,
      image_id,
      title,
      description,
      link,
      created_at
    FROM announce
	ORDER BY created_at`

	query += " LIMIT $1 OFFSET $2"
	params = append(params, limit, offset)

	rows, err := mod.db.Query(query, params...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		rec := AnnounceRecord{}

		err = rows.Scan(
			&rec.Id,
			&rec.RegionId,
			&rec.ImageId,
			&rec.Title,
			&rec.Description,
			&rec.Link,
			&rec.CreatedAt,
		)

		if err != nil {
			continue
		}

		data := Announce{
			Id:       rec.Id,
			RegionId: rec.RegionId,
			ImageId:  rec.ImageId,
		}

		if rec.Title.Valid {
			data.Title = rec.Title.String
		}
		if rec.Description.Valid {
			data.Description = rec.Description.String
		}
		if rec.Link.Valid {
			data.Link = rec.Link.String
		}
		if rec.CreatedAt.Valid {
			createdAt, _ := time.Parse(TIMESTAMP_FORMAT, rec.CreatedAt.String)
			data.CreatedAt = createdAt
		}

		items = append(items, data)
	}

	return items, nil
}
