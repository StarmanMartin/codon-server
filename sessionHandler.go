package main

import (
	"net/http"    
    "fmt"
    "strings"
    "errors"
    "github.com/gorilla/sessions"
)

var (
    store = sessions.NewCookieStore([]byte("something-very-secret"))
)

const (
    codonlistKey = "clist"
    sessionName = "session-name"
)

func loadList(r *http.Request) (string, error) {
    session, err := store.Get(r, sessionName)
    if err != nil {
		return "", err
	}
    
	oldList, has := session.Values[codonlistKey]
	if !has {
		return "", errors.New("no List")
	}
    
	sOldList := fmt.Sprint(oldList)
	sOldList = strings.Trim(sOldList, " ")
	return sOldList, nil
}

func loadArrayList(r *http.Request) ([]string, error) {
    list, err := loadList(r)
    if err != nil {
        return make([]string, 0), err
    }
    
    return strings.Split(list, " "), err
}

func saveList(r *http.Request, w http.ResponseWriter, newList string) (error){
    session, err := store.Get(r, sessionName)
    if err != nil {
		return err
	}
    
	session.Values[codonlistKey] = newList
    
    return session.Save(r, w)
}

func saveAsArrayList(r *http.Request, w http.ResponseWriter, newList []string) (error){
    return saveList(r, w, strings.Join(newList, " "))
}

func removeList(r *http.Request, w http.ResponseWriter) error{
    session, _ := store.Get(r, sessionName)
	delete(session.Values, codonlistKey)
	return session.Save(r, w)
}