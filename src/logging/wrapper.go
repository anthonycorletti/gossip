package logging

// TODO: support text coloring (to terminal only, most likely)

import (
	"log"
)

type Router struct {
	Directory string
	out       *log.Logger
	warn      *log.Logger
	err       *log.Logger
	web       *log.Logger
}

func (self *Router) Web() *log.Logger {
	return self.web
}

func (self *Router) Fatal(v ...interface{}) {
	self.err.Fatal(v...)
	self.out.Println("ERR: Fatal error sent to error.log")
}

func (self *Router) Fatalf(f string, v ...interface{}) {
	self.err.Fatalf(f, v...)
	self.out.Println("ERR: Fatal error sent to error.log")
}

func (self *Router) Fatalln(v ...interface{}) {
	self.err.Fatalln(v...)
	self.out.Println("ERR: Fatal error sent to error.log")
}

func (self *Router) Panic(v ...interface{}) {
	self.err.Panic(v...)
	self.out.Println("ERR: Panic error sent to error.log")
}

func (self *Router) Panicf(f string, v ...interface{}) {
	self.err.Panicf(f, v...)
	self.out.Println("ERR: Panic error sent to error.log")
}

func (self *Router) Panicln(v ...interface{}) {
	self.err.Panicln(v...)
	self.out.Println("ERR: Panic error sent to error.log")
}

func (self *Router) Print(v ...interface{}) {
	self.out.Print(v...)
}

func (self *Router) Printf(f string, v ...interface{}) {
	self.out.Printf(f, v...)
}

func (self *Router) Println(v ...interface{}) {
	self.out.Println(v...)
}

func (self *Router) Warn(v ...interface{}) {
	self.warn.Print(v...)
}

func (self *Router) Warnf(f string, v ...interface{}) {
	self.warn.Printf(f, v...)
}

func (self *Router) Warnln(v ...interface{}) {
	self.warn.Println(v...)
}

func (self *Router) Error(v ...interface{}) {
	self.err.Print(v...)
}

func (self *Router) Errorf(f string, v ...interface{}) {
	self.err.Printf(f, v...)
}

func (self *Router) Errorln(v ...interface{}) {
	self.err.Println(v...)
}

// Print out a "banner" to the logs
func (self *Router) Banner(text string) {
	self.out.Println("\n")
	self.out.Println("************************************************************")
	self.out.Printf("**  %s\n", text)
	self.out.Println("************************************************************")
	self.out.Println("\n")
}
