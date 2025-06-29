package models

import "time"

type CreateTask struct {
	ID          int    `json:"ID"`
	Title       string `json:"Title" validate:"required"`
	Description string `json:"Description"`
	Status      string `json:"Status" validate:"oneof='Pending' 'In Progress' 'Completed'"`
	DueDate     string `json:"DueDate"`
	CreatedAt   string `json:"CreatedAt"`
	UpdatedAt   string `json:"UpdatedAt"`
	APIStatus   string `json:"APIStatus"`
	APIError    string `json:"APIError"`
}

type Users struct {
	ID           int    `json:"ID"`
	UserName     string `json:"UserName"`
	PasswordHash string `json:"PasswordHash"`
	CreatedAt    string `json:"CreatedAt"`
	UpdatedAt    string `json:"UpdatedAt"`
}

type TaskQueryParams struct {
	Page          int    `validate:"gte=1" schema:"page"`
	Limit         int    `validate:"gte=1,lte=50" schema:"limit"`
	Status        string `schema:"status"`
	DueDateAfter  string `schema:"due_date_after"`
	DueDateBefore string `schema:"due_date_before"`
	SortBy        string `validate:"omitempty,oneof=created_at updated_at due_date" schema:"sort_by"`
	SortOrder     string `validate:"omitempty,oneof=asc desc" schema:"sort_order"`
}

type GetAllTaskResp struct {
	Tasks     []CreateTask `json:"tasks"`
	Page      int          `json:"page"`
	Limit     int          `json:"limit"`
	Total     int          `json:"total"`
	APIStatus string       `json:"APIStatus"`
	APIError  string       `json:"APIError"`
}

type LoginDetails struct {
	Username string `json:"username" valid:"required,trim,lower"`
	Password string `json:"password" valid:"required,trim,lower"`
}

type Tasks struct {
	ID          int       `gorm:"column:ID;primaryKey;autoIncrement" json:"ID"`
	Title       string    `gorm:"column:Title;type:varchar(200);not null" json:"Title"`
	Description string    `gorm:"column:Description;type:text" json:"Description"`
	Status      string    `gorm:"column:Status;type:status_enum;default:'Pending';not null" json:"Status"`
	Duedate     time.Time `gorm:"column:Duedate" json:"Duedate"`
	CreatedAt   time.Time `gorm:"column:CreatedAt;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CreatedAt"`
	UpdatedAt   time.Time `gorm:"column:UpdatedAt;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UpdatedAt"`
	CreatedBy   string    `gorm:"column:CreatedBy;type:varchar(100);not null" json:"CreatedBy"`
	UpdatedBy   string    `gorm:"column:UpdatedBy;type:varchar(100);not null" json:"UpdatedBy"`
}

type Response struct {
	APIStatus string `json:"APIStatus"`
	APIError  string `json:"APIError"`
	Tasks
}
