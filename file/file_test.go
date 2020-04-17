package file

import (
	"fmt"
	"testing"
)

func TestIsExist(t *testing.T) {
	fmt.Println(IsExist("~/Workspaces/../a.php"))
	fmt.Println(IsExist("/Users/hulingjie/Workspaces/../a.php"))
}
