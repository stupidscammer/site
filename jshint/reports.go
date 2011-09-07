package reports;

import (
    "appengine"
    "appengine/datastore"
    "http"
    "time"
    "template"
)

type Snippet struct {
    Code string
    Opts string
    Date datastore.Time
}

var templates = make(map[string]*template.Template)

func init() {
    templates["report"] = template.New(nil)
    err := templates["report"].ParseFile("jshint/templates/report.html")
    if err != nil {
        panic("Can't parse template report.html")
    }

    for _, tmpl := range []string{"500", "404"} {
        templates[tmpl] = template.New(nil)
        templates[tmpl].SetDelims("[", "]")
        err = templates[tmpl].ParseFile("jshint/templates/" + tmpl + ".html")
        if err != nil {
            panic("Can't parse template " + tmpl + ".html")
        }
    }

    http.HandleFunc("/reports/save/", save)
    http.HandleFunc("/reports/", show)
}

func renderTemplate(w http.ResponseWriter, name string, snippet *Snippet) {
    err := templates[name].Execute(w, snippet)
    if err != nil {
        http.Error(w, err.String(), 500)
    }
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
        c.Criticalf(err.String())
        renderTemplate(w, "500", nil)
        return
    }

    w.Header().Set("Location", "/reports/" + key.Encode())
    w.WriteHeader(http.StatusFound)
}

const lenPath = len("/reports/")

func show(w http.ResponseWriter, r *http.Request) {
    var snippet Snippet
    c := appengine.NewContext(r)

    key, err := datastore.DecodeKey(r.URL.Path[lenPath:])
    if err != nil {
        c.Criticalf(err.String())
        renderTemplate(w, "404", nil)
        return
    }

    datastore.Get(c, key, &snippet)
    renderTemplate(w, "report", &snippet)
}