# Fetch camera bills API

## Overview

The project provide HTTP endpoint that returns billing information for cameras within a specified date range and save result to csv file.

## Setup

1. Clone the project
2. Use `go mod tidy` to install all dependencies
3. Set environment variable in `.env` file
   - `HOST` : database host
   - `POSTGRES_PORT` : database port
   - `POSTGRES_USER` : database username
   - `POSTGRES_PASSWORD` : database password
   - `DB_NAME` : database name

## API

### Get camera bill data

- HTTP POST request: `http://localhost:8080/`
- cURL :
  ```sh
  curl --location 'http://localhost:8080/' \
  --header 'Content-Type: application/json' \
  --data '{
      "start" : "30-10-2023 00:00:00",
      "end" : "29-07-2024 00:00:00",
      "timeZone" : "7"
  }'
  ```
- Body:

  - `start`: The start date and time in the format `dd-MM-yyyy hh:mm:ss`. Example: `30-10-2023 00:00:00`
  - `end` : The end date and time in the format `dd-MM-yyyy hh:mm:ss`. Example: `29-07-2024 00:00:00`
  - `timeZone`: The time zone offset from UTC. Example: `-7` is `UTC-7` and `7` is `UTC+7`

- This endpoint fetch the data from the database using this query:

  ```sh
  select transaction_id, purchase_date_time, camera_sn, package_type
  from cam_bills
  where transaction_id is not null
  and payment_method = 'VIETTELPAY'
  and purchase_date_time > $1 AND purchase_date_time < $2
  and package_type in (select code from package_service where period > 3 and expired > 2595000)
  ```

  $1 and $2 are replaced with start and end date

- The query result will be saved to a csv file in csv folder. The date time in file name will be in UTC+0
