package main

import (
	"go_grpc/pb"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var employees = []pb.Employee{
	{
		Id:        210,
		Number:    211,
		FirstName: "xx",
		LastName:  "xx1",
		MonthSalary: &pb.MonthSalary{
			Basic: 200,
			Bonus: 125.5,
		},
		Status: pb.EmployeeStatus_NORMAL,
		LastModfied: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
	{
		Id:        310,
		Number:    311,
		FirstName: "asd",
		LastName:  "wefewf",
		MonthSalary: &pb.MonthSalary{
			Basic: 300,
			Bonus: 5.5,
		},
		Status: pb.EmployeeStatus_NORMAL,
		LastModfied: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
	{
		Id:        410,
		Number:    411,
		FirstName: "wwy",
		LastName:  "wyq",
		MonthSalary: &pb.MonthSalary{
			Basic: 4566,
			Bonus: 100,
		},
		Status: pb.EmployeeStatus_NORMAL,
		LastModfied: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
}
