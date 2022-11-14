// Code generated by github.com/gnolang/gno. DO NOT EDIT.

//go:build gno
// +build gno

package nft

import (
	"github.com/gnolang/gno/stdlibs/stdshim"
	"strconv"

	"gno.tools/p/demo/avl"
	"gno.tools/p/demo/grc/grc721"
)

type token struct {
	grc721.IGRC721 // implements the GRC721 interface

	tokenCounter int
	tokens       *avl.Tree // grc721.TokenID -> *NFToken{}
	operators    *avl.Tree // owner std.Address -> operator std.Address
}

type NFToken struct {
	owner    std.Address
	approved std.Address
	tokenID  grc721.TokenID
	data     string
}

var gToken = &token{}

func GetToken() *token { return gToken }

func (grc *token) nextTokenID() grc721.TokenID {
	grc.tokenCounter++
	s := strconv.Itoa(grc.tokenCounter)
	return grc721.TokenID(s)
}

func (grc *token) getToken(tid grc721.TokenID) (*NFToken, bool) {
	_, token, ok := grc.tokens.Get(string(tid))
	if !ok {
		return nil, false
	}
	return token.(*NFToken), true
}

func (grc *token) Mint(to std.Address, data string) grc721.TokenID {
	tid := grc.nextTokenID()
	newTokens, _ := grc.tokens.Set(string(tid), &NFToken{
		owner:   to,
		tokenID: tid,
		data:    data,
	})
	grc.tokens = newTokens
	return tid
}

func (grc *token) BalanceOf(owner std.Address) (count int64) {
	panic("not yet implemented")
}

func (grc *token) OwnerOf(tid grc721.TokenID) std.Address {
	token, ok := grc.getToken(tid)
	if !ok {
		panic("token does not exist")
	}
	return token.owner
}

// XXX not fully implemented yet.
func (grc *token) SafeTransferFrom(from, to std.Address, tid grc721.TokenID) {
	grc.TransferFrom(from, to, tid)
	// When transfer is complete, this function checks if `_to` is a smart
	// contract (code size > 0). If so, it calls `onERC721Received` on
	// `_to` and throws if the return value is not
	// `bytes4(keccak256("onERC721Received(address,address,uint256,bytes)"))`.
	// XXX ensure "to" is a realm with onERC721Received() signature.
}

func (grc *token) TransferFrom(from, to std.Address, tid grc721.TokenID) {
	caller := std.GetCallerAt(2)
	token, ok := grc.getToken(tid)
	// Throws if `_tokenId` is not a valid NFT.
	if !ok {
		panic("token does not exist")
	}
	// Throws unless `msg.sender` is the current owner, an authorized
	// operator, or the approved address for this NFT.
	if caller != token.owner && caller != token.approved {
		_, operator, ok := grc.operators.Get(token.owner.String())
		if !ok || caller != operator.(std.Address) {
			panic("unauthorized")
		}
	}
	// Throws if `_from` is not the current owner.
	if from != token.owner {
		panic("from is not the current owner")
	}
	// Throws if `_to` is the zero address.
	if to == "" {
		panic("to cannot be empty")
	}
	// Good.
	token.owner = to
}

func (grc *token) Approve(approved std.Address, tid grc721.TokenID) {
	caller := std.GetCallerAt(2)
	token, ok := grc.getToken(tid)
	// Throws if `_tokenId` is not a valid NFT.
	if !ok {
		panic("token does not exist")
	}
	// Throws unless `msg.sender` is the current owner,
	// or an authorized operator.
	if caller != token.owner {
		_, operator, ok := grc.operators.Get(token.owner.String())
		if !ok || caller != operator.(std.Address) {
			panic("unauthorized")
		}
	}
	// Good.
	token.approved = approved
}

// XXX make it work for set of operators.
func (grc *token) SetApprovalForAll(operator std.Address, approved bool) {
	caller := std.GetCallerAt(2)
	newOperators, _ := grc.operators.Set(caller.String(), operator)
	grc.operators = newOperators
}

func (grc *token) GetApproved(tid grc721.TokenID) std.Address {
	token, ok := grc.getToken(tid)
	// Throws if `_tokenId` is not a valid NFT.
	if !ok {
		panic("token does not exist")
	}
	return token.approved
}

// XXX make it work for set of operators
func (grc *token) IsApprovedForAll(owner, operator std.Address) bool {
	_, operator2, ok := grc.operators.Get(owner.String())
	if !ok {
		return false
	}
	return operator == operator2.(std.Address)
}
