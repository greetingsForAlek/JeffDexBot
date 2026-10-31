package models

// character struct
type Character struct {
	ID int `json:"id"` // the charcater id
	Name string `json:"name"` // the character's name
	ImageURL string `json:"image_url"` // the image url
	Rarity string `json:"rarity"` // the character's rarity
}

// yeah thats all.