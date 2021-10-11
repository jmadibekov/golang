# Weather App

This weather app fetches the real-time current temperatures of cities Almaty, Nur-Sultan, Moscow, London and New-York and outputs them in JSON file.

## Running Locally

You need to have [Golang](https://golang.org/) (the version I'm using is 1.17.1) installed in your computer to run the app.

Run the following in the same directory where `weather.go` is located:

```sh
go run weather.go
```

Then it will take few seconds to fetch and output the data into `weather.json`. Note that the app outputs some commenting logs into stdout. 

## Packages used

- **`encoding/json`:** Used to decode and encode JSON into native types in Go.

- **`io/ioutil`:** Used to store the data into a file.

- **`log`:** Used to log comments (or possible errors) into stdout.

- **`net/http`:** Used to make HTTP requests to the external Weather API.

- **`time`:** Used to get the current localtime in the app.

Go 1.17.1 is the default language of the app.
