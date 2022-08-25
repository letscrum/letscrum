package mysql

import (
	"context"
	"errors"

	v1alpha1 "github.com/daocloud/skoala/app/hive/internal/model/v1alpha1"
	"github.com/daocloud/skoala/pkg/gormutil"
	metav1alpha1 "github.com/daocloud/skoala/pkg/meta/v1alpha1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type bookDao struct {
	db *gorm.DB
}

func newBookDao(d *hiveDao) *bookDao {
	return &bookDao{d.db}
}

func (u *bookDao) Get(_ context.Context, uid string, _ metav1alpha1.GetOptions) (*v1alpha1.Book, error) {
	book := &v1alpha1.Book{}

	err := u.db.Where("uid = ?", uid).First(&book).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error()) //nolint:wrapcheck // as is
		}

		return nil, status.Error(codes.Unknown, err.Error()) //nolint:wrapcheck // as is
	}

	return book, nil
}

func (u *bookDao) List(_ context.Context, pageNum, pageSize int32, _ metav1alpha1.ListOptions) (*v1alpha1.BookList, error) {
	results := &v1alpha1.BookList{}

	d := u.db.Offset(int(gormutil.GetPageOffset(pageNum, pageSize))).
		Limit(int(pageSize)).
		Order("id desc").
		Find(&results.Items).
		Offset(-1).
		Limit(-1).
		Count(&results.TotalCount)

	return results, d.Error
}

func (u *bookDao) Create(_ context.Context, book *v1alpha1.Book, _ metav1alpha1.CreateOptions) error {
	return u.db.Create(book).Error
}

func (u *bookDao) Update(_ context.Context, book *v1alpha1.Book, _ metav1alpha1.UpdateOptions) error {
	return u.db.Save(book).Error
}

func (u *bookDao) Delete(_ context.Context, uid string, _ metav1alpha1.DeleteOptions) error {
	return u.db.Where("uid = ?", uid).Delete(&v1alpha1.Book{}).Error
}
