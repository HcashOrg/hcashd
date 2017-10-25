# Testing Summary of package `blockchain`

## Overview   
This is a report serves to summarize unit tests in this package. All utility functions assisting testing are placed in `sim_net_utils.go`, and relevant definitions are displayed in `sim_net_params.go`. For those reported as **FAILED** below, the targeted functions are all prefixed with `DNW` (Does Not Work).  

## Status  
### `chain`   
+ **FAILED**   
  - `TestBlockchainFunctions`  

### `fullblocks`   
+ **FAILED**  
  - `TestFullBlocks`   

### `fullblocksstakeversion`  
+ **FAILED**  
  - `TestStakeVersion`  

### `stakeversion`
+ **PASSED**    
+ **FAILED** due to the key block returned from the called routine `BlockChain::getPrevKeyNodeFromNode` is `nil`   
  - `TestCalcStakeVersionCorners`   
  - `TestCalcStakeVersionByNode`   
  - `TestIsStakeMajorityVersion`   
  - `TestLarge`  

### `thresholdstate`   
+ **FAILED**  
  - `TestThresholdState`: index out of range     

### `validate`  
+ **FAILED**   
  - `TestCheckWorklessBlockSanity`: no signature mismatch, lacking `*BlockChain`  
  - `TestCheckWorklessBlockSanity`: no signature mismatch, lacking `*BlockChain`     
  - `TestBlockValidationRules`: index out of range
