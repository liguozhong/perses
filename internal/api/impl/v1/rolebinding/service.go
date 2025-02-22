// Copyright 2021 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rolebinding

import (
	"fmt"

	apiInterface "github.com/perses/perses/internal/api/interface"
	"github.com/perses/perses/internal/api/interface/v1/role"
	"github.com/perses/perses/internal/api/interface/v1/user"
	"github.com/perses/perses/internal/api/shared/rbac"

	"github.com/perses/perses/internal/api/interface/v1/rolebinding"
	"github.com/perses/perses/internal/api/shared"
	databaseModel "github.com/perses/perses/internal/api/shared/database/model"
	"github.com/perses/perses/internal/api/shared/schemas"
	"github.com/perses/perses/pkg/model/api"
	v1 "github.com/perses/perses/pkg/model/api/v1"
	"github.com/sirupsen/logrus"
)

type service struct {
	rolebinding.Service
	dao     rolebinding.DAO
	roleDAO role.DAO
	userDAO user.DAO
	rbac    rbac.RBAC
	sch     schemas.Schemas
}

func NewService(dao rolebinding.DAO, roleDAO role.DAO, userDAO user.DAO, rbac rbac.RBAC, sch schemas.Schemas) rolebinding.Service {
	return &service{
		dao:     dao,
		rbac:    rbac,
		roleDAO: roleDAO,
		userDAO: userDAO,
		sch:     sch,
	}
}

func (s *service) Create(_ apiInterface.PersesContext, entity api.Entity) (interface{}, error) {
	if object, ok := entity.(*v1.RoleBinding); ok {
		return s.create(object)
	}
	return nil, shared.HandleBadRequestError(fmt.Sprintf("wrong entity format, attempting roleBinding format, received '%T'", entity))
}

func (s *service) create(entity *v1.RoleBinding) (*v1.RoleBinding, error) {
	// Update the time contains in the entity
	entity.Metadata.CreateNow()
	if err := s.validateRoleBinding(entity); err != nil {
		return nil, err
	}
	if err := s.dao.Create(entity); err != nil {
		return nil, err
	}
	// Refreshing RBAC cache as the role binding can add or remove new permissions to concerned users
	if err := s.rbac.Refresh(); err != nil {
		logrus.WithError(err).Error("failed to refresh RBAC cache")
	}
	return entity, nil
}

func (s *service) Update(_ apiInterface.PersesContext, entity api.Entity, parameters apiInterface.Parameters) (interface{}, error) {
	if object, ok := entity.(*v1.RoleBinding); ok {
		return s.update(object, parameters)
	}
	return nil, shared.HandleBadRequestError(fmt.Sprintf("wrong entity format, attempting roleBinding format, received '%T'", entity))
}

func (s *service) update(entity *v1.RoleBinding, parameters apiInterface.Parameters) (*v1.RoleBinding, error) {
	if entity.Metadata.Name != parameters.Name {
		logrus.Debugf("name in roleBinding %q and name from the http request %q don't match", entity.Metadata.Name, parameters.Name)
		return nil, shared.HandleBadRequestError("metadata.name and the name in the http path request don't match")
	}
	if len(entity.Metadata.Project) == 0 {
		entity.Metadata.Project = parameters.Project
	} else if entity.Metadata.Project != parameters.Project {
		logrus.Debugf("project in roleBinding %q and project from the http request %q don't match", entity.Metadata.Project, parameters.Project)
		return nil, shared.HandleBadRequestError("metadata.project and the project name in the http path request don't match")
	}

	// find the previous version of the RoleBinding
	oldEntity, err := s.dao.Get(parameters.Project, parameters.Name)
	if err != nil {
		return nil, err
	}

	if err := s.validateRoleBinding(entity); err != nil {
		return nil, err
	}

	// If you do want to change the role for a binding, you need to remove the binding object and create a replacement.
	// More info at: https://github.com/perses/perses/blob/main/docs/authorization.md#rolebinding-and-globalrolebinding-update-restriction
	if entity.Spec.Role != oldEntity.Spec.Role {
		return nil, shared.HandleBadRequestError("spec.role can't be updated")
	}

	entity.Metadata.Update(oldEntity.Metadata)
	if updateErr := s.dao.Update(entity); updateErr != nil {
		logrus.WithError(updateErr).Errorf("unable to perform the update of the roleBinding %q, something wrong with the database", entity.Metadata.Name)
		return nil, updateErr
	}
	// Refreshing RBAC cache as the role binding can add or remove new permissions to concerned users
	if err := s.rbac.Refresh(); err != nil {
		logrus.WithError(err).Error("failed to refresh RBAC cache")
	}
	return entity, nil
}

func (s *service) Delete(_ apiInterface.PersesContext, parameters apiInterface.Parameters) error {
	if err := s.dao.Delete(parameters.Project, parameters.Name); err != nil {
		return err
	}
	// Refreshing RBAC cache as the role binding can add or remove new permissions to concerned users
	if err := s.rbac.Refresh(); err != nil {
		logrus.WithError(err).Error("failed to refresh RBAC cache")
	}
	return nil
}

func (s *service) Get(_ apiInterface.PersesContext, parameters apiInterface.Parameters) (interface{}, error) {
	return s.dao.Get(parameters.Project, parameters.Name)
}

func (s *service) List(_ apiInterface.PersesContext, q databaseModel.Query, _ apiInterface.Parameters) (interface{}, error) {
	return s.dao.List(q)
}

// Validating role and subjects are existing
func (s *service) validateRoleBinding(roleBinding *v1.RoleBinding) error {
	if _, err := s.roleDAO.Get(roleBinding.Metadata.Project, roleBinding.Spec.Role); err != nil {
		return shared.HandleBadRequestError(fmt.Sprintf("role %q doesn't exist", roleBinding.Spec.Role))
	}

	for _, subject := range roleBinding.Spec.Subjects {
		if subject.Kind == v1.KindUser {
			if _, err := s.userDAO.Get(subject.Name); err != nil {
				if databaseModel.IsKeyNotFound(err) {
					return shared.HandleBadRequestError(fmt.Sprintf("user subject name %q doesn't exist", subject.Name))
				}
				logrus.WithError(err).Errorf("unable to find the user with the name %q", subject.Name)
				return err
			}
		}
	}
	return nil
}
