//ftl:module httpingress
package httpingress

import (
	"context"
	"fmt"

	"ftl/builtin"
)

type GetRequest struct {
	UserID string `alias:"userId"`
	PostID string `alias:"postId"`
}

type Nested struct {
	GoodStuff string `alias:"good_stuff"`
}

type GetResponse struct {
	Message string `alias:"random"`
	Nested  Nested `alias:"nested"`
}

// Example: curl -i http://localhost:8892/ingress/http/users/123/posts?postId=456
//
//ftl:verb
//ftl:ingress http GET /http/users/{userID}/posts
func Get(ctx context.Context, req builtin.HttpRequest[GetRequest]) (builtin.HttpResponse[GetResponse], error) {
	return builtin.HttpResponse[GetResponse]{
		Headers: map[string][]string{"Get": {"Header from FTL"}},
		Body: GetResponse{
			Message: fmt.Sprintf("Got userId %s and postId %s", req.Body.UserID, req.Body.PostID),
			Nested:  Nested{GoodStuff: "Nested Good Stuff"},
		},
	}, nil
}

type PostRequest struct {
	UserID int `alias:"user_id"`
	PostID int `alias:"post_id"`
}

type PostResponse struct {
	Success bool `alias:"success"`
}

// Example: curl -i --json '{"user_id": 123, "post_id": 345}' http://localhost:8892/ingress/http/users
//
//ftl:verb
//ftl:ingress http POST /http/users
func Post(ctx context.Context, req builtin.HttpRequest[PostRequest]) (builtin.HttpResponse[PostResponse], error) {
	return builtin.HttpResponse[PostResponse]{
		Status:  201,
		Headers: map[string][]string{"Post": {"Header from FTL"}},
		Body:    PostResponse{Success: true},
	}, nil
}

type PutRequest struct {
	UserID string `alias:"userId"`
	PostID string `alias:"postId"`
}

type PutResponse struct{}

// Example: curl -X PUT http://localhost:8892/ingress/http/users/123 -d '{"postID": "123"}'
//
//ftl:verb
//ftl:ingress http PUT /http/users/{userID}
func Put(ctx context.Context, req builtin.HttpRequest[PutRequest]) (builtin.HttpResponse[PutResponse], error) {
	return builtin.HttpResponse[PutResponse]{
		Headers: map[string][]string{"Put": {"Header from FTL"}},
		Body:    PutResponse{},
	}, nil
}

type DeleteRequest struct {
	UserID string `alias:"userId"`
}

type DeleteResponse struct{}

// Example: curl -X DELETE http://localhost:8892/ingress/http/users/123
//
//ftl:verb
//ftl:ingress http DELETE /http/users/{userID}
func Delete(ctx context.Context, req builtin.HttpRequest[DeleteRequest]) (builtin.HttpResponse[DeleteResponse], error) {
	return builtin.HttpResponse[DeleteResponse]{
		Headers: map[string][]string{"Put": {"Header from FTL"}},
		Body:    DeleteResponse{},
	}, nil
}

type HtmlRequest struct{}

//ftl:verb
//ftl:ingress http GET /http/html
func Html(ctx context.Context, req builtin.HttpRequest[HtmlRequest]) (builtin.HttpResponse[string], error) {
	return builtin.HttpResponse[string]{
		Headers: map[string][]string{"Content-Type": {"text/html; charset=utf-8"}},
		Body:    "<html><body><h1>HTML Page From FTL 🚀!</h1></body></html>",
	}, nil
}

// Example: curl -X POST http://localhost:8892/ingress/http/bytes -d 'Your data here'
//
//ftl:verb
//ftl:ingress http POST /http/bytes
func Bytes(ctx context.Context, req builtin.HttpRequest[[]byte]) (builtin.HttpResponse[[]byte], error) {
	return builtin.HttpResponse[[]byte]{
		Headers: map[string][]string{"Content-Type": {"application/octet-stream"}},
		Body:    req.Body,
	}, nil
}
