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

	// Initialize handlers
	bankHandler := handlers.NewBankHandler()
	branchHandler := handlers.NewBranchHandler()
	customerHandler := handlers.NewCustomerHandler()
	accountHandler := handlers.NewAccountHandler(accountService)
	loanHandler := handlers.NewLoanHandler(loanService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

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
			banks.GET("/:id", bankHandler.GetBank)
		}

		// Branch routes
		branches := api.Group("/branches")
		{
			branches.POST("", branchHandler.CreateBranch)
			branches.GET("/:id", branchHandler.GetBranch)
			branches.GET("/bank/:bank_id", branchHandler.GetBranchesByBank)
		}

		// Customer routes
		customers := api.Group("/customers")
		{
			customers.POST("", customerHandler.CreateCustomer)
			customers.GET("/:id", customerHandler.GetCustomer)
			customers.DELETE("/:id", customerHandler.DeleteCustomer)
		}

		// Account routes
		accounts := api.Group("/accounts")
		{
			accounts.POST("", accountHandler.CreateAccount)
			accounts.GET("/:id", accountHandler.GetAccount)
			accounts.POST("/:id/customers", accountHandler.AddCustomerToAccount)
			accounts.DELETE("/:id/customers/:customer_id", accountHandler.RemoveCustomerFromAccount)
		}

		// Loan routes
		loans := api.Group("/loans")
		{
			loans.POST("", loanHandler.CreateLoan)
			loans.GET("/:id/details", loanHandler.GetLoanDetails)
			loans.POST("/:id/repay", loanHandler.RepayLoan)
		}

		// Transaction routes
		transactions := api.Group("/transactions")
		{
			transactions.POST("/:id", transactionHandler.CreateTransaction)
			transactions.GET("/account/:id", transactionHandler.GetAccountTransactions)
			transactions.GET("/loan/:loan_id", transactionHandler.GetLoanTransactions)
		}
	}

	return router
}
