///////////////////////////////////////////////////////////////
// File containing the model structs for the Cache Coin Token,
// and blockchain inputs and outputs 
///////////////////////////////////////////////////////////////                           

package main

/*---------------------------------------------------------------------------
SMART CONTRACT MODELS FOR LIQUIDITY POOLS
-----------------------------------------------------------------------------*/  


// Definition of the model of the smart contract //
type DataProtocolSmartContract struct {}



// Definition the output for the smart contract //
type Output struct {
	ID                	    string    					`json:"ID"`
	DID                	    string    					`json:"DID"`
	UpdateWallets           map[string]MultiWallet      `json:"UpdateWallets"`
	UpdateUsers             map[string]Actor            `json:"UpdateUsers"`
}

// Definition of an actor in the Cache Ecosystem //
type Actor struct {
	PublicId        string          `json:"PublicId"`
	Role            string          `json:"Role"`
	Privacy         map[string]bool `json:"Privacy"`
	TargetTIDs      map[string]bool `json:"TargetTIDs"`       
}

// Definition of the encryption object with DIDs //
type Encryption struct {
	PublicId        string `json:"PublicId"`
	DID             string `json:"DID"`
}

// Definition of the privacy modifier object //
type PrivacyModifier struct {
	PublicId        string `json:"PublicId"`
	BusinessId      string `json:"BusinessId"`
	Enabled         bool   `json:"Enabled"`
}

// Definition of the insigth discovery object //
type InsightDiscovery struct {
	Business_Id     string   `json:"Business_Id"`
	DID_list        []string `json:"DID_list"`
	ID_list         []string `json:"ID_list"`
}

// Definition of the insigth purchase object //
type InsightPurchase struct {
	DID_list        		[]string 				`json:"DID_list"`
	Price           		float64  				`json:"Price"`
	Business_Id     		string   				`json:"Business_Id"`
	InsightDistribution     map[string]float64 		`json:"InsightDistribution"`
}

// Definition of the insigth target object //
type InsightTarget struct {
	TID_list        		[]string 				`json:"TID_list"`
	Business_Id     		string   				`json:"Business_Id"`
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
