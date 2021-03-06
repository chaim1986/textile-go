package db

import (
	"database/sql"
	"path"
	"sync"
	"time"

	_ "github.com/mutecomm/go-sqlcipher"
	"github.com/op/go-logging"

	"github.com/textileio/textile-go/repo"
)

var log = logging.MustGetLogger("db")

type SQLiteDatastore struct {
	config   repo.Config
	settings repo.ConfigurationStore
	photos   repo.PhotoStore
	albums   repo.AlbumStore
	db       *sql.DB
	lock     *sync.Mutex
}

func Create(repoPath, password string) (*SQLiteDatastore, error) {
	var dbPath string
	dbPath = path.Join(repoPath, "datastore", "mainnet.db")
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if password != "" {
		p := "pragma key='" + password + "';"
		conn.Exec(p)
	}
	l := new(sync.Mutex)
	sqliteDB := &SQLiteDatastore{
		config: &ConfigDB{
			db:   conn,
			lock: l,
			path: dbPath,
		},
		settings: NewConfigurationStore(conn, l),
		photos:   NewPhotoStore(conn, l),
		albums:   NewAlbumStore(conn, l),
		db:       conn,
		lock:     l,
	}

	return sqliteDB, nil
}

func (d *SQLiteDatastore) Ping() error {
	return d.db.Ping()
}

func (d *SQLiteDatastore) Close() {
	d.db.Close()
}

func (d *SQLiteDatastore) Config() repo.Config {
	return d.config
}

func (d *SQLiteDatastore) Settings() repo.ConfigurationStore {
	return d.settings
}

func (d *SQLiteDatastore) Photos() repo.PhotoStore {
	return d.photos
}

func (d *SQLiteDatastore) Albums() repo.AlbumStore {
	return d.albums
}

func (d *SQLiteDatastore) Copy(dbPath string, password string) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	var cp string
	stmt := "select name from sqlite_master where type='table'"
	rows, err := d.db.Query(stmt)
	if err != nil {
		log.Errorf("error in copy: %s", err)
		return err
	}
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		tables = append(tables, name)
	}
	if password == "" {
		cp = `attach database '` + dbPath + `' as plaintext key '';`
		for _, name := range tables {
			cp = cp + "insert into plaintext." + name + " select * from main." + name + ";"
		}
	} else {
		cp = `attach database '` + dbPath + `' as encrypted key '` + password + `';`
		for _, name := range tables {
			cp = cp + "insert into encrypted." + name + " select * from main." + name + ";"
		}
	}

	_, err = d.db.Exec(cp)
	if err != nil {
		return err
	}

	return nil
}

func (d *SQLiteDatastore) InitTables(password string) error {
	return initDatabaseTables(d.db, password)
}

func initDatabaseTables(db *sql.DB, password string) error {
	var sqlStmt string
	if password != "" {
		sqlStmt = "PRAGMA key = '" + password + "';"
	}
	sqlStmt += `
	PRAGMA user_version = 0;
	create table config (key text primary key not null, value blob);
	create table photos (cid text primary key not null, lastCid text, album text not null, name text not null, ext text not null, created integer, added integer not null, latitude real, longitude real, local integer not null);
    create index index_album_added on photos (album, added);
	create index index_album_local on photos (album, local);
    create table albums (id text primary key not null, key blob not null, mnemonic text not null, name text not null);
    create unique index index_name on albums (name);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

type ConfigDB struct {
	db   *sql.DB
	lock *sync.Mutex
	path string
}

func (c *ConfigDB) Init(password string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return initDatabaseTables(c.db, password)
}

func (c *ConfigDB) Configure(creationDate time.Time) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert or replace into config(key, value) values(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec("creationDate", creationDate.Format(time.RFC3339))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *ConfigDB) GetCreationDate() (time.Time, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var t time.Time
	stmt, err := c.db.Prepare("select value from config where key=?")
	if err != nil {
		return t, err
	}
	defer stmt.Close()
	var creationDate []byte
	err = stmt.QueryRow("creationDate").Scan(&creationDate)
	if err != nil {
		return t, err
	}
	return time.Parse(time.RFC3339, string(creationDate))
}

func (c *ConfigDB) IsEncrypted() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	pwdCheck := "select count(*) from sqlite_master;"
	_, err := c.db.Exec(pwdCheck) // Fails if wrong password is entered
	if err != nil {
		return true
	}
	return false
}
