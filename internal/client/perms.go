package client

import (
	"context"
	"git.miem.hse.ru/1206/app/errs"
	perms "git.miem.hse.ru/1206/app/permissions-v2"
	"strconv"
)

type PermissionsClient struct {
	*perms.PermissionsClient
}

func NewPermissionsClient(cfg *perms.PermissionsDbConfig) (*PermissionsClient, error) {
	pc, err := perms.NewPermissionsClient(context.Background(), cfg)
	if err != nil {
		return nil, errs.NewInternal(err)
	}
	return &PermissionsClient{PermissionsClient: pc}, nil
}

func (c *PermissionsClient) CheckCanUploadMaterials(ctx context.Context, userId int64) (bool, error) {
	return c.CheckPermission(
		ctx,
		perms.DefinitionPlatform,
		perms.CPIS,
		perms.PermissionPlatform_UploadMaterials,
		perms.DefinitionUser,
		strconv.FormatInt(userId, 10),
	)
}

func (c *PermissionsClient) CheckCanEditMaterial(ctx context.Context, userId int64, materialId int64) (bool, error) {
	return c.CheckPermission(
		ctx,
		perms.DefinitionMaterial,
		strconv.FormatInt(materialId, 10),
		perms.PermissionMaterial_Edit,
		perms.DefinitionUser,
		strconv.FormatInt(userId, 10),
	)
}

func (c *PermissionsClient) CheckCanViewMaterial(ctx context.Context, userId int64, materialId int64) (bool, error) {
	return c.CheckPermission(
		ctx,
		perms.DefinitionMaterial,
		strconv.FormatInt(materialId, 10),
		perms.PermissionMaterial_View,
		perms.DefinitionUser,
		strconv.FormatInt(userId, 10),
	)
}

func (c *PermissionsClient) SetNewMaterialRelations(ctx context.Context, ownerId int64, materialId int64) error {
	err := c.SetRelation(
		ctx,
		perms.DefinitionMaterial,
		strconv.FormatInt(materialId, 10),
		perms.DefinitionUser,
		strconv.FormatInt(ownerId, 10),
		perms.RelationMaterial_Owner,
	)
	if err != nil {
		return err
	}
	err = c.SetRelation(
		ctx,
		perms.DefinitionMaterial,
		strconv.FormatInt(materialId, 10),
		perms.DefinitionPlatform,
		perms.CPIS,
		perms.RelationMaterial_Platform,
	)
	return err
}

func (c *PermissionsClient) DeleteMaterialRelations(ctx context.Context, materialId int64) error {
	return c.DeleteAllResourceRelations(
		ctx,
		perms.DefinitionMaterial,
		strconv.FormatInt(materialId, 10),
	)
}

func (c *PermissionsClient) ListAvailableMaterialIds(ctx context.Context, userId int64) ([]int64, error) {
	return c.LookupRelations(
		ctx,
		perms.DefinitionMaterial,
		perms.PermissionMaterial_View,
		perms.DefinitionUser,
		strconv.FormatInt(userId, 10),
	)
}
