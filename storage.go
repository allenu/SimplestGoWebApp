package simplestgowebapp

import (
    "cloud.google.com/go/storage"
    "context"
    "encoding/json"
    "fmt"
    "github.com/google/uuid"
    "google.golang.org/appengine/file"
    "io/ioutil"
    "os"
    "strings"
)

func UploadContentToStore(ctx context.Context, content string) (string, error) {
    bucketName := os.Getenv("bucket_name")
    if bucketName == "" {
        defaultBucketName, _ := file.DefaultBucketName(ctx)
        bucketName = defaultBucketName
    }

    client, err := storage.NewClient(ctx)

    if err != nil {
        return "", err
    } else {
        bkt := client.Bucket(bucketName)

        postId := uuid.New().String()
        obj := bkt.Object("payloads/" + postId + ".json")

        objectWriter := obj.NewWriter(ctx)
        objectWriter.ContentType = "application/json"

        if _, err := fmt.Fprintf(objectWriter, content); err != nil {
            return "", err
        }

        // Close, just like writing a file.
        if err := objectWriter.Close(); err != nil {
            return "", err
        }

        return postId, nil
    }
}

func ReadContentFromStore(ctx context.Context, postId string) (Post, error) {
    bucketName := os.Getenv("bucket_name")
    if bucketName == "" {
        defaultBucketName, _ := file.DefaultBucketName(ctx)
        bucketName = defaultBucketName
    }

    client, err := storage.NewClient(ctx)

    var post Post

    if err == nil {
        bkt := client.Bucket(bucketName)

        obj := bkt.Object("payloads/" + postId + ".json")

        objectReader, err := obj.NewReader(ctx)
        if err == nil {
            defer objectReader.Close()

            slurp, err := ioutil.ReadAll(objectReader)

            if err == nil {
                rawString := string(slurp)

                postDecoder := json.NewDecoder(strings.NewReader(rawString))
                err = postDecoder.Decode(&post)
            }
        }
    }

    return post, err
}

