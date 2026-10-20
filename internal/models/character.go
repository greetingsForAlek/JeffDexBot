package models

type Character struct {
	ID int `json:"id"`
	Name string `json:"name"`
	ImageURL string `json:"image_url"`
	Rarity string `json:"rarity"`
}