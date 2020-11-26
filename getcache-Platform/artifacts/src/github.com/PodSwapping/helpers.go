
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
retrieveRiskParameters: this function returns the risk parameters for a given POD in a Risk Parameter
             	        object. Input is:
Token              string                // Token to retrieve POD rist parameters
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) retrieveRiskParameters( stub shim.ChaincodeStubInterface, 
								   token string) (RiskParameters, error) {

	// Retrieve risk parameters from Blockchain //
	risk := RiskParameters{}
	riskBytes, err := stub.GetState( IndexRisk + token )
	if err != nil {
		err2 := errors.New ("ERROR: GETTING THE RISK PARAM THE POD. " + err.Error() )
		return risk, err2 }
	json.Unmarshal( riskBytes, &risk )		
	return risk, nil
}	

/* -------------------------------------------------------------------------------------------------
retrievePoolList: this function returns the list of all the liquidity pools created. No input:
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) retrievePoolList( stub shim.ChaincodeStubInterface ) ( 
	                                                 map[string]bool, error ) {
	// Retrieve list of pools from Blockchain //
	pool_list := make( map[string]bool )
	poolBytes, err := stub.GetState( IndexPoolList )
	if err != nil {
		err2 := errors.New ( "ERROR: GETTING THE POOL LIST. " + err.Error() )
		return pool_list, err2 }
	json.Unmarshal( poolBytes, &pool_list )		
	return pool_list, nil
}	

/* -------------------------------------------------------------------------------------------------
retrievePoolInfo: this function returns the info for a given Liquidity Pools. It takes as input
                  the id of the pool.
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) retrievePoolInfo( stub shim.ChaincodeStubInterface, 
												pool_id string) ( LiquidityTokenPool, error ) {
	// Retrieve a liquidity pool from Blockchain //
	pool := LiquidityTokenPool{}
	poolBytes, err := stub.GetState( IndexPool + pool_id )
	json.Unmarshal( poolBytes, &pool )
	if err != nil {
		err2 := errors.New( "ERROR: GETTING THE POD LIST. " + err.Error() )
		return pool, err2 }
	return pool, nil
}	

/* -------------------------------------------------------------------------------------------------
retrieveUserWallet: this function returns the wallet of an User. Input is:
PublicId            string             // Public ID of the user to retrieve the wallet
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) retrieveUserWallet( stub shim.ChaincodeStubInterface, 
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
retrievePodList: this function returns the list of Pods. No inputs for this function
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) retrievePodList( stub shim.ChaincodeStubInterface) (
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

func (t *PodSwappingSmartContract) retrievePodInfo( stub shim.ChaincodeStubInterface, 
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

func (t *PodSwappingSmartContract) updatePodList( stub shim.ChaincodeStubInterface, 
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

func (t *PodSwappingSmartContract) updatePod( stub shim.ChaincodeStubInterface, 
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
updateLiquidityPool: this function updates a Liquidity Pool on Blockchain. Inputs: the updated pool.
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) updateLiquidityPool( stub shim.ChaincodeStubInterface, 
												pool LiquidityTokenPool) ( []byte, error ) {
	// Store pool on Blockchain //
	poolBytes, _ := json.Marshal( pool )
	err := stub.PutState( IndexPool + pool.Id, poolBytes )
	if err != nil {
		err2 := errors.New( "ERROR: UPDATING THE STATE OF POOL " +
							 pool.Id + ". " + err.Error() )
		return poolBytes, err2 }
	return poolBytes, nil
}

/* -------------------------------------------------------------------------------------------------
updateUserWallet: this function updates an User Wallet on Blockchain. 
                  Inputs: the updated user wallet.
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) updateUserWallet( stub shim.ChaincodeStubInterface, 
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
investmentSwap: this function returns the amount that the user gets by investing a certain 
                amount on the POD (with funding coin).
pod              POD          		            // Pod to compute swap
rateChange       map[string]float64             // Rate of change of tokens with BC
amount           float64                        // Amount desired to swap
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) investmentSwap( stub shim.ChaincodeStubInterface, 
			            pod POD, rateChange map[string]float64, amount float64) ( POD, float64 ) {
	
	// Constant of the conversion //
	C := pod.Invariant_cte
	X_cte := C / pod.InitialSupply

	// Compute liquidity funds //
	X_collateral_BC := 0.
	for token, amount_token := range(pod.Pools.Collateral_Pool) {
		X_collateral_BC = X_collateral_BC + amount_token*rateChange[token] }

	X_collateral_BC = math.Min( X_collateral_BC, 
								(1+pod.Interest)*pod.Principal*rateChange[pod.Token] )
	X_liquidity_t := ( pod.Pools.Interest_Pool + pod.Pools.Exchange_Pool + 
					   pod.Pools.Funding_Pool ) * rateChange[pod.Token] + X_collateral_BC

	invest_amount := X_liquidity_t * amount / ( pod.Pools.POD_Token_Pool - X_cte - amount ) 
	invest_amount = invest_amount / rateChange[pod.Token]
	return pod, invest_amount
}



/* -------------------------------------------------------------------------------------------------
interestSwap: this function returns the equivalent to be paid on interest from a payment of 
              interest in funding coin;
pod              POD          		            // Pod to compute swap
rateChange       map[string]float64             // Rate of change of tokens with BC
amount           float64                        // Amount desired to swap
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) interestSwap( stub shim.ChaincodeStubInterface, 
	pod POD, rateChange map[string]float64, amount float64) ( float64 ) {

	// Constant of the conversion //
	C := pod.Invariant_cte
	X_cte := C / pod.InitialSupply
	amount_BC := amount*rateChange[pod.Token]

	// Compute liquidity funds //
	X_collateral_BC := 0.
	for token, amount_token := range(pod.Pools.Collateral_Pool) {
		X_collateral_BC = X_collateral_BC + amount_token*rateChange[token] }

	X_collateral_BC = math.Min( X_collateral_BC, 
					 (1+pod.Interest)*pod.Principal*rateChange[pod.Token] )
	X_liquidity_t := ( pod.Pools.Interest_Pool + pod.Pools.Exchange_Pool + 
					   pod.Pools.Funding_Pool ) * rateChange[pod.Token] + X_collateral_BC

	interest_amount := ( pod.Pools.POD_Token_Pool - X_cte ) * amount_BC / (  X_liquidity_t + amount_BC ) 
	return interest_amount 
}

/* -------------------------------------------------------------------------------------------------
swapFunctionCrypto: this function returns the amount that the user gets by swapping a certain 
                    amount on the pod token.
pod              POD          		            // Pod to compute swap
rateChange       map[string]float64             // Rate of change of tokens with BC
amount           float64                        // Amount desired to swap
token            string                         // Token desired to get
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) swapFunctionCrypto( stub shim.ChaincodeStubInterface, 
	  		pod POD, rateChange map[string]float64, amount float64, token string) ( POD, float64 ) {

	// Constant of the conversion //
	C := pod.Invariant_cte
	X_cte := C / pod.InitialSupply

	// Compute liquidity funds //
	X_collateral_BC := 0.
	for token, amount_pool := range(pod.Pools.Collateral_Pool) {
		X_collateral_BC = X_collateral_BC + amount_pool*rateChange[token] }
	
	X_collateral_BC = math.Min( X_collateral_BC, 
					  (1+pod.Interest)*pod.Principal*rateChange[pod.Token] )
	X_liquidity_t := ( pod.Pools.Interest_Pool + pod.Pools.Exchange_Pool + 
					   pod.Pools.Funding_Pool ) * rateChange[pod.Token] + X_collateral_BC
	swap_amount := X_liquidity_t * amount / ( pod.Pools.POD_Token_Pool - X_cte + amount ) 
	swap_amount = swap_amount / rateChange[token]

	return pod, swap_amount
}

/* -------------------------------------------------------------------------------------------------
swapLiquidityPool: this function is used to compute an swap a pod token by the pool token.
pool             LiquidityTokenPool          	// Liquidity pool where the swap takes place
swap_amount      float64                        // Amount of the Liquidity Pool to swap
fee_amount       float64                        // Amount charged as fee
pod_amount       float64                        // Amount of pod token to introduce in the pool
pod_token        string                         // Id of the pod token to introduce on pool
RateChange       map[string]float64             // Rate of change of tokens
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) swapLiquidityPool( stub shim.ChaincodeStubInterface, 
					 			   pool LiquidityTokenPool, swap_amount float64, fee_amount float64,
								   pod_amount float64, pod_token string, 
								   RateChange map[string]float64 ) ( LiquidityTokenPool, error ) {

	// Check that the pool contains enough funds //
	if pool.State.Reserve < swap_amount-fee_amount {
		return pool, errors.New( "ERROR: THE LIQUIDITY POOL " + pool.Id + "DOES NOT HOLD " + 
								 "ENOUGH FUNDS OF TOKEN" + pool.Token ) }

	// Compute swapping of tokens //
	pool.State.Reserve = pool.State.Reserve - ( swap_amount - fee_amount )
	pod_token_reserve, isInList := pool.State.Liq_Tokens[ pod_token ]
	if !isInList { pod_token_reserve = 0. }
	pod_token_reserve = pod_token_reserve + pod_amount
	pool.State.Liq_Tokens[ pod_token ] = pod_token_reserve

	// Compute new reserve ratio //
	pod_value_BC := 0.
	for token, amount := range( pool.State.Liq_Tokens ) {
		pod_value_BC = pod_value_BC + amount * RateChange[token] }
	reserve_BC := pool.State.Reserve * RateChange[pool.Token]
	deposited_BC := pool.State.Deposited * RateChange[pool.Token]
	pool.State.ReserveRatio = ( pod_value_BC + reserve_BC ) / deposited_BC

	// Check that Reserve Ratio keeps above the Min Reserve Ratio //
	if pool.State.ReserveRatio < pool.MinReserveRatio {
		return pool, errors.New( "ERROR: POD TOKENS NOT ACCEPTED. THE RESERVE RATIO " +
								 "DROPS BELOW THE MIN RESERVE RATIO" ) }
	return pool, nil
}

