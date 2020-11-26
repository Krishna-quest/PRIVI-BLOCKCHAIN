package main

import (
	//"time"
	"encoding/json"
	"fmt"
	//"errors"
	//"strconv"
	//"math"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/rs/xid"
	//"github.com/hyperledger/fabric/common/util"
)

/* -------------------------------------------------------------------------------------------------
Init:  this function register PRIVI as the Admin of the POD Insurance Smart Contract. It initialises
       the lists of PODs and the Indexes used. Args: array containing a string:
PrivateKeyID           string   // Private Key of the admin of the smart contract
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) Init(stub shim.ChaincodeStubInterface) pb.Response {

	_, args := stub.GetFunctionAndParameters()
	// Store in the state of the smart contract the Private Key of Admin //
	err1 := stub.PutState(IndexAdmin, []byte(args[0]))
	if err1 != nil {
		return shim.Error( "ERROR: SETTING THE ADMIN PRIVATE KEY: " +
							err1.Error()) }
	if args[1] == "UPGRADE" { return shim.Success(nil) }

	// Initialise pod token parameters index as empty //
	insurance_list, _ := json.Marshal(make(map[string]bool))
	err2 := stub.PutState(IndexInsuranceList, insurance_list)
	if err2 != nil {
		return shim.Error( "ERROR: INITIALISING THE INSURANCE LIST: " +
							err2.Error()) }

	// Initialise parameters of smart contract as empty //
	parameters, _ := json.Marshal(Parameters{})
	err3:= stub.PutState(IndexParameters, parameters)
	if err3 != nil {
		return shim.Error( "ERROR: INITIALISING THE PARAMETER INDEX: " +
							err3.Error()) }
	
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
 The Invoke method is called as a result of an application request to run the Smart Contract ""
 The calling application program has also specified the particular smart contract function to be called
-------------------------------------------------------------------------------------------------*/

func (t *PodInsurance) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// Retrieve function and arguments //
	function, args := stub.GetFunctionAndParameters()

	// Retrieve caller of the function //
	caller, err1 := CallerCN(stub)
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE CALLER OF THE TRANSFER " +
						   "FUNCTION. " + err1.Error()) }
	// Retrieve admin of the smart contract //
	adminBytes, err2 := stub.GetState(IndexAdmin)
	if err2 != nil {
		return shim.Error( "ERROR: GETTING THE ADMIN OF THE SMART " +
						   "CONTRACT. " + err2.Error()) }
	admin := string(adminBytes)
	// Verify that the caller of the function is admin //
	if caller != admin {
		return shim.Error("ERROR: CALLER " + caller + " DOES NOT HAVE PERMISSION ") }

	// Call the proper function //
	switch function {
		case "retrieveInsuranceList":
			insurance_list, err := t.retrieveInsuranceList(stub)
			if err != nil { return shim.Error( err.Error() ) }
			insuranceListBytes, _ := json.Marshal( insurance_list )
			return shim.Success( insuranceListBytes )
	    case "retrieveInsuranceInfo":
			insurance_pool, err := t.retrieveInsuranceInfo(stub, args[0])
			if err != nil { return shim.Error( err.Error() ) }
			insurancePoolBytes, _ := json.Marshal( insurance_pool )
			return shim.Success( insurancePoolBytes )
		case "getParameters":
			parameters, err := t.getParameters( stub )
			if err != nil { return shim.Error( err.Error() ) }
			parametersBytes, _ := json.Marshal( parameters )
			return shim.Success( parametersBytes )
		case "updateParameters":
			return t.updateParameters(stub, args)
		case "initiateInsurancePool":
			return t.initiateInsurancePool(stub, args)
		case "investInsurancePool":
			return t.investInsurancePool(stub, args)
		case "subscribeInsurancePool":
			return t.subscribeInsurancePool(stub, args)
		case "withdrawInsurancePool":
			return t.withdrawInsurancePool(stub, args)
			
	}
	return shim.Error("Incorrect function name: " + function)
}


/* -------------------------------------------------------------------------------------------------
updateParameters:  this function is called to update the parameters of the POD Insurance Smart
				   Contract for a given token. Args is an array containing a json with:
MintingAmount		 float64   // Amount of minting privi insurance tokens (PI)
MintingFrequency     int64     // Frequency (in days) for minting PI tokens
VotingProtocol		 int64     // Number of days to vote a Proposal in the smart contract
MajorityProtocol     float64   // Pct of consensus needed to accept a change on the Smart Contract
ClaimConsensus       float64   // Pct of consensus needed to validate a claim 
ClaimVotingTime      int64     // Number of days that insurers has to vote for a claim
CourtVotingTime      int64     // Number of days that Members of DC has to vote
MinCovPct            float64   // Minimum pct required to the Guarantor to cover on principal
WithdrawCovPct       float64   // Threshold on the coverage pct which allows Guarantor to withdraw
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) updateParameters( stub shim.ChaincodeStubInterface, 
										 args []string ) pb.Response {

	// Retrieve the input information of Risk Parameters //
	parameters := Parameters{}
	err1 := json.Unmarshal( []byte(args[0]), &parameters )
	if err1 != nil { return shim.Error( "ERROR: RETRIEVING THE INPUT " + err1.Error()) }

	// Update parameters on the state of the Smart Contract //
	parametersBytes, _ := json.Marshal( parameters )
	err2 := stub.PutState( IndexParameters, parametersBytes )
	if err2 != nil {
		return shim.Error("ERROR: UPDATING PARAMETERS " + err2.Error()) }
	return shim.Success( parametersBytes )
}


/* -------------------------------------------------------------------------------------------------
initiateInsurancePool: this function initialises a new Insurance Pool for a NFT POD by a guarantor.
                       Args is an array containing a json with two fields:
Guarantor            string                // Id of the guarantor of the pool
PodId                string                // Id of the NFT POD being insured
Token                string  			   // Token of the 
Duration             int64                 // Duration of the Insurance
Frequency            int64                 // Frequency of payments of the Membership Fee
Valuation            float64     		   // Valuation of each of the NFT POD tokens by the Guarantor
FeeInscription       float64               // Fee paid for the inscription in an insurance pool
FeeMembership        float64               // Fee paid for a membership in an insurance pool
Coverage             int64                 // Number of pod tokens desired to cover
Deposit              float64               // Initial deposit of the Guarantor
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) initiateInsurancePool( stub shim.ChaincodeStubInterface,
									    	  args []string) pb.Response {
	
	// Output to update //
	update_wallets := make( map[string]MultiWallet )
	update_pods := make( map[string]POD )
	update_insurance := make( map[string]InsurancePool )
	transactions := []Transfer{}
	date := getTimeNow()

	// Retrieve the input information for Insurance Initialisation //
	insurance_pool := InsurancePool{}
	err1 := json.Unmarshal( []byte(args[0]), &insurance_pool )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) }

	pod, wallet, _, err_input := t.retrieveObjects( 
								stub, "", insurance_pool.Guarantor, 
								insurance_pool.PodId )
	if err_input != nil { return shim.Error(err_input.Error()) }

	ID := xid.New().String() + xid.New().String()

	// Check that the caller has the role of guarantor //
	actor, err2 := t.retrieveUserRole(stub, insurance_pool.Guarantor)
	if err2 != nil { shim.Error(err2.Error()) }
	if actor.Role != GUARANTOR_ROLE {
		shim.Error( "ONLY USERS WITH ROLE GUARANTOR CAN CREATE AN " + 
					"INSURANCE POOL.") }
					
	// Compute amount that is covered by Guarantor //
	balance_token := wallet.Balances[insurance_pool.Token]
	amount_coverage := float64( insurance_pool.Coverage ) *
	                   insurance_pool.Valuation

	// Check that the deposited amount is enought //
	if amount_coverage > insurance_pool.Deposit {
		return shim.Error( "THE GUARANTOR NEEDS TO DEPOSIT MORE FUNDS " + 
		            		"TO SATISFY THE COVERAGE PCT.") }
	if balance_token.Amount < insurance_pool.Deposit {
		return shim.Error( "THE GUARANTOR DOES NOT HOLD SUFFICIENT FUNDS " +
					       "TO CREATE THE INSURANCE POOL.") }
	
    // Compute transactions //
	balance_token.Amount = balance_token.Amount - insurance_pool.Deposit
	insurance_transfer := Transfer{
		Type: "insurance_creation", Token: insurance_pool.Token, 
		From: insurance_pool.Guarantor, To: "Insurance Pool " + ID,
		Amount: insurance_pool.Deposit, Id: xid.New().String(),
		Date: date }
	transactions = append( transactions, insurance_transfer )
	wallet.Balances[insurance_pool.Token] = balance_token	
	wallet.Transaction = []Transfer{ insurance_transfer }	

	// Charge claiming fee //
	owner_wallet := t.retrieveUserWallet( stub, pod.Creator )
	balance_owner := owner_wallet.Balances[ pod.Token ]
	amount_claiming := amount_coverage * CLAIMING_FEE
	if balance_owner.Amount < amount_claiming {
		return shim.Error( "OWNER DOES NOT HOLD SUFFICIENT FUNDS " +
						   "TO PAY THE CLAIMING FEE.") }
	balance_owner.Amount = balance_owner.Amount - amount_claiming
	claiming_transfer := Transfer{
		Type: "insurance_claiming_fee", Token: pod.Token, 
		From: pod.Creator, To: "Claiming Pool " + pod.PodId,
		Amount: amount_claiming, Id: xid.New().String(),
		Date: date }
	transactions = append( transactions, claiming_transfer )
	owner_wallet.Balances[ pod.Token ] = balance_owner
	owner_wallet.Transaction = []Transfer{ claiming_transfer }	
	
	// Update pod // 
	pod.State.ClaimingPool = pod.State.ClaimingPool + amount_claiming
	pod.State.Status = "INSURED"
	pod.Guarantors[ID] = insurance_pool.Coverage

	// Create new insurance pool and store in Blockchain //
	insurers := make( map[string]Insurer )
	insurers[ insurance_pool.Guarantor ] = Insurer{
					Amount: insurance_pool.Deposit, Date: date }
	clients := make( map[string]Client )
	clients[ pod.Creator ] = Client{ 
		               Amount: insurance_pool.Coverage, Date: date }
	state := InsuranceState{ Insurers: insurers, Clients: clients,
							 InsuredAmount: insurance_pool.Coverage,
		                     CoveragePool: insurance_pool.Deposit }
	insurance_pool.Id = ID
	insurance_pool.Date = getTimeNow()
	insurance_pool.State = state
	insurance_pool.Status = "ACTIVE"

	// Add in list of pod investors //
	pod.TotalInsurance = pod.TotalInsurance + insurance_pool.Deposit
	amount, isInList := pod.State.InsuredInvestors[pod.InvestorId]
	if !isInList { amount = 0 } 
	pod.State.InsuredInvestors[pod.InvestorId] = amount + 
											insurance_pool.Coverage

	// Add insurance pool to insurance pool list //
	insurance_list, err5 := t.retrieveInsuranceList(stub)
	if err5 != nil { shim.Error(err5.Error()) }
	insurance_list[ID] = true
	err6 := t.updateInsuranceList( stub, insurance_list )
	if err6 != nil  { return shim.Error(err6.Error()) }

	// Update blockchain //
	update_insurance[insurance_pool.Id] = insurance_pool
	update_pods[pod.PodId] = pod
	update_wallets[wallet.PublicId] = wallet
	update_wallets[owner_wallet.PublicId] = owner_wallet
	outputBytes, err_output := t.createOutput( stub, update_insurance,
											   update_wallets,
		                                       update_pods, transactions )
	if err_output != nil { return shim.Error(err_output.Error()) }
			
	return shim.Success( outputBytes )
}

/* -------------------------------------------------------------------------------------------------
investInsurancePool: this function is called when a user invest on an insurance pool.
Id                   string                // Id of the insurance pool to invest
InvestorId           string                // Id of the user that wants to invest
Amount               float64  			   // Amount desired to invest in the pool
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) investInsurancePool( stub shim.ChaincodeStubInterface,
											args []string) pb.Response {
	
	// Output to update //
	update_wallets := make( map[string]MultiWallet )
	update_pods := make( map[string]POD )
	update_insurance := make( map[string]InsurancePool )
	transactions := []Transfer{}
	date := getTimeNow()

	// Retrieve the input information //
	input := Invester{}
	err1 := json.Unmarshal( []byte(args[0]), &input )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) }
	pod, wallet, insurance, err_input := t.retrieveObjects( stub,
		                                 input.Id, input.InvestorId, "" )
    if err_input != nil { return shim.Error(err_input.Error()) }

	// Get info // 
	balance_investor := wallet.Balances[ insurance.Token ]
	if balance_investor.Amount < input.Amount {
		return shim.Error( "ERROR: INVESTOR DOES NOT HOLD ENOUGH FUNDS " +
						   "TO DEPOSIT IN THE INSURANCE POOL." ) }
	balance_investor.Amount = balance_investor.Amount - input.Amount

	// Do the transfer to the insurance pool //
	insurers := insurance.State.Insurers
	insurer, isInList := insurers[ input.InvestorId ]
	if !isInList { 
		insurer = Insurer{ Amount: input.Amount, Date: date } }
	insurer.Amount = insurer.Amount + input.Amount
	insurers[input.InvestorId] = insurer
	insurance.State.Insurers = insurers
	insurance.State.CoveragePool = insurance.State.CoveragePool + 
								   input.Amount
	pod.TotalInsurance = pod.TotalInsurance + input.Amount

	// Update wallet with transfer //
	invest_transfer := Transfer{ 
		Type: "insurance_investment", Token: insurance.Token, 
		From: input.InvestorId, To: "Insurance Pool " + insurance.Id,
		Amount: input.Amount, Id: xid.New().String(), Date: date }
	transactions = append( transactions, invest_transfer )
	wallet.Transaction = []Transfer{ invest_transfer }
	wallet.Balances[ insurance.Token ] = balance_investor

	// Update blockchain //
	update_insurance[insurance.Id] = insurance
	update_wallets[wallet.PublicId] = wallet
	update_pods[pod.PodId] = pod
	outputBytes, err_output := t.createOutput( stub, update_insurance, 
								   			   update_wallets,
		                                       update_pods, transactions )
	if err_output != nil { return shim.Error(err_output.Error()) }
	return shim.Success( outputBytes )
}

/* -------------------------------------------------------------------------------------------------
subscribeInsurancePool: this function is called when a pod investor wants to subscribe to the
                        insurance pool to cover its NFT pod tokens.
Id                      string                // Id of the insurance pool to subscribe
ClientId                string                // Id of the pod investor that wants coverage
Coverage                int64  			      // Amount of pod tokens desired to insure
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) subscribeInsurancePool( stub shim.ChaincodeStubInterface,
											   args []string) pb.Response {
	
	// Output to update //
	update_wallets := make( map[string]MultiWallet )
	update_pods := make( map[string]POD )
	update_insurance := make( map[string]InsurancePool )
	transactions := []Transfer{}
	date := getTimeNow()

	// Retrieve the input information //
	input := Subscriber{}
	err1 := json.Unmarshal( []byte(args[0]), &input )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) }
	pod, wallet, insurance, err_input := t.retrieveObjects( stub, input.Id, 
															input.ClientId, "" )
	
	if err_input != nil { return shim.Error(err_input.Error()) }
	balance_NFT := wallet.BalancesNFT[ pod.PodId ]
	balance_token := wallet.Balances[ insurance.Token ]

	// Get that the amount to insure is lesser than the amount held by the client // 
	clients := insurance.State.Clients
	client, isInList := clients[input.ClientId]  
	if !isInList { client = Client{ Amount: 0, Date: date} }
	client.Amount = client.Amount + input.Coverage
	if balance_NFT.Amount < client.Amount {
		return shim.Error( "ERROR: IT IS NOT POSSIBLE TO PURCHASE " +
						   "MORE COVERAGE THAN THE HELD AMOUNT OF NFT TOKENS." ) }
	clients[input.ClientId] = client
	insurance.State.Clients = clients
	insurance.State.InsuredAmount = insurance.State.InsuredAmount + input.Coverage

	// Update wallet with transfer with charged amount for inscription //
	charge_inscription := float64(input.Coverage) * insurance.Valuation *
	                      insurance.FeeInscription				  
	charge_claiming_fee := float64(input.Coverage) * insurance.Valuation *
						   CLAIMING_FEE
	if charge_inscription + charge_claiming_fee > balance_token.Amount {
		return shim.Error( "ERROR: CLIENT DOES NOT HOLD FUNDS TO " +
						   "PURCHASE THE REQUIRED INSURANCE." ) }
	balance_token.Amount = balance_token.Amount - charge_inscription -
						   charge_claiming_fee
	invest_transfer := Transfer{
		Type: "insurance_purchase", Token: insurance.Token, 
		From: input.ClientId, To: "Insurance Pool " + insurance.Id,
		Amount: charge_inscription, Id: xid.New().String(), Date: date }
	transactions = append( transactions, invest_transfer )
	claiming_transfer := Transfer{
		Type: "insurance_claiming_fee", Token: insurance.Token, 
		From: input.ClientId, To: "Claiming Pool " + pod.PodId,
		Amount: charge_claiming_fee, Id: xid.New().String(), Date: date }
	transactions = append( transactions, claiming_transfer )
	wallet.Transaction = []Transfer{ invest_transfer }
	wallet.Balances[ insurance.Token ] = balance_token

	// Update pod insurance //
	amount, isInList := pod.State.InsuredInvestors[pod.InvestorId]
	if !isInList { amount = 0 } 
	pod.State.InsuredInvestors[pod.InvestorId] = amount + 
												 input.Coverage
	pod.State.ClaimingPool = pod.State.ClaimingPool + charge_claiming_fee
	

	// Transfer charged amount to insurance investors proportionally //
	for investor_id, investor := range(insurance.State.Insurers) {
		proportion := investor.Amount / insurance.State.CoveragePool
		payment := charge_inscription * proportion
		wallet_investor, err2 := t.retrieveUserWallet( stub, investor_id )
		if err2 != nil { return shim.Error(err2.Error()) }
		// Transfer return to investor //
		balance_investor := wallet_investor.Balances[insurance.Token]
		balance_investor.Amount = balance_investor.Amount + payment
		payment_transfer := Transfer{
			Type: "insurance_inscription_get", Token: insurance.Token, 
			From: "Insurance Pool " + insurance.Id, To: investor_id, 
			Amount: payment, Id: xid.New().String(), Date: date }
		wallet_investor.Balances[insurance.Token] = balance_investor
		wallet_investor.Transaction = []Transfer{ payment_transfer }
		transactions = append( transactions, payment_transfer )
		update_wallets[investor_id] = wallet_investor
	}

	// Update blockchain //
	update_insurance[insurance.Id] = insurance
	update_wallets[wallet.PublicId] = wallet
	outputBytes, err_output := t.createOutput( stub, update_insurance, 
								   			   update_wallets,
		                                       update_pods, transactions )
	if err_output != nil { return shim.Error(err_output.Error()) }
	return shim.Success( outputBytes )
}

/* -------------------------------------------------------------------------------------------------
withdrawInsurancePool: this function is called when a user wants to withdraw some amount from the
                       insurance pool.
Id                   string                // Id of the insurance pool to withdraw
InvestorId           string                // Id of the user that wants to withdraw
Amount               float64  			   // Amount desired to withdraw 
------------------------------------------------------------------------------------------------- */

func (t *PodInsurance) withdrawInsurancePool( stub shim.ChaincodeStubInterface,
											  args []string) pb.Response {

	// Output to update //
	update_wallets := make( map[string]MultiWallet )
	update_pods := make( map[string]POD )
	update_insurance := make( map[string]InsurancePool )
	transactions := []Transfer{}
	date := getTimeNow()

	// Retrieve the input information //
	input := Invester{}
	err1 := json.Unmarshal( []byte(args[0]), &input )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) }
	pod, wallet, insurance, err_input := t.retrieveObjects( stub, input.Id, 
														  input.InvestorId, "" )
	if err_input != nil { return shim.Error(err_input.Error()) }						

	// Check that user is in the list of insurance investors //
	amount := checkWithdrawEnabled( insurance, pod, input.InvestorId )
	if amount <= input.Amount {
		return shim.Error( "ERROR: WITHDRAWING IS NOT ALLOWED." ) }
	invested, isInList := insurance.State.Insurers[ input.InvestorId ]
	if !isInList {
		return shim.Error( "ERROR: ID " + input.InvestorId + 
						   " IS NOT IN THE INSURER LIST." ) }
    if invested.Amount < input.Amount {
		return shim.Error( "ERROR: ID " + input.InvestorId + 
						   " DID NOT INVEST THE AMOUNT DESIRED TO WITHDRAW." )  }
	invested.Amount = invested.Amount - input.Amount
	insurance.State.CoveragePool = insurance.State.CoveragePool - input.Amount
	insurance.State.Insurers[ input.InvestorId ] = invested
	pod.TotalInsurance = pod.TotalInsurance - input.Amount
	if invested.Amount == 0 {
		delete(insurance.State.Insurers, input.InvestorId ) }
												
	// Withdraw funds // 
	balance := wallet.Balances[insurance.Token]
	balance.Amount = balance.Amount + input.Amount
	withdrawal_transfer := Transfer{
		Type: "withdraw_insurance_funds", Token: insurance.Token, 
		From: "Insurance Pool " + insurance.Id, To: input.InvestorId, 
		Amount: input.Amount, Id: xid.New().String(), Date: date }
	wallet.Balances[insurance.Token] = balance
	wallet.Transaction = []Transfer{ withdrawal_transfer }
	transactions = append( transactions, withdrawal_transfer )

	// Update blockchain //
	update_insurance[insurance.Id] = insurance
	update_wallets[wallet.PublicId] = wallet
	update_pods[pod.PodId] = pod
	outputBytes, err_output := t.createOutput( stub, update_insurance, 
								   			   update_wallets,
		                                       update_pods, transactions )
	if err_output != nil { return shim.Error(err_output.Error()) }
	return shim.Success( outputBytes )
}



/* -------------------------------------------------------------------------------------------------
claimPOD: this function is called when the pod owner makes a claim for the NFT Pod. This sends a
		  notification to all the Guarantors that have a certain amount of times to validate.
		  The decision if under consensus of Guarantors.
PodId                string                // Id of the pod submitting the claim
------------------------------------------------------------------------------------------------- */
/*
func (t *PodInsurance) claimPOD( stub shim.ChaincodeStubInterface,
							args []string) pb.Response {
		(1) Retrieve POD and activate claim counter.
	}
*/


/* -------------------------------------------------------------------------------------------------
claimVote: this function is called when the pod owner makes a claim for the NFT Pod. This sends a
		  notification to all the Guarantors that have a certain amount of times to validate.
		  The decision if under consensus of Guarantors.
PodId                string                // Id of the pod submitting the claim
------------------------------------------------------------------------------------------------- */
/*
func (t *PodInsurance) claimVote( stub shim.ChaincodeStubInterface,
	args []string) pb.Response {
		(1) Retrieve POD and activate claim counter.
	}
*/






/* -------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------- */

func main() {
	err := shim.Start(&PodInsurance{})
	if err != nil {
		fmt.Errorf("Error starting Pod Swapping chaincode: %s", err)
	}
}
