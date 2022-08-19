package go_pg_goose_tests_test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
)

const (
	dbName = "postgres"
	dbUser = "postgres"
	dbPass = "pg_pass"
)

type BasicSuite struct {
	suite.Suite

	db    *sql.DB
	pool  *dockertest.Pool
	dbRef *dockertest.Resource
}

// Set up the postgres container
func (s *BasicSuite) SetupSuite() {
	pool, err := dockertest.NewPool("")
	s.Require().NoError(err)
	s.pool = pool

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", dbName),
			fmt.Sprintf("POSTGRES_USER=%s", dbUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", dbPass),
			"listen_address=*",
		},
	})
	s.Require().NoError(err)
	s.dbRef = resource

	s.Require().NoError(pool.Retry(func() error {
		var err error
		s.db, err = sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPass, resource.GetHostPort("5432/tcp"), dbName))
		if err != nil {
			return err
		}

		return s.db.Ping()
	}))

	// Set the dialect for goose, only needs to happen ones
	s.Require().NoError(goose.SetDialect("postgres"))
}

// Remove container and volumes after suite is done running
func (s *BasicSuite) TearDownSuite() {
	s.Require().NoError(s.pool.Purge(s.dbRef))
}

// Run migrations and populate database
func (s *BasicSuite) SetupTest() {
	start := time.Now()
	defer func() {
		fmt.Printf("Finished BasicSuite.SetupTest in %s\n", time.Since(start))
	}()

	s.Require().NoError(goose.Up(s.db, "migrations"))

	var err error

	// This is obviously very manual, just a demonstration
	_, err = s.db.Exec("INSERT INTO users VALUES (0, 'admin', 'admin', 'user', 'admin@org.org')")
	s.Require().NoError(err)
	_, err = s.db.Exec("INSERT INTO users VALUES (1, 'jboy', 'john', 'boy', 'jon.boy@someplace.com')")
	s.Require().NoError(err)
}

// Invert all migrations after test runs. Another option would be to manually truncate tables.
func (s *BasicSuite) TearDownTest() {
	cwd, err := os.Getwd()
	s.Require().NoError(err)

	s.Require().NoError(goose.DownTo(s.db, filepath.Join(cwd, "migrations"), 0))
}

func (s *BasicSuite) TestUsersExist() {
	rows, err := s.db.Query("SELECT * FROM users")
	s.Require().NoError(err)

	var seen int
	for rows.Next() {
		seen += 1
	}

	s.Require().Equal(2, seen)
}

func (s *BasicSuite) TestPostsEmpty() {
	rows, err := s.db.Query("SELECT * FROM posts")
	s.Require().NoError(err)

	var seen int
	for rows.Next() {
		seen += 1
	}

	s.Require().Equal(0, seen)
}

func TestBasic(t *testing.T) {
	suite.Run(t, &BasicSuite{})
}
