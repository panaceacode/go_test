package repository

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"go_test/webook/internal/domain"
	"go_test/webook/internal/repository/cache"
	"go_test/webook/internal/repository/dao"
	"time"
)

type CachedUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

var (
	ErrorDuplicateEmail = dao.ErrorDuplicateEmail
	ErrorUserNotFound   = dao.ErrorUserNotFound
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
	}
}

func (repo *UserRepository) UpdateNonZeroFields(ctx *gin.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(user))
}

func (repo *UserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}

func (repo *UserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	//if ctx.Value("x-stress") != true {
	//	du, err := repo.cache.Get(ctx, uid)
	//	// 只要 err 为 nil，就返回
	//	if err == nil {
	//		return du, nil
	//	}
	//}
	//du, err := repo.cache.Get(ctx, uid)
	//// 只要 err 为 nil，就返回
	//if err == nil {
	//	return du, nil
	//}

	// 检测限流/熔断/降级标记位
	//if ctx.Value("downgrade") == "true" {
	//	return du, errors.New("触发降级，不再查询数据库")
	//}

	// err 不为 nil，就要查询数据库
	// err 有两种可能
	// 1. key 不存在，说明 redis 是正常的
	// 2. 访问 redis 有问题。可能是网络有问题，也可能是 redis 本身就崩溃了

	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	du := repo.toDomain(u)
	//go func() {
	//	err = repo.cache.Set(ctx, du)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}()

	//err = repo.cache.Set(ctx, du)
	//if err != nil {
	//	// 网络崩了，也可能是 redis 崩了
	//	log.Println(err)
	//}
	return du, nil
}

func (repo *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:        u.Id,
		Email:     u.Email.String,
		Phone:     u.Phone.String,
		Password:  u.Password,
		AboutMe:   u.AboutMe,
		Nickname:  u.Nickname,
		Birthday:  time.UnixMilli(u.Birthday),
		CreatedAt: time.UnixMilli(u.CreatedAt),
	}
}
