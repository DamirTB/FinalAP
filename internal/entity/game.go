package entity

import (
	"time"
	"errors"
)

type Game struct {
	ID        	int64     `json:"id"`                       
	CreatedAt 	time.Time `json:"-"`                        
	Name     	string    `json:"name"`                    
	Price      	int32     `json:"price"`           
	Genres    	[]string  `json:"genres,omitempty"`         
	Version   	int32     `json:"version"`                  
}

var (
	ErrRecordNotFound = errors.New("record (row, entry) not found")
	ErrEditConflict   = errors.New("edit conflict")
)
