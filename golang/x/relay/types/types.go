package types

import (
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
)

// Hash256Digest 32-byte double-sha2 digest
type Hash256Digest = btcspv.Hash256Digest

// Hash160Digest is a 20-byte ripemd160+sha2 hash
type Hash160Digest = btcspv.Hash160Digest

// RawHeader is an 80-byte raw header
type RawHeader = btcspv.RawHeader

// BitcoinHeader is a parsed Bitcoin header
type BitcoinHeader = btcspv.BitcoinHeader

// SPVProof is the base struct for an SPV proof
type SPVProof = btcspv.SPVProof

// // Link is a link in the chain
// type Link struct {
// 	Digest string `json:"digest"`
// 	Parent string `json:"parent"`
// }
