
package main
import (
	"time" 
	"encoding/json"
	"fmt"
	"bytes"
	"math"
	"strconv"
	"github.com/rs/xid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/common/util"
	pb "github.com/hyperledger/fabric/protos/peer"
)



/* -------------------------------------------------------------------------------------------------
Init:  this function is called at PRIVI Blockchain Deployment and initialises the Coin Balance 
	   Smart Contract. This smart contract is the responsible to manage the balances of the
	   different tokens powered in PRIVI Ecosystem. The initialisation sets the Private Key ID
	   of the admin in the ledger. Args: array containing a string:
PrivateKeyID           string   // Private Key of the admin of the smart contract
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) Init( stub shim.ChaincodeStubInterface ) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	if args[1] == "UPGRADE" { return shim.Success(nil) }

	// Store in the state of the smart contract the Private Key of Admin //
	err := stub.PutState( IndexAdmin, []byte( args[0]) )
	if err != nil {
		return shim.Error( "ERROR: SETTING THE ADMIN PRIVATE KEY: " +
	                       err.Error() )
	}

	// Initialise list of tokens in the Smart Contract as empty //
	token_list, _ := json.Marshal( []string{} )
	err2 := stub.PutState( IndexTokenList, token_list )
	if err2 != nil {
		return shim.Error( "ERROR: INITIALISING THE TOKEN LIST: " +
	                       err2.Error() )
	}

	// Initialise list of users in the Smart Contract as empty //
	user_list, _ := json.Marshal( []string{} )
	err3 := stub.PutState( IndexUserList, user_list )
	if err3 != nil {
		return shim.Error( "ERROR: INITIALISING THE USER LIST: " +
	                       err3.Error() )
	}

	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
Invoke:  this function is the router of the different functions supported by the Smart Contract.
		 It receives the input from the controllers and ensure the correct calling of the functions.
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	
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
		case "getUserList":
			return t.getUserList(stub, args)
		case "registerWallet":
	     	return t.registerWallet(stub, args)
		case "balanceOf":
			return t.balanceOf(stub, args)
		case "transfer":
			return t.transfer(stub, args)
		case "multitransfer":
			return t.multitransfer(stub, args)
		case "spendFunds":
			return t.spendFunds(stub, args)
		case "getWallets":
			return t.getWallets(stub, args)
		case "updateMultiwallets":
			return t.updateMultiwallets(stub, args)
		case "getHistory":
			return t.getHistory(stub, args)
		case "mint":
			return t.mint(stub, args)
		case "burn":
			return t.burn(stub, args)
		case "swap":
			return t.swap(stub, args)
		case "withdraw":
			return t.withdraw(stub, args)
	}

	// If function does not exist, retrieve error//
	return shim.Error( "ERROR: INCORRECT FUNCTION NAME " +
	                   function)
}


/* -------------------------------------------------------------------------------------------------
registerToken:  this function is called when a new token is listed on the PRIVI Blockchain. It can 
				only be called by Admin. At deployment, PRIVI Coin and Base Coin are created in the 
				system. Input: array containing a json with fields:
Name           string   // Private Key of the admin of the smart contract
Symbol         string   // Public Key of the admin of the smart contract
Supply         string   // Amount of coins sold in the IEO sale
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) registerToken( stub shim.ChaincodeStubInterface,
												  args []string ) pb.Response {

	// Retrieve arguments of the input in a Token Model struct //
	if len(args) != 2 {
		return shim.Error( "ERROR: REGISTER TOKEN FUNCTION SHOULD BE CALLED " +
		                   "WITH 2 ARGUMENTS." ) }
	token := Token{}
	json.Unmarshal([]byte(args[0]), &token)
	admin_public_id := args[1] 
	
	update_wallets := make( map[string]MultiWallet )
	// Check if Token is already registered on Blockchain //
	token_check, err1 := stub.GetState( IndexToken+token.Symbol )
	if err1 != nil {
		return shim.Error( "ERROR: CHECKING IF TOKEN IS ALREADY " +
		                   "REGISTERED. " +  err1.Error() )
	} else if token_check != nil {
		return shim.Error( "ERROR: TOKEN " +  token.Symbol +
						   " IS ALREADY LISTED." ) }
	          
	// List new token on Blockchain //
	tokenBytes, _ := json.Marshal(token)
	err2 := stub.PutState(IndexToken+token.Symbol, tokenBytes)
	if err2 != nil {
		return shim.Error( "ERROR: SETTING TOKEN STATE TO " +
		                   "BLOCKCHAIN. " + err2.Error() ) }

	// Add token to the token list of the smart contract //
	tokenListBytes, err3 := stub.GetState( IndexTokenList )
	tokenList := []string{}
	json.Unmarshal(tokenListBytes, &tokenList)
	tokenList = append(tokenList, token.Symbol)
	tokenListBytes2, _ := json.Marshal(tokenList)
	err4 := stub.PutState( IndexTokenList, tokenListBytes2 )
	if (err3 != nil || err4 != nil) {
		return shim.Error( "ERROR: ADDING TOKEN " + token.Symbol + 
		                   " THE TOKEN LIST: " ) }

	// Retrieve list of all users registered in Blockchain //
	userListBytes, err6 := stub.GetState( IndexUserList )
	if err6 != nil {
		return shim.Error( "ERROR: RETRIEVING THE USER LIST" ) }
	user_list := []string{}
	json.Unmarshal( userListBytes, &user_list )

	// Initialise values to run the iteration throughout all users//
	var userWalletBytes []byte
	var user_wallet MultiWallet
	collaterals := make( map[string]float64 )
	user_balance := Balance{ 
		Token: token.Symbol, Amount: 0., Credit_Amount: 0., 
		Staking_Amount: 0., Borrowing_Amount: 0., Type: "CRYPTO",
		PRIVI_lending: 0., PRIVI_borrowing: 0.,
		Collateral_Amount: 0., Collaterals: collaterals }
    admin_balance := Balance{ 
		Token: token.Symbol, Amount: token.Supply, Credit_Amount: 0., 
		Staking_Amount: 0., Borrowing_Amount: 0., Type: "CRYPTO",
		PRIVI_lending: 0., PRIVI_borrowing: 0., 
		Collateral_Amount: 0., Collaterals: collaterals }

	// Loop through all users registered on Blockchain //
	for i, _ := range user_list {
		id := user_list[i]
		userBytes, err_loop := stub.GetState( IndexWallet + id )
		if err_loop != nil {
			return shim.Error( "ERROR: RETRIEVING THE STATE OF " +
							   " USER: " + id + " " + err_loop.Error() ) }
		json.Unmarshal(userBytes, &user_wallet)
		// Check if balance of token is already registered //
		balances := user_wallet.Balances
		_, registered := balances[ token.Symbol ]
		if !registered {
			balances[ token.Symbol ] = user_balance
			// If id is admin id, set all the supply //
			if id == admin_public_id { 
				balances[ token.Symbol ] = admin_balance }
			user_wallet.Balances = balances
			userWalletBytes, _ = json.Marshal(user_wallet)
			update_wallets[ id ] = user_wallet
			err_loop = stub.PutState( IndexWallet+id, userWalletBytes )
			if err_loop != nil {
				return shim.Error( "ERROR SETTING TOKEN " + token.Symbol +
								   " BALANCE FOR USER: " + id + " " + 
								   err_loop.Error() ) }
		}
	}

	// Prepare output object with updates //
	output := Output{ UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
getTokenList:  this function returns the list of tokens listed in the PRIVI Blockchain. This 
               function does not need to be called with any argument.
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) getTokenList( stub shim.ChaincodeStubInterface,
												 args []string ) pb.Response {
	// Retrieve list of tokens from the Blockchain //
	tokenList, err := stub.GetState( IndexTokenList )
	if err != nil {
		return shim.Error( "ERROR: RETRIEVING THE TOKEN LIST" ) }

	return shim.Success( tokenList )
}

/* -------------------------------------------------------------------------------------------------
getUserList:  this function returns the list of users registered in the PRIVI Blockchain. This 
              function does not need to be called with any argument.
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) getUserList( stub shim.ChaincodeStubInterface,
												 args []string ) pb.Response {
	// Retrieve list of users from the Blockchain //
	userList, err := stub.GetState( IndexUserList )
	if err != nil {
		return shim.Error( "ERROR: RETRIEVING THE USER LIST" ) }
	return shim.Success( userList )
}


/* -------------------------------------------------------------------------------------------------
registerBalance:  this function register a wallet for a new user in the Smart Contract. It first 
                  checks that the user does not already hold a wallet. Args:
args[0]              string   // Id of the user 
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) registerWallet( stub shim.ChaincodeStubInterface, 
	                                                args []string ) pb.Response {

	// Retrieve Public Id Key and check that this is not already registered //
	PublicId := args[0]
	wallet, err1 := stub.GetState( IndexWallet + PublicId )
	if err1 != nil {
		return shim.Error( "ERROR: CHECKING IF WALLET IS ALREADY " +
		                   "REGISTERED. " +  err1.Error() )
	} else if wallet != nil {
		return shim.Error( "ERROR: WALLET FOR USER " +  PublicId +
						   " IS ALREADY REGISTERED." ) }
						   
	update_wallets := make( map[string]MultiWallet )

	// Retrieve list of tokens from the Blockchain //
	tokenListBytes, err := stub.GetState( IndexTokenList )
	if err != nil {
		return shim.Error( "ERROR: RETRIEVING THE TOKEN LIST" ) 
	}
	token_list := []string{}
	json.Unmarshal( tokenListBytes, &token_list )

	// Create a Wallet with all tokens initialised to 0 //
	balances := make( map[string]Balance )
	collaterals := make( map[string]float64 )
	PRIVI := make( map[string]bool )
	for i, _ := range token_list {
		balances[ token_list[i] ] = Balance{
			Token: token_list[i], Amount: 0., Credit_Amount: 0.,
			Staking_Amount: 0.,  Borrowing_Amount: 0., Type: "CRYPTO",
			PRIVI_lending: 0., PRIVI_borrowing: 0.,
			PRIVIcreditLend: PRIVI, PRIVIcreditBorrow: PRIVI,
			Collateral_Amount: 0., Collaterals: collaterals }
	}
	newWallet := MultiWallet{}
	newWallet.PublicId = PublicId
	newWallet.Balances = balances
	newWallet.TrustScore = 0.5
	newWallet.EndorsementScore = 0.5
	newWallet.BalancesFT = make( map[string]BalanceFT )
	newWallet.BalancesNFT = make( map[string]BalanceNFT )

	// Store state of the new balance on Blockchain //
	newWalletBytes, _ := json.Marshal(newWallet)
	err2 := stub.PutState(IndexWallet + PublicId, newWalletBytes)
	if err2 != nil {
		return shim.Error( "ERROR: STORING STATE THE WALLET " +
		                   "ON BLOCKCHAIN. " + err2.Error() ) } 

	// Add new user to the state of the Smart Contract //
	userListBytes, err3 := stub.GetState( IndexUserList )
	userList := []string{}
	json.Unmarshal(userListBytes, &userList)
	userList = append(userList, PublicId)
	userListBytes2, _ := json.Marshal(userList)
	err4 := stub.PutState( IndexUserList, userListBytes2 )
	if (err3 != nil || err4 != nil) {
		return shim.Error( "ERROR: ADDING USER " + PublicId + 
						   " TO THE USER LIST" ) }
				
	// Prepare outputs with the updates to make //
	update_wallets[ newWallet.PublicId ] = newWallet
	output := Output{ UpdateWallets: update_wallets }

	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
balanceOf: This function is called to get the balance on the wallet of a given actor. The caller of 
		   this function can only be Cache Admin or wallet owner. 
		   Args: array containing
args[0]               string   // Public Id Key of the actor
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) balanceOf( stub shim.ChaincodeStubInterface, 
	                                          args []string ) pb.Response {

	// Retrieve information from the input //
	if len(args) != 1 {
		return shim.Error( "ERROR: BALANCEOF FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENT." ) }
	publicId := args[0]

	// Retrieve balance information //
	balance, err := stub.GetState( IndexWallet + publicId )
	if err != nil {
		return shim.Error( "ERROR: FAILED TO GET THE BALANCE " + 
						   "INFORMATION." + err.Error() )
	} else if balance == nil {
		return shim.Error(  "ERROR: THE BALANCE FOR " +  publicId +
							" DOES NOT EXIST." ) } 
	return shim.Success( balance )
}


/* -------------------------------------------------------------------------------------------------
transfer: This function is called to transfer a given token from one wallet to another one.
Token              string   // Symbol of token to transfer	  
From               string   // Id of the sender
To                 string   // Id of the receiver
Amount             float64  // Amount that is being sent
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) transfer( stub shim.ChaincodeStubInterface, 
	                                        args []string ) pb.Response {

	// Retrieve information from the input in a Transfer object //
	if len(args) != 1 {
		return shim.Error( "ERROR: TRANSFER FUNCTION SHOULD BE CALLED " +
		                   "WITH ONE ARGUMENT." ) }
	transfer := Transfer{}
	input := []byte(args[0])
	json.Unmarshal(input, &transfer)
	date := getTimeNow()

	update_wallets := make( map[string]MultiWallet )

	// Verify the sender and receiver are not the same //
	if ( transfer.From==transfer.To ) {
		return shim.Error( "ERROR: SENDER AND RECEIVER CANNOT BE THE " +
		                   " SAME IN A TRANSFER." )
	}

	// Retrieve balance of the sender //
	senderBytes, err1 := stub.GetState( IndexWallet + transfer.From )
	if err1 != nil {
		return shim.Error( "ERROR: FAILED TO GET SENDER OF THE TRANSFER. " + 
		                   err1.Error() ) } 
	sender := MultiWallet{}
	err2 := json.Unmarshal(senderBytes, &sender)
	if err2 != nil {
		return shim.Error( "ERROR: RETRIEVING THE SENDER WALLET. " +
		                    err2.Error() )
	}

	// Retrieve balancer for receiver //
	receiverBytes, err3 := stub.GetState( IndexWallet + transfer.To )
	if err3 != nil {
		return shim.Error( "ERROR: FAILED TO GET RECEIVER OF THE TRANSFER. " + 
		                   err3.Error() ) }
	receiver := MultiWallet{}
	err4 := json.Unmarshal(receiverBytes, &receiver)
	if err4 !=nil {
		return shim.Error( "ERROR: RETRIEVING THE RECEIVER WALLET. " +
		                    err4.Error() )
	}

	// Check that sender has sufficient funds and receiver overflows //
	balance_sender := sender.Balances[ transfer.Token ]
	balance_receiver := receiver.Balances[ transfer.Token ]
	allow_transfer := balance_sender.Amount 
	if ( transfer.Amount >  allow_transfer ) {
		return shim.Error( "ERROR: INSUFFICIENT FUNDS IN SENDER WALLET." )
	}
	if ( balance_receiver.Amount + transfer.Amount < balance_receiver.Amount ) {
		return shim.Error( "ERROR: OVERFLOW ON RECEIVER WALLET." )
	}

	// Transfer funds from sender wallet to receiver wallet //
	balance_sender.Amount = balance_sender.Amount - transfer.Amount
	balance_receiver.Amount = balance_receiver.Amount + transfer.Amount
	sender.Balances[transfer.Token] = balance_sender
	receiver.Balances[transfer.Token] = balance_receiver
	transfer.Id = xid.New().String()
	transfer.Date = date
	transfer.Type = "transfer_send"
	sender.Transaction = []Transfer{ transfer }
	transfer.Type = "transfer_receive"
	receiver.Transaction = []Transfer{ transfer }

	// Update state of wallet on Blockchain //
	senderAsBytes, _ := json.Marshal(sender)
	receiverAsBytes, _ := json.Marshal(receiver)
	err6 := stub.PutState(IndexWallet + transfer.From, senderAsBytes)
	err7 := stub.PutState(IndexWallet + transfer.To, receiverAsBytes)
	if ( err6 != nil || err7 != nil) {
		return shim.Error( "ERROR: STORING STATE OF BALANCE " +
						   "ON BLOCKCHAIN" ) }
						   
	// Return output with updates //
	update_wallets[ transfer.From ] = sender
	update_wallets[ transfer.To ] = receiver
	output := Output{ UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal(output)
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
multitransfer: This function is called to perform a multitransfer between different actors in 
               one call to blockchain. Args: is a list of transfer types.
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) multitransfer( stub shim.ChaincodeStubInterface, 
											      args []string ) pb.Response {
	
	// Set of users participating in transfer and whose state should be updated //
	transaction_users := make(map[string]MultiWallet)
	date := getTimeNow()

	// Iterate throught all the transactions on the multitransfer call //
	for _, arg := range args {
		// Retrieve trasnfer from the list // 
		transfer := Transfer{}
		json.Unmarshal( []byte(arg), &transfer )
		From := transfer.From
		To := transfer.To
		Amount := transfer.Amount
		Token := transfer.Token
		if ( From == To || Amount==0 ) { continue }
		
		// Check if sender is already in transaction users list. Otherwise add to set //
		_, inList1 := transaction_users[From]
		if !inList1 {
			sender := MultiWallet{}
			senderBytes, err1 := stub.GetState( IndexWallet + From )
			json.Unmarshal(senderBytes, &sender)
			if err1 != nil {
				return shim.Error( "ERROR: GETTING THE INFORMATION " +
				                   "FROM SENDER " + From + ". " + err1.Error() )} 
			sender.Transaction = []Transfer{}
			transaction_users[From] = sender
		}

		// Check if receiver is already in transaction users list. Otherwise add to set //
	    _, inList2 := transaction_users[To]
		if !inList2 {
			receiver := MultiWallet{}
			receiverBytes, err2 := stub.GetState( IndexWallet + To )
			json.Unmarshal(receiverBytes, &receiver)
			if err2 != nil {
				return shim.Error( "ERROR: GETTING THE INFORMATION " +
				                   "FROM SENDER " + From + ". " + err2.Error() ) } 
			receiver.Transaction = []Transfer{}
			transaction_users[To] = receiver
		}

		sender := transaction_users[From]
		receiver := transaction_users[To]

		// Check that the sender holds the amount to send and transfer funds //
		sender_balance := sender.Balances[Token]
		receiver_balance := receiver.Balances[Token]
		if ( Amount > sender_balance.Amount ){
			return shim.Error( "ERROR: INSUFFICIENT FUNDS IN SENDER " +
		                       "WALLET " + From )} 
		sender_balance.Amount = sender_balance.Amount - Amount
		receiver_balance.Amount = receiver_balance.Amount + Amount
		sender.Balances[Token] = sender_balance
		receiver.Balances[Token] = receiver_balance
		transfer.Date = date
		transfer.Id = xid.New().String()
		transfer.Date = date
		transfer.Type = "transfer_send"
		sender.Transaction = append( sender.Transaction, transfer )
		transfer.Type = "transfer_receive"
		receiver.Transaction = append( receiver.Transaction, transfer )
		
		transaction_users[To] = receiver
		transaction_users[From] = sender
	}

	// Update States of all the users that did some transaction //
	for ID, transactions := range(transaction_users) {
		user_update, _ := json.Marshal( transactions )
		err3 := stub.PutState( IndexWallet + ID, user_update )
		if err3 != nil {
			return shim.Error( "ERROR: UPDATING STATE OF USER " + ID + 
		                       ". " + err3.Error() ) }
	}

	// Return output with updates //
	output := Output{ UpdateWallets: transaction_users }
	outputBytes, _ := json.Marshal(output)
	return shim.Success( outputBytes )
}

/* -------------------------------------------------------------------------------------------------
getWallets:  this function returns the list of users with its respective Multiwallet object. It
             does not need to be called with any argument.
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) getWallets( stub shim.ChaincodeStubInterface,
											   args []string ) pb.Response {

	// Retrieve list of users from the Blockchain //
	userListBytes, err1 := stub.GetState( IndexUserList )
	if err1 != nil {
		return shim.Error( "ERROR: RETRIEVING THE USER LIST. " + 
						    err1.Error() ) 
	}
	user_list := []string{}
	json.Unmarshal( userListBytes, &user_list )

	userWallets := []MultiWallet{}
	// Loop through all users and get Wallets //
	for _, id := range user_list {
		userBytes, err2 := stub.GetState(  IndexWallet + id  )
		if err2 != nil {
			return shim.Error( "ERROR: RETRIEVING THE WALLET OF " +
							   " USER: " + id + " " + err2.Error() ) }
		user_wallet := MultiWallet{}
		json.Unmarshal(userBytes, &user_wallet)
		user_wallet.Transaction = []Transfer{}
		userWallets = append( userWallets, user_wallet )
	}

	result, _ := json.Marshal( userWallets )
	return shim.Success( result )
}

/* -------------------------------------------------------------------------------------------------
updateMultiwallets: This function is called by the TraditionalLending smart contract with the 
				   	multiwallet updates after requesting for a loan. It is called with an array 
				    of MultiWallets.
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) updateMultiwallets( stub shim.ChaincodeStubInterface,
													   args []string ) pb.Response {
	// Loop through all the users to update mutiwallet //
	for _, arg := range args {
		walletBytes := []byte(arg)
		wallet := MultiWallet{}
		json.Unmarshal(walletBytes, &wallet)
		// Update state of the MultiWallet for user // 
		err := stub.PutState( IndexWallet + wallet.PublicId, walletBytes )
		if err != nil {
			return shim.Error( "ERROR: UPDATING MULTIWALLETS FOR USER. " + wallet.PublicId +
								err.Error() )
		}
	}
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
spendFunds: This function is called when a user wants to spend funds with some tokens. This function
			should be called with the following arguments:
PublicId           string      // ID of the fund spender
ProviderId         string      // ID of the provider 	
Token              string      // Symbol of token to transfer	
Amount             float64     // Amount desired to spend
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) spendFunds( stub shim.ChaincodeStubInterface,
										       args []string ) pb.Response {
	//// Retrieve information from the input //
	if len(args) != 1 {
		return shim.Error( "ERROR: GETHISTORY FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENT." ) }
	spending := Spending{}
	json.Unmarshal([]byte(args[0]), &spending)

	// Retrieve mutiwallet of spender //
	walletBytes, err1 := stub.GetState( IndexWallet + spending.PublicId )
	if err1 != nil {
		return shim.Error( "ERROR: RETRIEVING MULTIWALLET FOR USER: " + spending.PublicId +
						   ". " + err1.Error() ) }
	spender_wallet := MultiWallet{}
	json.Unmarshal(walletBytes, &spender_wallet)
	spender_balance := spender_wallet.Balances[spending.Token]

	// Retrieve mutiwallet of provider //
	walletBytes2, err2 := stub.GetState( IndexWallet + spending.ProviderId )
	if err2 != nil {
		return shim.Error( "ERROR: RETRIEVING MULTIWALLET FOR USER: " + spending.ProviderId +
						   ". " + err2.Error() ) }
	provider_wallet := MultiWallet{}
	json.Unmarshal(walletBytes2, &provider_wallet)
	provider_balance := provider_wallet.Balances[spending.Token]

	// Get all PRIVI Loans of user //
	PRIVI_credits := []PRIVIloan{}
	total_credit := 0.
	total_credit_discount := 0.
	for priviId, _ := range(spender_balance.PRIVIcreditBorrow) {
		privi_loan := PRIVIloan{}
		chainCodeArgs := util.ToChaincodeArgs( "getPRIVIcredit", priviId )
		response := stub.InvokeChaincode( PRIVI_CREDIT_CHAINCODE, chainCodeArgs, 
											CHANNEL_NAME )
		if response.Status != shim.OK {
			return shim.Error( "ERROR INVOKING THE PRIVICREDIT CHAINCODE TO " +
								"GET THE PRIVI CREDIT INFO OF: " + priviId ) }
		json.Unmarshal(response.Payload, &privi_loan)
		PRIVI_credits = append( PRIVI_credits, privi_loan )
		borrower_credit := privi_loan.State.Borrowers[ spending.PublicId ]
		total_credit = total_credit + borrower_credit.Amount
		total_credit_discount = total_credit_discount + 
		                        borrower_credit.Amount * (1-privi_loan.P_premium) }

	transactions := []Transfer{}
	PREMIUMS := map[string]Premium{}
	// Check if user needs to user PRIVI Credit Funds // 
	amount_without_credit := math.Max(spender_balance.Amount - total_credit, 0.)
	if spending.Amount > amount_without_credit {
		if amount_without_credit + total_credit_discount < spending.Amount {
			return shim.Error( "ERROR: THE SPENDER " + spending.PublicId + " DOES NOT HAVE " +
							   "ENOUGH FUNDS TO SPEND.") }
		// Charge first the amount without credit //
		transfer := Transfer{
			Type: "spending", Token: spending.Token, From: spending.PublicId, 
			To: spending.ProviderId, Amount: amount_without_credit }
		transactions = append( transactions, transfer )
		// Credit needed to pay remaining //
		credit_needed := spending.Amount - amount_without_credit
		for i, priviLoan := range(PRIVI_credits) {
			if credit_needed <= 0. {break;}
			spender_PRIVI, _ := priviLoan.State.Borrowers[ spending.PublicId ]
			// Check if it is enough with this credit //
			premium_charged := credit_needed * priviLoan.P_premium
			amount_charged := credit_needed + premium_charged
			if spender_PRIVI.Amount < amount_charged {
				premium_charged = spender_PRIVI.Amount * priviLoan.P_premium
				amount_charged = spender_PRIVI.Amount
			}
			credit_needed = credit_needed - (amount_charged-premium_charged)
			transfer_premium := Transfer{
				Type: "PRIVI_premium", Token: spending.Token, From: spending.PublicId, 
				To: "PRIVI Pool " + priviLoan.LoanId, Amount: premium_charged }
			transactions = append( transactions, transfer_premium )
			// Update PRIVI Loan with new information //
			premium_list := priviLoan.State.PremiumList
			premium_id := xid.New().String()
			premium_list[premium_id] = Premium{ ProviderId: spending.ProviderId, 
				PremiumId: premium_id, Risk_Pct: 0., Premium_Amount: premium_charged }
			//new_premium, _ := json.Marshal(premium_list[premium_id])
			PREMIUMS[priviLoan.LoanId] = premium_list[premium_id]
			spender_PRIVI.Amount = spender_PRIVI.Amount - amount_charged
			spender_balance.Amount = spender_balance.Amount - premium_charged
			priviLoan.State.Total_Premium = priviLoan.State.Total_Premium + premium_charged
			priviLoan.State.Total_Coverage = priviLoan.State.Total_Coverage + premium_charged
			priviLoan.State.PremiumList = premium_list
			priviLoan.State.Borrowers[ spending.PublicId ] = spender_PRIVI 
			PRIVI_credits[i] = priviLoan }
	}

	// Transfer funds from spender to provider //
	spender_balance.Amount = spender_balance.Amount - spending.Amount
	provider_balance.Amount = provider_balance.Amount + spending.Amount
	transfer := Transfer{
		Type: "spending", Token: spending.Token, From: spending.PublicId, 
		To: spending.ProviderId, Amount: spending.Amount }
	transactions = append( transactions, transfer )

	// Update state of wallet for both Spender on Blockchain //
	spender_wallet.Balances[spending.Token] = spender_balance
	spender_wallet.Transaction = transactions
	spenderAsBytes, _ := json.Marshal(spender_wallet)
	err3 := stub.PutState(IndexWallet + spending.PublicId, spenderAsBytes)
	if err3 != nil {
		return shim.Error( "ERROR: STORING SPENDER ON BLOCKCHAIN " + err3.Error() ) }
	
	// Update state of wallet for both Provider on Blockchain //
	provider_wallet.Balances[spending.Token] = provider_balance
	provider_wallet.Transaction = transactions
	providerAsBytes, _ := json.Marshal(provider_wallet)
	err4 := stub.PutState(IndexWallet + spending.ProviderId, providerAsBytes)
	if err4 != nil {
		return shim.Error( "ERROR: STORING PROVIDER ON BLOCKCHAIN " + err4.Error() ) }

	// Update PRIVI credits of user by invoking PRIVI Coin Chaincode //
	update_privi_credits := []string{"updatePRIVIcredits"}
	for _, credit := range(PRIVI_credits) {
		creditBytes, _ := json.Marshal( credit )
		update_privi_credits = append( update_privi_credits, string(creditBytes)) }
	multiChainCodeArgs := ToChaincodeArgs( update_privi_credits )
	response2 := stub.InvokeChaincode( PRIVI_CREDIT_CHAINCODE, multiChainCodeArgs, 
										CHANNEL_NAME )
	if response2.Status != shim.OK {
		return shim.Error( "ERROR INVOKING THE UPDATEMULTIWALLET CHAINCODE TO " +
							"UPDATE THE WALLET OF USERs" )
	}
	premiumsBytes, _ := json.Marshal(PREMIUMS)
	return shim.Success(premiumsBytes)
}


/* -------------------------------------------------------------------------------------------------
getHistory: This function returns the transaction history of a given Public Key Id from a given
            timestamp. Args: is an array containing a json with the following attributes:
PublicId           string   // Id of the user to retrieve the history
Timestamp          string   // Timestamp of time from which retrieve the history
------------------------------------------------------------------------------------------------- */

func (t *CoinBalanceSmartContract) getHistory( stub shim.ChaincodeStubInterface,
	                                           args []string ) pb.Response {

	// Retrieve information from the input //
	if len(args) != 1 {
		return shim.Error( "ERROR: GETHISTORY FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENT." )
	}
	history := History{}
	input := []byte(args[0])
	json.Unmarshal(input, &history)


	// Retrieve iterator of History for the User Public Key //
	resultsIterator, err1 := stub.GetHistoryForKey(IndexWallet + history.PublicId)
	if err1 != nil {
		return shim.Error( "ERROR: RETRIEVING THE HISTORY FROM BLOCKCHAIN. " +
						   "ERROR WAS: " + err1.Error())
	}
	defer resultsIterator.Close()

	// Create response. Filter for timestamps greater than the given one //
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err2 := resultsIterator.Next()
		if err2 != nil {
			return shim.Error( "ERROR: GETTING NEXT ITERATOR. " +
						       "ERROR WAS: " + err2.Error())
		}
		//Check if txn time is greater than the query time //
		time_txn :=time.Unix( response.Timestamp.Seconds, 
			                  int64(response.Timestamp.Nanos))
		if ( time_txn.Unix() < history.Timestamp || 
		     response.IsDelete) {continue;}

		// Add a comma before array members, suppress it for first one
		if bArrayMemberAlreadyWritten == true {buffer.WriteString(",")}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		data := make(map[string]interface{})
		json.Unmarshal(response.Value, &data)
		buffer.WriteString(", \"Value\":")
		dataAsBytes, _ := json.Marshal(data["Transaction"])
		buffer.WriteString(string(dataAsBytes))

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString( strconv.FormatInt(time_txn.Unix(),10) )
		buffer.WriteString("\"")
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
} 

/* -------------------------------------------------------------------------------------------------
mint: This function is called to perform a minting for a given token. Initially, the minted amount
	  is transferred to the admin. Args: is an array containing a json with the following attributes
	  and the admin public key in the second argument.
Token              string   // Symbol of the token to mint
Amount             float64  // Amount of new tokens to mint
Admin_Id           string   // Public Key ID of admin
------------------------------------------------------------------------------------------------- */


func (t *CoinBalanceSmartContract) mint( stub shim.ChaincodeStubInterface, 
	                                     args []string) pb.Response {

	// Retrieve information from the input //
	if len(args) != 2 {
		return shim.Error( "ERROR: MINT FUNCTION SHOULD BE CALLED " +
						   "WITH 2 ARGUMENT." )
	}
	minter := Minter{}
	input := []byte(args[0])
	json.Unmarshal(input, &minter)
	admin_public_id := args[1]

	// Get state of the token from the Ledger //
	token := Token{}
	tokenBytes, err1 := stub.GetState( IndexToken+minter.Token )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING STATE OF TOKEN " + 
						   minter.Token + ". " + err1.Error() ) 
	}
	json.Unmarshal(tokenBytes, &token)

	// Retrieve wallet of the admin //
	adminWalletBytes, err_2 := stub.GetState( IndexWallet + admin_public_id )
	if err_2 != nil {
		return shim.Error( "ERROR: RETRIEVING THE STATE OF " +
						   " ADMIN. " + err_2.Error() ) 
	}
	admin_wallet := MultiWallet{}
	json.Unmarshal(adminWalletBytes, &admin_wallet)

	// Transfer minted amount to admin //
	balance_token := admin_wallet.Balances[minter.Token]
	balance_token.Amount = balance_token.Amount + minter.Amount
	admin_wallet.Balances[minter.Token] = balance_token
	adminWalletBytes2, _ := json.Marshal(admin_wallet)
	err_3 := stub.PutState( IndexWallet+admin_public_id, adminWalletBytes2 )
	if err_3 != nil {
		return shim.Error( "ERROR TRANSFERING MINTED TOKEN " + token.Symbol +
						   " TO ADMIN BALANCE. " + err_3.Error() ) }

	// Mint amount of tokens and store value on Blockchain // 
	token.Supply = token.Supply + minter.Amount
	minting, _ := json.Marshal(token)
	err4 := stub.PutState(IndexToken+minter.Token, minting)
	if err4 != nil {
		return shim.Error( "ERROR: UPDATING TOKEN STATE " + 
							minter.Token + ". " + err4.Error() ) }
							
	return shim.Success(nil)
}


/* -------------------------------------------------------------------------------------------------
burn: This function is called to perform a burning for a given token. The burned amount is removed 
      from the admin balance. Args: is an array containing a json with the following attributes
	  and the admin public key in the second argument.
Token              string   // Symbol of the token to burn
Amount             float64  // Amount of new tokens to burn
Admin_Id           string   // Public Key ID of admin
------------------------------------------------------------------------------------------------- */


func (t *CoinBalanceSmartContract) burn( stub shim.ChaincodeStubInterface, 
	                                     args []string) pb.Response {

	// Retrieve information from the input //
	if len(args) != 2 {
		return shim.Error( "ERROR: BURN FUNCTION SHOULD BE CALLED " +
						   "WITH 2 ARGUMENT." ) }
	burner := Burner{}
	input := []byte(args[0])
	json.Unmarshal(input, &burner)
	admin_public_id := args[1]

	// Get state of the token from the Ledger //
	token := Token{}
	tokenBytes, err1 := stub.GetState( IndexToken+burner.Token )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING STATE OF TOKEN " + 
		                   burner.Token + ". " + err1.Error() ) }
	json.Unmarshal(tokenBytes, &token)

	// Retrieve wallet of the admin //
	adminWalletBytes, err_2 := stub.GetState( IndexWallet + admin_public_id )
	if err_2 != nil {
		return shim.Error( "ERROR: RETRIEVING THE STATE OF " +
						   " ADMIN. " + err_2.Error() ) }
	admin_wallet := MultiWallet{}
	json.Unmarshal(adminWalletBytes, &admin_wallet)

	// Check that the amount in Admin Wallet is enough to burn //
	balance_token := admin_wallet.Balances[burner.Token]
	if balance_token.Amount < burner.Amount {
		return shim.Error( "ERROR: ADMIN DOES NOT HOLD ENOUGH " +
		                   "TOKENS TO BURN THE AMOUNT GIVEN." ) } 
	
	// Burn tokens from Admin Balance //
	balance_token.Amount = balance_token.Amount - burner.Amount
	admin_wallet.Balances[burner.Token] = balance_token
	adminWalletBytes2, _ := json.Marshal(admin_wallet)
	err_3 := stub.PutState( IndexWallet+admin_public_id, adminWalletBytes2 )
	if err_3 != nil {
		return shim.Error( "ERROR BURNING TOKEN " + token.Symbol +
						   " FROM ADMIN BALANCE: " + err_3.Error() )  }
	
	// Burn amount of tokens and store value on Blockchain // 
	token.Supply = token.Supply - burner.Amount
	burning, _ := json.Marshal(token)
	err4 := stub.PutState(IndexToken+burner.Token, burning)
	if err4 != nil {
		return shim.Error( "ERROR: UPDATING TOKEN STATE " + 
		                    burner.Token + ". " + err4.Error() ) }
	return shim.Success(nil)
}


/* -------------------------------------------------------------------------------------------------
swap: This function is called to perform a swapping of user's token wallets by the same amount of
	  Fabric tokens version. It mints new tokens and transfer it to the user. 
	  Args: is an array containing a json with the following attributes
Token              string   // Symbol of the token to swap
Amount             float64  // Amount of tokens to swap
PublicId          string   // Public Key ID of user
------------------------------------------------------------------------------------------------- */


func (t *CoinBalanceSmartContract) swap( stub shim.ChaincodeStubInterface, 
										 args []string) pb.Response {
	// Retrieve information from the input //
	if len(args) != 1 {
		return shim.Error( "ERROR: MINT FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENT." ) }
	swapper := Swapper{}
	json.Unmarshal( []byte(args[0]), &swapper )

	// Get state of the token from the Ledger //
	token := Token{}
	tokenBytes, err1 := stub.GetState( IndexToken+swapper.Token )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING STATE OF TOKEN " + 
							swapper.Token + ". " + err1.Error() ) }
	json.Unmarshal(tokenBytes, &token)

	// Retrieve wallet of the user //
	userWalletBytes, err_2 := stub.GetState( IndexWallet + swapper.PublicId )
	if err_2 != nil {
		return shim.Error( "ERROR: RETRIEVING THE STATE OF USER " +  
		                    swapper.PublicId + ". " + err_2.Error())  }
	user_wallet := MultiWallet{}
	json.Unmarshal(userWalletBytes, &user_wallet)

	// Transfer swapping amount to user //
	balance_token := user_wallet.Balances[swapper.Token]
	balance_token.Amount = balance_token.Amount + swapper.Amount
	user_wallet.Balances[swapper.Token] = balance_token
	
	// Add transfer object //
	id := xid.New().String()
	date := getTimeNow()
	swap_transfer := Transfer{ 
		Type: "swap_ethereum", Token: swapper.Token, To: user_wallet.PublicId, Date: date,
		From: "Ethereum Wallet " + user_wallet.PublicId, Amount: swapper.Amount, Id: id }
	user_wallet.Transaction = []Transfer{ swap_transfer }
	
	userWalletBytes2, _ := json.Marshal(user_wallet)
	err_3 := stub.PutState( IndexWallet+swapper.PublicId , userWalletBytes2 )
	if err_3 != nil {
		return shim.Error( "ERROR TRANSFERING SWAPPING TOKEN " + token.Symbol +
						   " TO USER BALANCE. " + err_3.Error() ) }

	// Mint amount of tokens in the system and update state // 
	token.Supply = token.Supply + swapper.Amount
	minting, _ := json.Marshal(token)
	err4 := stub.PutState(IndexToken+swapper.Token, minting)
	if err4 != nil {
		return shim.Error( "ERROR: UPDATING TOKEN STATE " + 
							swapper.Token + ". " + err4.Error() ) }
	
	// Prepare output object with updates //
	update_wallets := make( map[string]MultiWallet )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	output := Output{ UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
withdraw: This function is called to perform a swapping back of user's token in Fabric version by the 
          same amount on original token version. It burns the Fabric tokens of the user.
	  	  Args: is an array containing a json with the following attributes
Token              string   // Symbol of the token to withdraw
Amount             float64  // Amount of tokens to withdraw
PublicId           string   // Public Key ID of user
------------------------------------------------------------------------------------------------- */


func (t *CoinBalanceSmartContract) withdraw( stub shim.ChaincodeStubInterface, 
											 args []string) pb.Response {
	// Retrieve information from the input //
	if len(args) != 1 {
		return shim.Error( "ERROR: WITHDRAW FUNCTION SHOULD BE CALLED " +
						   "WITH ONE ARGUMENT." )
	}
	withdrawer := Withdrawer{}
	input := []byte(args[0])
	json.Unmarshal(input, &withdrawer)

	// Get state of the token from the Ledger //
	token := Token{}
	tokenBytes, err1 := stub.GetState( IndexToken+withdrawer.Token )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING STATE OF TOKEN " + 
							withdrawer.Token + ". " + err1.Error() ) 
	}
	json.Unmarshal(tokenBytes, &token)

	// Retrieve wallet of the user //
	userWalletBytes, err_2 := stub.GetState( IndexWallet + withdrawer.PublicId )
	if err_2 != nil {
		return shim.Error( "ERROR: RETRIEVING THE STATE OF USER " +  
							withdrawer.PublicId + ". " + err_2.Error() ) 
	}
	user_wallet := MultiWallet{}
	json.Unmarshal(userWalletBytes, &user_wallet)

	// Burning withdrawing amount from user balance //
	balance_token := user_wallet.Balances[withdrawer.Token]
	balance_token.Amount = balance_token.Amount - withdrawer.Amount
	user_wallet.Balances[withdrawer.Token] = balance_token

	// Add transfer object //
	id := xid.New().String()
	date := getTimeNow()
	withdraw_transfer := Transfer{ 
		Type: "withdraw_ethereum", Token: withdrawer.Token, From: user_wallet.PublicId, Date: date,
		To: "Ethereum Wallet " + user_wallet.PublicId, Amount: withdrawer.Amount, Id: id }
	user_wallet.Transaction = []Transfer{ withdraw_transfer }
	
	userWalletBytes2, _ := json.Marshal(user_wallet)
	err_3 := stub.PutState( IndexWallet+withdrawer.PublicId, userWalletBytes2 )
	if err_3 != nil {
		return shim.Error( "ERROR WITHDRAWING TOKEN " + token.Symbol +
						   " TO USER BALANCE. " + err_3.Error() ) }

	// Burn amount of tokens in the sytem and update state // 
	token.Supply = token.Supply - withdrawer.Amount
	burning, _ := json.Marshal(token)
	err4 := stub.PutState( IndexToken+withdrawer.Token, burning )
	if err4 != nil {
	return shim.Error( "ERROR: UPDATING TOKEN STATE " + 
						withdrawer.Token + ". " + err4.Error() ) }

	// Prepare output object with updates //
	update_wallets := make( map[string]MultiWallet )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	output := Output{ UpdateWallets: update_wallets }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------- */

func main() {
	err := shim.Start(&CoinBalanceSmartContract{})
	if err != nil {
		fmt.Errorf("Error starting Token chaincode: %s", err)
	}
}
