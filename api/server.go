package api

import (
	"net/http"

	"github.com/saryginrodion/pr_review_assignment_service/api/context"
	"github.com/saryginrodion/pr_review_assignment_service/api/routes"
	"github.com/saryginrodion/stackable"
)

func HttpServer(stack stackable.Stackable[context.SharedState, context.LocalState], addr string) *http.Server {
	http.Handle("GET /", stack.AddUniqueHandler(routes.GetIndex))

	// Teams
	http.Handle("POST /team/add", stack.AddUniqueHandler(routes.TeamAdd))
	http.Handle("GET /team/get", stack.AddUniqueHandler(routes.TeamGet))

	// Users
	http.Handle("POST /users/setIsActive", stack.AddUniqueHandler(routes.UserSetIsActive))

	s := &http.Server{
		Addr: addr,
	}

	return s
}
