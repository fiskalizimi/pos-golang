package fiskalizimi

import (
	"fiskalizimi/proto"
	"time"
)

func GetCitizenCoupon() *proto.CitizenCoupon {
	return &proto.CitizenCoupon{
		BusinessId: 1,
		PosId:      1,
		CouponId:   1234,
		Type:       proto.CouponType_Sale,
		Time:       time.Date(2024, time.September, 24, 6, 11, 29, 0, time.Local).Unix(),
		Total:      1820,
		TaxGroups: []*proto.TaxGroup{
			{TaxRate: "C", TotalForTax: 450, TotalTax: 0},
			{TaxRate: "D", TotalForTax: 320, TotalTax: 26},
			{TaxRate: "E", TotalForTax: 1050, TotalTax: 189},
		},
		TotalTax: 215,
	}
}

func GetPosCoupon() *proto.PosCoupon {
	return &proto.PosCoupon{
		BusinessId:     1,
		CouponId:       1234,
		BranchId:       3,
		Location:       "Prishtine",
		OperatorId:     "Kushtrimi",
		PosId:          1,
		ApplicationId:  1,
		VerificationNo: "1234567890123456",
		Type:           proto.CouponType_Sale,
		Time:           time.Date(2024, time.September, 24, 6, 11, 29, 0, time.Local).Unix(),
		Items: []*proto.CouponItem{
			{Name: "uje rugove", Price: 150, Unit: "cope", Quantity: 3, Total: 450, TaxRate: "C", Type: "TT"},
			{Name: "sendviq", Price: 300, Unit: "cope", Quantity: 2, Total: 600, TaxRate: "E", Type: "TT"},
			{Name: "buke", Price: 80, Unit: "cope", Quantity: 4, Total: 320, TaxRate: "D", Type: "TT"},
			{Name: "machiato e madhe", Unit: "cope", Price: 150, Quantity: 3, Total: 450, TaxRate: "E", Type: "TT"},
		},
		Payments: []*proto.Payment{
			{Type: proto.PaymentType_Cash, Amount: 500},
			{Type: proto.PaymentType_CreditCard, Amount: 1000},
			{Type: proto.PaymentType_Voucher, Amount: 320},
		},
		Total: 1820,
		TaxGroups: []*proto.TaxGroup{
			{TaxRate: "C", TotalForTax: 450, TotalTax: 0},
			{TaxRate: "D", TotalForTax: 320, TotalTax: 26},
			{TaxRate: "E", TotalForTax: 1050, TotalTax: 189},
		},
		TotalTax:   215,
		TotalNoTax: 1605,
	}
}
