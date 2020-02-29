# puffadder



Package `puffadder` is a project database built from scratch that is designed to store and update current bank accounts.



The main features that allow this are:

* Event source architecture so that you have a history of activity on each account.

* Choice of homemade algorithm (btree, hash_table, ordered_array).

* Transactions between accounts are atomic.

* The database is durable, it stores each action to a file (like a commit log in other databases)
  * Side note: I havn't looked into making sure it is flushed to permenated storeage like postgres (fsync)

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
curl 'localhost:8090/view-current-account?accountName=joe' 
curl 'localhost:8090/add-money?accountName=joe&addAmount=100'
curl 'localhost:8090/withdraw-money?accountName=john&subtractAmount=10'
curl 'localhost:8090/transfet?fromAccount=joe&toAccount&subtractAmount=10'
```

Test example of how to create an account in code, can be found in [server_test.go](server_test.go).

```go

func createAccountSetUp(accountName string) (*httptest.ResponseRecorder, *http.Request) {
	params := "?accountName=" + accountName
	url := path.Join(baseURL+"create-account", params)
	request, _ := http.NewRequest("POST", url, nil)
	recorder := httptest.NewRecorder()
	return recorder, request
}

func TestCreateAccount(t *testing.T) {
	setUp()
	recorder, request := createAccountSetUp("joe")
	createAccount(recorder, request)
	response := recorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Creating account failed")
	}
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

Tests:
* [algorithm_test.go](algorithm_test.go)
* [server_test.go](server_test.go)

## Internals

Show how the project works internally.

## Performance

Use the different algorithms in test.

Also think about setting up concurrency.
