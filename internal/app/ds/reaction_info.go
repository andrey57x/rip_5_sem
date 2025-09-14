package ds

type ReactionInfo struct {
	ID                 int
	Title              string
	Reagent            string
	Product            string
	ConversationFactor float32
	ImgLink            string

	OutputMass float32
	InputMass  float32
}
