package models

type PurchaseOrder struct {
	JobIDNo            *string `json:"job_id_no"`
	SalesTeam          *string `json:"sales_team"`
	ProjectManager     *string `json:"project_manager"`
	Purchasing         *string `json:"purchasing"`
	CustomerPO         *string `json:"customer_po"`
	JobAmount          *int    `json:"job_amount"`
	PeriodStart        *string `json:"period_start"`
	PeriodEnd          *string `json:"period_end"`
	Customer           *string `json:"customer"`
	ProductCode        *string `json:"product_code"`
	ProductDescription *string `json:"product_description"`
	Ordered            *int    `json:"ordered"`
	Received           *int    `json:"received"`
	Remain             *int    `json:"remain"`
	Currency           *string `json:"currency"`
	UnitListPrice      *int    `json:"unit_list_price"`
	ExtendListPrice    *int    `json:"extend_list_price"`
	DiscountPercent    *int    `json:"discount_percent"`
	DiscountAmount     *int    `json:"discount_amount"`
	ExtendUnitNetPrice *int    `json:"extend_unit_net_price"`
	ExtendNetPrice     *int    `json:"extend_net_price"`
	DeliveryDate       *string `json:"delivery_date"`
	Status             *string `json:"status"`
}
