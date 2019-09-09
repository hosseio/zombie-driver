// +build wireinject

package bootstrap

import (
	netHttp "net/http"

	"github.com/heetch/jose-odg-technical-test/driver-location/pkg"

	"github.com/heetch/jose-odg-technical-test/driver-location/driver-location"
	"github.com/heetch/jose-odg-technical-test/driver-location/driver-location/internal"

	"github.com/nsqio/go-nsq"

	"github.com/chiguirez/cromberbus"

	"github.com/heetch/jose-odg-technical-test/driver-location/driver-location/messaging"

	"github.com/google/wire"
	"github.com/heetch/jose-odg-technical-test/driver-location/driver-location/cache"
	"github.com/heetch/jose-odg-technical-test/driver-location/driver-location/http"
)

var HttpSet = wire.NewSet(
	http.NewServer,
	http.NewRouter,
	http.NewLocationController,
	http.NewHealthController,
)

var CacheSet = wire.NewSet(
	cache.NewRedisDriver,
)

var AppSet = wire.NewSet(
	driver_location.NewLocationsByDriverAndTimeQueryService,
	driver_location.NewTransformer,
	wire.Bind(new(cromberbus.CommandHandlerResolver), cromberbus.MapHandlerResolver{}),
	wire.Bind(new(cromberbus.CommandBus), cromberbus.CromberBus{}),
	InitializeCromberbus,
	InitializeMapHandlerResolver,
	driver_location.NewCreateLocationCommandHandler,
	driver_location.NewDriverBuilder,
	wire.Bind(new(driver_location.LocationsByDriverAndTimeGetter), driver_location.LocationsByDriverAndTimeQueryService{}),
	wire.Bind(new(pkg.EventDispatcher), new(pkg.EventDispatcherMock)),
	newEventDispatcherMock,
)

func newEventDispatcherMock() *pkg.EventDispatcherMock {
	return &pkg.EventDispatcherMock{
		DispatchFunc: func(domainEvent []pkg.DomainEvent) {
			return
		},
	}
}

var MessagingSet = wire.NewSet(
	messaging.NewNsqConsumer,
	messaging.NewCreateDriverLocationHandler,
)

func InitializeServer(cfg Config) (*netHttp.Server, error) {
	wire.Build(
		HttpSet,
		AppSet,
		serverAddr,
		InitializeRedisDriver,
		wire.Bind(new(domain.LocationView), cache.RedisDriver{}),
	)

	return &netHttp.Server{}, nil
}

func serverAddr(cfg Config) http.ServerAddr {
	return http.ServerAddr(cfg.Server.Addr)
}

func InitializeCromberbus(handlerResolver cromberbus.CommandHandlerResolver) cromberbus.CromberBus {
	return cromberbus.NewCromberBus(handlerResolver)
}

func InitializeMapHandlerResolver(
	createLocationCommandHandler driver_location.CreateLocationCommandHandler,
) cromberbus.MapHandlerResolver {
	mapHandlerResolver := cromberbus.NewMapHandlerResolver()
	mapHandlerResolver.AddHandler(new(driver_location.CreateLocationCommand), createLocationCommandHandler)

	return mapHandlerResolver
}

func InitializeRedisDriver(cfg Config) (cache.RedisDriver, error) {
	wire.Build(
		CacheSet,
		AppSet,
		redisAddr,
	)

	return cache.RedisDriver{}, nil
}

func redisAddr(cfg Config) cache.RedisAddr {
	return cache.RedisAddr(cfg.Redis)
}

func InitializeCreateDriverLocationNsqConsumer(cfg Config) (messaging.NsqConsumer, error) {
	wire.Build(
		MessagingSet,
		AppSet,
		getNsqAddr,
		getTopicAddr,
		getChannelAddr,
		InitializeRedisDriver,
		wire.Bind(new(nsq.Handler), messaging.CreateDriverLocationHandler{}),
		wire.Bind(new(domain.DriverView), cache.RedisDriver{}),
		wire.Bind(new(domain.DriverRepository), cache.RedisDriver{}),
	)

	return messaging.NsqConsumer{}, nil
}

func getNsqAddr(cfg Config) messaging.NsqAddr {
	return messaging.NsqAddr(cfg.Nsq.Addr)
}

func getTopicAddr(cfg Config) messaging.TopicAddr {
	return messaging.TopicAddr(cfg.Nsq.Topic)
}

func getChannelAddr(cfg Config) messaging.ChannelAddr {
	return messaging.ChannelAddr(cfg.Nsq.Channel)
}
