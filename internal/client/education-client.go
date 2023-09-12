package client

import (
	"context"
	"git.miem.hse.ru/1206/app"
	"git.miem.hse.ru/1206/material-library/internal/models/tplibrary"
	"git.miem.hse.ru/1206/proto/pb"
	"github.com/friendsofgo/errors"
)

type EducationClient struct {
	grpcClient    *app.GRPCClient
	edu           pb.CalendarClient
	difLevelCache map[int64]string
	subjectsCache map[int64]string
}

func NewEducationClient(cfg *app.GRPCConfig) (*EducationClient, error) {
	cl, err := app.NewGRPCClient(cfg)
	if err != nil {
		return nil, err
	}
	pb.NewCalendarClient(cl.Conn)
	return &EducationClient{
		grpcClient: cl,
		difLevelCache: map[int64]string{
			1: "A1",
			2: "A2",
		},
		subjectsCache: map[int64]string{
			1: "Математика",
			2: "Биология",
		},
	}, nil
}

func (c *EducationClient) ValidateSubjectId(ctx context.Context, subjectId int64) (bool, error) {
	if _, ok := c.subjectsCache[subjectId]; ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (c *EducationClient) ValidateDifficultcyLevelId(ctx context.Context, difficultcyLevelId int64) (bool, error) {
	if _, ok := c.difLevelCache[difficultcyLevelId]; ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (c *EducationClient) GetSubjectById(ctx context.Context, subjectId int64) (tplibrary.IdName, error) {
	if sName, ok := c.subjectsCache[subjectId]; ok {
		return tplibrary.IdName{
			Id:   subjectId,
			Name: sName,
		}, nil
	}
	return tplibrary.IdName{}, errors.New("not found")
}

func (c *EducationClient) GetDifficultcyLevelById(ctx context.Context, difficultcyLevelId int64) (tplibrary.IdName, error) {
	if sName, ok := c.difLevelCache[difficultcyLevelId]; ok {
		return tplibrary.IdName{
			Id:   difficultcyLevelId,
			Name: sName,
		}, nil
	}
	return tplibrary.IdName{}, errors.New("not found")
}

func (c *EducationClient) GetSubjectIdToName(ctx context.Context) (map[int64]string, error) {
	return c.subjectsCache, nil
}

func (c *EducationClient) GetDifficultcyLevelIdToName(ctx context.Context) (map[int64]string, error) {
	return c.difLevelCache, nil
}
