package apitypes

import "Backend/internal/app/ds"

// ReactionJSON model
// @Description Model for reaction
// @Tags reactions
type ReactionJSON struct {
	ID                 int     `json:"id"`
	Title              string  `json:"title"`
	Reagent            string  `json:"reagent"`
	Product            string  `json:"product"`
	ConversationFactor float32 `json:"conversation_factor"`
	ImgLink            string  `json:"img_link"`
	Description        string  `json:"description"`
}

func ReactionToJSON(r ds.Reaction) ReactionJSON {
	return ReactionJSON{
		ID:                 r.ID,
		Title:              r.Title,
		Reagent:            r.Reagent,
		Product:            r.Product,
		ConversationFactor: r.ConversationFactor,
		ImgLink:            "http://localhost:9000/img/" + r.ImgLink,
		Description:        r.Description,
	}
}

func ReactionFromJSON(r ReactionJSON) ds.Reaction {
	return ds.Reaction{
		Title:              r.Title,
		Reagent:            r.Reagent,
		Product:            r.Product,
		ConversationFactor: r.ConversationFactor,
		Description:        r.Description,
		IsDelete: false,
	}
}
