package go_gdrive

import (
	"errors"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"io"
)

type gdrive struct {
	googleDriveService *drive.Service
	spreadSheetService *sheets.Service
}

type gdriveInterface interface {
	CreateDirectory(Name string, parentDirectory string)(*drive.File,error)
	UploadFile(mimeType string, fileName string, parentDirectory string,content io.Reader)(*drive.File,error)
	CreateNewFile(mimeType string, fileName string, parentDirectory string)(*drive.File,error)
	GetDirectoryList()([]*drive.File,error)
	GetFilesWithQuery(query string)([]*drive.File,error)
}

func (g *googleApi)GetDriveServices()(gdriveInterface,error){
	service, err := drive.New(g.client)
	return &gdrive{googleDriveService:service},err
}

func (g *googleApi)GetSpreadSheetServices()(SpreadSheetInterface,error){
	service,err := sheets.New(g.client)
	return &spreadSheet{spreadSheetService:service},err
}

func (g *gdrive)UploadFile(mimeType string, fileName string, parentDirectory string,content io.Reader)(*drive.File,error){
	f := &drive.File{
		MimeType: mimeType,
		Name:     fileName,
		Parents:  []string{parentDirectory},
	}
	file, err := g.googleDriveService.Files.Create(f).Media(content).Do()
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not Upload file: %v\n", err))
		return nil, err
	}
	return file, err
}
func (g *gdrive)CreateNewFile(mimeType string, fileName string, parentDirectory string)(*drive.File,error){
	f := &drive.File{
		MimeType: mimeType,
		Name:     fileName,
		Parents:  []string{parentDirectory},
	}
	file, err := g.googleDriveService.Files.Create(f).Do()
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not create file: %v\n", err))
		return nil, err
	}
	return file, err
}
func (g *gdrive)CreateDirectory(Name string, parentDirectory string)(*drive.File,error){
	d := &drive.File{
		Name:     Name,
		MimeType: GOOGLE_DIRECTORY,
		Parents:  []string{parentDirectory},
	}
	file, err := g.googleDriveService.Files.Create(d).Do()
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not create directory: %v\n", err))
		return nil, err
	}
	return file, nil
}

func(g *gdrive)GetDirectoryList()([]*drive.File,error){
	listfile := g.googleDriveService.Files.List()
	filesList,err:= listfile.Q(GOOGLE_DIRECTORY_Q_MIME_TYPE).Do()
	if err != nil {
		return nil, err
	}
	return filesList.Files, err
}
func(g *gdrive)GetFilesWithQuery(query string)([]*drive.File,error){
	listfile := g.googleDriveService.Files.List()
	filesList,err:= listfile.Q(query).Spaces("drive").Do()
	if err != nil {
		return nil, err
	}
	return filesList.Files, err
}