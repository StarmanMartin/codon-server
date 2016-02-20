package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/starmanmartin/simple-router"
	"github.com/starmanmartin/simple-router/view"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

var indexTmpl, errTmpl *template.Template

func resetHandler(w http.ResponseWriter, r *router.Request) (isNext bool, err error) {
	session, _ := store.Get(r.Request, "session-name")
	delete(session.Values, "clist")
	session.Save(r.Request, w)
	w.Write([]byte("sucess"))

	return false, nil

}

func uploadNewListHandler(w http.ResponseWriter, r *router.Request) (isNext bool, err error) {
    session, _ := store.Get(r.Request, "session-name")
	delete(session.Values, "clist")
	session.Save(r.Request, w)
	return uploadHandler(w, r)
}

func uploadHandler(w http.ResponseWriter, r *router.Request) (isNext bool, err error) {
	r.ParseForm()
	session, _ := store.Get(r.Request, "session-name")
	codon := r.Form.Get("list")
	oldList, has := session.Values["clist"]
	var list []string
	if !has {
		list = strings.Split(codon, " ")
		session.Values["clist"] = codon
	} else {
		sOldList := fmt.Sprint(oldList)
		sOldList = strings.Trim(sOldList, " ")

		list = strings.Split(sOldList, " ")
		if len(codon) > 0 {
			if _, has = indexOf(list, codon); !has {
				w.Write([]byte("Error"))
				return false, nil
			}

			list = append(list, codon)
			codon = fmt.Sprintf("%s %s", oldList, codon)
			codon = strings.Trim(codon, " ")
			session.Values["clist"] = codon
		}
	}

	graph, err := NewCodonGraph(list)

	if err != nil {
		w.Write([]byte("Error"))
		return false, nil
	}

	session.Save(r.Request, w)
	graph.OrderNodes()
	graph.FindIfCircular()
	graph.IsSelfComplementary()
	myjson, err := json.Marshal(graph)

	w.Write([]byte(myjson))

	return false, nil
}

func removeHandler(w http.ResponseWriter, r *router.Request) (isNext bool, err error) {
	r.ParseForm()
	session, _ := store.Get(r.Request, "session-name")
	codon := r.Form.Get("list")
	oldList, has := session.Values["clist"]
	var list []string
	if !has {
		w.Write([]byte("Error"))
		return false, nil
	}

	sOldList := fmt.Sprint(oldList)

	if len(codon) == 3 {
		codon = strings.Replace(sOldList, codon, "", -1)
		codon = strings.Trim(codon, " ")
		codon = strings.Replace(codon, "  ", " ", -1)
		list = strings.Split(codon, " ")
		session.Values["clist"] = codon
	} else {
		w.Write([]byte("Error"))
		return false, nil
	}

	graph, err := NewCodonGraph(list)

	if err != nil {
		w.Write([]byte("Error"))
		return false, nil
	}

	session.Save(r.Request, w)
	graph.OrderNodes()
	graph.FindIfCircular()
	graph.IsSelfComplementary()
	myjson, err := json.Marshal(graph)

	w.Write([]byte(myjson))

	return false, nil
}

func indexHandler(w http.ResponseWriter, r *router.Request) (isNext bool, err error) {
	indexTmpl.ExecuteTemplate(w, "base", "")
	return
}

func errorFunc(err error, w http.ResponseWriter, r *router.Request) {
	errTmpl.ExecuteTemplate(w, "base", err)
}

func notFound(w http.ResponseWriter, r *router.Request) {
	errTmpl.ExecuteTemplate(w, "base", "404")
}

func iniWebRouter() {
	app := router.NewRouter()
	app.Public("/public")
	app.Get("/", indexHandler)
	app.Post("/newgraph", uploadHandler)
	app.Post("/removecodon", removeHandler)
	app.Post("/reset", resetHandler)
	app.Post("/newlist", uploadNewListHandler)

	router.ErrorHandler = errorFunc
	router.NotFoundHandler = notFound

	view.ViewPath = "views"
	indexTmpl = view.ParseTemplate("index", "index.html")
	errTmpl = view.ParseTemplate("error", "error.html")

	http.ListenAndServe(":8080", app)
}

func main() {
	iniWebRouter()
}
