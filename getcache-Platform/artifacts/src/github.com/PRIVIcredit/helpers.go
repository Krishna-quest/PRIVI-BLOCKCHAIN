
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
getPRIVIcreditList:  this function returns the list of PRIVI credits active in the Smart Contract. 
                     This function does not need to be called with any argument.
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) getPRIVIcreditList( stub shim.ChaincodeStubInterface ) (
													   map[string]bool, error ) {
	// Retrieve list of PRIVI Credits from the Blockchain //
	privi_list := make( map[string]bool )
	creditListBytes, _ := stub.GetState( IndexCreditList )
	err := json.Unmarshal( creditListBytes, &privi_list )
	if err != nil {
		return privi_list, errors.New( "ERROR: RETRIEVING THE PRIVI CREDIT LIST" ) }
	return privi_list, nil
}

/* -------------------------------------------------------------------------------------------------
retrieveUserWallet: this function returns the wallet of an User. Input is:
PublicId            string             // Public ID of the user to retrieve the wallet
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) retrieveUserWallet( stub shim.ChaincodeStubInterface, 
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
updateUserWallet: this function updates an User Wallet on Blockchain. 
                  Inputs: the updated user wallet.
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) updateUserWallet( stub shim.ChaincodeStubInterface, 
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

/* -------------------------------------------------------------------------------------------------
getRiskParameters:  this function gets the risk parameters of a given token.
		       	    Args: array containing 
token               string   // Token to retrieve the risk params
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) getRiskParameters( stub shim.ChaincodeStubInterface, 
												      token string) ( RiskParameters, error ) {
	// Retrieve PRIVI Credit information //
	risk_parameters := RiskParameters{}
	riskParamsBytes, err := stub.GetState( IndexRisk + token )
	json.Unmarshal( riskParamsBytes, &risk_parameters )
	if err != nil {
		return risk_parameters, errors.New( "ERROR: RETRIEVING PARAMS OF TOKEN. " + err.Error() ) }
	return risk_parameters, nil 
}
/* -------------------------------------------------------------------------------------------------
getPRIVIcredit:  this function returns the state of a particular PRIVI Credit.
		       	 Args: array containing 
args[0]               string   // ID of the PRIVI Credit
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) getPRIVIcredit( stub shim.ChaincodeStubInterface, 
												   loan_id string) ( PRIVIloan, error ) {
	// Retrieve PRIVI Credit information //
	privi_loan := PRIVIloan{}
	priviLoanBytes, err := stub.GetState( IndexPriviCredit + loan_id )
	json.Unmarshal( priviLoanBytes, &privi_loan )
	if err != nil {
		return privi_loan, errors.New( "ERROR: RETRIEVING THE PRIVI CREDIT" + err.Error() ) }
	return privi_loan, nil 
}

/* -------------------------------------------------------------------------------------------------
updateUserWallet: this function updates an User Wallet on Blockchain. 
                  Inputs: the updated user wallet.
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) updatePRIVIcredit( stub shim.ChaincodeStubInterface, 
													  privi_loan PRIVIloan) ( error ) {					
	// Update PRIVI loan on Blockchain //
	priviBytes, _ := json.Marshal( privi_loan )
	err := stub.PutState( IndexPriviCredit + privi_loan.LoanId, priviBytes )
	if err != nil {
		return errors.New( "ERROR: UPDATING PRIVI CREDIT " + privi_loan.LoanId + ". " +
							err.Error() ) }
	return nil
}