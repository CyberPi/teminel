package exec

import (
	"fmt"
	"os/exec"

	"source.cyberpi.de/go/teminel/utils"
)

func Shell(args ...string) error {
	stdout, err := exec.Command("sh", append([]string{"-c"}, args...)...).Output()
	if err != nil {
		return err
	}
	fmt.Println(utils.StringsToAny(append([]string{"Run:"}, args...)...)...)
	for _, line := range stdout {
		fmt.Println(line)
	}
	return nil
}
