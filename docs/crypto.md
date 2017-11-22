# Cryptography

## Overview  
This section mainly deals with documentation about cryptography in hcashd project, including  
+ unit tests: reports progress for tests under different packages    
+ ...  

## Unit Tests  
### package `hcashec/edwards`  
| file  | brief |  
|-----:|:-----|  
| primitives_test.go  | conversion among big integers, encoded bytes, field elements and curve points |  
| chiphering_test.go  | ECDH simulation, encryption/decryption on elliptic curve, error handling  |     
| ecdsa_benchmark_test.go | benchmarking of signing/verifying |   
| ecdsa_test.go | signing/verifying on bad pubkeys and signatures, validates signing/verifying against golden implementation, check key serialization/deserialization |    
| curve_test.go | curve points addition/recovery/multiplication |  
| threshold_schnorr_test.go | Schnorr threshold signature scheme  |     
