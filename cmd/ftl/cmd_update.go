package main

import (
	"context"

	"connectrpc.com/connect"

	"github.com/TBD54566975/ftl/backend/common/model"
	ftlv1 "github.com/TBD54566975/ftl/protos/xyz/block/ftl/v1"
	"github.com/TBD54566975/ftl/protos/xyz/block/ftl/v1/ftlv1connect"
)

type updateCmd struct {
	Replicas   int32                `short:"n" help:"Number of replicas to deploy." default:"1"`
	Deployment model.DeploymentName `arg:"" help:"Deployment to update."`
}

func (u *updateCmd) Run(ctx context.Context, client ftlv1connect.ControllerServiceClient) error {
	_, err := client.UpdateDeploy(ctx, connect.NewRequest(&ftlv1.UpdateDeployRequest{
		DeploymentName: u.Deployment.String(),
		MinReplicas:    u.Replicas,
	}))
	if err != nil {
		return err
	}
	return nil
}
