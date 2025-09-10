package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Reaction struct {
	ID                 int
	Title              string
	Reagent            string
	Product            string
	ConversationFactor float32
	ImgLink            string
	Description        string
}

func (r *Repository) GetReactions() ([]Reaction, error) {
	reactions := []Reaction{
		{
			ID:                 1,
			Title:              "Получение диоксида серы",
			Reagent:            "Пирит",
			Product:            "Диоксид серы",
			ConversationFactor: 1.067,
			ImgLink:            "so2.png",
			Description:        "Диоксид серы получают обжигом пирита на воздухе: сульфиды окисляются до диоксида серы и образуется оксид железа. Газ на промышленных установках очищают от пыли и загрязнений и используют для производства триоксида серы и серной кислоты.",
		},
		{
			ID:                 2,
			Title:              "Контактный метод получения серного ангидрида",
			Reagent:            "Диоксид серы",
			Product:            "Серный ангидрид",
			ConversationFactor: 1.25,
			ImgLink:            "so3.png",
			Description:        "Диоксид серы окисляют кислородом на катализаторе при повышенной температуре, превращая его в триоксид серы. Реакция экзотермична, протекает в несколько ступеней для увеличения выхода, и используется в контактном методе получения серного ангидрида.",
		},
		{
			ID:                 3,
			Title:              "Производство серной кислоты",
			Reagent:            "Серный ангидрид",
			Product:            "Серная кислота",
			ConversationFactor: 1.225,
			ImgLink:            "h2so4.png",
			Description:        "Триоксид серы не смешивают напрямую с водой, поэтому его сначала растворяют в концентрированной серной кислоте, образуя олеум. Затем олеум осторожно разбавляют водой, получая серную кислоту нужной концентрации для промышленного применения.",
		},
		{
			ID:                 4,
			Title:              "Синтез аммиака",
			Reagent:            "Азот",
			Product:            "Аммиак",
			ConversationFactor: 1.214,
			ImgLink:            "nh3.png",
			Description:        "Азот и водород при высоком давлении и температуре проходят через железный катализатор, образуя аммиак. Реакция обратима, поэтому непрореагировавшие газы возвращают обратно. Этот процесс известен как метод Габера-Боша для синтеза аммиака.",
		},
		{
			ID:                 5,
			Title:              "Производство метанола",
			Reagent:            "Монооксид углерода",
			Product:            "Метанол",
			ConversationFactor: 1.143,
			ImgLink:            "ch3oh.png",
			Description:        "Метанол получают гидрированием синтез-газа (смеси оксидов углерода и водорода) на медных катализаторах при повышенном давлении и умеренных температурах. В промышленности используют синтез-газ, получаемый паровым риформингом природного газа или газификацией угля и биомассы.",
		},
	}
	if len(reactions) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return reactions, nil
}

func (r *Repository) GetReaction(id int) (Reaction, error) {
	reactions, err := r.GetReactions()
	if err != nil {
		return Reaction{}, err
	}

	for _, reaction := range reactions {
		if reaction.ID == id {
			return reaction, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Reaction{}, fmt.Errorf("заказ не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}

func (r *Repository) GetReactionsByTitle(title string) ([]Reaction, error) {
	reactions, err := r.GetReactions()
	if err != nil {
		return []Reaction{}, err
	}

	var result []Reaction
	for _, reaction := range reactions {
		if strings.Contains(strings.ToLower(reaction.Title), strings.ToLower(title)) {
			result = append(result, reaction)
		}
	}

	return result, nil
}

func (r *Repository) GetCalculationReactions(id int) ([]Reaction, error) {
	manyToMany := map[int][]int{
		1: {1, 2},
	}

	reactionIds, ok := manyToMany[id]
	if !ok {
		return []Reaction{}, nil
	}

	reactions, err := r.GetReactions()
	if err != nil {
		return []Reaction{}, err
	}

	var result []Reaction
	for _, reaction := range reactions {
		for _, testId := range reactionIds {
			if reaction.ID == testId {
				result = append(result, reaction)
				break
			}
		}
	}
	return result, nil
}

func (r *Repository) CurrentCalculation() (int, int, error) {
	// функция возвращает количество услуг в корзине и id текущей корзины
	var id = 1

	reactions, err := r.GetCalculationReactions(id)
	if err != nil {
		return 0, 0, err
	}

	return len(reactions), id, nil
}
