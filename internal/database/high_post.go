package database

import (
	"database/sql"
	"escapade/internal/models"
	"time"

	//
	_ "github.com/lib/pq"
)

// CreateThread handle thread creation
func (db *DataBase) CreatePost(posts []models.Post, slug string) (createdPosts []models.Post, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	createdPosts = []models.Post{}

	var thatThread models.Thread

	if thatThread, err = db.threadFindByIDorSlug(tx, slug); err != nil {
		return
	}

	t := time.Now()
	for _, post := range posts {
		// if returnForum, err = db.postConfirmUnique(tx, forum); err != nil {
		// 	return
		// }

		if post.Author, err = db.userCheckID(tx, post.Author); err != nil {
			return
		}

		if post, err = db.postCreate(tx, post, thatThread, t); err != nil {
			return
		}

		createdPosts = append(createdPosts, post)
	}
	err = tx.Commit()
	return
}

func (db *DataBase) GetPosts(slug string, limit int, existLimit bool, t time.Time, existTime bool, sort string, desc bool) (returnPosts []models.Post, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	// if _, err = db.findForumBySlug(tx, slug); err != nil {
	// 	err = re.ErrorForumNotExist()
	// 	return
	// }

	var thatThread models.Thread

	if thatThread, err = db.threadFindByIDorSlug(tx, slug); err != nil {
		return
	}

	if sort == "flat" {
		if returnPosts, err = db.postsGetFlat(tx, thatThread, slug, limit, existLimit, t, existTime, desc); err != nil {
			return
		}
	}

	err = tx.Commit()
	return
}
