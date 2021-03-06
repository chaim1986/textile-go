package db

import (
	"database/sql"
	"strconv"
	"sync"
	"time"

	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/photos"
)

type PhotoDB struct {
	modelStore
}

func NewPhotoStore(db *sql.DB, lock *sync.Mutex) repo.PhotoStore {
	return &PhotoDB{modelStore{db, lock}}
}

func (c *PhotoDB) Put(set *repo.PhotoSet) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into photos(cid, lastCid, album, name, ext, created, added, latitude, longitude, local) values(?,?,?,?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}

	localInt := 0
	if set.IsLocal {
		localInt = 1
	}

	defer stmt.Close()
	_, err = stmt.Exec(
		set.Cid,
		set.LastCid,
		set.AlbumID,
		set.MetaData.Name,
		set.MetaData.Ext,
		int(set.MetaData.Created.Unix()),
		int(set.MetaData.Added.Unix()),
		set.MetaData.Latitude,
		set.MetaData.Longitude,
		localInt,
	)
	if err != nil {
		tx.Rollback()
		log.Errorf("error in db exec: %s", err)
		return err
	}
	tx.Commit()
	return nil
}

func (c *PhotoDB) GetPhotos(offsetId string, limit int, query string) []repo.PhotoSet {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offsetId != "" {
		q := ""
		if query != "" {
			q = query + " and "
		}
		stm = "select * from photos where " + q + "added<(select added from photos where cid='" + offsetId + "') order by added desc limit " + strconv.Itoa(limit) + " ;"
	} else {
		q := ""
		if query != "" {
			q = "where " + query + " "
		}
		stm = "select * from photos " + q + "order by added desc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *PhotoDB) GetPhoto(cid string) *repo.PhotoSet {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from photos where cid='" + cid + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *PhotoDB) DeletePhoto(cid string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.db.Exec("delete from photos where cid=?", cid)
	return nil
}

func (c *PhotoDB) handleQuery(stm string) []repo.PhotoSet {
	var ret []repo.PhotoSet
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var cid, lastCid, album, name, ext string
		var createdInt, addedInt int
		var latitude, longitude float64
		var localInt int
		if err := rows.Scan(&cid, &lastCid, &album, &name, &ext, &createdInt, &addedInt, &latitude, &longitude, &localInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		created := time.Unix(int64(createdInt), 0)
		added := time.Unix(int64(addedInt), 0)
		local := false
		if localInt == 1 {
			local = true
		}
		photo := repo.PhotoSet{
			Cid:     cid,
			LastCid: lastCid,
			AlbumID: album,
			MetaData: photos.Metadata{
				Name:      name,
				Ext:       ext,
				Created:   created,
				Added:     added,
				Latitude:  latitude,
				Longitude: longitude,
			},
			IsLocal: local,
		}
		ret = append(ret, photo)
	}
	return ret
}
