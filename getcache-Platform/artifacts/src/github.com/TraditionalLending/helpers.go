
/*--------------------------------------------------------------------------
----------------------------------------------------------------------------
   HELPER FUNCTIONS CALLED SEVERAL TIMES ON MAIN SMART CONTRACT FUNCIOTNS
----------------------------------------------------------------------------
-------------------------------------------------------------------------- */

package main
import (
	"errors"
	//"math"
	//"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/common/util"
	//pb "github.com/hyperledger/fabric/protos/peer"
)


/* -------------------------------------------------------------------------------------------------
checkTokenListed: this function if a token is already listed in a smart contract. Input is:
token           string           // Token to verify
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) checkTokenListed( stub shim.ChaincodeStubInterface,
	                                                        token string ) ( map[string]bool, error ) {
	// Retrieve token list from blockchain //
	token_list := map[string]bool{}
	tokenListBytes, err1 := stub.GetState( IndexTokenList )
	if err1 != nil {
		err := errors.New( "ERROR: FAILED TO GET THE TOKEN LENDING POOLS " + 
					       "INFORMATION." + err1.Error() )
		return token_list, err }
	json.Unmarshal(tokenListBytes, &token_list)

	// Check if the token is already listed //
	_, isInList := token_list[ token]
	if !isInList {
		err := errors.New( "ERROR: LENDING POOL FOR TOKEN" +  token + 
						   " IS NOT LISTED." ) 
		return token_list, err}
	return token_list, nil
}

/* -------------------------------------------------------------------------------------------------
retrieveTokenPool:  this function returns the state of the Lending Pool for the Token in question.
		            Args: array containing
args[0]               string   // Symbol of the token to release the information	   
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) retrieveTokenPool( stub shim.ChaincodeStubInterface,
															 token string ) (LendingPool, error) {

	lending_pool := LendingPool{}
	// Check that Token is registered on Blockchain //
	_, err1 := t.checkTokenListed( stub, token )
	if err1 != nil { return lending_pool, err1 }

	// Retrieve token lending pool information //
	poolsBytes, err2 := stub.GetState( IndexLendingPools + token )
	if err2 != nil {
		err := errors.New( "ERROR: FAILED TO GET THE TOKEN LENDING POOL " + 
						   "INFORMATION." + err2.Error() ) 
		return lending_pool, err }
	json.Unmarshal( poolsBytes, &lending_pool )
	return lending_pool, nil
}

/* -------------------------------------------------------------------------------------------------
retrieveUserWallet: this function returns the wallet of an User. Input is:
PublicId            string             // Public ID of the user to retrieve the wallet
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) retrieveUserWallet( stub shim.ChaincodeStubInterface, 
													   PublicId string) (MultiWallet, error) {

	// Retrieve wallet of an user from Blockchain //
	user_wallet := MultiWallet{}
	chainCodeArgs := util.ToChaincodeArgs( "balanceOf", PublicId )
	response := stub.InvokeChaincode( COIN_BALANCE_CHAINCODE, chainCodeArgs, 
									  CHANNEL_NAME )
	if response.Status != shim.OK {
		err := errors.New ( "ERROR INVOKING THE COINBALANCE CHAINCODE TO " +
							"GET THE BALANCE OF USER: " + PublicId ) 
		return user_wallet, err }
	json.Unmarshal(response.Payload, &user_wallet)
	return user_wallet, nil
}

/* -------------------------------------------------------------------------------------------------
updatePool: this function updates a Pool on Blockchain. Inputs: the updated pool.
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) updatePool( stub shim.ChaincodeStubInterface, 
													  pool LendingPool) ( []byte, error ) {
	// Store pool on Blockchain //
	poolBytes, _ := json.Marshal( pool )
	err := stub.PutState( IndexLendingPools + pool.Token, poolBytes )
	if err != nil {
		err2 := errors.New( "ERROR: UPDATING THE STATE OF POOL FOR TOKEN " +
							pool.Token + ". " + err.Error() )
		return poolBytes, err2 }
	return poolBytes, nil
}

/* -------------------------------------------------------------------------------------------------
updateUserWallet: this function updates an User Wallet on Blockchain. 
                  Inputs: the updated user wallet.
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) updateUserWallet( stub shim.ChaincodeStubInterface, 
													        user_wallet MultiWallet) ( error ) {					
	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	userWalletBytes, _ := json.Marshal( user_wallet )
	multiChainCodeArgs := util.ToChaincodeArgs( "updateMultiwallets", 
												string(userWalletBytes) )
	response := stub.InvokeChaincode( COIN_BALANCE_CHAINCODE, multiChainCodeArgs, 
									  CHANNEL_NAME )
	if response.Status != shim.OK {
		err := errors.New( "ERROR INVOKING THE UPDATEMULTIWALLET CHAINCODE TO " +
						   "UPDATE THE WALLET OF USER: " + user_wallet.PublicId )
		return err }
	return nil
}
