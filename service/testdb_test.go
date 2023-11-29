package service

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

type testDB struct {
	User string
	Pass string
	Host string

	DB *sql.DB
}

func (tdb *testDB) getRootURI() string {
	return fmt.Sprintf("postgres://%v:%v@%v?sslmode=disable", tdb.User, tdb.Pass, tdb.Host)
}

func (tdb *testDB) createDatabase() (err error) {
	db, err := sql.Open("postgres", tdb.getRootURI())
	if err != nil {
		return fmt.Errorf("could not open postgres database (%v) error (%w)", tdb.getRootURI(), err)
	}
	defer func() {
		err2 := db.Close()
		if err2 != nil {
			err = errors.Join(err, fmt.Errorf("could not close postgres db (%w)", err2))
		}
	}()

	_, err = db.Exec("CREATE DATABASE petinfoservice;")
	if err != nil {
		return fmt.Errorf("could not create petinfoservice database (%w)", err)
	}

	return nil
}

func (tdb *testDB) destroyDatabase() (err error) {
	db, err := sql.Open("postgres", tdb.getRootURI())
	if err != nil {
		return fmt.Errorf("could not open postgres database (%v) error (%w)", tdb.getRootURI(), err)
	}
	defer func() {
		err2 := db.Close()
		if err2 != nil {
			err = errors.Join(err, fmt.Errorf("could not close postgres db (%w)", err2))
		}
	}()

	_, err = db.Exec("DROP DATABASE petinfoservice;")
	if err != nil {
		return fmt.Errorf("could not delete database during cleanup stage (%w)", err)
	}

	return nil
}

func (tdb *testDB) loadSchema() (err error) {
	connURI := fmt.Sprintf("postgres://%v:%v@%v/petinfoservice?sslmode=disable", tdb.User, tdb.Pass, tdb.Host)

	db, err := sql.Open("postgres", connURI)
	if err != nil {
		return fmt.Errorf("could not open petinfoservice database (%v) error (%w)", connURI, err)
	}
	defer func() {
		// Close database connection only if an error occured
		// This prevent a lingering database connection if caller does not
		// tdb.Close() (which they shouldn't call if tdb.Open() fails anyway)
		if err == nil {
			return
		}

		err2 := db.Close()
		if err2 != nil {
			err = errors.Join(err, fmt.Errorf("could not close postgres db (%w)", err2))
		}
	}()

	schema, err := os.ReadFile("../database/schema.sql")
	if err != nil {
		return fmt.Errorf("could not read ../database/schema.sql error (%w)", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("could not create schema in database (%w)", err)
	}

	tdb.DB = db

	return nil
}

func (tdb *testDB) Open() (err error) {
	if err = tdb.createDatabase(); err != nil {
		return err
	}

	if err = tdb.loadSchema(); err != nil {
		return err
	}

	return nil
}

func (tdb *testDB) Close() (err error) {
	if err = tdb.DB.Close(); err != nil {
		return fmt.Errorf("could not close petinfoservice database connection (%w)", err)
	}

	// TODO: if an error occurs in closing the DB connection used by tests, we never delete
	// the database that we created (should not matter that much with test containers)

	if err = tdb.destroyDatabase(); err != nil {
		return err
	}

	return nil
}
