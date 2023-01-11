package rpc

import (
	"context"
	"reflect"
	"sync"
)

type handler struct {
	reg           *serviceRegistry
	unsubscribeCb *callback
	idgen         func() ID
	// respWait       map[string]*requestOp          // active client requests
	// clientSubs     map[string]*ClientSubscription // active client subscriptions
	callWG         sync.WaitGroup  // pending call goroutines
	rootCtx        context.Context // canceled by close()
	cancelRoot     func()
	conn           *httpConn
	allowSubscribe bool

	subLock    sync.Mutex
	serverSubs map[ID]*Subscription
}

func newHandler(connCtx context.Context, conn *httpConn, idgen func() ID, reg *serviceRegistry) *handler {
	rootCtx, cancelRoot := context.WithCancel(connCtx)
	h := &handler{
		reg:   reg,
		idgen: idgen,
		conn:  conn,
		// respWait:       make(map[string]*requestOp),
		// clientSubs:     make(map[string]*ClientSubscription),
		rootCtx:        rootCtx,
		cancelRoot:     cancelRoot,
		allowSubscribe: true,
		serverSubs:     make(map[ID]*Subscription),
	}

	return h
}

func (h *handler) handleBatch(msgs []*jsonrpcMessage) {
	answers := make([]*jsonrpcMessage, 0, len(msgs))
	for _, msg := range msgs {
		callb := h.reg.callback(msg.Method)
		if callb == nil {
			answers = append(answers, msg.errorResponse(&methodNotFoundError{method: msg.Method}))
			continue
		}

		args, err := parsePositionalArguments(msg.Params, callb.argTypes)
		if err != nil {
			answers = append(answers, msg.errorResponse(&invalidParamsError{err.Error()}))
			continue
		}
		answers = append(answers, h.runMethod(context.Background(), msg, callb, args))
	}

	if len(answers) > 0 {
		h.conn.writeJSON(context.Background(), answers)
	}
}

func (h *handler) handleMsg(msg *jsonrpcMessage) {
	callb := h.reg.callback(msg.Method)
	if callb == nil {
		h.conn.writeJSON(context.Background(), msg.errorResponse(&methodNotFoundError{method: msg.Method}))
		return
	}

	args, err := parsePositionalArguments(msg.Params, callb.argTypes)
	if err != nil {
		h.conn.writeJSON(context.Background(), msg.errorResponse(&invalidParamsError{err.Error()}))
		return
	}

	answer := h.runMethod(context.Background(), msg, callb, args)

	h.conn.writeJSON(context.Background(), answer)
}

func (h *handler) runMethod(ctx context.Context, msg *jsonrpcMessage, callb *callback, args []reflect.Value) *jsonrpcMessage {
	result, err := callb.call(ctx, msg.Method, args)
	if err != nil {
		return msg.errorResponse(err)
	}
	return msg.response(result)
}

type callProc struct {
	ctx context.Context
	// notifiers []*Notifier
}
