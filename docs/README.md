# `hcashd`: daemon program for `hcash`  

## Contents  
+ [mining](#mining)  

## Mining <a name="mining" />   
While working on some PoW, `hcashd` actually tries to build a block out of a prepared `BlockTemplate` consisting of

+ `Block`: block header + plain txs + stake txs, detailed as [Fig](#block-struct)    
+ `Fee`: tx fees vector paid by each tx in `Block`   
+ `SigOpCounts`: TBC   
+ `Height`: block height of current block   
+ `KeyHeight`: number of key blocks before current block   
+ `ValidPayAddress`: TBC   
+ `GenerateKey`: TBC    

![block structure][block-struct]
[block-struct]: images/block.png "block structure"  
