# Walk Thru

## Dependencies 

**Stripe:**
[Stripe Go](github.com/stripe/stripe-go)
[Stripe Go Charge](github.com/stripe/stripe-go/charge)
[Stripe Go Form](github.com/stripe/stripe-go/form)
[Stripe Go Token](github.com/stripe/stripe-go/token)

**Glog:**
[Golang Glog](github.com/golang/glog)

**MUX:**
[Gorilla Mux](github.com/gorilla/mux)

**Postgresql:**
[PQ](github.com/lib/pq)
[PQ Oid](github.com/lib/pq/oid)
[PQ Scram](github.com/lib/pq/scram)

## Config
**Enviroment Variables:**
`WEB_HOST, WEB_PORT, DATABASE_HOST, DATABASE_PORT, DATABASE_USER, DATABASE_PASSWORD, DATABASE_DBNAME, STRIPE_KEY`
</br>
</br>
If Enviroment Variables are not exported in the working directory, the values can be fetched through the **config.json** file.
</br>
</br>
**Sample Config File:**
```json
{
  "dbHost": "postgres",
  "dbPort": "5432",
  "dbUser": "postgres",
  "dbPassword": "StrongPassword",
  "dbName": "postgres",
  "webHost": "",
  "webPort": "8000",
  "stripeKey": "sk_test_##############"
}
```

## Docker
Project uses the `golang:alpine` image.
</br></br>
To run docker image:
</br>
`docker build -t webserver .` 
</br>
`docker run -p "8000:8000" webserver`
</br>
</br>
Where `8000` is the WEB_PORT.

## Docker-Compose
To run Docker Compose:
</br>
`docker-compose build`
</br>
`docker-compose up`

