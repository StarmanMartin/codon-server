package main

import (
	"net/http"

	"github.com/starmanmartin/codon-resarch/ctools"
	"github.com/starmanmartin/simple-router"
)

func initRouter() {
	sub := router.NewSubRouter("/check")
	sub.Post("/:tool", handleTool)
}

func handleTool(w http.ResponseWriter, r *router.Request) (isNext bool, err error) {
	r.ParseForm()
	
    list, err := loadArrayList(r.Request)
	if err != nil {
		w.Write([]byte("Error"))
		return false, nil
	}


	switch r.RouteParams["tool"] {
	case "shift":
        saveList(r.Request, w, ctools.ShiftLeft(list))
	case "shiftcodon":
        saveList(r.Request, w, ctools.ShiftCodonLeft(list))
	case "fill_comp":
        saveList(r.Request, w, ctools.FillComlements(list))
	case "remove_comp":
        saveList(r.Request, w, ctools.RemoveComlements(list))
    case "shuffle":
        saveList(r.Request, w, ctools.Shuffle(list))
	default:
		w.Write([]byte("Error"))
	}

	return true, nil
}
