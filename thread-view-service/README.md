# Thread View Service

This service is responsible for providing the UI to view, create, and edit threads.

## Endpoints

| Method | Path | Description |
| ------ | ---- | ----------- |
| GET | /thread/list | Get all threads |
| GET | /thread/get | Get a form to search for a thread by it's ID. |
| POST | /thread/getid | Shows a single thread by it's ID. |
| GET | /thread/create | Get a form to create a new thread. |
| POST | /thread/create | Takes the form input and sends it to the view-service api. |
