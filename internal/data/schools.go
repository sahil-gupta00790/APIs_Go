package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"Apis_go.sahil.net/internal/validator"
	"github.com/lib/pq"
)

type School struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Level     string    `json:"level"`
	Contact   string    `json:"contact"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email,omitempty"`
	Website   string    `json:"website,omitempty"`
	Address   string    `json:"address"`
	Mode      []string  `json:"mode"`
	Version   int32     `json:"version"`
}

func ValidateSchool(v *validator.Validator, input *School) {
	//use check method to execute validation check
	v.Check(input.Name != "", "name", "must be provided")
	v.Check(len(input.Name) <= 200, "name", "must not be more than 200 bytes")

	v.Check(input.Level != "", "level", "must be provided")
	v.Check(len(input.Level) <= 200, "level", "must not be more than 200 bytes")

	v.Check(input.Contact != "", "contact", "must be provided")
	v.Check(len(input.Contact) <= 200, "contact", "must not be more than 200 bytes")

	v.Check(input.Phone != "", "phone", "must be provided")
	v.Check(validator.Matches(input.Phone, validator.PhoneRX), "phone", "must be valid phone number")

	v.Check(input.Email != "", "email", "must be provided")
	v.Check(validator.Matches(input.Email, validator.EmailRX), "email", "must be valid email address")

	v.Check(input.Website != "", "website", "must be provided")
	v.Check(validator.ValidWebsite(input.Website), "website", "must be valid phone number")

	v.Check(input.Address != "", "address", "must be provided")
	v.Check(len(input.Address) <= 500, "address", "must not be more than 500 bytes")

	v.Check(input.Mode != nil, "mode", "must be provided")
	v.Check(len(input.Mode) >= 1, "mode", "must contain atleast one mode")
	v.Check(len(input.Mode) <= 5, "mode", "must contain atmost five mode")
	v.Check(validator.Unique(input.Mode), "mode", "must not contain duplicate entries")

}

// define a schoolModel which wraps a sql.DB connection pool
type SchoolModel struct {
	DB *sql.DB
}

// Insert( allows us to create a new school)
func (m SchoolModel) Insert(school *School) error {
	query := `INSERT INTO schools (name,level,contact,phone,email,website,address,mode)
	VALUES( $1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING id,created_at,version`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//cleanup to prevent memroy leaks
	defer cancel()

	//collect the data fields into a slice
	args := []interface{}{school.Name, school.Level, school.Contact, school.Phone, school.Email, school.Website, school.Address, school.Mode}
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&school.ID, &school.CreatedAt, &school.Version)
}

// Get() alows to retriece a specific school
func (m SchoolModel) Get(id int64) (*School, error) {
	//Ensure the valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	//create the query
	query :=
		`SELECT id,created_at,name,level,contact,phone,email,website,address,mode,version
	FROM schools
	WHERE id=$1
	`
	//declare a school var. to hold the returned data
	var school School
	//create a context
	//time is started here
	//so if no open connection m, the qeury time is not calcualted but the time it waits for open connection of db is also counted
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//cleanup to prevent memroy leaks
	defer cancel()

	//execute the query using queryROw
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&school.ID, &school.CreatedAt, &school.Name, &school.Level, &school.Contact, &school.Phone, &school.Email, &school.Website, &school.Address, pq.Array(&school.Mode), &school.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	//success
	return &school, nil
}

//Simlultaneous request to api might make things worse
//How to set it right?
//Optimistic locking(version number)
//

// Update() allows us to edit/alter a specific school
func (m SchoolModel) Update(school *School) error {
	//Create a query for new updated school data
	query := `
	UPDATE schools
	SET name=$1 , level=$2,contact=$3,phone=$4,email=$5,website=$6,address=$7,mode=$8,version=version+1
	WHERE id=$9
	AND version = $10
	RETURNING version`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//cleanup to prevent memroy leaks
	defer cancel()

	args := []interface{}{
		school.Name,
		school.Level,
		school.Contact,
		school.Phone,
		school.Email,
		school.Website,
		school.Address,
		pq.Array(school.Mode),
		school.ID,
		school.Version,
	}
	//check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&school.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflig
		default:
			return err
		}

	}
	return nil
}

// Delete() removes a specific school
func (m SchoolModel) Delete(id int64) error {
	//Ensure the valid id
	if id < 1 {
		return ErrRecordNotFound
	}
	//creating the delete query
	query := `
	DELETE FROM schools
	WHERE id=$1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//cleanup to prevent memroy leaks
	defer cancel()

	//execute the quret
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	//Check how many rules were affected by the delete operation
	//call the rowsaffected method on the results varriable
	rowsAddected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	//check if no rows are affected
	if rowsAddected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

//the getall mathod returns a list of all the schools sorted by id

func (m SchoolModel) GetAll(name string, level string, mode []string, filters Filters) ([]*School, error) {
	//construct the query
	query := `SELECT id,created_at,name,level,contact,phone,email,website,address,mode,version 
	FROM schools
	WHERE(to_tsvector('simple',name) @@ plainto_tsquery('simple',$1) OR $1 = '')
	AND (to_tsvector('simple',level) @@ plainto_tsquery('simple' , $2) OR $2 = '')
	AND (mode @> $3 OR $3 = '{}') 
	ORDER BY id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, name, level, pq.Array(mode))
	if err != nil {
		return nil, err
	}
	//close the result set
	defer rows.Close()
	//initialize an empty slice to hold the school data
	schools := []*School{}
	for rows.Next() {
		var school School
		//scan for the values form the row into the school struct
		err := rows.Scan(
			&school.ID,
			&school.CreatedAt,
			&school.Name,
			&school.Level,
			&school.Contact,
			&school.Phone,
			&school.Email,
			&school.Website,
			&school.Address,
			pq.Array(&school.Mode),
			&school.Version,
		)
		if err != nil {
			return nil, err
		}
		//add the school to the slice
		schools = append(schools, &school)
	}
	//check for errors after looping thouigh resuilt set
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return schools, err
}
