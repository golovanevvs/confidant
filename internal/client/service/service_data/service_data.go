package service_data

import (
	"context"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/client/service/service_security"
)

type ITransportData interface {
}

type IRepositoryData interface {
	AddNote(ctx context.Context, data *model.NoteEnc) (err error)
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
