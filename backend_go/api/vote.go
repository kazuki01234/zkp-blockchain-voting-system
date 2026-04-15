package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazuki01234/zkp-blockchain-voting-system/backend_go/blockchain"
	"github.com/kazuki01234/zkp-blockchain-voting-system/backend_go/crypto"
)

type VoteRequest struct {
	VoterPublicKey string `json:"voter_public_key"`
	VoteData       string `json:"vote_data"`
	Proof          string `json:"proof"`
	Signature      string `json:"signature"`
}

func VoteHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()

	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	voterHash := crypto.HashPublicKey(req.VoterPublicKey)
	message := req.Proof + "|" + voterHash

	if !crypto.VerifySignature(req.VoterPublicKey, message, req.Signature) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	proofValid := crypto.VerifyProof(req.Proof)

	if !proofValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ZKP proof"})
		return
	}

	if blockchain.BlockchainInstance.HasVoted(voterHash) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This public key has already voted"})
		return
	}

	tx := blockchain.Transaction{
		VoterHash:  voterHash,
		VoteData:   req.VoteData,
		Proof:     req.Proof,
		ProofValid: proofValid,
	}

	_, err := blockchain.BlockchainInstance.AddBlock([]blockchain.Transaction{tx})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote added to blockchain"})
}