
/*--------------------------------------------------------------------------
----------------------------------------------------------------------------
   HELPER FUNCTIONS CALLED SEVERAL TIMES ON MAIN SMART CONTRACT FUNCIOTNS
----------------------------------------------------------------------------
-------------------------------------------------------------------------- */

package main
import (
	"errors"
	"math"
	//"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/common/util"
	//pb "github.com/hyperledger/fabric/protos/peer"
)

/* -------------------------------------------------------------------------------------------------
retrieveUserRole: this function returns the role of an User. Input is:
PublicId            string             // Public ID of the user to retrieve Role
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) retrieveUserRole( stub shim.ChaincodeStubInterface, 
								    	 PublicId string) (Actor, error) {

	// Retrieve wallet of an user from Blockchain //
	actor := Actor{}
	chainCodeArgs := util.ToChaincodeArgs( "getUser", PublicId )
	response := stub.InvokeChaincode( DATA_PROTOCOL_CHAINCODE, chainCodeArgs, 
									  CHANNEL_NAME )
	if response.Status != shim.OK {
		err := errors.New ( "ERROR INVOKING THE DATAPROTOCOL CHAINCODE TO " +
							"GET THE BALANCE OF USER: " + PublicId ) 
		return actor, err }
	json.Unmarshal(response.Payload, &actor)
	return actor, nil
}	

/* -------------------------------------------------------------------------------------------------
retrieveUserWallet: this function returns the wallet of an User. Input is:
PublicId            string             // Public ID of the user to retrieve the wallet
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) retrieveUserWallet( stub shim.ChaincodeStubInterface, 
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

func (t *PodInsurance) updateUserWallet( stub shim.ChaincodeStubInterface, 
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
retrieveInsuranceList: this function returns the list of insurance pools. No inputs needed.
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) retrieveInsuranceList( stub shim.ChaincodeStubInterface) (
	                                          map[string]bool, error ) {
	// Retrieve a pod from Blockchain //
	insurance_list := make( map[string]bool )
	insuranceListBytes, err := stub.GetState( IndexInsuranceList )
	json.Unmarshal( insuranceListBytes, &insurance_list )
	if err != nil {
		err2 := errors.New ("ERROR: GETTING THE INSURANCE LIST. " + err.Error() )
		return insurance_list, err2 }
	return insurance_list, nil
}

/* -------------------------------------------------------------------------------------------------
retrieveInsuranceInfo: this function returns the information of a given insurance pool. It takes 
                       as input the ID for the insurance pool.
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) retrieveInsuranceInfo( stub shim.ChaincodeStubInterface, 
	                           			      pool_id string) ( InsurancePool, error ) {
	// Retrieve wallet of an user from Blockchain //
	insurance_pool := InsurancePool{}
	insurancePoolBytes, err := stub.GetState( IndexInsurance + pool_id )
	json.Unmarshal( insurancePoolBytes, &insurance_pool )
	if err != nil {
		err2 := errors.New( "ERROR: GETTING THE POD LIST. " + err.Error() )
		return insurance_pool, err2 }
	return insurance_pool, nil
}	

/* -------------------------------------------------------------------------------------------------
updateInsuranceList: this function updates the list of insurance pools on Blockchain. 
                     Inputs: the updated insurance pool list.
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) updateInsuranceList( stub shim.ChaincodeStubInterface, 
							                insurance_list map[string]bool) ( error ) {
	// Store pod list on Blockchain //
	insuranceListBytes, _ := json.Marshal( insurance_list )
	err := stub.PutState( IndexInsuranceList, insuranceListBytes )
	if err != nil {
		err2 := errors.New( "ERROR: UPDATING THE POD LIST ON " +
							"BLOCKCHAIN. " + err.Error())
		return err2 }
	return nil
}

/* -------------------------------------------------------------------------------------------------
updateInsurance: this function updates a Pod on Blockchain. Inputs: the updated pod.
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) updateInsurance( stub shim.ChaincodeStubInterface, 
								        pool InsurancePool) ( error ) {
	// Store pod on Blockchain //
	insurancePoolBytes, _ := json.Marshal( pool )
	err := stub.PutState( IndexInsurance + pool.Id, insurancePoolBytes )
	if err != nil {
		err2 := errors.New( "ERROR: UPDATING THE STATE OF POD " +
							 pool.Id + ". " + err.Error() )
		return err2 }
	return nil
}

/* -------------------------------------------------------------------------------------------------
retrievePodInfo: this function returns the information of a given pod. It takes as input the
                 ID for the pod.
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) retrievePodInfo( stub shim.ChaincodeStubInterface, 
	                           	 		pod_id string) ( POD, error ) {
	// Retrieve NFT pod from Blockchain //
	pod := POD{}
	chainCodeArgs := util.ToChaincodeArgs( "retrievePodInfo", pod_id )
	response := stub.InvokeChaincode( POD_NFT_CHAINCODE, chainCodeArgs, 
									  CHANNEL_NAME )
	if response.Status != shim.OK {
		err := errors.New( "ERROR INVOKING THE PODNFT CHAINCODE TO " +
						   "GET THE POD INFO OF: " + pod_id )
		return pod, err }
	json.Unmarshal(response.Payload, &pod)
	return pod, nil
}

/* -------------------------------------------------------------------------------------------------
updatePod: this function updates  a Pod on Blockchain. Inputs: the updated pod.
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) updatePod( stub shim.ChaincodeStubInterface, 
								  pod POD) ( error ) {
	// Update pod on Blockchain // 
	podBytes, _ := json.Marshal( pod )
	multiChainCodeArgs := util.ToChaincodeArgs( "updatePod", string(podBytes) )
	response := stub.InvokeChaincode( POD_NFT_CHAINCODE, multiChainCodeArgs, 
									  CHANNEL_NAME )
	if response.Status != shim.OK {
		err := errors.New( "ERROR INVOKING THE PODNFT CHAINCODE TO " +
						   "UPDATE THE POD: " + pod.PodId )
		return err }
	return nil
}

/* -------------------------------------------------------------------------------------------------
retrieveObjects: this function returns the objects needed in the functions.
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) retrieveObjects( stub shim.ChaincodeStubInterface, 
							  	        insurance_id string,
								        wallet_id string, pod_id string ) ( POD, 
								        MultiWallet, InsurancePool, error ) {
	pod := POD{}
	wallet := MultiWallet{}
	insurance := InsurancePool{}
	err := errors.New("")

	// Retrieve insurance pool from blockchain // 
	if insurance_id != "" {
		insurance, err = t.retrieveInsuranceInfo( stub, insurance_id )
		if err != nil { return pod, wallet, insurance, err } }

	// Retrieve pod from blockchain // 
	if pod_id == "" { pod_id = insurance.PodId }
	pod, err = t.retrievePodInfo( stub, pod_id )
	if err != nil { return pod, wallet, insurance, err }
	
	// Retrieve user wallet from blockchain // 
	wallet, err = t.retrieveUserWallet( stub, wallet_id )
	if err != nil { return pod, wallet, insurance, err }
	
	return  pod, wallet, insurance, nil
}

/* -------------------------------------------------------------------------------------------------
createOutput: this function updates states and returns the object to update.
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) createOutput( stub shim.ChaincodeStubInterface, 
							         update_insurance map[string]InsurancePool,
							         update_wallets map[string]MultiWallet,
									 update_pods map[string]POD, 
									 transactions []Transfer ) ( []byte, error ) {
	outputBytes := []byte{}
	// Add new insurance pool to blockchain //
	for _, insurance := range(update_insurance) {
		err1 := t.updateInsurance(stub, insurance)
		if err1 != nil { return outputBytes, err1 } }
		
	// Update pod on Blockchain //
	for _, pod := range(update_pods) {
		err2 := t.updatePod(stub, pod)
		if err2 != nil { return outputBytes, err2 } }
	
	// Update guarantor wallet //
	for _, wallet := range(update_wallets) {
		err3 := t.updateUserWallet(stub, wallet)
		if err3 != nil { return outputBytes, err3 } }
	
	// Output of the result //
	output := Output{
		UpdatePods: update_pods, UpdateWallets: update_wallets,
		UpdateInsurance: update_insurance, Transaction: transactions }
	outputBytesGood, _ := json.Marshal(output)

	return outputBytesGood, nil
}

/* -------------------------------------------------------------------------------------------------
createPodBalance: this function updates a Pod on Blockchain. Inputs: the updated pod.
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) getParameters( stub shim.ChaincodeStubInterface ) ( 
	                                  Parameters, error ) {
	// Retrieve Parameters list from the Smart Contract //
	parameters := Parameters{}
	paramBytes, err := stub.GetState( IndexParameters )
	json.Unmarshal(paramBytes, &parameters)
	if err != nil {
		return parameters, errors.New("ERROR: GETTING PARAMETERS " + err.Error()) }
	return parameters, nil
}

/* -------------------------------------------------------------------------------------------------
checkWithdrawEnabled: this function checks if the participant can withdraw and returns the max
                      quantity that can withdrawn.
------------------------------------------------------------------------------------------------- */

func checkWithdrawEnabled( insurance InsurancePool, pod POD, id string ) ( float64 ) {

	amount := 0.
	// Case (1): withdrawer is the guarantor //
	if id == insurance.Guarantor {
		deposited := insurance.State.Insurers[ id ]
		coverage_guarantor := float64( insurance.Coverage ) *
							  insurance.Valuation
		amount = math.Max( 0., deposited.Amount - coverage_guarantor )
		return amount }
	
	// Case (2): withdrawer is an insurance investor //
	total_valuation := float64( pod.Supply ) * insurance.Valuation
	ratio_withdraw :=  pod.TotalInsurance / total_valuation - 1.
	amount = math.Max( 0., pod.TotalInsurance * ratio_withdraw )

	return amount
}

/* -------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------- */
