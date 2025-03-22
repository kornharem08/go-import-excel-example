package models

type PurchaseOrder struct {
	JobIDNo            *string `json:"job_id_no"`
	Type               *string `json:"type"`
	SalesTeam          *string `json:"sales_team"`
	ProjectManager     *string `json:"project_manager"`
	Customer           *string `json:"customer"`
	ProductCode        *string `json:"product_code"`
	ProductDescription *string `json:"product_description"`
	Ordered            *int    `json:"ordered"`
	Received           *int    `json:"received"`
	Remain             *int    `json:"remain"`
	PR                 *string `json:"pr"`
	PRDate             *string `json:"pr_date"`
	PO                 *string `json:"po"`
	PODate             *string `json:"po_date"`
	Distribution       *string `json:"distribution"`
	PaymentTerm        *string `json:"payment_term"`
	RequestDate        *string `json:"request_date"`
	DeliveryDate       *string `json:"delivery_date"`
	Status             *string `json:"status"`
}
