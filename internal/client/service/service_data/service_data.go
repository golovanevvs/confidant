package service_data

import (
	"context"
	"time"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/client/service/service_security"
)

type ITransportData interface {
}

type IRepositoryData interface {
	GetDataTitles(ctx context.Context, groupID int) (dataTitles [][]byte, err error)
	GetDataTypes(ctx context.Context, groupID int) (dataTypes []string, err error)
	GetDataIDAndType(ctx context.Context, groupID int, dataTitle string) (dataID int, dataType string, err error)
	AddNote(ctx context.Context, data model.NoteEnc) (err error)
	GetNote(ctx context.Context, dataID int) (data model.NoteEnc, err error)
	AddPass(ctx context.Context, data model.PassEnc) (err error)
	GetPass(ctx context.Context, dataID int) (data model.PassEnc, err error)
	AddCard(ctx context.Context, data model.CardEnc) (err error)
	GetCard(ctx context.Context, dataID int) (data model.CardEnc, err error)
	AddFile(ctx context.Context, data model.FileEnc) (err error)
	GetFile(ctx context.Context, dataID int) (data model.FileEnc, err error)
	GetFileForSave(ctx context.Context, dataID int) (data []byte, err error)
	GetDataIDs(ctx context.Context, groupServerIDs []int) (dataServerIDs []int, dataNoServerIDs []int, err error)
	GetDataDates(ctx context.Context, dataIDs []int) (dataDates map[int]time.Time, err error)
	SaveDatas(ctx context.Context, datas []model.Data) (err error)
	SaveDataFile(ctx context.Context, dataID int, file []byte) (err error)
	GetDatasByIDs(ctx context.Context, dataIDs []int) (datas []model.Data, err error)
	UpdateDataIDsOnServer(ctx context.Context, newDataIDs map[int]int) (err error)
	GetDataFile(ctx context.Context, dataID int) (idOnServer int, file []byte, err error)
}

type IServiceSecurity interface {
	Encrypt(data []byte) (encryptedData []byte, err error)
	Decrypt(data []byte) (decryptedData []byte, err error)
	EncryptFile(filepath string) (encryptedFile []byte, err error)
	DecryptFile(data []byte, filepath string) (err error)
}

type ServiceData struct {
	tr ITransportData
	rp IRepositoryData
	ss IServiceSecurity
}

func New(tr ITransportData, rp IRepositoryData) *ServiceData {
	return &ServiceData{
		tr: tr,
		rp: rp,
		ss: service_security.New(),
	}
}
