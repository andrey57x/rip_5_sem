package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/ds"
	minioInclude "Backend/internal/app/minio"

	"github.com/gin-gonic/gin"
)

func (r *Repository) GetReactions() ([]ds.Reaction, error) {
	var reactions []ds.Reaction
	err := r.db.Where("is_delete = false").Find(&reactions).Error
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	if err != nil {
		return nil, ErrorNotFound
	}
	if len(reactions) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return reactions, nil
}

func (r *Repository) GetReaction(id int) (ds.Reaction, error) {
	reaction := ds.Reaction{}
	if id < 0 {
		return ds.Reaction{}, errors.New("invalid id, it must be >= 0")
	}
	sub := r.db.Where("id = ? and is_delete = ?", id, false).Find(&reaction)
	if sub.Error != nil {
		return ds.Reaction{}, sub.Error
	}
	if sub.RowsAffected == 0 {
		return ds.Reaction{}, ErrorNotFound
	}
	err := sub.First(&reaction).Error
	if err != nil {
		return ds.Reaction{}, err
	}
	return reaction, nil
}

func (r *Repository) GetReactionsByTitle(title string) ([]ds.Reaction, error) {
	var reactions []ds.Reaction
	err := r.db.Where("title ILIKE ? and is_delete = ?", "%"+title+"%", false).Find(&reactions).Error
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (r *Repository) CreateReaction(reactionJSON apitypes.ReactionJSON) (ds.Reaction, error) {
	reaction := apitypes.ReactionFromJSON(reactionJSON)
	if reaction.ConversationFactor <= 0 {
		return ds.Reaction{}, errors.New("invalid conversation factor")
	}
	err := r.db.Create(&reaction).First(&reaction).Error
	if err != nil {
		return ds.Reaction{}, err
	}
	return reaction, nil
}

func (r *Repository) ChangeReaction(id int, reactionJSON apitypes.ReactionJSON) (ds.Reaction, error) {
	reaction := ds.Reaction{}
	if id < 0 {
		return ds.Reaction{}, errors.New("invalid id, it must be >= 0")
	}
	err := r.db.Where("id = ? and is_delete = ?", id, false).First(&reaction).Error
	if err != nil {
		return ds.Reaction{}, ErrorNotFound
	}
	if reactionJSON.ConversationFactor <= 0 {
		return ds.Reaction{}, errors.New("invalid conversation factor")
	}
	err = r.db.Model(&reaction).Updates(apitypes.ReactionFromJSON(reactionJSON)).Error
	if err != nil {
		return ds.Reaction{}, err
	}
	return reaction, nil
}

func (r *Repository) DeleteReaction(id int) error {
	reaction := ds.Reaction{}
	if id < 0 {
		return errors.New("invalid id, it must be >= 0")
	}

	err := r.db.Where("id = ? and is_delete = ?", id, false).First(&reaction).Error
	if err != nil {
		return ErrorNotFound
	}
	if reaction.ImgLink != "" {
		err = minioInclude.DeleteObject(context.Background(), r.mc, minioInclude.GetImgBucket(), reaction.ImgLink)
		if err != nil {
			return err
		}
	}

	err = r.db.Model(&ds.Reaction{}).Where("id = ?", id).Update("is_delete", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddReactionToCalculation(calculationID int, reactionID int) error {
	var reaction ds.Reaction
	if err := r.db.First(&reaction, reactionID).Error; err != nil {
		return ErrorNotFound
	}

	var calculation ds.MassCalculation
	if err := r.db.First(&calculation, calculationID).Error; err != nil {
		return err
	}
	reactionCalculation := ds.ReactionCalculation{}
	result := r.db.Where("reaction_id = ? and calculation_id = ?", reactionID, calculationID).Find(&reactionCalculation)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return nil
	}
	return r.db.Create(&ds.ReactionCalculation{
		ReactionID:    reactionID,
		CalculationID: calculationID,
	}).Error
}

func CalculateMass(mass float32, conversationFactor, outputKoef float32) (float32, error) {
	if conversationFactor == 0 {
		return 0, errors.New("invalid conversation factor")
	}
	if outputKoef == 0 || outputKoef > 1 {
		return 0, errors.New("invalid output koeficient")
	}
	return mass / conversationFactor / outputKoef, nil
}

func (r *Repository) UploadImage(ctx *gin.Context, reactionID int, file *multipart.FileHeader) (ds.Reaction, error) {
	reaction, err := r.GetReaction(reactionID)
	if err != nil {
		return ds.Reaction{}, ErrorNotFound
	}

	fileName, err := minioInclude.UploadImage(ctx, r.mc, minioInclude.GetImgBucket(), file, reactionID)
	if err != nil {
		return ds.Reaction{}, err
	}

	reaction.ImgLink = fileName
	err = r.db.Save(&reaction).Error
	if err != nil {
		return ds.Reaction{}, err
	}
	return reaction, nil
}

func (r *Repository) GetModeratorAndCreatorLogin(calculation ds.MassCalculation) (string, string, error) {
	var creator ds.User
	var moderator ds.User

	err := r.db.Where("uuid = ?", calculation.CreatorID).First(&creator).Error
	if err != nil {
		return "", "", err
	}

	var moderatorLogin string
	if calculation.ModeratorID.Valid {
		err = r.db.Where("uuid = ?", calculation.ModeratorID).First(&moderator).Error
		if err != nil {
			return "", "", err
		}
		moderatorLogin = moderator.Login
	}

	return creator.Login, moderatorLogin, nil
}
