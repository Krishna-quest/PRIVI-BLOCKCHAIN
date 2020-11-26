package main

import (
	//"time"
	"encoding/json"
	"fmt"
	"errors"
	//"strconv"
	"math"
	"github.com/rs/xid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric/common/util"
)


/* -------------------------------------------------------------------------------------------------
Init:  this function register PRIVI as the Admin of the POD Swapping Smart Contract. It initialises
       the lists of PODs and the Indexes used. Args: array containing a string:
PrivateKeyID           string   // Private Key of the admin of the smart contract
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) Init( stub shim.ChaincodeStubInterface ) pb.Response {

	_, args := stub.GetFunctionAndParameters()
	// Store in the state of the smart contract the Private Key of Admin //
	err1 := stub.PutState( IndexAdmin, []byte(args[0])  )
	if err1 != nil {
		return shim.Error( "ERROR: SETTING THE ADMIN PRIVATE KEY: " +
						   err1.Error() ) }
	if args[1] == "UPGRADE" { return shim.Success(nil) }

	// Initialise list of PODs in the Smart Contract as empty //
	pod_list, _ := json.Marshal( make(map[string]bool) )
	err2 := stub.PutState( IndexPodList, pod_list )
	if err2 != nil {
		return shim.Error( "ERROR: INITIALISING THE POD LIST: " +
	                       err2.Error() ) }
	// Initialise risk parameters index as empty //
	risk, _ := json.Marshal( RiskParameters{} )
	err3 := stub.PutState( IndexRisk, risk )
	if err3 != nil {
		return shim.Error( "ERROR: INITIALISING THE RISK PARAMETERS: " +
						   err3.Error() ) }
	// Initialise pod token parameters index as empty //
	pool_list, _ := json.Marshal( make(map[string]bool) )
	err4 := stub.PutState( IndexPoolList, pool_list )
	if err4 != nil {
		return shim.Error( "ERROR: INITIALISING THE POOL LIST: " +
						   err4.Error() ) }	   
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
 The Invoke method is called as a result of an application request to run the Smart Contract ""
 The calling application program has also specified the particular smart contract function to be called
-------------------------------------------------------------------------------------------------*/

func (t *PodSwappingSmartContract) Invoke( stub shim.ChaincodeStubInterface ) pb.Response {

	// Retrieve function and arguments //
	function, args := stub.GetFunctionAndParameters()

	// Retrieve caller of the function //
	caller, err1 := CallerCN(stub)
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE CALLER OF THE TRANSFER " +
						   "FUNCTION. " + err1.Error() ) }
	// Retrieve admin of the smart contract //
	adminBytes, err2 := stub.GetState( IndexAdmin )
	if err2 != nil {
		return shim.Error( "ERROR: GETTING THE ADMIN OF THE SMART " +
						   "CONTRACT. " + err2.Error() ) }
	admin := string( adminBytes )
	// Verify that the caller of the function is admin //
	if caller != admin {
		return shim.Error( "ERROR: CALLER " + caller + " DOES NOT HAVE PERMISSION " ) }

	// Call the proper function //
	switch function {
		case "updateRiskParameters":
			return t.updateRiskParameters(stub, args)
		case "retrievePodList":
			pod_list, err := t.retrievePodList(stub)
			if err != nil { return shim.Error( err.Error() ) }
			podlistBytes, _ := json.Marshal( pod_list )
			return shim.Success( podlistBytes )
		case "retrievePoolList":
			pool_list, err := t.retrievePoolList(stub)
			if err != nil { return shim.Error( err.Error() ) }
			poolListBytes, _ := json.Marshal( pool_list )
			return shim.Success( poolListBytes )
		case "retrievePodInfo":
			pod, err := t.retrievePodInfo(stub, args[0])
			if err != nil { return shim.Error( err.Error() ) }
			podBytes, _ := json.Marshal( pod )
			return shim.Success( podBytes )
		case "retrievePoolInfo":
			pool, err := t.retrievePoolInfo(stub, args[0])
			if err != nil { return shim.Error( err.Error() ) }
			poolBytes, _ := json.Marshal( pool )
			return shim.Success( poolBytes )
		case "createLiquidityPool":
			return t.createLiquidityPool(stub, args)
		case "depositLiquidity":
			return t.depositLiquidity(stub, args)
		case "withdrawLiquidity":
			return t.withdrawLiquidity(stub, args)
		case "managerLiquidity":
			return t.managerLiquidity(stub, args)
		case "initiatePOD":
			return t.initiatePOD(stub, args)
		case "deletePOD":
			return t.deletePOD(stub, args)
		case "investPOD":
			return t.investPOD(stub, args)
		case "interestPOD":
			result, err := t.interestPOD(stub, args)
			if err != nil { return shim.Error( err.Error() ) }
			return shim.Success( result )
		case "swapPOD":
			return t.swapPOD(stub, args)
		case "managerPOD":
			return t.managerPOD(stub, args)
		case "liquidatePOD":
			result, err := t.liquidatePOD(stub, args)
			if err != nil { return shim.Error( err.Error() ) }
			return shim.Success( result )
		case "checkPODLiquidation":
			return t.checkPODLiquidation(stub, args)	
	}
	return shim.Error("Incorrect function name: " + function)
}


/* -------------------------------------------------------------------------------------------------
updateRiskParameters:  this function is called to update the risk parameters of the POD Smart
				 	   Contract for a given token. Args is an array containing a json with:
Token                 string    // Token Symbol (args[0])
Interest_min		  float64   // Minimum interest rate allowed for a POD
Interest_max		  float64   // Maximum interest rate allowed for a POD
Pct_supply_lower      float64   // Lower pct allowed for the supply
Pct_supply_upper      float64   // Upper pct allowed for the supply
Liquidation_min       float64   // Threshold limit to liquidate POD
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) updateRiskParameters( stub shim.ChaincodeStubInterface, 
													     args []string) pb.Response {

	// Retrieve the input information of Risk Parameters //
	token := args[0]
	risk := RiskParameters{}
	err1 := json.Unmarshal( []byte(args[1]), &risk )
	if err1 != nil {
		return shim.Error( "ERROR: RETRIEVING THE INPUT " + err1.Error() ) 
	}
	// Initialise/Update risk parameters for the Token //
	riskParamsBytes2, _ := json.Marshal( risk )
	err2 := stub.PutState( IndexRisk + token, riskParamsBytes2 )
	if err2 != nil {
		return shim.Error( "ERROR: UPDATING RISK PARAMS " + err2.Error() ) 
	}
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
initiateLiquidityPool: this function initialises a new liquidity pool for a given token.
                       Args is an array containing a json with the following fields:
Token                string  			   // Token of the loan principal
CreatorId            string                // Id of the pool creator
MinReserveRatio      float64               // Minimum reserve ratio targeted
InitialAmount        float64  			   // Initial amount deposited on the reserve
WithdrawalTime       int64                 // Number of days to be able to withdraw without fee
WithdrawalFee        float64               // Instantaneous withdrawal fee (at t=0)
MinTrustScore        float64               // Minimum Trust Score required
MinEndScore          float64               // Minimum Endorsement Score required
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) createLiquidityPool( stub shim.ChaincodeStubInterface, 
														args []string) pb.Response {

	// Retrieve the input information for Pool Initialisation //
	liq_pool := LiquidityTokenPool{}
	err1 := json.Unmarshal( []byte(args[0]), &liq_pool )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	liq_pool.Id = xid.New().String() + xid.New().String() 
	date := getTimeNow()

	// Retrieve wallet of pool creator //
	wallet_creator, err2 := t.retrieveUserWallet(stub, liq_pool.CreatorId)
	if err2 != nil { shim.Error (err2.Error() ) }
	balance_pool := wallet_creator.Balances[liq_pool.Token] 

	// Check that creator has enough funds and compute transfers //
	if balance_pool.Amount < liq_pool.InitialAmount {
		return shim.Error( "ERROR: POOL CREATOR DOES NOT HOLD ENOUGH FUNDS" ) }
	balance_pool.Amount = balance_pool.Amount - liq_pool.InitialAmount
	pool_transfer := Transfer{ 
		Type: "liquidity_pool_creation", Token: liq_pool.Token, 
		Amount: liq_pool.InitialAmount, From: liq_pool.CreatorId, 
		To: "Liquidity Pool " + liq_pool.Id, Id: xid.New().String(),
	    Date: date }
	wallet_creator.Balances[liq_pool.Token] = balance_pool
	wallet_creator.Transaction = []Transfer{pool_transfer}

	// Create the new liquidity pool // 
	liq_providers := make( map[string]Provider )
	liq_providers[liq_pool.CreatorId] = Provider{ 
		Amount: liq_pool.InitialAmount, WithdrawDay: 0 }
	pool_state := PoolState{
		LiqProviders: liq_providers, LiqProvidersNum: 1, 
		Reserve: liq_pool.InitialAmount, Deposited: liq_pool.InitialAmount,
		Liq_Tokens: make( map[string]float64 ), ReserveRatio: 1. }
	liq_pool.State = pool_state
	liq_pool.Date = getTimeNow()

	// Add new liquidity pool on Blockchain //
	_, err3 := t.updateLiquidityPool( stub, liq_pool )
	if err3 != nil { return shim.Error( err3.Error() ) }

	// Retrieve pool list and add new token pool //
	pool_list, err4 := t.retrievePoolList( stub )
	if err4 != nil { return shim.Error( err4.Error() ) }
	pool_list[liq_pool.Id] = true
	liqPoolListBytes, _ := json.Marshal( pool_list )
	err5 := stub.PutState( IndexPoolList, liqPoolListBytes )
	if err5 != nil {
		return shim.Error( "ERROR: UPDATING LIQUIDITY POOL LIST. " + 
							err5.Error() ) }
	
	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	err6 := t.updateUserWallet( stub, wallet_creator )
	if err6 != nil { return shim.Error( err6.Error() ) }

	// Generate output //
	update_pools := make( map[string]LiquidityTokenPool )
	update_wallets := make( map[string]MultiWallet )
	update_pools[liq_pool.Id] = liq_pool
	update_wallets[wallet_creator.PublicId] =  wallet_creator
	output := Output{
		UpdatePools: update_pools, UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes ) 
}

/* -------------------------------------------------------------------------------------------------
depositLiquidity: this function is used when a Liquidity Provider wants to deposit funds on a
                      Liquidity pool. Args is an array containing a json with the following fields:
LiqPoolId            string  			   // Id of the pool to deposit the funds
LiqProviderId        string                // Id of the liquidity provider
Amount               float64               // Amount to deposit
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) depositLiquidity( stub shim.ChaincodeStubInterface, 
													 args []string) pb.Response {

	// Retrieve the input information //
	input := LiquidityDeposit{}
	err1 := json.Unmarshal( []byte(args[0]), &input )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	date := getTimeNow()

	// Retrieve liquity pool  //
	liq_pool, err2 := t.retrievePoolInfo(stub, input.LiqPoolId)
	if err2 != nil { shim.Error(err2.Error()) }

	// Retrieve wallet of liquidity provider //
	wallet_provider, err3 := t.retrieveUserWallet(stub, input.LiqProviderId)
	if err3 != nil { shim.Error (err3.Error() ) }
	balance_pool := wallet_provider.Balances[liq_pool.Token]

	// Check that LP has enough funds and compute transfers //
	if balance_pool.Amount < input.Amount {
		return shim.Error( "ERROR: LIQUIDITY PROVIDER DOES NOT HOLD ENOUGH FUNDS" ) }
	balance_pool.Amount = balance_pool.Amount - input.Amount
	pool_transfer := Transfer{ 
		Type: "liquidity_pool_provide", Token: liq_pool.Token,  Amount: input.Amount,
		From: input.LiqProviderId, To: "Liquidity Pool " + input.LiqPoolId, 
		Id: xid.New().String(), Date: date }
	wallet_provider.Balances[liq_pool.Token] = balance_pool
	wallet_provider.Transaction = []Transfer{pool_transfer}

	// Add liquidity provider to the pool // 
	liq_providers := liq_pool.State.LiqProviders
	provider, isInList := liq_providers[input.LiqProviderId]
	if !isInList{ 
		liq_pool.State.LiqProvidersNum = liq_pool.State.LiqProvidersNum+1
		provider = Provider{Amount: 0., WithdrawDay: 0} }
	provider.Amount = provider.Amount + input.Amount
	liq_providers[input.LiqProviderId] = provider 

	// Update liquidity Pool //
	liq_pool.State.LiqProviders = liq_providers
	liq_pool.State.Deposited = liq_pool.State.Deposited+input.Amount
	liq_pool.State.Reserve = liq_pool.State.Reserve+input.Amount

	// Update liquidity pool on Blockchain //
	_, err4 := t.updateLiquidityPool( stub, liq_pool )
	if err4 != nil { return shim.Error( err4.Error() ) }

	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	err5 := t.updateUserWallet( stub, wallet_provider )
	if err5 != nil { return shim.Error( err5.Error() ) }

	// Generate output //
	update_pools := make( map[string]LiquidityTokenPool )
	update_wallets := make( map[string]MultiWallet )
	update_pools[liq_pool.Id] = liq_pool
	update_wallets[wallet_provider.PublicId] = wallet_provider
	output := Output{
		UpdatePools: update_pools, UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes ) 
}


/* -------------------------------------------------------------------------------------------------
withdrawLiquidity: this function is used when a Liquidity Provider wants to withdraw funds from
                   a Liquidity pool. Args is an array containing a json with the following fields:
LiqPoolId            string  			   // Id of the pool to deposit the funds
LiqProviderId        string                // Id of the liquidity provider
Amount               float64               // Amount to withdraw
RateChange           map[string]float64    // Rate of Change of Tokens with BC
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) withdrawLiquidity( stub shim.ChaincodeStubInterface, 
													  args []string ) pb.Response {
	
	// Retrieve the input information for POD Initialisation //
	input := LiquidityWithdrawal{}
	err1 := json.Unmarshal( []byte(args[0]), &input )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	date := getTimeNow()
	update_pods := make( map[string]POD )
	update_pools := make( map[string]LiquidityTokenPool )
	update_wallets := make( map[string]MultiWallet )

	// Retrieve liquity pool  //
	liq_pool, err2 := t.retrievePoolInfo(stub, input.LiqPoolId)
	if err2 != nil { shim.Error(err2.Error()) }

	// Retrieve wallet of liquidity provider //
	wallet_provider, err3 := t.retrieveUserWallet(stub, input.LiqProviderId)
	if err3 != nil { shim.Error (err3.Error() ) }
	balance_pool := wallet_provider.Balances[liq_pool.Token]

	// Check that LP deposited the amount that he wants to withdraw //
	liq_provider, isInList := liq_pool.State.LiqProviders[ input.LiqProviderId ]
	if !isInList {
		return shim.Error( "ERROR: THE ID " + input.LiqProviderId + " IS NOT A" +
						   " LIQUIDITY PROVIDER FOR THE POOL " + input.LiqPoolId ) }
	if liq_provider.Amount < input.Amount {
		return shim.Error( "ERROR: THE LP " + input.LiqProviderId  + " DID NOT" +
						   " DEPOSITED THAT FUNDS ON POOL " + input.LiqPoolId ) }
	
	// Retrieve model parameters //
	transactions := []Transfer{}
	r := - math.Log( 1 + liq_pool.WithdrawalFee )
	T := liq_pool.WithdrawalTime
	time := liq_provider.WithdrawDay
	reserve := liq_pool.State.Reserve

	// Take first amount from deposit //
	reserve_fee_rate := math.Exp( -r * math.Max(float64(T-time), 0.) / float64(T) ) - 1
	reserve_amount := math.Min( reserve, input.Amount )
	reserve_fee := reserve_amount * reserve_fee_rate

	liq_pool.State.Reserve = liq_pool.State.Reserve - reserve_amount + reserve_fee
	balance_pool.Amount = balance_pool.Amount + (reserve_amount-reserve_fee)
	withdrawal_transfer := Transfer{
		Type: "liquidity_pool_withdraw", Amount: reserve_amount - reserve_fee, 
		Token: liq_pool.Token, From: "Liquidity Pool " + input.LiqPoolId, 
		To: input.LiqProviderId, Id: xid.New().String(), Date: date }
	fee_transfer := Transfer{
		Type: "liquidity_pool_fee_withdraw", Amount: reserve_fee, Token: 
		liq_pool.Token,From: input.LiqProviderId,
		To: "Liquidity Pool " + input.LiqPoolId, 
		Id: xid.New().String(), Date: date }
	transactions = append( transactions, withdrawal_transfer )
	transactions = append( transactions, fee_transfer )

	// If deposit is not enough, take from token pools, proportionally //
	amount_left_BC := ( input.Amount - reserve_amount ) * input.RateChange[liq_pool.Token]
	if amount_left_BC > 0 {
		pod_liquidity_BC := 0.
		liq_tokens := liq_pool.State.Liq_Tokens
		for token, amount := range(liq_tokens){
			pod_liquidity_BC = pod_liquidity_BC + amount * input.RateChange[token] }

		if pod_liquidity_BC < amount_left_BC {
			return shim.Error( "ERROR: NOT ENOUGH LIQUIDITY TO WITHDRAW FUNDS" ) }
		
		// Distribute mix of all pod tokens proportionally //
		for token, amount := range(liq_tokens) {
			amount_pod := amount * amount_left_BC / pod_liquidity_BC
			amount_pod_fee := amount_pod * reserve_fee_rate
			liq_tokens[token] = math.Max( 0, liq_tokens[token] - amount_pod +
											 amount_pod_fee )
			
			// Retrieve POD to update //
			pod, err4 := t.retrievePodInfo(stub, token)
			if err4 != nil { shim.Error(err4.Error() ) }
			pod_list, err5 := t.retrievePodList(stub)
			if err5 != nil { shim.Error(err5.Error() ) }
			active, isInList := pod_list[token]
			if !isInList || active == false {
				return shim.Error( "ERROR: POD ID IS NOT REGISTERED OR LIQUIDATED" ) }

			// Check that user has the balance of the pod token registered on wallet //
			investment_amount, isInList := pod.State.Investors[input.LiqProviderId]
			balance_swap, _ := wallet_provider.Balances[ token ]
			if !isInList {
				balance_swap = Balance{ 
					Token: token, Amount: 0., Credit_Amount: 0., 
					Staking_Amount: 0., Borrowing_Amount: 0., 
					PRIVI_lending: 0., PRIVI_borrowing: 0., 
					Collateral_Amount: 0., Type: "POD",
					Collaterals: make( map[string]float64) }
				investment_amount = 0.
				pod.State.InvestorNum = pod.State.InvestorNum + 1 
			}

			pod.State.Investors[input.LiqProviderId] = investment_amount +
													   amount_pod - amount_pod_fee
			balance_swap.Amount = balance_swap.Amount + (amount_pod-amount_pod_fee)
			wallet_provider.Balances[ token ] = balance_swap

			// Compute transfers //
			pod_transfer := Transfer{
				Type: "liquidity_pool_withdraw", Amount: amount_pod - amount_pod_fee, 
				Token: token, From: "Liquidity Pool " + input.LiqPoolId, 
				To: input.LiqProviderId, Id: xid.New().String(), Date: date }  
			fee_pood_transfer := Transfer{
				Type: "liquidity_pool_fee_withdraw", Amount: amount_pod_fee, Token: token,
				From: input.LiqProviderId, To: "Liquidity Pool " + input.LiqPoolId,
			    Id: xid.New().String(), Date: date  }
			transactions = append( transactions, pod_transfer )
			transactions = append( transactions, fee_pood_transfer ) 

			// Update Pod State on Blockchain //   
			update_pods[pod.PodId] = pod                     
			_, err6 := t.updatePod( stub, pod )
			if err6 != nil { return shim.Error( err6.Error() ) }
		}
		liq_pool.State.Liq_Tokens = liq_tokens
	}

	// Update liquidity Pool //
	liq_provider.Amount = liq_provider.Amount - input.Amount
	liq_pool.State.LiqProviders[ input.LiqProviderId ] = liq_provider
	if liq_provider.Amount == 0. { 
		delete(liq_pool.State.LiqProviders, input.LiqProviderId) }
	liq_pool.State.Deposited = liq_pool.State.Deposited-input.Amount

	// Update liquidity pool on Blockchain //
	update_pools[liq_pool.Id] = liq_pool
	_, err7 := t.updateLiquidityPool( stub, liq_pool )
	if err7 != nil { return shim.Error( err7.Error() ) }

	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	update_wallets[wallet_provider.PublicId] = wallet_provider
	err8 := t.updateUserWallet( stub, wallet_provider )
	if err8 != nil { return shim.Error( err8.Error() ) }

	// Generate output //
	output := Output{
		UpdatePools: update_pools, UpdateWallets: update_wallets,
	    UpdatePods: update_pods }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes ) 
}

/* -------------------------------------------------------------------------------------------------
managerLiquidity: this function is called daily on a cron job to update all the Withdrawal days of
                  the LP. It is called without any argument.
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) managerLiquidity( stub shim.ChaincodeStubInterface, 
					                                 args []string ) pb.Response {
	// Retrieve pool list //
	update_pools := make( map[string]LiquidityTokenPool )
	pool_list, err1 := t.retrievePoolList( stub )
	if err1 != nil { return shim.Error( err1.Error() ) }

	// Loop through all active Liquidity Pools and update withdrawal day //
	for pool_id, isActive := range(pool_list) {
		// Retrieve liquity pool  //
		pool, err2 := t.retrievePoolInfo(stub, pool_id)
		if err2 != nil { shim.Error(err2.Error()) }
		if !isActive {continue}
		liq_providers := pool.State.LiqProviders
		for provider_id, provider := range(liq_providers) {
			provider.WithdrawDay = provider.WithdrawDay + 1
			liq_providers[provider_id] = provider }
		pool.State.LiqProviders = liq_providers
		// Update liquidity pool on Blockchain //
	    _, err3 := t.updateLiquidityPool( stub, pool )
		if err3 != nil { return shim.Error( err3.Error() ) }
		update_pools[pool_id] = pool
	}

	// Generate output //
	output := Output{ UpdatePools: update_pools }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes ) 
}


/* -------------------------------------------------------------------------------------------------
initiatePOD: this function initialises a new POD with the parameters described below. 
             Args is an array containing three json with the following fields:
Creator              string                // Id of the creator of the POD
Token                string  			   // Token of the loan principal
Duration             int64                 // Duration of the POD (in days)
Payments             int64  			   // Number of payments of the pod
Principal         	 float64  			   // Principal amount targeted to raise 
Interest             float64  			   // Interest Rate of the POD
P_liquidation        float64               // CCR threshold for pod liquidation
InitialSupply        float64               // Initial supply of Pod Tokens created
Collaterals          map[string]float64    // Collateral deposited ( args[1] )
RateChange           map[string]float64    // Rate of changes of Collaterals ( args[2] )
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) initiatePOD( stub shim.ChaincodeStubInterface, 
	                                  			args []string) pb.Response {

	// Retrieve the input information for POD Initialisation //
	pod := POD{}
	err1 := json.Unmarshal( []byte(args[0]), &pod )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	collaterals := make( map[string]float64  )
	json.Unmarshal( []byte(args[1]), &collaterals )
	rateChange := make( map[string]float64  )
	json.Unmarshal( []byte(args[2]), &rateChange )
	date := getTimeNow()
	update_pods := make( map[string]POD )
	update_wallets := make( map[string]MultiWallet )

	// Retrieve the risk parameters of POD for the token //
	risk_params, err2 := t.retrieveRiskParameters(stub, pod.Token)
	if err2 != nil { shim.Error(err2.Error() ) }
	
	// Retrieve wallet of POD Creator //
	wallet_creator, err3 := t.retrieveUserWallet(stub, pod.Creator)
	if err3 != nil { shim.Error(err3.Error() ) }

	// Retrieve pod list //
	pod_list, err4 := t.retrievePodList(stub)
	if err4 != nil { shim.Error(err4.Error() ) }

	// Check that the Pod conditions satisfy the risk parameters requirements //
	if pod.InitialSupply < pod.Principal * risk_params.Pct_supply_lower ||
	   pod.InitialSupply > pod.Principal * risk_params.Pct_supply_upper {
		return shim.Error( "ERROR: INITIAL SUPPLY FOR THE TOKEN POD SHOULD BE " +
						   "BETWEEN THE BOUNDS REQUIRED." ) }
	if pod.Interest < risk_params.Interest_min || 
	   pod.Interest > risk_params.Interest_max {
		return shim.Error( "ERROR: THE POD INTEREST SHOULD BE BETWEEN " +
	                       "THE BOUNDS REQUIRED.") } 
	if pod.P_liquidation < risk_params.Liquidation_min {
		return shim.Error( "ERROR: THE THRESHOLD FOR LIQUIDATION SHOULD BE" +
						   "GREATER THAN THE MINIMIM RISK PARAMETER.") }
	pod.PodId = xid.New().String() + xid.New().String() 
	
	// Check that the creator holds the collateral and transfer it //
	transactions := []Transfer{}
	collateral_pool := make( map[string]float64  )
	collateral_BC := 0.
	for token, amount := range(collaterals) {
		balance_token := wallet_creator.Balances[token]
		if balance_token.Amount < amount {
			return shim.Error( "ERROR: USER " + pod.Creator + "DOES NOT " +
							   "HOLD ENOUGH COLLATERAL FOR TOKEN " + token) }
		// Update user balance //
		balance_token.Amount = balance_token.Amount - amount
		col_transfer := Transfer{
			Type: "POD_creation", Token: token, Amount: amount,
			From: pod.Creator, To: "POD Pool " + pod.PodId,
		    Id: xid.New().String(), Date: date }
		transactions = append( transactions, col_transfer )
		// Update pollateral pool //
		collateral_pool[token] = amount
		collateral_BC = collateral_BC + amount*rateChange[token]
		wallet_creator.Balances[token] = balance_token
	}

	// Check that collateral condition is satisfied and set invariant cte //
	principal_BC := pod.Principal * rateChange[pod.Token]
	CCR := collateral_BC / principal_BC 
	if CCR < pod.P_liquidation {
		return shim.Error( "ERROR: THE POD DOES NOT SATISFY THE " +
						   "MINIMUM COLLATERAL CONDITION.") }
	collateral_BC_liquidity := math.Min( collateral_BC, 
							(1+pod.Interest)*pod.Principal*rateChange[pod.Token] )
	pod.Invariant_cte = 0.5 * collateral_BC_liquidity * pod.InitialSupply

	// Initialise state and pools of the POD //
	state := PODstate{
		Liq_Pools: make( map[string]float64 ), 
		Investors: make( map[string]float64 ), 
		InvestorNum: 0, Status: "INITIATED",
		POD_Day: 0, Debt: 0., MissingPayments: 0,
		FundsRaised: 0., SupplyReleased: 0. }
	pod_pools := PODpools{
		Funding_Pool: 0., Exchange_Pool: 0., 
		Collateral_Pool: collateral_pool, 
		Interest_Pool: 0., 
		POD_Token_Pool: pod.InitialSupply }

	// Create new pod and store in Blockchain and add pod to PodList //
	pod.Date = getTimeNow()
	pod.State = state
	pod.Pools = pod_pools
	pod_list[pod.PodId] = true
	_, err5 := t.updatePod( stub, pod )
	if err5 != nil { return shim.Error( err5.Error() ) }
	err6 := t.updatePodList( stub, pod_list )
	if err6 != nil { return shim.Error( err6.Error() ) }
	update_pods[pod.PodId] = pod
	
	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	wallet_creator.Transaction = transactions
	err7 := t.updateUserWallet( stub, wallet_creator )
	if err7 != nil { return shim.Error( err7.Error() ) }
	update_wallets[wallet_creator.PublicId] = wallet_creator

	// Generate output //
	output := Output{ 
		 UpdatePods: update_pods, UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes ) 
}


/* -------------------------------------------------------------------------------------------------
deletePOD: this function is used to delete a POD if the state is INITIATED (ie, not OPEN). It can 
           only be called by user. Args is an array containing a json with the following fields:
PublicId              string                // Id of the creator of the POD
PodId                 string  			   // Id of the POD to delete
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) deletePOD( stub shim.ChaincodeStubInterface, 
											  args []string) pb.Response {

	// Retrieve the input information for POD Initialisation //
	deletion := PodDeletion{}
	err1 := json.Unmarshal( []byte(args[0]), &deletion )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	date := getTimeNow()
	update_pods := make( map[string]POD )
	update_wallets := make( map[string]MultiWallet )
	
	// Retrieve pod list //
	pod_list, err2 := t.retrievePodList(stub)
	if err2 != nil { shim.Error(err2.Error() ) }
	_, isInList := pod_list[deletion.PodId]
	
	// Retrieve pod //
	pod, err3 := t.retrievePodInfo(stub, deletion.PodId)
	if err3 != nil { shim.Error(err3.Error() ) }

	// Check that Pod is in state "INITIATED" and the caller is the Creator //
	if !isInList {
		return shim.Error( "ERROR: POD ID IS NOT REGISTERED" ) }
	if pod.State.Status != "INITIATED" {
		return shim.Error( "ERROR: POD CANNOT BE DELETED IF IS NOT IN INITIATED STATE" ) }
	if pod.Creator != deletion.PublicId {
		return shim.Error( "ERROR: THE CALLER IS NOT THE CREATOR OF THE POD." ) }

	// Retrieve wallet of POD Creator //
	wallet_creator, err4 := t.retrieveUserWallet(stub, deletion.PublicId)
	if err4 != nil { shim.Error(err4.Error() ) }

	// Return collateral to Pod Creator //
	transactions := []Transfer{}
	for token, amount := range(pod.Pools.Collateral_Pool) {
		balance_token := wallet_creator.Balances[token]
		balance_token.Amount = balance_token.Amount + amount
		col_transfer := Transfer{
			Type: "POD_deletion", Token: token, 
			Amount: amount, From: "POD Pool " + pod.PodId, To: pod.Creator,
		    Id: xid.New().String(), Date: date }
		transactions = append( transactions, col_transfer )
		wallet_creator.Balances[token] = balance_token }

	// Deactivate Pod //
	pod_list[ deletion.PodId ] = false
	pod.Pools.Collateral_Pool = make( map[string]float64 )
	pod.State.Status = "REMOVED"
	_, err5 := t.updatePod( stub, pod )
	if err5 != nil { return shim.Error( err5.Error() ) }
	err6 := t.updatePodList( stub, pod_list )
	if err6 != nil { return shim.Error( err6.Error() ) }
	update_pods[pod.PodId] = pod
	
	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	wallet_creator.Transaction = transactions
	err7 := t.updateUserWallet( stub, wallet_creator )
	if err7 != nil { return shim.Error( err7.Error() ) }
	update_wallets[wallet_creator.PublicId] = wallet_creator

	// Generate output //
	output := Output{ 
		 UpdatePods: update_pods, UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes ) 
}


/* -------------------------------------------------------------------------------------------------
investPOD: this function is called when an Investor wants to invest some amount in a given POD
           It receives pod tokens Args is an array containing a json with the following fields:
InvestorId              string                // Id of the creator of the POD
PodId                   string  			  // Id of the POD to delete
RateChange              string                // Rate of change of tokens with BC
Amount                  float64               // Amount of pod tokens to buy
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) investPOD( stub shim.ChaincodeStubInterface, 
											  args []string) pb.Response {

	// Retrieve the input information //
	investment := PodInvestment{}
	err1 := json.Unmarshal( []byte(args[0]), &investment )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	rateChange := investment.RateChange
	date := getTimeNow()
	update_pods := make( map[string]POD )
	update_wallets := make( map[string]MultiWallet )

	// Retrieve pod //
	pod, err2 := t.retrievePodInfo(stub, investment.PodId)
	if err2 != nil { shim.Error(err2.Error() ) }

	// Retrieve pod list //
	pod_list, err3 := t.retrievePodList(stub)
	if err3 != nil { shim.Error(err3.Error() ) }
	_, isInList := pod_list[investment.PodId]
	if !isInList {
		return shim.Error( "ERROR: POD ID IS NOT REGISTERED" ) }

	// Get the cost of buying that amount of pod tokens //
	pod, charged_token := t.investmentSwap( stub, pod, rateChange, investment.Amount )

	// Retrieve wallet of POD Investor //
	wallet_investor, err4 := t.retrieveUserWallet(stub, investment.InvestorId)
	if err4 != nil { shim.Error(err4.Error() ) }
	balance_investor := wallet_investor.Balances[ pod.Token ]
	if balance_investor.Amount < charged_token {
		return shim.Error( "ERROR: INVESTOR " + investment.InvestorId + "DOES NOT HOLD " +
		                   "SUFFICIENT FUNDS IN WALLET. ") }

	transactions := []Transfer{}
	// Transfer the funding token charged amount to the investor//
	balance_investor.Amount = balance_investor.Amount - charged_token
	outstanding := math.Max( pod.State.FundsRaised + charged_token - pod.Principal, 0.)
	pod.State.FundsRaised = pod.State.FundsRaised + charged_token - outstanding
	pod.State.Debt = pod.State.Debt + charged_token - outstanding
	pod.Pools.Funding_Pool = pod.Pools.Funding_Pool + charged_token - outstanding
	pod.Pools.Exchange_Pool = pod.Pools.Exchange_Pool + outstanding
	wallet_investor.Balances[ pod.Token ] = balance_investor
	invest_transfer := Transfer{
		Type: "POD_investment", Token: pod.Token, Amount: charged_token,
		From: investment.InvestorId, To: "POD Pool " + pod.PodId,
		Id: xid.New().String(), Date: date }
	transactions = append( transactions, invest_transfer )
	
	// Transfer the pod token purchased to the investor //
	investment_amount, isInList := pod.State.Investors[investment.InvestorId]
	balance_pod, _ := wallet_investor.Balances[ pod.PodId ]
	if !isInList { 
		// Create pod balance //
		balance_pod = Balance{ 
			Token: pod.PodId, Amount: 0., Credit_Amount: 0., 
			Staking_Amount: 0., Borrowing_Amount: 0., Type: "POD",
			PRIVI_lending: 0., PRIVI_borrowing: 0., 
			Collateral_Amount: 0., 
			Collaterals: make( map[string]float64) } 
		// Register investor //
		investment_amount = 0.
		pod.State.InvestorNum = pod.State.InvestorNum + 1 
	}
	pod.State.Investors[investment.InvestorId] = investment_amount +
									             investment.Amount 
	balance_pod.Amount = balance_pod.Amount + investment.Amount 
	wallet_investor.Balances[ pod.PodId ] = balance_pod
	pod.Pools.POD_Token_Pool = pod.Pools.POD_Token_Pool - investment.Amount
	pod.State.SupplyReleased = pod.State.SupplyReleased + investment.Amount 
	pod_token_transfer := Transfer{
		Type: "POD_investment", Token: pod.PodId, Amount: investment.Amount,
		From: "POD Pool " + pod.PodId, To: investment.InvestorId,
	    Id: xid.New().String(), Date: date  }
	transactions = append( transactions, pod_token_transfer )

	// Update Pod State on Blockchain //
	pod.State.Status = "ACTIVE"
	_, err5 := t.updatePod( stub, pod )
	if err5 != nil { return shim.Error( err5.Error() ) }
	update_pods[pod.PodId] = pod

	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	wallet_investor.Transaction = transactions
	err6 := t.updateUserWallet( stub, wallet_investor )
	if err6 != nil { return shim.Error( err6.Error() ) }
	update_wallets[wallet_investor.PublicId] = wallet_investor

	// Generate output //
	output := Output{ 
		 UpdatePods: update_pods, UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes ) 
}

/* -------------------------------------------------------------------------------------------------
interestPOD: this function is called when a POD should pay credit. 
PodId                   string  			  // Id of the POD to delete
RateChange              string                // Rate of change of tokens with BC
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) interestPOD( stub shim.ChaincodeStubInterface, 
											    args []string) ([]byte, error) {

	// Initialise the output //
	output := OutputLiquidation{}
	output.Liquidated = false
	output.Pod = POD{}
	output.UpdatePools = []LiquidityTokenPool{} 
	output.UpdateUsers = []MultiWallet{} 
	output_error, _ := json.Marshal( output )
	date := getTimeNow()

	// Retrieve the input information //
	input := PodLiquidator{}
	err1 := json.Unmarshal( []byte(args[0]), &input )
	if err1 != nil {
		return output_error, errors.New( "ERROR: GETTING THE INPUT INFORMATION. " + 
		                                  err1.Error() ) }
	podId := input.PodId
	rateChange := input.RateChange

	// Retrieve pod //
	pod, err2 := t.retrievePodInfo(stub, podId)
	if err2 != nil { return output_error, err2 }

	// Retrieve pod list //
	pod_list, err3 := t.retrievePodList(stub)
	if err3 != nil { return output_error, err3 }
	_, isInList := pod_list[podId]
	if !isInList {
		return output_error, errors.New( "ERROR: POD ID IS NOT REGISTERED" ) }

	// Retrieve interest to pay and check owner has funds to pay //
	interest_amount := pod.Interest * pod.State.FundsRaised
	wallet_pod_owner, err4 := t.retrieveUserWallet(stub, pod.Creator)
	if err4 != nil { return output_error, err4 }
	balance_pod_owner := wallet_pod_owner.Balances[pod.Token]
	if balance_pod_owner.Amount < interest_amount { 
		pod.State.MissingPayments = pod.State.MissingPayments + 1
		pod.State.Debt = pod.State.Debt + interest_amount
		output.Pod = pod
		outputBytes, _ := json.Marshal( output )
		return outputBytes, nil }
	output.UpdateUsers = append( output.UpdateUsers, wallet_pod_owner )    

	// Pod Owner holds fund, discount it from wallet //
	balance_pod_owner.Amount = balance_pod_owner.Amount - interest_amount
	pod.Pools.Interest_Pool = pod.Pools.Interest_Pool + interest_amount
	transaction_pay_interest := Transfer{
		Type: "POD_pay_interest", Token: pod.Token, Amount: interest_amount,
		From: pod.Creator, To: "POD Pool " + pod.PodId,
	    Id: xid.New().String(), Date: date }
	wallet_pod_owner.Balances[pod.Token] = balance_pod_owner
	wallet_pod_owner.Transaction = []Transfer{ transaction_pay_interest }
	err5 := t.updateUserWallet( stub, wallet_pod_owner )
	if err5 != nil { return output_error, err5 }

	// Interest to be given in pod tokens //
	token_interest := t.interestSwap( stub, pod, rateChange, 
									  interest_amount )

	// Divide interest between pod token holders. Give it in POD token //
	for investor_id, amount_invested := range(pod.State.Investors) {
		proportion := amount_invested / pod.State.FundsRaised 
		interest_get := proportion*token_interest
		// Retrieve investor wallet //
		wallet_investor, err6 := t.retrieveUserWallet(stub, investor_id)
		if err6 != nil { return output_error, err6 }
		balance_investor := wallet_investor.Balances[pod.Token] 
		balance_investor.Amount = balance_investor.Amount + interest_get
	    pod.Pools.POD_Token_Pool  = pod.Pools.POD_Token_Pool - interest_get
		transaction_get_interest := Transfer{
			Type: "POD_get_interest", Token: pod.Token, Amount: interest_get,
			From: "POD Pool " + pod.PodId, To: investor_id,
		    Id: xid.New().String(), Date: date }
		wallet_investor.Balances[pod.Token] = balance_investor
		wallet_investor.Transaction = []Transfer{ transaction_get_interest }
		// Update user wallet // 
		err7 := t.updateUserWallet( stub, wallet_investor )
		if err7 != nil { return output_error, err7 }
		output.UpdateUsers = append( output.UpdateUsers, wallet_investor )
	}
	// Divide interest between liquidity pools. Give it in POD token //
	for pool_id, amount_pool := range(pod.State.Liq_Pools) {
		proportion := amount_pool / pod.State.FundsRaised
		interest_get_pool :=  proportion*token_interest
		// Retrieve pool //
		pool, err8 := t.retrievePoolInfo(stub, pool_id)
		if err8 != nil { shim.Error(err8.Error()) }
		for LP_id, LP := range(pool.State.LiqProviders) {
			proportion_LP := LP.Amount / pool.State.Deposited
			interest_get := proportion_LP*interest_get_pool
			// Retrieve LP wallet //
			wallet_LP, err9 := t.retrieveUserWallet(stub, LP_id)
			if err9 != nil { return output_error, err9 }
			balance_LP := wallet_LP.Balances[pod.Token] 
			balance_LP.Amount = balance_LP.Amount + interest_get
			pod.Pools.POD_Token_Pool  = pod.Pools.POD_Token_Pool - interest_get
			transaction_get_interest := Transfer{
				Type: "POD_get_interest", Token: pod.Token, Amount: interest_get,
				From: "POD Pool " + pool.Id, To: LP_id,
			    Id: xid.New().String(), Date: date }
			wallet_LP.Balances[pod.Token] = balance_LP
			wallet_LP.Transaction = []Transfer{ transaction_get_interest }
			// Update user wallet // 
			err10 := t.updateUserWallet( stub, wallet_LP )
			if err10 != nil { return output_error, err10 }
			output.UpdateUsers = append( output.UpdateUsers, wallet_LP ) }
	}

	// Update Pod State on Blockchain //
	_, err11 := t.updatePod( stub, pod )
	if err11 != nil {  return output_error, err11 }
	output.Pod = pod
	outputBytes, _ := json.Marshal( output )
	return outputBytes, nil
}

/* -------------------------------------------------------------------------------------------------
swapPOD: this function is called when an Investor wants to swap some amount of pod tokens by
         another token. Args is an array containing a json with the following fields:
PodId            string          			// Id of the pod token
LiqPoolId        string                     // Id of the liquidity pool
InvestorId       string                     // Id of the investor 
Amount           float64        		    // Amount of pod token to swap
RateChange       map[string]float64         // Rate of change list with BC
Type             string                     // Type of the desired coin 
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) swapPOD( stub shim.ChaincodeStubInterface, 
								            args []string ) pb.Response {
	// Retrieve the input information //
	swap := PodSwapping{}
	err1 := json.Unmarshal( []byte(args[0]), &swap )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	date := getTimeNow()

	// Retrieve pod list //
	pod_list, err2 := t.retrievePodList(stub)
	if err2 != nil { shim.Error(err2.Error() ) }
	active, isInList := pod_list[swap.PodId]
	if !isInList {
		return shim.Error( "ERROR: POD ID IS NOT REGISTERED." ) }
	if !active {
		return shim.Error( "ERROR: POD ID IS NOT ACTIVE." ) }

	// Retrieve token POD //
	pod, err3 := t.retrievePodInfo(stub, swap.PodId)
	if err3 != nil { shim.Error(err3.Error() ) }

	// Retrieve pool to compute the swap and check token coincides //
	pool, err4 := t.retrievePoolInfo(stub, swap.LiqPoolId)
	if err4 != nil { shim.Error(err4.Error()) }
 
	// Retrieve wallet of POD Investor and check that has pod tokens //
	wallet_investor, err5 := t.retrieveUserWallet(stub, swap.InvestorId)
	if err5 != nil { shim.Error(err5.Error() ) }
	balance_pod, isInList := wallet_investor.Balances[ pod.PodId ]
	if !isInList {
		return shim.Error( "ERROR: INVESTOR " + swap.InvestorId + "DOES NOT HOLD " +
						   "POD TOKENS IN WALLET. ") }
	if balance_pod.Amount < swap.Amount {
		return shim.Error( "ERROR: INVESTOR " + swap.InvestorId + "DOES NOT HOLD " +
						   "SUFFICIENT POD TOKENS IN WALLET. ") }

	// Check that user has the balance of the desired token registered on wallet //
	balance_swap, isInList := wallet_investor.Balances[ pool.Token ]
	if !isInList {
		balance_swap = Balance{ 
			Token: swap.PodId, Amount: 0., Credit_Amount: 0., 
			Staking_Amount: 0., Borrowing_Amount: 0., 
			PRIVI_lending: 0., PRIVI_borrowing: 0., 
			Collateral_Amount: 0., Type: swap.Type,
			Collaterals: make( map[string]float64) } }
	
	// Compute swapping amount and check that there is funds on Lending Pool // 
	transactions := []Transfer{}
	fee_amount := 0.
	if swap.Type == "CRYPTO" {
		// Compute amount to get in the desired coin and funds and ratio //
		pod, swap_amount := t.swapFunctionCrypto( stub, pod, swap.RateChange, 
												  swap.Amount, pool.Token )
		fee_amount := swap_amount*pool.Fee
		var err6 error
		pool, err6 = t.swapLiquidityPool( stub, pool, swap_amount, fee_amount,
										  swap.Amount, pod.PodId, swap.RateChange )
		if err6 != nil { shim.Error(err6.Error()) }
		// Make transfers on user wallet //
		balance_pod.Amount = balance_pod.Amount - swap.Amount
		balance_swap.Amount = balance_swap.Amount + (swap_amount-fee_amount)
		transaction_give := Transfer{
			Type: "POD_swap_give", Token: swap.PodId, Amount: swap.Amount,
			From: swap.InvestorId, To: "Liquidity Pool " + pool.Id,
		    Id: xid.New().String(), Date: date }
		transaction_get := Transfer{
				Type: "POD_swap_get", Token: pool.Token, 
				Amount: (swap_amount-fee_amount),
				From: "Liquidity Pool " + pool.Id, To: swap.InvestorId }
		transaction_fee := Transfer{
			Type: "POD_swap_fee", Token: swap.PodId, Amount: fee_amount,
			From: swap.InvestorId, To: "Liquidity Pool " + pool.Id, 
		    Id: xid.New().String(), Date: date }
		transactions = append( transactions, transaction_give )
		transactions = append( transactions, transaction_get )
		transactions = append( transactions, transaction_fee )
		// Update wallet and lending pool //
		wallet_investor.Balances[ pod.PodId ] = balance_pod
		wallet_investor.Balances[ pool.Token ] = balance_swap
		wallet_investor.Transaction = transactions
		// Update Pod State on Blockchain //
		_, isInList :=  pod.State.Liq_Pools[ pool.Id ]
		if !isInList { pod.State.Liq_Pools[ pool.Id ] = 0. }
		pod.State.Liq_Pools[ swap.LiqPoolId ] = pod.State.Liq_Pools[ swap.LiqPoolId ] +
												 swap.Amount
		pod.State.Investors[ swap.InvestorId ] = pod.State.Investors[ swap.InvestorId ] -
										          swap.Amount
		// Check if investor should be deleted //
		if pod.State.Investors[ swap.InvestorId ] == 0 {
			delete( pod.State.Investors, swap.InvestorId )
			pod.State.InvestorNum = pod.State.InvestorNum-1 }
	}

	// Distribute Fee amount between LP //
	for LP_id, LP := range(pool.State.LiqProviders) {
		wallet_LP, err7 := t.retrieveUserWallet(stub, LP_id)
		if err7 != nil { shim.Error(err7.Error() ) }
		balance_LP := wallet_LP.Balances[ pool.Token ]
		proportion := LP.Amount / pool.State.Deposited
		fee_received := fee_amount * proportion
		balance_LP.Amount = balance_LP.Amount + fee_received
		transaction_LP := Transfer{
			Type: "liquidity_pool_swap_fee", Token: pool.Token, 
			Amount: fee_received,
			From: "Liquidity Pool " + pool.Id, To: LP_id,
		    Id: xid.New().String(), Date: date }
		// Update LP wallet on blockchain //
		wallet_LP.Balances[ pool.Token ] = balance_LP
		wallet_LP.Transaction = []Transfer{ transaction_LP }
		err8 := t.updateUserWallet( stub, wallet_LP )
		if err8 != nil { return shim.Error( err8.Error() ) }	
	}
	
	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	err9 := t.updateUserWallet( stub, wallet_investor )
	if err9 != nil { return shim.Error( err9.Error() ) }

	// Update Pod State on Blockchain //                           
	_, err10 := t.updatePod( stub, pod )
	if err10 != nil { return shim.Error( err10.Error() ) }

	// Update Pool State on Blockchain //
	_, err11 := t.updateLiquidityPool( stub, pool )
	if err11 != nil { return shim.Error( err11.Error() ) }

	// Prepare output //
	output := OutputSwapping{ UpdatePool: pool, UpdatePod: pod, 
		                      UpdateUser: wallet_investor }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
managerPOD: this function is called daily on a cron job to update the day of the POD and check if
            it is expiration day to raise the liquidation. It is called without any argument.
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) managerPOD( stub shim.ChaincodeStubInterface, 
											   args []string ) pb.Response {
	// Retrieve the input information //
	input := PodLiquidator{}
	err0 := json.Unmarshal( []byte(args[0]), &input )
	if err0 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err0.Error() ) }
	rateChange := input.RateChange

	// Retrieve pod list //
	pod_list, err1 := t.retrievePodList( stub )
	if err1 != nil { return shim.Error( err1.Error() ) }

	// Loop through all active Liquidity Pools and update day //
	output := map[string]OutputLiquidation{}
	
	for pod_id, isActive := range(pod_list) {
		// Initialise Pod Output //
		output_pod := OutputLiquidation{
			Liquidated: false, Pod: POD{},
			UpdatePools: []LiquidityTokenPool{}, 
			UpdateUsers: []MultiWallet{} }

		// Retrieve liquity pool  //
		pod, err2 := t.retrievePodInfo(stub, pod_id)
		if err2 != nil { shim.Error(err2.Error()) }
		if !isActive {continue}
		pod_day :=  pod.State.POD_Day + 1
		pod.State.POD_Day = pod_day
		output_pod.Pod = pod
		
		// CHECK IF PAYING INTEREST DAY //
		step := int64( pod.Duration / pod.Payments )
		if math.Mod( float64(pod_day), float64(step)) == 0 {
			args_input := PodLiquidator{
				PodId: pod.PodId, RateChange: rateChange }
			argsInputBytes, _ := json.Marshal( args_input )
			outputBytes, err := t.interestPOD(stub, 
				                  []string{string(argsInputBytes)} )
			if err != nil { shim.Error(err.Error()) }
			output_interest := OutputLiquidation{}
			json.Unmarshal( outputBytes, &output_interest )
			output_pod.Pod.State.POD_Day = pod_day
			output_pod = output_interest
		}

		// CHECK IF LIQUIDATION DAY // 
		if pod.State.POD_Day == pod.Duration {
			liquidator := PodLiquidator{
				PodId: pod.PodId, RateChange: rateChange }
			liquidatorBytes, _ := json.Marshal( liquidator )
			outputBytes, err3 := t.liquidatePOD( stub, 
				                append( []string{ string(liquidatorBytes) }, "EXPIRED") )
			if err3 != nil { shim.Error(err3.Error() ) }
			output_liq := OutputLiquidation{}
			json.Unmarshal( outputBytes, &output_liq )
			output_pod.Pod = output_liq.Pod
			output_pod.Pod.State.POD_Day = pod_day
			for _, pool := range(output_liq.UpdatePools) {
				output_pod.UpdatePools = append( output_pod.UpdatePools, pool ) }
			for _, user := range(output_liq.UpdateUsers) {
				output_pod.UpdateUsers = append( output_pod.UpdateUsers, user ) }
		}

		// Update pod on Blockchain //
		_, err4 := t.updatePod( stub, output_pod.Pod )
		if err4 != nil { return shim.Error( err4.Error() ) }
		output[pod_id] = output_pod
	}

	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
liquidatePOD: this function is called when a POD is being liquidated. I could happen at expiration
              date or if the collateral falls bellow the threshold levels.
PodId 			string 			 		// Id of the POD to revise liquidation
RateChange      map[string]float64      // Rates of changes of all currencies with BC
Type            string                  // Type of liquidation: liquidation or expiration
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) liquidatePOD( stub shim.ChaincodeStubInterface,
												 args []string ) ([]byte, error) {
	// Retrieve the input information //
	liquidator := PodLiquidator{}
	output := OutputLiquidation{}
	empty_output, _ := json.Marshal( output )
	err1 := json.Unmarshal( []byte(args[0]), &liquidator )
	if err1 != nil {
		return empty_output, errors.New( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	rateChange := liquidator.RateChange
	date := getTimeNow()

	// Retrieve pod list //
	pod_list, err2 := t.retrievePodList(stub)
	if err2 != nil { return empty_output, err2 }
	active, isInList := pod_list[liquidator.PodId]
	if !isInList {
		return empty_output, errors.New( "ERROR: POD ID IS NOT REGISTERED." ) }
	if !active {
		return empty_output, errors.New( "ERROR: POD ID IS NOT ACTIVE." ) }

	// Retrieve token POD //
	pod, err3 := t.retrievePodInfo(stub, liquidator.PodId)
	if err3 != nil { return empty_output, err3 }

	// Compute Collateral equivalent in BC //
	X_collateral_BC := 0.
	for token, amount := range(pod.Pools.Collateral_Pool) {
		X_collateral_BC = X_collateral_BC + amount*rateChange[token] }
	X_principal_BC := pod.Principal*rateChange[pod.Token] 
	X_collateral_BC = math.Min( X_collateral_BC, (1+pod.Interest)*X_principal_BC )

	// Initialise pools and users to update, and pct of collateral to take //
	output.UpdatePools = []LiquidityTokenPool{}
	output.UpdateUsers = []MultiWallet{}
	X_debt_BC := pod.State.Debt * rateChange[pod.Token]
	collateral_pct := math.Max(0, X_debt_BC-X_collateral_BC)/X_collateral_BC
	pod.State.Debt = math.Max(0, X_debt_BC-X_collateral_BC) /
	                 rateChange[pod.Token]
	// Retrieve pod owner wallet //
	wallet_owner, err4 := t.retrieveUserWallet(stub, pod.Creator)
	if err4 != nil { return empty_output, err4 }
	owner_trans := []Transfer{}

	// Take part of the collateral of the Pod Owner //
	for token, amount := range(pod.Pools.Collateral_Pool) {
		col_pod_owner := amount*(1-collateral_pct)
		balance_owner_tok := wallet_owner.Balances[token]
		balance_owner_tok.Amount = balance_owner_tok.Amount+col_pod_owner
		transaction_owner := Transfer{
			Type: "POD_liquidation_collateral", Token: token, 
			Amount: col_pod_owner, From: "POD Pool " + pod.PodId, 
			To: pod.Creator, Id: xid.New().String(), Date: date }
		owner_trans = append( owner_trans, transaction_owner ) 
		wallet_owner.Balances[token] = balance_owner_tok }
	// Update MultiWallet of owner //
	wallet_owner.Transaction = owner_trans
	err5 := t.updateUserWallet( stub, wallet_owner )
	if err5 != nil { return empty_output, err5 }
	output.UpdateUsers = append( output.UpdateUsers, wallet_owner )

	/////////////////// (1) UPDATE INVESTORS //////////////////////////////////
	for investor, amount_holder := range(pod.State.Investors) {
		// Retrieve investor wallet //
		wallet_investor, err6 := t.retrieveUserWallet(stub, investor)
		if err6 != nil { return empty_output, err6 }
		proportion := amount_holder / pod.State.SupplyReleased
		investor_trans := []Transfer{}
		// ---> Distribute collateral //
		for token, amount := range(pod.Pools.Collateral_Pool) {
			col_investor := amount*collateral_pct*proportion
			balance_investor_tok := wallet_investor.Balances[token]
			balance_investor_tok.Amount = balance_investor_tok.Amount+col_investor 
			transaction_investor := Transfer{
				Type: "POD_liquidation_collateral", Token: token, 
				Amount: col_investor, From: "POD Pool " + pod.PodId, 
				To: investor, Id: xid.New().String(), Date: date }
			investor_trans = append( investor_trans, transaction_investor )
			wallet_investor.Balances[token] = balance_investor_tok }
		// ---> Distribute Rest of Liquidity in the Pool //
		rest_amount := ( pod.Pools.Funding_Pool + pod.Pools.Exchange_Pool +
						 pod.Pools.Interest_Pool ) * proportion
		balance_investor := wallet_investor.Balances[pod.Token]
		balance_investor.Amount = balance_investor.Amount+rest_amount 
		transaction_investor_funds := Transfer{
			Type: "POD_liquidation_funds", Token: pod.Token, Amount: rest_amount,
			From: "POD Pool " + pod.PodId, To: investor, Id: xid.New().String(), Date: date }
		investor_trans = append( investor_trans, transaction_investor_funds )
		wallet_investor.Balances[ pod.Token ] = balance_investor
		// ---> Burn pod Tokens  //
		delete(wallet_investor.Balances, pod.PodId)
		transaction_burn := Transfer{
			Type: "POD_burning_token", Token: pod.PodId, 
			Amount: pod.State.Investors[investor],
			From: investor, To: "", Id: xid.New().String(), Date: date }
		pod.State.Investors[investor] = 0.
		investor_trans = append( investor_trans, transaction_burn )
		wallet_investor.Transaction = investor_trans
		// Update MultiWallet of user //
		err7 := t.updateUserWallet( stub, wallet_investor )
		if err7 != nil { return empty_output, err7 }
		output.UpdateUsers = append( output.UpdateUsers, wallet_investor )
	}
	
	////////////////// (2) UPDATE LIQUIDITY POOLS /////////////////////////////
	for pool_id, amount_pool := range(pod.State.Liq_Pools) {
		// Retrieve pool from blockchain //
		pool, err8 := t.retrievePoolInfo(stub, pool_id)
		if err8 != nil { return empty_output, err8  }
		proportion := amount_pool / pod.State.SupplyReleased
		// ---> Distribute collateral //
		for token, amount := range(pod.Pools.Collateral_Pool) {
			col_pool := amount*collateral_pct*proportion
			amount_pool, isInList := pool.State.Liq_Tokens[token]
			if !isInList { amount_pool = 0 }
			pool.State.Liq_Tokens[token] = amount_pool + col_pool }
		// ---> Distribute Rest of Liquidity in the Pool //
		rest_amount := ( pod.Pools.Funding_Pool + pod.Pools.Exchange_Pool +
						 pod.Pools.Interest_Pool ) * proportion
		amount_pool, isInList := pool.State.Liq_Tokens[pod.Token]
		if !isInList { amount_pool = 0 }
		pool.State.Liq_Tokens[pod.Token] = amount_pool + rest_amount
		// Update Pool State on Blockchain //
		delete( pool.State.Liq_Tokens, pod.PodId )
		_, err9 := t.updateLiquidityPool( stub, pool )
		if err9 != nil { return empty_output, err9  }
		output.UpdatePools = append( output.UpdatePools, pool )		
	}

	// Empty POD pools // 
	pod_list[liquidator.PodId] = false
	pod.Pools.Collateral_Pool = make( map[string]float64 )
	pod.Pools.Funding_Pool = 0.
	pod.Pools.Exchange_Pool = 0.
	pod.Pools.Interest_Pool = 0.
	pod.Pools.POD_Token_Pool = 0.
	pod.State.Status = args[1]

	// Update Pod State on Blockchain //                           
	_, err10 := t.updatePod( stub, pod )
	if err10 != nil { return empty_output, err10  }

	output.Pod = pod
	output.Liquidated = true
	ouputBytes, _ := json.Marshal( output )

	return ouputBytes, nil

}
	
/* -------------------------------------------------------------------------------------------------
checkPODLiquidation: this function is called when a POD is identified to be under the collateral
                     liquidity threshold. Input: array containing a json with fields:
PodId 			string 			 		// Id of the POD to revise liquidation
RateChange      map[string]float64      // Rates of changes of all currencies with BC
------------------------------------------------------------------------------------------------- */

func (t *PodSwappingSmartContract) checkPODLiquidation( stub shim.ChaincodeStubInterface,
													    args []string ) pb.Response {
	// Retrieve the input information //
	liquidator := PodLiquidator{}
	err1 := json.Unmarshal( []byte(args[0]), &liquidator )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	output := OutputLiquidation{}
	rateChange := liquidator.RateChange

	// Retrieve pod list //
	pod_list, err2 := t.retrievePodList(stub)
	if err2 != nil { shim.Error(err2.Error() ) }
	active, isInList := pod_list[liquidator.PodId]
	if !isInList {
		return shim.Error( "ERROR: POD ID IS NOT REGISTERED." ) }
	if !active {
		return shim.Error( "ERROR: POD ID IS NOT ACTIVE." ) }

	// Retrieve token POD //
	pod, err3 := t.retrievePodInfo(stub, liquidator.PodId)
	if err3 != nil { shim.Error(err3.Error() ) }

	// Compute Collateral equivalent in BC //
	X_collateral_BC := 0.
	for token, amount := range(pod.Pools.Collateral_Pool) {
		X_collateral_BC = X_collateral_BC + amount*rateChange[token] }
	X_principal_BC := pod.Principal*rateChange[pod.Token] 
	X_collateral_BC = math.Min( X_collateral_BC, (1+pod.Interest)*X_principal_BC )

	// Compute CCR and check condition //
	CCR := X_collateral_BC / X_principal_BC
	if CCR >= pod.P_liquidation {
		output.Liquidated = false
		outputBytes, _ := json.Marshal( output )
		return shim.Success( outputBytes ) }

	outputBytes, err4 := t.liquidatePOD( stub, append( args, "LIQUIDATED") )
	if err4 != nil { shim.Error(err4.Error() ) }

	return shim.Success( outputBytes )												  
}




/* -------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------- */

func main() {
	err := shim.Start(&PodSwappingSmartContract{})
	if err != nil {
		fmt.Errorf("Error starting Pod Swapping chaincode: %s", err)
	}
}

