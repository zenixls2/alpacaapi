package alpacaapi

type AccountStatus string

const (
	AccountStatusOnBoarding       AccountStatus = "ONBOARDING"
	AccountStatusSubmissionFailed AccountStatus = "SUBMISSION_FAILED"
	AccountStatusSubmitted        AccountStatus = "SUBMITTED"
	AccountStatusAccountUpdated   AccountStatus = "ACCOUNT_UPDATED"
	AccountStatusApprovalPending  AccountStatus = "APPROVAL_PENDING"
	AccountStatusActive           AccountStatus = "ACTIVE"
	AccountStatusRejected         AccountStatus = "REJECTED"
)

type OptionsLevel int

const (
	OptionsDisabled       OptionsLevel = 0
	OptionsSecurePutCall  OptionsLevel = 1
	OptionsLongPutCall    OptionsLevel = 2
	OptionsSpreadStraddle OptionsLevel = 3
)
