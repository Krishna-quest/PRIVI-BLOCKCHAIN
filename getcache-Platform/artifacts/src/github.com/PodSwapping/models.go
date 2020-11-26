///////////////////////////////////////////////////////////////
// File containing the model structs for the POD definition,
// and blockchain inputs and outputs 
///////////////////////////////////////////////////////////////                           

package main


/*---------------------------------------------------------------------------
SMART CONTRACT MODELS FOR LIQUIDITY POOLS
-----------------------------------------------------------------------------*/  

// Definition of the model of the smart contract //
type PodSwappingSmartContract struct {}

// Define instance of a Pod Swap on Blockchain //
type Output struct { 
	Liquidated             bool                           `json:"Liquidated"`
	UpdatePods             map[string]POD                 `json:"UpdatePods"`
	UpdatePools		       map[string]LiquidityTokenPool  `json:"UpdatePools"`
	UpdateWallets          map[string]MultiWallet         `json:"UpdateWallets"`
}


// Define model of a Liquidity Provider //
type Provider struct {
	Amount 					float64   					 `json:"Amount"`
	WithdrawDay             int64                        `json:"WithdrawDay"`
}

// Define model of the state of a Liquidity Token Pool //
type PoolState struct {
	LiqProviders            map[string]Provider          `json:"LiqProviders"` 
	LiqProvidersNum         int64                        `json:"LiqProvidersNum"`
	Reserve 				float64   					 `json:"Reserve"`
	Deposited               float64                      `json:"Deposited"`
	Liq_Tokens              map[string]float64           `json:"Liq_Tokens"`
	ReserveRatio            float64                      `json:"ReserveRatio"`
}

// Define instance of a Liquidity Token Pool on Blockchain //
type LiquidityTokenPool struct {
	Id                      string 			 			 `json:"LiqPoolId"`
	CreatorId               string 			 			 `json:"CreatorId"`
	Token                   string 			 			 `json:"Token"`
	MinReserveRatio         float64                      `json:"MinReserveRatio"`  
	InitialAmount           float64                      `json:"InitialAmount"`
	State                   PoolState                    `json:"State"`
	Date                    string                       `json:"Date"`
	Fee                     float64                      `json:"Fee"`
	WithdrawalTime          int64                        `json:"WithdrawalTime"`
	WithdrawalFee           float64                      `json:"WithdrawalFee"`
	MinTrustScore           float64                      `json:"MinTrustScore"`
	MinEndScore             float64                      `json:"MinEndScore"`
}

// Define instance of a Liquidity Token Deposit on Blockchain //
type LiquidityDeposit struct {
	LiqProviderId 			string 			 			 `json:"LiqProviderId"`
	LiqPoolId          	    string         			     `json:"LiqPoolId"`
	Amount          		float64         			 `json:"Amount"`
}

// Define instance of a Liquidity Token Withdrawal on Blockchain //
type LiquidityWithdrawal struct {
	LiqProviderId 			string 			 			 `json:"LiqProviderId"`
	LiqPoolId          	    string         			     `json:"LiqPoolId"`
	Amount          		float64         			 `json:"Amount"`
	RateChange              map[string]float64           `json:"RateChange"`
}

	
/*---------------------------------------------------------------------------
SMART CONTRACT MODELS FOR POD AND SWAPPING
-----------------------------------------------------------------------------*/  

// Define model parameter of risk //
type RiskParameters struct {
	Pct_supply_lower		float64                      `json:"Pct_supply_lower"`
	Pct_supply_upper		float64                      `json:"Pct_supply_upper"`
	Interest_min		    float64                      `json:"Interest_min"`
	Interest_max		    float64                      `json:"Interest_max"`
	Liquidation_min		    float64                      `json:"Liquidation_min"`
}

// Define model of the state of a POD //
type PODstate struct {
	Liq_Pools               map[string]float64           `json:"Liq_Pools"`
	Investors				map[string]float64  		 `json:"Investors"`
	InvestorNum             int64                        `json:"InvestorNum"`
	Status 					string 						 `json:"Status"`
	POD_Day 				int64   					 `json:"POD_Day"`
	Debt                    float64                      `json:"Debt"`
	MissingPayments         int64                        `json:"MissingPayments"`
	FundsRaised             float64                      `json:"FundsRaised"`
	SupplyReleased          float64                      `json:"SupplyReleased"`
}

// Define model for liquidity pools in the POD //
type PODpools struct {
	Funding_Pool 			float64						 `json:"Funding_Pool"`
	Exchange_Pool           float64                      `json:"Exchange_Pool"`
	Collateral_Pool 		map[string]float64 		     `json:"Collateral_Pool"`
	Interest_Pool           float64						 `json:"Interest_Pool"`
	POD_Token_Pool          float64                      `json:"POD_Token_Pool"`
}

// Define our struct to store the conditions of the POD in the Blockchain //
type POD struct {
	PodId 				 	string 			 			 `json:"PodId"`
	Creator                 string                       `json:"Creator"`
	Token                   string           			 `json:"Token"`
	Duration 			    int64						 `json:"Duration"`
	Payments 			    int64						 `json:"Payments"`
	Principal 			    float64						 `json:"Principal"`
	Interest 				float64 					 `json:"Interest"`
	P_liquidation 			float64 					 `json:"P_liquidation"`
	InitialSupply           float64                      `json:"InitialSupply"`
	Invariant_cte           float64                      `json:"Invariant_cte"`
	Date 					string						 `json:"Date"`
	State                   PODstate                     `json:"State"`	
	Pools                   PODpools                     `json:"Pools"`	 
}

// Define instance of a Pod Deletion on Blockchain //
type PodDeletion struct {
	PublicId 				string 			 			 `json:"PublicId"`
	PodId          	        string         			     `json:"PodId"`
} 

// Define instance of a Pod Investment on Blockchain //
type PodInvestment struct {
	InvestorId 				string 			 			 `json:"InvestorId"`
	PodId          	        string         			     `json:"PodId"`
	Amount          		float64         			 `json:"Amount"`
	RateChange              map[string]float64           `json:"RateChange"`
}

// Define instance of a Pod Swap on Blockchain //
type PodSwapping struct {
	InvestorId 				string 			 			 `json:"InvestorId"`
	LiqPoolId               string                       `json:"LiqPoolId"`
	PodId          	        string         			     `json:"PodId"`
	Type                    string                       `json:"Type"`
	Amount          		float64         			 `json:"Amount"`
	RateChange              map[string]float64           `json:"RateChange"`
}

// Define instance of a Pod Swap on Blockchain //
type PodLiquidator struct {
	PodId 				 	string 			 			 `json:"PodId"`
	RateChange              map[string]float64           `json:"RateChange"`
}

// Define instance of a Pod Swap on Blockchain //
type OutputSwapping struct {
	UpdatePool 				LiquidityTokenPool 			 `json:"UpdatePool"`
	UpdatePod               POD                          `json:"UpdatePod"`
	UpdateUser              MultiWallet                  `json:"UpdateUser"`
}

// Define instance of a Pod Swap on Blockchain //
type OutputLiquidation struct { 
	Liquidated             bool                           `json:"Liquidated"`
	Pod                    POD                            `json:"Pod"`
	UpdatePools		       []LiquidityTokenPool 		  `json:"UpdatePools"`
	UpdateUsers            []MultiWallet        		  `json:"UpdateUsers"`
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
