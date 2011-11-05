package hvadsigeralex

import (
    "http"
    "template"
    "rand"
    "strconv"
    "appengine"
    "strings"
    "hvadsigeralex/model"
)

func init() {
    http.HandleFunc("/", randomStatus)
    http.HandleFunc("/status/", singleStatus)
    http.HandleFunc("/primeCache", primeCache)
}

func singleStatus(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
	statusId, _ := strconv.Atoui64( strings.TrimLeft(r.URL.Path, "/status/") )
	c.Debugf("Trying to fetch status %s from cache", statusId)
	status, err := model.GetStatusById(c, statusId)
	if err == nil {
		renderStatus(w, status)
	} else {
		renderError(w, "Could not find status! <a href=\"/\">Reload me.<a>")
	}

}

func randomStatus(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
	statusList, statusErr := model.GetStatuses(c)
	if statusErr == nil && len(statusList) > 0 {
  	status := statusList[rand.Intn(len(statusList)-1)]
		renderStatus(w, status)
	} else {
		renderError(w, "No statuses found, sorry.")
	}
}

func renderStatus(w http.ResponseWriter, status model.Status) {
	data := map[string] string {
  	"bodyClass": "col"+strconv.Itoa(rand.Intn(4)),
  	"status": status.Message,
		"directLink": "<a href=\"/status/"+status.Id+"\">#"+status.Id+"</a>"}
	data["extraCSS"] = calcExtraCSS(data["status"])
	renderPage(w, data)
}

func renderError(w http.ResponseWriter, errMessage string) {
	data := map[string] string {
  	"bodyClass": "colError",
  	"status": errMessage,
		"directLink": ""}
	data["extraCSS"] = calcExtraCSS(data["status"])
	renderPage(w, data)
}

func renderPage(w http.ResponseWriter, data map[string] string) {
	mainPageTemplate, templateErr := template.ParseFile("hvadsigeralex/templates/index.html")
	if templateErr != nil {
		http.Error(w, templateErr.String(), http.StatusInternalServerError)
	}
	err := mainPageTemplate.Execute(w, data)
  if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
  }
}

func calcExtraCSS(text string) (string){
  switch x := len(text); {
    case x < 20: return "font-size: 2.1em";
    case x < 40: return "font-size: 2.0em";
    case x < 60: return "font-size: 1.5em";
    case x < 80: return "font-size: 1.4em";
    case x < 100: return "font-size: 1.25em";
    case x < 140: return "font-size: 1.1em";
    case x < 160: return "font-size: 1.05em";
    case x < 180: return "font-size: 1.0em";
    case x < 250: return "font-size: 0.85em";
  }
  return "font-size: 0.8em;"
}

func primeCache(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  err := model.ForceUpdateStatuses(c)
  if err != nil {
      http.Error(w, "An error occured while updating the cache", http.StatusInternalServerError)
  }
}
