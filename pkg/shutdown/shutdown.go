package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

type Shutdown struct {
	ctx chan os.Signal
}

// New 构造
func New() *Shutdown {
	sd := &Shutdown{
		ctx: make(chan os.Signal, 1),
	}
	return sd.WithSignals(syscall.SIGINT, syscall.SIGTERM)
}

func (sd *Shutdown) WithSignals(signals ...syscall.Signal) *Shutdown {
	for _, s := range signals {
		//监听指定信号
		signal.Notify(sd.ctx, s)
	}

	return sd
}

func (sd *Shutdown) Close(funcs ...func()) {
	//监听
	<-sd.ctx
	//取消监听
	signal.Stop(sd.ctx)

	//回收函数
	for _, f := range funcs {
		f()
	}
}
