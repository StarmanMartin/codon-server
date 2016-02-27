package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/starmanmartin/simple-router"
)

var (
	totla *regexp.Regexp
)

func init() {
	totla = regexp.MustCompile(`([AUGC]{3})`)
}

func initRouter() {
	sub := router.NewSubRouter("/check")
	sub.Post("/:tool", handleTool)
}

func handleTool(w http.ResponseWriter, r *router.Request) (isNext bool, err error) {
	r.ParseForm()
	session, _ := store.Get(r.Request, "session-name")
	oldList, has := session.Values["clist"]
	if !has {
		w.Write([]byte("Error"))
		return false, nil
	}

	sOldList := fmt.Sprint(oldList)
	sOldList = strings.Trim(sOldList, " ")
	list := strings.Split(sOldList, " ")

	switch r.RouteParams["tool"] {
	case "shift":
		tempStringList := strings.Join(list, "") + string(list[0][0])
		session.Values["clist"] = totla.ReplaceAllString(tempStringList[1:], "$1 ")
	    session.Save(r.Request, w)
    default:
		w.Write([]byte("Error"))
	}

	return true, nil
}
