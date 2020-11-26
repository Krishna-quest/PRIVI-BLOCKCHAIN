///////////////////////////////////////////////////////////////
// File containing the model structs for the POD definition,
// and blockchain inputs and outputs 
///////////////////////////////////////////////////////////////                           

package main


/*---------------------------------------------------------------------------
SMART CONTRACT MAIN MODELS
-----------------------------------------------------------------------------*/  

// Definition of the model of the smart contract //
type PodNFT struct {}

// Model of outputs //
type Output struct {
	UpdatePods              map[string]POD               `json:"UpdatePods"`
	UpdateWallets 		    map[string]MultiWallet		 `json:"UpdateWallets"`
	UpdateInsurance         map[string]InsurancePool     `json:"UpdateInsurance"`
	Transaction             []Transfer                   `json:"Transaction"`
}

// Model of output NFT //
type OutputNFT struct {
	UpdatePods              []POD                        `json:"UpdatePods"`
	UpdateWallets 		    []MultiWallet		 	     `json:"UpdateWallets"`
	Transaction             []Transfer                   `json:"Transaction"`
}

/*---------------------------------------------------------------------------
SMART CONTRACT MODELS FOR POD TOKEN EXCHANGE
-----------------------------------------------------------------------------*/ 

// Model of a Market Order on the pod //
type MarketOrder struct {
	Trader                  string                       `json:"Trader"`
	Amount                  int64                        `json:"Amount"`
	Price                   float64                      `json:"Price"`
}

// Model of an Order Book on the pod //
type OrderBook struct {
	Buy                     map[string]MarketOrder       `json:"Buy"`
	Sell                    map[string]MarketOrder       `json:"Sell"`
	FundingPool             float64 					 `json:"FundingPool"`
}

// Model to delete a Market Order //
type NewOrder struct {
	Trader                  string                       `json:"Trader"`
	PodId 				 	string 			 			 `json:"PodId"`
	Amount                  int64      				     `json:"Amount"`
	Price                   float64      				 `json:"Price"`
}

// Model to delete a Market Order //
type DeleteOrder struct {
	Trader                  string                       `json:"Trader"`
	PodId 				 	string 			 			 `json:"PodId"`
	OrderId                 string      				 `json:"OrderId"`
}

// Model to nitiate/submit a vote for claim on the Pod //
type Voter struct {
	PodId 				 	string 			 			 `json:"PodId"`
	InvestorId              string      				 `json:"InvestorId"`
	Vote                    bool                         `json:"Vote"`
}


/*---------------------------------------------------------------------------
SMART CONTRACT MODELS FOR NFT PODS 
-----------------------------------------------------------------------------*/ 


// Define model of a voting  //
type VotingProcess struct {
	StartDate               string				         `json:"StartDate"`
	Votes                   map[string]bool              `json:"Votes"`
	Duration                float64                      `json:"Duration"`   
}

// Define model of a claim //
type ClaimingProcess struct {
	Status                  string                       `json:"Status"`
	Votation                VotingProcess				 `json:"Votation"`
	VotingYes               []string                     `json:"VotingYes"`
	VotingNo                []string                     `json:"VotingNo"`
	VotingNothing           []string                     `json:"VotingNothing"`   
}

// Define model of the state of a POD //
type PODstate struct {
	InsuredInvestors        map[string]float64           `json:"InsuredInvestors"`
    AmountOffer             int64                        `json:"AmountOffer"`
	AmountDemand            int64                        `json:"AmountDemand"`
	Status                  string                       `json:"Status"`
	ClaimingPool            float64                      `json:"ClaimingPool"`
}

// Define our struct to store the conditions of the POD in the Blockchain //
type POD struct {
	PodId 				 	string 			 			 `json:"PodId"`
	Royalty  				float64   					 `json:"Royalty"`
	Creator                 string                       `json:"Creator"`
	Token                   string           			 `json:"Token"`
	Supply                  int64                        `json:"Supply"`
	Date 					string						 `json:"Date"`
	State                   PODstate                     `json:"State"`	
	VotingTime              float64                      `json:"VotingTime"`
	VerifTime               float64                      `json:"VerifTime"`
	CourtTime               float64                      `json:"CourtTime"`
	Claiming                ClaimingProcess              `json:"Claiming"`
	Guarantors              map[string]int64             `json:"Guarantors"`	
	TotalInsurance          float64                      `json:"TotalInsurance"`	
	OrderBook               OrderBook                    `json:"OrderBook"`
}

// Model to buy POD tokens from a selling market order //
type PodBuy struct {
	Trader                  string                       `json:"Trader"`
	PodId 				 	string 			 			 `json:"PodId"`
	OrderId                 string      				 `json:"OrderId"` 
	Amount                  int64                        `json:"Amount"`
}



/*---------------------------------------------------------------------------
SMART CONTRACT MODELS FOR INSURANCE POOLS
-----------------------------------------------------------------------------*/  

// Define model of a insurance pool insurer //
type Insurer struct {
	Amount                	float64                       `json:"Amount"`
	Date                    string                        `json:"Date"`
}

// Define model of a insurance pool client //
type Client struct {
	Amount                	int64                         `json:"Amount"`
	Date                    string                        `json:"Date"`
}

// Define model of the state of a Liquidity Token Pool //
type InsuranceState struct { 
	Insurers                map[string]Insurer        `json:"Insurers"`
	Clients                 map[string]Client         `json:"Clients"` 
	InsuredAmount           int64                     `json:"InsuredAmount"`   
	CoveragePool            float64                   `json:"CoveragePool"`
}

// Define instance of a Liquidity Token Pool on Blockchain //
type InsurancePool struct {
	Id                      string 			 			 `json:"Id"`
	Guarantors              string 			 	         `json:"Guarantors"`
	PodId                   string                       `json:"PodId"`
	Token                   string 			 			 `json:"Token"`
	Duration                int64                        `json:"Duration"`
	Frequency               int64                        `json:"Frequency"`
	State                   InsuranceState               `json:"State"`
	Status                  string                       `json:"Status"`
	Date                    string                       `json:"Date"`
	FeeInscription          float64                      `json:"FeeInscription"`
	FeeMembership           float64                      `json:"FeeMembership"`
	Valuation               float64                      `json:"Valuation"`
	Coverage                int64                        `json:"Coverage"`
	Deposit                 float64                      `json:"Deposit"`
}

/*---------------------------------------------------------------------------
COIN BALANCE INVOKATIONS 
-----------------------------------------------------------------------------*/  

// Definition of the user Balance for a given token //
type Balance struct {
	Token                	string    					`json:"Token"`
	Type                    string                      `json:"Type"`
	Amount 		         	float64   					`json:"Amount"`
	Credit_Amount        	float64   					`json:"Credit_Amount"`
	Staking_Amount       	float64   					`json:"Staking_Amount"`
	Borrowing_Amount    	float64   					`json:"Borrowing_Amount"`
	PRIVI_lending           float64                     `json:"PRIVI_lending"`
	PRIVIcreditLend     	map[string]bool             `json:"PRIVIcreditLend"`
	PRIVI_borrowing         float64                     `json:"PRIVI_borrowing"`
	PRIVIcreditBorrow     	map[string]bool             `json:"PRIVIcreditBorrow"`
	Collateral_Amount    	float64   					`json:"Collateral_Amount"`
	Collaterals          	map[string]float64	    	`json:"Collaterals"`
} 

// Definition of the user Balance for FT tokens //
type BalanceFT struct {
	Token                	string    					`json:"Token"`
	Type                    string                      `json:"Type"`
	Amount 		         	float64   					`json:"Amount"`
	PRIVI_lending           float64                     `json:"PRIVI_lending"`
	PRIVIcreditLend     	map[string]float64          `json:"PRIVIcreditLend"`
	PRIVI_borrowing         float64                     `json:"PRIVI_borrowing"`
	PRIVIcreditBorrow     	map[string]float64          `json:"PRIVIcreditBorrow"`
	Collateral_Amount    	float64   					`json:"Collateral_Amount"`
	Collaterals          	map[string]float64	    	`json:"Collaterals"`
} 

// Definition of the user Balance for NFT tokens //
type BalanceNFT struct {
	Token                	string    					`json:"Token"`
	Type                    string                      `json:"Type"`
	Amount 		         	int64    					`json:"Amount"`
	PRIVI_lending           int64                       `json:"PRIVI_lending"`
	PRIVIcreditLend     	map[string]int64            `json:"PRIVIcreditLend"`
	PRIVI_borrowing         int64                       `json:"PRIVI_borrowing"`
	PRIVIcreditBorrow     	map[string]int64            `json:"PRIVIcreditBorrow"`
	Collateral_Amount    	int64   					`json:"Collateral_Amount"`
	Collaterals          	map[string]float64	    	`json:"Collaterals"`
} 

// Definition of a multi wallet containing all balances //
type MultiWallet struct {
	PublicId          		string                		`json:"PublicId"`
	Balances          		map[string]Balance    		`json:"Balances"`
	BalancesFT          	map[string]BalanceFT    	`json:"BalancesFT"`
	BalancesNFT          	map[string]BalanceNFT    	`json:"BalancesNFT"`
	TrustScore              float64                     `json:"TrustScore"`
	EndorsementScore        float64                     `json:"EndorsementScore"`
	Transaction             []Transfer            		`json:"Transaction"`
}

// Definition of a Token Transfer //
type Transfer struct {
	Type      				 string  					`json:"Type"`
	Token      				 string  					`json:"Token"`
	From      				 string  					`json:"From"`
	To        				 string  					`json:"To"`
	Amount     				 float64	  			    `json:"Amount"`
	Id         				 string       				`json:"Id"`
	Date                     string                     `json:"Date"`
} 

/*---------------------------------------------------------------------------
-----------------------------------------------------------------------------*/  
