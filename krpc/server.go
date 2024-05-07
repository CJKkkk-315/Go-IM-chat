package KgoRpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)
import "GOchat/krpc/codec"

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber    int
	CodecType      codec.Type
	ConnectTimeout time.Duration
	HandleTimeout  time.Duration
}

var DefaultOption = &Option{
	MagicNumber:    MagicNumber,
	CodecType:      codec.GobType,
	ConnectTimeout: time.Second * 10,
}

type Server struct {
	serviceMap    sync.Map
	HandleTimeout time.Duration
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go s.ServeConn(conn)
	}
}

func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

func (s *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	opt := &Option{}
	if err := json.NewDecoder(conn).Decode(opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	f := codec.NewCodeFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	s.HandleTimeout = opt.HandleTimeout
	s.ServeCodec(f(conn))
}

var invalidRequest = struct{}{}

func (s *Server) ServeCodec(c codec.Codec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := s.readRequest(c)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			s.sendResponse(c, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go s.handleRequest(c, req, sending, wg)
	}
	wg.Wait()
	_ = c.Close()
}

type request struct {
	h            *codec.Header
	argv, replyv reflect.Value
	mtype        *methodType
	svc          *service
}

func (s *Server) readRequest(c codec.Codec) (*request, error) {
	req := &request{}
	var err error
	h := &codec.Header{}
	if err := c.ReadHeader(h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server read header error: ", err)
			return req, nil
		}
		return req, err
	}

	req.h = h

	req.svc, req.mtype, err = s.findService(req.h.ServiceMethod)
	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReplyv()

	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}
	if err = c.ReadBody(argvi); err != nil {
		log.Println("rpc server: read body err:", err)
		return req, err
	}
	return req, nil

}

func (s *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

func (s *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	called := make(chan struct{})
	sent := make(chan struct{})
	go func() {
		err := req.svc.call(req.mtype, req.argv, req.replyv)
		called <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			s.sendResponse(cc, req.h, invalidRequest, sending)
			sent <- struct{}{}
			return
		}
		s.sendResponse(cc, req.h, req.replyv.Interface(), sending)
		sent <- struct{}{}
	}()

	if s.HandleTimeout == 0 {
		<-called
		<-sent
		return
	}
	select {
	case <-time.After(s.HandleTimeout):
		req.h.Error = fmt.Sprintf("rpc server: request handle timeout: expect within %s", s.HandleTimeout)
		s.sendResponse(cc, req.h, invalidRequest, sending)
	case <-called:
		<-sent
	}
}

func (s *Server) Register(rcvr interface{}) error {
	service := newService(rcvr)
	if _, dup := s.serviceMap.LoadOrStore(service.name, service); dup {
		return errors.New("rpc: service already defined: " + service.name)
	}
	return nil
}

func Register(rcvr interface{}) error { return DefaultServer.Register(rcvr) }

func (s *Server) findService(serviceMethod string) (svc *service, mtype *methodType, err error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc server: service/method request ill-formed: " + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	svci, ok := s.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}
	svc = svci.(*service)
	mtype = svc.method[methodName]
	if mtype == nil {
		err = errors.New("rpc server: can't find method " + methodName)
	}
	return
}
