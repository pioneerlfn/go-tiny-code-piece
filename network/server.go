pacakge main

import "net"


func main()  {
	ln, err := net.Listen("tcp", 9000)
	panicOnErr(err)
	
	



}