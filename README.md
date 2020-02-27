# puffadder



Package `puffadder` is a project database built from scratch that is designed to store and update current bank accounts.



The main features that allow this are:

* Event source architecture so that you have a history of activity on each account.

* Choice of homemade algorithm (btree, hash_table, ordered_array).

* Transactions between accounts are atomic.

* All transactions are ATOMIC, therefore they are saved to disk and if the server crashes nothing is lost.



---



* [Install](#install)

* [Run](#run)

* [Examples](#examples)

* [Structure](#structure)

* [Internals](#internals)

* [Performance](#performance)


---



## Install



With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:



```sh

go get -u github.com/zoomie/puffadder

```

## Run

In the project direcotry run the following
```sh
go build .
./puffadder
```


## Examples

Make withdraw, add-money and transactions POST requests
```sh

curl 'localhost:8090/create-account?accountName=joe'
curl 'localhost:8090/add-money?accountName=joe&addAmount=100'
curl 'localhost:8090/view-current-account?accountName=joe' 
curl 'localhost:8090/withdraw-money?accountName=john&subtractAmount=10'
curl 'localhost:8090/transfet?fromAccount=joe&subtractAmount=10'
```

Or using one of the handlers directly:

```go

func main() {
	x := "add working example in here"
}
```

## Structure

The server and entry point for the database is found in:
* [main.go](main.go)

Algorithms that store value in memory:
* [btree.go](btree.go)
* [hash_table.go](hash_table.go)
* [ordered_array.go](ordered_array.go)

The file that persists the data to disk and updates the current value using an algorithm:
* [update_file_data.go](update_file_data.go)

## Internals

Show how the project works internally.

## Performance

Use the different algorithms in test.

Also think about setting up concurrency.
