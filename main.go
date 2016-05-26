package main

import (
	"encoding/json"
	"fmt"
	"log"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
    "errors"
	"github.com/starmanmartin/codon-resarch"
	"github.com/starmanmartin/codon-resarch/ctools"
	"github.com/starmanmartin/codon-resarch/cThree"
	"github.com/starmanmartin/goconfig"
	"github.com/starmanmartin/simple-router"
	"github.com/starmanmartin/simple-router/view"
	"github.com/starmanmartin/simple-router/request"
)

var indexTmpl, errTmpl *template.Template

func resetHandler(w http.ResponseWriter, r *request.Request) (isNext bool, err error) {
	removeList(r.Request, w)
	w.Write([]byte("sucess"))

	return false, nil
}

func permutateListHandler(w http.ResponseWriter, r *request.Request) (isNext bool, err error) {
	r.ParseForm()
   
	list, err := loadArrayList(r.Request)
	if err != nil {
		return false, err
	}
    
    
	rule := r.Form.Get("rule")
	list, err = ctools.PermutateCodons(list, rule)
	if err != nil {
		return false, err
	}

	saveAsArrayList(r.Request, w, list)

	return true, nil
}

func uploadNewListHandler(w http.ResponseWriter, r *request.Request) (isNext bool, err error) {
	err = removeList(r.Request, w)
	return true, err
}

func uploadHandler(w http.ResponseWriter, r *request.Request) (isNext bool, err error) {
    r.ParseForm()
    
    oldList, err := loadList(r.Request)
	codon := r.Form.Get("list")
	var list []string
	if err != nil {
		if len(codon) <= 2 {
			w.Write([]byte("Empty"))
			return false, nil
		}
		list = strings.Split(codon, " ")
	} else {
		sOldList := fmt.Sprint(oldList)
		sOldList = strings.Trim(sOldList, " ")

		list = strings.Split(sOldList, " ")
		if len(codon) > 0 {
			if _, has := indexOf(list, codon); !has {
				return false, errors.New("codon dublicated")
			}

			list = append(list, codon)
			codon = fmt.Sprintf("%s %s", oldList, codon)
			codon = strings.Trim(codon, " ")
		} else {
            codon = sOldList
        }
	}

	graph, err := resarch.NewCodonGraph(list)

	if err != nil {
		return false, err
	}

	saveList(r.Request, w, codon)
	graph.OrderNodes()
	graph.FindIfCircular()
	graph.IsSelfComplementary()
	myjson, err := json.Marshal(graph)

	w.Write([]byte(myjson))

	return false, nil
}

func removeHandler(w http.ResponseWriter, r *request.Request) (isNext bool, err error) {
	r.ParseForm()
    
	oldList, err := loadList(r.Request)
	if err != nil {
		return false, err
	}


	codon := r.Form.Get("list")
	sOldList := fmt.Sprint(oldList)

	if len(codon) == 3 {
		codon = strings.Replace(sOldList, codon, "", -1)
		codon = strings.Trim(codon, " ")
		codon = strings.Replace(codon, "  ", " ", -1)
	} else {
		return false, errors.New("No codon")
	}

	list := strings.Split(codon, " ")
	graph, err := resarch.NewCodonGraph(list)

	if err != nil {
		return false, err
	}

	saveList(r.Request, w, codon)
	graph.OrderNodes()
	graph.FindIfCircular()
	graph.IsSelfComplementary()
	myjson, err := json.Marshal(graph)

	w.Write([]byte(myjson))

	return false, nil
}

func indexHandler(w http.ResponseWriter, r *request.Request) (isNext bool, err error) {
	indexTmpl.ExecuteTemplate(w, "base", "")
	return
}

func errorFunc(err error, w http.ResponseWriter, r *request.Request) {
	errTmpl.ExecuteTemplate(w, "base", err)
}

func notFound(w http.ResponseWriter, r *request.Request) {
	errTmpl.ExecuteTemplate(w, "base", "404")
}

func errorFuncXHR(err error, w http.ResponseWriter, r *request.Request) {
	w.Write([]byte("Error"))
}

func notFoundXHR(w http.ResponseWriter, r *request.Request) {
	w.Write([]byte("Error"))
}


func iniWebRouter() {
	cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if err := goconfig.InitConficOnce(cwd+"/config/config.json", cwd+"/config/param.json"); err != nil {
		panic(err)
		return
	}
    
	app := router.NewRouter()
	app.Public("/public")
	app.Get("/", indexHandler)
    
    app = router.NewXHRRouter()
	app.Post("/newgraph", uploadHandler)
	app.Post("/removecodon", removeHandler)
	app.Post("/reset", resetHandler)
	app.Post("/newlist", uploadNewListHandler, uploadHandler)
	app.Post("/permutate", permutateListHandler, uploadHandler)
	
    initRouter()
	app.Post("/check/*", uploadHandler)
	
    router.ErrorHandler = errorFunc
	router.NotFoundHandler = notFound
    router.XHRErrorHandler = errorFuncXHR
	router.XHRNotFoundHandler = notFoundXHR

	view.ViewPath = "views"
	indexTmpl = view.ParseTemplate("index", "index.html")
	errTmpl = view.ParseTemplate("error", "error.html")
	port, _ := goconfig.GetString("port")
	log.Println("Listening on port:" , port)
	log.Println(cThree.GetSingleSwitchPath())
	
	http.ListenAndServe(":"+port, app)
	
	
}

func main() {
	iniWebRouter()
}
