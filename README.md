# Hcash

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/HcashOrg/hcashd)

## Introduction

Blockchain technique has been thoroughly supported by researchers in both fields of distributed system and cryptography, regarding its delicate consensus model, scalability, efficiency, security, robustness, as well as privacy properties.
After the testimony of the market over eight years, blockchain has been considered as a dependable technique of public awareness. 

Apart from the blockchain, DAG technique has also been leveraged in few existing cryptocurrencies, merited for its overwhelming throughput, rapid transaction confirmation (especially in case of massive participation), by replacing the blockchain with a directed acyclic graph of transactions.
However, DAG-based public-chain systems have been well supported by neither a rigid and convincing theoretical analysis nor a sophisticated trial of the market. Moreover, problematic issues exist in DAG-based cryptocurrencies, like IOTA's security properties' high dependence on the transaction frequency, and Byteball's dependence on few ''witness'' nodes which leads to a potential centralization.
To avoid risks of DAG-based techniques, we adopt blockchain as the basis of our Hcash's underlying consensus scheme.


In 2008, for the first time in history, a fully decentralized monetary network based on a novel technique known as the blockchain was proposed by Satoshi Nakamoto in his documentation description of bitcoin.
After bitcoin's proposal, various decentralized cryptocurrencies have been theoretically proposed or developed. Also, massive enhancements to the blockchain have been devised and put into practice of existing cryptocurrencies.
Some of the most attractive innovations include Ethereum (by supporting smart contracts with the extension of the underlying script),
CryptoNotes, ZeroCash (by enhancing privacy via ring signatures or non-interactive zero-knowledge proofs), DASH, Decred (by a hybrid consensus and the implementation of a basic DAO), IOTA, Byteball (by improving the throughput with a DAG structure), Bitcoin-NG (by an novel primitive of a blockchain of two levels), and ''sidechain'' techniques (providing a linkability among different cryptocurrencies).


Cryptocurrencies implementing the proof-of-work (PoW), which is the foundation of bitcoin's consensus scheme, have been merited for many advantages, including a trustworthy sustainability (after sustaining sufficient practical trials), strong robustness against certain malicious participants, delicate incentive-compatibility, and a support of participants' dynamical joining and leaving.
On the other hand, PoW has also been indifferent due to its waste of resources and a potential centralization of hash power.
Therefore, alternative consensus primitives have been introduced to replace (fully or partially) PoW, like proof-of-stake (PoS).
However, pure PoS (a total substitution from PoW to PoS) is also controversial for its sustainability, due to its lack of practical trials and the risk of time-stamp forgery.


Another disadvantage of bitcoin is its limited throughput of transactions. Specifically, merely 7 transactions can be tackled by the existing network per second. That is, approximately only 600 thousand transactions can be tackled every day. This is a notorious drawback of bitcoin's future participation. Current improvement to this include:

1.  shortening the block interval;
2.  extending the block size;
3.  a chain structure of multiple levels;
4.  introducing ''the lightning network'';
5.  a ledger of a DAG framework.

Within them:

1. compromises a certain stability, which has been proven by ETH's short block interval (20s to 30s). To face this issue, a controversial protocol called ''GHOST'' has to be introduced to ETH.
2. causes a great communication burden to the network.
3. is utilized in Bitcoin-NG, its main principle is: the block proposed by miners after solving a hash puzzle is called a ''keyblock''. After the proposal of a keyblock, the proposer can release few ''microblock''s, before other miner's proposal of another keyblock. However, Bitcoin-NG is also problematic for its vulnerability facing a selfish mining and a potential attack by keyblock proposer's spawning massive microblocks, in which case the convergence property is undermined by chain forks, and the network is overloaded. 
4. provides efficiency, targeting at transactions with small value and high frequency.
5. extends the throughput significantly, while not well supported by a convincing theoretical proof or a long-term practice. Also, certain security risks exist in 5).

## Technology Summary

To date, pure PoW (without a hybrid of other consensus primitives like PoS) and a hybrid of PoW and PoS (like DASH and Decred) are still mainstreams of existing cryptocurrencies. However, the efficiency and throughput are limited in these mechanisms. 
Hcash aims at a secure, efficient, robust and hence dependable public chain cryptocurrency. Moreover, new characters will be released step by step by us including post-quantum techniques, smart contracts, a strong protection of wallets, and a linkability among blockchain based and DAG based distributed ledgers, etc.
To achieve this, two major characteristics are provided ahead in this release of Hcash:

1. With an alternative double-chain framework of the blockchain (differs from that of Bitcoin-NG), significant improvement of the efficiency is offered without compromising the security.
2. With a hybrid consensus, miners and stakeholders are incentivized to take part in the consensus, thereby enhancing the security and flexibility, and providing a mechanism that supports DAO for future protocol updating and project investments.

Adopting merits of both Decred and Bitcoin-NG, we devise a brand new consensus mechanism consisting of their key methodologies and our newly proposed innovations.
Firstly, with the methodology from Bitcoin-NG's keyblock/microblock structure, we offer a novel chain structure consisting two levels of blocks. In this way, the aforementioned hazard concerning the convergence property is faced. Specifically, the difficulty of the PoW hash puzzle is replaced by two difficulties. During miner's solving a hash puzzle, a keyblock can be proposed once a harder difficulty is met, and a microblock in case of a lower difficulty. In this way, the malicious power has to solve enough puzzles to spawn microblocks, hence the throughput and efficiency is enhanced without a compromise of security or robustness.
Moreover, to face the selfish mining, strengthen the robustness (against ''the 51\% attack''), and offer the flexibility (supporting both PoW and PoS mining), we adopt Decred's mechanism of voting tickets (a practical and flexible PoS scheme).
With the combination of two innovations above, a novel hybrid consensus scheme is formed, where keyblocks should be confirmed by certain voting tickets and both PoW and PoS miners participate in the consensus and play an important role.
Based on such a hybrid scheme, we implement a DAO that offers PoW and PoS miners' future decisions concerning protocol updating and project investments.
Meanwhile, our scheme supports ''the segregated witness'', which facilitates future post-quantum signature schemes and the lightning network.
The framework of our newly devised blockchain structure is shown in the figure below. 

<p align="center">
	<img src ="pic/Figure1.jpg" />
	<br/>
	Figure 1: The blockchain structure in our scheme
	<br/>
</p>
     

In the following table, comparisons are made between Hcash and few other existing distributed ledgers, including throughputs of Hcash with different intervals of block generations. The current release corresponds to the line marked with the bold font.

|             |Keyblock Average Interval|Block Size|Microblock Average Interval|Transaction Size|Throughput      |
|:----------- |:-----------------------:|:--------:|:-------------------------:|:--------------:|:--------------:|
|BTC          | 10min                   | 1MB      |                           | 250B           | 6.99TPS        |
|BTC(extended)| 10min                   | 2MB      |                           | 250B           | 13.98TPS       |
|BCC          | 10min                   | 8MB      |                           | 250B           | 55.92TPS       |
|Decred       | 5min                    |1.25MB    |                           | 250B           | 17.48TPS       |
|__Hcash__    | __5min__                | __2MB__  |    __18.75 sec__          | __250B__       | __447.39TPS__  |
|Hcash        | 5min                    | 8MB      |      18.75 sec            | 250B           | 1789.57TPS     |
<p align="center">
	<br/>
	Table 1: Comparisons between Hcash and other distributerd ledgers 
	<br/>
</p>

The following table offers an approximated chance of adversary's successfully undermining the system (in case of different PoW power and PoS capabilities, α: proportion of adversary's hash power, β: proportion of adversary's PoS capabilities).

<p align="center">
	<img src ="pic/Table2.png" />
	<br/>
	Table 2:  Probability of adversary's attack success with α fraction of total hash power and β fraction of total stake
	<br/>
</p>

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
