package main

import "testing"

func Test_splitLine(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{"要坚持党的领导、坚持中国特⾊社会主义制度、贯彻中国特⾊社会主义法治理论。坚定不移⾛${中法道路}，最根本的是坚持${党}的领导；${中社制度}是${中法体系}的根本制度基础；${中法理论}是${中法体系}的理论指导和学理⽀撑；"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitLine(tt.args.str)
		})
	}
}
