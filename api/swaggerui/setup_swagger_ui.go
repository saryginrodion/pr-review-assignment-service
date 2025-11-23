package swaggerui

import (
	"net/http"
)

const swaggerHTML = `<!DOCTYPE html>
<html>
<head>
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
<script>
window.onload = () => {
  SwaggerUIBundle({
    url: "/openapi.yaml",
    dom_id: "#swagger-ui"
  });
};
</script>
</body>
</html>`

func SetupSwaggerUI() {
	http.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(swaggerHTML))
	})

	http.Handle("GET /openapi.yaml", http.FileServer(http.Dir("./openapi/")))
}
