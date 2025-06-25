package saves

import (
	"backend/services/db"
	"database/sql"
	"errors"
	"fmt"
	"shared/services/clock"
	"shared/services/uuid"
	"shared/utils/httperrors"
	"time"
)

type saveMetaFactory struct {
	Clock       clock.Clock
	UUIDFactory uuid.Factory
}

func newSaveMetaFactory(
	clock clock.Clock,
	uuidFactory uuid.Factory,
) SaveMetaFactory {
	return &saveMetaFactory{
		Clock:       clock,
		UUIDFactory: uuidFactory,
	}
}

func (saveMetaFactory *saveMetaFactory) New(name SaveName) SaveMeta {
	return NewSaveMeta(
		SaveId(saveMetaFactory.UUIDFactory.NewUUID().String()),
		saveMetaFactory.Clock.Now(),
		name,
	)
}

type savesMetaRepo struct {
	tx                 db.Tx
	isServiceAvailable bool
	dateFormat         clock.DateFormat
}

func newSavesMetaRepo(
	tx db.Tx,
	dateFormat clock.DateFormat,
) SavesMetaRepo {
	return &savesMetaRepo{
		tx:                 tx,
		isServiceAvailable: tx.Ok(),
		dateFormat:         dateFormat,
	}
}

func (repo *savesMetaRepo) Upsert(meta SaveMeta) error {
	if !repo.isServiceAvailable {
		return httperrors.Err503
	}
	_, err := repo.tx.Tx().Exec(`
		INSERT INTO saves (id, created, last_modified, name)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            last_modified = excluded.last_modified,
            name = excluded.name;`,
		meta.Id,
		meta.Created.Format(repo.dateFormat.String()),
		meta.LastModified.Format(repo.dateFormat.String()),
		meta.Name,
	)
	if err != nil {
		return errors.Join(httperrors.Err503, err)
	}
	return nil
}

func (repo *savesMetaRepo) Delete(id SaveId) error {
	if !repo.isServiceAvailable {
		return httperrors.Err503
	}
	_, err := repo.tx.Tx().Exec(`DELETE FROM saves WHERE id = ?`, id)
	if err != nil {
		return errors.Join(httperrors.Err503, err)
	}
	return nil
}

func (repo *savesMetaRepo) ListSaves(query ListSavesQuery) ([]SaveMeta, error) {
	if !repo.isServiceAvailable {
		return nil, httperrors.Err503
	}
	if errs := query.Valid(); len(errs) != 0 {
		return nil, errors.Join(append(errs, httperrors.Err400)...)
	}
	var orderBy string
	switch query.OrderedBy {
	case OrderedByCreated:
		orderBy = "created"
	case OrderedByLastModified:
		orderBy = "last_modified"
	default:
		return nil, httperrors.Err501
	}

	var orderDir string
	switch query.SortOrder {
	case AscOrder:
		orderDir = "ASC"
	case DescOrder:
		orderDir = "DESC"
	default:
		return nil, httperrors.Err501
	}

	limit := query.SavesPerPage
	offset := int(query.CurrentPage * limit)

	rows, err := repo.tx.Tx().Query(
		fmt.Sprintf(`SELECT id, created, last_modified, name FROM saves ORDER BY %s %s LIMIT ? OFFSET ?`, orderBy, orderDir),
		limit,
		offset)
	if err != nil {
		return nil, errors.Join(httperrors.Err503, err)
	}
	defer rows.Close()

	var saves []SaveMeta
	for rows.Next() {
		var meta SaveMeta
		var createdStr, lastModifiedStr string
		if err := rows.Scan(&meta.Id, &createdStr, &lastModifiedStr, &meta.Name); err != nil {
			return nil, errors.Join(httperrors.Err500, err)
		}
		if meta.Created, err = time.Parse(repo.dateFormat.String(), createdStr); err != nil {
			return nil, errors.Join(httperrors.Err500, err)
		}
		if meta.LastModified, err = time.Parse(repo.dateFormat.String(), lastModifiedStr); err != nil {
			return nil, errors.Join(httperrors.Err500, err)
		}

		saves = append(saves, meta)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(httperrors.Err500, err)
	}

	return saves, nil
}

func (repo *savesMetaRepo) SavesPages(query ListSavesQuery) (int, error) {
	if !repo.isServiceAvailable {
		return 0, httperrors.Err503
	}
	if errs := query.Valid(); len(errs) != 0 {
		return 0, errors.Join(append(errs, httperrors.Err400)...)
	}
	var total int
	err := repo.tx.Tx().QueryRow(`SELECT COUNT(*) FROM saves`).Scan(&total)
	if err != nil {
		return 0, errors.Join(httperrors.Err503, err)
	}

	pages := (total + int(query.SavesPerPage) - 1) / int(query.SavesPerPage)
	return pages, nil
}

func (repo *savesMetaRepo) GetById(id SaveId) (SaveMeta, error) {
	if !repo.isServiceAvailable {
		return SaveMeta{}, httperrors.Err503
	}
	var meta SaveMeta
	var createdStr, lastModifiedStr string

	err := repo.tx.Tx().QueryRow(`SELECT id, created, last_modified, name FROM saves WHERE id = ?`, id).
		Scan(&meta.Id, &createdStr, &lastModifiedStr, &meta.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return SaveMeta{}, httperrors.Err404 // Not found error (assuming defined)
		}
		return SaveMeta{}, errors.Join(httperrors.Err503, err)
	}

	meta.Created, err = time.Parse(time.RFC3339, createdStr)
	if err != nil {
		return SaveMeta{}, errors.Join(httperrors.Err500, err)
	}
	meta.LastModified, err = time.Parse(time.RFC3339, lastModifiedStr)
	if err != nil {
		return SaveMeta{}, errors.Join(httperrors.Err503, err)
	}

	return meta, nil
}
