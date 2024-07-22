package handlers

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/kylerequez/go-upload-example/src/models"
	"github.com/kylerequez/go-upload-example/src/repositories"
	"github.com/kylerequez/go-upload-example/src/views"
)

type UploadHandler struct {
	ur *repositories.UploadRepository
}

func NewUploadHandler(ur *repositories.UploadRepository) *UploadHandler {
	return &UploadHandler{
		ur: ur,
	}
}

// Route initializer
func (uh *UploadHandler) InitRoutes(server *fiber.App) {
	// Views routes
	views := server.Group("")
	views.Get("/", uh.GetUploadPage)

	// Api routes
	api := server.Group("/api/v1/upload")
	api.Get("/:id", uh.GetFile)
	api.Post("", uh.UploadFile)
	api.Delete("/:id", uh.DeleteFile)
}

// Get the upload page
func (uh *UploadHandler) GetUploadPage(c fiber.Ctx) error {
	msg := views.UploadMessages{
		Message: "",
		Errors:  make(map[string]string),
	}
	pageTitle := "Go Upload Example"

	if c.Method() == fiber.MethodGet {
		files, err := uh.ur.GetAllUploads()
		if err != nil {
			msg.Errors["get-uploads-err"] = err.Error()
			return Render(c, views.UploadView(pageTitle, msg, nil, nil))
		}

		return Render(c, views.UploadView(pageTitle, msg, &files, nil))
	}
	return nil
}

// Retrieves the file and displays it
func (uh *UploadHandler) GetFile(c fiber.Ctx) error {
	msg := views.UploadMessages{
		Message: "",
		Errors:  make(map[string]string),
	}

	if c.Method() == fiber.MethodGet {
		id := c.Params("id")
		if id == "" {
			msg.Errors["get-file-err"] = "id must be valid"
		}

		oid, err := uuid.Parse(id)
		if err != nil {
			msg.Errors["get-file-err"] = err.Error()
			return Render(c, views.FileShowcase(nil, msg))
		}

		file, err := uh.ur.GetUploadById(oid)
		if err != nil {
			msg.Errors["get-file-err"] = err.Error()
			return Render(c, views.FileShowcase(nil, msg))
		}

		return Render(c, views.FileShowcase(file, msg))
	}
	return nil
}

// Upload the file to the /uploads folder and saves the info inside the database
func (uh *UploadHandler) UploadFile(c fiber.Ctx) error {
	msg := views.UploadMessages{
		Message: "",
		Errors:  make(map[string]string),
	}

	if c.Method() == fiber.MethodPost {
		files, err := uh.ur.GetAllUploads()
		if err != nil {
			msg.Errors["get-uploads-err"] = err.Error()
			return Render(c, views.Uploads(&files, nil, msg))
		}

		file, err := c.FormFile("upload")
		if err != nil {
			msg.Errors["upload-err"] = err.Error()
			return Render(c, views.Uploads(&files, nil, msg))
		}

		fileSpl := strings.Split(file.Filename, ".")
		if len(fileSpl) < 2 {
			msg.Errors["upload-err"] = "file must have an extension"
			return Render(c, views.Uploads(&files, nil, msg))
		}

		newFile := models.File{
			Name: fileSpl[0],
			Size: file.Size,
			Type: fileSpl[1],
		}

		err = ValidateFile(newFile, uh.ur)
		if err != nil {
			msg.Errors["upload-err"] = err.Error()
			return Render(c, views.Uploads(&files, nil, msg))
		}

		err = uh.ur.CreateFile(newFile)
		if err != nil {
			msg.Errors["upload-err"] = err.Error()
			return Render(c, views.Uploads(&files, nil, msg))
		}

		if err := c.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename)); err != nil {
			msg.Errors["upload-err"] = err.Error()
			return Render(c, views.Uploads(&files, nil, msg))
		}

		files, err = uh.ur.GetAllUploads()
		if err != nil {
			msg.Errors["get-uploads-err"] = err.Error()
			return Render(c, views.Uploads(&files, nil, msg))
		}

		msg.Message = "file has been uploaded successfully"
		return Render(c, views.Uploads(&files, nil, msg))
	}
	return nil
}

// Deletes the file inside the /uploads directory and deletes it from the database
func (uh *UploadHandler) DeleteFile(c fiber.Ctx) error {
	msg := views.UploadMessages{
		Message: "",
		Errors:  make(map[string]string),
	}

	if c.Method() == fiber.MethodDelete {
		files, err := uh.ur.GetAllUploads()
		if err != nil {
			msg.Errors["get-uploads-err"] = err.Error()
			return Render(c, views.Uploads(&files, nil, msg))
		}

		id := c.Params("id")
		if id == "" {
			msg.Errors["delete-upload-err"] = "id must not be empty"
			return Render(c, views.Uploads(&files, nil, msg))
		}

		oid, err := uuid.Parse(id)
		if err != nil {
			msg.Errors["delete-upload-err"] = err.Error()
			return Render(c, views.Uploads(&files, nil, msg))
		}

		file, err := uh.ur.GetUploadById(oid)
		if err != nil {
			msg.Errors["delete-upload-err"] = err.Error()
			return Render(c, views.UploadsContainer(&files, nil, msg))
		}

		if err := DeleteFile(file); err != nil {
			msg.Errors["delete-upload-err"] = err.Error()
			return Render(c, views.Uploads(&files, nil, msg))
		}

		if err := uh.ur.DeleteUploadById(oid); err != nil {
			msg.Errors["delete-upload-err"] = err.Error()
			return Render(c, views.Uploads(&files, nil, msg))
		}

		files, err = uh.ur.GetAllUploads()
		if err != nil {
			msg.Errors["get-uploads-err"] = err.Error()
			return Render(c, views.Uploads(&files, nil, msg))
		}

		msg.Message = "successfully deleted a file"
		return Render(c, views.Uploads(&files, nil, msg))
	}
	return nil
}

func ValidateFile(file models.File, ur *repositories.UploadRepository) error {
	if len(file.Name) > 50 {
		return errors.New("file name must not exceed 50 characters")
	}

	if len(file.Type) > 10 {
		return errors.New("file type must not exceed 50 characters")
	}

	isAccepted := false
	acceptedTypes := []string{"jpeg", "jpg", "png", "gif", "bmp", "tiff", "pdf"}
	for i := range acceptedTypes {
		if acceptedTypes[i] == file.Type {
			isAccepted = true
			break
		}
	}
	if !isAccepted {
		return errors.New("file extension is not accepted")
	}

	if file.Size > 200000000 {
		return errors.New("file is too large (must not exceed 25 MB)")
	}

	isExists, err := ur.GetUploadByName(file.Name)
	if err != nil {
		return err
	}

	if isExists != nil {
		return errors.New("file exists!")
	}
	return nil
}

func DeleteFile(file *models.File) error {
	store := "./uploads"
	filepath := fmt.Sprintf("%s/%s.%s", store, file.Name, file.Type)

	err := os.Remove(filepath)
	if err != nil {
		return err
	}
	return nil
}
