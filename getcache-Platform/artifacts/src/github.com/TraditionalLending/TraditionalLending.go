package main

import (
	"encoding/json"
	"fmt"
	"math"
	"github.com/rs/xid"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


/* -------------------------------------------------------------------------------------------------
Init:  this function is called at PRIVI Blockchain Deployment and initialises the Traditional 
	   Lending Smart Contract. This smart contract is the responsible to manage the liquidity 
	   Lending pools and implement the traditional lending crypto model. 
	   Args: array containing a string:
PrivateKeyID           string   // Private Key of the admin of the smart contract
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) Init( stub shim.ChaincodeStubInterface ) pb.Response {
												_, args := stub.GetFunctionAndParameters()

	// Store in the state of the smart contract the Private Key of Admin //
	err1 := stub.PutState( IndexAdmin, []byte( args[0]) )
	if err1 != nil {
		return shim.Error( "ERROR: SETTING THE ADMIN PRIVATE KEY: " +
	                        err1.Error() ) }
	if args[1] == "UPGRADE" { return shim.Success(nil) }
 
	// Initialise list of tokens in the Smart Contract as empty //
	token_list, _ := json.Marshal( []string{} )
	err2 := stub.PutState( IndexTokenList, token_list )
	if err2 != nil {
		return shim.Error( "ERROR: INITIALISING THE TOKEN LIST: " +
						    err2.Error() ) }
	
	// Initialise list of borrowers in the Smart Contract as empty //
	borrowing_list, _ := json.Marshal( []string{} )
	err3 := stub.PutState( IndexBorrowingList, borrowing_list )
	if err3 != nil {
		return shim.Error( "ERROR: INITIALISING THE BORROWING LIST: " +
						    err3.Error() ) }
	
	// Initialise list of staking users in the Smart Contract as empty //
	staking_list, _ := json.Marshal( []string{} )
	err4 := stub.PutState( IndexStakingList, staking_list )
	if err4 != nil {
		return shim.Error( "ERROR: INITIALISING THE STAKING LIST: " +
	                       err4.Error() ) }

	return shim.Success(nil)
}


/* -------------------------------------------------------------------------------------------------
Invoke:  this function is the router of the different functions supported by the Smart Contract.
		 It receives the input from the controllers and ensure the correct calling of the functions.
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// Retrieve function and arguments //
	function, args := stub.GetFunctionAndParameters()

	// Retrieve caller of the function //
	caller, err1 := CallerCN(stub)
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE CALLER OF THE TRANSFER " +
						   "FUNCTION. " + err1.Error() )
	}
	// Retrieve admin of the smart contract //
	adminBytes, err2 := stub.GetState( IndexAdmin )
	if err2 != nil {
		return shim.Error( "ERROR: GETTING THE ADMIN OF THE SMART " +
						   "CONTRACT. " + err2.Error() )
	}
	admin := string( adminBytes )
	// Verify that the caller of the function is admin //
	if caller != admin {
		return shim.Error( "ERROR: CALLER " + caller + " DOES NOT HAVE PERMISSION " )
	}

	// Call the proper function //
	switch function {
		case "registerToken":
			return t.registerToken(stub, args)
		case "getTokenList":
			return t.getTokenList(stub, args)
		case "getTokenPool":
			return t.getTokenPool(stub, args)
		case "getDemandRatios":
			return t.getDemandRatios(stub, args)
		case "borrowFunds":
			return t.borrowFunds(stub, args)
		case "stakeToken":
			return t.stakeToken(stub, args)
		case "unStakeToken":
			return t.unStakeToken(stub, args)
		case "payInterests":
			return t.payInterests(stub, args)
		case "depositCollateral":
			return t.depositCollateral(stub, args)
		case "withdrawCollateral":
			return t.withdrawCollateral(stub, args)
		case "repayFunds":
			return t.repayFunds(stub, args)
		case "checkLiquidation":
			return t.checkLiquidation(stub, args)
		case "updateRiskParameters":
			return t.updateRiskParameters(stub, args)
		case "updatePoolBytes":
			return t.updatePoolBytes(stub, args)
	}

	// If function does not exist, retrieve error//
	return shim.Error( "ERROR: INCORRECT FUNCTION NAME " + function )
}


/* -------------------------------------------------------------------------------------------------
registerToken:  this function is called when a new token is listed on the PRIVI Blockchain. It can 
				only be called by Admin. At deployment, PRIVI Coin and Base Coin are created in the 
				system. Input: array containing a json with fields:
Token           string   // Symbol of a new token registered in the smart contract
Reserve         string   // Initial Reserve in the Lending Pool for the Token
RiskParameters  Risk     // Json containing the risk parameters of the token
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) registerToken( stub shim.ChaincodeStubInterface,
												 		 args []string ) pb.Response {

	// Retrieve arguments of the input in a Token Model struct //
	if len(args) != 1 {
		return shim.Error( "ERROR: REGISTER TOKEN FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENTS." ) }
	pool := LendingPool{}
	json.Unmarshal([]byte(args[0]), &pool)

	// Check if Token is already registered on Blockchain //
	token_list, err1 := t.checkTokenListed( stub, pool.Token )
	_, isListed := token_list[ pool.Token ]
	if err1 != nil && isListed { return shim.Error( err1.Error()) }
	if isListed { 
		return shim.Error( "ERROR: LENDING POOL FOR TOKEN" +  pool.Token + 
						   " IS ALREADY LISTED." ) }

	// Initialise Lending Pool for the Token and store state on Blockchain //
	pool.Loaned = 0.
	pool.Staked = 0.
	pool.Collateral = 0.
	poolsBytes2, _ := json.Marshal( pool )
	err2 := stub.PutState( IndexLendingPools + pool.Token, poolsBytes2 )
	if err2 != nil {
		return shim.Error( "ERROR: UPDATING LENDING POOLS STATE TO " +
						   "BLOCKCHAIN. " + err2.Error() ) }

	// Add token to the token list of the smart contract //
	token_list[ pool.Token ] = true
	tokenListBytes2, _ := json.Marshal( token_list )
	err3 := stub.PutState( IndexTokenList, tokenListBytes2 )
	if err3 != nil  {
		return shim.Error( "ERROR: ADDING TOKEN " + pool.Token + 
						   " THE TOKEN LIST: " ) }

	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
getTokenList:  this function returns the list of tokens listed in the Smart Contract. This 
               function does not need to be called with any argument.
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) getTokenList( stub shim.ChaincodeStubInterface,
													    args []string ) pb.Response {
	// Retrieve list of tokens from the Blockchain //
	tokenList, err := stub.GetState( IndexTokenList )
	if err != nil {
		return shim.Error( "ERROR: RETRIEVING THE TOKEN LIST" ) }
	return shim.Success( tokenList )
}


/* -------------------------------------------------------------------------------------------------
getDemandRatios:  this function returns the list of demand ratios for each token. This 
              	  function does not need to be called with any argument.
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) getDemandRatios( stub shim.ChaincodeStubInterface,
													       args []string ) pb.Response {
	// Retrieve list of tokens from Blockchain //
	token_list := map[string]bool{}
	tokenListBytes, err1 := stub.GetState( IndexTokenList )
	if err1 != nil { return shim.Error( err1.Error()) }
	json.Unmarshal( tokenListBytes, &token_list )
 
	// Compute Demand Ratios // 
	demand_ratios := map[string]float64{}
	for token, isListed := range(token_list) {
		if !isListed { continue }
		token_pool, err2 := t.retrieveTokenPool( stub, token )
		if err2 != nil { return shim.Error( err2.Error())  }
		demand := 0.
		if token_pool.Loaned + token_pool.Reserve != 0 {
			demand = token_pool.Loaned / ( token_pool.Loaned + token_pool.Reserve ) }
		demand_ratios[token] = demand
	}

	demandBytes, _ := json.Marshal( demand_ratios )
	return shim.Success( demandBytes )
}

/* -------------------------------------------------------------------------------------------------
getTokenPool:  this function returns the state of the Lending Pool for the Token in question.
		       Args: array containing
args[0]               string   // Symbol of the token to release the information	   
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) getTokenPool( stub shim.ChaincodeStubInterface,
														args []string ) pb.Response {
	// Retrieve information from the input //
	if len(args) != 1 {
		return shim.Error( "ERROR: GETTOKENPOOL FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENT." ) }

	// Check that Token is registered on Blockchain //
	_, err1 := t.checkTokenListed( stub, args[0] )
	if err1 != nil { return shim.Error( err1.Error()) }

	// Retrieve token lending pool information //
	poolsBytes, err2 := stub.GetState( IndexLendingPools + args[0] )
	if err2 != nil {
		return shim.Error( "ERROR: FAILED TO GET THE TOKEN LENDING POOL " + 
						   "INFORMATION." + err2.Error() ) }
	return shim.Success( poolsBytes )
}


/* -------------------------------------------------------------------------------------------------
borrowFunds:  this function is called when a borrower wants to borrow funds.
              Input: array containing a json with fields:
PublicId        string    				// Id of the user requesting the loan
Token           string   			    // Symbol of the lending token
Amount          float64  			    // Amount requested to borrow
Collaterals     map[string]float64      // Amount of collateral for each symbol
RateChange      map[string]float64      // Rates of changes of all currencies with BC
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) borrowFunds( stub shim.ChaincodeStubInterface,
												       args []string ) pb.Response {
	// Retrieve arguments of the input in a Request Token Model struct //
	if len(args) != 1 {
		return shim.Error( "ERROR: LENDFUNDS FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENT." )
	}
	request := RequestLoan{}
	json.Unmarshal([]byte(args[0]), &request)
	date := getTimeNow()

	// Check that Token is registered on Blockchain //
	_, err1 := t.checkTokenListed( stub, request.Token )
	if err1 != nil { return shim.Error( err1.Error()) }

	// Retrieve loan token pool //
	pool_loan, err2 := t.retrieveTokenPool( stub, request.Token )
	if err2 != nil { return shim.Error( err2.Error()) }

	// Check that there are enough reserves to satisfy amount requested //
	if request.Amount > pool_loan.Reserve {
		return shim.Error( "ERROR: THERE ARE NOT ENOUGH FUNDS ON THE LENDING " + 
		                   "POOL TO SATISFY DEMAND." ) }

	// Retrieve the MultiWallet of user. Invoke CoinBalance Smart Contract // 
	user_wallet, err3 := t.retrieveUserWallet( stub, request.PublicId )
	if err3 != nil { return shim.Error( err3.Error()) }

	// Compute collateral level and transfer funds from user to Lending Pools //
	transactions := []Transfer{}
	collateral_BC := 0.
	balance_loan := user_wallet.Balances[request.Token]
	for token, amount := range(request.Collaterals) {
		// Balance and Pool of token to update //
		balance_token := user_wallet.Balances[token]
		pool_token, err4 := t.retrieveTokenPool( stub, request.Token )
		if err4 != nil { return shim.Error( err4.Error()) }
		//Check if collateral is in list and add collateral //
		col, isInList := balance_loan.Collaterals[token]
		if !isInList { balance_loan.Collaterals[token] = 0. }
		balance_loan.Collaterals[token] = col + amount
		collateral_BC = collateral_BC +
						balance_loan.Collaterals[token]  * request.RateChange[token]
		// Check user token balance //
		if amount > balance_token.Amount {
			return shim.Error( "ERROR: THE USER DOES NOT THE AMOUNT OF COLLATERAL " +
								  "FOR TOKEN " + token ) }
	    // Update Pools and Balances //
		balance_token.Amount = balance_token.Amount - amount 
		balance_token.Collateral_Amount = balance_token.Collateral_Amount + amount
		pool_token.Reserve = pool_token.Reserve + amount
		pool_token.Collateral = pool_token.Collateral + amount
		col_transfer := Transfer{
			Type: "traditional_collateral_deposit", Token: token, From: request.PublicId,
			To: "Lending Pool " + token, Amount: amount, 
			Id: xid.New().String(), Date: date }
		transactions = append( transactions, col_transfer )
		// Update balance and pool //
		user_wallet.Balances[token] = balance_token
		_, err5 := t.updatePool( stub, pool_token )
		if err5 != nil { return shim.Error( err5.Error()) }
	}

	// Check if Collateral Level is above the thresholds // 
	risk_params := pool_loan.RiskParameters
	principal_BC := request.Amount * request.RateChange[request.Token]
	CCR := collateral_BC / principal_BC
	if CCR < risk_params.Col_required {
		return shim.Error( "ERROR: THE COLLATERAL IS NOT ENOUGH TO SATISFY THE " +
		                   "CCR REQUIREMENT" ) }

	// If the collateral condition is satisfied transfer funds to user //
	balance_loan.Amount = balance_loan.Amount + request.Amount 
	balance_loan.Borrowing_Amount = balance_loan.Borrowing_Amount + request.Amount 
	pool_loan.Reserve = pool_loan.Reserve - request.Amount 
	pool_loan.Loaned = pool_loan.Loaned + request.Amount 
	loan_transfer := Transfer{
		Type: "traditional_borrowing", Token: request.Token, 
		From: "Lending Pool " + request.Token, To: request.PublicId, 
		Amount: request.Amount, Id: xid.New().String(), Date: date }
	transactions = append( transactions, loan_transfer )
	
	// Update Loan token Pool // 
	_, err6 := t.updatePool( stub, pool_loan )
	if err6 != nil { return shim.Error( err6.Error()) }

	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	user_wallet.Balances[request.Token] = balance_loan
	user_wallet.Transaction = transactions
	err7 := t.updateUserWallet( stub, user_wallet )
	if err7 != nil { return shim.Error( err7.Error()) }

	// Generate response for outputs to update //
	update_wallets := make( map[string]MultiWallet )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	output := Output{ UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}

/* -------------------------------------------------------------------------------------------------
stakeToken:  this function is called when a users wants to stake a certain amount of token.
             Input: array containing a json with fields:
PublicId        string    				// Id of the user requesting the loan
Token           string   			    // Symbol of the lending token
Amount          float64  			    // Amount requested to borrow
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) stakeToken( stub shim.ChaincodeStubInterface,
													  args []string ) pb.Response {
	// Retrieve arguments of the input in a Request Token Model struct //
	if len(args) != 1 {
		return shim.Error( "ERROR: STAKETOKEN FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENT." ) }
	staking := StakeToken{}
	json.Unmarshal([]byte(args[0]), &staking)
	date := getTimeNow()

	// Check that Token is registered on Blockchain //
	_, err1 := t.checkTokenListed( stub, staking.Token )
	if err1 != nil { return shim.Error( err1.Error()) }

	// Retrieve staking token pool //
	pool_stake, err2 := t.retrieveTokenPool( stub, staking.Token )
	if err2 != nil { return shim.Error( err2.Error()) }

	// Retrieve the MultiWallet of user. Invoke CoinBalance Smart Contract // 
	user_wallet, err3 := t.retrieveUserWallet( stub, staking.PublicId )
	if err3 != nil { return shim.Error( err3.Error()) }
	balance_token := user_wallet.Balances[ staking.Token ]

	// Check that user has the amount of token to stake //
	if balance_token.Amount < staking.Amount {
		return shim.Error( "ERROR: USER " + staking.PublicId + "DOES NOT" +
						   "HOLD THE FUNDS TO STAKE " + staking.Token  ) }

	// Update user balance and lending pool //
	balance_token.Amount = balance_token.Amount - staking.Amount 
	balance_token.Staking_Amount = balance_token.Staking_Amount + staking.Amount 
	pool_stake.Reserve = pool_stake.Reserve + staking.Amount 
	pool_stake.Staked = pool_stake.Staked + staking.Amount 
	staking_transfer := Transfer{
		Type: "staking", Token: staking.Token, 
		From: staking.PublicId, To: "Lending Pool " + staking.Token, 
		Amount: staking.Amount, Id: xid.New().String(), Date: date }
	transactions := []Transfer{staking_transfer}

	// Update Staking token Pool // 
	_, err4 := t.updatePool( stub, pool_stake )
	if err4 != nil { return shim.Error( err4.Error()) }

	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	user_wallet.Balances[staking.Token] = balance_token
	user_wallet.Transaction = transactions
	err5 := t.updateUserWallet( stub, user_wallet )
	if err5 != nil { return shim.Error( err5.Error()) }

	// Generate response for outputs to update //
	update_wallets := make( map[string]MultiWallet )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	output := Output{ UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}

/* -------------------------------------------------------------------------------------------------
unStakeToken:  	this function is called when a users wants to unstake a certain amount of token.
             	Input: array containing a json with fields:
PublicId        string    				// Id of the user requesting the loan
Token           string   			    // Symbol of the lending token
Amount          float64  			    // Amount requested to borrow
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) unStakeToken( stub shim.ChaincodeStubInterface,
														args []string ) pb.Response {

	// Retrieve arguments of the input in a Request Token Model struct //
	if len(args) != 1 {
		return shim.Error( "ERROR: STAKETOKEN FUNCTION SHOULD BE CALLED " +
							"WITH ONE ARGUMENT." ) }
	unstaking := StakeToken{}
	json.Unmarshal([]byte(args[0]), &unstaking)
	date := getTimeNow()

	// Check that Token is registered on Blockchain //
	_, err1 := t.checkTokenListed( stub, unstaking.Token )
	if err1 != nil { return shim.Error( err1.Error()) }

	// Retrieve staking token pool //
	pool_unstake, err2 := t.retrieveTokenPool( stub, unstaking.Token )
	if err2 != nil { return shim.Error( err2.Error()) }

	// Retrieve the MultiWallet of user. Invoke CoinBalance Smart Contract // 
	user_wallet, err3 := t.retrieveUserWallet( stub, unstaking.PublicId )
	if err3 != nil { return shim.Error( err3.Error()) }
	balance_token := user_wallet.Balances[ unstaking.Token ]

	// Check that user has the amount to unstake and funds on the Pool//=
	if balance_token.Staking_Amount < unstaking.Amount {
		return shim.Error( "ERROR: USER " + unstaking.PublicId + "CAN NOT" +
						   "UNSTAKE TOKEN MORE THAN IT HAS " + unstaking.Token) }
	if pool_unstake.Reserve < unstaking.Amount {
		return shim.Error( "ERROR: TOKEN " + pool_unstake.Token + "HAS NOT" +
						   "ENOUGH RESERVE" ) }

	// Update user balance and lending pool //
	balance_token.Amount = balance_token.Amount + unstaking.Amount 
	balance_token.Staking_Amount = balance_token.Staking_Amount - unstaking.Amount 
	pool_unstake.Reserve = pool_unstake.Reserve - unstaking.Amount 
	pool_unstake.Staked = pool_unstake.Staked - unstaking.Amount 
	unstaking_transfer := Transfer{
			Type: "unstaking", Token: unstaking.Token, 
			From: "Lending Pool " + unstaking.Token, To: unstaking.PublicId, 
			Amount: unstaking.Amount, Id: xid.New().String(), Date: date }
	transactions := []Transfer{unstaking_transfer}
	user_wallet.Balances[unstaking.Token] = balance_token

	// Update Staking token Pool // 
	_, err4 := t.updatePool( stub, pool_unstake )
	if err4 != nil { return shim.Error( err4.Error()) }

	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	user_wallet.Balances[unstaking.Token] = balance_token
	user_wallet.Transaction = transactions
	err5 := t.updateUserWallet( stub, user_wallet )
	if err5 != nil { return shim.Error( err5.Error()) }

	// Generate response for outputs to update //
	update_wallets := make( map[string]MultiWallet )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	output := Output{ UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
payInterests:  this function is called scheduled in a daily basis and makes the borrowers pay the 
			   interest of the loan and staking holders receive the staking interest. This takes an
			   array containing a json with the following fields:
LendingInterest      map[string]float64  	   // Daily lending interest for each token
StakingInterest      map[string]float64  	   // Daily staking interest for each token
RateChange      	 map[string]float64        // Rates of changes of all currencies with BC
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) payInterests( stub shim.ChaincodeStubInterface,
													    args []string ) pb.Response {
	// Retrieve arguments of the input //
	if len(args) != 1 {
		return shim.Error( "ERROR: PAYINTERESTS FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENT." ) }
	interest_rates := InterestRates{}
	json.Unmarshal( []byte(args[0]), &interest_rates )
	RateChange := interest_rates.RateChange
	date := getTimeNow()						
	
	// Retrieve all the MultiWallet of users. Invoke CoinBalance Smart Contract // 
	chainCodeArgs := util.ToChaincodeArgs("getWallets", "")
	response := stub.InvokeChaincode( COIN_BALANCE_CHAINCODE, chainCodeArgs, 
									  CHANNEL_NAME )
	if response.Status != shim.OK {
		return shim.Error( "ERROR INVOKING THE GETWALLETS CHAINCODE TO " +
						   "GET THE WALELT OF USERS." )
	}
	user_wallets := []MultiWallet{}
	json.Unmarshal( response.Payload, &user_wallets )

	// Retrieve list of tokens from the Blockchain //
	tokenListBytes, err1 := stub.GetState( IndexTokenList )
	if err1 != nil {
		return shim.Error( "ERROR: RETRIEVING THE TOKEN LIST. " +
		                    err1.Error() ) }
	token_list := []string{}
	json.Unmarshal(tokenListBytes, &token_list)

	// LOOP THROUGH ALL TOKENS //
	for _, token := range(token_list) {
		P_lending := interest_rates.LendingInterest[token]
		P_staking := interest_rates.StakingInterest[token]
		// CHECK THAT THERE ARE ENOUGH FUNDS TO PAY INTERESTS //
		token_pool, err2 := t.retrieveTokenPool( stub, token )
		if err2 != nil { return shim.Error( err2.Error()) }
		funds_to_pay := token_pool.Staked * P_staking
		funds_available := token_pool.Reserve +
		                   token_pool.Loaned * P_lending
		if funds_to_pay > funds_available {
			return shim.Error( "ERROR: NOT ENOUGHT FUNDS TO PAY " + 
						   	   "INTEREST RATES." ) }
		// LOOP THROUGH USER WALLETS //
		for i, wallet := range(user_wallets) {
			//////////////// (1) CHARGE LENDING INTEREST  ///////////////////////
			balance_token := wallet.Balances[token]
			charge_amount := balance_token.Borrowing_Amount * P_lending
			// Charge user. First case: it has the interest amount in the balance //
			if charge_amount > 0. && balance_token.Amount > charge_amount {
				balance_token.Amount = balance_token.Amount - charge_amount
				token_pool.Reserve = token_pool.Reserve + charge_amount 
				interest_transfer := Transfer{
					Type: "lending_interest", Token: token, 
					From: wallet.PublicId, To: "Lending Pool " + token, 
					Amount: charge_amount, Id: xid.New().String(), Date: date   }
				wallet.Transaction = append( wallet.Transaction, interest_transfer )
				wallet.Balances[token] = balance_token
			// Charge user. Second case: discount equivalent from collaterals //
			} else if charge_amount > 0. {
				charge_amount_BC := charge_amount*RateChange[token]
				for token_col, amount := range(balance_token.Collaterals) {
					col_BC := amount * RateChange[token_col]
					left_BC := math.Max( 0, col_BC - charge_amount_BC ) // collateral left in BC
					charged_BC := col_BC - left_BC // collateral charged in BC
					charge_amount_BC = charge_amount_BC - charged_BC // remaining to charge
					left_rate := left_BC / col_BC // left rate 
					// Update User Wallets and Pool Liquidity accordingly //
					balance_token_col := wallet.Balances[token_col]
					token_col_pool, err3 := t.retrieveTokenPool( stub, token_col )
					if err3 != nil { return shim.Error( err3.Error()) }
					balance_token.Collaterals[token_col] = amount * left_rate
					balance_token_col.Collateral_Amount = balance_token_col.Collateral_Amount - 
														  amount * (1-left_rate)
					token_col_pool.Reserve = token_col_pool.Reserve + amount * (1-left_rate)
					// Transfer object //
					interest_transfer := Transfer{
						Type: "collateral_interest", Token: token_col, 
						From: wallet.PublicId, To: "Lending Pool " + token_col, 
						Amount: amount * (1-left_rate) }
					wallet.Balances[token_col] = balance_token_col
					wallet.Transaction = append( wallet.Transaction, interest_transfer )
					_, err4 := t.updatePool( stub, token_pool )
					if err4 != nil { return shim.Error( err4.Error()) } 
					// Check if there is stil amount to charge user //
					if charge_amount_BC <= 0. {break;}
				} 
			}
			/////////////////// (2) PAY STAKING INTEREST  ////////////////////////
			if balance_token.Staking_Amount > 0. {
				staking_payment := balance_token.Staking_Amount * P_staking
				if staking_payment > token_pool.Reserve {
					return shim.Error( "ERROR: NOT ENOUGHT FUNDS IN THE LENDING" + 
						   	  		   "POOL OF " + token + " TO PAY STAKING INTEREST." )
				}
				balance_token.Amount =  balance_token.Amount + staking_payment
				token_pool.Reserve = token_pool.Reserve - staking_payment
				staking_transfer := Transfer{
					Type: "staking_interest", Token: token, 
					From:"Lending Pool " + token, To: wallet.PublicId, 
					Amount: staking_payment, Id: xid.New().String(), Date: date }
				wallet.Transaction = append( wallet.Transaction, staking_transfer )
			}
			// Update user wallet //
			wallet.Balances[token] = balance_token
			user_wallets[i] = wallet
		}
		// Update token Pool // 
		_, err5 := t.updatePool( stub, token_pool )
		if err5 != nil { return shim.Error( err5.Error()) }
	}

	// Update wallets of users on blockchain //
	update_wallets := make( map[string]MultiWallet )
	for _, wallet := range(user_wallets) {
		update_wallets[ wallet.PublicId ] = wallet
		err6 := t.updateUserWallet( stub, wallet )
		if err6 != nil { return shim.Error( err6.Error()) } }

	// Generate response for outputs to update //
	output := Output{ UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
depositCollateral:  this function is called when a borrower wants to deposit additional collateral.
              		Input: array containing a json with fields:
PublicId        string    				// Id of the user depositing collateral
Token           string   			    // Symbol of the borrowing token
Collaterals     map[string]float64      // Amount of collateral deposited for each symbol
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) depositCollateral( stub shim.ChaincodeStubInterface,
													   		 args []string ) pb.Response {
	// Retrieve arguments of the input in a Deposit Token Model struct //
	if len(args) != 1 {
	return shim.Error( "ERROR: DEPOSITCOLLATERAL FUNCTION SHOULD BE CALLED " +
						"WITH ONE ARGUMENT." ) }
	deposit := CollateralDeposit{}
	json.Unmarshal([]byte(args[0]), &deposit)
	date := getTimeNow()

	// Check that Token is registered on Blockchain //
	_, err1 := t.checkTokenListed( stub, deposit.Token )
	if err1 != nil { return shim.Error( err1.Error()) }

	// Retrieve the MultiWallet of user. Invoke CoinBalance Smart Contract // 
	user_wallet, err3 := t.retrieveUserWallet( stub, deposit.PublicId )
	if err3 != nil { return shim.Error( err3.Error()) }
	balance_loan := user_wallet.Balances[deposit.Token]

	// Transfer funds from user to Lending Pools //
	transactions := []Transfer{}
	for token, amount := range(deposit.Collaterals) {
		// Balance and Pool of token to update //
		balance_token := user_wallet.Balances[token]
		pool_token, err4 := t.retrieveTokenPool( stub, token )
		if err4 != nil { return shim.Error( err4.Error()) }
		//Check if collateral is in list and add collateral //
		col, isInList := balance_loan.Collaterals[token]
		if !isInList { balance_loan.Collaterals[token] = 0. }
		balance_loan.Collaterals[token] = col + amount
		// Check user token balance is enough //
		if amount > balance_token.Amount {
			return shim.Error( "ERROR: THE USER DOES NOT HOLD THE AMOUNT OF " +
			                   "COLLATERAL FOR TOKEN " + token ) }
		balance_token.Amount = balance_token.Amount - amount 
		balance_token.Collateral_Amount = balance_token.Collateral_Amount + amount
		pool_token.Reserve = pool_token.Reserve + amount
		pool_token.Collateral = pool_token.Collateral + amount
		col_transfer := Transfer{
			Type: "traditional_collateral_deposit", Token: token, 
			From: deposit.PublicId, To: "Lending Pool " + token, Amount: amount,
			Id: xid.New().String(), Date: date }
		transactions = append( transactions, col_transfer )
		// Update balance and pool //
		user_wallet.Balances[token] = balance_token
		_, err5 := t.updatePool( stub, pool_token )
		if err5 != nil { return shim.Error( err5.Error()) }
	}

	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	user_wallet.Balances[deposit.Token] = balance_loan
	user_wallet.Transaction = transactions
	err6 := t.updateUserWallet( stub, user_wallet )
	if err6 != nil { return shim.Error( err6.Error()) }

	// Generate response for outputs to update //
	update_wallets := make( map[string]MultiWallet )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	output := Output{ UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
withdrawCollateral:  this function is called when a borrower wants to withdraw collateral.
              		 Input: array containing a json with fields:
PublicId        string    				// Id of the user depositing collateral
Token           string   			    // Symbol of the borrowing token
Collaterals     map[string]float64      // Amount of collateral deposited for each symbol
RateChange      map[string]float64      // Rates of changes of all currencies with BC
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) withdrawCollateral( stub shim.ChaincodeStubInterface,
															  args []string ) pb.Response {

	// Retrieve arguments of the input in a Wirdrawal Token Model struct //
	if len(args) != 1 {
		return shim.Error( "ERROR: WITHDRAWCOLLATERAL FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENT." ) }
	withdrawal := CollateralWithdrawal{}
	json.Unmarshal([]byte(args[0]), &withdrawal)
	date := getTimeNow()

	// Check that Token is registered on Blockchain //
	_, err1 := t.checkTokenListed( stub, withdrawal.Token )
	if err1 != nil { return shim.Error( err1.Error()) }

	// Retrieve staking token pool //
	pool_loan, err2 := t.retrieveTokenPool( stub, withdrawal.Token )
	if err2 != nil { return shim.Error( err2.Error()) }

	// Retrieve the MultiWallet of user. Invoke CoinBalance Smart Contract // 
	user_wallet, err3 := t.retrieveUserWallet( stub, withdrawal.PublicId )
	if err3 != nil { return shim.Error( err3.Error()) }
	balance_loan := user_wallet.Balances[withdrawal.Token]

	// Compute collateral level and transfer funds from user to Lending Pools //
	transactions := []Transfer{}
	collateral_BC := 0.
	for token, amount := range(withdrawal.Collaterals) {
		// Balance and Pool of token to update //
		balance_token := user_wallet.Balances[token]
		pool_token, err4 := t.retrieveTokenPool( stub, token )
		if err4 != nil { return shim.Error( err4.Error()) }
		//Check that user has the amount to withdraw  //
		collateral, isInList := balance_loan.Collaterals[token]
		if collateral < amount || isInList == false {
			return shim.Error( "ERROR: THE USER DOES NOT HAVE THE AMOUNT OF COLLATERAL " +
		                   	   "FOR TOKEN " + token + " TO WITHDRAW." ) }
		// Transfer funds to user //
		if pool_token.Reserve < amount {
			return shim.Error( "ERROR: THE FUNDS ON THE TOKEN LENDING POOL ARE NOT " +
		                   	   "SUFFICIENT TO WITHDRAW TOKEN " + token ) }
		balance_loan.Collaterals[token] = collateral - amount
		balance_token.Amount = balance_token.Amount + amount 
		pool_token.Reserve = pool_token.Reserve - amount
		pool_token.Collateral = pool_token.Collateral - amount
		col_transfer := Transfer{
			Type: "traditional_collateral_withdraw", Token: token, 
			From:  "Lending Pool " + token, To: withdrawal.PublicId, Amount: amount,
			Id: xid.New().String(), Date: date }
		transactions = append( transactions, col_transfer )
		// Update collateral equivalent in BC and user balance of token //
		collateral_BC = collateral_BC +
						balance_loan.Collaterals[token] * withdrawal.RateChange[token]
		user_wallet.Balances[token] = balance_token
		_, err5 := t.updatePool( stub, pool_token )
		if err5 != nil { return shim.Error( err5.Error()) }
	}

	// Check if the withdrawal is allowed by the risk parameters //
	principal_BC := balance_loan.Borrowing_Amount * 
					withdrawal.RateChange[ withdrawal.Token ]
	CCR := collateral_BC / principal_BC
	risk_params := pool_loan.RiskParameters
	if CCR < risk_params.Col_withdrawal {
		return shim.Error( "ERROR: CCR BELOW THE WITHDRAWAL THRESHOLD. " +
						   "USER CANNOT WITHDRAWAL THE REQUESTED COLLATERAL." ) }

	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	user_wallet.Balances[withdrawal.Token] = balance_loan
	user_wallet.Transaction = transactions
	err7 := t.updateUserWallet( stub, user_wallet )
	if err7 != nil { return shim.Error( err7.Error()) }

	// Generate response for outputs to update //
	update_wallets := make( map[string]MultiWallet )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	output := Output{ UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}

/* -------------------------------------------------------------------------------------------------
repayFunds:  this function is called when a borrower wants to repay some funds of the loan.
              	Input: array containing a json with fields:
PublicId        string    		 // Id of the user repaying fimds
Token           string   		 // Symbol of the borrowing token
Amount          float64     	 // Amount of loan to repay
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) repayFunds( stub shim.ChaincodeStubInterface,
													  args []string ) pb.Response {
	// Retrieve arguments of the input in a Repay Model struct //
	if len(args) != 1 {
		return shim.Error( "ERROR: REPAYFUNDS FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENT." ) }
	repay := RepayFunds{}
	json.Unmarshal([]byte(args[0]), &repay)
	date := getTimeNow()

	// Check that Token is registered on Blockchain //
	_, err1 := t.checkTokenListed( stub, repay.Token )
	if err1 != nil { return shim.Error( err1.Error()) }

	// Retrieve staking token pool //
	pool_loan, err2 := t.retrieveTokenPool( stub, repay.Token )
	if err2 != nil { return shim.Error( err2.Error()) }

	// Retrieve the MultiWallet of user. Invoke CoinBalance Smart Contract // 
	user_wallet, err3 := t.retrieveUserWallet( stub, repay.PublicId )
	if err3 != nil { return shim.Error( err3.Error()) }
	balance_loan := user_wallet.Balances[repay.Token]

	// Ensure that the amount to repay is less than the borrowed amount //
	transactions := []Transfer{}
	repayment := math.Min( balance_loan.Borrowing_Amount, repay.Amount )
	balance_loan.Amount = balance_loan.Amount - repayment
	balance_loan.Borrowing_Amount = balance_loan.Borrowing_Amount - repayment
	pool_loan.Reserve = pool_loan.Reserve + repayment
	pool_loan.Loaned = math.Max( pool_loan.Loaned - repayment, 0. )
	repay_transfer := Transfer{
		Type: "repay_traditional", Token: repay.Token, From: repay.PublicId, 
		To: "Lending Pool " + repay.Token, Amount: repayment, 
	    Id: xid.New().String(), Date: date }
	transactions = append( transactions, repay_transfer )

	// Update token Pool // 
	_, err4 := t.updatePool( stub, pool_loan )
	if err4 != nil { return shim.Error( err4.Error()) }

	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	user_wallet.Balances[repay.Token] = balance_loan
	user_wallet.Transaction = transactions
	err5 := t.updateUserWallet( stub, user_wallet )
	if err5 != nil { return shim.Error( err5.Error()) }

	// Generate response for outputs to update //
	update_wallets := make( map[string]MultiWallet )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	output := Output{ UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}
	
/* -------------------------------------------------------------------------------------------------
checkLiquidation: this function is called a loan is identified to be undercollaterised and that may
                  be consequently closed. Input: array containing a json with fields:
PublicId        string    		 		// Id of the user to revise collateral
Token           string   				// Symbol of the borrowing token
RateChange      map[string]float64      // Rates of changes of all currencies with BC
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) checkLiquidation( stub shim.ChaincodeStubInterface,
													  		args []string ) pb.Response {
	// Retrieve arguments of the input in a Liquidator Model struct //
	if len(args) != 1 {
		return shim.Error( "ERROR: CHECKLIQUIDATION FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENT." )
	}
	liquidator := Liquidator{}
	json.Unmarshal([]byte(args[0]), &liquidator)
	date := getTimeNow()

	// Check that Token is registered on Blockchain //
	_, err1 := t.checkTokenListed( stub, liquidator.Token )
	if err1 != nil { return shim.Error( err1.Error()) }

	// Retrieve token pool //
	pool_loan, err2 := t.retrieveTokenPool( stub, liquidator.Token )
	if err2 != nil { return shim.Error( err2.Error()) }

	// Retrieve the MultiWallet of user. Invoke CoinBalance Smart Contract // 
	user_wallet, err3 := t.retrieveUserWallet( stub, liquidator.PublicId )
	if err3 != nil { return shim.Error( err3.Error()) }
	balance_loan := user_wallet.Balances[liquidator.Token]

	// Compute collateral equivalent in Base Coin (BC) //
	collateral_BC := 0.
	for token, amount := range(balance_loan.Collaterals) {
		collateral_BC = collateral_BC + 
		                amount*liquidator.RateChange[token] }
	principal_BC := balance_loan.Borrowing_Amount *
					liquidator.RateChange[liquidator.Token]
					
	// Compute CCR level and check if the loan should be liquidated //
	risk_params := pool_loan.RiskParameters
	CCR := collateral_BC / principal_BC
	if CCR >= risk_params.Col_liquidation {
		output := Output{ Liquidated: "NO",
						  UpdateWallets: make(map[string]MultiWallet) }
		outputBytes, _ := json.Marshal( output )
		return shim.Success( outputBytes ) }

	// Liquidate Loan. Return outstanding collateral (if any) //
	transactions := []Transfer{}
	pct_discount := math.Min( 1./ CCR, 1. )
	for token, amount := range(balance_loan.Collaterals) {
		pool_token, err4 := t.retrieveTokenPool( stub, token )
		if err4 != nil { return shim.Error( err4.Error()) }
		balance_token := user_wallet.Balances[token] 
		amount_returned := amount*(1 -pct_discount)
		// Check that there are enough funds in Token Pool //
		if amount_returned >  pool_token.Reserve {
			return shim.Error( "ERROR: INSUFFICIENT FUNDS TO RETURN " +
							   "COLLATERAL AMOUNT FOR TOKEN " + token ) } 
		balance_token.Amount = balance_token.Amount + amount_returned
		balance_token.Collateral_Amount = 0.
		pool_token.Reserve = pool_token.Reserve - amount_returned
		pool_token.Collateral = math.Max( pool_token.Collateral - amount, 0.)
		// Transaction of outstanding collateral to User //
		liquidation_transfer := Transfer{
			Type: "traditional_liquidated", Token: token, From: "Lending Pool " + token, 
			To: liquidator.PublicId, Amount: amount_returned, 
			Id: xid.New().String(), Date: date  }
		transactions = append( transactions, liquidation_transfer)
		_, err5 := t.updatePool( stub, pool_token )
		if err5 != nil { return shim.Error( err5.Error()) }
	}

	// Set lending amount to 0 and discount the Lending from the Token Pool //
	pool_loan.Loaned = math.Max( pool_loan.Loaned - balance_loan.Amount, 0. )
	balance_loan.Borrowing_Amount = 0.
	balance_loan.Collaterals = make( map[string]float64 )

	// Update token Pool // 
	_, err6 := t.updatePool( stub, pool_loan )
	if err6 != nil { return shim.Error( err6.Error()) }

	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	user_wallet.Balances[liquidator.Token] = balance_loan
	user_wallet.Transaction = transactions
	err7 := t.updateUserWallet( stub, user_wallet )
	if err7 != nil { return shim.Error( err7.Error()) }

	// Generate response for outputs to update //
	update_wallets := make( map[string]MultiWallet )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	output := Output{ Liquidated: "YES",
					  UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes ) 													  
}


/* -------------------------------------------------------------------------------------------------
updateRiskParameters:  	this function is called when risk parameters need to be modified. It can 
						only be called by Admin. This takes an array containing a json with the 
						following fields:
Token           string   // Symbol of a new token registered in the smart contract
RiskParameters  Risk     // Json containing the risk parameters of the token
------------------------------------------------------------------------------------------------- */
func (t *TraditionalLendingSmartContract) updateRiskParameters( stub shim.ChaincodeStubInterface,
																args []string ) pb.Response {
	// Retrieve arguments of the input in a Risk Model struct //
	if len(args) != 2 {
		return shim.Error( "ERROR: UPDATE RISK PARAMETERS TOKEN FUNCTION SHOULD BE CALLED " +
							"WITH TWO ARGUMENTS." ) }
	update := Risks{}
	json.Unmarshal([]byte(args[1]), &update)

	// Check that Token is registered on Blockchain //
	_, err1 := t.checkTokenListed( stub, args[0] )
	if err1 != nil { return shim.Error( err1.Error()) }

	// Retrieve token pool //
	pool_loan, err2 := t.retrieveTokenPool( stub, args[0] )
	if err2 != nil { return shim.Error( err2.Error()) }
	
	// Update token Pool // 
	pool_loan.RiskParameters = update
	_, err3 := t.updatePool( stub, pool_loan )
	if err3 != nil { return shim.Error( err3.Error()) }

	return shim.Success( nil )
}

/* -------------------------------------------------------------------------------------------------
updatePoolBytes: this function updates a Pool on Blockchain. Inputs: the updated pool in Bytes.
------------------------------------------------------------------------------------------------- */

func (t *TraditionalLendingSmartContract) updatePoolBytes( stub shim.ChaincodeStubInterface, 
														   args []string) pb.Response {
	// Store pool on Blockchain //
	pool := LendingPool{}
	poolBytes := []byte(args[0])
	json.Unmarshal( poolBytes, &pool )
	err := stub.PutState( IndexLendingPools + pool.Token, poolBytes )
	if err != nil {
		return shim.Error( "ERROR: UPDATING THE STATE OF POOL FOR TOKEN " +
						    pool.Token + ". " + err.Error() ) }
	return shim.Success( nil )
}


/* -------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------- */

func main() {
	err := shim.Start(&TraditionalLendingSmartContract{})
	if err != nil {
		fmt.Errorf("Error starting Token chaincode: %s", err)
	}
}

