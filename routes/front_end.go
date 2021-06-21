package routes

import (
	"html/template"
	"net/http"

	"github.com/Skjaldbaka17/quotes-api/structs"
)

func Home(rw http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.ParseFiles("./tmpl/home.html"))
	q, _ := getRandomQuoteFromDb(&structs.Request{})
	err := templates.ExecuteTemplate(rw, "home.html", q)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	// content, _ := ioutil.ReadFile("./templates/home.html")
	// t, _ := template.New("Home").Parse(string(content))
	// q, err := getRandomQuoteFromDb(&structs.Request{})
	// if err != nil {
	// 	rw.WriteHeader(http.StatusInternalServerError)
	// 	log.Printf("Got error when querying DB in GetRandomQuote: %s", err)
	// 	json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
	// 	return
	// }
	// t.Execute(os.Stdout, q)
	// rw.Header()
}
