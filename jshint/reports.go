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
    Data string
    Date datastore.Time
}

func init() {
    http.HandleFunc("/reports/save/", save)
    http.HandleFunc("/reports/show/", show)
}

func save(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    snippet := Snippet{
        Code: r.FormValue("code"),
        Data: r.FormValue("data"),
        Date: datastore.SecondsToTime(time.Seconds()), // Now
    }

    key, _ := datastore.Put(c, datastore.NewIncompleteKey("snippet"), &snippet)
    fmt.Fprintf(w, key.Encode())
}

func show(w http.ResponseWriter, r *http.Request) {
    var snippet Snippet
    c := appengine.NewContext(r)

    key, err := datastore.DecodeKey(r.FormValue("k"))
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