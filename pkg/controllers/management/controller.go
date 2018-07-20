package management

import (
	"context"

	"github.com/rancher/rancher/pkg/controllers/management/auth"
	"github.com/rancher/types/config"
)

func Register(ctx context.Context, management *config.ManagementContext) {
	// auth handlers need to run early to create namespaces that back clusters and projects
	// also, these handlers are purely in the mgmt plane, so they are lightweight compared to those that interact with machines and clusters
	auth.RegisterEarly(ctx, management)

	// Register last
	auth.RegisterLate(ctx, management)
}
