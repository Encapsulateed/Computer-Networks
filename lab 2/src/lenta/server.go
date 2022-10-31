package main

import (
	log "github.com/mgutz/logxi/v1"
	"html/template"
	"net/http"
)

const INDEX_HTML = `
    <!doctype html>
    <html lang="ru">
        <head>
            <meta charset="utf-8">
            <title>ТОП 10 КРИПТОВОЛЮТ</title>
        </head>
        <body>
		{{if .}}
	        {{range .}}
  			
  			<div>
				<img src={{.Img}}>
				<a href= "{{.Ref}}" >{{.Title}} {{.StrVol}}</a>
				
			</div>

			<br/>
		{{end}}
{{else}}
Не удалось загрузить новости!
{{end}}
         
        </body>
    </html>
    `

var indexHtml = template.Must(template.New("index").Parse(INDEX_HTML))

func serveClient(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path

	log.Info("got request", "Method", request.Method, "Path", path)
	if path != "/" && path != "/index.html" {
		log.Error("invalid path", "Path", path)
		response.WriteHeader(http.StatusNotFound)
	} else if err := indexHtml.Execute(response, downloadNews()); err != nil {
		log.Error("HTML creation failed", "error", err)
	} else if err == nil {
		clear()
	} else {
		log.Info("response sent to client successfully")
	}
}

func main() {
	http.HandleFunc("/", serveClient)
	log.Info("starting listener")

	log.Error("listener failed", "error", http.ListenAndServe("127.0.0.1:0228", nil))
}
