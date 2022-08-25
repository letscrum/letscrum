package dao

import (
	"context"

	"github.com/daocloud/skoala/app/hive/internal/model/v1alpha1"
	metav1alpha1 "github.com/daocloud/skoala/pkg/meta/v1alpha1"
)

// BookDao represents the interface of book dao.
type BookDao interface {
	Get(ctx context.Context, uid string, opts metav1alpha1.GetOptions) (*v1alpha1.Book, error)
	List(ctx context.Context, pageNum, pageSize int32, opts metav1alpha1.ListOptions) (*v1alpha1.BookList, error)
	Create(ctx context.Context, book *v1alpha1.Book, opts metav1alpha1.CreateOptions) error
	Update(ctx context.Context, book *v1alpha1.Book, opts metav1alpha1.UpdateOptions) error
	Delete(ctx context.Context, uid string, opts metav1alpha1.DeleteOptions) error
}
