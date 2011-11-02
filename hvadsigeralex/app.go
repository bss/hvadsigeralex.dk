package hvadsigeralex

import (
    "http"
    "template"
    "rand"
    "strconv"
    "appengine"
    "hvadsigeralex/model"
)

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/primeCache", primeCache)
}

func root(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  
  var status string = ""
  statusList, statusErr := model.GetStatuses(c)
  if statusErr != nil {
    status = "..."
  } else if len(statusList) > 0 {
    status = statusList[rand.Intn(len(statusList)-1)].Message    
  }
  data := map[string] string {
      "bodyClass": "col"+strconv.Itoa(rand.Intn(5)),
      "status": status}
  data["extraCSS"] = calcExtraCSS(data["status"])
  
  err := mainPageTemplate.Execute(w, data)
  if err != nil {
      http.Error(w, err.String(), http.StatusInternalServerError)
  }
}

func calcExtraCSS(text string) (string){
  switch x := len(text); {
    case x < 20: return "font-size: 1.9em";
    case x < 40: return "font-size: 1.9em";
    case x < 60: return "font-size: 1.6em";
    case x < 80: return "font-size: 1.4em";
    case x < 100: return "font-size: 1.4em";
    case x < 140: return "font-size: 1.2em";
    case x < 160: return "font-size: 1.1em";
    case x < 180: return "font-size: 1.0em";
    case x < 250: return "font-size: 0.9em";
  }
  return "font-size: 0.8em;"
}

var mainPageTemplate = template.Must(template.New("MainPage").Parse(mainPageHTML))

const mainPageHTML = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Hvad siger Alex?</title>
  <link rel="stylesheet" href="/css/main.css" type="text/css" />
  <meta property="og:image" content="http://www.hvadsigeralex.dk/img/head_col.png"/> 
</head>
<body class="{{html .bodyClass}}">
  <div id="bubbleImg">&nbsp;
    <h1 style="{{html .extraCSS}}">{{html .status}}</h1>
    <a id="head" href="http://alexbp.dk">&nbsp;</a>
  </div>
  <div id="footer">Foto: <a href="http://www.gadang.dk">Frederik Holmgaard</a></div>
</body>
</html>
`

func primeCache(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  err := model.ForceUpdateStatuses(c)
  if err != nil {
      http.Error(w, err.String(), http.StatusInternalServerError)
  }
}
