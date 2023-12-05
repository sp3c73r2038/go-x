package httpx

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"code\":0}"))
	})
	http.ListenAndServe(":8000", nil)
}
