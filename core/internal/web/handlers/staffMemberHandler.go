package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/jyotirmoydotdev/openfy/database"
)

func ValidatePermissions(permissions []string) bool {
	var validPermissions = map[string]bool{
		"APPLICATIONS":                             true,
		"BEACONS":                                  true,
		"BILLING_APPLICATION_CHARGES":              true,
		"CHANNELS":                                 true,
		"CONTENT":                                  true,
		"CONTENT_ENTRIES_DELETE":                   true,
		"CONTENT_ENTRIES_EDIT":                     true,
		"CONTENT_ENTRIES_VIEW":                     true,
		"CONTENT_MODELS_DELETE":                    true,
		"CONTENT_MODELS_EDIT":                      true,
		"CONTENT_MODELS_VIEW":                      true,
		"CREATE_STORE_CREDIT_ACCOUNT_TRANSACTIONS": true,
		"CUSTOM_PIXELS_MANAGEMENT":                 true,
		"CUSTOM_PIXELS_VIEW":                       true,
		"CUSTOMERS":                                true,
		"DASHBOARD":                                true,
		"DELETE_PRODUCTS":                          true,
		"DOMAINS":                                  true,
		"DRAFT_ORDERS":                             true,
		"CREATE_AND_EDIT_DRAFT_ORDERS":             true,
		"APPLY_DISCOUNTS_TO_DRAFT_ORDERS":          true,
		"MARK_DRAFT_ORDERS_AS_PAID":                true,
		"SET_PAYMENT_TERMS_FOR_DRAFT_ORDERS":       true,
		"DELETE_DRAFT_ORDERS":                      true,
		"EDIT_ORDERS":                              true,
		"EDIT_PRIVATE_APPS":                        true,
		"EDIT_PRODUCT_COST":                        true,
		"EDIT_PRODUCT_PRICE":                       true,
		"EDIT_THEME_CODE":                          true,
		"GIFT_CARDS":                               true,
		"LINKS":                                    true,
		"LOCATIONS":                                true,
		"MANAGE_DELIVERY_SETTINGS":                 true,
		"MANAGE_INVENTORY":                         true,
		"MANAGE_POLICIES":                          true,
		"MANAGE_PRODUCT_TAGS":                      true,
		"MANAGE_PRODUCTS":                          true,
		"MANAGE_STORE_CREDIT_SETTINGS":             true,
		"MANAGE_TAXES_SETTINGS":                    true,
		"MARKETING":                                true,
		"MARKETING_SECTION":                        true,
		"METAOBJECTS_DELETE":                       true,
		"METAOBJECTS_EDIT":                         true,
		"METAOBJECTS_VIEW":                         true,
		"METAOBJECT_DEFINITIONS_DELETE":            true,
		"METAOBJECT_DEFINITIONS_EDIT":              true,
		"METAOBJECT_DEFINITIONS_VIEW":              true,
		"MERGE_CUSTOMERS":                          true,
		"ORDERS":                                   true,
		"OVERVIEWS":                                true,
		"PAGES":                                    true,
		"PAY_DRAFT_ORDERS_BY_CREDIT_CARD":          true,
		"PAY_ORDERS_BY_CREDIT_CARD":                true,
		"PAY_ORDERS_BY_VAULTED_CARD":               true,
		"PREFERENCES":                              true,
		"PRODUCTS":                                 true,
		"REFUND_ORDERS":                            true,
		"REPORTS":                                  true,
		"TRANSLATIONS":                             true,
		"THEMES":                                   true,
		"VIEW_ALL_SHOPIFY_CREDIT_TRANSACTIONS":     true,
		"VIEW_BALANCE_BANK_ACCOUNTS":               true,
		"VIEW_PRIVATE_APPS":                        true,
		"VIEW_PRODUCT_COSTS":                       true,
		"VIEW_STORE_CREDIT_ACCOUNT_TRANSACTIONS":   true,
		"APPLY_DISCOUNTS_TO_ORDERS":                true,
		"FULFILL_AND_SHIP_ORDERS":                  true,
		"BUY_SHIPPING_LABELS":                      true,
		"RETURN_ORDERS":                            true,
		"MANAGE_ABANDONED_CHECKOUTS":               true,
		"CANCEL_ORDERS":                            true,
		"DELETE_ORDERS":                            true,
		"MANAGE_ORDERS_INFORMATION":                true,
		"SET_PAYMENT_TERMS_FOR_ORDERS":             true,
		"MARK_ORDERS_AS_PAID":                      true,
		"CAPTURE_PAYMENTS_FOR_ORDERS":              true,
		"VIEW_COMPANIES":                           true,
		"CREATE_AND_EDIT_COMPANIES":                true,
		"DELETE_COMPANIES":                         true,
		"MANAGE_COMPANY_LOCATION_ASSIGNMENTS":      true,
		"THIRD_PARTY_MONEY_MOVEMENT":               true,
		"EXPORT_CUSTOMERS":                         true,
		"EXPORT_DRAFT_ORDERS":                      true,
		"EXPORT_ORDERS":                            true,
		"EXPORT_PRODUCTS":                          true,
		"SHOPIFY_PAYMENTS_ACCOUNTS":                true,
		"SHOPIFY_PAYMENTS_TRANSFERS":               true,
		"STAFF_AUDIT_LOG_VIEW":                     true,
		"STAFF_MANAGEMENT_UPDATE":                  true,
		"APPLICATIONS_BILLING":                     true,
		"ATTESTATION_AUTHORITY":                    true,
		"AUTHENTICATION_MANAGEMENT":                true,
		"BALANCE_BANK_ACCOUNTS_MANAGEMENT":         true,
		"BILLING_CHARGES":                          true,
		"BILLING_INVOICES_PAY":                     true,
		"BILLING_INVOICES_VIEW":                    true,
		"BILLING_PAYMENT_METHODS_MANAGE":           true,
		"BILLING_PAYMENT_METHODS_VIEW":             true,
		"BILLING_SETTINGS":                         true,
		"BILLING_SUBSCRIPTIONS":                    true,
		"CAPITAL":                                  true,
		"CUSTOMER_PRIVATE_DATA":                    true,
		"ERASE_CUSTOMER_DATA":                      true,
		"REQUEST_CUSTOMER_DATA":                    true,
		"DOMAINS_MANAGEMENT":                       true,
		"DOMAINS_TRANSFER_OUT":                     true,
		"ENABLE_PRIVATE_APPS":                      true,
		"EXPERIMENTS_MANAGEMENT":                   true,
		"GDPR_ACTIONS":                             true,
		"MANAGE_ALL_SHOPIFY_CREDIT_CARDS":          true,
		"MANAGE_TAP_TO_PAY":                        true,
		"PAYMENT_SETTINGS":                         true,
		"UPGRADE_TO_PLUS_PLAN":                     true,
		"SHOPIFY_PAYMENTS":                         true,
		"SQLITE_BULK_DATA_TRANSFER":                true,
		"STAFF_API_PERMISSION_MANAGEMENT":          true,
		"STAFF_MANAGEMENT":                         true,
		"STAFF_MANAGEMENT_ACTIVATION":              true,
		"STAFF_MANAGEMENT_CREATE":                  true,
		"STAFF_MANAGEMENT_DELETE":                  true,
		"SUPPORT_METHODS":                          true,
		"COLLABORATOR_REQUEST_MANAGEMENT":          true,
		"COLLABORATOR_REQUEST_SETTINGS":            true,
		"VIEW_PRICE_LISTS":                         true,
		"DELETE_PRICE_LISTS":                       true,
		"CREATE_AND_EDIT_PRICE_LISTS":              true,
		"VIEW_CATALOGS":                            true,
		"DELETE_CATALOGS":                          true,
		"CREATE_AND_EDIT_CATALOGS":                 true,
	}
	for _, permission := range permissions {
		if !validPermissions[permission] {
			return false
		}
	}
	return true
}

func SavePermission(ctx *gin.Context) {
	var newRequest struct {
		ID          string   `json:"id"`
		Permissions []string `json:"Permissions"`
	}
	if err := ctx.ShouldBindJSON(&newRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !ValidatePermissions(newRequest.Permissions) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid Permission detected",
		})
		return
	}
	_, err := db.GetStaffMemberDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	// TODO: Store the permission in the database
	// NOTE: If the sender is owner or STAFF_API_PERMISSION_MANAGEMENT
	//       is in the database before allow to change
}
