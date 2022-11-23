# SleepStats
 
This is the backend service for the final project in my statistics class. It's a simple CRUD api that is also capable of calculating the final summary of the data.

The null hypothesis for the project was:
> 34.8% of U.S. adults sleep for less than 7 hours per night.

## Requirements 
- Postgres

## My deployment

I deployed this app as a docker container running on a Raspberry Pi 4 8gb. The 'example-docker-compose.yml' file provides the database as a container. The only necessary environment variable needed for this backend service is included in the '.env.example' file.

## Load Testing

Running a k6 loadtest against the 'summary' endpoint (requiring a call to the database) hosted on a Raspberry Pi 4 8gb resulted in:
- About 1200 simultaneous users before any threshold failures (given p(99) < 2 seconds)
  