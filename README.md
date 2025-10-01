# CSV Email Checker

A simple Go backend that:

- Uploads a CSV file, parses it, and appends a `flag` column (`true` if any field is a valid email, else `false`).
- Stores the processed CSV and allows downloading by ID.

---

## Endpoints

### 1. Upload

POST /API/upload
Form field: file (CSV file)


**Response (200 OK):**
```json
{
  "id": "a225eb00-0907-4273-92ca-5faadeefae5f"
}

Error :

{
  "error": "{failure-reason}"
}

```

### 2. Download

GET /API/download/{id}

- 200 OK → returns processed CSV as a file
- 423 Locked → job still in progress
- 400 Bad Request → invalid ID


## Run Locally
```
cd csv-email-checker
go mod tidy
go run main.go
```

- Server runs on http://localhost:8080


### Results

```
ashish@DESKTOP-56M6SBL:/csv-email-checker$ echo "name,email,age
ABC,abc@example.com,25
PQR,pqrexample.com,30
XYZ,xyz@example.com,28" > sample.csv

ashish@DESKTOP-56M6SBL:/csv-email-checker$ cat sample.csv
name,email,age
ABC,abc@example.com,25
PQR,pqrexample.com,30
XYZ,xyz@example.com,28

ashish@DESKTOP-56M6SBL:/csv-email-checker$ curl -X POST -F "file=@sample.csv" http://localhost:8080/API/upload
{"id":"5b1acef4-2e7e-4fa2-8a73-ac04ecccbb94"}

ashish@DESKTOP-56M6SBL:/csv-email-checker$ curl -X GET http://localhost:8080/API/download/3e4e8a4d-18ac-4b65-8797-6c2e2e3c5d1a -o processed.csv
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    23  100    23    0     0   4126      0 --:--:-- --:--:-- --:--:--  4600

ashish@DESKTOP-56M6SBL:/csv-email-checker$ cat processed.csv
{"error":"invalid id"}

ashish@DESKTOP-56M6SBL:/csv-email-checker$ curl -X GET http://localhost:8080/API/download/5b1acef4-2e7e-4fa2-8a73-ac04ecccbb94 -o processed.csv
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   104  100   104    0     0   6677      0 --:--:-- --:--:-- --:--:--  6933


ashish@DESKTOP-56M6SBL:/csv-email-checker$ cat processed.csv
name,email,age,flag
ABC,abc@example.com,25,true
PQR,pqrexample.com,30,false
XYZ,xyz@example.com,28,true

```
