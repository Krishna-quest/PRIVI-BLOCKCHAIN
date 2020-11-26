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
Init:  this function register PRIVI as the Admin of the POD NFT Smart Contract. It initialises
       the lists of PODs and the Indexes used. Args: array containing a string:
PrivateKeyID           string   // Private Key of the admin of the smart contract
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) Init(stub shim.ChaincodeStubInterface) pb.Response {

	_, args := stub.GetFunctionAndParameters()
	// Store in the state of the smart contract the Private Key of Admin //
	err1 := stub.PutState(IndexAdmin, []byte(args[0]))
	if err1 != nil {
		return shim.Error("ERROR: SETTING THE ADMIN PRIVATE KEY: " +
			err1.Error()) }
			
	if args[1] == "UPGRADE" { return shim.Success(nil) }

	// Initialise list of PODs in the Smart Contract as empty //
	pod_list, _ := json.Marshal(make(map[string]bool))
	err2 := stub.PutState(IndexPodList, pod_list)
	if err2 != nil {
		return shim.Error("ERROR: INITIALISING THE POD LIST: " +
			err2.Error()) }

	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
 The Invoke method is called as a result of an application request to run the Smart Contract ""
 The calling application program has also specified the particular smart contract function to be called
-------------------------------------------------------------------------------------------------*/

func (t *PodNFT) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// Retrieve function and arguments //
	function, args := stub.GetFunctionAndParameters()

	// Retrieve caller of the function //
	caller, err1 := CallerCN(stub)
	if err1 != nil {
		return shim.Error("ERROR: GETTING THE CALLER OF THE TRANSFER " +
			"FUNCTION. " + err1.Error()) }
	// Retrieve admin of the smart contract //
	adminBytes, err2 := stub.GetState(IndexAdmin)
	if err2 != nil {
		return shim.Error("ERROR: GETTING THE ADMIN OF THE SMART " +
			"CONTRACT. " + err2.Error()) }
	admin := string(adminBytes)
	// Verify that the caller of the function is admin //
	if caller != admin {
		return shim.Error("ERROR: CALLER " + caller + " DOES NOT HAVE PERMISSION ") }

	// Call the proper function //
	switch function {
		case "retrievePodList":
			pod_list, err := t.retrievePodList(stub)
			if err != nil {
				return shim.Error(err.Error()) }
			podlistBytes, _ := json.Marshal(pod_list)
			return shim.Success(podlistBytes)
		case "retrievePodInfo":
			pod, err := t.retrievePodInfo(stub, args[0])
			if err != nil {
				return shim.Error(err.Error()) }
			podBytes, _ := json.Marshal(pod)
			return shim.Success(podBytes)
		case "updatePod":
			pod := POD{}
			json.Unmarshal( []byte(args[0]), &pod )
			podBytes, err := t.updatePod(stub, pod)
			if err != nil { return shim.Error(err.Error()) }
			return shim.Success( podBytes )
		
		// EXCHANGE NFT FUNCTIONS //
		case "initiatePodNFT":
			return t.initiatePodNFT(stub, args)
		case "newBuyOrder":
			return t.newBuyOrder(stub, args)
		case "deleteBuyOrder":
			return t.deleteBuyOrder(stub, args)
		case "newSellOrder":
			return t.newSellOrder(stub, args)
		case "deleteSellOrder":
			return t.deleteSellOrder(stub, args)
		case "getOrderBook":
			return t.getOrderBook(stub, args)
		case "buyPodNFT":
			return t.buyPodNFT(stub, args)
		case "sellPodNFT":
			return t.sellPodNFT(stub, args)

		// CLAIMING FUNCTIONS //
		case "initiateClaimProposal":
			return t.initiateClaimProposal(stub, args)
		case "voteClaimProposal":
			return t.voteClaimProposal(stub, args)
		case "endProposalVoting":
			return t.endProposalVoting(stub, args)
		case "voteClaimValidation":
			return t.voteClaimValidation(stub, args)
		case "endValidationVotation":
			return t.endValidationVotation(stub, args)
		case "voteDigitalCourt":
			return t.voteDigitalCourt(stub, args)


	}
	return shim.Error("Incorrect function name: " + function)
}

/* -------------------------------------------------------------------------------------------------
initiatePODNFT: this function initialises a new NFT POD with the parameters described below.
                Args is an array containing a json with two fields:
Creator              string                // Id of the creator of the POD
Royalty              float64               // Percentage of royalty of the POD
Token                string  			   // Token accepted to sell a pod token
Offers               map[int64]float64     // List of initial selling offers of holder (args[1])
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) initiatePodNFT( stub shim.ChaincodeStubInterface,
								 args []string) pb.Response {

	// Retrieve the input information for POD Initialisation //
	pod := POD{}
	err1 := json.Unmarshal([]byte(args[0]), &pod)
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) }
	offers := make(map[int64]float64)
	json.Unmarshal([]byte(args[1]), &offers)
	date := getTimeNow()

	// Retrieve pod list //
	pod_list, err2 := t.retrievePodList(stub)
	if err2 != nil { shim.Error(err2.Error()) }
	pod.PodId = xid.New().String() + xid.New().String()

	// Generate initial Order Book //
	sell := make(map[string]MarketOrder)
	amount_offer := int64(0)
	for amount, price := range offers {
		id := xid.New().String()
		amount_offer = amount_offer + amount
		sell[id] = MarketOrder{
			Amount: amount, Price: price, Trader: pod.Creator } }

	orderBook := OrderBook{
		Buy: make(map[string]MarketOrder), Sell: sell, FundingPool: 0. }

	// Generate initial state for the NFT POD //
	state := PODstate{
		InsuredInvestors: make(map[string]float64),
		Status: "ACTIVE",
		AmountOffer: amount_offer,
		AmountDemand: 0 }

	// Create new pod and store in Blockchain and add pod to PodList //
	pod.Date = date
	pod.State = state
	pod.Supply = amount_offer
	pod.Guarantors = make(map[string]bool)
	pod.TotalInsurance = 0.
	pod.Claiming = ClaimingProcess{ Status: "NONE" }
	pod.OrderBook = orderBook
	pod_list[pod.PodId] = true

	// Store new pod on Blockchain //
	_, err3 := t.updatePod(stub, pod)
	if err3 != nil { return shim.Error(err3.Error()) }
	err4 := t.updatePodList(stub, pod_list)
	if err4 != nil { return shim.Error(err4.Error()) }

	// Update pod creator wallet //
	wallet_creator, err5 := t.retrieveUserWallet(stub, pod.Creator)
	if err5 != nil {
		shim.Error(err5.Error()) }
	wallet_creator = t.createPodBalance( stub, wallet_creator, pod.PodId,
		                			     amount_offer )
	err6 := t.updateUserWallet(stub, wallet_creator)
	if err6 != nil {
		return shim.Error(err6.Error()) }

	txn := Transfer{
		Type: "NFT_Pod_Creation", Amount: float64(amount_offer), 
		Token: pod.PodId, 
		From: "Pod NFT " + pod.PodId, To: pod.Creator,
		Id: xid.New().String(), Date: date }

	// Output of the result //
	output := OutputNFT{
		UpdatePods:    []POD{pod},
		UpdateWallets: []MultiWallet{wallet_creator},
	    Transaction:   []Transfer{txn} }
	outputBytes, _ := json.Marshal(output)
	return shim.Success(outputBytes)
}

/* -------------------------------------------------------------------------------------------------
getOrderBook: get the order book of a given pod. Args is an array containing a json with:
PodId                string                // Id of the POD
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) getOrderBook( stub shim.ChaincodeStubInterface,
							   args []string) pb.Response {
	// Retrieve pod //
	pod, err1 := t.retrievePodInfo(stub, args[0])
	if err1 != nil { shim.Error(err1.Error()) }

	outputBytes, _ := json.Marshal( pod.OrderBook )
	return shim.Success(outputBytes)
}

/* -------------------------------------------------------------------------------------------------
newBuyOrder: this function submits a buying market order in the order book of NFT POD tokens.
             Args is an array containing a json with two fields:
PodId                string                // Id of the POD
Trader               string  			   // Id of the creator of the market order
Amount               int64                 // Amount of pod tokens willing to buy
Price                float64               // Buying price per POD token
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) newBuyOrder( stub shim.ChaincodeStubInterface,
							  args []string) pb.Response {

	// Retrieve input information in a deletion object //
	order := NewOrder{}
	err1 := json.Unmarshal([]byte(args[0]), &order)
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
				            err1.Error()) }
	order_id := xid.New().String() 
	date := getTimeNow()
	transaction := []Transfer{}

	// Retrieve pod list //
	pod_list, err2 := t.retrievePodList(stub)
	if err2 != nil { shim.Error(err2.Error()) }
	_, isInList := pod_list[order.PodId]
	if !isInList {
		return shim.Error("ERROR: POD ID IS NOT REGISTERED") }

	// Retrieve pod //
	pod, err3 := t.retrievePodInfo(stub, order.PodId)
	if err3 != nil { shim.Error(err3.Error()) }

	// Retrieve trader wallet //
	wallet_trader, err4 := t.retrieveUserWallet(stub, order.Trader)
	if err4 != nil { shim.Error(err4.Error()) }
	balance_trader_token := wallet_trader.Balances[pod.Token]

	// Check that the trader has enough funds to buy it //
	amount := order.Price * float64(order.Amount)
	if balance_trader_token.Amount < amount {
		return shim.Error( "ERROR: TRADER " + order.Trader + "HAS NOT " +
						   "ENOUGH FUNDS TO SUBMIT BUYING MARKET ORDER." ) }
	balance_trader_token.Amount = balance_trader_token.Amount - amount
	order_transfer := Transfer{
		Type: "NFT_new_buy_order", Amount: amount, Token: pod.Token,
		From: order.Trader, To: "Exchange NFT " + pod.PodId, 
		Id: xid.New().String(), Date: date }
	transaction = append(transaction, order_transfer)
	wallet_trader.Balances[pod.Token] = balance_trader_token
	wallet_trader.Transaction = []Transfer{order_transfer}

	// Input new buying order on the Order Book of the Pod //
	new_order := MarketOrder{
		Trader: order.Trader, Amount: order.Amount, Price: order.Price}
	pod.OrderBook.Buy[order_id] = new_order
	pod.State.AmountDemand = pod.State.AmountDemand + order.Amount
	pod.OrderBook.FundingPool = pod.OrderBook.FundingPool + amount

	// Update Pod State on Blockchain //
	_, err5 := t.updatePod(stub, pod)
	if err5 != nil { return shim.Error(err5.Error()) }

	// Update MultiWallet of trader //
	err6 := t.updateUserWallet(stub, wallet_trader)
	if err6 != nil { return shim.Error(err6.Error()) }

	// Output of the result //
	output := OutputNFT{
		UpdatePods:    []POD{pod},
		UpdateWallets: []MultiWallet{wallet_trader},
	    Transaction: transaction }
	outputBytes, _ := json.Marshal(output)
	return shim.Success(outputBytes)
}


/* -------------------------------------------------------------------------------------------------
deleteBuyOrder: this function deletes a submitted buying market order from the order book.
                Args is an array containing a json with two fields:
PodId                string                // Id of the POD
OrderId              string                // Delete market order
Trader               string  			   // Id of the creator of the market order
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) deleteBuyOrder( stub shim.ChaincodeStubInterface,
								 args []string) pb.Response {

	// Retrieve input information in a deletion object //
	delete_order := DeleteOrder{}
	err1 := json.Unmarshal([]byte(args[0]), &delete_order)
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) }
	order_id := delete_order.OrderId
	date := getTimeNow()
	transaction := []Transfer{}

	// Retrieve pod list //
	pod_list, err2 := t.retrievePodList(stub)
	if err2 != nil { shim.Error(err2.Error()) }
	_, isInList := pod_list[delete_order.PodId]
	if !isInList {
		return shim.Error("ERROR: POD ID IS NOT REGISTERED") }

	// Retrieve pod //
	pod, err3 := t.retrievePodInfo(stub, delete_order.PodId)
	if err3 != nil { shim.Error(err3.Error()) }

	// Check that order is in Pod Order Book //
	market_order, isInList := pod.OrderBook.Buy[order_id]
	if !isInList {
		return shim.Error("ERROR: THE MARKET ORDER " + order_id +
			" IS NOT IN THE ORDER BOOK OF THE POD.") }

	// Check that the caller of the function is the creator of the Order //
	if market_order.Trader != delete_order.Trader {
		return shim.Error(	"ERROR: THE CALLER " + delete_order.Trader +
							" IS NOT THE CREATOR OF THE ORDER." ) }
	wallet_trader, err4 := t.retrieveUserWallet(stub, delete_order.Trader)
	if err4 != nil { shim.Error(err4.Error()) }
	balance_trader_token := wallet_trader.Balances[pod.Token]

	// Give back the order amount to investor and delete order //
	amount := market_order.Price * float64(market_order.Amount)
	balance_trader_token.Amount = balance_trader_token.Amount + amount
	delete_transfer := Transfer{
		Type: "NFT_delete_buy_order", Amount: amount, Token: pod.Token,
		From: "Exchange NFT " + pod.PodId, To: delete_order.Trader,
		Id: xid.New().String(), Date: date }
	transaction = append(transaction, delete_transfer)
	wallet_trader.Balances[pod.Token] = balance_trader_token
	wallet_trader.Transaction = []Transfer{delete_transfer}
	pod.State.AmountDemand = pod.State.AmountDemand - market_order.Amount
	pod.OrderBook.FundingPool = pod.OrderBook.FundingPool - amount
	delete(pod.OrderBook.Buy, order_id)

	// Update Pod State on Blockchain //
	_, err5 := t.updatePod(stub, pod)
	if err5 != nil { return shim.Error(err5.Error()) }

	// Update MultiWallet of trader //
	err6 := t.updateUserWallet(stub, wallet_trader)
	if err6 != nil { return shim.Error(err6.Error()) }

	// Output of the result //
	output := OutputNFT{
		UpdatePods:    []POD{pod},
		UpdateWallets: []MultiWallet{wallet_trader},
	    Transaction: transaction }
	outputBytes, _ := json.Marshal(output)
	return shim.Success(outputBytes)
}

/* -------------------------------------------------------------------------------------------------
newSellOrder: this function submits a selling market, order in the order book of NFT POD tokens.
              Args is an array containing a json with two fields:
PodId                string                // Id of the POD
Trader               string  			   // Id of the creator of the market order
Amount               int64                 // Amount of pod tokens willing to sell
Price                float64               // Selling price per POD token
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) newSellOrder( stub shim.ChaincodeStubInterface,
							   args []string) pb.Response {

	// Retrieve input information in a deletion object //
	order := NewOrder{}
	err1 := json.Unmarshal([]byte(args[0]), &order)
	if err1 != nil {
		return shim.Error("ERROR: GETTING THE INPUT INFORMATION. " +
			err1.Error()) }
	order_id := xid.New().String()
	date := getTimeNow()
	transaction := []Transfer{}

	// Retrieve pod list //
	pod_list, err2 := t.retrievePodList(stub)
	if err2 != nil { shim.Error(err2.Error()) }
	_, isInList := pod_list[order.PodId]
	if !isInList {
		return shim.Error("ERROR: POD ID IS NOT REGISTERED") }

	// Retrieve pod //
	pod, err3 := t.retrievePodInfo(stub, order.PodId)
	if err3 != nil { shim.Error(err3.Error()) }

	// Retrieve trader wallet //
	wallet_trader, err4 := t.retrieveUserWallet(stub, order.Trader)
	if err4 != nil { shim.Error(err4.Error()) }
	balance_trader_pod := wallet_trader.BalancesNFT[pod.PodId]

	// Check that the trader has enough funds to buy it //
	if balance_trader_pod.Amount < order.Amount {
		return shim.Error( "ERROR: TRADER " + order.Trader + "HAS NOT " +
			               "ENOUGH FUNDS TO SUBMIT SELLING MARKET ORDER.") }
	amount_balance := balance_trader_pod.Amount - order.Amount
	balance_trader_pod.Amount = amount_balance
	order_transfer := Transfer{
		Type: "NFT_new_sell_order", Amount: float64(order.Amount),
		Token: pod.Token, From: order.Trader, To: "Exchange NFT " + pod.PodId,
		Id: xid.New().String(), Date: date }
	transaction = append(transaction, order_transfer)
	wallet_trader.BalancesNFT[pod.PodId] = balance_trader_pod
	wallet_trader.Transaction = []Transfer{order_transfer}

	// Input new selling order on the Order Book of the Pod //
	new_order := MarketOrder{
		Trader: order.Trader, Amount: order.Amount, Price: order.Price }
	pod.OrderBook.Sell[order_id] = new_order
	pod.State.AmountOffer = pod.State.AmountOffer + order.Amount

	// Update Pod State on Blockchain //
	_, err5 := t.updatePod(stub, pod)
	if err5 != nil { return shim.Error(err5.Error()) }

	// Update MultiWallet of trader //
	err6 := t.updateUserWallet(stub, wallet_trader)
	if err6 != nil { return shim.Error(err6.Error()) }

	// Output of the result //
	output := OutputNFT{
		UpdatePods:    []POD{pod},
		UpdateWallets: []MultiWallet{wallet_trader},
	    Transaction: transaction }
	outputBytes, _ := json.Marshal(output)
	return shim.Success(outputBytes)
}

/* -------------------------------------------------------------------------------------------------
deleteSellOrder: this function deletes a submitted buying market order from the order book.
                 Args is an array containing a json with two fields:
PodId                string                // Id of the POD
OrderId              string                // Delete market order
Trader               string  			   // Id of the creator of the market order
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) deleteSellOrder( stub shim.ChaincodeStubInterface,
								  args []string) pb.Response {

	// Retrieve input information in a deletion object //
	delete_order := DeleteOrder{}
	err1 := json.Unmarshal([]byte(args[0]), &delete_order)
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) }
	order_id := delete_order.OrderId
	date := getTimeNow()
	transaction := []Transfer{}

	// Retrieve pod list //
	pod_list, err2 := t.retrievePodList(stub)
	if err2 != nil { shim.Error(err2.Error()) }
	_, isInList := pod_list[delete_order.PodId]
	if !isInList {
		return shim.Error("ERROR: POD ID IS NOT REGISTERED") }

	// Retrieve pod //
	pod, err3 := t.retrievePodInfo(stub, delete_order.PodId)
	if err3 != nil { shim.Error(err3.Error()) }

	// Check that order is in Pod Order Book //
	market_order, isInList := pod.OrderBook.Sell[order_id]
	if !isInList {
		return shim.Error( "ERROR: THE MARKET ORDER " + order_id +
						   " IS NOT IN THE ORDER BOOK OF THE POD.") }

	// Check that the caller of the function is the creator of the Order //
	if market_order.Trader != delete_order.Trader {
		return shim.Error( "ERROR: THE CALLER " + delete_order.Trader +
						   " IS NOT THE CREATOR OF THE ORDER.") }
	wallet_trader, err4 := t.retrieveUserWallet(stub, delete_order.Trader)
	if err4 != nil { shim.Error(err4.Error()) }
	balance_trader_pod := wallet_trader.BalancesNFT[pod.PodId]

	// Give back the order amount to investor and delete order //
	amount_balance := balance_trader_pod.Amount + market_order.Amount
	balance_trader_pod.Amount = amount_balance
	delete_transfer := Transfer{
		Type: "NFT_delete_sell_order", Amount: float64(market_order.Amount),
		Token: pod.PodId, From: "Exchange NFT " + pod.PodId,
		To: delete_order.Trader, Id: xid.New().String(), Date: date }
	transaction = append(transaction, delete_transfer)
	wallet_trader.BalancesNFT[pod.Token] = balance_trader_pod
	wallet_trader.Transaction = []Transfer{delete_transfer}
	pod.State.AmountOffer = pod.State.AmountOffer - market_order.Amount
	delete(pod.OrderBook.Sell, order_id)

	// Update Pod State on Blockchain //
	_, err5 := t.updatePod(stub, pod)
	if err5 != nil { return shim.Error(err5.Error()) }

	// Update MultiWallet of trader //
	err6 := t.updateUserWallet(stub, wallet_trader)
	if err6 != nil { return shim.Error(err6.Error()) }

	// Output of the result //
	output := OutputNFT{
		UpdatePods:    []POD{pod},
		UpdateWallets: []MultiWallet{wallet_trader},
	    Transaction:   transaction }
	outputBytes, _ := json.Marshal(output)
	return shim.Success(outputBytes)
}

/* -------------------------------------------------------------------------------------------------
buyPodNFT: this function applies for buying a pod token from a given selling offer in the Pod
           Orderbook.  Args is an array containing a json with two fields:
PodId                string                // Id of the POD
OrderId              string                // Selling Order Id to buy pod token
Trader               string  			   // Id of the creator of the market order
Amount               int64                 // Quantity of pod tokens to buy
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) buyPodNFT( stub shim.ChaincodeStubInterface,
							args []string) pb.Response {

	// Retrieve input information in a deletion object //
	input := PodBuy{}
	err1 := json.Unmarshal([]byte(args[0]), &input)
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) }
	date := getTimeNow()
	transaction := []Transfer{}

	// Retrieve pod list //
	pod_list, err2 := t.retrievePodList(stub)
	if err2 != nil { shim.Error(err2.Error()) }
	_, isInList := pod_list[input.PodId]
	if !isInList {
		return shim.Error("ERROR: POD ID IS NOT REGISTERED") }

	// Retrieve pod //
	pod, err3 := t.retrievePodInfo(stub, input.PodId)
	if err3 != nil { shim.Error(err3.Error()) }

	// Check that selling market order exists and has enough tokens //
	order, isInList := pod.OrderBook.Sell[input.OrderId]
	if !isInList {
		return shim.Error( "ERROR: ORDER ID " + input.OrderId +
			              " NOT REGISTERED IN THE ORDER BOOK OF THE POD.") }
	if order.Amount < input.Amount {
		return shim.Error( "ERROR: THE ORDER ID " + input.OrderId +
			               " DOES NOT OFFER ENOUGH TOKENS.") }

	// Retrieve wallets from seller and buyer//
	wallet_buyer, err4 := t.retrieveUserWallet(stub, input.Trader)
	if err4 != nil { shim.Error(err4.Error()) }
	balance_buyer_pod, isInList := wallet_buyer.BalancesNFT[pod.PodId]
	if !isInList {
		wallet_buyer = t.createPodBalance( stub, wallet_buyer, pod.PodId, 0. )
		balance_buyer_pod = wallet_buyer.BalancesNFT[pod.PodId] }
	balance_buyer_token := wallet_buyer.Balances[pod.Token]

	wallet_seller, err5 := t.retrieveUserWallet(stub, order.Trader)
	if err5 != nil { shim.Error(err5.Error()) }
	balance_seller_token := wallet_seller.Balances[pod.Token]

	// Check that buyer holds funds to buy tokens and transfer them //
	amount_buy := float64(input.Amount) * order.Price
	if amount_buy > balance_buyer_token.Amount {
		return shim.Error( "ERROR: USER " + input.Trader + 
						   " DOES NOT HOLD FUNDS TO BUY TOKENS.") }
	balance_buyer_token.Amount = balance_buyer_token.Amount - amount_buy
	balance_buyer_pod.Amount = balance_buyer_pod.Amount + input.Amount
	balance_seller_token.Amount = balance_seller_token.Amount + amount_buy
	order.Amount = order.Amount - input.Amount
	wallet_buyer.BalancesNFT[pod.PodId] = balance_buyer_pod

	wallet_buyer.Balances[pod.Token] = balance_buyer_token
	wallet_seller.Balances[pod.Token] = balance_seller_token
	
	// Update pod information //
	pod.State.AmountOffer = pod.State.AmountOffer - input.Amount
	pod.OrderBook.Sell[input.OrderId] = order
	if order.Amount == 0 { 
		delete(pod.OrderBook.Sell, input.OrderId) } 
	_, err6 := t.updatePod(stub, pod)
	if err6 != nil { return shim.Error(err6.Error()) }

	// Compute tranfers and add to wallet history //
	transfer_buyer_get := Transfer{
		Type: "NFT_buy_token_get", Amount: float64(input.Amount),
		Token: pod.PodId, From: "Exchange NFT " + pod.PodId,
		To: order.Trader, Id: xid.New().String(), Date: date }
	transaction = append(transaction, transfer_buyer_get)
	transfer_buyer_pay := Transfer{
		Type: "NFT_buy_token_give", Amount: amount_buy, Token: pod.Token, 
		From: order.Trader, To: "Exchange NFT " + pod.PodId,
		Id: xid.New().String(), Date: date }
	transaction = append(transaction, transfer_buyer_pay)
	transfer_seller_get := Transfer{
		Type: "NFT_sell_token_get", Amount: amount_buy, Token: pod.Token, 
		From: "Exchange NFT " + pod.PodId, To: order.Trader,
		Id: xid.New().String(), Date: date }
	transaction = append(transaction, transfer_seller_get)

	// Update buyer wallet //
	wallet_buyer.Transaction = []Transfer{ 
						 transfer_buyer_get, transfer_buyer_pay }
	err7 := t.updateUserWallet(stub, wallet_buyer)
	if err7 != nil { return shim.Error(err7.Error()) }

	// Update seller wallet //
	wallet_seller.Transaction = []Transfer{ transfer_seller_get }
	err8 := t.updateUserWallet(stub, wallet_seller)
	if err8 != nil { return shim.Error(err8.Error()) }

	// Output of the result //
	output := OutputNFT{
		UpdatePods:    []POD{pod},
		UpdateWallets: []MultiWallet{wallet_buyer, wallet_seller},
	    Transaction: transaction }
	outputBytes, _ := json.Marshal(output)
	return shim.Success(outputBytes)
}

/* -------------------------------------------------------------------------------------------------
sellPodNFT: this function applies for buying a pod token from a given selling offer in the Pod
            Orderbook.  Args is an array containing a json with two fields:
PodId                string                // Id of the POD
OrderId              string                // Buying Order Id to sell pod token
Trader               string  			   // Id of the creator of the market order
Amount               int64                 // Quantity of pod tokens to sell
------------------------------------------------------------------------------------------------- */

func (t *PodNFT) sellPodNFT( stub shim.ChaincodeStubInterface,
							 args []string) pb.Response {

	// Retrieve input information in a deletion object //
	input := PodBuy{}
	err1 := json.Unmarshal([]byte(args[0]), &input)
	if err1 != nil {
		return shim.Error( "ERROR: GETTING THE INPUT INFORMATION. " +
							err1.Error()) }
	date := getTimeNow()
	transaction := []Transfer{}

	// Retrieve pod list //
	pod_list, err2 := t.retrievePodList(stub)
	if err2 != nil { shim.Error(err2.Error()) }
	_, isInList := pod_list[input.PodId]
	if !isInList {
		return shim.Error("ERROR: POD ID IS NOT REGISTERED") }

	// Retrieve pod //
	pod, err3 := t.retrievePodInfo(stub, input.PodId)
	if err3 != nil { shim.Error(err3.Error()) }

	// Check that selling market order exists and has enough tokens //
	order, isInList := pod.OrderBook.Buy[input.OrderId]
	if !isInList {
		return shim.Error( "ERROR: THE ORDER ID " + input.OrderId +
			               " NOT REGISTERED IN THE ORDER BOOK OF THE POD.") }
	if order.Amount < input.Amount {
		return shim.Error( "ERROR: THE ORDER ID " + input.OrderId +
			               " DOES NOT OFFER ENOUGH TOKENS.") }

	// Retrieve wallets from seller and buyer //
	wallet_buyer, err4 := t.retrieveUserWallet(stub, order.Trader)
	if err4 != nil { shim.Error(err4.Error()) }
	balance_buyer_pod, isInList := wallet_buyer.BalancesNFT[pod.PodId]
	if !isInList {
		wallet_buyer = t.createPodBalance( stub, wallet_buyer, pod.PodId, 0. )
		balance_buyer_pod = wallet_buyer.BalancesNFT[pod.PodId] }
	
	wallet_seller, err5 := t.retrieveUserWallet(stub, input.Trader)
	if err5 != nil { shim.Error(err5.Error()) }
	balance_seller_token := wallet_seller.Balances[pod.Token]
	balance_seller_pod, isInList := wallet_seller.BalancesNFT[pod.PodId]
	if !isInList {
		wallet_seller = t.createPodBalance( stub, wallet_seller, pod.PodId, 0. )
		balance_seller_pod = wallet_seller.BalancesNFT[pod.PodId] }

	// Check that seller holds enough pod tokens to sell //
	amount_sell := float64(input.Amount) * order.Price
	if balance_seller_pod.Amount < input.Amount {
		return shim.Error( "ERROR: SELLER DOES NOT HOLD ENOUGH" +
		                   " POD TOKENS TO SELL.") }
	balance_buyer_pod.Amount = balance_buyer_pod.Amount + input.Amount
	balance_seller_pod.Amount = balance_seller_pod.Amount - input.Amount
	balance_seller_token.Amount = balance_seller_token.Amount + amount_sell
	order.Amount = order.Amount - input.Amount
	pod.OrderBook.FundingPool = pod.OrderBook.FundingPool - amount_sell

	wallet_buyer.BalancesNFT[pod.PodId] = balance_buyer_pod
	wallet_seller.BalancesNFT[pod.PodId] = balance_seller_pod
	wallet_seller.Balances[pod.Token] = balance_seller_token
	
	// Update pod information //
	pod.State.AmountDemand = pod.State.AmountDemand - input.Amount
	pod.OrderBook.Buy[input.OrderId] = order
	if order.Amount == 0 { 
		delete(pod.OrderBook.Sell, input.OrderId) } 
	_, err6 := t.updatePod(stub, pod)
	if err6 != nil { return shim.Error(err6.Error()) }

	// Compute tranfers and add to wallet history //
	transfer_seller_get := Transfer{
		Type: "NFT_sell_token_get", Amount: amount_sell, Token: pod.Token,
		From: "Exchange NFT " + pod.PodId, To: input.Trader,
		Id: xid.New().String(), Date: date }
	transaction = append(transaction, transfer_seller_get)
	transfer_seller_pay := Transfer{
		Type: "NFT_sell_token_give", Amount: float64(input.Amount), Token: pod.PodId, 
		From: input.Trader, To: "Exchange NFT " + pod.PodId,
		Id: xid.New().String(), Date: date }
	transaction = append(transaction, transfer_seller_pay)
	transfer_buyer_get := Transfer{
		Type: "NFT_buy_token_get", Amount: float64(input.Amount), Token: pod.PodId, 
		From: "Exchange NFT " + pod.PodId, To: order.Trader,
		Id: xid.New().String(), Date: date }
	transaction = append(transaction, transfer_buyer_get)

	// Update buyer wallet //
	wallet_buyer.Transaction = []Transfer{ transfer_buyer_get}
	err7 := t.updateUserWallet(stub, wallet_buyer)
	if err7 != nil { return shim.Error(err7.Error()) }

	// Update seller wallet //
	wallet_seller.Transaction = []Transfer{ transfer_seller_pay, 
		                                    transfer_seller_get }
	err8 := t.updateUserWallet(stub, wallet_seller)
	if err8 != nil { return shim.Error(err8.Error()) }

	// Output of the result //
	output := OutputNFT{
		UpdatePods: []POD{pod},
		UpdateWallets: []MultiWallet{wallet_buyer, wallet_seller},
	    Transaction: transaction }
	outputBytes, _ := json.Marshal(output)
	return shim.Success(outputBytes)
}




/* -------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------- */

func main() {
	err := shim.Start(&PodNFT{})
	if err != nil {
		fmt.Errorf("Error starting Pod Swapping chaincode: %s", err)
	}
}
