# Cryptography

## Overview  
This section mainly deals with documentation about cryptography in hcashd project, including  
+ unit tests  
+ ...  

## Unit Tests  
Progress  
+ [x] conversion among big integers, encoded bytes, field elements and curve points in `hcashec/edwards/primitives_test.go`   
+ [x] ECDH simulation, encryption/decryption on elliptic curve, error handling in `hcashec/edwards/chiphering_test.go`   
+ [x] benchmarking of signing/verifying in `hcashec/edwards/ecdsa_benchmark_test.go`   
+ [x] signing/verifying on bad pubkeys and signatures, validates signing/verifying against golden implementation, check key serialization/deserialization in `hcashec/edwards/ecdsa_test.go`  
+ [x] curve points addition/recovery/multiplication in `hcashec/edwards/curve_test.go`  
+ [ ] Schnorr threshold signature scheme in `hcashec/edwards/schnorr_threshold_test.go`   
  - test on bad public keys ain't understood
  - test on bad public nonces ain't understood
