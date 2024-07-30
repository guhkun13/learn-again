package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"kecilin-id/local-mks-api/config"
	"kecilin-id/local-mks-api/lib"

	"kecilin-id/local-mks-api/internal/dto"
	"kecilin-id/local-mks-api/internal/helper"
	"kecilin-id/local-mks-api/internal/service"
)

type CompressionHandlerImpl struct {
	Env                *config.EnvironmentVariable
	CompressionService service.CompressionService
	ScheduleService    service.ScheduleService
	MonitoringService  service.MonitoringService
}

func NewCompressionHandlerImpl(env *config.EnvironmentVariable,
	compressionService service.CompressionService,
	scheduleService service.ScheduleService,
	monitoringService service.MonitoringService,
) CompressionHandler {
	return &CompressionHandlerImpl{
		Env:                env,
		CompressionService: compressionService,
		ScheduleService:    scheduleService,
		MonitoringService:  monitoringService,
	}
}

// FindAll godoc
// @Summary     	Find all compression
// @Description 	-
// @Security 		BearerAuth
// @Tags        	Compression
// @Accept      	json
// @Produce     	json
// @Param 		 	limit 							query 		number 		false 		"limit result"
// @Param 		 	page 							query 		number 		false 		"page"
// @Param 		 	sortBy 							query 		string 		false 		"sort by"
// @Param 		 	sortOrder 						query 		string 		true 		"sort order"		Enums(desc,asc)
// @Param 		 	repetition 						query 		string 		false 		"Repetition" 		Enums(one_time,periodic)
// @Success     	200  			{object}  		lib.APIResponsePaginated{data=dto.CompressionResponses}
// @Failure     	400  			{object}  		lib.HTTPError
// @Failure     	500  			{object}  		lib.HTTPError
// @Router      	/v1/compression/ 	[get]
func (h *CompressionHandlerImpl) FindAll(ctx *gin.Context) {
	// log.Info().Msg("CompressionHandler.FindAll")
	var err error

	filter := h.buildFindFilter(ctx)
	err = h.validateFilter(filter)
	if err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	filterPagination, err := FindAllFilter(ctx)
	if err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	filter.Filtering = filterPagination
	resPaginated, err := h.CompressionService.FindAll(filter)
	if err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	data := resPaginated.Compressions.CreateResponse()
	pagination := resPaginated.GetPaginationResponse()
	lib.RespondSuccessPaginated(ctx, http.StatusOK, lib.MsgOk, data, pagination)
}

// FindById godoc
// @Summary     	Find one compression
// @Description 	by id
// @Security 		BearerAuth
// @Tags        	Compression
// @Accept      	json
// @Produce     	json
// @Param 			id 			path 			string 		true 		"Compression Id"
// @Success     	200  		{object}  		lib.APIResponse{data=dto.CompressionResponse}
// @Failure     	400  		{object}  		lib.HTTPError
// @Failure     	500  		{object}  		lib.HTTPError
// @Router      	/v1/compression/{id} [get]
func (h *CompressionHandlerImpl) FindById(ctx *gin.Context) {
	// log.Info().Msg("CompressionHandler.FindById")
	var err error

	commpressionId := ctx.Param("id")

	compression, err := h.CompressionService.FindById(commpressionId)
	if err != nil {
		lib.RespondError(ctx, http.StatusNotFound, err.Error(), err)
		return
	}

	if compression.Id == "" {
		lib.RespondError(ctx, http.StatusNotFound, lib.ErrRecordNotFound, err)
		return
	}

	data := compression.CreateResponse()
	lib.RespondSuccess(ctx, http.StatusOK, lib.MsgOk, data)
}

// DeleteById godoc
// @Summary      	Delete one compression
// @Description  	by id
// @Security 	 	BearerAuth
// @Tags         	Compression
// @Accept       	json
// @Produce      	json
// @Param 		 	id 			path 			string 		true 	"Compression Id"
// @Success      	200  		{object}  		lib.APIResponse
// @Failure      	400  		{object}  		lib.HTTPError
// @Failure      	500  		{object}  		lib.HTTPError
// @Router       	/v1/compression/{id} 	[delete]
func (h *CompressionHandlerImpl) DeleteById(ctx *gin.Context) {
	// log.Info().Msg("CompressionHandler.DeleteById")
	var err error

	commpressionId := ctx.Param("id")

	compression, err := h.CompressionService.FindById(commpressionId)
	if err != nil {
		lib.RespondError(ctx, http.StatusNotFound, err.Error(), err)
		return
	}

	if compression.Id == "" {
		lib.RespondError(ctx, http.StatusNotFound, lib.ErrRecordNotFound, err)
		return
	}

	err = h.CompressionService.DeleteById(compression.Id)
	if err != nil {
		lib.RespondError(ctx, http.StatusInternalServerError, lib.ErrDeleteFailed, err)
		return
	}

	countDeletedCompression := 1
	countDeletedCron := 0

	if compression.Repetition == lib.Periodic {
		err = h.ScheduleService.DeleteCron(compression.Id)
		if err != nil {
			log.Error().Err(err).Msg("DeleteCron failed")
		} else {
			countDeletedCron += 1
		}
	}

	data := dto.DeleteCompressionResponse{
		Success:            true,
		DeletedCompression: countDeletedCompression,
		DeletedCron:        countDeletedCron,
	}
	lib.RespondSuccess(ctx, http.StatusOK, lib.MsgOk, data)
}

/** Multiple **/

// FindByIds godoc
// @Summary     	Find compressions by multiple id
// @Description 	Find records
// @Security 		BearerAuth
// @Tags        	Compression
// @Accept      	json
// @Produce     	json
// @Param 			ids 		path 			string 		true 		"Compression Ids"
// @Success     	200  		{object}  		lib.APIResponse{data=dto.CompressionResponses}
// @Failure     	400  		{object}  		lib.HTTPError
// @Failure     	500  		{object}  		lib.HTTPError
// @Router      	/v1/compression/many/{ids} [get]
func (h *CompressionHandlerImpl) FindByIds(ctx *gin.Context) {
	log.Info().Msg("CompressionHandler.FindByIds")
	var err error

	commpressionIds := ctx.Param("ids")
	log.Info().Str("commpressionIds", commpressionIds).Msg("payload")

	compressions, err := h.CompressionService.FindByIds(commpressionIds)
	if err != nil {
		lib.RespondError(ctx, http.StatusNotFound, err.Error(), err)
		return
	}

	data := compressions.CreateResponse()
	lib.RespondSuccess(ctx, http.StatusOK, lib.MsgOk, data)
}

// DeleteByIds godoc
// @Summary      	Delete compressions
// @Description  	by ids
// @Security 	 	BearerAuth
// @Tags         	Compression
// @Accept       	json
// @Produce      	json
// @Param 		 	ids 		path 			string 				true 		"Compression Ids"
// @Success      	200  		{object}  		lib.APIResponse
// @Failure      	400  		{object}  		lib.HTTPError
// @Failure      	500  		{object}  		lib.HTTPError
// @Router       	/v1/compression/many/{ids} 	[delete]
func (h *CompressionHandlerImpl) DeleteByIds(ctx *gin.Context) {
	var err error

	commpressionIds := ctx.Param("ids")

	compressions, err := h.CompressionService.FindByIds(commpressionIds)
	if err != nil {
		lib.RespondError(ctx, http.StatusNotFound, err.Error(), err)
		return
	}

	if len(compressions) < 1 || compressions == nil {
		err = errors.New("no records found")
		lib.RespondError(ctx, http.StatusNotFound, err.Error(), err)
		return
	}

	var foundCompressionIds []string
	for _, c := range compressions {
		foundCompressionIds = append(foundCompressionIds, c.Id)
	}

	_, err = h.CompressionService.DeleteByIds(foundCompressionIds)
	if err != nil {
		lib.RespondError(ctx, http.StatusInternalServerError, lib.ErrDeleteFailed, err)
		return
	}

	countDeletedCron := 0
	for _, compression := range compressions {
		if compression.Repetition == lib.Periodic {
			err = h.ScheduleService.DeleteCron(compression.Id)
			if err != nil {
				log.Error().Err(err).Msg("DeleteCron failed")
			} else {
				countDeletedCron += 1
			}
		}
	}

	data := dto.DeleteCompressionResponse{
		Success:            true,
		DeletedCompression: len(compressions),
		DeletedCron:        countDeletedCron,
	}
	lib.RespondSuccess(ctx, http.StatusOK, lib.MsgOk, data)
}

func (h *CompressionHandlerImpl) buildFindFilter(ctx *gin.Context) dto.FilterFindCompression {
	filter := dto.FilterFindCompression{}

	repetition := ctx.Query(dto.FilterRepetition)

	if repetition != "" {
		filter.Repetition = repetition
	}

	return filter
}

func (h *CompressionHandlerImpl) validateFilter(filter dto.FilterFindCompression) error {
	var errMessages []string

	if filter.Repetition != "" && !helper.IsValidRepetition(filter.Repetition) {
		errMessages = append(errMessages, fmt.Sprintf("repetition %s is not valid", filter.Repetition))
	}

	if len(errMessages) > 0 {
		return errors.New(lib.ListToString(errMessages))
	}

	return nil
}
