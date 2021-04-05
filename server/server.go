package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"smicro/logs"
	"smicro/middleware"
	"smicro/registry"
	_ "smicro/registry/etcd"
	"smicro/util"
)

type SmicroServer struct {
	*grpc.Server
	limiter        *rate.Limiter
	register       registry.Registry
	userMiddleware []middleware.Middleware
}

var smicroServer = &SmicroServer{
	Server: grpc.NewServer(),
}

func Use(m ...middleware.Middleware) {
	smicroServer.userMiddleware = append(smicroServer.userMiddleware, m...)
}

func Init(serviceName string) (err error) {
	err = InitConfig(serviceName)
	if err != nil {
		return
	}

	//初始化限流器
	if smicroConf.Limit.SwitchOn {
		smicroServer.limiter = rate.NewLimiter(rate.Limit(smicroConf.Limit.QPSLimit),
			smicroConf.Limit.QPSLimit)
	}

	initLogger()

	//初始化注册中心
	err = initRegister(serviceName)
	if err != nil {
		logs.Error(context.TODO(), "init register failed, err:%v", err)
		return
	}

	err = initTrace(serviceName)
	if err != nil {
		logs.Error(context.TODO(), "init tracing failed, err:%v", err)
	}
	return
}

func initTrace(serviceName string) (err error) {

	if !smicroConf.Trace.SwitchOn {
		return
	}

	return middleware.InitTrace(serviceName, smicroConf.Trace.ReportAddr,
		smicroConf.Trace.SampleType, smicroConf.Trace.SampleRate)
}

func initLogger() (err error) {
	filename := fmt.Sprintf("%s/%s.log", smicroConf.Log.Dir, smicroConf.ServiceName)
	outputer, err := logs.NewFileOutputer(filename)
	if err != nil {
		return
	}

	level := logs.GetLogLevel(smicroConf.Log.Level)
	logs.InitLogger(level, smicroConf.Log.ChanSize, smicroConf.ServiceName)
	logs.AddOutputer(outputer)

	if smicroConf.Log.ConsoleLog {
		logs.AddOutputer(logs.NewConsoleOutputer())
	}
	return
}

func initRegister(serviceName string) (err error) {

	if !smicroConf.Regiser.SwitchOn {
		return
	}

	ctx := context.TODO()
	registryInst, err := registry.InitRegistry(ctx,
		smicroConf.Regiser.RegisterName,
		registry.WithAddrs([]string{smicroConf.Regiser.RegisterAddr}),
		registry.WithTimeout(smicroConf.Regiser.Timeout),
		registry.WithRegistryPath(smicroConf.Regiser.RegisterPath),
		registry.WithHeartBeat(smicroConf.Regiser.HeartBeat),
	)
	if err != nil {
		logs.Error(ctx, "init registry failed, err:%v", err)
		return
	}

	smicroServer.register = registryInst
	service := &registry.Service{
		Name: serviceName,
	}

	ip, err := util.GetLocalIP()
	if err != nil {
		return
	}
	service.Nodes = append(service.Nodes, &registry.Node{
		IP:   ip,
		Port: smicroConf.Port,
	},
	)

	registryInst.Register(context.TODO(), service)
	return
}

func Run() {
	/*
		if smicroConf.Prometheus.SwitchOn {
			go func() {
				http.Handle("/metrics", promhttp.Handler())
				addr := fmt.Sprintf("0.0.0.0:%d", smicroConf.Prometheus.Port)
				log.Fatal(http.ListenAndServe(addr, nil))
			}()
		}*/

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", smicroConf.Port))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	smicroServer.Serve(lis)
}

func GRPCServer() *grpc.Server {
	return smicroServer.Server
}

func BuildServerMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var mids []middleware.Middleware

	mids = append(mids, middleware.AccessLogMiddleware)
	if smicroConf.Prometheus.SwitchOn {
		mids = append(mids, middleware.PrometheusServerMiddleware)
	}

	if smicroConf.Limit.SwitchOn {
		mids = append(mids, middleware.NewRateLimitMiddleware(smicroServer.limiter))
	}

	if smicroConf.Trace.SwitchOn {
		mids = append(mids, middleware.TraceServerMiddleware)
	}

	if len(smicroServer.userMiddleware) != 0 {
		mids = append(mids, smicroServer.userMiddleware...)
	}

	m := middleware.Chain(middleware.PrepareMiddleware, mids...)
	return m(handle)
}
