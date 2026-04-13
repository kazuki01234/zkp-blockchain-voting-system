package crypto

import "regexp"

func GenerateProof() string {
	panic("Proof generation is client-side only")
}

func VerifyProof(proof string) bool {
	matched, _ := regexp.MatchString("^[0-9a-f]{64}$", proof)
	return matched
}
