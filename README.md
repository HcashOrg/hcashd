## Starting Hcashd
Hcashd is a Hypercash full node implementation written in Go (golang).

This acts as a chain daemon for the [Hypercash](https://h.cash) cryptocurrency. Hcashd maintains the entire past transactional ledger of Hypercash and allows relaying of transactions to other Hypercash nodes across the world.
The installation of hcashd requires Go 1.7 or newer.
* Glide 
	Glide is used to manage project dependencies and provide reproducible builds. To install:
	```
	go get -u github.com/Masterminds/glide
	```
* Build and Installation
	For a first time installation, the project and dependency sources can be obtained manually with git and glide (create directories as needed):
	```
	git clone https://github.com/HcashOrg/hcashd $GOPATH/src/github.com/HcashOrg/hcashd
	cd $GOPATH/src/github.com/HcashOrg/hcashd
	glide install
	go install $(glide nv)
	```
    To update an existing source tree, pull the latest changes and install the matching dependencies:
    ```
	cd $GOPATH/src/github.com/HcashOrg/hcashd
	git pull
	glide install
	go install $(glide nv)
    ```


## License

hcashd is licensed under the [copyfree](http://copyfree.org) ISC License.
