package routes

import (
	"bank-management-system/internal/handlers"
	"bank-management-system/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Initialize services
	accountService := services.NewAccountService()
	loanService := services.NewLoanService()
	transactionService := services.NewTransactionService()
	jointAccountService := services.NewJointAccountService()

	// Initialize handlers
	bankHandler := handlers.NewBankHandler()
	branchHandler := handlers.NewBranchHandler()
	customerHandler := handlers.NewCustomerHandler()
	accountHandler := handlers.NewAccountHandler(accountService)
	loanHandler := handlers.NewLoanHandler(loanService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	jointAccountHandler := handlers.NewJointAccountHandler(jointAccountService)

	// checking the server is running
	router.GET("/check", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// APIs routes
	api := router.Group("/api")
	{
		// Bank routes
		banks := api.Group("/banks")
		{
			banks.POST("", bankHandler.CreateBank)
			banks.GET("", bankHandler.GetAllBanks)
			banks.GET("/:id", bankHandler.GetBank)
		}

		// Branch routes
		branches := api.Group("/branches")
		{
			branches.POST("", branchHandler.CreateBranch)
			branches.GET("", branchHandler.GetAllBranches)
			branches.GET("/:id", branchHandler.GetBranch)
			branches.GET("/bank/:bank_id", branchHandler.GetBranchesByBank)
		}

		// Customer routes
		customers := api.Group("/customers")
		{
			customers.POST("", customerHandler.CreateCustomer)
			customers.GET("", customerHandler.GetAllCustomers)
			customers.GET("/:id", customerHandler.GetCustomer)
			customers.PUT("/:id", customerHandler.UpdateCustomer)
		}

		// Account routes
		accounts := api.Group("/accounts")
		{
			accounts.POST("", accountHandler.CreateAccount)
			accounts.GET("/:id", accountHandler.GetAccount)
			accounts.GET("/customer/:customer_id", accountHandler.GetCustomerAccounts)
			accounts.POST("/:id/deposit", accountHandler.Deposit)
			accounts.POST("/:id/withdraw", accountHandler.Withdraw)
		}

		// Joint Account routes
		jointAccounts := api.Group("/joint-accounts")
		{
			jointAccounts.POST("", jointAccountHandler.CreateJointAccount)
			jointAccounts.GET("/:account_id", jointAccountHandler.GetAccountDetails)
			jointAccounts.GET("/:account_id/owners", jointAccountHandler.GetAccountOwners)
			jointAccounts.GET("/customer/:customer_id", jointAccountHandler.GetCustomerJointAccounts)
			jointAccounts.POST("/:account_id/customers", jointAccountHandler.AddCustomerToAccount)
			jointAccounts.DELETE("/:account_id/customers/:customer_id", jointAccountHandler.RemoveCustomerFromAccount)
			jointAccounts.PUT("/:account_id/customers/:customer_id/role", jointAccountHandler.UpdateCustomerRole)
		}

		// Loan routes
		loans := api.Group("/loans")
		{
			loans.POST("", loanHandler.CreateLoan)
			loans.GET("/:id", loanHandler.GetLoan)
			loans.GET("/:id/details", loanHandler.GetLoanDetails)
			loans.GET("/customer/:customer_id", loanHandler.GetCustomerLoans)
			loans.POST("/:id/repay", loanHandler.RepayLoan)

		}

		// Transaction routes
		transactions := api.Group("/transactions")
		{
			transactions.GET("/account/:id", transactionHandler.GetAccountTransactions)
			transactions.GET("/customer/:customer_id", transactionHandler.GetCustomerTransactions)
			transactions.GET("/loan/:loan_id", transactionHandler.GetLoanTransactions)
		}
	}

	return router
}
