package database

type Admin struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	AccountOwner bool   `json:"accountOwner"`
	Locale       string `json:"locale"`
	Permissions  struct {
		Home struct {
			Home bool `json:"home"`
		} `json:"home"`
		Orders struct {
			View                   bool `json:"view"`
			ManageOrderInformation bool `json:"manageOrderInformation"`
			EditLineItems          bool `json:"editLineItems"`
			ApplyDiscounts         bool `json:"applyDiscounts"`
			SetPaymentTerms        bool `json:"setPaymentTerms"`
			ChargeCreditCard       bool `json:"chargeCreditCard"`
			ChargeVaultedCard      bool `json:"chargeVaultedCard"`
			MarkAsPaid             bool `json:"markAsPaid"`
			CapturePayments        bool `json:"capturePayments"`
			FullfillAndShip        bool `json:"fullfillAndShip"`
			BuyShippingLabels      bool `json:"buyShippingLabels"`
			Return                 bool `json:"return"`
			Refund                 bool `json:"refund"`
			Cancel                 bool `json:"cancel"`
			Export                 bool `json:"export"`
			Delete                 bool `json:"delete"`
			AbandonedCheckouts     struct {
				Manage bool `json:"manage"`
			}
		} `json:"orders"`
		DraftOrders struct {
			View             bool `json:"view"`
			CreateAndEdit    bool `json:"createAndEdit"`
			ApplyDiscounts   bool `json:"applyDiscounts"`
			SetPaymentTerms  bool `json:"setPaymentTerms"`
			ChargeCreditCard bool `json:"chargeCreditCard"`
			MarkAsPaid       bool `json:"markAsPaid"`
			Export           bool `json:"export"`
			Delete           bool `json:"delete"`
		} `json:"draftOrders"`
		Products struct {
			View          bool `json:"view"`
			ViewCost      bool `json:"viewCost"`
			CreateAndEdit struct {
				CreateAndEdit bool `json:"createAndEdit"`
				EditCost      bool `json:"editCost"`
				EditPrice     bool `json:"editPrice"`
			} `json:"createAndEdit"`
			Export    bool `json:"export"`
			Delete    bool `json:"delete"`
			Inventory struct {
				Manage bool `json:"manage"`
			} `json:"inventory"`
		} `json:"products"`
		GiftCards struct {
			ViewCreateAndDelete bool `json:"viewCreateAndDelete"`
		} `json:"giftCards"`
		Content struct {
			MetaobjectDefinitions struct {
				View          bool `json:"view"`
				CreateAndEdit bool `json:"createAndEdit"`
				Delete        bool `json:"delete"`
			} `json:"metaobjectDefinitions"`
			Entries struct {
				View          bool `json:"view"`
				CreateAndEdit bool `json:"createAndEdit"`
				Delete        bool `json:"delete"`
			} `json:"entries"`
		} `json:"content"`
		Customers struct {
			ErasePersonalData bool `json:"erasePersonalData"`
			RequestData       bool `json:"requestData"`
			Export            bool `json:"export"`
			Merge             bool `json:"merge"`
		} `json:"customers"`
		Analytics struct {
			Reports    bool `json:"reports"`
			Dashboards bool `json:"dashboards"`
		} `json:"analytics"`
		Marketing struct {
			ViewCreateAndDeleteCampaigns bool `json:"viewCreateAndDeleteCampaigns"`
		} `json:"marketing"`
		Discounts struct {
			ViewCreateAndDelete bool `json:"viewCreateAndDelete"`
		} `json:"discounts"`
		OnlineStore struct {
			BlogPostAndPages bool `json:"blogPostAndPages"`
			Navigation       bool `json:"navigation"`
		} `json:"onlineStore"`
		Staff struct {
			EditAddAndRemovePermissions bool `json:"editAddAndRemovePermissions"`
			ExternalLoginServices       bool `json:"externalLoginServices"`
			RevokeAccessToken           bool `json:"revokeAccessToken"`
			Collaborators               struct {
				ManageRequests bool `json:"manageRequests"`
			} `json:"collaborators"`
		} `json:"staff"`
		StoreSettings struct {
			ManageSettings     bool `json:"manageSettings"`
			TaxesAndDuties     bool `json:"taxesAndDuties"`
			Locations          bool `json:"locations"`
			Domains            bool `json:"domains"`
			ViewCustomerEvents bool `json:"viewCustomerEvents"`
			StorePolicies      bool `json:"storePolicies"`
		} `json:"storeSettings"`
	} `json:"permissions"`
}

var Admins []Admin
var AdminSecrets = make(map[string]string)
var AdminIDCounter int
