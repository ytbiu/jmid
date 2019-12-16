package jmid

import (
	"context"
	"io"
	l "log"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/log"
	"github.com/ytbiu/eventbus"
)

const CFG_CHANGED_TOPIC = "jaeger-cfg-changed"

var closer io.Closer

func Init(cfg *config.Configuration, watchChange bool) error {
	if watchChange {
		// 监听jaeger config变更事件
		eventbus.PubSub().Sub(context.TODO(), CFG_CHANGED_TOPIC, func(event interface{}) {
			newCfg, ok := event.(*config.Configuration)
			if !ok {
				l.Println("sub() event is not the type *config.Configuration")
				return
			}

			if closer != nil {
				closer.Close()
			}

			if err := Init(newCfg, false); err != nil {
				l.Println(err)
			}
		})

	}

	tracer, clo, err := cfg.NewTracer(
		config.Logger(log.StdLogger),
	)
	if err != nil {
		return err
	}

	closer = clo
	opentracing.SetGlobalTracer(tracer)
	return nil
}
