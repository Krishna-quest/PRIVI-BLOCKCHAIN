///////////////////////////////////////////////////////////////
// File containing the model structs for the Cache Coin Token,
// and blockchain inputs and outputs 
///////////////////////////////////////////////////////////////                           

package main

/*---------------------------------------------------------------------------
SMART CONTRACT MODELS
-----------------------------------------------------------------------------*/  

// Definition of the model of the smart contract //
type TraditionalLendingSmartContract struct {}

// Definition the output for the smart contract //
type Output struct {
	Liquidated              string                      `json:"Liquidated"`
	UpdateWallets           map[string]MultiWallet      `json:"UpdateWallets"`
}

// Definition of Risk Parameters in Blockchain //
type Risks struct {
	Col_required      		float64  				    `json:"Col_required"`
	Col_withdrawal     		float64  				    `json:"Col_withdrawal"`
	Col_liquidation    		float64 	   			    `json:"Col_liquidation"`
	P_limit            		float64 	  			    `json:"P_limit"` 
}

// Definition of Lending Token Pools //
type LendingPool struct {
	Token            		string      				`json:"Token"`
	Reserve          		float64     				`json:"Reserve"`
	Loaned            		float64    					`json:"Loaned"`
	Staked            	    float64    			    	`json:"Staked"`
	Collateral      	    float64    			    	`json:"Collateral"`
	RiskParameters   		Risks      					`json:"RiskParameters"`
}

// Definition of Lending Request Amount //
type RequestLoan struct {
	PublicId          		string    		   			`json:"PublicId"`
	Token             		string    		   			`json:"Token"`
	Amount            		float64             		`json:"Amount"`
	Collaterals      		map[string]float64   		`json:"Collaterals"`
	RateChange       		map[string]float64   		`json:"RateChange"`
}

// Definition of Stake Toquen Request //
type StakeToken struct {
	PublicId         		string    		  			`json:"PublicId"`
	Token            		string    		   			`json:"Token"`
	Amount            		float64             		`json:"Amount"`
}

// Definition of Interest Rates as input for payInterests //
type InterestRates struct {
	LendingInterest     	map[string]float64   		`json:"LendingInterest"`
	StakingInterest    		map[string]float64    		`json:"StakingInterest"`
	RateChange          	map[string]float64   		`json:"RateChange"`
}

// Definition of a Collateral Deposit request //
type CollateralDeposit struct {
	PublicId         		string    		   			`json:"PublicId"`
	Token            		string    		  		 	`json:"Token"`
	Collaterals      		map[string]float64   		`json:"Collaterals"`
}

// Definition of a Collateral Withdrawal request //
type CollateralWithdrawal struct {
	PublicId         		string    		   			`json:"PublicId"`
	Token             		string    		   			`json:"Token"`
	Collaterals       		map[string]float64   		`json:"Collaterals"`
	RateChange       		map[string]float64   		`json:"RateChange"`
}

// Definition of a Collateral Withdrawal request //
type RepayFunds struct {
	PublicId         		string    		  			`json:"PublicId"`
	Token             		string    		   			`json:"Token"`
	Amount            		float64              		`json:"Amount"`
}

// Definition of a the model to check a potential liquidation //
type Liquidator struct {		
	PublicId          		string    		  		    `json:"PublicId"`
	Token             		string    		   			`json:"Token"`
	RateChange        		map[string]float64   		`json:"RateChange"`
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

