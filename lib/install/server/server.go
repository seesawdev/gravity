package server

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	installpb "github.com/gravitational/gravity/lib/install/proto"
	"github.com/gravitational/gravity/lib/ops"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// New returns a new instance of the installer server.
// Use Serve to make server start listening
func New(ctx context.Context) *Server {
	localCtx, cancel := context.WithCancel(ctx)
	grpcServer := grpc.NewServer()
	server := &Server{
		FieldLogger: log.WithField(trace.Component, "installer:service"),
		parentCtx:   ctx,
		ctx:         localCtx,
		cancel:      cancel,
		// operator:    config.Operator,
		rpc: grpcServer,
		// TODO(dmitri): arbitrary channel buffer size
		eventsC: make(chan Event, 100),
	}
	installpb.RegisterAgentServer(grpcServer, server)
	return server
}

// Serve starts the server using the specified engine
func (r *Server) Serve(executor Executor, listener net.Listener) error {
	r.executor = executor
	return trace.Wrap(r.rpc.Serve(listener))
}

// Stop stops the server gracefully
func (r *Server) Stop() {
	r.rpc.GracefulStop()
}

// Uninstall aborts the operation and cleans up the state
// Implements installpb.AgentServer
func (r *Server) Uninstall(ctx context.Context, req *installpb.UninstallRequest) (*installpb.UninstallResponse, error) {
	// TODO
	return &installpb.UninstallResponse{}, nil
}

// Shutdown shuts down the installer.
// Implements installpb.AgentServer
func (r *Server) Shutdown(ctx context.Context, req *installpb.ShutdownRequest) (*installpb.ShutdownResponse, error) {
	// The caller should be blocked at least as long as the wizard process is closing.
	// TODO(dmitri): find out how this returns to the caller and whether it would make sense
	// to split the shut down into several steps with wizard shutdown to be invoked as part of Shutdown
	// and the rest - from a goroutine so the caller is not receiving an error when the server stops
	// serving
	r.executor.Stop(ctx)
	return &installpb.ShutdownResponse{}, nil
}

// Execute executes the installation using the specified engine
// Implements installpb.AgentServer
func (r *Server) Execute(req *installpb.ExecuteRequest, stream installpb.Agent_ExecuteServer) error {
	r.executeOnce.Do(func() {
		r.serveWG.Add(1)
		go func() {
			if err := r.executor.Execute(); err != nil {
				r.WithError(err).Info("Failed to execute.")
				if err := r.sendError(err); err != nil {
					// TODO: only exit if unable to send the error.
					// Otherwise, the client will shut down the server at
					// the most appropriate time
				}
			}
			r.serveWG.Done()
			r.executor.Stop(r.parentCtx)
		}()
	})
	for {
		select {
		case event := <-r.eventsC:
			resp := &installpb.ProgressResponse{}
			if event.Progress != nil {
				resp.Message = event.Progress.Message
			} else if event.Error != nil {
				resp.Errors = append(resp.Errors, &installpb.Error{Message: event.Error.Error()})
			}
			err := stream.Send(resp)
			if err != nil {
				return trace.Wrap(err)
			}
			if resp.Complete {
				return nil
			}
		case <-stream.Context().Done():
			return trace.Wrap(stream.Context().Err())
		case <-r.parentCtx.Done():
			return trace.Wrap(r.parentCtx.Err())
		case <-r.ctx.Done():
			// Clean exit
			return nil
		}
	}
	return nil
}

// Run schedules f to run until the server is stopped
func (r *Server) Run(f func()) {
	r.serveWG.Add(1)
	go func() {
		f()
		r.serveWG.Done()
	}()
}

// Send streams the specified progress event to the client.
// The method will not block - event will be dropped if it cannot be published
// (subject to internal channel buffer capacity)
func (r *Server) Send(event Event) error {
	select {
	case r.eventsC <- event:
		// Pushed the progress event
		return nil
	case <-r.parentCtx.Done():
		return nil
	case <-r.ctx.Done():
		return nil
	default:
		r.WithField("event", event).Warn("Failed to publish event.")
		return trace.BadParameter("failed to publish event")
	}
}

// Executor wraps a potentially failing operation
type Executor interface {
	Execute() error
	// Stop signals the executor that it should abort the operation
	Stop(context.Context) error
}

// Server implements the installer gRPC server
type Server struct {
	log.FieldLogger
	// parentCtx specifies the external context.
	// If cancelled, all operations abort with the corresponding error
	parentCtx context.Context
	// ctx defines the local server context used to cancel internal operation
	ctx    context.Context
	cancel context.CancelFunc

	executor Executor
	eventsC  chan Event
	// rpc is the fabric to communicate to the server client prcess
	rpc *grpc.Server

	executeOnce sync.Once
	serveWG     sync.WaitGroup
}

// String formats this event for readability
func (r Event) String() string {
	var buf bytes.Buffer
	fmt.Print(&buf, "event(")
	if r.Progress != nil {
		fmt.Fprintf(&buf, "progress(completed=%v, message=%v),",
			r.Progress.Completion, r.Progress.Message)
	}
	if r.Error != nil {
		fmt.Fprintf(&buf, "error(%v)", r.Error.Error())
	}
	fmt.Print(&buf, ")")
	return buf.String()
}

// Event describes the installer progress step
type Event struct {
	// Progress describes the operation progress
	Progress *ops.ProgressEntry
	// Error specifies the error if any
	Error error
}

// RunProgressLoop starts progress loop for the specified operation
func (r *Server) RunProgressLoop(operator ops.Operator, operationKey ops.SiteOperationKey) {
	r.serveWG.Add(1)
	go func() {
		r.WithField("operation", operationKey.OperationID).Info("Start progress feedback loop.")
		ticker := time.NewTicker(1 * time.Second)
		defer func() {
			ticker.Stop()
			r.serveWG.Done()
		}()
		var lastProgress *ops.ProgressEntry
		for {
			select {
			case <-ticker.C:
				progress, err := operator.GetSiteOperationProgress(operationKey)
				if err != nil {
					r.WithError(err).Warn("Failed to query operation progress.")
					continue
				}
				if lastProgress != nil && lastProgress.IsEqual(*progress) {
					continue
				}
				r.Send(Event{Progress: progress})
				lastProgress = progress
				if progress.IsCompleted() {
					return
				}
			case <-r.parentCtx.Done():
				return
			case <-r.ctx.Done():
				return
			}
		}
	}()
}

func (r *Server) sendError(err error) error {
	return trace.Wrap(r.Send(Event{Error: err}))
}
