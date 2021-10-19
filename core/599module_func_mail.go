package core

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"strings"
)

// 邮件发送
func (fns *InternalFunctionSet) Mailer() Value {
	msg := gomail.NewMessage()
	obj := &Mailer{msg: msg}
	return newClass("Mailer", &obj)
}

type Mailer struct {
	msg       *gomail.Message
	src       string
	mhost     string
	mport     int
	musername string
	mpassword string
}

func (m *Mailer) SetFrom(f string) {
	m.src = f
	m.msg.SetHeader("From", f)
}
func (m *Mailer) From(f string) {
	m.SetFrom(f)
}

func (m *Mailer) SetTo(t string) {
	m.msg.SetHeader("To", t)
}
func (m *Mailer) To(t string) {
	m.SetTo(t)
}

func (m *Mailer) SetCc(c string) {
	index := strings.Index(c, "@")
	nickname := m.src[:index]
	m.msg.SetAddressHeader("Cc", c, nickname)
}
func (m *Mailer) Cc(c string) {
	m.SetCc(c)
}

func (m *Mailer) SetSubject(sj string) {
	m.msg.SetHeader("Subject", sj)
}
func (m *Mailer) Subject(sj string) {
	m.SetSubject(sj)
}

func (m *Mailer) SetBody(body string) {
	m.msg.SetBody("text/html", body)
}
func (m *Mailer) Body(body string) {
	m.SetBody(body)
}

func (m *Mailer) SetAttach(path string) {
	m.msg.Attach(path)
}
func (m *Mailer) Attach(path string) {
	m.SetAttach(path)
}

func (m *Mailer) SetHost(mhost string) {
	m.mhost = mhost
}
func (m *Mailer) Host(mhost string) {
	m.SetHost(mhost)
}

func (m *Mailer) SetPort(mport int) {
	m.mport = mport
}
func (m *Mailer) Port(mport int) {
	m.SetPort(mport)
}

func (m *Mailer) SetUsername(musername string) {
	m.musername = musername
}
func (m *Mailer) Username(musername string) {
	m.SetUsername(musername)
}

func (m *Mailer) SetPassword(mpassword string) {
	m.mpassword = mpassword
}
func (m *Mailer) Password(mpassword string) {
	m.SetPassword(mpassword)
}

// 邮件发送
func (m *Mailer) Send() {
	host := m.parseHost()
	port := m.parsePort()
	fmt.Println("host", host)
	fmt.Println("port", port)

	d := gomail.NewDialer(host, port, m.musername, m.mpassword)

	if err := d.DialAndSend(m.msg); err != nil {
		runtimeExcption(err)
	}
}

var mailHostMap = map[string]string{
	"sohu.com": "smtp.sohu.com",
	"163.com":  "smtp.163.com",
}
var mailPortMap = map[string]int{
	"sohu.com": 25,
	"163.com":  465,
}

func (m *Mailer) parseHost() string {
	if m.mhost != "" {
		return m.mhost
	}
	index := strings.Index(m.src, "@")
	key := m.src[index+1:]
	host, ok := mailHostMap[key]
	if !ok {
		runtimeExcption("email host is not found for", m.src)
	}
	return host
}

func (m *Mailer) parsePort() int {
	if m.mport != 0 {
		return m.mport
	}

	index := strings.Index(m.src, "@")
	key := m.src[index+1:]
	port, ok := mailPortMap[key]
	if !ok {
		runtimeExcption("email port is not found for", m.src)
	}
	return port
}
