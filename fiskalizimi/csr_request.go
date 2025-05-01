package main

type CsrRequest struct {
	Country      string `json:"country"`
	BusinessName string `json:"businessName"`
	Nui          uint64 `json:"nui"`
	BranchId     uint64 `json:"branchId"`
	PosId        uint64 `json:"posId"`
}
