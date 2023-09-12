package domain

import "git.miem.hse.ru/1206/app/storage/stpg"

type Usecase struct {
	st    stpg.PGer
	perms Permissions
	s3    S3
	edu   Education
}

func NewUsecase(p Permissions, s3 S3, edu Education) *Usecase {
	return &Usecase{
		st:    stpg.Gist(),
		perms: p,
		s3:    s3,
		edu:   edu,
	}
}
