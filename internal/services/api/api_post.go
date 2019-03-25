package api

import (
	data "escapade/internal/database"
	"escapade/internal/models"
	"net/http"
)

func (h *Handler) CreatePosts(rw http.ResponseWriter, r *http.Request) {
	const place = "CreatePosts"
	var (
		posts []models.Post
		err   error
		slug  string
	)

	rw.Header().Set("Content-Type", "application/json")

	if slug, err = getSlug(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if posts, err = getPosts(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if len(posts) > 0 {
		if posts, err = h.DB.CreatePost(posts, slug); err != nil {
			//if err.Error() == re.ErrorForumUserNotExist().Error() {
			rw.WriteHeader(http.StatusNotFound)
			sendErrorJSON(rw, err, place)
			// } else {
			// 	rw.WriteHeader(http.StatusConflict)
			// 	sendSuccessJSON(rw, forum, place)
			// }
			printResult(err, http.StatusBadRequest, place)
			return
		}
	}
	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, posts, place)
	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) GetPosts(rw http.ResponseWriter, r *http.Request) {
	const place = "GetPosts"
	var (
		posts      []models.Post
		slug       string
		sort       string
		limit      int
		since      int
		err        error
		existLimit bool
		existSince bool
		desc       bool
		qgc        data.QueryGetConditions
	)

	rw.Header().Set("Content-Type", "application/json")

	if slug, err = getSlug(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if existLimit, limit, err = getLimit(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if existSince, since, err = getIDmin(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	desc = getDesc(r)
	sort = getSort(r)

	qgc.InitPost(existSince, since, existLimit, limit, desc)

	if posts, err = h.DB.GetPosts(slug, qgc, sort); err != nil {
		//if err.Error() == re.ErrorForumUserNotExist().Error() {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		// } else {
		// 	rw.WriteHeader(http.StatusConflict)
		// 	sendSuccessJSON(rw, forum, place)
		// }
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, posts, place)
	printResult(err, http.StatusOK, place)
	return
}
