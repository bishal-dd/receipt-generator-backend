// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreateBulkService struct {
	Description string  `json:"description"`
	Rate        float64 `json:"rate"`
	Quantity    int     `json:"quantity"`
	Amount      float64 `json:"amount"`
	ID          *string `json:"id,omitempty"`
}

type CreateProduct struct {
	Name      string  `json:"name"`
	UnitPrice float64 `json:"unit_price"`
	Quantity  *int    `json:"quantity,omitempty"`
	UserID    string  `json:"user_id"`
}

type CreateProfile struct {
	CompanyName            *string  `json:"company_name,omitempty"`
	LogoImage              *string  `json:"logo_image,omitempty"`
	PhoneNo                *string  `json:"phone_no,omitempty"`
	Email                  *string  `json:"email,omitempty"`
	Address                *string  `json:"address,omitempty"`
	Currency               *string  `json:"currency,omitempty"`
	Tax                    *float64 `json:"tax,omitempty"`
	PhoneNumberCountryCode string   `json:"phone_number_country_code"`
	City                   *string  `json:"city,omitempty"`
	Title                  *string  `json:"title,omitempty"`
	SignatureImage         *string  `json:"signature_image,omitempty"`
	UserID                 string   `json:"user_id"`
}

type CreateReceipt struct {
	ReceiptName      *string  `json:"receipt_name,omitempty"`
	RecipientName    *string  `json:"recipient_name,omitempty"`
	RecipientPhone   *string  `json:"recipient_phone,omitempty"`
	RecipientEmail   *string  `json:"recipient_email,omitempty"`
	RecipientAddress *string  `json:"recipient_address,omitempty"`
	IsReceiptSend    bool     `json:"is_receipt_send"`
	ReceiptNo        string   `json:"receipt_no"`
	PaymentMethod    string   `json:"payment_method"`
	PaymentNote      *string  `json:"payment_note,omitempty"`
	UserID           string   `json:"user_id"`
	Date             string   `json:"date"`
	TotalAmount      *float64 `json:"total_amount,omitempty"`
}

type CreateService struct {
	Description string  `json:"description"`
	Rate        float64 `json:"rate"`
	Quantity    int     `json:"quantity"`
	Amount      float64 `json:"amount"`
	ReceiptID   string  `json:"receipt_id"`
}

type CreateUser struct {
	ID string `json:"id"`
}

type DownloadPDF struct {
	ReceiptName      *string              `json:"receipt_name,omitempty"`
	RecipientName    *string              `json:"recipient_name,omitempty"`
	RecipientPhone   *string              `json:"recipient_phone,omitempty"`
	RecipientEmail   *string              `json:"recipient_email,omitempty"`
	RecipientAddress *string              `json:"recipient_address,omitempty"`
	ReceiptNo        string               `json:"receipt_no"`
	PaymentMethod    string               `json:"payment_method"`
	PaymentNote      *string              `json:"payment_note,omitempty"`
	IsReceiptSend    bool                 `json:"is_receipt_send"`
	UserID           string               `json:"user_id"`
	OrginazationID   string               `json:"orginazation_id"`
	Date             string               `json:"date"`
	Services         []*CreateBulkService `json:"Services,omitempty"`
}

type Mutation struct {
}

type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor,omitempty"`
	EndCursor       *string `json:"endCursor,omitempty"`
}

type Product struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	UnitPrice float64 `json:"unit_price"`
	Quantity  *int    `json:"quantity,omitempty"`
	UserID    string  `json:"user_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

type Profile struct {
	ID                     string  `json:"id"`
	CompanyName            *string `json:"company_name,omitempty"`
	LogoImage              *string `json:"logo_image,omitempty"`
	PhoneNo                *string `json:"phone_no,omitempty"`
	Email                  *string `json:"email,omitempty"`
	Address                *string `json:"address,omitempty"`
	City                   *string `json:"city,omitempty"`
	Title                  *string `json:"title,omitempty"`
	SignatureImage         *string `json:"signature_image,omitempty"`
	Currency               string  `json:"currency"`
	Tax                    float64 `json:"tax"`
	PhoneNumberCountryCode string  `json:"phone_number_country_code"`
	UserID                 string  `json:"user_id"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              *string `json:"updated_at,omitempty"`
	DeletedAt              *string `json:"deleted_at,omitempty"`
}

type Query struct {
}

type Receipt struct {
	ID               string     `json:"id"`
	ReceiptName      *string    `json:"receipt_name,omitempty"`
	RecipientName    *string    `json:"recipient_name,omitempty"`
	RecipientPhone   *string    `json:"recipient_phone,omitempty"`
	RecipientEmail   *string    `json:"recipient_email,omitempty"`
	RecipientAddress *string    `json:"recipient_address,omitempty"`
	ReceiptNo        string     `json:"receipt_no"`
	UserID           string     `json:"user_id"`
	Date             string     `json:"date"`
	TotalAmount      *float64   `json:"total_amount,omitempty"`
	SubTotalAmount   *float64   `json:"sub_total_amount,omitempty"`
	TaxAmount        *float64   `json:"tax_amount,omitempty"`
	PaymentMethod    string     `json:"payment_method"`
	PaymentNote      *string    `json:"payment_note,omitempty"`
	IsReceiptSend    bool       `json:"is_receipt_send"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        *string    `json:"updated_at,omitempty"`
	DeletedAt        *string    `json:"deleted_at,omitempty"`
	Services         []*Service `json:"Services,omitempty"`
}

type ReceiptConnection struct {
	Edges      []*ReceiptEdge `json:"edges"`
	PageInfo   *PageInfo      `json:"pageInfo"`
	TotalCount int            `json:"totalCount"`
}

type ReceiptEdge struct {
	Cursor string   `json:"cursor"`
	Node   *Receipt `json:"node"`
}

type SearchReceipt struct {
	Receipts   []*Receipt `json:"receipts"`
	TotalCount int        `json:"totalCount"`
	FoundCount int        `json:"foundCount"`
}

type SendReceiptPDFToEmail struct {
	ReceiptName      *string              `json:"receipt_name,omitempty"`
	RecipientName    *string              `json:"recipient_name,omitempty"`
	RecipientPhone   *string              `json:"recipient_phone,omitempty"`
	RecipientEmail   string               `json:"recipient_email"`
	RecipientAddress *string              `json:"recipient_address,omitempty"`
	ReceiptNo        string               `json:"receipt_no"`
	PaymentMethod    string               `json:"payment_method"`
	PaymentNote      *string              `json:"payment_note,omitempty"`
	IsReceiptSend    bool                 `json:"is_receipt_send"`
	UserID           string               `json:"user_id"`
	OrginazationID   string               `json:"orginazation_id"`
	Date             string               `json:"date"`
	Services         []*CreateBulkService `json:"Services,omitempty"`
}

type SendReceiptPDFToWhatsApp struct {
	ReceiptName      *string              `json:"receipt_name,omitempty"`
	RecipientName    *string              `json:"recipient_name,omitempty"`
	RecipientPhone   string               `json:"recipient_phone"`
	RecipientEmail   *string              `json:"recipient_email,omitempty"`
	RecipientAddress *string              `json:"recipient_address,omitempty"`
	ReceiptNo        string               `json:"receipt_no"`
	PaymentMethod    string               `json:"payment_method"`
	PaymentNote      *string              `json:"payment_note,omitempty"`
	UserID           string               `json:"user_id"`
	IsReceiptSend    bool                 `json:"is_receipt_send"`
	OrginazationID   string               `json:"orginazation_id"`
	Date             string               `json:"date"`
	Services         []*CreateBulkService `json:"Services,omitempty"`
}

type Service struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Rate        float64 `json:"rate"`
	Quantity    int     `json:"quantity"`
	Amount      float64 `json:"amount"`
	ReceiptID   string  `json:"receipt_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
}

type UpdateProduct struct {
	ID        string   `json:"id"`
	Name      *string  `json:"name,omitempty"`
	UnitPrice *float64 `json:"unit_price,omitempty"`
	Quantity  *int     `json:"quantity,omitempty"`
}

type UpdateProfile struct {
	ID                     string   `json:"id"`
	CompanyName            *string  `json:"company_name,omitempty"`
	LogoImage              *string  `json:"logo_image,omitempty"`
	PhoneNo                *string  `json:"phone_no,omitempty"`
	Email                  *string  `json:"email,omitempty"`
	Address                *string  `json:"address,omitempty"`
	Currency               *string  `json:"currency,omitempty"`
	PhoneNumberCountryCode *string  `json:"phone_number_country_code,omitempty"`
	Tax                    *float64 `json:"tax,omitempty"`
	City                   *string  `json:"city,omitempty"`
	Title                  *string  `json:"title,omitempty"`
	SignatureImage         *string  `json:"signature_image,omitempty"`
}

type UpdateReceipt struct {
	ID               string   `json:"id"`
	ReceiptName      *string  `json:"receipt_name,omitempty"`
	RecipientName    *string  `json:"recipient_name,omitempty"`
	RecipientPhone   *string  `json:"recipient_phone,omitempty"`
	RecipientEmail   *string  `json:"recipient_email,omitempty"`
	RecipientAddress *string  `json:"recipient_address,omitempty"`
	ReceiptNo        *string  `json:"receipt_no,omitempty"`
	PaymentMethod    *string  `json:"payment_method,omitempty"`
	PaymentNote      *string  `json:"payment_note,omitempty"`
	IsReceiptSend    *bool    `json:"is_receipt_send,omitempty"`
	UserID           *string  `json:"user_id,omitempty"`
	Date             *string  `json:"date,omitempty"`
	TotalAmount      *float64 `json:"total_amount,omitempty"`
}

type UpdateService struct {
	ID          string   `json:"id"`
	Description *string  `json:"description,omitempty"`
	Rate        *float64 `json:"rate,omitempty"`
	Quantity    *int     `json:"quantity,omitempty"`
	Amount      *float64 `json:"amount,omitempty"`
}

type User struct {
	ID        string   `json:"id"`
	Mode      string   `json:"mode"`
	UseCount  int      `json:"use_count"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt *string  `json:"updated_at,omitempty"`
	DeletedAt *string  `json:"deleted_at,omitempty"`
	Profile   *Profile `json:"Profile,omitempty"`
}

type UserConnection struct {
	Edges      []*UserEdge `json:"edges"`
	PageInfo   *PageInfo   `json:"pageInfo"`
	TotalCount int         `json:"totalCount"`
}

type UserEdge struct {
	Cursor string `json:"cursor"`
	Node   *User  `json:"node"`
}
