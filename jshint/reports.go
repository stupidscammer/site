package reports;

import (
    "appengine"
    "appengine/datastore"
    "http"
    "time"
    "template"
    "strconv"
)

type Report struct {
    Code []byte
    Opts []byte
    Date datastore.Time
}

var templates = make(map[string]*template.Template)

func init() {
    templates["report"] = template.New(nil)
    if err := templates["report"].ParseFile("jshint/templates/report.html"); err != nil {
        panic("Can't parse report.html")
    }

    for _, tmpl := range []string{"500", "404"} {
        templates[tmpl] = template.New(nil)
        templates[tmpl].SetDelims("[", "]")
        if err := templates[tmpl].ParseFile("jshint/templates/" + tmpl + ".html"); err != nil {
            panic("Can't parse " + tmpl + ".html")
        }
    }

    http.HandleFunc("/reports/save/", save)
    http.HandleFunc("/reports/", show)
    http.HandleFunc("/", notFound)
}

func renderTemplate(w http.ResponseWriter, name string, report *Report) {
    err := templates[name].Execute(w, report)
    if err != nil {
        http.Error(w, err.String(), 500)
    }
}

func save(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    report := Report{
        Code: []byte(r.FormValue("code")),
        Opts: []byte(r.FormValue("data")),
        Date: datastore.SecondsToTime(time.Seconds()), // Now
    }

    key, err := datastore.Put(c, datastore.NewIncompleteKey("report"), &report)
    if err != nil {
        c.Criticalf(err.String())
        renderTemplate(w, "500", nil)
        return
    }

    w.Header().Set("Location", "/reports/" + strconv.Itoa64(key.IntID()))
    w.WriteHeader(http.StatusFound)
}

const lenPath = len("/reports/")

func show(w http.ResponseWriter, r *http.Request) {
    var report Report
    c := appengine.NewContext(r)

    intID, _ := strconv.Atoi64(r.URL.Path[lenPath:])
    key := datastore.NewKey("report", "", intID, nil)
    if err := datastore.Get(c, key, &report); err != nil {
        c.Criticalf(err.String())
        renderTemplate(w, "404", nil)
        return
    }
    renderTemplate(w, "report", &report)
}

func notFound(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "404", nil)
}