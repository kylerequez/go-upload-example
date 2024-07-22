package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/kylerequez/go-upload-example/src/models"
)

type UploadRepository struct {
	DB   *sql.DB
	Name string
}

func NewUploadRepository(db *sql.DB, name string) *UploadRepository {
	return &UploadRepository{
		DB:   db,
		Name: name,
	}
}

func (ur *UploadRepository) GetAllUploads() ([]models.File, error) {
	sql := fmt.Sprintf(`
			SELECT
				id,
				filename,
				filesize,
				filetype,
				createdAt
			FROM
				%s;
		`, ur.Name)

	stmt, err := ur.DB.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []models.File{}
	for rows.Next() {
		file := models.File{}

		err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Size,
			&file.Type,
			&file.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}
	defer rows.Close()

	return files, err
}

func (ur *UploadRepository) CreateFile(file models.File) error {
	sql := fmt.Sprintf(`
			INSERT INTO
				%s
			(
				filename,
				filesize,
				filetype
			)
			VALUES
			(
				$1,
				$2,
				$3
			);
		`, ur.Name)

	stmt, err := ur.DB.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		file.Name,
		file.Size,
		file.Type,
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count < 1 {
		return errors.New("no file were inserted in the database")
	}

	return nil
}

func (ur *UploadRepository) GetUploadByName(name string) (*models.File, error) {
	sql := fmt.Sprintf(`
			SELECT
				id,
				filename,
				filesize,
				filetype,
				createdAt
			FROM
				%s
			WHERE
				filename = $1
			LIMIT
				1;
		`, ur.Name)

	stmt, err := ur.DB.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	file := &models.File{}
	if rows.Next() {
		if err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Size,
			&file.Type,
			&file.CreatedAt,
		); err != nil {
			return nil, err
		}

		return file, nil
	}
	defer rows.Close()

	return nil, nil
}

func (ur *UploadRepository) GetUploadById(id uuid.UUID) (*models.File, error) {
	sql := fmt.Sprintf(`
			SELECT
				id,
				filename,
				filesize,
				filetype,
				createdAt
			FROM
				%s
			WHERE
				id = $1
			LIMIT
				1;
		`, ur.Name)

	stmt, err := ur.DB.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	file := &models.File{}
	for rows.Next() {
		if err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Size,
			&file.Type,
			&file.CreatedAt,
		); err != nil {
			return nil, err
		}

		return file, nil
	}
	defer rows.Close()

	return nil, nil
}

func (ur *UploadRepository) DeleteUploadById(id uuid.UUID) error {
	sql := fmt.Sprintf(
		`
			DELETE FROM
				%s
			WHERE
				id = $1;
		`, ur.Name)

	stmt, err := ur.DB.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count <= 0 {
		return errors.New("there was no file deleted")
	}

	return nil
}
