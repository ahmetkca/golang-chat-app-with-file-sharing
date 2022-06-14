# golang-chat-app-with-file-sharing
A Chat web application with file sharing capability

It uses MinIO for storing uploaded files. MinIO is a High Performance Object Storage. It is API compatible with Amazon S3 cloud storage service. I used `gocloud.dev/blob` golang package to avoid vendor-lock implementation which means instead of depending on AWS, GCP, or Azure `gocloud.dev` provides a common api to interface with any cloud provider.

Features to add in the future
-   Persistance chat history (even after stopping the server.)
-   Login/Register.
-   Preview for Uploaded File.
-   Group Chats.
-   Direct Messages (DM)

> Run `docker-compose up -d` to spin up a MinIO server for s3 compatible Object Storage

> Execute `setup_access_keys.sh` before starting the chat app to set the MinIO credentials