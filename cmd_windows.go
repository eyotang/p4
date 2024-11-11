//go:build windows
// +build windows

package p4

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
)

func watch(ctx context.Context, cmd *exec.Cmd, waitCh chan struct{}) error {
	if cmd == nil {
		return errors.New("cmd is nil")
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}

	go func() {
		select {
		case <-ctx.Done(): // 超时处理
			err := cmd.Process.Kill()
			if err != nil {
				fmt.Printf("kill error   : [%v]\n", err)
			}
			fmt.Printf("cmd '%s' killed\n", cmd.String())
			return
		case <-waitCh: // 正常结束
			//fmt.Println("wait ok, normal exit")
			return
		}
	}()

	return nil
}
