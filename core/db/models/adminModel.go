package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Admin struct {
	ID           string   `gorm:"column:id;primaryKey"`
	Username     string   `gorm:"column:username;index"`
	Password     string   `gorm:"column:password"`
	Email        string   `gorm:"column:email"`
	Name         string   `gorm:"column:name"`
	FirstName    string   `gorm:"column:firstName"`
	LastName     string   `gorm:"column:lastName"`
	AccountOwner bool     `gorm:"column:accountOwner"`
	Locale       string   `gorm:"column:locale"`
	Permissions  []string `gorm:"column:permissions;type:text[]"`
}

// type AdminIDCounter struct{
// 	ID
// }

// type Permissions struct {
// 	Home struct {
// 		Home bool `json:"home"`
// 	} `json:"home"`
// 	Orders struct {
// 		View                   bool `json:"view"`
// 		ManageOrderInformation bool `json:"manageOrderInformation"`
// 		EditLineItems          bool `json:"editLineItems"`
// 		ApplyDiscounts         bool `json:"applyDiscounts"`
// 		SetPaymentTerms        bool `json:"setPaymentTerms"`
// 		ChargeCreditCard       bool `json:"chargeCreditCard"`
// 		ChargeVaultedCard      bool `json:"chargeVaultedCard"`
// 		MarkAsPaid             bool `json:"markAsPaid"`
// 		CapturePayments        bool `json:"capturePayments"`
// 		FullfillAndShip        bool `json:"fullfillAndShip"`
// 		BuyShippingLabels      bool `json:"buyShippingLabels"`
// 		Return                 bool `json:"return"`
// 		Refund                 bool `json:"refund"`
// 		Cancel                 bool `json:"cancel"`
// 		Export                 bool `json:"export"`
// 		Delete                 bool `json:"delete"`
// 		AbandonedCheckouts     struct {
// 			Manage bool `json:"manage"`
// 		}
// 	} `json:"orders"`
// 	DraftOrders struct {
// 		View             bool `json:"view"`
// 		CreateAndEdit    bool `json:"createAndEdit"`
// 		ApplyDiscounts   bool `json:"applyDiscounts"`
// 		SetPaymentTerms  bool `json:"setPaymentTerms"`
// 		ChargeCreditCard bool `json:"chargeCreditCard"`
// 		MarkAsPaid       bool `json:"markAsPaid"`
// 		Export           bool `json:"export"`
// 		Delete           bool `json:"delete"`
// 	} `json:"draftOrders"`
// 	Products struct {
// 		View          bool `json:"view"`
// 		ViewCost      bool `json:"viewCost"`
// 		CreateAndEdit struct {
// 			CreateAndEdit bool `json:"createAndEdit"`
// 			EditCost      bool `json:"editCost"`
// 			EditPrice     bool `json:"editPrice"`
// 		} `json:"createAndEdit"`
// 		Export    bool `json:"export"`
// 		Delete    bool `json:"delete"`
// 		Inventory struct {
// 			Manage bool `json:"manage"`
// 		} `json:"inventory"`
// 	} `json:"products"`
// 	GiftCards struct {
// 		ViewCreateAndDelete bool `json:"viewCreateAndDelete"`
// 	} `json:"giftCards"`
// 	Content struct {
// 		MetaobjectDefinitions struct {
// 			View          bool `json:"view"`
// 			CreateAndEdit bool `json:"createAndEdit"`
// 			Delete        bool `json:"delete"`
// 		} `json:"metaobjectDefinitions"`
// 		Entries struct {
// 			View          bool `json:"view"`
// 			CreateAndEdit bool `json:"createAndEdit"`
// 			Delete        bool `json:"delete"`
// 		} `json:"entries"`
// 	} `json:"content"`
// 	Customers struct {
// 		ErasePersonalData bool `json:"erasePersonalData"`
// 		RequestData       bool `json:"requestData"`
// 		Export            bool `json:"export"`
// 		Merge             bool `json:"merge"`
// 	} `json:"customers"`
// 	Analytics struct {
// 		Reports    bool `json:"reports"`
// 		Dashboards bool `json:"dashboards"`
// 	} `json:"analytics"`
// 	Marketing struct {
// 		ViewCreateAndDeleteCampaigns bool `json:"viewCreateAndDeleteCampaigns"`
// 	} `json:"marketing"`
// 	Discounts struct {
// 		ViewCreateAndDelete bool `json:"viewCreateAndDelete"`
// 	} `json:"discounts"`
// 	OnlineStore struct {
// 		BlogPostAndPages bool `json:"blogPostAndPages"`
// 		Navigation       bool `json:"navigation"`
// 	} `json:"onlineStore"`
// 	Staff struct {
// 		EditAddAndRemovePermissions bool `json:"editAddAndRemovePermissions"`
// 		ExternalLoginServices       bool `json:"externalLoginServices"`
// 		RevokeAccessToken           bool `json:"revokeAccessToken"`
// 		Collaborators               struct {
// 			ManageRequests bool `json:"manageRequests"`
// 		} `json:"collaborators"`
// 	} `json:"staff"`
// 	StoreSettings struct {
// 		ManageSettings     bool `json:"manageSettings"`
// 		TaxesAndDuties     bool `json:"taxesAndDuties"`
// 		Locations          bool `json:"locations"`
// 		Domains            bool `json:"domains"`
// 		ViewCustomerEvents bool `json:"viewCustomerEvents"`
// 		StorePolicies      bool `json:"storePolicies"`
// 	} `json:"storeSettings"`
// }

// var Admins []Admin

// var AdminSecrets = make(map[string]string)
var AdminIDCounter int

type AdminSecrets struct {
	Username string `gorm:"column:username"`
	Secret   string `gorm:"column:secret"`
}

type AdminModel struct {
	db *gorm.DB
}

func NewAdminModel(db *gorm.DB) *AdminModel {
	return &AdminModel{db: db}
}

func (ad *AdminModel) Save(admin *Admin) error {
	return ad.db.Create(admin).Error
}
func (ad *AdminModel) SaveAdminSecret(adminSecret *AdminSecrets) error {
	return ad.db.Create(adminSecret).Error
}
func AdminExistByEmail(db *gorm.DB, email string) (bool, error) {
	var count int64
	if err := db.Model(&Admin{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
func AdminExistByUsername(db *gorm.DB, username string) (bool, error) {
	var count int64
	if err := db.Model(&Admin{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
func CheckAdminTableIsEmpty(db *gorm.DB) (bool, error) {
	var count int64
	if err := db.Model(&Admin{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}
func GetAdminHashedPasswordByUsername(db *gorm.DB, username string) (string, error) {
	var admin Admin
	if err := db.Model(&Admin{}).Where("username = ?", username).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("admin not found")
		}
		return "", fmt.Errorf("error fetching user: %v", err)
	}

	return admin.Password, nil
}
func GetSecretKeyByUsername(db *gorm.DB, username string) (string, error) {
	var admin AdminSecrets
	if err := db.Model(&AdminSecrets{}).Where("username = ?", username).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("admin not found")
		}
		return "", fmt.Errorf("error fetching user: %v", err)
	}

	return admin.Secret, nil
}
