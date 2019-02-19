// main.go
package simplestgowebapp

import (
    "bytes"
    "encoding/json"
    "github.com/google/uuid"
    "google.golang.org/appengine"
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "path"
    "strings"
    "time"
)

type PostRequest struct {
    Content string
}

type Post struct {
    Title string
}

type ReadResponse struct {
    DatabasePost
    Post
}

type DrawViewModel struct {
    PostId string
}

type ListResponse struct {
    Posts []DatabasePost
    EndOfList bool
}

func init() {
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/r/", readHandler)

    http.HandleFunc("/api/write", apiWriteHandler)
    http.HandleFunc("/api/read/", apiReadHandler)
    http.HandleFunc("/api/list/", apiListHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    fp := path.Join("templates", "welcome.html")
    if tmpl, err := template.ParseFiles(fp); err == nil {
        tmpl.Execute(w, nil)
    } else {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func readHandler(w http.ResponseWriter, r *http.Request) {
    urlPath := "/r/"
    postId := strings.TrimPrefix(r.URL.Path, urlPath)

    model := DrawViewModel{
        PostId: postId,
    }

    fp := path.Join("templates", "read.html")
    if tmpl, err := template.ParseFiles(fp); err == nil {
        tmpl.Execute(w, model)
    } else {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func apiListHandler(w http.ResponseWriter, r *http.Request) {
    urlPath := "/api/list/"
    dateStr := strings.TrimPrefix(r.URL.Path, urlPath)

    sinceTime := time.Now()

    if len(dateStr) > 0 {
        layout := "2006-01-02T15:04:05.999999Z"
        t, err := time.Parse(layout, dateStr)
        if err == nil {
            sinceTime = t
        } else {
            log.Printf("Couldn't detect time from %s - err: %s", dateStr, err.Error())
        }
    }

    appContext := appengine.NewContext(r)
    clientFetchCount := 5
    actualFetchCount := clientFetchCount + 1
    posts, err := ReadPosts(appContext, sinceTime, actualFetchCount)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Bad request"))
        w.Write([]byte(err.Error()))
        return
    } else {

        endOfList := (len(posts) < actualFetchCount)
        // Figure out how many posts to return to caller
        var returnCount int
        posts_len := len(posts)
        if posts_len == 0 {
            returnCount = 0
        } else if posts_len == actualFetchCount {
            returnCount = clientFetchCount
        } else {
            returnCount = posts_len
        }

        response := ListResponse {
            Posts: posts[:returnCount],
            EndOfList: endOfList,
        }
        js, _ := json.Marshal(response)
        w.Header().Set("Content-Type", "application/json")
        w.Write(js)
    }
}

func apiReadHandler(w http.ResponseWriter, r *http.Request) {
    urlPath := "/api/read/"
    postId := strings.TrimPrefix(r.URL.Path, urlPath)

    appContext := appengine.NewContext(r)

    databasePost, postErr := ReadPost(appContext, postId)
    if postErr != nil {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("Content not found"))
        return
    }

    post, err := ReadContentFromStore(appContext, databasePost.StorageId)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("Content not found"))
        return
    }

    response := ReadResponse{
        DatabasePost: databasePost,
        Post: post,
    }

    js, _ := json.Marshal(response)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Cache-Control", "max-age=2592000") // 30 days
    w.Write(js)
}

func apiWriteHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Unsupported request method\n"))
        return
    }

    // Read all of the body into a string. This is ugly, but we decode it into the UpdateRequest
    // and into a map[] so that we can read out the properties that were provided.
    defer r.Body.Close()
    body, readErr := ioutil.ReadAll(r.Body)
    if readErr != nil {
        w.Write([]byte("Error reading body\n"))
        return
    }

    decoder := json.NewDecoder(bytes.NewReader(body))

    var postRequest PostRequest
    err := decoder.Decode(&postRequest)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Poorly formatted request\n"))
    } else {
        // TODO: Validation: Ensure it's actually a real post
        var post Post
        postDecoder := json.NewDecoder(strings.NewReader(postRequest.Content))
        postErr := postDecoder.Decode(&post)
        if postErr != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Poorly formatted Post\n"))
            return
        }

        // Validate it's good content
        if len(post.Title) > 100 {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Content was too long\n\n"))
            return
        }

        // Validate the content here
        if post.Title == "" {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Poorly formatted Post -- bad Json\n"))
            return
        }

        // Write to Store first
        appContext := appengine.NewContext(r)
        storageId, err := UploadContentToStore(appContext, postRequest.Content)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Now create the database entry
        postId := uuid.New().String()
        created := time.Now()
        databasePost := DatabasePost{
            PostId: postId,
            Created: created,
            StorageId: storageId,
        }
        err = InsertPost(appContext, postId, databasePost)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        } else {
            post := ReadResponse{
                DatabasePost: databasePost,
                Post: post,
            }

            js, _ := json.Marshal(post)
            w.Header().Set("Content-Type", "application/json")
            w.Write(js)
        }
    }
}

