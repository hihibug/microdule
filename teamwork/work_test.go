package teamwork_test

import (
	"testing"

	"github.com/hihibug/microdule/teamwork"
)

func TestTeamWork(t *testing.T) {
	tw := teamwork.NewTeamwork()

	// tw.Reginster("test", func() {
	// 	log.Println("test")
	// 	time.Sleep(10 * time.Second)
	// 	panic("err")
	// }).HandleClose(func() {
	// 	log.Println("test close")
	// })

	if err := tw.Start(); err != nil {
		tw.Close()
	}
}
