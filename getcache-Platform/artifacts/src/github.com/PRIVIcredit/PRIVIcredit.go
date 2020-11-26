package main

import (
	//"time"
	"encoding/json"
	"fmt"
	"errors"
	"strconv"
	"math"
	"github.com/rs/xid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric/common/util"
)


/* -------------------------------------------------------------------------------------------------
Init:  this function register PRIVI as the Admin of the Network at the deployment of the 
       PRIVI Blockchain. Args: array containing a string:
PrivateKeyID           string   // Private Key of the admin of the smart contract
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) Init( stub shim.ChaincodeStubInterface ) pb.Response {

	_, args := stub.GetFunctionAndParameters()
	// Store in the state of the smart contract the Private Key of Admin //
	err1 := stub.PutState( IndexAdmin, []byte(args[0])  )
	if err1 != nil {
		return shim.Error( "ERROR: SETTING THE ADMIN PRIVATE KEY: " +
						   err1.Error() ) }
	if args[1] == "UPGRADE" { return shim.Success(nil) }
	// Initialise list of businesses in the Smart Contract as empty //
	credit_list, _ := json.Marshal( make(map[string]bool) )
	err2 := stub.PutState( IndexCreditList, credit_list )
	if err2 != nil {
		return shim.Error( "ERROR: INITIALISING THE CREDIT LIST: " +
	                       err2.Error() ) }
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
 The Invoke method is called as a result of an application request to run the Smart Contract ""
 The calling application program has also specified the particular smart contract function to be called
-------------------------------------------------------------------------------------------------*/

func (t *PRIVIcreditSmartContract) Invoke( stub shim.ChaincodeStubInterface ) pb.Response {

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
		case "initiatePRIVIcredit":
			return t.initiatePRIVIcredit(stub, args)
		case "getPRIVIcreditList":
			privi_list, err := t.getPRIVIcreditList(stub)
			if err != nil { return shim.Error( err.Error()) }
			priviListBytes, _ := json.Marshal( privi_list )
			return shim.Success( priviListBytes )
		case "getPRIVIcredit":
			privi_credit, err := t.getPRIVIcredit(stub, args[0])
			if err != nil { return shim.Error( err.Error()) }
			priviCreditBytes, _ := json.Marshal( privi_credit )
			return shim.Success( priviCreditBytes )
		case "getRiskParameters":
			risk_parameters, err := t.getRiskParameters(stub, args[0])
			if err != nil { return shim.Error( err.Error()) }
			riskParamsBytes, _ := json.Marshal( risk_parameters )
			return shim.Success( riskParamsBytes )
		case "updatePRIVIcredits":
			return t.updatePRIVIcredits(stub, args)
		case "modifyPRIVIparameters":
			return t.modifyPRIVIparameters(stub, args)
		case "withdrawFunds":
			return t.withdrawFunds(stub, args)
		case "depositFunds":
			return t.depositFunds(stub, args)
		case "borrowFunds":
			return t.borrowFunds(stub, args)
		case "assumePRIVIrisk":
			return t.assumePRIVIrisk(stub, args)
		case "managePRIVIcredits":
			return t.managePRIVIcredits(stub, args)
			
	}
	return shim.Error("Incorrect function name: " + function)
}


/* -------------------------------------------------------------------------------------------------
updateRiskParameters:  this function is called to update the risk parameters of the PRIVI Credit Smart
				  	   Contract for a given token. Args is an array containing a json with:
Token                string    // Token Symbol (args[0])
Interest_min		 float64   // Minimum interest rate allowed for PRIVI Credit
Interest_max		 float64   // Maximum interest rate allowed for PRIVI Credit
P_incentive_min      float64   // Minimum incentive rate allowed for PRIVI Credit
P_incentive_max      float64   // Maximum incentive rate allowed for PRIVI Credit
P_premium_min        float64   // Minimum premium rate allowed for PRIVI Credit
P_premium_max        float64   // Minimum premium rate allowed for PRIVI Credit
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) updateRiskParameters( stub shim.ChaincodeStubInterface, 
													     args []string) pb.Response {

	// Retrieve the input information of Risk Parameters //
	token := args[0]
	risk_token := RiskParameters{}
	err1 := json.Unmarshal( []byte(args[1]), &risk_token )
	if err1 != nil {
		return shim.Error( "ERROR: RETRIEVING THE INPUT " + err1.Error() ) }

	// Update parameters on Blockchain //
	riskParamsBytes, _ := json.Marshal( risk_token )
	err2 := stub.PutState( IndexRisk + token, riskParamsBytes )
	if err2 != nil {
		return shim.Error( "ERROR: UPDATING RISK PARAMS " + err2.Error() ) }
	return shim.Success(nil)
}


/* -------------------------------------------------------------------------------------------------
initiatePRIVIcredit  this function initialises a PRIVI credit loan with the parameters described
					 below. Args is an array containing a json with the following fields:
Creator              string                // Id of the creator of the PRIVI Credit
Token                string  			   // Token of the loan principal
MaxFunds         	 float64  			   // Max Funds allowed to be deposited 
Interest             float64  			   // Credit Interest of the Loan
Payments             int64   			   // Number of payments of the credit interest
Duration             int64                 // Duration of the Loan (in days)
P_incentive          float64               // Percentage of incentive
P_premium            float64  			   // Percentage of premium
TrustScore 			 float64               // Trust Score required for borrowers
EndorsementScore 	 float64               // Endorsement Score required for borrowers
InitialDeposit       float64               // Initial deposit of lending funds (in args[1])
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) initiatePRIVIcredit( stub shim.ChaincodeStubInterface, 
	                                  					args []string) pb.Response {

	// Retrieve the input information of PRIVI Credit Initialisation //
	priviLoan := PRIVIloan{}
	err1 := json.Unmarshal( []byte(args[1]), &priviLoan )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	initial_deposit, _ := strconv.ParseFloat( args[0], 64 )
	date := getTimeNow()

	// Retrieve the risk parameters of the PRIVI credit //
	risk_token, err2 := t.getRiskParameters( stub, priviLoan.Token )
	if err2 != nil { return shim.Error( err2.Error()) }

	// Retrieve wallet of PRIVI Credit Creator. Invoke CoinBalance Smart Contract //
	user_wallet, err3 := t.retrieveUserWallet( stub, priviLoan.Creator )
	if err3 != nil { return shim.Error( err3.Error()) }
	balance_creator := user_wallet.Balances[ priviLoan.Token ]

	// Check that the creator holds the funds and risk parameters satisfy conditions //
	if balance_creator.Amount < initial_deposit {
		return shim.Error( "ERROR: PRIVI CREDIT CREATOR " + priviLoan.Creator +
						   "DOES NOT HOLD ENOUGHT FUNDS FOR LENDING." ) }
	if priviLoan.Interest < risk_token.Interest_min || 
	   priviLoan.Interest > risk_token.Interest_max {
		return shim.Error( "ERROR: THE PRIVI CREDIT INTEREST SHOULD BE BETWEEN " +
	                       "THE BOUNDS REQUIRED.") } 
	if priviLoan.P_incentive < risk_token.P_incentive_min || 
	   priviLoan.P_incentive > risk_token.P_incentive_max {
		return shim.Error( "ERROR: THE PRIVI INCENTIVE RATE SHOULD BE BETWEEN " +
	                       "THE BOUNDS REQUIRED.") }
	if priviLoan.P_premium < risk_token.P_premium_min || 
	   priviLoan.P_premium > risk_token.P_premium_max {
		return shim.Error( "ERROR: THE PRIVI PREMIUM RATE SHOULD BE BETWEEN " +
	                       "THE BOUNDS REQUIRED.") }

	// Initialise state of the PRIVI loan //
	lenders := make( map[string]Lender )
	lender := Lender{ LenderId: priviLoan.Creator, 
					  Amount: initial_deposit, JoiningDay: 0. }
	lenders[priviLoan.Creator] = lender
	state := PRIVIstate{
		Funds: initial_deposit, Loaned: 0., LenderNum: 1, 
		BorrowerNum: 0, ProviderNum: 0, Status: "OFFER", 
		Loan_Day: 0, Total_Premium: 0., Total_Coverage: 0., 
		PremiumList: make( map[string]Premium ), Lenders: lenders,
		Borrowers: make( map[string]Borrower ), 
		Collaterals: make( map[string]float64 ) }

	// Complete PRIVI Loan Information and transfer funds //
	priviLoan.LoanId = xid.New().String() + xid.New().String() 
	priviLoan.Date = getTimeNow()
	priviLoan.State = state
	transactions := []Transfer{}
	balance_creator.Amount = balance_creator.Amount - initial_deposit
	balance_creator.PRIVI_lending = balance_creator.PRIVI_lending +
									initial_deposit

	PRIVI_transfer := Transfer{
		Type: "PRIVI_credit_creation", Token: priviLoan.Token, 
		From: priviLoan.Creator, To: "PRIVI Pool " + priviLoan.LoanId,
		Amount: initial_deposit, Id: xid.New().String(),
	    Date: date }

	// Update MultiWallet of user by invoking CoinBalance Chaincode //
	user_wallet.Transaction = append( transactions, PRIVI_transfer )
	balance_creator.PRIVIcreditLend[ priviLoan.LoanId ] = true
	user_wallet.Balances[ priviLoan.Token ] = balance_creator
	err4 := t.updateUserWallet( stub, user_wallet )
	if err4 != nil { return shim.Error( err4.Error()) }
	
	// Store new PRIVI Credit in PRIVI List  //
	credit_list, err5 := t.getPRIVIcreditList( stub )
	if err5 != nil { return shim.Error( err5.Error() ) }
	credit_list[ priviLoan.LoanId ] = true
	creditListBytes, _ := json.Marshal( credit_list )
	err6 := stub.PutState( IndexCreditList, creditListBytes )
	if err6 != nil {
		return shim.Error( "ERROR: UPDATING THE CREDIT LIST ON BLOCKCHAIN " +
							err6.Error() ) }

	// Store new PRIVI Credit on Blockchain //
	err7 := t.updatePRIVIcredit( stub, priviLoan )
	if err7 != nil { return shim.Error( err7.Error()) }

	// Generate response for outputs to update //
	update_wallets := make( map[string]MultiWallet )
	update_privi := make( map[string]PRIVIloan )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	update_privi[ priviLoan.LoanId ] = priviLoan
	output := Output{
		UpdateWallets: update_wallets, UpdateLoans: update_privi }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
updatePRIVIcredits: This function is called by the spend function of the CoinBalance smart contract
					with the PRIVI credit updates. 
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) updatePRIVIcredits( stub shim.ChaincodeStubInterface,
													   args []string ) pb.Response {
	// Loop through all the users to update mutiwallet //
	for _, arg := range args {
		privi_loan := PRIVIloan{}
		json.Unmarshal( []byte(arg), &privi_loan )
		// Update state of the PRIVI Credit // 
		err := t.updatePRIVIcredit( stub, privi_loan )
		if err != nil { return shim.Error( err.Error()) }
	}
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
modifyPRIVIparameters: this function is used to update the parameters of the PRIVI credit by the 
					   Creator if the loan has not been applied yet.
LoanId               string                // Id of the PRIVI Credit
Creator              string                // Id of the creator of the PRIVI Credit
MaxFunds         	 float64  			   // Max Funds allowed to be deposited 
Interest             float64  			   // Credit Interest of the Loan
Payments             float64  			   // Number of payments of the credit interest
Duration             float64               // Duration of the Loan (in days)
P_incentive          float64               // Percentage of incentive
P_premium            float64  			   // Percentage of premium
TrustScore 			 float64               // Trust Score required for borrowers
EndorsementScore 	 float64               // Endorsement Score required for borrowers
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) modifyPRIVIparameters( stub shim.ChaincodeStubInterface, 
												          args []string) pb.Response {

	// Retrieve the input information of PRIVI Credit Initialisation //
	priviLoan := PRIVIloan{}
	err1 := json.Unmarshal( []byte(args[0]), &priviLoan )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	loanId := priviLoan.LoanId

	// Retrieve PRIVI Credit information //
	privi_loan_old, err2 := t.getPRIVIcredit( stub, loanId )
	if err2 != nil { return shim.Error( err2.Error()) }

	// Check that the user is creator of the PRIVI Credit //
	if priviLoan.Creator != privi_loan_old.Creator  {
		return shim.Error( "ERROR: ONLY THE CREATOR CAN MODIFY LOAN." ) }
	if privi_loan_old.State.Status != "OFFER" {
		return shim.Error( "ERROR: THE LOAN SHOULD BE IN STATE OFFER." ) }

	// Retrieve the risk parameters of the PRIVI credit //
	risk_token, err3 := t.getRiskParameters( stub, privi_loan_old.Token )
	if err3 != nil { return shim.Error( err3.Error()) }

	// Check risk parameters satisfy conditions //
	if priviLoan.Interest < risk_token.Interest_min || 
	   priviLoan.Interest > risk_token.Interest_max {
		return shim.Error( "ERROR: THE PRIVI CREDIT INTEREST SHOULD BE BETWEEN " +
	                       "THE BOUNDS REQUIRED.") } 
	if priviLoan.P_incentive < risk_token.P_incentive_min || 
	   priviLoan.P_incentive > risk_token.P_incentive_max {
		return shim.Error( "ERROR: THE PRIVI INCENTIVE RATE SHOULD BE BETWEEN " +
	                       "THE BOUNDS REQUIRED.") }
	if priviLoan.P_premium < risk_token.P_premium_min || 
	   priviLoan.P_premium > risk_token.P_premium_max {
		return shim.Error( "ERROR: THE PRIVI PREMIUM RATE SHOULD BE BETWEEN " +
						   "THE BOUNDS REQUIRED.") }
	if priviLoan.MaxFunds < privi_loan_old.State.Funds {
		return shim.Error( "ERROR: THERE MAXFUNDS CANNOT BE REDUCED TO A QUANTITY" +
						   "INFERIOR TO THE FUNDS PRESENT IN THE CREDIT") }
	if priviLoan.Duration < privi_loan_old.State.Loan_Day {
		return shim.Error( "ERROR: THERE DURATION CANNOT BE REDUCED TO A QUANTITY" +
						   "INFERIOR TO THE DAYS ALREADY SPENT.") }
		   
	// Update PRIVI Loan info //
	priviLoan.State = privi_loan_old.State
	priviLoan.Date = privi_loan_old.Date
	priviLoan.Token = privi_loan_old.Token

	// Return updated PRIVI Credit //
	err4 := t.updatePRIVIcredit( stub, priviLoan )
	if err4 != nil { return shim.Error( err4.Error()) }

	// Generate response for outputs to update //
	update_privi := make( map[string]PRIVIloan )
	update_privi[ priviLoan.LoanId ] = priviLoan
	output := Output{ UpdateLoans: update_privi }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}



/* -------------------------------------------------------------------------------------------------
withdrawFunds: this function is called when the PRIVI Credit creator wants to withdraw funds if the
               loan has not been applied yet.
UserId        string   // ID of the user requesting to withdraw funds
LoanId        string   // ID of the PRIVI Credit
Amount        float64  // Amount to deposit in the PRIVI Credit
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) withdrawFunds( stub shim.ChaincodeStubInterface, 
												  args []string) pb.Response {

	// Retrieve the input information of PRIVI Credit Initialisation //
	withdraw := WithdrawFunds{}
	err1 := json.Unmarshal( []byte(args[0]), &withdraw )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	date := getTimeNow()

	// Retrieve PRIVI Credit information //
	privi_loan, err2 := t.getPRIVIcredit( stub, withdraw.LoanId )
	if err2 != nil { return shim.Error( err2.Error()) }

	// Check that the user is the Creator and that the amount to withdraw is present //
	if withdraw.UserId != privi_loan.Creator {
		return shim.Error( "ERROR: THE USER " + withdraw.UserId + " IS NOT THE CREATOR.") }
	if withdraw.Amount > privi_loan.State.Funds {
		return shim.Error( "ERROR: THE AMOUNT TO WITHDRAW IS GREATER THAN THE AMOUNT IN CREDIT") }
	if privi_loan.State.Status != "OFFER" {
		return shim.Error( "ERROR: THE LOAN SHOULD BE IN STATE OFFER." ) }

	// Retrieve wallet of the Lender Depositor and check the funds //
	user_wallet, err3 := t.retrieveUserWallet( stub, privi_loan.Creator )
	if err3 != nil { return shim.Error( err3.Error()) }
	balance_creator := user_wallet.Balances[ privi_loan.Token ]

	// Transfer the withdrawal funds on the PRIVI Credit to borrower //
	privi_loan.State.Funds = privi_loan.State.Funds - withdraw.Amount
	creator := privi_loan.State.Lenders[withdraw.UserId]
	creator.Amount = creator.Amount - withdraw.Amount
	privi_loan.State.Lenders[withdraw.UserId] = creator
	balance_creator.Amount = balance_creator.Amount + withdraw.Amount
	balance_creator.PRIVI_lending = balance_creator.PRIVI_lending - withdraw.Amount

	transactions := []Transfer{}
	PRIVI_withdraw := Transfer{
		Type: "PRIVI_credit_withdraw", Token: privi_loan.Token, 
		From: "PRIVI Pool " + privi_loan.LoanId, To: withdraw.UserId,
		Amount: withdraw.Amount, Id: xid.New().String(),
	    Date: date  }
	transactions = append( transactions, PRIVI_withdraw )

	// Update wallet of the user on blockchain //
	user_wallet.Balances[ privi_loan.Token ] = balance_creator
	user_wallet.Transaction = transactions
	err4 := t.updateUserWallet( stub, user_wallet )
	if err4 != nil { return shim.Error( err4.Error()) }

	// Return updated PRIVI Credit //
	err5 := t.updatePRIVIcredit( stub, privi_loan )
	if err5 != nil { return shim.Error( err5.Error()) }

	// Generate response for outputs to update //
	update_wallets := make( map[string]MultiWallet )
	update_privi := make( map[string]PRIVIloan )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	update_privi[ privi_loan.LoanId ] = privi_loan
	output := Output{
		UpdateWallets: update_wallets, UpdateLoans: update_privi }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
depositFunds: this function is called when a lender deposit funds on the PRIVI Credit.
              Args: an array containing a json with the following fields:
LenderId      string   // ID of the lender depositing funds
LoanId        string   // ID of the PRIVI Credit
Amount        float64  // Amount to deposit in the PRIVI Credit
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) depositFunds( stub shim.ChaincodeStubInterface, 
												 args []string) pb.Response {

	// Retrieve the input information of PRIVI Credit Initialisation //
	deposit := DepositFunds{}
	err1 := json.Unmarshal( []byte(args[0]), &deposit )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	date := getTimeNow()

	// Retrieve PRIVI Credit information //
	privi_loan, err2 := t.getPRIVIcredit( stub, deposit.LoanId )
	if err2 != nil { return shim.Error( err2.Error()) }

	// Retrieve wallet of the Lender Depositor and check the funds //
	user_wallet, err3 := t.retrieveUserWallet( stub, deposit.LenderId )
	if err3 != nil { return shim.Error( err3.Error()) }
	balance_lender := user_wallet.Balances[ privi_loan.Token ]
	if balance_lender.Amount < deposit.Amount {
		return shim.Error( "ERROR: LENDER DEPOSITOR " + deposit.LenderId +
						   " DOES NOT HOLD ENOUGHT FUNDS FOR LENDING." ) }

	// Check that the amount deposited does not exceed the max allowed //
	privi_loan.State.Funds = privi_loan.State.Funds + deposit.Amount
	if privi_loan.State.Funds > privi_loan.MaxFunds {
		return shim.Error( "ERROR: THE AMOUNT DEPOSITED EXCEEDS THE MAXIMUM " +
						   "AMOUNT FIXED BY PRIVI CREDIT CREATOR" ) }

	// Transfer the deposited funds on the PRIVI Credit //
	lenders := privi_loan.State.Lenders
	lender, isInList := lenders[ deposit.LenderId ] 
	if !isInList { 
		lender = Lender{ Amount: 0.,
			LenderId: deposit.LenderId, JoiningDay: privi_loan.State.Loan_Day }
		privi_loan.State.LenderNum = privi_loan.State.LenderNum + 1 }
	lender.Amount = lender.Amount + deposit.Amount
	lenders[ deposit.LenderId ] = lender
	privi_loan.State.Lenders = lenders

	balance_lender.Amount = balance_lender.Amount - deposit.Amount
	balance_lender.PRIVI_lending = balance_lender.PRIVI_lending + deposit.Amount
	balance_lender.PRIVIcreditLend[ privi_loan.LoanId ] = true
	transactions := []Transfer{}
	PRIVI_transfer := Transfer{
		Type: "PRIVI_credit_deposit", Token: privi_loan.Token, 
		From: deposit.LenderId, To: "PRIVI Pool " + privi_loan.LoanId,
		Amount: deposit.Amount, Id: xid.New().String(),
	    Date: date }
	transactions = append( transactions, PRIVI_transfer )

	// Update wallet of the user on blockchain //
	balance_lender.PRIVIcreditLend[ deposit.LoanId ] = true
	user_wallet.Balances[ privi_loan.Token ] = balance_lender
	user_wallet.Transaction = transactions
	err4 := t.updateUserWallet( stub, user_wallet )
	if err4 != nil { return shim.Error( err4.Error()) }

	// Return updated PRIVI Credit //
	err5 := t.updatePRIVIcredit( stub, privi_loan )
	if err5 != nil { return shim.Error( err5.Error()) }

	// Generate response for outputs to update //
	update_wallets := make( map[string]MultiWallet )
	update_privi := make( map[string]PRIVIloan )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	update_privi[ privi_loan.LoanId ] = privi_loan
	output := Output{
		UpdateWallets: update_wallets, UpdateLoans: update_privi }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
borrowFunds: this function is called when a borrower wants to borrow funds from a PRIVI Credit.
             Args: an array containing a json with the following fields:
BorrowerId      		string         // ID of the borrower borrowing funds
LoanId         			string         // ID of the PRIVI Credit
Amount         			float64        // Amount to borrow in the PRIVI Credit
Collaterals             float64        // Collaterals
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) borrowFunds( stub shim.ChaincodeStubInterface, 
												args []string) pb.Response {

	// Retrieve the input information //
	borrow := BorrowFunds{}
	err1 := json.Unmarshal( []byte(args[0]), &borrow )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	date := getTimeNow()

	// Retrieve PRIVI Credit information //
	privi_loan, err2 := t.getPRIVIcredit( stub, borrow.LoanId )
	if err2 != nil { return shim.Error( err2.Error()) }

	// Check that the borrower is not in the lender list //
	_, isInLenders := privi_loan.State.Lenders[ borrow.BorrowerId ]
	if isInLenders {
		return shim.Error( "ERROR: LENDERS CANNOT BORROW FUNDS FROM THEIR CREDITS") } 

	// Retrieve wallet of the Borrower and check the funds //
	user_wallet, err3 := t.retrieveUserWallet( stub, borrow.BorrowerId )
	if err3 != nil { return shim.Error( err3.Error()) }
	balance_borrower := user_wallet.Balances[ privi_loan.Token ]

	// Check that the Borrower satisfy the conditions and check if funds available //
	if user_wallet.TrustScore < privi_loan.TrustScore ||
	   user_wallet.EndorsementScore < privi_loan.EndorsementScore {
		return shim.Error( "ERROR: BORROWER " + borrow.LoanId +
						   "DOES NOT SATISFY THE REQUIREMENTS OF PRIVI CREDIT") }
	if borrow.Amount > privi_loan.State.Funds - privi_loan.State.Loaned  {
		return shim.Error( "ERROR: NO ENOUGH FUNDS TO SATISFY BORROWER REQUEST" ) }

	
	// Check that borrower has the collaterals and deposit it //
	transactions := []Transfer{}
	for token, amount := range(borrow.Collaterals) {
		balance_token := user_wallet.Balances[token]
		if balance_token.Amount < amount {
			return shim.Error( "BORROWER DOES NOT HOLD THE AMOUNT OF " +
							   "TOKEN  " + token + " TO DEPOSIT AS COLLATERAL" ) }
		balance_token.Amount = balance_token.Amount - amount
		user_wallet.Balances[token] = balance_token
		// Check if that collateral is in list //
		collateral, isInList := privi_loan.State.Collaterals[ token ]
		if !isInList { collateral = 0. }
		privi_loan.State.Collaterals[ token ] = collateral + amount
		// Check if that borrower is already in list //
		borrowers := privi_loan.State.Borrowers
		borrower, isInList := borrowers[ borrow.BorrowerId ] 
		if !isInList { 
			borrower = Borrower{ 
				Amount: 0., JoiningDay: privi_loan.State.Loan_Day,
				BorrowerId: borrow.BorrowerId, MissingPayments: 0, TotalPayments: 0,
				Debt: 0., TrustScore: user_wallet.TrustScore, 
				EndorsementScore: user_wallet.EndorsementScore, 
				Collaterals: make( map[string]float64 )  }
			privi_loan.State.BorrowerNum = privi_loan.State.BorrowerNum + 1 }
		// Check if borrower already deposited that collateral //
		collateral_borrower, isInList := borrower.Collaterals[ token ]
		if !isInList { collateral_borrower = 0. }
		borrower.Collaterals[ token ] = collateral_borrower + amount
		borrowers[ borrow.BorrowerId ] = borrower
		// Update transfer of coin //
		PRIVI_col_transfer := Transfer{
			Type: "PRIVI_credit_collateral", Token: token, 
			From: borrow.BorrowerId, To: "PRIVI Pool " + borrow.LoanId,
			Amount: amount, Id: xid.New().String(),
			Date: date }
		transactions = append( transactions, PRIVI_col_transfer )
	}

	// Transfer the borrowed funds on the PRIVI Credit //
	balance_borrower.Amount = balance_borrower.Amount + borrow.Amount
	balance_borrower.PRIVI_borrowing = balance_borrower.PRIVI_borrowing + borrow.Amount
	balance_borrower.PRIVIcreditBorrow[ privi_loan.LoanId ] = true
	privi_loan.State.Loaned = privi_loan.State.Loaned + borrow.Amount
													  
	borrower2 := privi_loan.State.Borrowers[ borrow.BorrowerId ]
	borrower2.Amount = borrower2.Amount + borrow.Amount
	borrower2.Debt = borrower2.Debt + borrow.Amount
	privi_loan.State.Borrowers[ borrow.BorrowerId ] = borrower2
	
	PRIVI_borrow_transfer := Transfer{
		Type: "PRIVI_credit_borrowing", Token: privi_loan.Token, 
		From: "PRIVI Pool " + borrow.LoanId,  To: borrow.BorrowerId, 
		Amount: borrow.Amount, Id: xid.New().String(),
	    Date: date }
	transactions = append( transactions, PRIVI_borrow_transfer )

	// Update wallet of the user on blockchain //
	balance_borrower.PRIVIcreditBorrow[ borrow.LoanId ] = true
	user_wallet.Balances[ privi_loan.Token ] = balance_borrower
	user_wallet.Transaction = transactions
	err4 := t.updateUserWallet( stub, user_wallet )
	if err4 != nil { return shim.Error( err4.Error()) }

	// Return updated PRIVI Credit //
	privi_loan.State.Status = "ISSUED"
	err5 := t.updatePRIVIcredit( stub, privi_loan )
	if err5 != nil { return shim.Error( err5.Error()) }

	// Generate response for outputs to update //
	update_wallets := make( map[string]MultiWallet )
	update_privi := make( map[string]PRIVIloan )
	update_wallets[ user_wallet.PublicId ] = user_wallet
	update_privi[ privi_loan.LoanId ] = privi_loan
	output := Output{
		UpdateWallets: update_wallets, UpdateLoans: update_privi }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}

/* -------------------------------------------------------------------------------------------------
assumePRIVIrisk: this function is called when a provider wants to assume risk on a PRIVI Loan.
			Args: an array containing a json with the following fields:
LoanId                  string          // ID of the PRIVI Credit
PremiumId               string          // ID of the premium
ProviderId      		string          // ID of the borrower borrowing funds
Risk_Pct         		float64         // ID of the PRIVI Credit
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) assumePRIVIrisk( stub shim.ChaincodeStubInterface, 
											        args []string) pb.Response {

	// Retrieve the input information //
	risk := AssumeRisk{}
	err1 := json.Unmarshal( []byte(args[0]), &risk )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " + err1.Error() ) }
	date := getTimeNow()

	// Retrieve PRIVI Credit information //
	privi_loan, err2 := t.getPRIVIcredit( stub, risk.LoanId )
	if err2 != nil { return shim.Error( err2.Error()) }

	// Check that the Provider Id is in the list of providers of the Credit //
	premium, isInList := privi_loan.State.PremiumList[risk.PremiumId]
	if !isInList {
		return shim.Error( "ERROR: THE PREMIUM ID " + risk.PremiumId +
						   " IS NOT IN THE LIST OF PROVIDERS OF THE PRIVI CREDIT." ) }

	// Retrieve wallet of the Provider and check the funds //
	provider_wallet, err3 := t.retrieveUserWallet( stub, risk.ProviderId )
	if err3 != nil { return shim.Error( err3.Error()) }
	balance_provider := provider_wallet.Balances[ privi_loan.Token ]

	// Check that Provider has the risking amount and add to the Pool //
	premium.Risk_Pct = risk.Risk_Pct
	risking_amount := premium.Risk_Pct * premium.Premium_Amount
	if balance_provider.Amount < risking_amount {
		return shim.Error( "ERROR: THE PROVIDER " + risk.ProviderId +
						   " DOES NOT HOLD ENOUGH FUNDS TO RISK.") }

	balance_provider.Amount = balance_provider.Amount - risking_amount
	balance_provider.PRIVI_lending = balance_provider.PRIVI_lending + risking_amount
	balance_provider.PRIVIcreditLend[ privi_loan.LoanId ] = true
	privi_loan.State.Total_Coverage = privi_loan.State.Total_Coverage + risking_amount
	privi_loan.State.PremiumList[risk.PremiumId] = premium
	PRIVI_risk_transfer := Transfer{
		Type: "PRIVI_risk_taking", Token: privi_loan.Token, 
		From: risk.ProviderId, To: "PRIVI Pool " + privi_loan.LoanId, 
		Amount: risking_amount, Id: xid.New().String(),
	    Date: date }
	transactions := []Transfer{ PRIVI_risk_transfer }

	// Update wallet of the user on blockchain //
	provider_wallet.Balances[ privi_loan.Token ] = balance_provider
	provider_wallet.Transaction = transactions
	err4 := t.updateUserWallet( stub, provider_wallet )
	if err4 != nil { return shim.Error( err4.Error()) }

	// Return updated PRIVI Credit //
	err5 := t.updatePRIVIcredit( stub, privi_loan )
	if err5 != nil { return shim.Error( err5.Error()) }

	// Generate response for outputs to update //
	update_wallets := make( map[string]MultiWallet )
	update_privi := make( map[string]PRIVIloan )
	update_wallets[ provider_wallet.PublicId ] = provider_wallet
	update_privi[ privi_loan.LoanId ] = privi_loan
	output := Output{
		UpdateWallets: update_wallets, UpdateLoans: update_privi }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}

/* -------------------------------------------------------------------------------------------------
interestPayment: this function is called to pay the interest of a loan. It receives as input
                 a PRIVILoan, updates its state and return it.
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) interestPayment( stub shim.ChaincodeStubInterface, 
	 update_wallets map[string]MultiWallet, privi_loan PRIVIloan ) (map[string]MultiWallet, PRIVIloan, error) {

	// Compute interest charged //
	P_interest := privi_loan.Interest / float64(privi_loan.Duration)
	total_amount_charged := 0.
	date := getTimeNow()
	err := errors.New("")

	// CHARGE INTEREST TO BORROWERS // 
	for borrower_id, borrower := range( privi_loan.State.Borrowers ) {
		// Retrieve wallet of the Borrower and check the funds //
		user_wallet, isInList := update_wallets[borrower_id]
		if !isInList {
			user_wallet, err = t.retrieveUserWallet( stub, borrower_id )
			if err != nil { return update_wallets, privi_loan, err }
			user_wallet.Transaction = []Transfer{}
			update_wallets[borrower_id] = user_wallet
		 }
	    balance_borrower := user_wallet.Balances[ privi_loan.Token ]
		// Charge amount //
		charge := borrower.Amount * P_interest
		if balance_borrower.Amount > charge {
			balance_borrower.Amount = balance_borrower.Amount - charge
			interest_transfer := Transfer{
				Type: "PRIVI_interest_charge", Token: privi_loan.Token, 
				From: borrower_id, To: "PRIVI Pool " + privi_loan.LoanId, 
				Amount: charge, Id: xid.New().String(),
				Date: date }
			user_wallet.Transaction = append( user_wallet.Transaction, 
				                              interest_transfer )
			user_wallet.Balances[privi_loan.Token] = balance_borrower
			total_amount_charged = total_amount_charged + charge
		// No Funds on Wallet //
		} else {
			borrower.MissingPayments = borrower.MissingPayments+1
			borrower.Debt = borrower.Debt+charge }
			borrower.TotalPayments = borrower.TotalPayments+1
		// Update Multiwallet and loan //
		privi_loan.State.Borrowers[ borrower_id ] = borrower
		update_wallets[borrower_id] = user_wallet 
	}
	
	// PAY INTEREST TO LENDERS // 
	for lender_id, lender := range( privi_loan.State.Lenders ) {
		user_wallet, isInList := update_wallets[lender_id]
		if !isInList {
			user_wallet, err = t.retrieveUserWallet( stub, lender_id )
			if err != nil { return update_wallets, privi_loan, err }
			user_wallet.Transaction = []Transfer{}
			update_wallets[lender_id] = user_wallet 
		}
		balance_lender := user_wallet.Balances[ privi_loan.Token ]
		// Paid amount //
		proportion := lender.Amount / privi_loan.State.Funds
		paid := total_amount_charged * proportion
		balance_lender.Amount = balance_lender.Amount + paid
		interest_transfer := Transfer{
			Type: "PRIVI_interest_payment", Token: privi_loan.Token, 
			From: "PRIVI Pool " + privi_loan.LoanId, To: lender_id,
			Amount: paid, Id: xid.New().String(),
			Date: date }
		user_wallet.Transaction = append( user_wallet.Transaction, 
			                              interest_transfer )
		// Update Multiwallet and loan //
		update_wallets[lender_id] = user_wallet 
	}
	return update_wallets, privi_loan, nil 
}

/* -------------------------------------------------------------------------------------------------
managePRIVIcredits: this function is called when a provider wants to assume risk on a PRIVI Loan.
			        Args: an array containing a json with the following fields:
LoanId                  string          // ID of the PRIVI Credit
PremiumId               string          // ID of the premium
ProviderId      		string          // ID of the borrower borrowing funds
Risk_Pct         		float64         // ID of the PRIVI Credit
------------------------------------------------------------------------------------------------- */

func (t *PRIVIcreditSmartContract) managePRIVIcredits( stub shim.ChaincodeStubInterface, 
													   args []string) pb.Response {
	// Retrieve the list of Active loans //
	privi_list, err1 := t.getPRIVIcreditList( stub )
	if err1 != nil { return shim.Error( err1.Error()) }

	update_wallets := make( map[string]MultiWallet )
	update_privi_loans := make( map[string]PRIVIloan )
	payment_credits := []string{}
	err := errors.New("")

	for privi_id, active := range(privi_list) {
		if !active {continue}
		// Retrieve PRIVI credit //
		privi_loan, err2 := t.getPRIVIcredit( stub, privi_id )
		if err2 != nil { return shim.Error( err2.Error()) }

		// Update Loan Day //
		loan_day := privi_loan.State.Loan_Day+1
		privi_loan.State.Loan_Day = loan_day

		// Update PRIVI credit and pay interest if it is loan day //
		step := int64( privi_loan.Duration / privi_loan.Payments )
		if math.Mod( float64(loan_day), float64(step)) == 0 {
			update_wallets, privi_loan, err = t.interestPayment(stub, 
										      update_wallets, privi_loan )
			if err != nil { shim.Error(err.Error()) }
			payment_credits = append(payment_credits, privi_id) }
		
		// Update state of Loan //
		err3 := t.updatePRIVIcredit( stub, privi_loan )
		if err3 != nil { return shim.Error( err3.Error()) }
		update_privi_loans[privi_id] = privi_loan
	}

	// Update MultiWallet of users by invoking CoinBalance Chaincode //
	for _, user_wallet := range(update_wallets) {
		err4 := t.updateUserWallet( stub, user_wallet )
		if err4 != nil { return shim.Error( err4.Error()) } }

	// Generate response for outputs to update //
	output := Output{
		UpdateWallets: update_wallets, UpdateLoans: update_privi_loans }
	outputBytes, _ := json.Marshal( output )
	return shim.Success( outputBytes )
}


/* -------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------- */

func main() {
	err := shim.Start(&PRIVIcreditSmartContract{})
	if err != nil {
		fmt.Errorf("Error starting Token chaincode: %s", err)
	}
}

