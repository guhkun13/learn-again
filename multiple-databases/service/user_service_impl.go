package service

import (
	"time"

	"github.com/rs/zerolog/log"

	config "github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/internal/dto"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
	"github.com/guhkun13/learn-again/multiple-databases/internal/repository"
	lib "github.com/guhkun13/learn-again/multiple-databases/lib"
)

type UserServiceImpl struct {
	Env                   *config.EnvironmentVariable
	CompressionRepository repository.CompressionRepository
}

func NewUserServiceImpl(
	env *config.EnvironmentVariable,
	compressionRepo repository.CompressionRepository,

) UserService {
	return &UserServiceImpl{
		Env:            env,
		UserRepository: userRepository,
	}
}

func (s *UserServiceImpl) FindAll(filter dto.FilterFindCompression) (res model.FindCompressionWithPagination, err error) {
	res, err = s.CompressionRepository.FindAll(filter)
	if err != nil {
		log.Error().Err(err).Msg("CompressionRepository.FindAll failed")
	}

	return
}

func (s *CompressionServiceImpl) FindById(id string) (res model.Compression, err error) {
	res, err = s.CompressionRepository.FindById(id)
	if err != nil {
		log.Error().Err(err).Msg("CompressionRepository.FindById failed")
	}

	return
}

func (s *CompressionServiceImpl) DeleteById(id string) (err error) {
	err = s.CompressionRepository.RemoveById(id)
	if err != nil {
		log.Error().Err(err).Msg("CompressionRepository.RemoveById failed")
	}

	return
}

func (s *CompressionServiceImpl) UploadCompress(req dto.UploadCompressRequest) (compressionId string, err error) {
	log.Info().Msg("CompressionServiceImpl.CreateUploadCompression")

	// insert to database compression
	newItem := model.Compression{
		Id:              lib.GenerateId(),
		UserId:          req.UserId,
		Repetition:      lib.OneTime,
		CreatedAt:       time.Now(),
		SourcePath:      req.SourcePath,
		DestinationPath: req.DestinationPath,
	}

	err = s.CompressionRepository.Insert(newItem)
	if err != nil {
		return
	}

	compressionId = newItem.Id

	return
}

func (s *CompressionServiceImpl) RegisterAutoCompress(req dto.RegisterAutoCompressRequest) (res model.Compression, err error) {
	res = model.Compression{
		Id:              lib.GenerateId(),
		UserId:          req.UserId,
		CreatedAt:       time.Now(),
		Repetition:      lib.Periodic,
		SourcePath:      req.SourcePath,
		DestinationPath: req.DestinationPath,
	}

	err = s.CompressionRepository.Insert(res)
	if err != nil {
		return
	}
	return
}

func (s *CompressionServiceImpl) UpdateAutoCompress(compression model.Compression, req dto.UpdateAutoCompressRequest) error {
	err := s.CompressionRepository.UpdateAutoCompressPath(compression, req)

	return err
}

// Multiple

func (s *CompressionServiceImpl) FindByIds(ids string) (res model.Compressions, err error) {
	ids = lib.RemoveWhiteSpaces(ids)
	arrIds := lib.StringToList(ids, ",")

	res, err = s.CompressionRepository.FindByIds(arrIds)

	return
}

func (s *CompressionServiceImpl) DeleteByIds(ids []string) (res int, err error) {
	res, err = s.CompressionRepository.RemoveByIds(ids)

	return
}
