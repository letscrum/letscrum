package dao

import (
	"context"

	"github.com/daocloud/skoala/app/hive/internal/model/v1alpha1"
	metav1alpha1 "github.com/daocloud/skoala/pkg/meta/v1alpha1"
)

// RegistryDao represents the dao interface for registry.
type RegistryDao interface {
	// list registries, if page < 0 query all
	List(ctx context.Context, page, pageSize int32, opts metav1alpha1.ListOptions) (*v1alpha1.RegistryList, error)
	Create(ctx context.Context, reg *v1alpha1.Registry, opts metav1alpha1.CreateOptions) error

	// Get get the registry by uid
	Get(ctx context.Context, uid string, opts metav1alpha1.GetOptions) (*v1alpha1.Registry, error)
	Update(ctx context.Context, reg *v1alpha1.Registry, opts metav1alpha1.UpdateOptions) error
	Delete(ctx context.Context, reg *v1alpha1.Registry, opts metav1alpha1.DeleteOptions) error
}
