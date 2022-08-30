package gzk

import (
	_ "fmt"
	_ "time"

	_ "github.com/go-zookeeper/zk"
)

func init() {

	// _, _, err := zk.Connect([]string{"127.0.0.1:2181", "127.0.0.1:2182", "127.0.0.1:2183"}, time.Second) //*10)

	// if err != nil {
	// 	panic(err)
	// }

	// children, stat, ch, err := c.ChildrenW("/")

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("%+v %+v\n", children, stat)

	// e := <-ch

	// fmt.Printf("%+v\n", e)

}
