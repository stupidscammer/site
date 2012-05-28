package reports

import (
	"appengine"
	"appengine/datastore"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Report struct {
	Code []byte
	Opts []byte
	Date time.Time
}

var templates = make(map[string]*template.Template)

func init() {
	for _, name := range []string{"500", "404", "report"} {
		tmpl := template.Must(template.New(name).ParseFiles("jshint/templates/" + name + ".html"))
		templates[name] = tmpl
	}

	http.HandleFunc("/reports/save/", save)
	http.HandleFunc("/reports/", show)
	http.HandleFunc("/", notFound)
}

func renderTemplate(w http.ResponseWriter, name string, report *Report) {
	err := templates[name].ExecuteTemplate(w, name+".html", report)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func save(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	report := Report{
		Code: []byte(r.FormValue("code")),
		Opts: []byte(r.FormValue("data")),
		Date: time.Now(),
	}

	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "report", nil), &report)
	if err != nil {
		c.Criticalf(err.Error())
		renderTemplate(w, "500", nil)
		return
	}

	w.Header().Set("Location", "/reports/"+strconv.FormatInt(key.IntID(), 10))
	w.WriteHeader(http.StatusFound)
}

const lenPath = len("/reports/")

func show(w http.ResponseWriter, r *http.Request) {
	var report Report
	c := appengine.NewContext(r)

	intID, _ := strconv.ParseInt(r.URL.Path[lenPath:], 10, 64)
	key := datastore.NewKey(c, "report", "", intID, nil)
	if err := datastore.Get(c, key, &report); err != nil {
		c.Criticalf(err.Error())
		renderTemplate(w, "404", nil)
		return
	}
	renderTemplate(w, "report", &report)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "404", nil)
}
