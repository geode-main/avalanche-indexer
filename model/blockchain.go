package model

type Blockchain struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SubnetID string `json:"subnet_id"`
	VMID     string `json:"vm_id"`
}
