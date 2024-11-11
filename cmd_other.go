//go:build !windows
// +build !windows

package p4

import (
	"context"
	"log"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
)

func watch(ctx context.Context, cmd *exec.Cmd, waitCh chan struct{}) error {
	if cmd == nil {
		return errors.New("cmd is nil")
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	go func() {
		select {
		case <-ctx.Done(): // 超时处理
			err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			if err != nil {
				log.Printf("kill error   : [%v]\n", err)
				return
			}
			log.Println("killed")
			return
		case <-waitCh: // 正常结束
			//log.Println("wait ok, normal exit")
			return
		}
	}()

	return nil
}
