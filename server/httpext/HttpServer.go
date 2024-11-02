package httpext

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/server/route"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
)

const defaultMemory = 32 << 20

type Session struct {
	ID    string
	Token string
}

type httpServer struct {
	//serverDesc          *serverDesc
	engine             *mux.Router
	route              constrain.IRoute
	redisClient        constrain.IRedis
	etcdClient         constrain.IEtcd
	errorHandleResult  constrain.IResultError
	router             *mux.Router
	beforeViewRender   constrain.IBeforeViewRender
	viewRender         constrain.IAfterViewRender
	notFoundViewRender constrain.IAfterViewRender
	middlewares        []constrain.IMiddleware
	defaultMiddleware  *httpMiddleware
	errorHandler       ErrorHandler
	serviceList        []constrain.IService
}

func (m *httpServer) ApiErrorHandle(result constrain.IResultError) {
	m.errorHandleResult = result
}

/*
	func (m *httpServer) Use(middleware constrain.IMiddleware) {
		m.middlewares = append(m.middlewares, middleware)
	}
*/
func (m *httpServer) Listen(microServerConfig *config.MicroServerConfig, callbacks ...func() error) error {
	/*if err := microServerConfig.Register(); err != nil {
		return err
	}
	defer func() {
		err := microServerConfig.UnRegister()
		if err != nil {
			log.Fatalln(err)
		}
	}()*/
	m.defaultMiddleware.serverName = microServerConfig.MicroServer.Name
	if m.etcdClient != nil {
		var err error
		microServerConfig, err = m.etcdClient.Register(microServerConfig)
		if err != nil {
			return err
		}
	}

	listenAddr := fmt.Sprintf("%s:%d", microServerConfig.IP, microServerConfig.Port)
	log.Printf("HttpServer Listen:%s", listenAddr)
	srv := &http.Server{
		Handler:      m.engine,
		Addr:         listenAddr,
		WriteTimeout: 3 * time.Minute,
		ReadTimeout:  3 * time.Minute,
	}

	go func() {
		errg := &errgroup.Group{}
		if len(callbacks) > 0 {
			for i := range callbacks {
				errg.Go(func() error {
					return callbacks[i]()
				})
			}
		}
		err := errg.Wait()
		if err != nil {
			log.Println(err)
		}
	}()

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Println(fmt.Sprintf("tcp:server:%s:%s:START_ERROR:%s", m.defaultMiddleware.serverName, listenAddr, err.Error()))
		return err
	}
	go func() {
		_, err := net.DialTimeout("tcp", listenAddr, time.Second*5)
		if err != nil {
			log.Println(fmt.Sprintf("tcp:server:%s:%s:START_ERROR:%s", m.defaultMiddleware.serverName, listenAddr, err.Error()))
		} else {
			log.Println(fmt.Sprintf("tcp:server:%s:%s:START_SUCCESS", m.defaultMiddleware.serverName, listenAddr))
		}
	}()
	return srv.Serve(ln)
}

type Option interface {
	apply(server *httpServer)
}

/*type serverDesc struct {
	serverName string
	ip         string
	port       int
}

func (s *serverDesc) apply(server *httpServer) {
	server.serverDesc = s
}*/

/*func WithServerDesc(serverName string, ip string, port int) Option {
	return &serverDesc{serverName: serverName, ip: ip, port: port}
}*/

type ErrorHandler func(ctx constrain.IContext, customizeViewRender constrain.IAfterViewRender, w http.ResponseWriter, r *http.Request, err error)

type emptyOption struct {
	applyFunc func(server *httpServer)
}

func newOption(apply func(server *httpServer)) *emptyOption {
	return &emptyOption{applyFunc: apply}
}
func (e *emptyOption) apply(server *httpServer) {
	e.applyFunc(server)
}

/*
	func WithEtcdOption(etcdClient constrain.IEtcd) Option {
		return newOption(func(server *httpServer) {
			server.etcdClient = etcdClient
		})
	}

	func WithRedisOption(redisClient constrain.IRedis) Option {
		return newOption(func(server *httpServer) {
			server.redisClient = redisClient
		})
	}
*/
func WithErrorHandlerOption(errorHandler ErrorHandler) Option {
	return newOption(func(server *httpServer) {
		server.errorHandler = errorHandler
	})
}
func WithBeforeViewRenderOption(beforeViewRender constrain.IBeforeViewRender) Option {
	return newOption(func(server *httpServer) {
		server.beforeViewRender = beforeViewRender
	})
}
func WithViewRenderOption(customizeViewRender constrain.IAfterViewRender) Option {
	return newOption(func(server *httpServer) {
		server.viewRender = customizeViewRender
	})
}
func WithService(service constrain.IService) Option {
	return newOption(func(server *httpServer) {
		server.serviceList = append(server.serviceList, service)
	})
}
func WithNotFoundViewRenderOption(customizeViewRender constrain.IAfterViewRender) Option {
	return newOption(func(server *httpServer) {
		server.notFoundViewRender = customizeViewRender
	})
}

func (m *httpServer) handlerFunc(beforeViewRender constrain.IBeforeViewRender, viewRender constrain.IAfterViewRender, response http.ResponseWriter, request *http.Request) (next bool, ctxValue *contexext.ContextValue) {

	var ctx constrain.IContext
	ctxValue = contexext.FromContext(request.Context())
	if ctxValue != nil {
		ctx = request.Context().(constrain.IContext)
		ctxValue = contexext.FromContext(ctx)
	} else {
		ctx = m.defaultMiddleware.CreateContext(m.etcdClient, m.redisClient, m.route, response, request) //CreateContext(m.redisClient, m.etcdClient, m.route, response, request)
		ctxValue = contexext.FromContext(ctx)
		ctxValue.Request = request.WithContext(ctx)
	}

	defer func() {
		var err error
		if mErr := recover(); mErr != nil {
			switch mErr.(type) {
			case error:
				err = mErr.(error)
			default:
				err = fmt.Errorf("%v", mErr)
			}
			ctx.Logger().Error("http-server", zap.Error(err))
			m.errorHandler(ctx, viewRender, ctxValue.Response, ctxValue.Request, err)
		}
	}()

	for i := range m.middlewares {
		middleware := m.middlewares[i]
		if err := middleware.Handle(ctx, m.route, beforeViewRender, viewRender, ctxValue.Response, ctxValue.Request); err != nil {
			ctx.Logger().Error("http-server", zap.Error(err))
			m.errorHandler(ctx, viewRender, ctxValue.Response, ctxValue.Request, err)
			return
		}
	}
	/*if ctxValue.Request.Method == http.MethodOptions {
		return
	}*/
	if err := m.defaultMiddleware.Handle(ctx, m.route, beforeViewRender, viewRender, ctxValue.Response, ctxValue.Request); err != nil {
		ctx.Logger().Error("http-server", zap.Error(err))
		m.errorHandler(ctx, viewRender, ctxValue.Response, ctxValue.Request, err)
		return
	}
	return true, ctxValue
}

func NewHttpServer(etcdClient constrain.IEtcd, redisClient constrain.IRedis, engine *mux.Router, router *mux.Router, mRoute constrain.IRoute, ops ...Option) *httpServer {
	s := &httpServer{router: router, etcdClient: etcdClient, redisClient: redisClient, route: mRoute, engine: engine, defaultMiddleware: &httpMiddleware{}}

	engine.HandleFunc("/debug/pprof/*any", func(writer http.ResponseWriter, request *http.Request) {
		http.DefaultServeMux.ServeHTTP(writer, request)
	})

	s.errorHandler = func(ctx constrain.IContext, customizeViewRender constrain.IAfterViewRender, w http.ResponseWriter, r *http.Request, err error) {
		contextValue := contexext.FromContext(ctx)
		if contextValue.IsApi {
			var ar *result.ActionResult
			if errors.As(err, &ar) {
				w.Header().Set("Code", fmt.Sprintf("%d", ar.Code))
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			var e error
			bytes, e := json.Marshal(result.NewError(err))
			if e != nil {
				log.Println(e)
			}
			w.Write(bytes)
		} else {
			err = customizeViewRender.Render(ctx, r, w, route.NewViewResult("404", map[string]any{"Error": err.Error()}))
			if err != nil {
				ctx.Logger().With(zap.NamedError("ErrorHandler", err))
			}
		}
	}

	for i := range ops {
		ops[i].apply(s)
	}

	s.defaultMiddleware.serviceList = s.serviceList

	if s.viewRender == nil {
		s.viewRender = &DefaultViewRender{}
	}

	if s.router != nil && s.route != nil {
		router.NotFoundHandler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			/*if s.notFoundViewRender == nil {
				if request.URL.Path == "/404" {
					writer.WriteHeader(http.StatusNotFound)
					t, _ := template.New("404").Parse("404")
					err := t.Execute(writer, nil)
					if err != nil {
						log.Println(err)
					}
				} else {
					http.Redirect(writer, request, "/404", http.StatusPermanentRedirect)
				}
			} else {
				s.handlerFunc(s.beforeViewRender, s.notFoundViewRender, writer, request)
			}*/
			if s.notFoundViewRender == nil {
				/*if request.URL.Path == "/404" {
					writer.WriteHeader(http.StatusNotFound)
					t, _ := template.New("404").Parse("404")
					err := t.Execute(writer, nil)
					if err != nil {
						log.Println(err)
					}
				} else {
					http.Redirect(writer, request, "/404", http.StatusPermanentRedirect)
				}*/
				s.handlerFunc(s.beforeViewRender, s.viewRender, writer, request)
			} else {
				s.handlerFunc(s.beforeViewRender, s.notFoundViewRender, writer, request)
			}

		})
		//router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		/*ctx := DefaultHttpMiddleware.CreateContext(s.redisClient, mRoute, w, r)
		w.WriteHeader(http.StatusNotFound)
		viewResult := route.NewViewResult("404", map[string]any{})
		viewBaseValue := reflect.ValueOf(viewResult).Elem().FieldByName("ViewBase")
		viewBase := viewBaseValue.Interface().(extends.ViewBase)
		viewBaseValue.Set(reflect.ValueOf(viewBase))
		if s.customizeViewRender != nil {
			if err := s.customizeViewRender.Render(ctx, r, w, viewResult); err != nil {
				ctx.Logger().Error("render", zap.Error(err))
			}
			return
		}
		vr := &viewRender{}
		if err := vr.Render(ctx, r, w, viewResult); err != nil {
			ctx.Logger().Error("render", zap.Error(err))
		}*/
		//})

		s.router.Use(func(next http.Handler) http.Handler {
			if s.route == nil {
				log.Println("没有启用路由功能，因为httpServer.route(constrain.IRoute)对象为空")
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			}
			return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				hasNext, ctxValue := s.handlerFunc(s.beforeViewRender, s.viewRender, response, request)
				if hasNext {
					next.ServeHTTP(ctxValue.Response, ctxValue.Request)
				}
			})
		})
	}

	return s
}
