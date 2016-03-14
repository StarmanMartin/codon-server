package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/starmanmartin/simple-router"
	"github.com/starmanmartin/codon-resarch/ctools"
    
)

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
		session.Values["clist"] = ctools.ShiftLeft(list)
	    session.Save(r.Request, w)
    default:
		w.Write([]byte("Error"))
	}

	return true, nil
}
