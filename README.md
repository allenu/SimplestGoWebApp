# SimplestGoWebApp

This is a basic Google Cloud Platform app written in Golang that illustrates the basics of a CRUD app:
- a simple front-end that lets you enter text and via JavaScript fetch() requets, hits two endpoints for writing and reading
- two REST endpoints for reading and writing
- one REST endpoint for getting the latest N posts (a date can be provided to list the last N posts from there,
  making it possible to paginate)
- Google Storage to store data
- how to process and validate JSON in golang

Use this project as a starting point for your next simple Golang web project!

If you have questions, message me or tweet me at https://twitter.com/ussherpress.

## Setup

- Create a new Google Cloud Platform project using the GCP console.
  https://console.cloud.google.com/

- While in GCP, navigate to Storage on the side and create a bucket. Pick a unique name.

- Generate GOOGLE_APPLICATION_CREDENTIALS json file. You'll need this to test locally. Without it, you will not be
  able to access the Google Storage buckets remotely.
  - Follow the instructions here (i.e. visit "Create service account key")
    https://cloud.google.com/docs/authentication/production
  - You may also find it under API & Services/Credentials/Create Credentials/Service account key.
    You might need to create a new service account first. (Name the service account name "storage-reader-writer"
    and pick Storage Object Creator and Storage Object Viewer from the Role dropdown.
  - Download the JSON file it generates. This is essentially a password for accessing your Google
    Storage bucket, so don't share it.

- Edit app.yaml
  - Enter a unique session secret. (Just come up with a long phrase. This will be used to encrypt the session data stored
    in the client browser.)
  - Enter a bucket_name. This is the bucket name you created earlier. If you forgot it, you can find out the bucket
    name using GCP and select Storage in the side menu.

## Testing it out

Run this to debug locally (provide the path to the JSON file downloaded earlier)
```
    export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your/GoogleCredentials.json
    dev_appserver.py app.yaml --require_indexes=yes
```

# TODO

- [ ] Add delete feature
- [ ] Add basic users
- [ ] Make this truly a SPA app. Clicks on /r/ links should just present a lightbox.

