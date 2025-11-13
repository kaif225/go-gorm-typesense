package models

type Images struct {
	ID              int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Category        string `gorm:"column:category;size=100;not null" json:"category,omitempty"`
	FileName        string `gorm:"column:file_name;size=100;not null"`
	S3URL           string `gorm:"column:s3_url;size=100;not null" json:"s3_url,omitempty"`
	S3Key           string `gorm:"column:s3_key;size=100;not null" json:"s3_key,omitempty"`
	UploadedAt      string `gorm:"column:uploaded_at;size=100;not null" json:"uploaded_at,omitempty"`
	TypesenseSynced bool   `gorm:"column:typesense_synced;default:false" json:"typesense_synced,omitempty"`
}

type TypesenseImage struct {
	ID         string `json:"id"`
	Category   string `json:"category"`
	FileName   string `json:"file_name"`
	S3Key      string `json:"s3_key"`
	S3URL      string `json:"s3_url"`
	UploadedAt string `json:"uploaded_at"`
}

type ImageResponse struct {
	Category   string `json:"category"`
	FileName   string `json:"file_name"`
	S3Key      string `json:"s3_key"`
	S3URL      string `json:"s3_url"`
	SignedURL  string `json:"signed_url"`
	UploadedAt string `json:"uploaded_at"`
}
