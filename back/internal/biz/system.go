// 要求底层提供“检查数据库连接”的能力

package biz

import "context"

type HealthRepository interface{ Ping(context.Context) error }
type SystemUsecase struct{ repo HealthRepository }

func NewSystemUsecase(repo HealthRepository) *SystemUsecase { return &SystemUsecase{repo: repo} }
func (u *SystemUsecase) Ready(ctx context.Context) error    { return u.repo.Ping(ctx) }
