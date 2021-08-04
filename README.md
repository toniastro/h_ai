## How to run the application
There are 2 main services i.e (Api and consumer) service(s).
Clone the application with 
`
git clone https://github.com/toniastro/h_ai.git
`

### Docker

```
cd h_ai
docker-compose up --build
```
> Note: By default the port number its being run on is **8100**.


### Languages/Frameworks used
- Golang
- Mongo DB
- RabbitMQ
- Redis

## Endpoints Description

### Create a Job

This makes a request to persist a job to be consumed at a later time.

```
    URL - /api/job
    Method - POST
    Response - (content-type = application/json)
``` 
```JSON
    {
       "object_id": "12223333"
    }
```

### Get Job Status

This makes a request to check the details of a job alongside its job status with the Job Id returned from creation of Job

```
    URL - /api/job/{job_id}
    Method - GET
    Response - (content-type = application/json)
``` 
```JSON
    {
      "data": {
        "_id": "610aa4b237f9575426dbaf71",
        "created_at": "2021-08-04T14:31:14.691Z",
        "updated_at": "2021-08-04T14:31:14.691Z",
        "status": "COMPLETED",
        "time_taken": 24,
        "object_id": "10",
        "job_id": "585dc2cb-a5bc-40af-8209-610ddb632975"
      },
      "message": "Job Details",
      "status": 200
    }
```