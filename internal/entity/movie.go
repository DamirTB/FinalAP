package entity

import (
	"time"
	"errors"
)

type Movie struct {
	ID        int64     `json:"id"`                       
	CreatedAt time.Time `json:"-"`                        
	Title     string    `json:"title"`                    
	Year      int32     `json:"year,omitempty"`           
	Runtime   int32     `json:"runtime,omitempty,string"` 
	Genres    []string  `json:"genres,omitempty"`         
	Version   int32     `json:"version"`                  
}

var (
	ErrRecordNotFound = errors.New("record (row, entry) not found")
	ErrEditConflict   = errors.New("edit conflict")
)

