package pachyderm

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/pachyderm/pachyderm/src/client/auth"
)

func (b *backend) revokePath() *framework.Path {

	return &framework.Path{
		Pattern: "revoke",
		Fields: map[string]*framework.FieldSchema{
			"user_token": &framework.FieldSchema{
				Type: framework.TypeString,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRevoke,
		},
	}
}

func (b *backend) pathRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	userToken := d.Get("user_token").(string)
	if len(userToken) == 0 {
		return nil, logical.ErrInvalidRequest
	}

	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if len(config.AdminToken) == 0 {
		return nil, errors.New("plugin is missing admin token")
	}

	err = b.revokeUserCredentials(ctx, userToken, config.AdminToken)
	if err != nil {
		return nil, err
	}

	// Compose the response
	// TODO: Not sure if this is the right way to return a successful response
	return &logical.Response{
		Auth: &logical.Auth{},
	}, nil
}

func (b *backend) revokeUserCredentials(ctx context.Context, userToken string, adminToken string) error {
	// This is where we'd make the actual pachyderm calls to create the user
	// token using the admin token. For now, for testing purposes, we just do an action that only an
	// admin could do

	// Setup a single use client w the given auth token
	pachClient := b.PachydermClient.WithCtx(ctx)
	pachClient.SetAuthToken(adminToken)

	_, err := b.PachydermClient.AuthAPIClient.ModifyAdmins(pachClient.Ctx(), &auth.ModifyAdminsRequest{
		Remove: []string{"tweetybird"},
	})

	if err != nil {
		return err
	}

	return nil
}
