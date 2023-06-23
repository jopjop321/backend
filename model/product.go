package product

type Product struct {
	ID           int    `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Image        string `json:"image"`
	Price_Cost   int    `json:"cost_price"`
	Price_Member int    `json:"member_price"`
	Price_Normal int    `json:"normal_price"`
	Sell         int    `json:"sell"`
	Amount       int    `json:"amount"`
}
