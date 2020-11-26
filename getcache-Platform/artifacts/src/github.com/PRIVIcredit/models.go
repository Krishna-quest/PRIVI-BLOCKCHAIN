///////////////////////////////////////////////////////////////
// File containing the model structs for the Cache Coin Token,
// and blockchain inputs and outputs 
///////////////////////////////////////////////////////////////                           

package main


/*---------------------------------------------------------------------------
SMART CONTRACT MODELS
-----------------------------------------------------------------------------*/  

// Definition of the model of the smart contract //
type PRIVIcreditSmartContract struct {}

// Definition the output for the smart contract //
type Output struct {
	UpdateWallets           map[string]MultiWallet      `json:"UpdateWallets"`
	UpdateLoans             map[string]PRIVIloan        `json:"UpdateLoans"`
}

// Define model parameter of risk //
type RiskParameters struct {
	Interest_min		    float64                      `json:"Interest_min"`
	Interest_max		    float64                      `json:"Interest_max"`
	P_incentive_min		    float64                      `json:"P_incentive_min"`
	P_incentive_max		    float64                      `json:"P_incentive_max"`
	P_premium_min		    float64                      `json:"P_premium_min"`
	P_premium_max		    float64                      `json:"P_premium_max"`
}

// Define instance of a Lender model on Blockchain //
type Lender struct {
	LenderId 				string 						 `json:"LenderId"`
	Amount                  float64         			 `json:"Amount"` 
	JoiningDay              int64           			 `json:"JoiningDay"` 
} 

// Define instance of a Borrower model on Blockchain //
type Borrower struct {
	BorrowerId 				string 			 			 `json:"LenderId"`
	Amount                  float64          			 `json:"Amount"` 
	JoiningDay              int64            			 `json:"JoiningDay"`
	MissingPayments         int64           			 `json:"MissingPayments"`  
	TotalPayments           int64            			 `json:"TotalPayments"`      
	Debt                    float64         			 `json:"Debt"`  
	TrustScore              float64         			 `json:"TrustScore"`
	EndorsementScore        float64          			 `json:"EndorsementScore"` 
	Collaterals 			map[string]float64 		     `json:"Collaterals"`  
} 

// Define instance of a Premium model on Blockchain //
type Premium struct {
	ProviderId 				string 			 			 `json:"ProviderId"`
	PremiumId               string                       `json:"PremiumId"`
	Risk_Pct                float64         			 `json:"Risk_Pct"`
	Premium_Amount          float64         			 `json:"Premium_Amount"`
} 

// Define model of the state of the loan //
type PRIVIstate struct {
	Funds 			        float64						 `json:"Funds"`
	Loaned                  float64                      `json:"Loaned"`
	LenderNum               int64                        `json:"LenderNum"`
	BorrowerNum             int64                        `json:"BorrowerNum"`
	ProviderNum             int64                        `json:"ProviderNum"`
	Status 					string 						 `json:"State"`
	Loan_Day 				int64   					 `json:"Loan_Day"`
	Total_Premium 			float64						 `json:"Total_Premium"`
	Total_Coverage 			float64 					 `json:"Total_Coverage"`
	PremiumList             map[string]Premium       	 `json:"PremiumList"`
	Lenders				    map[string]Lender 			 `json:"Lenders"`
	Borrowers 				map[string]Borrower  		 `json:"Borrowers"`
	Collaterals 			map[string]float64 		     `json:"Collaterals"`
}

// Define our struct to store the conditions of the loan in the Blockchain //
type PRIVIloan struct {
	LoanId 				 	string 			 			 `json:"LoanId"`
	Creator                 string                       `json:"Creator"`
	Token                   string           			 `json:"Token"`
	MaxFunds 			    float64						 `json:"MaxFunds"`
	Interest 				float64 					 `json:"Interest"`
	Payments                int64                        `json:"Payments"`
	Duration 				int64 						 `json:"Duration"`
	P_incentive 			float64 		 			 `json:"P_incentive"`
	P_premium 				float64 		 			 `json:"P_premium"`
	TrustScore              float64                      `json:"TrustScore"`
	EndorsementScore        float64                      `json:"EndorsementScore"`
	Date 					string						 `json:"Date"`
	State                   PRIVIstate                   `json:"State"`	
}

// Define instance of a Deposit Funds model on Blockchain //
type WithdrawFunds struct {
	UserId 					string 			 			 `json:"UserId"`
	LoanId                  string         			     `json:"LoanId"`
	Amount             	    float64         			 `json:"Amount"`
} 

// Define instance of a Deposit Funds model on Blockchain //
type DepositFunds struct {
	LenderId 				string 			 			 `json:"LenderId"`
	LoanId                  string         			     `json:"LoanId"`
	Amount             	    float64         			 `json:"Amount"`
} 

// Define instance of a Borrow Funds model on Blockchain //
type BorrowFunds struct {
	BorrowerId 			    string 			 			 `json:"BorrowerId"`
	LoanId                  string         			     `json:"LoanId"`
	Amount             	    float64         			 `json:"Amount"`
	Collaterals 			map[string]float64 		     `json:"Collaterals"`
} 

// Define instance of a Assume Risk model on Blockchain //
type AssumeRisk struct {
	ProviderId 			    string 			 			 `json:"ProviderId"`
	PremiumId 			    string 			 			 `json:"PremiumId"`
	LoanId                  string         			     `json:"LoanId"`
	Risk_Pct                float64         			 `json:"Risk_Pct"`
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
