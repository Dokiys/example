package main

import "testing"

func Test_splitLine(t *testing.T) {
	str := "要坚持党的领导、坚持中国特⾊社会主义制度、贯彻中国特⾊社会主义法治理论。坚定不移⾛${中法道路}，最根本的是坚持${党}的领导；${中社制度}是${中法体系}的根本制度基础；${中法理论}是${中法体系}的理论指导和学理⽀撑；"
	splitLine(str)
}
