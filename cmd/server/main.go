package main

import (
	"errors"
	"fmt"
	"go_grpc/pb"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const port = ":5001"

func main() {
	listen, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalln(err.Error())
	}
	//创建creds证书
	creds, err := credentials.NewServerTLSFromFile("server.pem", "server.key")
	if err != nil {
		log.Fatalln(err.Error())
	}
	//通过grpc传递creds证书
	options := []grpc.ServerOption{grpc.Creds(creds)}
	//创建server
	server := grpc.NewServer(options...)
	pb.RegisterEmployeeServiceServer(server, new(employeeService))

	log.Println("gRPC Server started ..." + port)
	//开启server监听listen端口号
	server.Serve(listen)

}

type employeeService struct{}

//GetByNo:通过员工编号找到员工
//一元消息传递
func (s *employeeService) GetByNo(ctx context.Context,
	req *pb.GetByNoRequest) (*pb.EmployeeResponse, error) {
	for _, e := range employees {
		if req.Number == e.Number {
			return &pb.EmployeeResponse{
				Employee: &e,
			}, nil
		}
	}

	return nil, errors.New("employee not found")
}

//二元消息传递
//服务端会将数据以streaming的形式传回
func (s *employeeService) GetAll(req *pb.GetAllRequest,
	stream pb.EmployeeService_GetAllServer) error {

	for _, e := range employees {
		//stream.send会将数据一块块的传给客户端
		stream.Send(&pb.EmployeeResponse{
			Employee: &e,
		})
		time.Sleep(time.Second)
	}
	return nil
}

//client 以stream的形式传输图片给服务端
func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())

	if ok {
		//通过metadata获取并输出employee的number
		fmt.Printf("Employee: %s\n", md["number"][0])
	}

	img := []byte{}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			//输出文件大小
			fmt.Printf("File Size: %d\n", len(img))
			//告诉客户端已经成功接收
			return stream.SendAndClose(&pb.AddPhotoResponse{
				IsOk: true,
			})
		}
		//输出每次接收的一小块的大小
		fmt.Printf("File received: %d\n", len(data.Data))
		img = append(img, data.Data...)
	}
}

func (s *employeeService) Save(context.Context,
	*pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	return nil, nil
}

//双向传送stream
func (s *employeeService) SaveAll(
	stream pb.EmployeeService_SaveAllServer) error {
	for {
		empReq, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		employees = append(employees, *empReq.Employee)
		stream.Send(&pb.EmployeeResponse{
			Employee: empReq.Employee,
		})
	}

	for _, emp := range employees {
		fmt.Println(emp)
	}

	return nil
}
