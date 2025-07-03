package models

import "time"

type CreateTask struct {
	ID          int    `json:"ID"`
	Title       string `json:"Title" validate:"required" valid:"trim"`
	Description string `json:"Description" valid:"trim"`
	Status      string `json:"Status" validate:"oneof='Pending' 'In Progress' 'Completed'" valid:"trim"`
	DueDate     string `json:"DueDate" valid:"trim"`
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
	SortBy        string `validate:"omitempty,oneof=createdat updatedat duedate" schema:"sort_by"`
	SortOrder     string `validate:"omitempty,oneof=asc desc" schema:"sort_order"`
}

type GetAllTaskResp struct {
	Tasks     []Tasks `json:"tasks"`
	Page      int     `json:"page"`
	Limit     int     `json:"limit"`
	Total     int     `json:"total"`
	APIStatus string  `json:"APIStatus"`
	APIError  string  `json:"APIError"`
}

type LoginDetails struct {
	Username string `json:"username" valid:"required,trim,lower"`
	Password string `json:"password" valid:"required,trim,lower"`
}

type Tasks struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"ID" valid:"trim"`
	Title       string    `gorm:"column:title;type:varchar(200);not null" json:"Title" valid:"trim"`
	Description string    `gorm:"column:description;type:text" json:"Description" valid:"trim"`
	Status      string    `gorm:"column:status;type:status_enum;default:'Pending';not null" json:"Status" valid:"trim"`
	Duedate     time.Time `gorm:"column:duedate" json:"Duedate" valid:"trim"`
	CreatedAt   time.Time `gorm:"column:createdat;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CreatedAt" valid:"trim"`
	UpdatedAt   time.Time `gorm:"column:updatedat;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UpdatedAt" valid:"trim"`
	CreatedBy   string    `gorm:"column:createdby;type:varchar(100);not null" json:"-" valid:"trim"`
	UpdatedBy   string    `gorm:"column:updatedby;type:varchar(100);not null" json:"-" valid:"trim"`
}

type Response struct {
	APIStatus string `json:"APIStatus"`
	APIError  string `json:"APIError"`
	Tasks
}
