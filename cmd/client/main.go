package main

import (
	"context"
	"fmt"
	"go_grpc/pb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const port = ":5001"

func main() {
	creds, err := credentials.NewClientTLSFromFile("server.pem", "www.test.com")
	if err != nil {
		log.Fatal(err.Error())
	}
	//接着设置options
	options := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost"+port, options...)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()
	client := pb.NewEmployeeServiceClient(conn)
	fmt.Println("Client Server started...")
	getAll(client)
	// saveAll(client)
	// getByNo(client)
}

func getByNo(client pb.EmployeeServiceClient) {
	res, err := client.GetByNo(context.Background(), &pb.GetByNoRequest{Number: 211})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(res.Employee)
}

func saveAll(client pb.EmployeeServiceClient) {
	employees := []pb.Employee{
		{
			Id:        201,
			Number:    202,
			FirstName: "aa",
			LastName:  "x3",
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
			Id:        301,
			Number:    302,
			FirstName: "a2d",
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
			Id:        401,
			Number:    402,
			FirstName: "w2",
			LastName:  "w2wq",
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

	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	//我们不知道什么时候服务器会把数据发回，我们不能在这阻塞，采用goroutine
	finshChannel := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				finshChannel <- struct{}{}
				break
			}
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Println(res.Employee)
		}
	}()

	for _, e := range employees {
		err := stream.Send(&pb.EmployeeRequest{
			Employee: &e,
		})
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	stream.CloseSend()
	<-finshChannel
}

func addPhoto(client pb.EmployeeServiceClient) {
	imgFile, err := os.Open("WechatIMG3.jpeg")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer imgFile.Close()
	//metadata相当于报文的header，我们只需要把用户number放在header传输一次就可以了
	md := metadata.New(map[string]string{"number": "1994"})
	context := context.Background()
	context = metadata.NewOutgoingContext(context, md)

	stream, err := client.AddPhoto(context)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err.Error())
	}
	//循环分块传输数据
	for {
		chunk := make([]byte, 128*1024)
		chunkSize, err := imgFile.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err.Error())
		}
		if chunkSize < len(chunk) {
			chunk = chunk[:chunkSize]
		}
		//开始分块发送数据
		stream.Send(&pb.AddPhotoRequest{Data: chunk})

	}
	//closeandrec会向客户端发送一个信号EOF，等待服务端发回一个响应
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(res.IsOk)
}

func getAll(client pb.EmployeeServiceClient) {
	stream, err := client.GetAll(context.Background(), &pb.GetAllRequest{})
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		res, err := stream.Recv()
		//如果服务端数据发送结束，则为EOF
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(res.Employee)
	}
}
