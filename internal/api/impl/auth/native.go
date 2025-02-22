// Copyright 2023 The Perses Authors
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

package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/perses/perses/internal/api/interface/v1/user"
	"github.com/perses/perses/internal/api/shared"
	"github.com/perses/perses/internal/api/shared/crypto"
	databaseModel "github.com/perses/perses/internal/api/shared/database/model"
	"github.com/perses/perses/internal/api/shared/route"
	"github.com/perses/perses/internal/api/shared/utils"
	"github.com/perses/perses/pkg/model/api"
)

type nativeEndpoint struct {
	dao             user.DAO
	jwt             crypto.JWT
	tokenManagement tokenManagement
}

func newNativeEndpoint(dao user.DAO, jwt crypto.JWT) route.Endpoint {
	return &nativeEndpoint{
		dao:             dao,
		jwt:             jwt,
		tokenManagement: tokenManagement{jwt: jwt},
	}
}

func (e *nativeEndpoint) CollectRoutes(g *route.Group) {
	g.POST(fmt.Sprintf("/%s/%s", utils.AuthKindNative, utils.PathLogin), e.auth, true)
}

func (e *nativeEndpoint) auth(ctx echo.Context) error {
	body := &api.Auth{}
	if err := ctx.Bind(body); err != nil {
		return shared.HandleBadRequestError(err.Error())
	}
	usr, err := e.dao.Get(body.Login)
	if err != nil {
		if databaseModel.IsKeyNotFound(err) {
			return shared.HandleBadRequestError("wrong login or password ")
		}
		return shared.InternalError
	}

	if !crypto.ComparePasswords(usr.Spec.NativeProvider.Password, body.Password) {
		return shared.HandleBadRequestError("wrong login or password ")
	}
	login := body.Login
	accessToken, err := e.tokenManagement.accessToken(login, ctx.SetCookie)
	if err != nil {
		return err
	}
	refreshToken, err := e.tokenManagement.refreshToken(login, ctx.SetCookie)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, api.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
