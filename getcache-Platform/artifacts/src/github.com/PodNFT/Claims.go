
package main

import (
	//"time"
	"encoding/json"
	//"errors"
	//"strconv"
	//"math"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/rs/xid"
	//"github.com/hyperledger/fabric/common/util"
)

/* -------------------------------------------------------------------------------------------------
initiateClaimProposal: this function is called when the a insured pod investors starts a claim
                       votation for the NFT Pod. 
PodId                string                // Id of the pod submitting the claim
Voter                string                // Id of the investor submittin the claim
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) initiateClaimProposal( stub shim.ChaincodeStubInterface,
								 	    args []string) pb.Response {

	// Retrieve input information in a deletion object //
	input := Voter{}
	err1 := json.Unmarshal( []byte(args[0]), &input )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) }
	date := getTimeNow()

	// Retrieve pod //
	pod, err2 := t.retrievePodInfo(stub, input.PodId)
	if err2 != nil { shim.Error(err2.Error()) }

	// Check that InvestorId has rights to claim (it is insured) //
	_, isInList := pod.State.InsuredInvestors[input.InvestorId]
	if !isInList {
		return shim.Error( "ERROR: THE INVESTOR HAS NOT RIGHTS TO START A " +
		                   "VOTING PROCESS.") }

	// Check that pod is not already in a claiming process //
	if pod.Claiming.Status != "NONE" {
		return shim.Error( "ERROR: THE POD IS ALREADY UNDER A CLAIM.") }
	
	// Submit Proposal //
	pod.Claiming.Status = "PROPOSAL"
	pod.Claiming.VotingYes = []string{ input.Voter }
	pod.Claiming.VotingNo = []string{}
	pod.Claiming.VotingNothing = []string{}
	votes := make( map[string]bool )
	votes[input.Voter] = true
	pod.Claiming.Votation = VotingProcess{ 
		StartDate: date, Votes: votes, Duration: pod.VotingTime }

	// Update pod on blockchain //
	_, err3 := t.updatePod(stub, pod)
	if err3 != nil { return shim.Error(err3.Error()) }
	txn := Transfer{
		Type: "Proposal_Claim_Init", Amount: 0., Token: "", 
		From: input.Voter, To: "Pod NFT " + pod.PodId,
		Id: xid.New().String(), Date: date }

	// Retrieve output //
	output := OutputNFT{ Transaction: []Transfer{ txn } }
	outputBytes, _ := json.Marshal(output)
	return shim.Success(outputBytes)
}

/* -------------------------------------------------------------------------------------------------
voteClaimProposal: this function is called when the a insured pod investor votes for a claim
                   proposal.
PodId                string                // Id of the pod submitting the claim
Voter                string                // Id of the investor submittin the claim
Vote                 bool                  // Vote of the insured pod investor
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) voteClaimProposal( stub shim.ChaincodeStubInterface,
								 	  args []string) pb.Response {

	// Retrieve input information in a deletion object //
	input := Voter{}
	err1 := json.Unmarshal( []byte(args[0]), &input )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) }
	date := getTimeNow()

	// Retrieve pod //
	pod, err2 := t.retrievePodInfo(stub, input.PodId)
	if err2 != nil { shim.Error(err2.Error()) }

	// Check that pod in a claiming process proposal //
	if pod.Claiming.Status != "PROPOSAL" {
		return shim.Error( "ERROR: THE POD IS NOT IN PROPOSAL PHASE.") }

	// Check that InvestorId has rights to claim (it is insured) //
	_, isInList := pod.State.InsuredInvestors[input.Voter]
	if !isInList {
		return shim.Error( "ERROR: THE INVESTOR HAS NOT RIGHTS TO VOTE") }
	
	// Submit Vote //
	votes := make( map[string]bool )
	pod.Claiming.Votation.Votes[input.Voter] = true

	// Update pod on blockchain //
	_, err3 := t.updatePod(stub, pod)
	if err3 != nil { return shim.Error(err3.Error()) }
	txn := Transfer{
		Type: "Proposal_Claim_Vote", Amount: 0., Token: "", 
		From: input.Voter, To: "Pod NFT " + pod.PodId,
		Id: xid.New().String(), Date: date }

	// Retrieve output //
	output := OutputNFT{ Transaction: []Transfer{ txn } }
	outputBytes, _ := json.Marshal(output)
	return shim.Success(outputBytes)
}

/* -------------------------------------------------------------------------------------------------
endProposalVoting: this function is called when the Claim Proposal is identified to have arrived to 
				   its end by the backend.
PodId                string                // Id of the pod submitting the claim
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) endProposalVoting( stub shim.ChaincodeStubInterface,
								    args []string) pb.Response {

	// Retrieve input information in a deletion object //
	pod_id := args[0]
	date := getTimeNow()

	// Retrieve pod //
	pod, err2 := t.retrievePodInfo(stub, pod_id)
	if err2 != nil { shim.Error(err2.Error()) }
	claiming_process := pod.Claiming

	// Check if the time for voting proposal is done //
	period := getTimePeriod( claiming_process.Votation.StartDate ) 
	if float64(period) < pod.Claiming.Votation.Duration {
		return shim.Error( "ERROR: THE VOTING TIME DID NOT END UP.")}

	// Check that pod in a claiming process proposal //
	if pod.State.ClaimStatus != "PROPOSAL" {
		return shim.Error( "ERROR: THE POD IS NOT IN PROPOSAL PHASE.") }

	// Compute total voting of insurers that vote //
	total_voting_power := 0.
	for id, vote := range(claiming_process.Votation.Votes) {
		vote := pod.State.InsuredInvestors[id]
		total_voting_power = total_voting_power + float64(vote) }
	
	// Votation result // 
	votation_pct := 0.
	for id, power := range(pod.State.InsuredInvestors) {
		voting_weight := float64(power) / total_voting_power
		vote, isInList := pod.Claiming.Votation.Votes[id]
		if !isInList {
			claiming_process.VotingNothing = append( 
							 claiming_process.VotingNothing, id ) 
		} else if vote {
			votation_pct = votation_pct + voting_weight
			claiming_process.VotingYes = append( 
				             claiming_process.VotingYes, id )
		} else {
			claiming_process.VotingNo = append( 
				             claiming_process.VotingNo, id ) }
	}

	// Case (1): consensus reached //
	txn := Transfer{}
	if votation_pct >= CONSENSUS_CLAIM_PROPOSAL {
		claiming_process.Status = "VALIDATION"
		claiming_process.Votation = VotingProcess{
			StartDate: date, Votes: make( map[string]bool ),
			Duration: pod.VerifTime }
		txn = Transfer{
			Type: "Proposal_Claim_Approved", Amount: 0., Token: "", 
			From: "Pod NFT " + pod.PodId, To: "Pod_Investors+Guarantors",
			Id: xid.New().String(), Date: date }
	// Case (2): claim not validated //
	} else { 
		claiming_process = ClaimingProcess{}
		txn = Transfer{
			Type: "Proposal_Claim_Denied", Amount: 0., Token: "", 
			From: "Pod NFT " + pod.PodId, To: "Pod_Investors",
			Id: xid.New().String(), Date: date } }

	// Update pod on blockchain //
	pod.Claiming = claiming_process
	_, err3 := t.updatePod(stub, pod)
	if err3 != nil { return shim.Error(err3.Error()) }

	// Retrieve output //
	output := OutputNFT{ Transaction: []Transfer{ txn } }
	outputBytes, _ := json.Marshal(output)
	return shim.Success(outputBytes)
}


/* -------------------------------------------------------------------------------------------------
voteClaimValidation: this function is called when a Guarantor votes for a claim after validations.
PodId                string                // Id of the pod submitting the claim
Voter                string                // Id of the guarantor verifying a claim
Vote                 bool                  // Vote of the insured pod investor
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) voteClaimValidation( stub shim.ChaincodeStubInterface,
								      args []string) pb.Response {

	// Retrieve input information in a deletion object //
	input := Voter{}
	err1 := json.Unmarshal( []byte(args[0]), &input )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) })
    date := getTimeNow()

	// Retrieve pod //
	pod, err2 := t.retrievePodInfo(stub, input.PodId)
	if err2 != nil { shim.Error(err2.Error()) }

	// Check that the caller has the role of guarantor //
	actor, err3 := t.retrieveUserRole(stub, insurance_pool.Guarantor)
	if err3 != nil { shim.Error(err3.Error()) }
	if actor.Role != GUARANTOR_ROLE {
		shim.Error( "ONLY USERS WITH ROLE GUARANTOR CAN VALIDATE A CLAIM") }


	// Check that pod in a claiming validation process //
	if pod.Claiming.Status != "VALIDATION" {
		return shim.Error( "ERROR: THE POD IS NOT IN VALIDATION PHASE.") }

	// Check that Guarantor has rights to claim (it is insured) //
	_, isInList := pod.State.Guarantors[input.Voter]
	if !isInList {
		return shim.Error( "ERROR: THE VOTER HAS NOT RIGHTS TO VOTE") }
	
	// Submit Vote //
	votes := make( map[string]bool )
	pod.Claiming.Votation.Votes[input.Voter] = true

	// Update pod on blockchain //
	_, err4 := t.updatePod(stub, pod)
	if err4 != nil { return shim.Error(err4.Error()) }
	txn := Transfer{
		Type: "Validation_Claim_Vote", Amount: 0., Token: "", 
		From: input.Voter, To: "Pod NFT " + pod.PodId,
		Id: xid.New().String(), Date: date }

	// Retrieve output //
	output := OutputNFT{ Transaction: []Transfer{ txn } }
	outputBytes, _ := json.Marshal(output)
	return shim.Success(outputBytes)
}

/* -------------------------------------------------------------------------------------------------
endValidationVotation: this function is called when the Claim Validation is identified to have arrived 
                       to its end by the backend.
PodId                string                // Id of the pod submitting the claim
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) endValidationVotation( stub shim.ChaincodeStubInterface,
								   		args []string) pb.Response {

	// Retrieve input information in a deletion object //
	pod_id := args[0]
	date := getTimeNow()
	update_wallets := make( map[string]MultiWallet )
	update_pods := make( map[string]POD )
	transaction := []Transfer{}

	// Retrieve pod //
	pod, err2 := t.retrievePodInfo(stub, pod_id)
	if err2 != nil { shim.Error(err2.Error()) }
	claiming_process := pod.Claiming

	// Check if the time for voting proposal is done //
	period := getTimePeriod( claiming_process.Votation.StartDate ) 
	if float64(period) < pod.Claiming.Votation.Duration {
		return shim.Error( "ERROR: THE VOTING TIME DID NOT END UP.")}

	// Check that pod is in a claiming validation process //
	if pod.State.ClaimStatus != "VALIDATION" {
		return shim.Error( "ERROR: THE POD IS NOT IN PROPOSAL PHASE.") }

	// Compute total voting of insurers that vote //
	total_voting_power := 0.
	for id, vote := range(claiming_process.Votation.Votes) {
		vote := pod.State.Guarantors[id]
		total_voting_power = total_voting_power + float64(vote) }
	
	// Votation result // 
	votation_pct := 0.
	for id, power := range(pod.State.Guarantors) {
		voting_weight := float64(power) / total_voting_power
		vote, isInList := pod.Claiming.Votation.Votes[id]
		if !isInList {
			claiming_process.VotingNothing = append( 
							 claiming_process.VotingNothing, id ) 
		} else if vote {
			votation_pct = votation_pct + voting_weight
			claiming_process.VotingYes = append( 
				             claiming_process.VotingYes, id )
		} else {
			claiming_process.VotingNo = append( 
				             claiming_process.VotingNo, id ) }
	}

	///////////////////////// RESULT OF VOTATION ////////////////////////////

	// Case (1): consensus reached to accept claim //
	if votation_pct >= CONSENSUS_CLAIM_VALIDATION {
		update_wallets, transactions = payoffVotation( pod, claiming_process, 
													   true) 
		// TODO: liquidate function for all insurance pools //
		pod.ClaimingProcess = ClaimingProcess{}
		pod.State.Status = "LIQUIDATED"
		update_pods[ pod.Id ] = pod
	// Case (2): consensus reached to no accept claim //
	} else if votation_pct <= 1. - CONSENSUS_CLAIM_VALIDATION { 
		update_wallets, transactions = payoffVotation( pod, claiming_process, 
													   flase) 
		pod.ClaimingProcess = ClaimingProcess{}
		update_pods[ pod.Id ] = pod
	// Case (3): no consensus, go to digital courts
	} else {
		claiming_process.Status = "COURT"
		claiming_process.Votation = VotingProcess{
			StartDate: date, Votes: make( map[string]bool ),
			Duration: pod.CourtTime }
		txn = Transfer{
			Type: "Claim_Digital_Court", Amount: 0., Token: "", 
			From: "Pod NFT " + pod.PodId, 
			To: "Pod_Investors+Guarantors+DigitalCourt",
			Id: xid.New().String(), Date: date }
		transactions = []Transfer{ txn }
		pod.Claiming = claiming_process
	}

	// Return and update output //
	outputBytes, err4 := t.createOutput( stub,
		 update_wallets, update_pods, transactions )
    if err4 != nil { shim.Error(err4.Error()) }
	return shim.Success(outputBytes)
}


/* -------------------------------------------------------------------------------------------------
payoffVotation: this function is called when the a Votation is resolved. Bad actors loss his share 
                of the Claiming pool and good actors get that part.
------------------------------------------------------------------------------------------------- */

func payoffVotation( pod POD, claiming_process ClaimingProcess, 
	                 decision bool ) ( map[string]MultiWallet, []Transfer ) {
	
	date := getTimeNow()
	update_wallets := make( map[string]MultiWallet )
	transactions := []Transfer{}

	// Get bad and good actors //
	bad_actors := claiming_process.VotingNothing
	good_actors := []string{}
	if decision { 
		bad_actors = append(bad_actors, claiming_process.VotingNo) 
		good_actors = claiming_process.VotingYes
	} else {
		bad_actors = append(bad_actors, claiming_process.VotingYes) 
		good_actors = claiming_process.VotingNo }

	// Voting power Guarantors // 
	voting_power_guarantors := 0.
	for _, power := range(pod.Guarantors) {
		voting_power_guarantors = voting_power_guarantors + power }
	
	// Voting power Guarantors // 
	voting_power_investors := 0.
	for _, power := range(pod.State.InsuredInvestors) {
		voting_power_investors = voting_power_investors + power }

	// Get claiming share from bad actors //
	claiming_share := 0.
	for _, actor := range(bad_actors) {
		pct = 0.
		amount, isInList := pod.Guarantors[ actor ]
		if isInList { pct = amount / voting_power_guarantors }
		else {
			amount, isInList = pod.State.InsuredInvestors[ actor ]
			pct = amount / voting_power_investors }
		claiming_share = claiming_share + pct * 0.5 }

	redeem := claiming_share * pod.State.ClaimingPool
	// Transfer funds to good actors //
	for id, actor := range(good_actors) {
		pct = 0.
		amount, isInList := pod.Guarantors[ actor ]
		if isInList { pct = amount / voting_power_guarantors }
		else {
			amount, isInList = pod.State.InsuredInvestors[ actor ]
			pct = amount / voting_power_investors }
		payoff = pct * redeem * 0.5
		// Get balance of user to pay //
		user_wallet := t.retrieveUserWallet( stub, id )
		balance_wallet := user_wallet.Balances[ pod.Token ]
		balance_wallet.Amount = balance_wallet.Amount + payoff
		txn = Transfer{
			Type: "Claim_Good_Vote", Amount: payoff, 
			Token: pod.Token, 
			From: "Claiming Pool " + pod.PodId, To: id,
			Id: xid.New().String(), Date: date }
		user_wallet.Balances[ pod.Token ] = balance_wallet
		user_wallet.Transaction = []Transfer{txn}
		transactions = append( transactions, txn )
		update_users[ id ] = user_wallet
	 }

	 return update_wallets, transactions
}

	
/* -------------------------------------------------------------------------------------------------
voteDigitalCourt: this function is called when a member of Digital Court votes for resolving a claim.
PodId                string                // Id of the pod submitting the claim
Voter                string                // Id of the guarantor verifying a claim
Vote                 bool                  // Vote of the insured pod investor
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) voteDigitalCourt( stub shim.ChaincodeStubInterface,
								   args []string) pb.Response {

	// Retrieve input information in a voting object //
	pod_id := args[0]
	
	// 	Retrieve input information in a deletion object //
	input := Voter{}
	err1 := json.Unmarshal( []byte(args[0]), &input )
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) })
	date := getTimeNow()

	// Retrieve pod //
	pod, err2 := t.retrievePodInfo(stub, input.PodId)
	if err2 != nil { shim.Error(err2.Error()) }

	// Check that the caller has the role of guarantor //
	actor, err3 := t.retrieveUserRole(stub, input.Voter)
	if err3 != nil { shim.Error(err3.Error()) }
	if actor.Role != COURTMEMBER_ROLE {
		shim.Error( "ONLY DIGITAL COURT MEMBERS CAN RESOLVE A CLAIM") }

	// Check that pod is in a digital court process //
	if pod.Claiming.Status != "COURT" {
		return shim.Error( "ERROR: THE POD IS NOT IN DIGITAL COURT.") }

	// Submit Vote //
	votes := make( map[string]bool )
	pod.Claiming.Votation.Votes[input.Voter] = true
	txn := Transfer{
		Type: "Digital_Court_Vote", Amount: 0., Token: "", 
		From: input.Voter, To: "Pod NFT " + pod.PodId,
		Id: xid.New().String(), Date: date }

	// Retrieve output //
	output := OutputNFT{ Transaction: []Transfer{ txn } }
	outputBytes, _ := json.Marshal(output)
	return shim.Success(outputBytes)
}


/* -------------------------------------------------------------------------------------------------
endDigitalCourtVotation: this function is called when a Digital Court Voting arrives to its end
                         and identified in backend.
PodId                string                // Id of the pod submitting the claim
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) endDigitalCourtVotation( stub shim.ChaincodeStubInterface,
								   		args []string) pb.Response {

	// Retrieve input information in a deletion object //
	pod_id := args[0]
	date := getTimeNow()
	update_wallets := make( map[string]MultiWallet )
	update_pods := make( map[string]POD )
	transaction := []Transfer{}

	// Retrieve pod //
	pod, err2 := t.retrievePodInfo(stub, pod_id)
	if err2 != nil { shim.Error(err2.Error()) }
	claiming_process := pod.Claiming

	// Check if the time for voting proposal is done //
	period := getTimePeriod( claiming_process.Votation.StartDate ) 
	if float64(period) < pod.Claiming.Votation.Duration {
		return shim.Error( "ERROR: THE VOTING TIME DID NOT END UP.")}

	// Check that pod is in a claiming validation process //
	if pod.State.ClaimStatus != "COURT" {
		return shim.Error( "ERROR: THE POD IS NOT IN COURT PHASE.") }
	
	// Votation veredict // 
	votation := 0.
	total_votes := 0.
	voting_yes := []string{}
	voting_no := []string{}
	for id, vote := range(pod.ClaimingProcess.Votes) {
		total_votes = total_votes + 1
		if vote {
			voting_yes = append(voting_yes, id)
			votation = votation + 1 
		} else { 
			voting_no = append(voting_no, id)  }
	}
	result := votation / total_votes

	///////////////////////// RESULT OF VOTATION ////////////////////////////

	// Case (1): claim accepted //
	if votation_pct >= CONSENSUS_DIGITAL_COURT {
		update_wallets, transactions = payoffDigitalCourt( pod, claiming_process, 
			                                               voting_yes, true) 
		// TODO: liquidate function for all insurance pools //
		pod.ClaimingProcess = ClaimingProcess{}
		pod.State.Status = "LIQUIDATED"
		update_pods[ pod.Id ] = pod
	// Case (2): claim declined //
	} else if votation_pct <= 1. - CONSENSUS_CLAIM_VALIDATION { 
		update_wallets, transactions = payoffDigitalCourt( pod, claiming_process, 
													       voting_no, flase) 
		pod.ClaimingProcess = ClaimingProcess{}
		update_pods[ pod.Id ] = pod
	} 

	// Return and update output //
	outputBytes, err4 := t.createOutput( stub,
		 update_wallets, update_pods, transactions )
    if err4 != nil { shim.Error(err4.Error()) }
	return shim.Success(outputBytes)
}

/* -------------------------------------------------------------------------------------------------
payoffVotation: this function is called when the a Digital Court case is resolved. The Claiming
				Pool is reedemed. Digital courts who are right, gets the part of the share of 
				bad actors. Good actors receive its part.
------------------------------------------------------------------------------------------------- */

func payoffDigitalCourt( pod POD, claiming_process ClaimingProcess, court []string{},
	                     decision bool ) ( map[string]MultiWallet, []Transfer ) {
	
	date := getTimeNow()
	update_wallets := make( map[string]MultiWallet )
	transactions := []Transfer{}

	// Get bad and good actors //
	bad_actors := claiming_process.VotingNothing
	good_actors := []string{}
	if decision { 
		bad_actors = append(bad_actors, claiming_process.VotingNo) 
		good_actors = claiming_process.VotingYes
	} else {
		bad_actors = append(bad_actors, claiming_process.VotingYes) 
		good_actors = claiming_process.VotingNo }

	// Voting power Guarantors // 
	voting_power_guarantors := 0.
	for _, power := range(pod.Guarantors) {
		voting_power_guarantors = voting_power_guarantors + power }
	
	// Voting power Guarantors // 
	voting_power_investors := 0.
	for _, power := range(pod.State.InsuredInvestors) {
		voting_power_investors = voting_power_investors + power }

	// Get claiming share from bad actors //
	claiming_share := 0.
	for _, actor := range(bad_actors) {
		pct = 0.
		amount, isInList := pod.Guarantors[ actor ]
		if isInList { pct = amount / voting_power_guarantors }
		else {
			amount, isInList = pod.State.InsuredInvestors[ actor ]
			pct = amount / voting_power_investors }
		claiming_share = claiming_share + pct * 0.5 }

	// Transfer funds to court members //
	reedem_court := claiming_share * pod.State.ClaimingPool 
	num := len( court )
	payoff_court := 0
	for _, id := range( court ) {
		payoff := reedem_court / num
		payoff_court = payoff_court + payoff
		user_wallet := t.retrieveUserWallet( stub, id )
		balance_wallet := user_wallet.Balances[ pod.Token ]
		balance_wallet.Amount = balance_wallet.Amount + payoff
		txn = Transfer{
			Type: "Court_Resolution_Good", Amount: payoff, 
			Token: pod.Token, 
			From: "Claiming Pool " + pod.PodId, To: id,
			Id: xid.New().String(), Date: date }
		user_wallet.Balances[ pod.Token ] = balance_wallet
		user_wallet.Transaction = []Transfer{txn}
		transactions = append( transactions, txn )
		update_users[ id ] = user_wallet
	}

	remaining_amount := math.Max( pod.State.ClaimingPool - reedem_court, 0. )
	// Transfer share to good actors //
	for id, actor := range(good_actors) {
		pct = 0.
		amount, isInList := pod.Guarantors[ actor ]
		if isInList { pct = amount / voting_power_guarantors }
		else {
			amount, isInList = pod.State.InsuredInvestors[ actor ]
			pct = amount / voting_power_investors }
		payoff := pct * remaining_amount * 0.5
		// Get balance of user to pay //
		user_wallet := t.retrieveUserWallet( stub, id )
		balance_wallet := user_wallet.Balances[ pod.Token ]
		balance_wallet.Amount = balance_wallet.Amount + payoff
		txn = Transfer{
			Type: "Claim_Good_Vote", Amount: payoff, 
			Token: pod.Token, 
			From: "Claiming Pool " + pod.PodId, To: id,
			Id: xid.New().String(), Date: date }
		user_wallet.Balances[ pod.Token ] = balance_wallet
		user_wallet.Transaction = []Transfer{txn}
		transactions = append( transactions, txn )
		update_users[ id ] = user_wallet
	 }

	 return update_wallets, transactions
}
