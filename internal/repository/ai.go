package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"

	"thomas.vn/apartment_service/internal/domain/model"
	"thomas.vn/apartment_service/internal/domain/service"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type AiRepository struct {
	logger  *xlogger.Logger
	client  *xhttp.HTTPClient
	fileSvc service.FileService
	url     string
}

func NewAiRepository(logger *xlogger.Logger, client *xhttp.HTTPClient, fileSvc service.FileService, url string) *AiRepository {
	return &AiRepository{
		logger:  logger,
		client:  client,
		fileSvc: fileSvc,
		url:     url,
	}
}

func (r *AiRepository) VerifyCV(attachFile string, jobDesc string) (int, model.VerifyResponse, error) {
	reqBody := map[string]interface{}{
		"jd_text":   jobDesc,
		"file_urls": []string{attachFile},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		r.logger.Error("Failed to marshal request body", xlogger.Error(err))
		return 0, model.VerifyResponse{}, err
	}

	opts := &xhttp.ClientRequestOptions{
		Method: xhttp.MethodPost,
		URL:    r.url + "v2/scan-cv",
		Body:   bytes.NewBuffer(jsonBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	var responseBody model.VerifyResponse
	err = r.client.SendAndParse(opts, &responseBody)
	if err != nil {
		r.logger.Error("Failed to send request", xlogger.Error(err))
		return 0, model.VerifyResponse{}, err
	}

	r.logger.Info("Response", xlogger.Object("result", responseBody))

	return 1, responseBody, nil
}

// has download file and push to api call ai service
func (r *AiRepository) VerifyCVDownload(attachFile string, jobDesc string) (int, model.VerifyResponse, error) {
	file, err := r.fileSvc.Download(attachFile)
	if err != nil {
		r.logger.Error("Failed to download file with url", xlogger.String("url", attachFile), xlogger.Error(err))
		return 0, model.VerifyResponse{}, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="files"; filename="%s"`, file.FileName))
	h.Set("Content-Type", file.ContentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		r.logger.Error("Failed to create part", xlogger.Error(err))
		return 0, model.VerifyResponse{}, fmt.Errorf("failed to create part: %w", err)
	}

	if _, err := io.Copy(part, bytes.NewReader(file.Content)); err != nil {
		r.logger.Error("Failed to copy file content", xlogger.Error(err))
		return 0, model.VerifyResponse{}, err
	}

	if err := writer.WriteField("jd_text", jobDesc); err != nil {
		r.logger.Error("Failed to write field", xlogger.Error(err))
		return 0, model.VerifyResponse{}, err
	}
	if err := writer.Close(); err != nil {
		r.logger.Error("Failed to close writer", xlogger.Error(err))
		return 0, model.VerifyResponse{}, err
	}

	opts := &xhttp.ClientRequestOptions{
		Method: xhttp.MethodPost,
		URL:    r.url + "v1/scan-cv",
		Body:   body,
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
			"Accept":       "application/json",
		},
	}

	var responseBody model.VerifyResponse
	err = r.client.SendAndParse(opts, &responseBody)
	if err != nil {
		r.logger.Error("Failed to send request", xlogger.Error(err))
		return 0, model.VerifyResponse{}, err
	}

	r.logger.Info("Response", xlogger.Object("result", responseBody))

	return 1, responseBody, nil
}
