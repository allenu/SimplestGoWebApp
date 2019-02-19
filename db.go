// db.go
package simplestgowebapp

import (
    "golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
    "time"
)

const postEntityName = "Post"

// What we store in Google Datastore
type DatabasePost struct {
    PostId string
    Created time.Time
    StorageId string      // UUID used for the Google Storage payloads/<UUID>.json file
}

func ReadPosts(appEngineContext context.Context, sinceTime time.Time, numPosts int) ([]DatabasePost, error) {
    q := datastore.NewQuery(postEntityName)
    q = q.Filter("Created < ", sinceTime).Order("-Created").Limit(numPosts)

    children := make([]DatabasePost, 0, numPosts)
    if _, err := q.GetAll(appEngineContext, &children); err == nil {
        return children, nil
    } else {
        return children, err
    }
}

func InsertPost(appEngineContext context.Context, postId string, post DatabasePost) error {
    key := datastore.NewKey(appEngineContext, postEntityName, postId, 0, nil)
    _, err := datastore.Put(appEngineContext, key, &post)

    return err
}

func ReadPost(appEngineContext context.Context, postId string) (DatabasePost, error) {
    key := datastore.NewKey(appEngineContext, postEntityName, postId, 0, nil)

    var post DatabasePost
    err := datastore.Get(appEngineContext, key, &post)
    if err != nil {
        return post, err
    }
    return post, nil
}

