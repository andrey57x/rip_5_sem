package repository

import (
	"fmt"

	"Backend/internal/app/ds"

	"github.com/sirupsen/logrus"
)

func (r *Repository) GetReactions() ([]ds.Reaction, error) {
	var reactions []ds.Reaction
	err := r.db.Where("is_delete = false").Find(&reactions).Error
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	if err != nil {
		return nil, err
	}
	if len(reactions) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return reactions, nil
}

func (r *Repository) GetReaction(id int) (ds.Reaction, error) {
	reaction := ds.Reaction{}
	err := r.db.Where("id = ? and is_delete = ?", id, false).First(&reaction).Error
	if err != nil {
		return ds.Reaction{}, err
	}
	return reaction, nil
}

func (r *Repository) GetReactionsByTitle(title string) ([]ds.Reaction, error) {
	var reactions []ds.Reaction
	err := r.db.Where("name ILIKE ? and is_delete = ?", "%"+title+"%", false).Find(&reactions).Error
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (r *Repository) GetCalculationReactions(id int) ([]ds.Reaction, ds.Calculation, error) {
	var calculation ds.Calculation
	creatorID := 1
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	err := r.db.Model(&ds.Calculation{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").First(&calculation).Error
	if err != nil {
		return []ds.Reaction{}, ds.Calculation{}, err
	}
	
	var reactions []ds.Reaction
	sub := r.db.Table("reaction_calculations").Select("reaction_id").Where("calculation_id = ?", calculation.ID)
	err = r.db.Where("id IN (?)", sub).Find(&reactions).Error
	if err != nil {
		return []ds.Reaction{}, ds.Calculation{}, err
	}

	return reactions, calculation, nil
}

// GetCartCount для получения количества услуг в заявке
func (r *Repository) GetCartCount() int {
	var calculationID uint
	var count int64
	creatorID := 1
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	err := r.db.Model(&ds.Calculation{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("id").First(&calculationID).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.ReactionCalculation{}).Where("calculation_id = ?", calculationID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting reactions in reaction_calculations:", err)
	}

	return int(count)
}
