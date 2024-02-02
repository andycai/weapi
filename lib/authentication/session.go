package authentication

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mysql/v2"
	"github.com/gofiber/storage/sqlite3"
)

var storedAuthenticationSession *session.Store

func SessionSetup(dbDriver string, db *sql.DB, dsn, tableName string) {
	if dbDriver == "mysql" {
		sessionMySQLStart(db, tableName)
	} else {
		sessionStart(dsn, tableName)
	}
}

func sessionMySQLStart(db *sql.DB, tableName string) {
	store := mysql.New(mysql.Config{
		Db:    db,
		Table: tableName,
	})

	authSession := session.New(session.Config{
		Storage:    store,
		Expiration: 24 * time.Hour,
	})

	storedAuthenticationSession = authSession
}

func sessionStart(dsn, tableName string) {
	store := sqlite3.New(sqlite3.Config{
		Database: dsn,
		Table:    tableName,
	})

	authSession := session.New(session.Config{
		Storage: store,
	})

	storedAuthenticationSession = authSession
}

func AuthStore(c *fiber.Ctx, userID uint) {
	session, err := storedAuthenticationSession.Get(c)
	if err != nil {
		panic(err)
	}

	session.Set("authentication", userID)
	if err := session.Save(); err != nil {
		panic(err)
	}
}

func AuthGet(c *fiber.Ctx) (bool, uint) {
	session, err := storedAuthenticationSession.Get(c)
	if err != nil {
		panic(err)
	}

	value := session.Get("authentication")
	if value == nil {
		return false, 0
	}

	return true, value.(uint)
}

func AuthDestroy(c *fiber.Ctx) {
	session, err := storedAuthenticationSession.Get(c)
	if err != nil {
		panic(err)
	}

	session.Delete("authentication")
	if err := session.Save(); err != nil {
		panic(err)
	}
}
