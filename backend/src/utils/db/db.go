package db

import (
	"backend/src/utils/httperrors"
	"backend/src/utils/logger"
	"backend/src/utils/services/scopecleanup"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ogiusek/ioc"
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
	filePath string
}

func Package(filePath string) Pkg {
	return Pkg{
		filePath: filePath,
	}
}

func (pkg Pkg) Register(c ioc.Dic) {
	ioc.RegisterSingleton(c, func(c ioc.Dic) DB {
		db, err := sql.Open("sqlite3", pkg.filePath)
		if err != nil {
			panic(err)
		}
		go func() {
			for {
				if err := db.Ping(); err != nil {
					// panic(err)
				}
				time.Sleep(time.Hour)
			}
		}()
		return NewDB(db, true)
	})

	ioc.RegisterScoped(c, func(c ioc.Dic) Tx {
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
			err := tx.Commit()
			logger.Error(errors.Join(ErrCommitFailed, err))
		})
		return NewTx(tx, true)
	})

}
