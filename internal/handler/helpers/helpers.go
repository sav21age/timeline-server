package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"timeline/internal/domain"

	"github.com/gin-gonic/gin"
	// "github.com/pkg/errors"
)

func RequiredIntParam(ctx *gin.Context, name string) (id int, err error) {
	id, err = strconv.Atoi(ctx.Param(name))
	if err != nil {
		msg := fmt.Sprintf("%s: %s", domain.MsgQueryParamInvalid, name)
		return 0, errors.New(msg)
	}

	return
}

func InSlice[T comparable](a T, list []T) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func SetPagination(page, per_page string, limitDefaultSlice []int, limitDefault int) (domain.Pagination, error) {
	pagination := domain.Pagination{}

	skip, errSkip := strconv.Atoi(page)
	limit, errLimit := strconv.Atoi(per_page)

	if errLimit == nil {
		ok := InSlice(int(limit), limitDefaultSlice)
		if !ok {
			msg := fmt.Sprintf("%s: per_page", domain.MsgQueryParamInvalid)
			return pagination, errors.New(msg)
		}
	} else {
		limit = limitDefault
	}

	if errSkip == nil {
		pagination.Skip = (skip - 1) * limit
		pagination.Limit = limit
	} else {
		msg := fmt.Sprintf("%s: page", domain.MsgQueryParamInvalid)
		return pagination, errors.New(msg)
	}

	return pagination, nil
}

func SetDatination(min_date, max_date string) domain.Datination {
	datination := domain.Datination{}

	// minDate, errMinDate := time.Parse(time.RFC3339, min_date)
	// maxDate, errMaxDate := time.Parse(time.RFC3339, max_date)

	min, errMinDate := strconv.ParseInt(min_date, 10, 64)
	minDate := time.UnixMilli(min)

	max, errMaxDate := strconv.ParseInt(max_date, 10, 64)
	maxDate := time.UnixMilli(max)

	if errMinDate == nil && errMaxDate == nil {
		datination.MinDate = minDate
		datination.MaxDate = maxDate
	}

	return datination
}
