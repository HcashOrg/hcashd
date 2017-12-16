# Hcash  

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/HcashOrg/hcashd)

## Contents  
+ [Official Web](#web)  
+ [Innovations](#innovations)  
	- [Novel Consensus Scheme](#novel-consensus-scheme)  
	- [Post-Quantum Features](#post-quantum-features)  
+ [Starting-Hcashd](#starting-hcashd)  
+ [Using HcashWallet UI Version](#using-hcashwallet-ui-version)  
+ [License](#license)  

<a name="web"></a>
## Official Web
https://h.cash

<a name="innovations"></a>
## Innovations 
To date, existing decentralized cryptocurrencies adopt either PoW consensus scheme or hybrid consensus model of PoW and PoS. However, these systems still encounter the issue of **very limited efficiency/throughput**. Meanwhile, upcoming quantum computers threaten existing classical cryptography which is the foundation of **blockchain security**. In particular, the quantum algorithm by [Shor](http://www.jstor.org/stable/2653075?seq=1#page_scan_tab_contents) for computing discrete logarithms breaks the ECDSA signature scheme used by almost all cryptocurrencies, such as Bitcoin, Ethereum, Decred and Monero. However, if post-quantum cryptographic schemes are equipped in these systems, the throughput of them will become worse and even unbearable.

Hcash project aims to build a **secure**, **efficient**, **robust** and **reliable** decentralized system. Highlighted features such as **newly-proposed hybrid consensus scheme**, **post-quantum digital signature**, **linkability** among various blockchain-based and DAG-based decentralized cryptocurrencies, **smart contract mechanism** and **post-quantum privacy-preserving scheme** will be proposed and implemented in Hcash eventually.

<a name="novel-consensus-scheme"></a>
### Novel Consensus Scheme  
To deal with the performance issue, we implement a novel hybrid consensus scheme with **strong robustness**, **high throughput** as well as **sufficient flexibility** in Hcash. On the one hand, with a newly-proposed two-layer framework of block chain, significant improvement of the efficiency is offered without compromising the security. On the other hand, with a hybrid consensus model, both PoW and PoS miners are incentivized to take part in the consensus process, thereby enhancing the security and flexibility of the consensus scheme, and providing a mechanism that supports basic DAO for future protocol updating and project investments.

For more details, please refer to our [specific report](docs/research/design-rationale-of-the-consensus-scheme-in-hcash.md).

<a name="post-quantum-features"></a>
### Post-Quantum Features  
To address security issues stemming from quantum computers, we design and implement post-quantum solutions in Hcash. Our proposals achieve the following 4 features:  
+ **Compatibility**: Compatible with existing ECDSA signature solution;   
+ **Flexibility**: Support multiple post-quantum signature solutions that are thoroughly analyzed, assessed and proved by international cryptography research institutions, meanwhile their security and performance must be outstanding;  
+ **Security**: the post-quantum solution must be proved secure in theory, and side-channel attack proof in practice;  
+ **High performance**: Signing and signature verification must be fast. Most importantly, the public key and signature must be short.

Please refer to our [detailed report](docs/research/design-rationale-of-post-quantum-features-in-hcash.md) for more information.

<a name="starting-hcashd"></a>
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

* Start running hcash full node service to synchrnoze blocks
	```
	hcashd
	```

* Start hcash solo mining
	```
	hcashctl setgenerate true x     # where x represents the number of CPU threads
	```

* Stop hcash solo mining
	```
	hcashctl setgenerate false
	```
<a name="using-hcashwallet-ui-version"></a>
## Using HcashWallet UI Version

HcashWallet UI version is a graphical wallet for HCASH. With it, you can send and receive HCASH, purchase tickets for PoS voting, get a history of all your transactions and more.

HcashWallet UI version is located here: https://github.com/HcashOrg/hcashwallet/releases. It could be extracted and used directly.

<a name="license"></a>
## License  
hcashd is licensed under the [copyfree](http://copyfree.org) ISC License.
