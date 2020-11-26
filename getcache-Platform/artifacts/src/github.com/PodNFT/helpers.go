
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
retrieveUserRole: this function returns the role of an User. Input is:
PublicId            string             // Public ID of the user to retrieve Role
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) retrieveUserRole( stub shim.ChaincodeStubInterface, 
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

func (t *PodNFT) retrieveUserWallet( stub shim.ChaincodeStubInterface, 
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

func (t *PodNFT) updateUserWallet( stub shim.ChaincodeStubInterface, 
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
retrievePodList: this function returns the list of Pods. No inputs for this function
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) retrievePodList( stub shim.ChaincodeStubInterface) (
	                                                map[string]bool, error ) {
	// Retrieve a pod from Blockchain //
	pod_list := make( map[string]bool )
	poolListBytes, err := stub.GetState( IndexPodList )
	json.Unmarshal( poolListBytes, &pod_list )
	if err != nil {
		err2 := errors.New ("ERROR: GETTING THE POD LIST. " + err.Error() )
		return pod_list, err2 }
	return pod_list, nil
}	

/* -------------------------------------------------------------------------------------------------
retrievePodInfo: this function returns the information of a given pod. It takes as input the
                 ID for the pod.
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) retrievePodInfo( stub shim.ChaincodeStubInterface, 
	                           						pod_id string) ( POD, error ) {
	// Retrieve wallet of an user from Blockchain //
	pod := POD{}
	podBytes, err := stub.GetState( IndexPods + pod_id )
	json.Unmarshal( podBytes, &pod )
	if err != nil {
		err2 := errors.New( "ERROR: GETTING THE POD LIST. " + err.Error() )
		return pod, err2 }
	return pod, nil
}	

/* -------------------------------------------------------------------------------------------------
updatePodList: this function updates the list of Pods on Blockchain. 
               Inputs: the updated pod list.
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) updatePodList( stub shim.ChaincodeStubInterface, 
								pod_list map[string]bool) ( error ) {
	// Store pod list on Blockchain //
	podListBytes, _ := json.Marshal( pod_list )
	err := stub.PutState( IndexPodList, podListBytes )
	if err != nil {
		err2 := errors.New( "ERROR: UPDATING THE POD LIST ON " +
							"BLOCKCHAIN. " + err.Error())
		return err2 }
	return nil
}

/* -------------------------------------------------------------------------------------------------
updatePod: this function updates a Pod on Blockchain. Inputs: the updated pod.
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) updatePod( stub shim.ChaincodeStubInterface, 
							pod POD) ( []byte, error ) {
	// Store pod on Blockchain //
	podBytes, _ := json.Marshal( pod )
	err := stub.PutState( IndexPods + pod.PodId, podBytes )
	if err != nil {
		err2 := errors.New( "ERROR: UPDATING THE STATE OF POD " +
							pod.PodId + ". " + err.Error() )
		return podBytes, err2 }
	return podBytes, nil
}

/* -------------------------------------------------------------------------------------------------
createOutput: this function updates states and returns the object to update.
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) createOutput( stub shim.ChaincodeStubInterface, 
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

func (t *PodNFT) createPodBalance( stub shim.ChaincodeStubInterface,  wallet MultiWallet, 
	                               id string, amount int64 ) ( MultiWallet ) {
	// Create new NFT balance //
	wallet.BalancesNFT[id] = BalanceNFT{
		Token: id, Amount: amount,
		PRIVIcreditLend: make(map[string]int64),
		PRIVIcreditBorrow: make(map[string]int64),
		PRIVI_lending: 0, PRIVI_borrowing: 0,
		Collateral_Amount: 0, Type: "NFT",
		Collaterals: make(map[string]float64) }
	return wallet
}

/* -------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------- */