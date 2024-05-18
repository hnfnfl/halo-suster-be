package dto

type Role string

const (
	IT    Role = "it"
	Nurse Role = "nurse"
)

type Sort string

const (
	ASC  Sort = "ASC"
	DESC Sort = "DESC"
)

type (
	ReqParamUserGet struct {
		UserID    string `json:"userId"`
		Limit     int    `json:"limit"`
		Offset    int    `json:"offset"`
		Name      string `json:"name"`
		NIP       string `json:"nip"`
		Role      Role   `json:"role"`
		CreatedAt Sort   `json:"createdAt"`
	}
	ResUserGet struct {
		UserID    string `json:"userId"`
		NIP       int    `json:"nip"`
		Name      string `json:"name"`
		CreatedAt string `json:"createdAt"`
	}
)
