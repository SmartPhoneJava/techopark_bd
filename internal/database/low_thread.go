package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"fmt"
	"time"

	//
	_ "github.com/lib/pq"
)

// createThread create thread
func (db *DataBase) threadCreate(tx *sql.Tx, thread *models.Thread) (createdThread models.Thread, err error) {

	query := `INSERT INTO Thread(slug, author, created, forum, message, title) VALUES
						 	($1, $2, $3, $4, $5, $6) 
						 RETURNING id, slug, author, created, forum, message, title;
						 `
	row := tx.QueryRow(query, thread.Slug, thread.Author, thread.Created,
		thread.Forum, thread.Message, thread.Title)

	createdThread = models.Thread{}
	if err = row.Scan(&createdThread.ID, &createdThread.Slug,
		&createdThread.Author, &createdThread.Created, &createdThread.Forum,
		&createdThread.Message, &createdThread.Title); err != nil {
		return
	}
	return
}

/*
SELECT 	a.FieldWidth, a.FieldHeight,
					a.MinsTotal, a.MinsFound,
					a.Finished, a.Exploded
	 FROM Player as p
		JOIN
			(
				SELECT player_id,
					FieldWidth, FieldHeight,
					MinsTotal, MinsFound,
					Finished, Exploded
					FROM Game Order by id
			) as a
			ON p.id = a.player_id and p.name like $1
			OFFSET $2 Limit $3

*/

// getThreads get threads
func (db *DataBase) threadsGetWithLimit(tx *sql.Tx, slug string, limit int) (foundThreads []models.Thread, err error) {

	query := `select id, slug, author, created, forum, message, title from
							Thread where forum like $1 Limit $2;
						 `

	var rows *sql.Rows

	if rows, err = tx.Query(query, slug, limit); err != nil {
		return
	}
	defer rows.Close()

	foundThreads = []models.Thread{}
	for rows.Next() {
		thread := models.Thread{}
		if err = rows.Scan(&thread.ID, &thread.Slug,
			&thread.Author, &thread.Created, &thread.Forum,
			&thread.Message, &thread.Title); err != nil {
			break
		}

		foundThreads = append(foundThreads, thread)
	}
	return
}

func (db *DataBase) threadsGet(tx *sql.Tx, slug string, limit int, lb bool, t time.Time, tb bool, desc bool) (foundThreads []models.Thread, err error) {

	fmt.Println("threadsGet got:", t.String())
	query := `select id, slug, author, created, forum, message, title from
							Thread where lower(forum) like lower($1)`

	if tb {
		if desc {
			query += ` and created <= $2`
			query += ` order by created desc`
		} else {
			query += ` and created >= $2`
			query += ` order by created`
		}
		if lb {
			query += ` Limit $3`
		}
	} else if lb {
		if desc {
			query += ` order by created desc`
		} else {
			query += ` order by created`
		}
		query += ` Limit $2`
	}

	var rows *sql.Rows

	if tb {
		if lb {
			rows, err = tx.Query(query, slug, t, limit)
		} else {
			rows, err = tx.Query(query, slug, t)
		}
	} else if lb {
		rows, err = tx.Query(query, slug, limit)
	} else {
		rows, err = tx.Query(query, slug)
	}

	if err != nil {
		return
	}
	defer rows.Close()

	foundThreads = []models.Thread{}
	for rows.Next() {

		thread := models.Thread{}
		if err = rows.Scan(&thread.ID, &thread.Slug,
			&thread.Author, &thread.Created, &thread.Forum,
			&thread.Message, &thread.Title); err != nil {
			break
		}
		foundThreads = append(foundThreads, thread)
	}
	return
}

// checkUser checks, is thread's author exists
func (db *DataBase) threadCheckUser(tx *sql.Tx, thread *models.Thread) (err error) {
	var thatUser models.User
	if thatUser, err = db.findUserByName(tx, thread.Author); err != nil {
		err = re.ErrorUserNotExist()
		return
	}
	thread.Author = thatUser.Nickname
	return
}

// threadCheckForum checks, is thread's forum exists
func (db *DataBase) threadCheckForum(tx *sql.Tx, thread *models.Thread) (err error) {
	var thatForum models.Forum
	if thatForum, err = db.findForumBySlug(tx, thread.Forum); err != nil {
		err = re.ErrorForumNotExist()
		return
	}
	thread.Forum = thatForum.Slug
	return
}
