package reports;

import (
    "appengine"
    "appengine/datastore"
    "fmt"
    "http"
    "time"
    "template"
)

type Snippet struct {
    Code string
    Opts string
    Date datastore.Time
}

const lenPath = len("/reports/")

func init() {
    http.HandleFunc("/reports/save/", save)
    http.HandleFunc("/reports/", show)
}

func save(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    snippet := Snippet{
        Code: r.FormValue("code"),
        Opts: r.FormValue("data"),
        Date: datastore.SecondsToTime(time.Seconds()), // Now
    }

    key, err := datastore.Put(c, datastore.NewIncompleteKey("snippet"), &snippet)
    if err != nil {
        fmt.Fprintf(w, err.String())
    }
    w.Header().Set("Location", "/reports/" + key.Encode())
    w.WriteHeader(http.StatusFound)
}

func show(w http.ResponseWriter, r *http.Request) {
    var snippet Snippet
    c := appengine.NewContext(r)

    uid := r.URL.Path[lenPath:]
    key, err := datastore.DecodeKey(uid)
    if err != nil {
        fmt.Fprintf(w, err.String())
    }

    datastore.Get(c, key, &snippet)
    tmpl := template.MustParseFile("jshint/templates/report.html", nil)
    err = tmpl.Execute(w, snippet)
    if err != nil {
        http.Error(w, err.String(), 500);
    }
}