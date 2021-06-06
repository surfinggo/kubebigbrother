package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	spgc "github.com/spongeprojects/client-go/client/clientset/versioned"
	spgi "github.com/spongeprojects/client-go/client/informers/externalversions"
	spgl "github.com/spongeprojects/client-go/client/listers/spongeprojects.com/v1alpha1"
	"github.com/spongeprojects/kubebigbrother/pkg/gormdb"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

type Config struct {
	Env     string
	Version string

	Addr       string
	DBDialect  string
	DBArgs     string
	GinDebug   bool
	Kubeconfig string
}

type App struct {
	Version string

	Addr string
	Env  string

	EventStore             event_store.Interface
	ChannelInformer        cache.SharedIndexInformer
	WatcherInformer        cache.SharedIndexInformer
	ClusterWatcherInformer cache.SharedIndexInformer
	ChannelLister          spgl.ChannelLister
	WatcherLister          spgl.WatcherLister
	ClusterWatcherLister   spgl.ClusterWatcherLister
	Router                 *gin.Engine
}

func SetupApp(config *Config) (*App, error) {
	if config == nil {
		config = &Config{}
	}

	app := &App{}
	app.Version = config.Version
	app.Addr = config.Addr
	app.Env = config.Env

	if config.GinDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := gormdb.New(config.DBDialect, config.DBArgs)
	if err != nil {
		return nil, errors.Wrap(err, "create db instance error")
	}

	app.EventStore = event_store.New(db)

	restConfig, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "get kube config error")
	}

	spgClientset, err := spgc.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "new clientset error")
	}

	spgInformerFactory := spgi.NewSharedInformerFactory(spgClientset, 12*time.Hour)
	channelsInformer := spgInformerFactory.Spongeprojects().V1alpha1().Channels().Informer()
	watchersInformer := spgInformerFactory.Spongeprojects().V1alpha1().Watchers().Informer()
	clusterwatchersInformer := spgInformerFactory.Spongeprojects().V1alpha1().ClusterWatchers().Informer()
	channelsLister := spgInformerFactory.Spongeprojects().V1alpha1().Channels().Lister()
	watchersLister := spgInformerFactory.Spongeprojects().V1alpha1().Watchers().Lister()
	clusterwatchersLister := spgInformerFactory.Spongeprojects().V1alpha1().ClusterWatchers().Lister()

	app.ChannelInformer = channelsInformer
	app.WatcherInformer = watchersInformer
	app.ClusterWatcherInformer = clusterwatchersInformer
	app.ChannelLister = channelsLister
	app.WatcherLister = watchersLister
	app.ClusterWatcherLister = clusterwatchersLister

	r := gin.New()
	r.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/healthz"},
		}),
		gin.Recovery(),
	)

	r.GET("/", app.HandlerIndex)
	r.Any("/healthz", app.HandlerHealthz)
	r.GET("/api/v1/healthz", app.HandlerHealthz)
	r.POST("/api/v1/callback-channel-test", app.HandlerCallbackChannelTest)
	r.GET("/api/v1/config", app.HandlerConfig)
	r.GET("/api/v1/events", app.HandlerEventList)
	r.GET("/api/v1/events/:id", app.HandlerEvent)

	r.HandleMethodNotAllowed = true

	app.Router = r

	return app, nil
}

func (app *App) Run(stopCh chan struct{}) error {
	go app.ChannelInformer.Run(stopCh)
	go app.WatcherInformer.Run(stopCh)
	go app.ClusterWatcherInformer.Run(stopCh)

	cache.WaitForCacheSync(stopCh, app.ChannelInformer.HasSynced)
	cache.WaitForCacheSync(stopCh, app.WatcherInformer.HasSynced)
	cache.WaitForCacheSync(stopCh, app.ClusterWatcherInformer.HasSynced)

	return app.Router.Run(app.Addr)
}
