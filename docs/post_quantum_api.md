# Hcashd Post quantum API

## Bliss algorithm

#### location: crypto/bliss/bliss.go 

```
func (sp blissDSA) GenerateKey(rand io.Reader) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey,
	error) {
	return sp.generateKey(rand)
}
func (sp blissDSA) Sign(priv hcashcrypto.PrivateKey, hash []byte) (hcashcrypto.Signature, error) {
	return sp.sign(priv, hash)
}
func (sp blissDSA) Verify(pub hcashcrypto.PublicKey, hash []byte, sig hcashcrypto.Signature) bool {
	return sp.verify(pub, hash, sig)
}
```

## LMS algorithm

####location: crypto/lms/lms.go

```
func (sp lmsDSA) GenerateKey(rand io.Reader) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey,
	error) {
	return sp.generateKey(rand)
}
func (sp lmsDSA) Sign(priv hcashcrypto.PrivateKey, hash []byte) (hcashcrypto.Signature, error) {
	return sp.sign(priv, hash)
}
func (sp lmsDSA) Verify(pub hcashcrypto.PublicKey, hash []byte, sig hcashcrypto.Signature) bool {
	return sp.verify(pub, hash, sig)
}
```





# Post quantuam library API

## Bliss Algorithm

#### Github: https://github.com/LoCCS/bliss

#### location: key.go

````
func GeneratePrivateKey(version int, entropy *sampler.Entropy) (*PrivateKey, error) {
  ...
}
````

#### location: sign.go

```
func GeneratePrivateKey(version int, entropy *sampler.Entropy) (*PrivateKey, error) {
  ...
}

func (key *PublicKey) Verify(msg []byte, sig *Signature) (bool, error) { 
  ...
}
```



## LMS Algorithm

#### Github: https://github.com/LoCCS/lms

It is different from Bliss or ECDSA that LMS doesn't hava GeneratePrivateKey function.

#### location: lms.go

```
func Sign(agent *MerkleAgent, hash []byte) (*lmots.PrivateKey, *MerkleSig, error) {
...
}

func Verify(root []byte, hash []byte, merkleSig *MerkleSig) bool { 
...
}
```

