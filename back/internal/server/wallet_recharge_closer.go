package server

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"dd_web/back/internal/biz"
)

// WalletRechargeCloser 是订单服务内的专用超时关单循环。多实例可同时运行，
// 仓储通过状态与 version CAS 保证同一充值单只有一个实例获得处理权。
type WalletRechargeCloser struct {
	wallet     *biz.WalletUsecase
	logger     *slog.Logger
	interval   time.Duration
	retryDelay time.Duration
	batchSize  int

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
}

func NewWalletRechargeCloser(wallet *biz.WalletUsecase, logger *slog.Logger, interval, retryDelay time.Duration, batchSize int) *WalletRechargeCloser {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	if retryDelay <= 0 {
		retryDelay = time.Minute
	}
	if batchSize <= 0 {
		batchSize = 100
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &WalletRechargeCloser{wallet: wallet, logger: logger, interval: interval, retryDelay: retryDelay, batchSize: batchSize}
}

func (w *WalletRechargeCloser) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	w.mu.Lock()
	w.cancel, w.done = cancel, done
	w.mu.Unlock()
	defer close(done)

	w.run(ctx)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			w.run(ctx)
		}
	}
}

func (w *WalletRechargeCloser) Stop(ctx context.Context) error {
	w.mu.Lock()
	cancel, done := w.cancel, w.done
	w.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	if done == nil {
		return nil
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *WalletRechargeCloser) run(ctx context.Context) {
	if w.wallet == nil {
		return
	}
	result, err := w.wallet.CloseExpiredRecharges(ctx, w.batchSize, w.retryDelay)
	if err != nil {
		if ctx.Err() == nil {
			w.logger.Error("close expired wallet recharges failed", "error", err)
		}
		return
	}
	if result.Scanned > 0 {
		w.logger.Info("processed expired wallet recharges", "scanned", result.Scanned, "closed", result.Closed, "credited", result.Credited, "retrying", result.Retrying, "conflicts", result.Conflicts)
	}
}
