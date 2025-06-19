package db

import (
	"backend/src/utils/httperrors"
	"backend/src/utils/logger"
	"backend/src/utils/services/scopecleanup"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"

	"github.com/ogiusek/ioc/v2"
)

// db

type DB struct {
	db *sql.DB
	ok bool
}

func NewDB(db *sql.DB, ok bool) DB {
	return DB{db: db, ok: ok}
}

func (db DB) DB() *sql.DB {
	return db.db
}

func (db DB) Ok() bool {
	return db.ok
}

// tx

var (
	ErrCommitFailed error = errors.Join(httperrors.Err500, errors.New("failed to commit changed"))
)

type Tx struct {
	tx *sql.Tx
	ok bool
}

func NewTx(tx *sql.Tx, ok bool) Tx {
	return Tx{tx: tx, ok: ok}
}

func (tx Tx) Tx() *sql.Tx {
	return tx.tx
}

func (tx Tx) Ok() bool {
	return tx.ok
}

// pkg

type Pkg struct {
	dbPath        string
	migrationsDir string
}

func Package(dbPath string, migrationsDir string) Pkg {
	return Pkg{
		dbPath:        dbPath,
		migrationsDir: migrationsDir,
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) DB {
		if err := os.MkdirAll(filepath.Dir(pkg.dbPath), os.ModePerm); err != nil {
			panic(fmt.Sprintf("error creating directories %s", err.Error()))
		}
		db, err := sql.Open("sqlite3", pkg.dbPath)
		if err != nil {
			panic(errors.Join(errors.New("opening database"), err))
		}
		driver, err := sqlite.WithInstance(db, &sqlite.Config{})
		if err != nil {
			panic(errors.Join(errors.New("creating driver"), err))
		}
		mig, err := migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", pkg.migrationsDir),
			"sqlite3",
			driver,
		)
		if err != nil {
			panic(errors.Join(errors.New("creating migration"), err))
		}
		if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
			panic(errors.Join(errors.New("running up migration"), err))
		}
		go func() {
			for {
				db.Ping()
				time.Sleep(time.Hour)
			}
		}()
		return NewDB(db, true)
	})

	ioc.RegisterScoped(b, func(c ioc.Dic) Tx {
		logger := ioc.Get[logger.Logger](c)
		db := ioc.Get[DB](c)
		if !db.Ok() {
			return NewTx(nil, false)
		}
		tx, err := db.DB().Begin()
		if err != nil {
			logger.Error(errors.Join(httperrors.Err503, err))
		}
		scopeCleanUp := ioc.Get[scopecleanup.ScopeCleanUp](c)
		scopeCleanUp.AddCleanListener(func(args scopecleanup.CleanUpArgs) {
			if args.Error != nil || err != nil {
				return
			}
			if err := tx.Commit(); err != nil {
				logger.Error(errors.Join(ErrCommitFailed, err))
			}
		})
		return NewTx(tx, true)
	})

}
