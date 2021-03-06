package usecase

import (
	"net/http"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type CreateVoiceParam struct {
	LessonID    int64   `json:"lessonID"`
	ElapsedTime float32 `json:"elapsedTime"`
	DurationSec float32 `json:"durationSec"`
}

func GetVoice(request *http.Request, lessonID int64, id int64) (domain.Voice, error) {
	ctx := request.Context()

	var voice domain.Voice
	voice.ID = id

	if _, err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return voice, err
	}

	if err := domain.GetVoice(ctx, lessonID, &voice); err != nil {
		return voice, err
	}

	return voice, nil
}

func GetVoices(request *http.Request, lessonID int64) ([]domain.Voice, error) {
	ctx := request.Context()

	var voices []domain.Voice

	if _, err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return nil, err
	}

	if err := domain.GetVoices(ctx, lessonID, &voices); err != nil {
		return nil, err
	}

	return voices, nil
}

// CreateVoiceAndBlankFile creates Voice and blank files of mp3 and wav.
func CreateVoiceAndBlankFile(request *http.Request, params *CreateVoiceParam) (infrastructure.SignedURL, error) {
	ctx := request.Context()

	var response infrastructure.SignedURL

	userID, err := currentUserAccessToLesson(ctx, request, params.LessonID)
	if err != nil {
		return response, err
	}

	voice := domain.Voice{
		UserID:      userID,
		ElapsedTime: params.ElapsedTime,
		DurationSec: params.DurationSec,
	}

	if err = domain.CreateVoice(ctx, params.LessonID, &voice); err != nil {
		return response, err
	}

	lessonID := strconv.FormatInt(params.LessonID, 10)
	voiceID := strconv.FormatInt(voice.ID, 10)

	mp3FileRequest := infrastructure.FileRequest{
		ID:          voiceID,
		Entity:      "voice",
		Extension:   "mp3",
		ContentType: "audio/mpeg",
	}

	filePath := lessonID + "/" + voiceID
	mp3URL, err := infrastructure.CreateBlankFileToGCS(ctx, filePath, "voice", mp3FileRequest)
	if err != nil {
		return response, err
	}

	response.FileID = voiceID
	response.SignedURL = mp3URL

	return response, nil
}
